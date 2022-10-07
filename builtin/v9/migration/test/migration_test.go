package migration

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-state-types/builtin/v9/migration"

	"github.com/filecoin-project/go-state-types/builtin"
	system9 "github.com/filecoin-project/go-state-types/builtin/v9/system"
	"github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"
)

func TestMigration(t *testing.T) {
	ctx := context.Background()
	bs := cbor.NewCborStore(NewSyncBlockStoreInMemory())
	adtStore := adt.WrapStore(ctx, bs)

	startRoot := makeInputTree(ctx, t, adtStore)

	oldStateTree, err := migration.LoadTree(adtStore, startRoot)
	require.NoError(t, err)

	oldSystemActor, found, err := oldStateTree.GetActor(builtin.SystemActorAddr)
	require.NoError(t, err)
	require.True(t, found, "system actor not found")

	var oldSystemState system9.State
	err = adtStore.Get(ctx, oldSystemActor.Head, &oldSystemState)
	require.NoError(t, err)

	oldManifestDataCid := oldSystemState.BuiltinActors
	oldManifest := manifest.Manifest{
		Version: 1,
		Data:    oldManifestDataCid,
	}
	require.NoError(t, oldManifest.Load(ctx, adtStore), "failed to load old manifest")

	newManifestCid, newManifestDataCid := makeTestManifest(t, adtStore, "fil/9/")
	log := TestLogger{TB: t}

	cache := migration.NewMemMigrationCache()
	_, err = migration.MigrateStateTree(ctx, adtStore, newManifestCid, startRoot, 200, migration.Config{MaxWorkers: 1}, log, cache)
	require.NoError(t, err)

	cacheRoot, err := migration.MigrateStateTree(ctx, adtStore, newManifestCid, startRoot, 200, migration.Config{MaxWorkers: 1}, log, cache)
	require.NoError(t, err)

	noCacheRoot, err := migration.MigrateStateTree(ctx, adtStore, newManifestCid, startRoot, 200, migration.Config{MaxWorkers: 1}, log, migration.NewMemMigrationCache())
	require.NoError(t, err)
	require.True(t, cacheRoot.Equals(noCacheRoot))

	// check that the system actor state was correctly updated

	newStateTree, err := migration.LoadTree(adtStore, cacheRoot)
	require.NoError(t, err)

	newSystemActor, found, err := newStateTree.GetActor(builtin.SystemActorAddr)
	require.NoError(t, err)
	require.True(t, found, "system actor not found")

	var newSystemState system9.State
	err = adtStore.Get(ctx, newSystemActor.Head, &newSystemState)
	require.NoError(t, err)
	require.Equal(t, newManifestDataCid, newSystemState.BuiltinActors)

	var newManifestData manifest.ManifestData
	err = adtStore.Get(ctx, newManifestDataCid, &newManifestData)
	require.NoError(t, err)

	// check all the CIDs

	cidsMap := make(map[cid.Cid]cid.Cid)
	for _, entry := range newManifestData.Entries {
		newCid := entry.Code
		oldCid, ok := oldManifest.Get(entry.Name)
		require.True(t, ok, "didn't find entry in old manifest")
		cidsMap[oldCid] = newCid
	}

	_ = oldStateTree.ForEach(func(addr address.Address, oldActor *migration.Actor) error {
		newActor, ok, err := newStateTree.GetActor(addr)
		require.NoError(t, err, "failed to get actor")
		require.True(t, ok, "didn't find actor: %s", addr)
		expectedCid, ok := cidsMap[oldActor.Code]
		require.True(t, ok, "didn't find code in cidsmap")
		require.Equal(t, expectedCid, newActor.Code)
		require.Equal(t, oldActor.Balance, newActor.Balance)
		require.Equal(t, oldActor.CallSeqNum, newActor.CallSeqNum)
		return nil
	})

	// TODO: Add any subsequent state change checks

}
