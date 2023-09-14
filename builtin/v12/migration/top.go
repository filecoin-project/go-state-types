package migration

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	market11 "github.com/filecoin-project/go-state-types/builtin/v11/market"
	system11 "github.com/filecoin-project/go-state-types/builtin/v11/system"
	adt11 "github.com/filecoin-project/go-state-types/builtin/v11/util/adt"
	"github.com/filecoin-project/go-state-types/builtin/v12/market"
	market12 "github.com/filecoin-project/go-state-types/builtin/v12/market"
	adt12 "github.com/filecoin-project/go-state-types/builtin/v12/util/adt"
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
			continue
		}
		if oldEntry.Name == manifest.MarketKey {
			market11Cid = oldEntry.Code
			deferredCodeIDs[market11Cid] = struct{}{}
			continue
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

	// The Market Actor
	market12Cid, ok := newManifest.Get(manifest.MarketKey)
	if !ok {
		return cid.Undef, xerrors.Errorf("code cid for market actor not found in new manifest")
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

	actorsOut, err := migration.RunMigration(ctx, cfg, cache, store, log, actorsIn, migrations, deferredCodeIDs)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to run migration: %w", err)
	}

	// Retrieve the sector deal IDs from the minerMigrator. These IDs are crucial for the migration of market actor state.
	sectorDealIDs, err := minerMigrator.sectorDeals.Root()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to flush actorsOut: %w", err)
	}

	// Check if sectorDealIDs is not a zero value before proceeding.
	if !sectorDealIDs.Defined() {
		return cid.Undef, xerrors.New("sectorDealIDs is a zero value, cannot proceed with migration")
	}

	//save sectorDealIdIndex in Cache
	if err = cache.Write(migration.MarketSectorIndexKey(), sectorDealIDs); err != nil {
		return cid.Undef, xerrors.Errorf("failed to write inkey to cache: %w", err)
	}

	var prevDealStates *adt12.Array
	prevDealStatesOk, prevDealStatesCid, err := cache.Read(migration.PrevMarketStatesAmtKey())
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to read prev market states amt key from cache: %w", err)
	}
	if prevDealStatesOk {
		prevDealStates, err = adt12.AsArray(adtStore, prevDealStatesCid, market.StatesAmtBitwidth)
	} else {
		prevDealStates, err = adt12.MakeEmptyArray(adtStore, market.StatesAmtBitwidth)
	}

	//migrate market.States
	//todo confirm bitwidth
	newDealStates, err := adt12.MakeEmptyArray(adtStore, market.StatesAmtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to create new adt empty array newDealStates: %w", err)
	}
	if dealStates, err := adt11.AsArray(adtStore, oldMarketState.States, market11.StatesAmtBitwidth); err != nil {
		return cid.Undef, xerrors.Errorf("failed to read oldMarketState.States: %w", err)
	} else {
		var oldDealState market11.DealState
		err = dealStates.ForEach(&oldDealState, func(dealID int64) error {
			var prevRunDealState market12.DealState
			var sectorNumber uint64
			found, err := prevDealStates.Get(uint64(dealID), &prevRunDealState)
			if err != nil {
				return xerrors.Errorf("failed to load sector number from previous deal states AMT: %d", dealID)

			}
			if found {
				// we have found this deal id in our previous deal state amt, so extract the sector Number from there
				sectorNumber = uint64(prevRunDealState.SectorNumber)
			} else {
				// Deal ID not found in previous run array, so it must be in dealToSectorIndex or error
				sectorNumberAny, ok := (*minerMigrator.dealToSectorIndex).Load(uint64(dealID))
				if !ok {
					return xerrors.Errorf("failed to load sector number for deal ID: %d", dealID)
				}
				sectorNumber, ok = sectorNumberAny.(uint64)
				if !ok {
					return xerrors.Errorf("failed to assert sectorNumberUint64 to uint64 for deal ID: %d", dealID)
				}
			}

			newDealState := market12.DealState{
				SectorNumber:     abi.SectorNumber(sectorNumber),
				SectorStartEpoch: oldDealState.SectorStartEpoch,
				LastUpdatedEpoch: oldDealState.LastUpdatedEpoch,
				SlashEpoch:       oldDealState.SlashEpoch,
			}
			newDealStates.Set(uint64(dealID), &newDealState)
			_ = newDealStates
			return nil
		})
	}
	newDealStatesCid, err := newDealStates.Root()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to get root of newDealStatesCid: %w", err)
	}

	if err = cache.Write(migration.PrevMarketStatesAmtKey(), newDealStatesCid); err != nil {
		return cid.Undef, xerrors.Errorf("failed to write prev market states amt key to cache: %w", err)
	}

	// Create the new state
	newMarketState := market12.State{
		Proposals:                     oldMarketState.Proposals,
		States:                        newDealStatesCid,
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
	fmt.Println("sector index: ", sectorDealIDs)

	newMarketStateCid, err := store.Put(ctx, &newMarketState)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to put new state: %w", err)
	}

	// Update the market actor in the state tree with the newly fetched sector deal IDs.
	// This ensures the market actor's state reflects the most recent sector deals.
	if err = actorsOut.SetActorV5(builtin.StorageMarketActorAddr, &builtin.ActorV5{
		Code:       market12Cid,
		Head:       newMarketStateCid, // Updated value
		CallSeqNum: oldMarketActor.CallSeqNum,
		Balance:    oldMarketActor.Balance,
		Address:    oldMarketActor.Address,
	}); err != nil {
		return cid.Undef, xerrors.Errorf("failed to set market actor: %w", err)
	}

	return actorsOut.Flush()
}
