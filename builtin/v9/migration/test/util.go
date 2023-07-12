package migration

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/filecoin-project/go-state-types/actors"

	"github.com/filecoin-project/go-state-types/rt"

	block "github.com/ipfs/go-block-format"
	ipldcbor "github.com/ipfs/go-ipld-cbor"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v8/account"
	"github.com/filecoin-project/go-state-types/builtin/v8/cron"
	_init "github.com/filecoin-project/go-state-types/builtin/v8/init"
	"github.com/filecoin-project/go-state-types/builtin/v8/market"
	"github.com/filecoin-project/go-state-types/builtin/v8/power"
	"github.com/filecoin-project/go-state-types/builtin/v8/reward"
	"github.com/filecoin-project/go-state-types/builtin/v8/verifreg"
	"github.com/filecoin-project/go-state-types/builtin/v9/system"
	"github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	"github.com/filecoin-project/go-state-types/cbor"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
)

func makeTestManifest(t *testing.T, store adt.Store, prefix string) (cid.Cid, cid.Cid) {
	builder := cid.V1Builder{Codec: cid.Raw, MhType: mh.IDENTITY}

	newManifestData := manifest.ManifestData{}
	for _, name := range manifest.GetBuiltinActorsKeys(actors.Version9) {
		codeCid, err := builder.Sum([]byte(fmt.Sprintf("%s%s", prefix, name)))
		if err != nil {
			t.Fatal(err)
		}

		newManifestData.Entries = append(newManifestData.Entries,
			manifest.ManifestEntry{
				Name: name,
				Code: codeCid,
			})
	}

	manifestDataCid, err := store.Put(context.Background(), &newManifestData)
	if err != nil {
		t.Fatal(err)
	}

	newManifest := manifest.Manifest{
		Version: 1,
		Data:    manifestDataCid,
	}

	manifestCid, err := store.Put(context.Background(), &newManifest)
	if err != nil {
		t.Fatal(err)
	}

	return manifestCid, manifestDataCid
}

func makeInputTree(ctx context.Context, t *testing.T, store adt.Store) cid.Cid {
	tree, err := builtin.NewTree(store)
	require.NoError(t, err, "failed to create empty actors tree")

	manifestCid, manifestDataCid := makeTestManifest(t, store, "fil/8/")
	var actorsManifest manifest.Manifest
	require.NoError(t, store.Get(ctx, manifestCid, &actorsManifest), "error reading actor manifest")
	require.NoError(t, actorsManifest.Load(ctx, store), "error loading actor manifest")

	accountCid, ok := actorsManifest.Get("account")
	require.True(t, ok, "didn't find account actor in manifest")

	systemCid, ok := actorsManifest.Get("system")
	require.True(t, ok, "didn't find system actor in manifest")

	systemState, err := system.ConstructState(store)
	require.NoError(t, err, "failed to construct system state")
	systemState.BuiltinActors = manifestDataCid
	initializeActor(ctx, t, tree, store, systemState, systemCid, builtin.SystemActorAddr, big.Zero())

	initCid, ok := actorsManifest.Get("init")
	require.True(t, ok, "didn't find init actor in manifest")

	initState, err := _init.ConstructState(store, "migrationtest")
	require.NoError(t, err)
	initializeActor(ctx, t, tree, store, initState, initCid, builtin.InitActorAddr, big.Zero())

	rewardCid, ok := actorsManifest.Get("reward")
	require.True(t, ok, "didn't find reward actor in manifest")

	rewardState := reward.ConstructState(abi.NewStoragePower(0))
	initializeActor(ctx, t, tree, store, rewardState, rewardCid, builtin.RewardActorAddr, big.Mul(big.NewInt(1_100_000_000), big.NewInt(1e18)))

	cronCid, ok := actorsManifest.Get("cron")
	require.True(t, ok, "didn't find cron actor in manifest")

	cronState := cron.ConstructState(cron.BuiltInEntries())
	initializeActor(ctx, t, tree, store, cronState, cronCid, builtin.CronActorAddr, big.Zero())

	powerCid, ok := actorsManifest.Get("storagepower")
	require.True(t, ok, "didn't find power actor in manifest")

	powerState, err := power.ConstructState(store)
	require.NoError(t, err)
	initializeActor(ctx, t, tree, store, powerState, powerCid, builtin.StoragePowerActorAddr, big.Zero())

	marketCid, ok := actorsManifest.Get("storagemarket")
	require.True(t, ok, "didn't find market actor in manifest")

	marketState, err := market.ConstructState(store)
	require.NoError(t, err)
	initializeActor(ctx, t, tree, store, marketState, marketCid, builtin.StorageMarketActorAddr, big.Zero())

	// this will need to be replaced with the address of a multisig actor for the verified registry to be tested accurately
	VerifregRoot, err := address.NewIDAddress(80)
	require.NoError(t, err, "failed to create verifreg root")
	initializeActor(ctx, t, tree, store, &account.State{Address: VerifregRoot}, accountCid, VerifregRoot, big.Zero())

	verifregCid, ok := actorsManifest.Get("verifiedregistry")
	require.True(t, ok, "didn't find verifreg actor in manifest")

	vrState, err := verifreg.ConstructState(store, VerifregRoot)
	require.NoError(t, err)
	initializeActor(ctx, t, tree, store, vrState, verifregCid, builtin.VerifiedRegistryActorAddr, big.Zero())

	// burnt funds
	initializeActor(ctx, t, tree, store, &account.State{Address: builtin.BurntFundsActorAddr}, accountCid, builtin.BurntFundsActorAddr, big.Zero())

	root, err := tree.Flush()
	require.NoError(t, err, "failed to flush actors tree")
	return root
}

func initializeActor(ctx context.Context, t testing.TB, tree *builtin.ActorTree, store adt.Store, state cbor.Marshaler, code cid.Cid, a address.Address, balance abi.TokenAmount) {
	stateCID, err := store.Put(ctx, state)
	require.NoError(t, err)
	actor := &builtin.ActorV4{
		Head:    stateCID,
		Code:    code,
		Balance: balance,
	}
	err = tree.SetActorV4(a, actor)
	require.NoError(t, err)
}

type BlockStoreInMemory struct {
	data map[cid.Cid]block.Block
}

func NewBlockStoreInMemory() *BlockStoreInMemory {
	return &BlockStoreInMemory{make(map[cid.Cid]block.Block)}
}

func (mb *BlockStoreInMemory) Get(ctx context.Context, c cid.Cid) (block.Block, error) {
	d, ok := mb.data[c]
	if ok {
		return d, nil
	}
	return nil, fmt.Errorf("not found")
}

func (mb *BlockStoreInMemory) Put(ctx context.Context, b block.Block) error {
	mb.data[b.Cid()] = b
	return nil
}

// Creates a new, empty IPLD store in memory.
func NewADTStore(ctx context.Context) adt.Store {
	return adt.WrapStore(ctx, ipldcbor.NewCborStore(NewBlockStoreInMemory()))

}

type SyncBlockStoreInMemory struct {
	bs *BlockStoreInMemory
	mu sync.Mutex
}

func NewSyncBlockStoreInMemory() *SyncBlockStoreInMemory {
	return &SyncBlockStoreInMemory{
		bs: NewBlockStoreInMemory(),
	}
}

func (ss *SyncBlockStoreInMemory) Get(ctx context.Context, c cid.Cid) (block.Block, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.bs.Get(ctx, c)
}

func (ss *SyncBlockStoreInMemory) Put(ctx context.Context, b block.Block) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.bs.Put(ctx, b)
}

type TestLogger struct {
	TB testing.TB
}

func (t TestLogger) Log(_ rt.LogLevel, msg string, args ...interface{}) {
	t.TB.Logf(msg, args...)
}

func MakeCID(input string, prefix *cid.Prefix) cid.Cid {
	data := []byte(input)
	if prefix == nil {
		c, err := abi.CidBuilder.Sum(data)
		if err != nil {
			panic(err)
		}
		return c
	}
	c, err := prefix.Sum(data)
	switch {
	case errors.Is(err, mh.ErrSumNotSupported):
		// multihash library doesn't support this hash function.
		// just fake it.
	case err == nil:
		return c
	default:
		panic(err)
	}

	sum := sha256.Sum256(data)
	hash, err := mh.Encode(sum[:], prefix.MhType)
	if err != nil {
		panic(err)
	}
	return cid.NewCidV1(prefix.Codec, hash)
}
