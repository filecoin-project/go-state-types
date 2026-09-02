package reward

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v19/util/smoothing"
	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
)

// Tests to match with Rust fil_actor_reward::serialization

func idAddress(t *testing.T, id uint64) address.Address {
	t.Helper()
	addr, err := address.NewIDAddress(id)
	require.NoError(t, err)
	return addr
}

func delegatedAddress(t *testing.T) address.Address {
	t.Helper()
	addr, err := address.NewDelegatedAddress(10, bytes.Repeat([]byte{0x11}, 20))
	require.NoError(t, err)
	return addr
}

func emptyStreamsRoot(t *testing.T) cid.Cid {
	t.Helper()
	data, err := hex.DecodeString("0171a0e40220d63b11132be58f8f498e5f8c46c4d26b89675b443ff1c47f1e7e3d3cb8d2dcaa")
	require.NoError(t, err)
	root, err := cid.Cast(data)
	require.NoError(t, err)
	return root
}

func weight(vStart uint64, slope int64, tStart abi.ChainEpoch, floor, cap uint64) WeightRecord {
	return WeightRecord{VStart: vStart, Slope: slope, TStart: tStart, Floor: floor, Cap: cap}
}

func streamIDPtr(id StreamID) *StreamID {
	return &id
}

func TestDeferredPayloadRustVectors(t *testing.T) {
	testCases := []struct {
		name     string
		hex      string
		decoded  cborUnmarshaler
		expected any
	}{
		{
			name:     "weight_record",
			hex:      "850220030004",
			decoded:  &WeightRecord{},
			expected: &WeightRecord{VStart: 2, Slope: -1, TStart: 3, Floor: 0, Cap: 4},
		},
		{
			name: "weight_records_payload",
			hex:  "818282018502200300048205850800070809",
			decoded: &SetWeightRecordsParams{
				Updates: []WeightRecordUpdate{},
			},
			expected: &SetWeightRecordsParams{Updates: []WeightRecordUpdate{
				{ID: 1, Weight: weight(2, -1, 3, 0, 4)},
				{ID: 5, Weight: weight(8, 0, 7, 8, 9)},
			}},
		},
		{
			name:     "register_stream_payload_implicit",
			hex:      "82850220030004f6",
			decoded:  &RegisterStreamPayload{},
			expected: &RegisterStreamPayload{Weight: weight(2, -1, 3, 0, 4)},
		},
		{
			name:    "register_stream_payload_explicit",
			hex:     "82850220030004824300c80181824200651b0de0b6b3a7640000",
			decoded: &RegisterStreamPayload{},
			expected: &RegisterStreamPayload{
				Weight: weight(2, -1, 3, 0, 4),
				Distribution: &DistributionInit{
					Writer: idAddress(t, 200),
					Shares: []RecipientShare{{Recipient: idAddress(t, 101), Share: Denom}},
				},
			},
		},
		{
			name:     "set_distribution_payload",
			hex:      "81420001",
			decoded:  &SetDistributionPayload{},
			expected: &SetDistributionPayload{Writer: idAddress(t, 1)},
		},
		{
			name:     "pending_write_nil_id",
			hex:      "84f60041810a",
			decoded:  &PendingWrite{},
			expected: &PendingWrite{Op: PendingWriteOpSetWeightRecords, Payload: []byte{0x81}, EffectiveEpoch: 10},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := hex.DecodeString(tc.hex)
			require.NoError(t, err)
			require.NoError(t, unmarshalPendingPayload(encoded, tc.decoded))
			require.Equal(t, tc.expected, tc.decoded)
		})
	}
}

// Go exposes both transparent Rust wrappers as *abi.StoragePower in Methods. A nil
// pointer is None; non-nil pointers encode the inner bigint without a tuple wrapper.
func testStoragePowerParams(t *testing.T) {
	power := func(value int64) *abi.StoragePower {
		power := abi.NewStoragePower(value)
		return &power
	}
	testCases := []struct {
		name   string
		params *abi.StoragePower
		hex    string
	}{
		{name: "none", params: nil, hex: "f6"},
		{name: "zero", params: power(0), hex: "40"},
		{name: "255", params: power(255), hex: "4200ff"},
		{name: "256", params: power(256), hex: "43000100"},
		{name: "negative_255", params: power(-255), hex: "4201ff"},
		{name: "negative_256", params: power(-256), hex: "43010100"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)
			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))

			if tc.params != nil {
				var decoded abi.StoragePower
				req.NoError(decoded.UnmarshalCBOR(&buf))
				req.Equal(*tc.params, decoded)
			}
		})
	}
}

func TestSerializationConstructorParams(t *testing.T) {
	_, ok := Methods[1].Method.(func(*abi.StoragePower) *abi.EmptyValue)
	require.True(t, ok)
	testStoragePowerParams(t)
}

func TestSerializationUpdateNetworkKPIParams(t *testing.T) {
	_, ok := Methods[4].Method.(func(*abi.StoragePower) *abi.EmptyValue)
	require.True(t, ok)
	testStoragePowerParams(t)
}

func TestSerializationAwardBlockRewardParams(t *testing.T) {
	testCases := []struct {
		name   string
		params AwardBlockRewardParams
		hex    string
	}{
		{
			name: "zero",
			params: AwardBlockRewardParams{
				Miner: idAddress(t, 100), Penalty: big.Zero(), GasReward: big.Zero(), WinCount: 0,
			},
			hex: "84420064404000",
		},
		{
			name: "negative_win_count",
			params: AwardBlockRewardParams{
				Miner:     delegatedAddress(t),
				Penalty:   abi.NewTokenAmount(255),
				GasReward: abi.NewTokenAmount(256),
				WinCount:  -1,
			},
			hex: "8456040a11111111111111111111111111111111111111114200ff4300010020",
		},
		{
			name: "maximum_win_count",
			params: AwardBlockRewardParams{
				Miner:     idAddress(t, 100),
				Penalty:   abi.NewTokenAmount(256),
				GasReward: abi.NewTokenAmount(255),
				WinCount:  int64(^uint64(0) >> 1),
			},
			hex: "84420064430001004200ff1b7fffffffffffffff",
		},
		{
			name: "minimum_win_count",
			params: AwardBlockRewardParams{
				Miner:     delegatedAddress(t),
				Penalty:   abi.NewTokenAmount(1),
				GasReward: abi.NewTokenAmount(1),
				WinCount:  -int64(^uint64(0)>>1) - 1,
			},
			hex: "8456040a11111111111111111111111111111111111111114200014200013b7fffffffffffffff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)
			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))

			var decoded AwardBlockRewardParams
			req.NoError(decoded.UnmarshalCBOR(&buf))
			req.Equal(tc.params, decoded)
		})
	}
}

func TestSerializationThisEpochRewardReturn(t *testing.T) {
	params := ThisEpochRewardReturn{
		ThisEpochRewardSmoothed: smoothing.NewEstimate(big.NewInt(1), big.NewInt(-2)),
		ThisEpochBaselinePower:  abi.NewStoragePower(256),
	}
	const expected = "8282520001000000000000000000000000000000005201020000000000000000000000000000000043000100"

	var buf bytes.Buffer
	require.NoError(t, params.MarshalCBOR(&buf))
	require.Equal(t, expected, hex.EncodeToString(buf.Bytes()))

	var decoded ThisEpochRewardReturn
	require.NoError(t, decoded.UnmarshalCBOR(&buf))
	require.Equal(t, params, decoded)
}

func TestSerializationStreamsState(t *testing.T) {
	testCases := []struct {
		name   string
		params StreamsState
		hex    string
	}{
		{
			name:   "empty",
			params: StreamsState{},
			hex:    "83808080",
		},
		{
			name: "populated",
			params: StreamsState{
				Streams: []Stream{
					{ID: 1, Weight: weight(0, 0, 0, 0, 0)},
					{
						ID:     2,
						Weight: weight(5, -2, 1, 4, 5),
						Distribution: &ExplicitDistribution{
							Writer:        idAddress(t, 100),
							Shares:        []RecipientShare{{Recipient: idAddress(t, 101), Share: 6}},
							Payable:       []RecipientAmount{{Recipient: idAddress(t, 102), Amount: abi.NewTokenAmount(7)}},
							ClaimedPeriod: []RecipientAmount{{Recipient: idAddress(t, 103), Amount: abi.NewTokenAmount(8)}},
						},
					},
				},
				Tombstones:    []Tombstone{{ID: 3, Payable: []RecipientAmount{{Recipient: idAddress(t, 104), Amount: abi.NewTokenAmount(9)}}}},
				PendingWrites: []PendingWrite{{ID: streamIDPtr(4), Op: PendingWriteOpRegisterStream, Payload: []byte{0x81, 0x01}, EffectiveEpoch: 10}},
			},
			hex: "83828301850000000000f6830285052101040584420064818242006506818242006642000781824200674200088182038182420068420009818404024281010a",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt StreamsState
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.params, rt)
		})
	}
}

func TestSerializationState(t *testing.T) {
	root := emptyStreamsRoot(t)
	testCases := []struct {
		name   string
		params State
		hex    string
	}{
		{
			name: "empty",
			params: State{
				CumsumBaseline:          big.Zero(),
				CumsumRealized:          big.Zero(),
				EffectiveBaselinePower:  big.Zero(),
				ThisEpochReward:         big.Zero(),
				ThisEpochRewardSmoothed: smoothing.NewEstimate(big.Zero(), big.Zero()),
				ThisEpochBaselinePower:  big.Zero(),
				TotalMintedReward:       big.Zero(),
				TotalBurnMinted:         big.Zero(),
				TotalExplicitMinted:     big.Zero(),
				SWAActor:                idAddress(t, 0),
				StreamsRoot:             root,
			},
			hex: "8f404000404082404040004040408000420000d82a5827000171a0e40220d63b11132be58f8f498e5f8c46c4d26b89675b443ff1c47f1e7e3d3cb8d2dcaa",
		},
		{
			name: "populated",
			params: State{
				CumsumBaseline:          big.NewInt(1),
				CumsumRealized:          big.NewInt(2),
				EffectiveNetworkTime:    3,
				EffectiveBaselinePower:  big.NewInt(4),
				ThisEpochReward:         abi.NewTokenAmount(5),
				ThisEpochRewardSmoothed: smoothing.NewEstimate(big.Zero(), big.Zero()),
				ThisEpochBaselinePower:  big.NewInt(6),
				Epoch:                   7,
				TotalMintedReward:       abi.NewTokenAmount(8),
				TotalBurnMinted:         abi.NewTokenAmount(9),
				TotalExplicitMinted:     abi.NewTokenAmount(10),
				Accrued:                 []StreamAccrual{{ID: 2, Amount: abi.NewTokenAmount(11)}},
				SWATimelockEpochs:       13,
				SWAActor:                idAddress(t, 1001),
				StreamsRoot:             root,
			},
			hex: "8f420001420002034200044200058240404200060742000842000942000a81820242000b0d4300e907d82a5827000171a0e40220d63b11132be58f8f498e5f8c46c4d26b89675b443ff1c47f1e7e3d3cb8d2dcaa",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt State
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.params, rt)
		})
	}
}

func TestSerializationSetWeightRecordsParams(t *testing.T) {
	updates := []WeightRecordUpdate{
		{ID: 23, Weight: weight(23, -24, 23, 0, 23)},
		{ID: 24, Weight: weight(24, -25, 24, 0, 24)},
		{ID: 255, Weight: weight(255, -256, 255, 0, 255)},
		{ID: 256, Weight: weight(256, -257, 256, 0, 256)},
		{ID: 65_535, Weight: weight(65_535, -65_536, 65_535, 0, 65_535)},
		{ID: 65_536, Weight: weight(65_536, -65_537, 65_536, 0, 65_536)},
		{ID: StreamID(^uint32(0)), Weight: weight(uint64(^uint32(0)), -(1 << 32), abi.ChainEpoch(^uint32(0)), 0, uint64(^uint32(0)))},
		{ID: 1 << 32, Weight: weight(1<<32, -((1 << 32) + 1), 1<<32, 0, 1<<32)},
	}
	testCases := []struct {
		name   string
		params SetWeightRecordsParams
		hex    string
	}{
		{name: "empty", params: SetWeightRecordsParams{}, hex: "8180"},
		{
			name:   "integer_boundaries",
			params: SetWeightRecordsParams{Updates: updates},
			hex:    "81888217851737170017821818851818381818180018188218ff8518ff38ff18ff0018ff8219010085190100390100190100001901008219ffff8519ffff39ffff19ffff0019ffff821a00010000851a000100003a000100001a00010000001a00010000821affffffff851affffffff3affffffff1affffffff001affffffff821b0000000100000000851b00000001000000003b00000001000000001b0000000100000000001b0000000100000000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt SetWeightRecordsParams
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.params, rt)
		})
	}
}

func TestSerializationStepWeightRecordsParams(t *testing.T) {
	testCases := []struct {
		name   string
		params StepWeightRecordsParams
		hex    string
	}{
		{name: "empty", params: StepWeightRecordsParams{}, hex: "8180"},
		{
			name:   "integer_boundaries",
			params: StepWeightRecordsParams{Updates: []WeightRecordUpdate{{ID: 1 << 32, Weight: weight(1<<32, -((1 << 32) + 1), 65_536, 256, Denom)}}},
			hex:    "8181821b0000000100000000851b00000001000000003b00000001000000001a000100001901001b0de0b6b3a7640000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt StepWeightRecordsParams
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.params, rt)
		})
	}
}

func TestSerializationRegisterStreamParams(t *testing.T) {
	testCases := []struct {
		name   string
		params RegisterStreamParams
		hex    string
	}{
		{
			name:   "implicit",
			params: RegisterStreamParams{ID: 24, Weight: weight(24, -24, 256, 0, 65_536), ActivationEpoch: 1 << 32},
			hex:    "84181885181837190100001a00010000f61b0000000100000000",
		},
		{
			name: "explicit",
			params: RegisterStreamParams{
				ID:     1 << 32,
				Weight: weight(1<<32, -(1 << 32), 65_536, 256, Denom),
				Distribution: &DistributionInit{
					Writer: delegatedAddress(t),
					Shares: []RecipientShare{{Recipient: idAddress(t, 1<<32), Share: Denom}},
				},
				ActivationEpoch: 65_536,
			},
			hex: "841b0000000100000000851b00000001000000003affffffff1a000100001901001b0de0b6b3a76400008256040a11111111111111111111111111111111111111118182460080808080101b0de0b6b3a76400001a00010000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt RegisterStreamParams
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.params, rt)
		})
	}
}

func TestSerializationRemoveStreamParams(t *testing.T) {
	testCases := []struct {
		name   string
		params RemoveStreamParams
		hex    string
	}{
		{name: "id_24", params: RemoveStreamParams{ID: 24}, hex: "811818"},
		{name: "id_256", params: RemoveStreamParams{ID: 256}, hex: "81190100"},
		{name: "id_65536", params: RemoveStreamParams{ID: 65_536}, hex: "811a00010000"},
		{name: "id_4294967296", params: RemoveStreamParams{ID: 1 << 32}, hex: "811b0000000100000000"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt RemoveStreamParams
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.params, rt)
		})
	}
}

func TestSerializationSetDistributionParams(t *testing.T) {
	testCases := []struct {
		name   string
		params SetDistributionParams
		hex    string
	}{
		{name: "id_writer", params: SetDistributionParams{ID: 24, Writer: idAddress(t, 1<<32)}, hex: "82181846008080808010"},
		{name: "delegated_writer", params: SetDistributionParams{ID: 256, Writer: delegatedAddress(t)}, hex: "8219010056040a1111111111111111111111111111111111111111"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt SetDistributionParams
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.params, rt)
		})
	}
}

func TestSerializationSetSharesParams(t *testing.T) {
	maxShares := make([]RecipientShare, MaxRecipients)
	for i := range maxShares {
		maxShares[i] = RecipientShare{Recipient: idAddress(t, 1000+uint64(i)), Share: Denom / MaxRecipients}
	}
	testCases := []struct {
		name   string
		params SetSharesParams
		hex    string
	}{
		{name: "empty", params: SetSharesParams{ID: 24}, hex: "82181880"},
		{
			name: "integer_boundaries",
			params: SetSharesParams{
				ID: 256,
				Shares: []RecipientShare{
					{Recipient: idAddress(t, 24), Share: 24},
					{Recipient: idAddress(t, 256), Share: 256},
					{Recipient: idAddress(t, 65_536), Share: 65_536},
					{Recipient: idAddress(t, 1<<32), Share: 1 << 32},
				},
			},
			hex: "821901008482420018181882430080021901008244008080041a0001000082460080808080101b0000000100000000",
		},
		{
			name:   "max_recipients",
			params: SetSharesParams{ID: 65_536, Shares: maxShares},
			hex:    "821a000100009840824300e8071b003782dace9d9000824300e9071b003782dace9d9000824300ea071b003782dace9d9000824300eb071b003782dace9d9000824300ec071b003782dace9d9000824300ed071b003782dace9d9000824300ee071b003782dace9d9000824300ef071b003782dace9d9000824300f0071b003782dace9d9000824300f1071b003782dace9d9000824300f2071b003782dace9d9000824300f3071b003782dace9d9000824300f4071b003782dace9d9000824300f5071b003782dace9d9000824300f6071b003782dace9d9000824300f7071b003782dace9d9000824300f8071b003782dace9d9000824300f9071b003782dace9d9000824300fa071b003782dace9d9000824300fb071b003782dace9d9000824300fc071b003782dace9d9000824300fd071b003782dace9d9000824300fe071b003782dace9d9000824300ff071b003782dace9d900082430080081b003782dace9d900082430081081b003782dace9d900082430082081b003782dace9d900082430083081b003782dace9d900082430084081b003782dace9d900082430085081b003782dace9d900082430086081b003782dace9d900082430087081b003782dace9d900082430088081b003782dace9d900082430089081b003782dace9d90008243008a081b003782dace9d90008243008b081b003782dace9d90008243008c081b003782dace9d90008243008d081b003782dace9d90008243008e081b003782dace9d90008243008f081b003782dace9d900082430090081b003782dace9d900082430091081b003782dace9d900082430092081b003782dace9d900082430093081b003782dace9d900082430094081b003782dace9d900082430095081b003782dace9d900082430096081b003782dace9d900082430097081b003782dace9d900082430098081b003782dace9d900082430099081b003782dace9d90008243009a081b003782dace9d90008243009b081b003782dace9d90008243009c081b003782dace9d90008243009d081b003782dace9d90008243009e081b003782dace9d90008243009f081b003782dace9d9000824300a0081b003782dace9d9000824300a1081b003782dace9d9000824300a2081b003782dace9d9000824300a3081b003782dace9d9000824300a4081b003782dace9d9000824300a5081b003782dace9d9000824300a6081b003782dace9d9000824300a7081b003782dace9d9000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt SetSharesParams
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.params, rt)
		})
	}
}

func TestSerializationCancelPendingParams(t *testing.T) {
	testCases := []struct {
		name   string
		params CancelPendingParams
		hex    string
	}{
		{name: "set_weights", params: CancelPendingParams{Op: PendingWriteOpSetWeightRecords}, hex: "82f600"},
		{name: "step_weights", params: CancelPendingParams{Op: PendingWriteOpStepWeightRecords}, hex: "82f601"},
		{name: "register", params: CancelPendingParams{ID: streamIDPtr(24), Op: PendingWriteOpRegisterStream}, hex: "82181802"},
		{name: "remove", params: CancelPendingParams{ID: streamIDPtr(256), Op: PendingWriteOpRemoveStream}, hex: "8219010003"},
		{name: "set_distribution", params: CancelPendingParams{ID: streamIDPtr(65_536), Op: PendingWriteOpSetDistribution}, hex: "821a0001000004"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt CancelPendingParams
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.params, rt)
		})
	}
}

func TestSerializationClaimParams(t *testing.T) {
	maxWallets := make([]address.Address, MaxRecipients)
	for i := range maxWallets {
		maxWallets[i] = idAddress(t, 1000+uint64(i))
	}
	testCases := []struct {
		name   string
		params ClaimParams
		hex    string
	}{
		{name: "empty_wallets", params: ClaimParams{ID: 1 << 32}, hex: "821b000000010000000080"},
		{
			name:   "mixed_protocols",
			params: ClaimParams{ID: 65_536, Wallets: []address.Address{idAddress(t, 1<<32), delegatedAddress(t)}},
			hex:    "821a00010000824600808080801056040a1111111111111111111111111111111111111111",
		},
		{
			name:   "max_wallets",
			params: ClaimParams{ID: 65_536, Wallets: maxWallets},
			hex:    "821a0001000098404300e8074300e9074300ea074300eb074300ec074300ed074300ee074300ef074300f0074300f1074300f2074300f3074300f4074300f5074300f6074300f7074300f8074300f9074300fa074300fb074300fc074300fd074300fe074300ff074300800843008108430082084300830843008408430085084300860843008708430088084300890843008a0843008b0843008c0843008d0843008e0843008f084300900843009108430092084300930843009408430095084300960843009708430098084300990843009a0843009b0843009c0843009d0843009e0843009f084300a0084300a1084300a2084300a3084300a4084300a5084300a6084300a708",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt ClaimParams
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.params, rt)
		})
	}
}

func TestSerializationClaimReturn(t *testing.T) {
	testCases := []struct {
		name   string
		params ClaimReturn
		hex    string
	}{
		{name: "empty", params: ClaimReturn{}, hex: "8180"},
		{
			name: "integer_boundaries",
			params: ClaimReturn{Amounts: []abi.TokenAmount{
				abi.NewTokenAmount(0),
				abi.NewTokenAmount(24),
				abi.NewTokenAmount(256),
				abi.NewTokenAmount(65_536),
				abi.NewTokenAmount(1 << 32),
			}},
			hex: "81854042001843000100440001000046000100000000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt ClaimReturn
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.params, rt)
		})
	}
}

func TestExportedMethodNumbers(t *testing.T) {
	expected := map[string]abi.MethodNum{
		"SetWeightRecords":  3362570548,
		"StepWeightRecords": 3951753085,
		"RegisterStream":    386660827,
		"RemoveStream":      1623858416,
		"SetDistribution":   3872725033,
		"SetShares":         2414422607,
		"CancelPending":     187585191,
		"Claim":             4045527845,
	}
	for name, number := range expected {
		require.Equal(t, number, builtin.MustGenerateFRCMethodNum(name))
		_, ok := Methods[number]
		require.True(t, ok, "method %s is not registered", name)
	}
}
