package power

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v19/util/smoothing"
	"github.com/stretchr/testify/require"
)

// Test to match with Rust fil_actor_power::serialization
func TestSerializationCurrentTotalPowerReturn(t *testing.T) {
	testCases := []struct {
		params CurrentTotalPowerReturn
		hex    string
	}{
		{
			params: CurrentTotalPowerReturn{
				RawBytePower:            abi.NewStoragePower(0),
				QualityAdjPower:         abi.NewStoragePower(0),
				PledgeCollateral:        abi.NewTokenAmount(0),
				QualityAdjPowerSmoothed: smoothing.NewEstimate(big.Zero(), big.Zero()),
			},
			// [byte[],byte[],byte[],[byte[],byte[]]]
			hex: "84404040824040",
		},
		{
			params: CurrentTotalPowerReturn{
				RawBytePower:            abi.NewStoragePower(1 << 20),
				QualityAdjPower:         abi.NewStoragePower(1 << 21),
				PledgeCollateral:        abi.NewTokenAmount(1 << 22),
				QualityAdjPowerSmoothed: smoothing.NewEstimate(big.NewInt(1<<23), big.NewInt(1<<24)),
			},
			// FilterEstimate BigInts have a precision shift of 128, so they end up larger than the others.
			// [byte[00100000],byte[00200000],byte[00400000],[byte[0080000000000000000000000000000000000000],byte[000100000000000000000000000000000000000000]]]
			hex: "844400100000440020000044004000008254008000000000000000000000000000000000000055000100000000000000000000000000000000000000",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt CurrentTotalPowerReturn
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.params, rt)
		})
	}
}
