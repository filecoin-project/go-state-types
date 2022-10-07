package smoothing

import (
	smoothing9 "github.com/filecoin-project/go-state-types/builtin/v9/util/smoothing"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
)

var (
	DefaultAlpha = smoothing9.DefaultAlpha
	DefaultBeta  = smoothing9.DefaultBeta

	ExtrapolatedCumSumRatioEpsilon = smoothing9.ExtrapolatedCumSumRatioEpsilon
)

type FilterEstimate = smoothing9.FilterEstimate

// Create a new filter estimate given two Q.0 format ints.
func NewEstimate(position, velocity big.Int) FilterEstimate {
	return smoothing9.NewEstimate(position, velocity)
}

// Extrapolate the CumSumRatio given two filters.
// Output is in Q.128 format
func ExtrapolatedCumSumOfRatio(delta abi.ChainEpoch, relativeStart abi.ChainEpoch, estimateNum, estimateDenom FilterEstimate) big.Int {
	return smoothing9.ExtrapolatedCumSumOfRatio(delta, relativeStart, estimateNum, estimateDenom)
}
