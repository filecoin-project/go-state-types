package verifreg

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
)

// Tests to match with Rust fil_actor_verifreg::serialization

func TestSerializationSlaimAllocationsParams(t *testing.T) {
	testCases := []struct {
		params ClaimAllocationsParams
		hex    string
	}{
		{
			/*
			 82                                                # array(2)
			   80                                              #   array(0)
			   f4                                              #   false
			*/
			params: ClaimAllocationsParams{Sectors: nil, AllOrNothing: false},
			hex:    "8280f4",
		},
		{
			/*
			 82                                                # array(2)
			   81                                              #   array(1)
			     83                                            #     array(3)
			       18 65                                       #       uint(101)
			       18 ca                                       #       uint(202)
			       80                                          #       array(0)
			   f5                                              #   true
			*/
			params: ClaimAllocationsParams{
				Sectors: []SectorAllocationClaims{{
					Sector:       101,
					SectorExpiry: 202,
					Claims:       nil,
				}},
				AllOrNothing: true,
			},
			hex: "828183186518ca80f5",
		},
		{
			/*
			 82                                                # array(2)
			   82                                              #   array(2)
			     83                                            #     array(3)
			       18 65                                       #       uint(101)
			       18 ca                                       #       uint(202)
			       82                                          #       array(2)
			         84                                        #         array(4)
			           19 012f                                 #           uint(303)
			           19 0194                                 #           uint(404)
			           d8 2a                                   #           tag(42)
			             49                                    #             bytes(9)
			               000181e20392200100                  #               "\x00\x01\x81â\x03\x92 \x01\x00"
			           19 01f9                                 #           uint(505)
			         84                                        #         array(4)
			           19 025e                                 #           uint(606)
			           19 02c3                                 #           uint(707)
			           d8 2a                                   #           tag(42)
			             49                                    #             bytes(9)
			               000181e20392200101                  #               "\x00\x01\x81â\x03\x92 \x01\x01"
			           19 0328                                 #           uint(808)
			     83                                            #     array(3)
			       19 012f                                     #       uint(303)
			       19 0194                                     #       uint(404)
			       80                                          #       array(0)
			   f5                                              #   true
			*/
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
