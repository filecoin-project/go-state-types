package power

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v16/util/smoothing"
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
				RampStartEpoch:          0,
				RampDurationEpochs:      0,
			},
			// [byte[],byte[],byte[],[byte[],byte[]],0,0]
			hex: "864040408240400000",
		},
		{
			params: CurrentTotalPowerReturn{
				RawBytePower:            abi.NewStoragePower(1 << 20),
				QualityAdjPower:         abi.NewStoragePower(1 << 21),
				PledgeCollateral:        abi.NewTokenAmount(1 << 22),
				QualityAdjPowerSmoothed: smoothing.NewEstimate(big.NewInt(1<<23), big.NewInt(1<<24)),
				RampStartEpoch:          25,
				RampDurationEpochs:      26,
			},
			// FilterEstimate BigInts have a precision shift of 128, so they end up larger than the others.
			// [byte[00100000],byte[00200000],byte[00400000],[byte[0080000000000000000000000000000000000000],byte[000100000000000000000000000000000000000000]],25,26]
			hex: "8644001000004400200000440040000082540080000000000000000000000000000000000000550001000000000000000000000000000000000000001819181a",
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
