package market

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v19/util/adt"
	"github.com/filecoin-project/go-state-types/test_util"
	cid "github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"
	cbg "github.com/whyrusleeping/cbor-gen"
)

// NextUpdateEpoch mirrors the actors' arithmetic; these expectations are computed by hand from
// actors/market/src/lib.rs and pin the quantization, including the offset a non-zero deal ID
// introduces and the two cases where the result is not rounded up.
func TestNextUpdateEpoch(t *testing.T) {
	const interval = abi.ChainEpoch(100)

	// Deal 0 has offset 0, so the result is the next multiple of the interval at or after start.
	require.Equal(t, abi.ChainEpoch(0), NextUpdateEpoch(0, interval, 0))
	require.Equal(t, abi.ChainEpoch(100), NextUpdateEpoch(0, interval, 1))
	require.Equal(t, abi.ChainEpoch(100), NextUpdateEpoch(0, interval, 100))
	require.Equal(t, abi.ChainEpoch(200), NextUpdateEpoch(0, interval, 101))

	// A non-zero ID shifts the grid by id % interval; 7 lands the visits on ...07.
	require.Equal(t, abi.ChainEpoch(7), NextUpdateEpoch(7, interval, 0))
	require.Equal(t, abi.ChainEpoch(7), NextUpdateEpoch(7, interval, 7))
	require.Equal(t, abi.ChainEpoch(107), NextUpdateEpoch(7, interval, 8))
	require.Equal(t, abi.ChainEpoch(107), NextUpdateEpoch(7, interval, 107))
	// An ID beyond the interval wraps.
	require.Equal(t, abi.ChainEpoch(107), NextUpdateEpoch(207, interval, 8))

	// Negative earliest: truncating division already rounds up, so it must not round again.
	require.Equal(t, abi.ChainEpoch(-100), NextUpdateEpoch(0, interval, -100))
	require.Equal(t, abi.ChainEpoch(0), NextUpdateEpoch(0, interval, -99))
	require.Equal(t, abi.ChainEpoch(7), NextUpdateEpoch(7, interval, -50))

	// The production interval, with a deal ID whose offset is non-trivial.
	require.Equal(t, abi.ChainEpoch(DealUpdatesInterval), NextUpdateEpoch(0, DealUpdatesInterval, 1))
	require.Equal(t, abi.ChainEpoch(1234), NextUpdateEpoch(1234, DealUpdatesInterval, 0))
	require.Equal(t, abi.ChainEpoch(DealUpdatesInterval+1234), NextUpdateEpoch(1234, DealUpdatesInterval, 1235))
}

const testDealID = abi.DealID(3)
const testStartEpoch = abi.ChainEpoch(1000)

// firstVisit is where a never-visited deal must be queued.
func firstVisit(t *testing.T) abi.ChainEpoch {
	t.Helper()
	return NextUpdateEpoch(testDealID, DealUpdatesInterval, testStartEpoch)
}

func testStore(t *testing.T) adt.Store {
	t.Helper()
	return adt.WrapStore(context.Background(), cbor.NewCborStore(test_util.NewBlockStoreInMemory()))
}

func testAddr(t *testing.T, id uint64) address.Address {
	t.Helper()
	a, err := address.NewIDAddress(id)
	require.NoError(t, err)
	return a
}

// putProviderSector records the provider->sector->deal mapping that an activated deal must have.
func putProviderSector(store adt.Store, st *State, provider address.Address, sector abi.SectorNumber, deal abi.DealID) error {
	providerID, err := address.IDFromAddress(provider)
	if err != nil {
		return err
	}
	sectorMap, err := adt.MakeEmptyMap(store, ProviderSectorsHamtBitwidth)
	if err != nil {
		return err
	}
	dealIDs := SectorDealIDs{deal}
	if err := sectorMap.Put(abi.UIntKey(uint64(sector)), &dealIDs); err != nil {
		return err
	}
	sectorMapRoot, err := sectorMap.Root()
	if err != nil {
		return err
	}
	outer, err := adt.AsMap(store, st.ProviderSectors, ProviderSectorsHamtBitwidth)
	if err != nil {
		return err
	}
	if err := outer.Put(abi.UIntKey(providerID), (*cbg.CborCid)(&sectorMapRoot)); err != nil {
		return err
	}
	st.ProviderSectors, err = outer.Root()
	return err
}

// invariantState builds a market state holding one published deal, queued at `queuedAt`, with
// the deal's state set by `dealState` (nil for a deal that has never been visited).
func invariantState(t *testing.T, queuedAt []abi.ChainEpoch, dealState *DealState, lastCron abi.ChainEpoch) (*State, adt.Store) {
	t.Helper()
	store := testStore(t)
	st, err := ConstructState(store)
	require.NoError(t, err)
	st.LastCron = lastCron
	st.NextID = testDealID + 1

	proposal := DealProposal{
		PieceCID:             cid.MustParse("baga6ea4seaaqa"),
		PieceSize:            abi.PaddedPieceSize(2048),
		Client:               testAddr(t, 100),
		Provider:             testAddr(t, 200),
		Label:                EmptyDealLabel,
		StartEpoch:           testStartEpoch,
		EndEpoch:             testStartEpoch + DealMinDuration,
		StoragePricePerEpoch: big.Zero(),
		ProviderCollateral:   big.Zero(),
		ClientCollateral:     big.Zero(),
	}
	proposals, err := adt.AsArray(store, st.Proposals, ProposalsAmtBitwidth)
	require.NoError(t, err)
	require.NoError(t, proposals.Set(uint64(testDealID), &proposal))
	st.Proposals, err = proposals.Root()
	require.NoError(t, err)

	if dealState != nil {
		states, err := adt.AsArray(store, st.States, StatesAmtBitwidth)
		require.NoError(t, err)
		require.NoError(t, states.Set(uint64(testDealID), dealState))
		st.States, err = states.Root()
		require.NoError(t, err)
	} else {
		// Never activated, so it is still pending.
		pcid, err := proposal.Cid()
		require.NoError(t, err)
		pending, err := adt.AsMap(store, st.PendingProposals, builtin.DefaultHamtBitwidth)
		require.NoError(t, err)
		require.NoError(t, pending.Put(abi.CidKey(pcid), &proposal))
		st.PendingProposals, err = pending.Root()
		require.NoError(t, err)
	}

	if dealState != nil {
		require.NoError(t, putProviderSector(store, st, proposal.Provider, dealState.SectorNumber, testDealID))
	}

	// Queue the deal at each requested epoch.
	outer, err := adt.AsMap(store, st.DealOpsByEpoch, builtin.DefaultHamtBitwidth)
	require.NoError(t, err)
	for _, epoch := range queuedAt {
		inner, err := adt.MakeEmptySet(store, builtin.DefaultHamtBitwidth)
		require.NoError(t, err)
		require.NoError(t, inner.Put(abi.UIntKey(uint64(testDealID))))
		innerRoot, err := inner.Root()
		require.NoError(t, err)
		require.NoError(t, outer.Put(abi.UIntKey(uint64(epoch)), (*cbg.CborCid)(&innerRoot)))
	}
	st.DealOpsByEpoch, err = outer.Root()
	require.NoError(t, err)

	return st, store
}

func requireDealOpsOK(t *testing.T, st *State, store adt.Store, currEpoch abi.ChainEpoch) {
	t.Helper()
	_, acc := CheckStateInvariants(st, store, big.Zero(), currEpoch)
	require.True(t, acc.IsEmpty(), "unexpected invariant failures: %v", acc.Messages())
}

// A never-visited deal queued exactly once at its first visit satisfies the rule.
func TestDealOpsNeverVisitedQueuedCorrectly(t *testing.T) {
	fv := firstVisit(t)
	st, store := invariantState(t, []abi.ChainEpoch{fv}, nil, fv-1)
	requireDealOpsOK(t, st, store, testStartEpoch)
}

// Queued, but at the wrong epoch: cron would visit it at an epoch the actor never schedules.
func TestDealOpsNeverVisitedQueuedAtWrongEpoch(t *testing.T) {
	fv := firstVisit(t)
	st, store := invariantState(t, []abi.ChainEpoch{fv + 1}, nil, fv-1)
	_, acc := CheckStateInvariants(st, store, big.Zero(), testStartEpoch)
	require.False(t, acc.IsEmpty())
	require.Contains(t, acc.Messages()[0], "must have exactly one deal op")
}

// Queued twice: the deal would be visited more than once.
func TestDealOpsNeverVisitedQueuedTwice(t *testing.T) {
	fv := firstVisit(t)
	st, store := invariantState(t, []abi.ChainEpoch{fv, fv + DealUpdatesInterval}, nil, fv-1)
	_, acc := CheckStateInvariants(st, store, big.Zero(), testStartEpoch)
	require.False(t, acc.IsEmpty())
	require.Contains(t, acc.Messages()[0], "must have exactly one deal op")
}

// Not queued at all: the deal would never be visited and never cleaned up.
func TestDealOpsNeverVisitedNotQueued(t *testing.T) {
	fv := firstVisit(t)
	st, store := invariantState(t, nil, nil, fv-1)
	_, acc := CheckStateInvariants(st, store, big.Zero(), testStartEpoch)
	require.False(t, acc.IsEmpty())
	require.Contains(t, acc.Messages()[0], "must have exactly one deal op")
}

// A legacy deal sits past its first visit and is rescheduled by cron itself, so any queue
// position is legitimate and the rule does not apply.
func TestDealOpsLegacyDealExempt(t *testing.T) {
	fv := firstVisit(t)
	state := &DealState{
		SectorNumber:     1,
		SectorStartEpoch: testStartEpoch,
		LastUpdatedEpoch: fv,
		SlashEpoch:       EpochUndefined,
	}
	// Queued somewhere unrelated to its first visit, and cron has run past that visit.
	st, store := invariantState(t, []abi.ChainEpoch{fv + 7*DealUpdatesInterval}, state, fv+1)
	requireDealOpsOK(t, st, store, testStartEpoch+DealUpdatesInterval)
}

// A deal whose first visit is exactly its start epoch, after cron has run there. The old rule
// ("every proposal with start_epoch >= current epoch must have a deal-ops entry") flagged this
// state; the deal has been visited and dropped, so there is nothing left to require.
func TestDealOpsFirstVisitAtStartEpochAfterCron(t *testing.T) {
	// Deal ID 0 has offset 0, and start 0 quantizes to 0, so the first visit is the start epoch.
	store := testStore(t)
	st, err := ConstructState(store)
	require.NoError(t, err)
	st.LastCron = 0
	st.NextID = 1

	proposal := DealProposal{
		PieceCID:             cid.MustParse("baga6ea4seaaqa"),
		PieceSize:            abi.PaddedPieceSize(2048),
		Client:               testAddr(t, 100),
		Provider:             testAddr(t, 200),
		Label:                EmptyDealLabel,
		StartEpoch:           0,
		EndEpoch:             DealMinDuration,
		StoragePricePerEpoch: big.Zero(),
		ProviderCollateral:   big.Zero(),
		ClientCollateral:     big.Zero(),
	}
	require.Equal(t, abi.ChainEpoch(0), NextUpdateEpoch(0, DealUpdatesInterval, proposal.StartEpoch))

	proposals, err := adt.AsArray(store, st.Proposals, ProposalsAmtBitwidth)
	require.NoError(t, err)
	require.NoError(t, proposals.Set(0, &proposal))
	st.Proposals, err = proposals.Root()
	require.NoError(t, err)

	states, err := adt.AsArray(store, st.States, StatesAmtBitwidth)
	require.NoError(t, err)
	require.NoError(t, states.Set(0, &DealState{
		SectorNumber:     1,
		SectorStartEpoch: 0,
		LastUpdatedEpoch: 0,
		SlashEpoch:       EpochUndefined,
	}))
	st.States, err = states.Root()
	require.NoError(t, err)
	require.NoError(t, putProviderSector(store, st, proposal.Provider, 1, 0))

	// Visited and dropped: nothing queued.
	_, acc := CheckStateInvariants(st, store, big.Zero(), abi.ChainEpoch(1))
	require.True(t, acc.IsEmpty(), "unexpected invariant failures: %v", acc.Messages())
}
