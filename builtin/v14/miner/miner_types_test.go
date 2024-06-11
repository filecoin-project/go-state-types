package miner

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/filecoin-project/go-state-types/abi"
	cid "github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
)

// Tests to match with Rust fil_actor_miner::serialization

func TestSerializationSlaimAllocationsParams(t *testing.T) {
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
