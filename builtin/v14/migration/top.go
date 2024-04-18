package migration

import (
	"context"

	adt14 "github.com/filecoin-project/go-state-types/builtin/v14/util/adt"

	system13 "github.com/filecoin-project/go-state-types/builtin/v13/system"
	account14 "github.com/filecoin-project/go-state-types/builtin/v14/account"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/filecoin-project/go-state-types/migration"

	"github.com/filecoin-project/go-address"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"
)

// MigrateStateTree Migrates the filecoin state tree starting from the global state tree and upgrading all actor state.
// The store must support concurrent writes (even if the configured worker count is 1).
func MigrateStateTree(ctx context.Context, store cbor.IpldStore, newManifestCID cid.Cid, actorsRootIn cid.Cid, priorEpoch abi.ChainEpoch, cfg migration.Config, log migration.Logger, cache migration.MigrationCache) (cid.Cid, error) {
	if cfg.MaxWorkers <= 0 {
		return cid.Undef, xerrors.Errorf("invalid migration config with %d workers", cfg.MaxWorkers)
	}

	adtStore := adt14.WrapStore(ctx, store)

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

	var systemState system13.State
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

	for _, oldEntry := range oldManifestData.Entries {
		newCodeCID, ok := newManifest.Get(oldEntry.Name)
		if !ok {
			return cid.Undef, xerrors.Errorf("code cid for %s actor not found in new manifest", oldEntry.Name)
		}
		migrations[oldEntry.Code] = migration.CachedMigration(cache, migration.CodeMigrator{OutCodeCID: newCodeCID})
	}

	// migrations that migrate both code and state, override entries in `migrations`

	// The System Actor
	newSystemCodeCID, ok := newManifest.Get(manifest.SystemKey)
	if !ok {
		return cid.Undef, xerrors.Errorf("code cid for system actor not found in new manifest")
	}

	migrations[systemActor.Code] = systemActorMigrator{OutCodeCID: newSystemCodeCID, ManifestData: newManifest.Data}

	if len(migrations)+len(deferredCodeIDs) != len(oldManifestData.Entries) {
		return cid.Undef, xerrors.Errorf("incomplete migration specification with %d code CIDs, need %d", len(migrations)+len(deferredCodeIDs), len(oldManifestData.Entries))
	}

	actorsOut, err := migration.RunMigration(ctx, cfg, cache, store, log, actorsIn, migrations)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to run migration: %w", err)
	}

	// FIP 0085.  f090 multisig => unkeyed account actor
	// Account state now points to id address requiring upgrade to release funds
	// TODO move to small file
	// 1. return error and new f090
	// 2. pass in actors in in order to get f090 to construct new one
	// 3. pass in actors out so we can flush it (maybe flush outside of fn? but I think we match other migration patterns by flushing in function? )
	f090ID, err := address.NewIDAddress(90)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to construct f090 id addr", err)
	}
	f090OldAct, found, err := actorsOut.GetActorV5(f090ID)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to get old f090 actor", err)
	}
	if !found {
		return cid.Undef, xerrors.Errorf("failed to find old f090 actor")
	}
	f090NewSt := account14.State{Address: f090ID} // State points to ID addr
	h, err := store.Put(ctx, f090NewSt)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to write new f090 state")
	}

	newAccountCodeCID, ok := newManifest.Get(manifest.AccountKey)
	if !ok {
		return cid.Undef, xerrors.Errorf("invalid manifest missing account code cid")
	}

	actorsOut.SetActorV5(f090ID, &builtin.ActorV5{
		// unchanged
		CallSeqNum: f090OldAct.CallSeqNum,
		Balance:    f090OldAct.Balance,
		Address:    f090OldAct.Address,

		// changed
		Code: newAccountCodeCID,
		Head: h,
	})

	outCid, err := actorsOut.Flush()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to flush actorsOut: %w", err)
	}

	return outCid, nil
}
