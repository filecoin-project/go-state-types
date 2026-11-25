package v20_test

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	account19 "github.com/filecoin-project/go-state-types/builtin/v19/account"
	cron19 "github.com/filecoin-project/go-state-types/builtin/v19/cron"
	datacap19 "github.com/filecoin-project/go-state-types/builtin/v19/datacap"
	init19 "github.com/filecoin-project/go-state-types/builtin/v19/init"
	market19 "github.com/filecoin-project/go-state-types/builtin/v19/market"
	power19 "github.com/filecoin-project/go-state-types/builtin/v19/power"
	reward19 "github.com/filecoin-project/go-state-types/builtin/v19/reward"
	system19 "github.com/filecoin-project/go-state-types/builtin/v19/system"
	adt19 "github.com/filecoin-project/go-state-types/builtin/v19/util/adt"
	verifreg19 "github.com/filecoin-project/go-state-types/builtin/v19/verifreg"
	v20 "github.com/filecoin-project/go-state-types/builtin/v20"
	migration20 "github.com/filecoin-project/go-state-types/builtin/v20/migration"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/filecoin-project/go-state-types/rt"
	"github.com/filecoin-project/go-state-types/test_util"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"
	cbg "github.com/whyrusleeping/cbor-gen"
)

func TestCheckStateInvariantsAfterFilPlusRemoval(t *testing.T) {
	ctx := context.Background()
	store := cbor.NewCborStore(test_util.NewSyncBlockStoreInMemory())
	adtStore := adt19.WrapStore(ctx, store)
	tree, err := builtin.NewTree(adtStore)
	require.NoError(t, err)

	// Use real v19 actor states so the checker must decode the migrated tree.
	oldCodes, newCodes := make(map[string]cid.Cid), make(map[string]cid.Cid)
	var oldData, newData manifest.ManifestData
	for i, name := range []string{
		"system", "init", "cron", "account", "storagepower", "storageminer",
		"storagemarket", "paymentchannel", "multisig", "reward",
		"verifiedregistry", "datacap", "placeholder", "evm", "eam", "ethaccount",
	} {
		oldTag, newTag := cbg.CborInt(i+100), cbg.CborInt(i+200)
		oldCodes[name], err = store.Put(ctx, &oldTag)
		require.NoError(t, err)
		oldData.Entries = append(oldData.Entries, manifest.ManifestEntry{Name: name, Code: oldCodes[name]})
		if name == "verifiedregistry" || name == "datacap" {
			continue
		}
		newCodes[name], err = store.Put(ctx, &newTag)
		require.NoError(t, err)
		newData.Entries = append(newData.Entries, manifest.ManifestEntry{Name: name, Code: newCodes[name]})
	}
	oldManifestData, err := store.Put(ctx, &oldData)
	require.NoError(t, err)
	newManifestData, err := store.Put(ctx, &newData)
	require.NoError(t, err)
	newManifest, err := store.Put(ctx, &manifest.Manifest{Version: 1, Data: newManifestData})
	require.NoError(t, err)

	install := func(id uint64, name string, state any, balance abi.TokenAmount) {
		t.Helper()
		addr, err := address.NewIDAddress(id)
		require.NoError(t, err)
		head, err := store.Put(ctx, state)
		require.NoError(t, err)
		require.NoError(t, tree.SetActorV5(addr, &builtin.ActorV5{
			Code: oldCodes[name], Head: head, Balance: balance,
		}))
	}

	rootAddress, err := address.NewIDAddress(80)
	require.NoError(t, err)
	rootKey, err := address.NewSecp256k1Address([]byte("FIP-0121 root account"))
	require.NoError(t, err)
	initState, err := init19.ConstructState(adtStore, "fip0121-test")
	require.NoError(t, err)
	addressMap, err := adt19.AsMap(adtStore, initState.AddressMap, builtin.DefaultHamtBitwidth)
	require.NoError(t, err)
	rootID := cbg.CborInt(80)
	require.NoError(t, addressMap.Put(abi.AddrKey(rootKey), &rootID))
	verifierKey, err := address.NewBLSAddress(make([]byte, 48))
	require.NoError(t, err)
	verifierID := cbg.CborInt(81)
	require.NoError(t, addressMap.Put(abi.AddrKey(verifierKey), &verifierID))
	initState.AddressMap, err = addressMap.Root()
	require.NoError(t, err)

	rewardState, err := reward19.ConstructState(adtStore, big.Zero())
	require.NoError(t, err)
	powerState, err := power19.ConstructState(adtStore)
	require.NoError(t, err)
	marketState, err := market19.ConstructState(adtStore)
	require.NoError(t, err)
	verifregState, err := verifreg19.ConstructState(adtStore, rootAddress)
	require.NoError(t, err)
	datacapState, err := datacap19.ConstructState(adtStore, builtin.VerifiedRegistryActorAddr, builtin.DefaultHamtBitwidth)
	require.NoError(t, err)

	install(0, "system", &system19.State{BuiltinActors: oldManifestData}, big.Zero())
	install(1, "init", initState, big.Zero())
	install(2, "reward", rewardState, reward19.StorageMiningAllocationCheck)
	install(3, "cron", cron19.ConstructState(cron19.BuiltInEntries()), big.Zero())
	install(4, "storagepower", powerState, big.Zero())
	install(5, "storagemarket", marketState, big.Zero())
	install(6, "verifiedregistry", verifregState, big.NewInt(3))
	install(7, "datacap", datacapState, big.NewInt(5))
	install(10, "eam", []struct{}{}, big.Zero())
	install(80, "account", &account19.State{Address: rootKey}, big.NewInt(7))
	install(81, "account", &account19.State{Address: verifierKey}, big.NewInt(11))
	install(99, "account", &account19.State{Address: builtin.BurntFundsActorAddr},
		big.Sub(big.Sub(builtin.TotalFilecoin, reward19.StorageMiningAllocationCheck), big.NewInt(26)))
	root, err := tree.Flush()
	require.NoError(t, err)

	priorEpoch := rewardState.Epoch - 1
	migratedRoot, err := migration20.MigrateStateTree(ctx, store, newManifest, root, priorEpoch,
		migration.Config{MaxWorkers: 2}, invariantMigrationLogger{t}, migration.NewMemMigrationCache())
	require.NoError(t, err)
	out, err := builtin.LoadTree(adtStore, migratedRoot)
	require.NoError(t, err)
	for _, id := range []uint64{6, 7, 80, 81} {
		addr, err := address.NewIDAddress(id)
		require.NoError(t, err)
		_, found, err := out.GetActorV5(addr)
		require.NoError(t, err)
		require.False(t, found, "removed actor %d must be absent", id)
	}

	acc, err := v20.CheckStateInvariants(out, priorEpoch, newCodes)
	require.NoError(t, err)
	require.Empty(t, acc.Messages())

	// The checker must still catch a lost balance after removing the Fil+ checks.
	burn, found, err := out.GetActorV5(builtin.BurntFundsActorAddr)
	require.NoError(t, err)
	require.True(t, found)
	burn.Balance = big.Sub(burn.Balance, big.NewInt(1))
	require.NoError(t, out.SetActorV5(builtin.BurntFundsActorAddr, burn))
	acc, err = v20.CheckStateInvariants(out, priorEpoch, newCodes)
	require.NoError(t, err)
	require.Len(t, acc.Messages(), 1)
	require.Contains(t, acc.Messages()[0], "total token balance")
}

type invariantMigrationLogger struct{ t *testing.T }

func (l invariantMigrationLogger) Log(_ rt.LogLevel, format string, args ...interface{}) {
	l.t.Logf(format, args...)
}
