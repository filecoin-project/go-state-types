package migration

import (
	"context"

	adt15 "github.com/filecoin-project/go-state-types/builtin/v15/util/adt"

	system14 "github.com/filecoin-project/go-state-types/builtin/v14/system"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/filecoin-project/go-state-types/migration"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"
)

// MigrateStateTree migrates the Filecoin state tree starting from the global state tree and upgrading all actor states.
// The store must support concurrent writes (even if the configured worker count is 1).
//
// FIP-0081 constants for the power actor state for pledge calculations that apply only to this migration:
//
// - powerRampStartEpoch: Epoch at which the new pledge calculation starts.
// - powerRampDurationEpochs: Number of epochs over which the new pledge calculation is ramped up.
func MigrateStateTree(
	ctx context.Context,
	store cbor.IpldStore,
	newManifestCID cid.Cid,
	actorsRootIn cid.Cid,
	priorEpoch abi.ChainEpoch,
	powerRampStartEpoch int64,
	powerRampDurationEpochs uint64,
	cfg migration.Config,
	log migration.Logger,
	cache migration.MigrationCache,
) (cid.Cid, error) {
	if cfg.MaxWorkers <= 0 {
		return cid.Undef, xerrors.Errorf("invalid migration config with %d workers", cfg.MaxWorkers)
	}

	if powerRampStartEpoch == 0 {
		return cid.Undef, xerrors.Errorf("powerRampStartEpoch must be non-zero")
	}
	if powerRampStartEpoch < 0 {
		return cid.Undef, xerrors.Errorf("powerRampStartEpoch must be non-negative")
	}
	if powerRampDurationEpochs == 0 {
		return cid.Undef, xerrors.Errorf("powerRampDurationEpochs must be non-zero")
	}

	adtStore := adt15.WrapStore(ctx, store)

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

	var systemState system14.State
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

	power14Cid := cid.Undef

	for _, oldEntry := range oldManifestData.Entries {
		if oldEntry.Name == manifest.PowerKey {
			power14Cid = oldEntry.Code
		}

		newCodeCID, ok := newManifest.Get(oldEntry.Name)
		if !ok {
			return cid.Undef, xerrors.Errorf("code cid for %s actor not found in new manifest", oldEntry.Name)
		}
		migrations[oldEntry.Code] = migration.CodeMigrator{OutCodeCID: newCodeCID}
	}

	// migrations that migrate both code and state, override entries in `migrations`

	// The System Actor

	newSystemCodeCID, ok := newManifest.Get(manifest.SystemKey)
	if !ok {
		return cid.Undef, xerrors.Errorf("code cid for system actor not found in new manifest")
	}

	migrations[systemActor.Code] = systemActorMigrator{OutCodeCID: newSystemCodeCID, ManifestData: newManifest.Data}

	// The Power Actor
	if power14Cid == cid.Undef {
		return cid.Undef, xerrors.Errorf("code cid for power actor not found in old manifest")
	}
	power15Cid, ok := newManifest.Get(manifest.PowerKey)
	if !ok {
		return cid.Undef, xerrors.Errorf("code cid for power actor not found in new manifest")
	}

	pm, err := newPowerMigrator(powerRampStartEpoch, powerRampDurationEpochs, power15Cid)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to create miner migrator: %w", err)
	}
	migrations[power14Cid] = *pm

	if len(migrations)+len(deferredCodeIDs) != len(oldManifestData.Entries) {
		return cid.Undef, xerrors.Errorf("incomplete migration specification with %d code CIDs, need %d", len(migrations)+len(deferredCodeIDs), len(oldManifestData.Entries))
	}

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
