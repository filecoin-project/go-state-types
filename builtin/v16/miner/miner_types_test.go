package miner

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	cid "github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
)

// Tests to match with Rust fil_actor_miner::serialization

// defaultCid as a Rust Default::default() value
var defaultCid = cid.MustParse("baeaaaaa")

func TestSerializationProveCommitSectorsNIParams(t *testing.T) {
	testCases := []struct {
		params ProveCommitSectorsNIParams
		hex    string
	}{
		{
			params: ProveCommitSectorsNIParams{
				Sectors:                  nil,
				AggregateProof:           nil,
				SealProofType:            abi.RegisteredSealProof_StackedDrg32GiBV1_1,
				AggregateProofType:       abi.RegisteredAggregationProof_SnarkPackV2,
				ProvingDeadline:          2,
				RequireActivationSuccess: false,
			},
			// [[],byte[],8,1,2,false]
			hex: "868040080102f4",
		},
		{
			params: ProveCommitSectorsNIParams{
				Sectors: []SectorNIActivationInfo{{
					SealingNumber: 1,
					SealerID:      2,
					SealedCID:     cid.MustParse("bagboea4seaaqa"),
					SectorNumber:  3,
					SealRandEpoch: 4,
					Expiration:    5,
				}},
				SealProofType:            abi.RegisteredSealProof_StackedDrg32GiBV1_2_Feat_NiPoRep,
				AggregateProof:           []byte{0xde, 0xad, 0xbe, 0xef},
				AggregateProofType:       abi.RegisteredAggregationProof_SnarkPackV2,
				ProvingDeadline:          6,
				RequireActivationSuccess: true,
			},
			// [[[1,2,bagboea4seaaqa,3,4,5]],byte[deadbeef],18,1,6,true]
			hex: "8681860102d82a49000182e2039220010003040544deadbeef120106f5",
		},
		{
			params: ProveCommitSectorsNIParams{
				Sectors: []SectorNIActivationInfo{
					{
						SealingNumber: 1,
						SealerID:      2,
						SealedCID:     cid.MustParse("bagboea4seaaqa"),
						SectorNumber:  3,
						SealRandEpoch: 4,
						Expiration:    5,
					},
					{
						SealingNumber: 6,
						SealerID:      7,
						SealedCID:     cid.MustParse("bagboea4seaaqc"),
						SectorNumber:  8,
						SealRandEpoch: 9,
						Expiration:    10,
					},
				},
				SealProofType:            abi.RegisteredSealProof_StackedDrg32GiBV1_2_Feat_NiPoRep,
				AggregateProof:           []byte{0xde, 0xad, 0xbe, 0xef},
				AggregateProofType:       abi.RegisteredAggregationProof_SnarkPackV2,
				ProvingDeadline:          11,
				RequireActivationSuccess: false,
			},
			// [[[1,2,bagboea4seaaqa,3,4,5],[6,7,bagboea4seaaqc,8,9,10]],byte[deadbeef],18,1,11,false]
			hex: "8682860102d82a49000182e20392200100030405860607d82a49000182e2039220010108090a44deadbeef12010bf4",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt ProveCommitSectorsNIParams
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.params, rt)
		})
	}
}

func TestSectorOnChainInfo(t *testing.T) {
	sectorKey := cid.MustParse("baga6ea4seaaqc")
	zero := big.Zero()

	testCases := []struct {
		sector   SectorOnChainInfo
		readHex  string
		writeHex string
	}{
		{
			sector: SectorOnChainInfo{
				SealProof:             -1,
				SealedCID:             defaultCid,
				DealWeight:            zero,
				VerifiedDealWeight:    zero,
				InitialPledge:         zero,
				ExpectedDayReward:     zero,
				ExpectedStoragePledge: zero,
				ReplacedDayReward:     zero,
				DailyFee:              zero,
			},
			// [0,-1,{"/":"baeaaaaa"},[],0,0,[],[],[],[],[],0,[],null,0,[]]
			readHex: "900020d82a45000100000080000040404040400040f60040",
			// same on write as read
			writeHex: "900020d82a45000100000080000040404040400040f60040",
		},
		{
			sector: SectorOnChainInfo{
				SectorNumber:          1,
				SealProof:             abi.RegisteredSealProof_StackedDrg32GiBV1_1,
				SealedCID:             cid.MustParse("bagboea4seaaqa"),
				DeprecatedDealIDs:     nil,
				Activation:            2,
				Expiration:            3,
				DealWeight:            big.NewInt(4),
				VerifiedDealWeight:    big.NewInt(5),
				InitialPledge:         filWhole(6),
				ExpectedDayReward:     filWhole(7),
				ExpectedStoragePledge: filWhole(8),
				PowerBaseEpoch:        9,
				ReplacedDayReward:     filWhole(10),
				SectorKeyCID:          nil,
				Flags:                 0,
				DailyFee:              filWhole(11),
			},
			// '[1,8,{"/":"bagboea4seaaqa"},[],2,3,[AAQ],[AAU],[AFNESDXsWAAA],[AGEk/umTvAAA],[AG8FtZ07IAAA],9,[AIrHIwSJ6AAA],null,0,[AJin2bgxTAAA]]'
			readHex: "900108d82a49000182e20392200100800203420004420005490053444835ec58000049006124fee993bc000049006f05b59d3b2000000949008ac7230489e80000f600490098a7d9b8314c0000",
			// same on write as read
			writeHex: "900108d82a49000182e20392200100800203420004420005490053444835ec58000049006124fee993bc000049006f05b59d3b2000000949008ac7230489e80000f600490098a7d9b8314c0000",
		},
		{
			sector: SectorOnChainInfo{
				SectorNumber:          1,
				SealProof:             abi.RegisteredSealProof_StackedDrg32GiBV1_1,
				SealedCID:             cid.MustParse("bagboea4seaaqa"),
				DeprecatedDealIDs:     nil,
				Activation:            2,
				Expiration:            3,
				DealWeight:            big.NewInt(4),
				VerifiedDealWeight:    big.NewInt(5),
				InitialPledge:         filWhole(6),
				ExpectedDayReward:     filWhole(7),
				ExpectedStoragePledge: filWhole(8),
				PowerBaseEpoch:        9,
				ReplacedDayReward:     filWhole(10),
				SectorKeyCID:          &sectorKey,
				Flags:                 SIMPLE_QA_POWER,
				DailyFee:              filWhole(11),
			},
			// [1,8,{"/":"bagboea4seaaqa"},[],2,3,[AAQ],[AAU],[AFNESDXsWAAA],[AGEk/umTvAAA],[AG8FtZ07IAAA],9,[AIrHIwSJ6AAA],{"/":"baga6ea4seaaqc"},1,[AJin2bgxTAAA]]
			readHex: "900108d82a49000182e20392200100800203420004420005490053444835ec58000049006124fee993bc000049006f05b59d3b2000000949008ac7230489e80000d82a49000181e2039220010101490098a7d9b8314c0000",
			// same on write as read
			writeHex: "900108d82a49000182e20392200100800203420004420005490053444835ec58000049006124fee993bc000049006f05b59d3b2000000949008ac7230489e80000d82a49000181e2039220010101490098a7d9b8314c0000",
		},
		{
			// old format stored on chain but materialised as the new format with a default value at the end
			sector: SectorOnChainInfo{
				SectorNumber:          1,
				SealProof:             abi.RegisteredSealProof_StackedDrg64GiBV1_1,
				SealedCID:             cid.MustParse("bagboea4seaaqa"),
				DeprecatedDealIDs:     nil,
				Activation:            2,
				Expiration:            3,
				DealWeight:            big.NewInt(4),
				VerifiedDealWeight:    big.NewInt(5),
				InitialPledge:         filWhole(6),
				ExpectedDayReward:     filWhole(7),
				ExpectedStoragePledge: filWhole(8),
				PowerBaseEpoch:        9,
				ReplacedDayReward:     filWhole(10),
				SectorKeyCID:          nil,
				Flags:                 SIMPLE_QA_POWER,
				DailyFee:              big.Int{}, // default, not present in the binary
			},
			// [1,9,{"/":"bagboea4seaaqa"},[],2,3,[AAQ],[AAU],[AFNESDXsWAAA],[AGEk/umTvAAA],[AG8FtZ07IAAA],9,[AIrHIwSJ6AAA],null,1]
			readHex: "8f0109d82a49000182e20392200100800203420004420005490053444835ec58000049006124fee993bc000049006f05b59d3b2000000949008ac7230489e80000f601",
			// extra field at the end on write, zero BigInt (bytes) for daily_fee
			// [1,9,{"/":"bagboea4seaaqa"},[],2,3,[AAQ],[AAU],[AFNESDXsWAAA],[AGEk/umTvAAA],[AG8FtZ07IAAA],9,[AIrHIwSJ6AAA],null,1,[]]
			writeHex: "900109d82a49000182e20392200100800203420004420005490053444835ec58000049006124fee993bc000049006f05b59d3b2000000949008ac7230489e80000f60140",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			req := require.New(t)

			// write
			var buf bytes.Buffer
			req.NoError(tc.sector.MarshalCBOR(&buf))
			req.Equal(tc.writeHex, hex.EncodeToString(buf.Bytes()))

			// read
			byts, err := hex.DecodeString(tc.readHex)
			req.NoError(err)
			var rt SectorOnChainInfo
			req.NoError(rt.UnmarshalCBOR(bytes.NewReader(byts)))
			req.Equal(tc.sector, rt)
		})
	}
}

func TestExpirationSet(t *testing.T) {
	zero := big.Zero()

	testCases := []struct {
		set      ExpirationSet
		readHex  string
		writeHex string
	}{
		{
			set: ExpirationSet{
				OnTimeSectors: bitfield.New(),
				EarlySectors:  bitfield.New(),
				OnTimePledge:  zero,
				ActivePower:   NewPowerPairZero(),
				FaultyPower:   NewPowerPairZero(),
				FeeDeduction:  zero,
			},
			// [[],[],[],[[],[]],[[],[]],[]]
			readHex: "8640404082404082404040",
			// same on write as read
			writeHex: "8640404082404082404040",
		},
		{
			set: ExpirationSet{
				OnTimeSectors: bfrt(0),
				EarlySectors:  bfrt(1),
				OnTimePledge:  filWhole(2),
				ActivePower:   NewPowerPair(big.NewInt(3), big.NewInt(4)),
				FaultyPower:   NewPowerPair(big.NewInt(5), big.NewInt(6)),
				FeeDeduction:  filWhole(7),
			},
			// [[DA],[GA],[ABvBbWdOyAAA],[[AAM],[AAQ]],[[AAU],[AAY]],[AGEk/umTvAAA]]
			readHex: "86410c411849001bc16d674ec80000824200034200048242000542000649006124fee993bc0000",
			// same on write as read
			writeHex: "86410c411849001bc16d674ec80000824200034200048242000542000649006124fee993bc0000",
		},
		{
			set: ExpirationSet{
				OnTimeSectors: bfrt(0),
				EarlySectors:  bfrt(1),
				OnTimePledge:  filWhole(2),
				ActivePower:   NewPowerPair(big.NewInt(3), big.NewInt(4)),
				FaultyPower:   NewPowerPair(big.NewInt(5), big.NewInt(6)),
				FeeDeduction:  big.Int{},
			},
			// [[DA],[GA],[ABvBbWdOyAAA],[[AAM],[AAQ]],[[AAU],[AAY]]]
			readHex: "85410c411849001bc16d674ec800008242000342000482420005420006",
			// [[DA],[GA],[ABvBbWdOyAAA],[[AAM],[AAQ]],[[AAU],[AAY]],[]]
			writeHex: "86410c411849001bc16d674ec80000824200034200048242000542000640",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			req := require.New(t)

			// write
			var buf bytes.Buffer
			req.NoError(tc.set.MarshalCBOR(&buf))
			req.Equal(tc.writeHex, hex.EncodeToString(buf.Bytes()))

			// read
			byts, err := hex.DecodeString(tc.readHex)
			req.NoError(err)
			var rt ExpirationSet
			req.NoError(rt.UnmarshalCBOR(bytes.NewReader(byts)))
			req.Equal(tc.set, rt)
		})
	}
}

func TestDeadline(t *testing.T) {
	zero := big.Zero()

	testCases := []struct {
		deadline Deadline
		hex      string
	}{
		{
			deadline: Deadline{
				Partitions:                        defaultCid,
				ExpirationsEpochs:                 defaultCid,
				PartitionsPoSted:                  bfrt(),
				EarlyTerminations:                 bfrt(),
				FaultyPower:                       NewPowerPairZero(),
				OptimisticPoStSubmissions:         defaultCid,
				SectorsSnapshot:                   defaultCid,
				PartitionsSnapshot:                defaultCid,
				OptimisticPoStSubmissionsSnapshot: defaultCid,
				LivePower:                         NewPowerPairZero(),
				DailyFee:                          zero,
			},
			// [baeaaaaa,baeaaaaa,[],[],0,0,[[],[]],baeaaaaa,baeaaaaa,baeaaaaa,baeaaaaa,[[],[]],[]]
			hex: "8dd82a450001000000d82a45000100000040400000824040d82a450001000000d82a450001000000d82a450001000000d82a45000100000082404040",
		},
		{
			deadline: Deadline{
				Partitions:                        cid.MustParse("bagboea4seaaqa"),
				ExpirationsEpochs:                 cid.MustParse("bagboea4seaaqc"),
				PartitionsPoSted:                  bfrt(0),
				EarlyTerminations:                 bfrt(1),
				LiveSectors:                       2,
				TotalSectors:                      3,
				FaultyPower:                       NewPowerPair(big.NewInt(4), big.NewInt(5)),
				OptimisticPoStSubmissions:         cid.MustParse("bagboea4seaaqe"),
				SectorsSnapshot:                   cid.MustParse("bagboea4seaaqg"),
				PartitionsSnapshot:                cid.MustParse("bagboea4seaaqi"),
				OptimisticPoStSubmissionsSnapshot: cid.MustParse("bagboea4seaaqk"),
				LivePower:                         NewPowerPair(big.NewInt(6), big.NewInt(7)),
				DailyFee:                          filWhole(8),
			},
			// [bagboea4seaaqa,bagboea4seaaqc,[DA],[GA],2,3,[[AAQ],[AAU]],bagboea4seaaqe,bagboea4seaaqg,bagboea4seaaqi,bagboea4seaaqk,[[AAY],[AAc]],[AG8FtZ07IAAA]]
			hex: "8dd82a49000182e20392200100d82a49000182e20392200101410c4118020382420004420005d82a49000182e20392200102d82a49000182e20392200103d82a49000182e20392200104d82a49000182e203922001058242000642000749006f05b59d3b200000",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.deadline.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt Deadline
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.deadline, rt)
		})
	}
}

func filWhole(i int64) abi.TokenAmount {
	return big.Mul(big.NewInt(i), builtin.TokenPrecision)
}

// bfrt makes a bitfield from a list of integers and round-trips it through encoding
// so that the internal representation is set properly for a deep equals() test
func bfrt(b ...uint64) bitfield.BitField {
	bf := bitfield.NewFromSet(b)
	var buf bytes.Buffer
	if err := bf.MarshalCBOR(&buf); err != nil {
		panic(err)
	}
	var rt bitfield.BitField
	if err := rt.UnmarshalCBOR(&buf); err != nil {
		panic(err)
	}
	return rt
}
