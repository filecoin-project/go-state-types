package migration

import (
	"context"
	"testing"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	market18 "github.com/filecoin-project/go-state-types/builtin/v18/market"
	adt18 "github.com/filecoin-project/go-state-types/builtin/v18/util/adt"
	verifreg18 "github.com/filecoin-project/go-state-types/builtin/v18/verifreg"
	market19 "github.com/filecoin-project/go-state-types/builtin/v19/market"
	"github.com/filecoin-project/go-state-types/migration"
	cid "github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"
	cbg "github.com/whyrusleeping/cbor-gen"
)

// populatedMarketState builds a v18 market state with a distinct CID or value in every field,
// so a field-shift during migration shows up as a wrong value rather than merely a missing one.
func populatedMarketState(t *testing.T, ctx context.Context, store cbor.IpldStore, pendingAllocs map[abi.DealID]verifreg18.AllocationId) market18.State {
	t.Helper()
	req := require.New(t)
	adtStore := adt18.WrapStore(ctx, store)

	pending, err := adt18.MakeEmptyMap(adtStore, builtin.DefaultHamtBitwidth)
	req.NoError(err)
	for dealID, allocID := range pendingAllocs {
		v := cbg.CborInt(allocID)
		req.NoError(pending.Put(abi.UIntKey(uint64(dealID)), &v))
	}
	pendingRoot, err := pending.Root()
	req.NoError(err)

	// Distinct, real CIDs so a field-shift shows up as the wrong value, not just a missing one.
	distinct := func(tag int64) cid.Cid {
		v := cbg.CborInt(tag)
		c, err := store.Put(ctx, &v)
		req.NoError(err)
		return c
	}

	return market18.State{
		Proposals:                     distinct(1),
		States:                        distinct(2),
		PendingProposals:              distinct(3),
		EscrowTable:                   distinct(4),
		LockedTable:                   distinct(5),
		NextID:                        abi.DealID(1234),
		DealOpsByEpoch:                distinct(6),
		LastCron:                      abi.ChainEpoch(5678),
		TotalClientLockedCollateral:   abi.NewTokenAmount(11),
		TotalProviderLockedCollateral: abi.NewTokenAmount(22),
		TotalClientStorageFee:         abi.NewTokenAmount(33),
		PendingDealAllocationIds:      pendingRoot,
		ProviderSectors:               distinct(7),
	}
}

// Every surviving field must arrive with its own value: the removed field sits mid-tuple, so
// a mistake shifts ProviderSectors and everything after it.
func TestMarketMigration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := require.New(t)
	store := cbor.NewMemCborStore()

	inState := populatedMarketState(t, ctx, store, nil)
	inHead, err := store.Put(ctx, &inState)
	req.NoError(err)

	outCodeCID := cid.MustParse("bafy2bzaca4aaaaaaaaaqk")
	migrator := marketMigrator{OutCodeCID: outCodeCID}
	req.Equal(outCodeCID, migrator.MigratedCodeCID())
	req.False(migrator.Deferred())

	result, err := migrator.MigrateState(ctx, store, migration.ActorMigrationInput{Address: address.TestAddress, Head: inHead})
	req.NoError(err)
	req.Equal(outCodeCID, result.NewCodeCID)

	var outState market19.State
	req.NoError(store.Get(ctx, result.NewHead, &outState))
	req.Equal(inState.Proposals, outState.Proposals)
	req.Equal(inState.States, outState.States)
	req.Equal(inState.PendingProposals, outState.PendingProposals)
	req.Equal(inState.EscrowTable, outState.EscrowTable)
	req.Equal(inState.LockedTable, outState.LockedTable)
	req.Equal(inState.NextID, outState.NextID)
	req.Equal(inState.DealOpsByEpoch, outState.DealOpsByEpoch)
	req.Equal(inState.LastCron, outState.LastCron)
	req.Equal(inState.TotalClientLockedCollateral, outState.TotalClientLockedCollateral)
	req.Equal(inState.TotalProviderLockedCollateral, outState.TotalProviderLockedCollateral)
	req.Equal(inState.TotalClientStorageFee, outState.TotalClientStorageFee)
	// The field that follows the removed one: this is what a shift would corrupt.
	req.Equal(inState.ProviderSectors, outState.ProviderSectors)
}

// A non-empty map is the ordinary case, not an error: verified deals published but not yet
// activated hold entries at any epoch. The migration must drop them and carry on, because
// failing here would halt the network upgrade on perfectly normal state.
func TestMarketMigrationDropsPendingAllocations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := require.New(t)
	store := cbor.NewMemCborStore()

	inState := populatedMarketState(t, ctx, store, map[abi.DealID]verifreg18.AllocationId{7: 70, 8: 80})
	inHead, err := store.Put(ctx, &inState)
	req.NoError(err)

	outCodeCID := cid.MustParse("bafy2bzaca4aaaaaaaaaqk")
	result, err := marketMigrator{OutCodeCID: outCodeCID}.
		MigrateState(ctx, store, migration.ActorMigrationInput{Address: address.TestAddress, Head: inHead})
	req.NoError(err)

	// Every surviving field still arrives intact; only the allocation bookkeeping is gone.
	var outState market19.State
	req.NoError(store.Get(ctx, result.NewHead, &outState))
	req.Equal(inState.Proposals, outState.Proposals)
	req.Equal(inState.ProviderSectors, outState.ProviderSectors)
	req.Equal(inState.NextID, outState.NextID)
}

// The v19 state must round-trip through CBOR without the removed field, and decode to the
// same values it was written with.
func TestMarketStateRoundTripsWithoutPendingAllocations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := require.New(t)
	store := cbor.NewMemCborStore()

	inState := populatedMarketState(t, ctx, store, nil)
	inHead, err := store.Put(ctx, &inState)
	req.NoError(err)
	result, err := marketMigrator{OutCodeCID: cid.MustParse("bafy2bzaca4aaaaaaaaaqk")}.
		MigrateState(ctx, store, migration.ActorMigrationInput{Address: address.TestAddress, Head: inHead})
	req.NoError(err)

	var outState market19.State
	req.NoError(store.Get(ctx, result.NewHead, &outState))
	rewritten, err := store.Put(ctx, &outState)
	req.NoError(err)
	req.Equal(result.NewHead, rewritten, "v19 market state does not round-trip")

	// A v18 decoder must not accept it: the tuple is one field shorter.
	var asV18 market18.State
	req.Error(store.Get(ctx, result.NewHead, &asV18))
}
