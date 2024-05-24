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
	apt := abi.RegisteredAggregationProof_SnarkPackV2

	testCases := []struct {
		params ProveCommitSectorsNIParams
		hex    string
	}{
		{
			params: ProveCommitSectorsNIParams{Sectors: nil, SealProofType: abi.RegisteredSealProof_StackedDrg32GiBV1_1},
			// [[],8,[],byte[],null,false]
			hex: "8680088040f6f4",
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
				SealProofType:            abi.RegisteredSealProof_StackedDrg32GiBV1_1_Feat_NiPoRep,
				SectorProofs:             nil,
				AggregateProof:           []byte{0xde, 0xad, 0xbe, 0xef},
				AggregateProofType:       &apt,
				RequireActivationSuccess: true,
			},
			// [[[1,2,bagboea4seaaqa,3,4,5]],18,[],byte[deadbeef],1,true]
			hex: "8681860102d82a49000182e20392200100030405128044deadbeef01f5",
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
				SealProofType:            abi.RegisteredSealProof_StackedDrg32GiBV1_1_Feat_NiPoRep,
				SectorProofs:             [][]byte{[]byte{0xde, 0xad}, []byte{0xbe, 0xef}},
				AggregateProof:           nil,
				AggregateProofType:       nil,
				RequireActivationSuccess: false,
			},
			// [[[1,2,bagboea4seaaqa,3,4,5],[6,7,bagboea4seaaqc,8,9,10]],18,[byte[dead],byte[beef]],byte[],null,false]
			hex: "8682860102d82a49000182e20392200100030405860607d82a49000182e2039220010108090a128242dead42beef40f6f4",
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
