package migration

import (
	"context"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	market11 "github.com/filecoin-project/go-state-types/builtin/v11/market"
	system11 "github.com/filecoin-project/go-state-types/builtin/v11/system"
	adt11 "github.com/filecoin-project/go-state-types/builtin/v11/util/adt"
	market12 "github.com/filecoin-project/go-state-types/builtin/v12/market"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/filecoin-project/go-state-types/migration"
)

// MigrateStateTree Migrates the filecoin state tree starting from the global state tree and upgrading all actor state.
// The store must support concurrent writes (even if the configured worker count is 1).
func MigrateStateTree(ctx context.Context, store cbor.IpldStore, newManifestCID cid.Cid, actorsRootIn cid.Cid, priorEpoch abi.ChainEpoch, cfg migration.Config, log migration.Logger, cache migration.MigrationCache) (cid.Cid, error) {
	if cfg.MaxWorkers <= 0 {
		return cid.Undef, xerrors.Errorf("invalid migration config with %d workers", cfg.MaxWorkers)
	}

	adtStore := adt11.WrapStore(ctx, store)

	// Load input and output state trees
	actorsIn, err := builtin.LoadTree(adtStore, actorsRootIn)
	if err != nil {
		return cid.Undef, xerrors.Errorf("loading state tree: %w", err)
	}

	// load old manifest data
	systemActor, ok, err := actorsIn.GetActorV5(builtin.SystemActorAddr)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to get system actor: %w", err)
	}

	if !ok {
		return cid.Undef, xerrors.New("didn't find system actor")
	}

	var systemState system11.State
	if err := store.Get(ctx, systemActor.Head, &systemState); err != nil {
		return cid.Undef, xerrors.Errorf("failed to get system actor state: %w", err)
	}

	var oldManifestData manifest.ManifestData
	if err := store.Get(ctx, systemState.BuiltinActors, &oldManifestData); err != nil {
		return cid.Undef, xerrors.Errorf("failed to get old manifest data: %w", err)
	}

	// load new manifest
	var newManifest manifest.Manifest
	if err := adtStore.Get(ctx, newManifestCID, &newManifest); err != nil {
		return cid.Undef, xerrors.Errorf("error reading actor manifest: %w", err)
	}

	if err := newManifest.Load(ctx, adtStore); err != nil {
		return cid.Undef, xerrors.Errorf("error loading actor manifest: %w", err)
	}

	// Maps prior version code CIDs to migration functions.
	migrations := make(map[cid.Cid]migration.ActorMigration)
	// Set of prior version code CIDs for actors to defer during iteration, for explicit migration afterwards.
	deferredCodeIDs := make(map[cid.Cid]struct{})

	miner11Cid := cid.Undef
	market11Cid := cid.Undef

	for _, oldEntry := range oldManifestData.Entries {
		if oldEntry.Name == manifest.MinerKey {
			miner11Cid = oldEntry.Code
		}
		if oldEntry.Name == manifest.MarketKey {
			market11Cid = oldEntry.Code
		}
		newCodeCID, ok := newManifest.Get(oldEntry.Name)
		if !ok {
			return cid.Undef, xerrors.Errorf("code cid for %s actor not found in new manifest", oldEntry.Name)
		}
		migrations[oldEntry.Code] = migration.CachedMigration(cache, migration.CodeMigrator{OutCodeCID: newCodeCID})
	}

	if !miner11Cid.Defined() {
		return cid.Undef, xerrors.Errorf("didn't find miner actor in old manifest")
	}

	if !market11Cid.Defined() {
		return cid.Undef, xerrors.Errorf("didn't find market actor in old manifest")
	}

	// migrations that migrate both code and state, override entries in `migrations`

	// The System Actor

	newSystemCodeCID, ok := newManifest.Get(manifest.SystemKey)
	if !ok {
		return cid.Undef, xerrors.Errorf("code cid for system actor not found in new manifest")
	}

	migrations[systemActor.Code] = systemActorMigrator{OutCodeCID: newSystemCodeCID, ManifestData: newManifest.Data}

	// The Miner Actor
	miner12Cid, ok := newManifest.Get(manifest.MinerKey)
	if !ok {
		return cid.Undef, xerrors.Errorf("code cid for miner actor not found in new manifest")
	}

	minerMigrator, err := newMinerMigrator(ctx, store, miner12Cid, cache)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to create miner migrator: %w", err)
	}
	migrations[miner11Cid] = migration.CachedMigration(cache, *minerMigrator)

	if len(migrations)+len(deferredCodeIDs) != len(oldManifestData.Entries) {
		return cid.Undef, xerrors.Errorf("incomplete migration specification with %d code CIDs, need %d", len(migrations), len(oldManifestData.Entries))
	}

	// Load the state of the market actor from v11 for migration purposes.
	oldMarketActor, ok, err := actorsIn.GetActorV5(builtin.StorageMarketActorAddr)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to get market actor: %w", err)
	}

	if !ok {
		return cid.Undef, xerrors.New("didn't find market actor")
	}

	var oldMarketState market11.State
	if err := store.Get(ctx, oldMarketActor.Head, &oldMarketState); err != nil {
		return cid.Undef, xerrors.Errorf("failed to get market actor state: %w", err)
	}

	actorsOut, err := migration.RunMigration(ctx, cfg, cache, store, log, actorsIn, migrations)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to run migration: %w", err)
	}

	// Retrieve the sector deal IDs from the minerMigrator. These IDs are crucial for the migration of market actor state.
	sectorDealIDs, err := minerMigrator.sectorDeals.Map.Root()
	if err != nil {
		return cid.Undef, err
	}

	// Check if sectorDealIDs is not a zero value before proceeding.
	if sectorDealIDs == cid.Undef {
		return cid.Undef, xerrors.New("sectorDealIDs is a zero value, cannot proceed with migration")
	}

	//save sectorDealIdIndex in Cache
	if err = cache.Write(migration.SectorIndexHamtKey(), sectorDealIDs); err != nil {
		return cid.Undef, xerrors.Errorf("failed to write inkey to cache: %w", err)
	}

	// Create the new state
	newMarketState := market12.State{
		Proposals:                     oldMarketState.Proposals,
		States:                        oldMarketState.States, // FIXME migrate deal states to include sector number
		PendingProposals:              oldMarketState.PendingProposals,
		EscrowTable:                   oldMarketState.EscrowTable,
		LockedTable:                   oldMarketState.LockedTable,
		NextID:                        oldMarketState.NextID,
		DealOpsByEpoch:                oldMarketState.DealOpsByEpoch,
		LastCron:                      oldMarketState.LastCron,
		TotalClientLockedCollateral:   oldMarketState.TotalClientLockedCollateral,
		TotalProviderLockedCollateral: oldMarketState.TotalProviderLockedCollateral,
		TotalClientStorageFee:         oldMarketState.TotalClientStorageFee,
		PendingDealAllocationIds:      oldMarketState.PendingDealAllocationIds,
		ProviderSectors:               sectorDealIDs, // Updated value
	}

	newMarketStateCid, err := store.Put(ctx, &newMarketState)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to put new state: %w", err)
	}

	// Update the market actor in the state tree with the newly fetched sector deal IDs.
	// This ensures the market actor's state reflects the most recent sector deals.
	if err = actorsOut.SetActorV5(builtin.StorageMarketActorAddr, &builtin.ActorV5{
		Code:       market11Cid,
		Head:       newMarketStateCid, // Updated value
		CallSeqNum: oldMarketActor.CallSeqNum,
		Balance:    oldMarketActor.Balance,
		Address:    oldMarketActor.Address,
	}); err != nil {
		return cid.Undef, xerrors.Errorf("failed to set market actor: %w", err)
	}

	return actorsOut.Flush()
}
