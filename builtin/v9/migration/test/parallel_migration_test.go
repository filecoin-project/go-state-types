package migration

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-state-types/builtin/v9/util/adt"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v9/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestParallelMigrationCalls(t *testing.T) {
	// Construct simple prior state tree over a synchronized store
	ctx := context.Background()
	log := TestLogger{TB: t}
	bs := NewSyncBlockStoreInMemory()

	// Run migration
	adtStore := adt.WrapStore(ctx, cbor.NewCborStore(bs))
	startRoot := makeInputTree(ctx, t, adtStore)
	newManifestCid, _ := makeTestManifest(t, adtStore, "fil/9/")
	endRootSerial, err := migration.MigrateStateTree(ctx, adtStore, newManifestCid, startRoot, abi.ChainEpoch(0), migration.Config{MaxWorkers: 1}, log, migration.NewMemMigrationCache())
	require.NoError(t, err)

	// Migrate in parallel
	var endRootParallel1, endRootParallel2 cid.Cid
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		var err1 error
		endRootParallel1, err1 = migration.MigrateStateTree(ctx, adtStore, newManifestCid, startRoot, abi.ChainEpoch(0), migration.Config{MaxWorkers: 2}, log, migration.NewMemMigrationCache())
		return err1
	})
	grp.Go(func() error {
		var err2 error
		endRootParallel2, err2 = migration.MigrateStateTree(ctx, adtStore, newManifestCid, startRoot, abi.ChainEpoch(0), migration.Config{MaxWorkers: 2}, log, migration.NewMemMigrationCache())
		return err2
	})
	require.NoError(t, grp.Wait())
	require.Equal(t, endRootSerial, endRootParallel1)
	require.Equal(t, endRootParallel1, endRootParallel2)
}
