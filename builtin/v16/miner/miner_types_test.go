package miner

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	cid "github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
)

// Test to match with Rust fil_actor_miner::serialization
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

// Test to match with Rust fil_actor_miner::serialization
func TestSerializationSectorOnChainInfo(t *testing.T) {
	testCases := []struct {
		params SectorOnChainInfo
		old    bool
		hex    string
	}{
		{
			params: SectorOnChainInfo{
				SealProof:             -1,
				SealedCID:             cid.MustParse("baeaaaaa"),
				DealWeight:            big.Zero(),
				VerifiedDealWeight:    big.Zero(),
				InitialPledge:         big.Zero(),
				ExpectedDayReward:     big.Zero(),
				ExpectedStoragePledge: big.Zero(),
				ReplacedDayReward:     big.Zero(),
				ProvingPeriodFee:      big.Zero(),
			},
			old: false,
			// [0,-1,{"/":"baeaaaaa"},[],0,0,[],[],[],[],[],0,[],null,0,[]]
			hex: "900020d82a45000100000080000040404040400040f60040",
		},
		{
			params: SectorOnChainInfo{
				SectorNumber:          1,
				SealProof:             abi.RegisteredSealProof_StackedDrg32GiBV1_1,
				SealedCID:             cid.MustParse("bagboea4seaaqa"),
				DealIDs:               nil,
				Activation:            2,
				Expiration:            3,
				DealWeight:            big.NewInt(4),
				VerifiedDealWeight:    big.NewInt(5),
				InitialPledge:         big.Mul(big.NewInt(6), builtin.TokenPrecision),
				ExpectedDayReward:     big.Mul(big.NewInt(7), builtin.TokenPrecision),
				ExpectedStoragePledge: big.Mul(big.NewInt(8), builtin.TokenPrecision),
				PowerBaseEpoch:        9,
				ReplacedDayReward:     big.Mul(big.NewInt(10), builtin.TokenPrecision),
				SectorKeyCID:          nil,
				Flags:                 0,
				ProvingPeriodFee:      big.Mul(big.NewInt(11), builtin.TokenPrecision),
			},
			old: false,
			// '[1,8,{"/":"bagboea4seaaqa"},[],2,3,[AAQ],[AAU],[AFNESDXsWAAA],[AGEk]TvAAA"}},[AG8FtZ07IAAA],9,[AIrHIwSJ6AAA],null,0,[AJin2bgxTAAA]]'
			hex: "900108d82a49000182e20392200100800203420004420005490053444835ec58000049006124fee993bc000049006f05b59d3b2000000949008ac7230489e80000f600490098a7d9b8314c0000",
		},
		{
			params: SectorOnChainInfo{
				SectorNumber:          1,
				SealProof:             abi.RegisteredSealProof_StackedDrg32GiBV1_1,
				SealedCID:             cid.MustParse("bagboea4seaaqa"),
				DealIDs:               nil,
				Activation:            2,
				Expiration:            3,
				DealWeight:            big.NewInt(4),
				VerifiedDealWeight:    big.NewInt(5),
				InitialPledge:         big.Mul(big.NewInt(6), builtin.TokenPrecision),
				ExpectedDayReward:     big.Mul(big.NewInt(7), builtin.TokenPrecision),
				ExpectedStoragePledge: big.Mul(big.NewInt(8), builtin.TokenPrecision),
				PowerBaseEpoch:        9,
				ReplacedDayReward:     big.Mul(big.NewInt(10), builtin.TokenPrecision),
				SectorKeyCID:          ptr(cid.MustParse("baga6ea4seaaqc")),
				Flags:                 1,
				ProvingPeriodFee:      big.Mul(big.NewInt(11), builtin.TokenPrecision),
			},
			old: false,
			// [1,8,{"/":"bagboea4seaaqa"},[],2,3,[AAQ],[AAU],[AFNESDXsWAAA],[AGEk]TvAAA"}},[AG8FtZ07IAAA],9,[AIrHIwSJ6AAA],{"/":"baga6ea4seaaqc"},1,[AJin2bgxTAAA]]
			hex: "900108d82a49000182e20392200100800203420004420005490053444835ec58000049006124fee993bc000049006f05b59d3b2000000949008ac7230489e80000d82a49000181e2039220010101490098a7d9b8314c0000",
		},
		{
			// old format stored on chain but materialised as the new format with a default value at the end
			params: SectorOnChainInfo{
				SectorNumber:          1,
				SealProof:             abi.RegisteredSealProof_StackedDrg64GiBV1_1,
				SealedCID:             cid.MustParse("bagboea4seaaqa"),
				DealIDs:               nil,
				Activation:            2,
				Expiration:            3,
				DealWeight:            big.NewInt(4),
				VerifiedDealWeight:    big.NewInt(5),
				InitialPledge:         big.Mul(big.NewInt(6), builtin.TokenPrecision),
				ExpectedDayReward:     big.Mul(big.NewInt(7), builtin.TokenPrecision),
				ExpectedStoragePledge: big.Mul(big.NewInt(8), builtin.TokenPrecision),
				PowerBaseEpoch:        9,
				ReplacedDayReward:     big.Mul(big.NewInt(10), builtin.TokenPrecision),
				SectorKeyCID:          nil,
				Flags:                 1,
				ProvingPeriodFee:      big.Mul(big.NewInt(0), builtin.TokenPrecision),
			},
			old: true,
			// [1,9,{"/":"bagboea4seaaqa"},[],2,3,[AAQ],[AAU],[AFNESDXsWAAA],[AGEk]TvAAA"}},[AG8FtZ07IAAA],9,[AIrHIwSJ6AAA],null,1]
			hex: "8f0109d82a49000182e20392200100800203420004420005490053444835ec58000049006124fee993bc000049006f05b59d3b2000000949008ac7230489e80000f601",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			req := require.New(t)

			if !tc.old {
				var buf bytes.Buffer
				req.NoError(tc.params.MarshalCBOR(&buf))
				req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			}
			var rt SectorOnChainInfo
			byts, err := hex.DecodeString(tc.hex)
			req.NoError(err)
			req.NoError(rt.UnmarshalCBOR(bytes.NewReader(byts)))
			req.Equal(tc.params, rt)
		})
	}
}

func ptr[T any](t T) *T {
	return &t
}
