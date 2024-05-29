package verifreg

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
)

// Tests to match with Rust fil_actor_verifreg::serialization

func TestSerializationClaimAllocationsParams(t *testing.T) {
	testCases := []struct {
		params ClaimAllocationsParams
		hex    string
	}{
		{
			params: ClaimAllocationsParams{Sectors: nil, AllOrNothing: false},
			// [[],false]
			hex: "8280f4",
		},
		{
			params: ClaimAllocationsParams{
				Sectors: []SectorAllocationClaims{{
					Sector:       101,
					SectorExpiry: 202,
					Claims:       nil,
				}},
				AllOrNothing: true,
			},
			// [[[101,202,[]]],true]
			hex: "828183186518ca80f5",
		},
		{
			params: ClaimAllocationsParams{
				Sectors: []SectorAllocationClaims{{
					Sector:       101,
					SectorExpiry: 202,
					Claims: []AllocationClaim{
						{
							Client:       303,
							AllocationId: 404,
							Data:         cid.MustParse("baga6ea4seaaqa"),
							Size:         505,
						},
						{
							Client:       606,
							AllocationId: 707,
							Data:         cid.MustParse("baga6ea4seaaqc"),
							Size:         808,
						},
					},
				},
					{Sector: 303, SectorExpiry: 404, Claims: nil},
				},
				AllOrNothing: true,
			},
			// [[[101,202,[[303,404,baga6ea4seaaqa,505],[606,707,baga6ea4seaaqc,808]]],[303,404,[]]],true]
			hex: "828283186518ca828419012f190194d82a49000181e203922001001901f98419025e1902c3d82a49000181e203922001011903288319012f19019480f5",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt ClaimAllocationsParams
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.params, rt)
		})
	}
}
