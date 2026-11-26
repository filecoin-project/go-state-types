package migration

import (
	"context"
	"fmt"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	account19 "github.com/filecoin-project/go-state-types/builtin/v19/account"
	datacap19 "github.com/filecoin-project/go-state-types/builtin/v19/datacap"
	init19 "github.com/filecoin-project/go-state-types/builtin/v19/init"
	system19 "github.com/filecoin-project/go-state-types/builtin/v19/system"
	verifreg19 "github.com/filecoin-project/go-state-types/builtin/v19/verifreg"
	init20 "github.com/filecoin-project/go-state-types/builtin/v20/init"
	system20 "github.com/filecoin-project/go-state-types/builtin/v20/system"
	adt20 "github.com/filecoin-project/go-state-types/builtin/v20/util/adt"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/filecoin-project/go-state-types/rt"
	"github.com/filecoin-project/go-state-types/test_util"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"
	cbg "github.com/whyrusleeping/cbor-gen"
)

// Keep these expectations independent of GetBuiltinActorsKeys, so a mistake in
// that helper cannot silently change both the migration fixture and its oracle.
var fip121OldActorNames = []string{
	manifest.AccountKey, manifest.CronKey, manifest.InitKey, manifest.MarketKey,
	manifest.MinerKey, manifest.MultisigKey, manifest.PaychKey, manifest.PowerKey,
	manifest.RewardKey, manifest.SystemKey, manifest.VerifregKey, manifest.DatacapKey,
	manifest.EvmKey, manifest.EamKey, manifest.PlaceholderKey, manifest.EthAccountKey,
}

var fip121NewActorNames = []string{
	manifest.AccountKey, manifest.CronKey, manifest.InitKey, manifest.MarketKey,
	manifest.MinerKey, manifest.MultisigKey, manifest.PaychKey, manifest.PowerKey,
	manifest.RewardKey, manifest.SystemKey, manifest.EvmKey, manifest.EamKey,
	manifest.PlaceholderKey, manifest.EthAccountKey,
}

type fip121MigrationFixture struct {
	t                  *testing.T
	ctx                context.Context
	store              cbor.IpldStore
	adtStore           adt20.Store
	actors             *builtin.ActorTree
	oldCodes           map[string]cid.Cid
	newCodes           map[string]cid.Cid
	oldManifestData    cid.Cid
	newManifestData    cid.Cid
	newManifest        cid.Cid
	newManifestEntries manifest.ManifestData
}

func newFIP121MigrationFixture(t *testing.T) *fip121MigrationFixture {
	t.Helper()
	f := &fip121MigrationFixture{
		t: t, ctx: context.Background(),
		store:    cbor.NewCborStore(test_util.NewSyncBlockStoreInMemory()),
		oldCodes: make(map[string]cid.Cid), newCodes: make(map[string]cid.Cid),
	}
	f.adtStore = adt20.WrapStore(f.ctx, f.store)
	var err error
	f.actors, err = builtin.NewTree(f.adtStore)
	require.NoError(t, err)

	var oldData manifest.ManifestData
	for i, name := range fip121OldActorNames {
		tag := cbg.CborInt(1000 + i)
		code := f.put(&tag)
		f.oldCodes[name] = code
		oldData.Entries = append(oldData.Entries, manifest.ManifestEntry{Name: name, Code: code})
	}
	for i, name := range fip121NewActorNames {
		tag := cbg.CborInt(2000 + i)
		code := f.put(&tag)
		f.newCodes[name] = code
		f.newManifestEntries.Entries = append(f.newManifestEntries.Entries, manifest.ManifestEntry{Name: name, Code: code})
	}
	f.oldManifestData = f.put(&oldData)
	f.newManifestData = f.put(&f.newManifestEntries)
	f.newManifest = f.put(&manifest.Manifest{Version: 1, Data: f.newManifestData})

	// Every surviving actor type is represented, including a multisig other than
	// actor 80 and an account at the neighboring ID 82. Opaque heads are sufficient
	// for actors whose state is unchanged.
	for _, actor := range []struct {
		id   uint64
		name string
	}{
		{0, manifest.SystemKey}, {1, manifest.InitKey}, {2, manifest.RewardKey},
		{3, manifest.CronKey}, {4, manifest.PowerKey}, {5, manifest.MarketKey},
		{6, manifest.VerifregKey}, {7, manifest.DatacapKey}, {10, manifest.EamKey},
		{80, manifest.MultisigKey}, {81, manifest.AccountKey}, {82, manifest.AccountKey},
		{99, manifest.AccountKey}, {100, manifest.AccountKey},
		{101, manifest.MinerKey}, {102, manifest.MultisigKey}, {103, manifest.PaychKey},
		{104, manifest.PlaceholderKey}, {105, manifest.EvmKey}, {106, manifest.EthAccountKey},
	} {
		tag := cbg.CborInt(3000 + actor.id)
		record := &builtin.ActorV5{
			Code: f.oldCodes[actor.name], Head: f.put(&tag),
			Balance: abi.NewTokenAmount(int64(actor.id + 1)),
		}
		if actor.name == manifest.AccountKey {
			record.CallSeqNum = actor.id + 10
		}
		if actor.name == manifest.EvmKey || actor.name == manifest.EthAccountKey || actor.name == manifest.PlaceholderKey {
			delegated, err := address.NewDelegatedAddress(10, []byte{byte(actor.id)})
			require.NoError(t, err)
			record.DelegatedAddress = &delegated
		}
		require.NoError(t, f.actors.SetActorV5(f.id(actor.id), record))
	}
	f.setHead(0, f.put(&system19.State{BuiltinActors: f.oldManifestData}))
	initState, err := init19.ConstructState(f.adtStore, "fip121-test")
	require.NoError(t, err)
	f.setHead(1, f.put(initState))
	verifregState, err := verifreg19.ConstructState(f.adtStore, f.id(80))
	require.NoError(t, err)
	verifiers, err := adt20.AsMap(f.adtStore, verifregState.Verifiers, builtin.DefaultHamtBitwidth)
	require.NoError(t, err)
	allowance := abi.NewStoragePower(456)
	require.NoError(t, verifiers.Put(abi.AddrKey(f.id(120)), &allowance))
	verifregState.Verifiers, err = verifiers.Root()
	require.NoError(t, err)
	f.setHead(6, f.put(verifregState))
	datacapState, err := datacap19.ConstructState(f.adtStore, f.id(6), builtin.DefaultHamtBitwidth)
	require.NoError(t, err)
	f.setHead(7, f.put(datacapState))
	return f
}

func (f *fip121MigrationFixture) id(id uint64) address.Address {
	f.t.Helper()
	a, err := address.NewIDAddress(id)
	require.NoError(f.t, err)
	return a
}

func (f *fip121MigrationFixture) put(value interface{}) cid.Cid {
	f.t.Helper()
	c, err := f.store.Put(f.ctx, value)
	require.NoError(f.t, err)
	return c
}

func (f *fip121MigrationFixture) actor(id uint64) *builtin.ActorV5 {
	f.t.Helper()
	a, found, err := f.actors.GetActorV5(f.id(id))
	require.NoError(f.t, err)
	require.True(f.t, found)
	return a
}

func (f *fip121MigrationFixture) setHead(id uint64, head cid.Cid) {
	f.t.Helper()
	a := f.actor(id)
	a.Head = head
	require.NoError(f.t, f.actors.SetActorV5(f.id(id), a))
}

func (f *fip121MigrationFixture) root() cid.Cid {
	f.t.Helper()
	root, err := f.actors.Flush()
	require.NoError(f.t, err)
	return root
}

func (f *fip121MigrationFixture) snapshot(root cid.Cid) map[address.Address]builtin.ActorV5 {
	f.t.Helper()
	tree, err := builtin.LoadTree(f.adtStore, root)
	require.NoError(f.t, err)
	actors := make(map[address.Address]builtin.ActorV5)
	require.NoError(f.t, tree.ForEachV5(func(addr address.Address, actor *builtin.ActorV5) error {
		actors[addr] = *actor
		return nil
	}))
	return actors
}

func (f *fip121MigrationFixture) migrate(root cid.Cid, workers uint, cache migration.MigrationCache) (cid.Cid, error) {
	f.t.Helper()
	return MigrateStateTree(f.ctx, f.store, f.newManifest, root, 100,
		migration.Config{MaxWorkers: workers}, fip121MigrationLogger{f.t}, cache)
}

type fip121MigrationLogger struct{ t *testing.T }

func (l fip121MigrationLogger) Log(_ rt.LogLevel, format string, args ...interface{}) {
	l.t.Logf(format, args...)
}

func (f *fip121MigrationFixture) requireMigrated(root, outRoot cid.Cid, before map[address.Address]builtin.ActorV5) {
	f.t.Helper()
	after := f.snapshot(outRoot)
	removed := make(map[address.Address]bool)
	burnt := before[f.id(99)].Balance
	for _, id := range []uint64{6, 7, 80, 81} {
		addr := f.id(id)
		if actor, found := before[addr]; found {
			burnt = big.Add(burnt, actor.Balance)
			removed[addr] = true
		}
		require.NotContains(f.t, after, addr)
	}
	require.Len(f.t, after, len(before)-len(removed))
	oldTotal, newTotal := big.Zero(), big.Zero()
	for addr, oldActor := range before {
		oldTotal = big.Add(oldTotal, oldActor.Balance)
		if removed[addr] {
			continue
		}
		newActor, found := after[addr]
		require.True(f.t, found, "surviving actor %s is missing", addr)
		expected := oldActor
		for name, code := range f.oldCodes {
			if code == oldActor.Code {
				expected.Code = f.newCodes[name]
				break
			}
		}
		if addr == f.id(99) {
			expected.Balance = burnt
		}
		if addr == f.id(0) {
			var state system20.State
			require.NoError(f.t, f.store.Get(f.ctx, newActor.Head, &state))
			require.Equal(f.t, f.newManifestData, state.BuiltinActors)
			require.NotEqual(f.t, oldActor.Head, newActor.Head)
			expected.Head = newActor.Head
		}
		require.Equal(f.t, expected, newActor, "unexpected change to actor %s", addr)
	}
	for _, actor := range after {
		newTotal = big.Add(newTotal, actor.Balance)
	}
	require.True(f.t, oldTotal.Equals(newTotal), "total FIL changed: %s -> %s", oldTotal, newTotal)
	var data manifest.ManifestData
	require.NoError(f.t, f.store.Get(f.ctx, f.newManifestData, &data))
	require.Len(f.t, data.Entries, 14)
	require.Equal(f.t, f.newManifestEntries, data)
	for _, entry := range data.Entries {
		require.NotEqual(f.t, manifest.VerifregKey, entry.Name)
		require.NotEqual(f.t, manifest.DatacapKey, entry.Name)
	}
	f.requireHistory(root, before)
}

func (f *fip121MigrationFixture) requireHistory(root cid.Cid, before map[address.Address]builtin.ActorV5) {
	f.t.Helper()
	require.Equal(f.t, before, f.snapshot(root), "migration changed the original state tree")
	for addr, actor := range before {
		var state cbg.Deferred
		require.NoError(f.t, f.store.Get(f.ctx, actor.Head, &state), "historical state for %s is unreadable", addr)
	}
	var systemState system19.State
	require.NoError(f.t, f.store.Get(f.ctx, before[f.id(0)].Head, &systemState))
	require.Equal(f.t, f.oldManifestData, systemState.BuiltinActors)
	var oldData manifest.ManifestData
	require.NoError(f.t, f.store.Get(f.ctx, systemState.BuiltinActors, &oldData))
	require.Len(f.t, oldData.Entries, len(fip121OldActorNames))
	for i, entry := range oldData.Entries {
		require.Equal(f.t, fip121OldActorNames[i], entry.Name)
		require.Equal(f.t, f.oldCodes[entry.Name], entry.Code)
	}
	if actor, found := before[f.id(6)]; found {
		var state verifreg19.State
		require.NoError(f.t, f.store.Get(f.ctx, actor.Head, &state))
		verifiers, err := adt20.AsMap(f.adtStore, state.Verifiers, builtin.DefaultHamtBitwidth)
		require.NoError(f.t, err)
		var allowance abi.StoragePower
		found, err := verifiers.Get(abi.AddrKey(f.id(120)), &allowance)
		require.NoError(f.t, err)
		require.True(f.t, found)
		require.Equal(f.t, abi.NewStoragePower(456), allowance)
	}
	if actor, found := before[f.id(7)]; found {
		var state datacap19.State
		require.NoError(f.t, f.store.Get(f.ctx, actor.Head, &state))
		_, err := adt20.AsMap(f.adtStore, state.Token.Balances, int(state.Token.HamtBitWidth))
		require.NoError(f.t, err)
	}
}

func (f *fip121MigrationFixture) requireFailure(wantError string) {
	f.t.Helper()
	root := f.root()
	before := f.snapshot(root)
	out, err := f.migrate(root, 2, migration.NewMemMigrationCache())
	require.ErrorContains(f.t, err, wantError)
	require.Equal(f.t, cid.Undef, out)
	f.requireHistory(root, before)
}

func TestFIP121MigrationBalancesAndState(t *testing.T) {
	for _, tc := range []struct {
		name     string
		balances [5]int64 // IDs 6, 7, 80, 81, 99.
	}{
		{"nonzero balances", [5]int64{11, 23, 37, 53, 101}},
		{"zero retired balances", [5]int64{0, 0, 0, 0, 101}},
		{"zero burn balance", [5]int64{11, 23, 37, 53, 0}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFIP121MigrationFixture(t)
			for i, id := range []uint64{6, 7, 80, 81, 99} {
				actor := f.actor(id)
				actor.Balance = abi.NewTokenAmount(tc.balances[i])
				require.NoError(t, f.actors.SetActorV5(f.id(id), actor))
			}
			root := f.root()
			before := f.snapshot(root)
			out, err := f.migrate(root, 2, migration.NewMemMigrationCache())
			require.NoError(t, err)
			f.requireMigrated(root, out, before)
		})
	}
}

func TestFIP121MigrationRootActorVariants(t *testing.T) {
	for _, kind := range []string{manifest.MultisigKey, manifest.AccountKey, "unknown code", "absent"} {
		t.Run(kind, func(t *testing.T) {
			f := newFIP121MigrationFixture(t)
			actor := f.actor(80)
			switch kind {
			case "absent":
				require.NoError(t, f.actors.Map.Delete(abi.AddrKey(f.id(80))))
			case "unknown code":
				tag := cbg.CborInt(9000)
				actor.Code = f.put(&tag)
				require.NoError(t, f.actors.SetActorV5(f.id(80), actor))
			default:
				actor.Code = f.oldCodes[kind]
				require.NoError(t, f.actors.SetActorV5(f.id(80), actor))
			}
			root := f.root()
			before := f.snapshot(root)
			out, err := f.migrate(root, 1, migration.NewMemMigrationCache())
			require.NoError(t, err)
			f.requireMigrated(root, out, before)
		})
	}
}

func TestFIP121MigrationRejectsMissingRequiredActor(t *testing.T) {
	for _, tc := range []struct {
		id        uint64
		wantError string
	}{
		{6, "Verifreg actor not found"},
		{7, "Datacap actor not found"},
		{99, "burnt funds actor not found"},
	} {
		t.Run(fmt.Sprint(tc.id), func(t *testing.T) {
			f := newFIP121MigrationFixture(t)
			require.NoError(t, f.actors.Map.Delete(abi.AddrKey(f.id(tc.id))))
			f.requireFailure(tc.wantError)
		})
	}
}

func TestFIP121MigrationRejectsRetiredCodeAtOtherAddress(t *testing.T) {
	for _, name := range []string{manifest.VerifregKey, manifest.DatacapKey} {
		t.Run(name, func(t *testing.T) {
			f := newFIP121MigrationFixture(t)
			tag := cbg.CborInt(9001)
			require.NoError(t, f.actors.SetActorV5(f.id(200), &builtin.ActorV5{
				Code: f.oldCodes[name], Head: f.put(&tag), Balance: abi.NewTokenAmount(47),
			}))
			f.requireFailure(fmt.Sprintf("actor with code %s has no registered migration function", f.oldCodes[name]))
		})
	}
}

func TestFIP121MigrationRejectsInvalidManifest(t *testing.T) {
	for _, tc := range []struct {
		name      string
		add       []string
		remove    string
		wantError string
	}{
		{"retains verifreg", []string{manifest.VerifregKey}, "", "code for deleted actor verifiedregistry found in new manifest"},
		{"retains datacap", []string{manifest.DatacapKey}, "", "code for deleted actor datacap found in new manifest"},
		{"retains both", []string{manifest.VerifregKey, manifest.DatacapKey}, "", "code for deleted actor"},
		{"missing retained type", nil, manifest.MinerKey, "code cid for storageminer actor not found in new manifest"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFIP121MigrationFixture(t)
			var data manifest.ManifestData
			for _, entry := range f.newManifestEntries.Entries {
				if entry.Name != tc.remove {
					data.Entries = append(data.Entries, entry)
				}
			}
			for _, name := range tc.add {
				data.Entries = append(data.Entries, manifest.ManifestEntry{Name: name, Code: f.oldCodes[name]})
			}
			f.newManifest = f.put(&manifest.Manifest{Version: 1, Data: f.put(&data)})
			f.requireFailure(tc.wantError)
		})
	}
}

func TestFIP121MigrationPreservesRootAccountMapping(t *testing.T) {
	f := newFIP121MigrationFixture(t)
	key, err := address.NewSecp256k1Address([]byte("fip121 root account"))
	require.NoError(t, err)
	rootActor := f.actor(80)
	rootActor.Code = f.oldCodes[manifest.AccountKey]
	rootActor.Head = f.put(&account19.State{Address: key})
	require.NoError(t, f.actors.SetActorV5(f.id(80), rootActor))

	var state init19.State
	require.NoError(t, f.store.Get(f.ctx, f.actor(1).Head, &state))
	addressMap, err := adt20.AsMap(f.adtStore, state.AddressMap, builtin.DefaultHamtBitwidth)
	require.NoError(t, err)
	value := cbg.CborInt(80)
	require.NoError(t, addressMap.Put(abi.AddrKey(key), &value))
	state.AddressMap, err = addressMap.Root()
	require.NoError(t, err)
	f.setHead(1, f.put(&state))

	root := f.root()
	before := f.snapshot(root)
	out, err := f.migrate(root, 2, migration.NewMemMigrationCache())
	require.NoError(t, err)
	f.requireMigrated(root, out, before)
	tree, err := builtin.LoadTree(f.adtStore, out)
	require.NoError(t, err)
	initActor, found, err := tree.GetActorV5(f.id(1))
	require.NoError(t, err)
	require.True(t, found)
	var migrated init20.State
	require.NoError(t, f.store.Get(f.ctx, initActor.Head, &migrated))
	require.Equal(t, state.AddressMap, migrated.AddressMap)
	resolved, found, err := migrated.ResolveAddress(f.adtStore, key)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, f.id(80), resolved)
	_, found, err = tree.GetActorV5(resolved)
	require.NoError(t, err)
	require.False(t, found)
	_, messages := init20.CheckStateInvariants(&migrated, tree, f.newCodes)
	require.True(t, messages.IsEmpty(), "%v", messages.Messages())
}

func TestFIP121MigrationDeterministicAcrossWorkersAndCache(t *testing.T) {
	f := newFIP121MigrationFixture(t)
	root := f.root()
	before := f.snapshot(root)
	cache := migration.NewMemMigrationCache()
	expected, err := f.migrate(root, 1, cache)
	require.NoError(t, err)
	f.requireMigrated(root, expected, before)
	for _, tc := range []struct {
		name    string
		workers uint
		cache   migration.MigrationCache
	}{
		{"one worker warm cache", 1, cache},
		{"four workers warm cache", 4, cache},
		{"four workers fresh cache", 4, migration.NewMemMigrationCache()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture := *f
			fixture.t = t
			out, err := fixture.migrate(root, tc.workers, tc.cache)
			require.NoError(t, err)
			require.Equal(t, expected, out)
			fixture.requireMigrated(root, out, before)
		})
	}
}

func TestFIP121MigrationActor81Variants(t *testing.T) {
	for _, kind := range []string{manifest.AccountKey, manifest.MultisigKey, "unknown code", "absent", "both 80 and 81 absent"} {
		t.Run(kind, func(t *testing.T) {
			f := newFIP121MigrationFixture(t)
			actor := f.actor(81)
			switch kind {
			case "absent", "both 80 and 81 absent":
				require.NoError(t, f.actors.Map.Delete(abi.AddrKey(f.id(81))))
				if kind == "both 80 and 81 absent" {
					require.NoError(t, f.actors.Map.Delete(abi.AddrKey(f.id(80))))
				}
			case "unknown code":
				tag := cbg.CborInt(9002)
				actor.Code = f.put(&tag)
				require.NoError(t, f.actors.SetActorV5(f.id(81), actor))
			default:
				actor.Code = f.oldCodes[kind]
				require.NoError(t, f.actors.SetActorV5(f.id(81), actor))
			}
			root := f.root()
			before := f.snapshot(root)
			out, err := f.migrate(root, 2, migration.NewMemMigrationCache())
			require.NoError(t, err)
			f.requireMigrated(root, out, before)
		})
	}
}

func TestFIP121MigrationPreservesActor81AccountMapping(t *testing.T) {
	f := newFIP121MigrationFixture(t)
	key, err := address.NewSecp256k1Address([]byte("fip121 actor 81 account"))
	require.NoError(t, err)
	actor := f.actor(81)
	actor.Head = f.put(&account19.State{Address: key})
	require.NoError(t, f.actors.SetActorV5(f.id(81), actor))

	var state init19.State
	require.NoError(t, f.store.Get(f.ctx, f.actor(1).Head, &state))
	addressMap, err := adt20.AsMap(f.adtStore, state.AddressMap, builtin.DefaultHamtBitwidth)
	require.NoError(t, err)
	value := cbg.CborInt(81)
	require.NoError(t, addressMap.Put(abi.AddrKey(key), &value))
	state.AddressMap, err = addressMap.Root()
	require.NoError(t, err)
	f.setHead(1, f.put(&state))

	root := f.root()
	before := f.snapshot(root)
	out, err := f.migrate(root, 2, migration.NewMemMigrationCache())
	require.NoError(t, err)
	f.requireMigrated(root, out, before)
	tree, err := builtin.LoadTree(f.adtStore, out)
	require.NoError(t, err)
	initActor, found, err := tree.GetActorV5(f.id(1))
	require.NoError(t, err)
	require.True(t, found)
	var migrated init20.State
	require.NoError(t, f.store.Get(f.ctx, initActor.Head, &migrated))
	require.Equal(t, state.AddressMap, migrated.AddressMap)
	resolved, found, err := migrated.ResolveAddress(f.adtStore, key)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, f.id(81), resolved)
	_, found, err = tree.GetActorV5(resolved)
	require.NoError(t, err)
	require.False(t, found)
	_, messages := init20.CheckStateInvariants(&migrated, tree, f.newCodes)
	require.True(t, messages.IsEmpty(), "%v", messages.Messages())
}
