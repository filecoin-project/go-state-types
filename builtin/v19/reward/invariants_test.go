package reward

import (
	"bytes"
	"context"
	"io"
	"sort"
	"strings"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v19/util/adt"
	"github.com/filecoin-project/go-state-types/test_util"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"
)

func validInvariantState(t *testing.T) (*State, *StreamsState, adt.Store) {
	t.Helper()
	store := adt.WrapStore(context.Background(), cbor.NewCborStore(test_util.NewBlockStoreInMemory()))
	st, err := ConstructState(store, big.Zero())
	require.NoError(t, err)
	st.EffectiveBaselinePower = big.Zero()

	pct := Denom / 100
	streams := &StreamsState{
		Streams: []Stream{
			{ID: 1, Weight: weight(95*pct, 0, st.Epoch, 50*pct, 95*pct)},
			{
				ID:     2,
				Weight: weight(5*pct, 0, st.Epoch, 5*pct, 10*pct),
				Distribution: &ExplicitDistribution{
					Writer: idAddress(t, 100),
					Shares: []RecipientShare{{Recipient: idAddress(t, 101), Share: Denom}},
				},
			},
		},
	}
	st.Accrued = []StreamAccrual{{ID: 2, Amount: big.Zero()}}
	putStreamsState(t, store, st, streams)
	return st, streams, store
}

func putStreamsState(t *testing.T, store adt.Store, st *State, streams *StreamsState) {
	t.Helper()
	root, err := store.Put(store.Context(), streams)
	require.NoError(t, err)
	st.StreamsRoot = root
}

type cborMarshaler interface {
	MarshalCBOR(io.Writer) error
}

func marshalPayload(t *testing.T, value cborMarshaler) []byte {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, value.MarshalCBOR(&buf))
	return buf.Bytes()
}

func requireInvariantMessage(t *testing.T, acc *builtin.MessageAccumulator, expected string) {
	t.Helper()
	require.Contains(t, strings.Join(acc.Messages(), "\n"), expected)
}

func amountRows(t *testing.T, first, count uint64) []RecipientAmount {
	t.Helper()
	rows := make([]RecipientAmount, count)
	for i := range rows {
		rows[i] = RecipientAmount{Recipient: idAddress(t, first+uint64(i)), Amount: abi.NewTokenAmount(1)}
	}
	sort.Slice(rows, func(i, j int) bool {
		return bytes.Compare(rows[i].Recipient.Bytes(), rows[j].Recipient.Bytes()) < 0
	})
	return rows
}

func TestCheckStateInvariants(t *testing.T) {
	st, _, store := validInvariantState(t)
	summary, acc := CheckStateInvariants(st, store, st.Epoch-1, StorageMiningAllocationCheck)
	require.Empty(t, acc.Messages())
	require.Equal(t, &StateSummary{StreamCount: 2}, summary)
}

func TestCheckStateInvariantsRejectsTopLevelCorruption(t *testing.T) {
	testCases := []struct {
		name     string
		mutate   func(*testing.T, *State)
		expected string
	}{
		{
			name: "negative minted total",
			mutate: func(_ *testing.T, st *State) {
				st.TotalMintedReward = abi.NewTokenAmount(-1)
			},
			expected: "total minted reward is negative",
		},
		{
			name: "decomposition exceeds total",
			mutate: func(_ *testing.T, st *State) {
				st.TotalMintedReward = abi.NewTokenAmount(1)
				st.TotalBurnMinted = abi.NewTokenAmount(1)
				st.TotalExplicitMinted = abi.NewTokenAmount(1)
			},
			expected: "burn 1 + explicit 1 exceeds total minted 1",
		},
		{
			name: "negative timelock",
			mutate: func(_ *testing.T, st *State) {
				st.SWATimelockEpochs = -1
			},
			expected: "SWA timelock is negative",
		},
		{
			name: "non-ID SWA actor",
			mutate: func(t *testing.T, st *State) {
				addr, err := address.NewDelegatedAddress(10, bytes.Repeat([]byte{1}, 20))
				require.NoError(t, err)
				st.SWAActor = addr
			},
			expected: "is not an ID address",
		},
		{
			name: "duplicate accrual ID",
			mutate: func(_ *testing.T, st *State) {
				st.Accrued = []StreamAccrual{{ID: 2, Amount: big.Zero()}, {ID: 2, Amount: big.Zero()}}
			},
			expected: "accrual rows are not strictly ordered",
		},
		{
			name: "negative accrual",
			mutate: func(_ *testing.T, st *State) {
				st.Accrued[0].Amount = abi.NewTokenAmount(-1)
			},
			expected: "accrual for stream 2 is negative",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			st, _, store := validInvariantState(t)
			tc.mutate(t, st)
			_, acc := CheckStateInvariants(st, store, st.Epoch-1, StorageMiningAllocationCheck)
			requireInvariantMessage(t, acc, tc.expected)
		})
	}
}

func TestCheckStateInvariantsRejectsStreamCorruption(t *testing.T) {
	testCases := []struct {
		name     string
		mutate   func(*testing.T, *State, *StreamsState)
		expected string
	}{
		{
			name: "invalid weight record",
			mutate: func(_ *testing.T, _ *State, streams *StreamsState) {
				streams.Streams[0].Weight.Floor = streams.Streams[0].Weight.Cap + 1
			},
			expected: "weight floor exceeds cap",
		},
		{
			name: "future aggregate weight overflow",
			mutate: func(_ *testing.T, _ *State, streams *StreamsState) {
				streams.Streams[1].Weight.Slope = 1
			},
			expected: "stream weights exceed DENOM",
		},
		{
			name: "too many streams",
			mutate: func(_ *testing.T, _ *State, streams *StreamsState) {
				for id := StreamID(3); id <= 9; id++ {
					streams.Streams = append(streams.Streams, Stream{ID: id, Weight: weight(0, 0, 0, 0, 0)})
				}
			},
			expected: "stream count exceeds maximum",
		},
		{
			name: "unordered stream IDs",
			mutate: func(_ *testing.T, _ *State, streams *StreamsState) {
				streams.Streams[0], streams.Streams[1] = streams.Streams[1], streams.Streams[0]
			},
			expected: "stream IDs are not ordered",
		},
		{
			name: "multiple implicit streams",
			mutate: func(_ *testing.T, _ *State, streams *StreamsState) {
				streams.Streams = append(streams.Streams, Stream{ID: 3, Weight: weight(0, 0, 0, 0, 0)})
			},
			expected: "multiple implicit streams",
		},
		{
			name: "invalid share total",
			mutate: func(_ *testing.T, _ *State, streams *StreamsState) {
				streams.Streams[1].Distribution.Shares[0].Share--
			},
			expected: "shares sum to",
		},
		{
			name: "zero payable amount",
			mutate: func(t *testing.T, _ *State, streams *StreamsState) {
				streams.Streams[1].Distribution.Payable = []RecipientAmount{{Recipient: idAddress(t, 102), Amount: big.Zero()}}
			},
			expected: "payable amount is not positive",
		},
		{
			name: "claimed recipient absent from shares",
			mutate: func(t *testing.T, st *State, streams *StreamsState) {
				streams.Streams[1].Distribution.ClaimedPeriod = []RecipientAmount{{Recipient: idAddress(t, 102), Amount: abi.NewTokenAmount(1)}}
				st.Accrued[0].Amount = abi.NewTokenAmount(10)
			},
			expected: "claimed-period recipient is absent from shares",
		},
		{
			name: "accrual ID mismatch",
			mutate: func(_ *testing.T, st *State, _ *StreamsState) {
				st.Accrued[0].ID = 3
			},
			expected: "accrual IDs do not match live explicit streams",
		},
		{
			name: "live tombstone ID",
			mutate: func(t *testing.T, _ *State, streams *StreamsState) {
				streams.Tombstones = []Tombstone{{ID: 2, Payable: []RecipientAmount{{Recipient: idAddress(t, 102), Amount: abi.NewTokenAmount(1)}}}}
			},
			expected: "stream ID is live and tombstoned",
		},
		{
			name: "empty tombstone",
			mutate: func(_ *testing.T, _ *State, streams *StreamsState) {
				streams.Tombstones = []Tombstone{{ID: 3}}
			},
			expected: "tombstone 3 is empty",
		},
		{
			name: "malformed pending payload",
			mutate: func(_ *testing.T, _ *State, streams *StreamsState) {
				streams.PendingWrites = []PendingWrite{{Op: PendingWriteOpSetWeightRecords, Payload: []byte{0xff}, EffectiveEpoch: 1}}
			},
			expected: "invalid streams state",
		},
		{
			name: "non-canonical pending target",
			mutate: func(_ *testing.T, _ *State, streams *StreamsState) {
				streams.PendingWrites = []PendingWrite{{Op: PendingWriteOpRemoveStream, Payload: []byte{0x80}, EffectiveEpoch: 1}}
			},
			expected: "non-canonical stream ID",
		},
		{
			name: "duplicate pending slot",
			mutate: func(_ *testing.T, _ *State, streams *StreamsState) {
				id := StreamID(2)
				streams.PendingWrites = []PendingWrite{
					{ID: &id, Op: PendingWriteOpRemoveStream, Payload: []byte{0x80}, EffectiveEpoch: 1},
					{ID: &id, Op: PendingWriteOpRemoveStream, Payload: []byte{0x80}, EffectiveEpoch: 2},
				}
			},
			expected: "duplicate pending slot",
		},
		{
			name: "pending registration reuses live ID",
			mutate: func(t *testing.T, _ *State, streams *StreamsState) {
				id := StreamID(2)
				payload := &RegisterStreamPayload{Weight: weight(0, 0, 1, 0, 0)}
				streams.PendingWrites = []PendingWrite{{ID: &id, Op: PendingWriteOpRegisterStream, Payload: marshalPayload(t, payload), EffectiveEpoch: 1}}
			},
			expected: "pending registration reuses stream ID 2",
		},
		{
			name: "tombstone reservations exceed cap",
			mutate: func(t *testing.T, _ *State, streams *StreamsState) {
				streams.Tombstones = []Tombstone{{ID: 3, Payable: amountRows(t, 1_000, 200)}}
				id := StreamID(2)
				streams.PendingWrites = []PendingWrite{{ID: &id, Op: PendingWriteOpRemoveStream, Payload: []byte{0x80}, EffectiveEpoch: 1}}
			},
			expected: "tombstone row reservation 264 exceeds maximum 256",
		},
		{
			name: "liability exceeds total explicit minted",
			mutate: func(_ *testing.T, st *State, _ *StreamsState) {
				st.Accrued[0].Amount = abi.NewTokenAmount(10)
			},
			expected: "liabilities 10 exceed total explicit minted 0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			st, streams, store := validInvariantState(t)
			tc.mutate(t, st, streams)
			putStreamsState(t, store, st, streams)
			_, acc := CheckStateInvariants(st, store, st.Epoch-1, StorageMiningAllocationCheck)
			requireInvariantMessage(t, acc, tc.expected)
		})
	}
}

func TestCheckStateInvariantsRequiresLiabilityBalance(t *testing.T) {
	st, streams, store := validInvariantState(t)
	st.TotalMintedReward = StorageMiningAllocationCheck
	st.TotalExplicitMinted = abi.NewTokenAmount(10)
	st.Accrued[0].Amount = abi.NewTokenAmount(10)
	putStreamsState(t, store, st, streams)

	_, acc := CheckStateInvariants(st, store, st.Epoch-1, big.Zero())
	requireInvariantMessage(t, acc, "reward balance 0 does not cover explicit-stream liabilities 10")
}
