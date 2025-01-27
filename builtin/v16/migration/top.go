package migration

import (
	"context"

	"github.com/filecoin-project/go-state-types/builtin/v16/util/adt"

	system15 "github.com/filecoin-project/go-state-types/builtin/v13/system"
	miner16 "github.com/filecoin-project/go-state-types/builtin/v16/miner"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/filecoin-project/go-state-types/migration"

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

	adtStore := adt.WrapStore(ctx, store)

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

	var systemState system15.State
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

	miner15Cid := cid.Undef

	for _, oldEntry := range oldManifestData.Entries {
		newCodeCID, ok := newManifest.Get(oldEntry.Name)
		if !ok {
			return cid.Undef, xerrors.Errorf("code cid for %s actor not found in new manifest", oldEntry.Name)
		}
		if oldEntry.Name == manifest.MinerKey {
			miner15Cid = oldEntry.Code
		}
		migrations[oldEntry.Code] = migration.CachedMigration(cache, migration.CodeMigrator{OutCodeCID: newCodeCID})
	}

	if miner15Cid == cid.Undef {
		return cid.Undef, xerrors.Errorf("could not find miner actor in old manifest")
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

	var emptyPreCommittedSectorsHamtCid cid.Cid
	if hamt, err := adt.MakeEmptyMap(adtStore, builtin.DefaultHamtBitwidth); err != nil {
		return cid.Undef, xerrors.Errorf("failed to create empty precommit clean up amount array: %w", err)
	} else {
		if emptyPreCommittedSectorsHamtCid, err = hamt.Root(); err != nil {
			return cid.Undef, xerrors.Errorf("failed to get root of empty precommit clean up amount array: %w", err)
		}
	}
	var emptyPrecommitCleanUpAmtCid cid.Cid
	if amt, err := adt.MakeEmptyArray(adtStore, miner16.PrecommitCleanUpAmtBitwidth); err != nil {
		return cid.Undef, xerrors.Errorf("failed to create empty precommit clean up amount array: %w", err)
	} else {
		if emptyPrecommitCleanUpAmtCid, err = amt.Root(); err != nil {
			return cid.Undef, xerrors.Errorf("failed to get root of empty precommit clean up amount array: %w", err)
		}
	}

	miner16Cid, ok := newManifest.Get(manifest.MinerKey)
	if !ok {
		return cid.Undef, xerrors.Errorf("code cid for miner actor not found in new manifest")
	}

	minerMig, err := newMinerMigrator(ctx, store, miner16Cid, emptyPreCommittedSectorsHamtCid, emptyPrecommitCleanUpAmtCid)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to create miner migrator: %w", err)
	}

	migrations[miner15Cid] = migration.CachedMigration(cache, minerMig)

	actorsOut, err := migration.RunMigration(ctx, cfg, cache, store, log, actorsIn, migrations)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to run migration: %w", err)
	}

	outCid, err := actorsOut.Flush()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to flush actorsOut: %w", err)
	}

	return outCid, nil
}
