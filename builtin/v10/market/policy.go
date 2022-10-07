package market

import (
	market9 "github.com/filecoin-project/go-state-types/builtin/v9/market"

	"github.com/filecoin-project/go-state-types/abi"
)

var ProviderCollateralSupplyTarget = market9.ProviderCollateralSupplyTarget
var DealMinDuration = market9.DealMinDuration
var DealMaxDuration = market9.DealMaxDuration
var MarketDefaultAllocationTermBuffer = market9.MarketDefaultAllocationTermBuffer

func DealDurationBounds(ps abi.PaddedPieceSize) (min abi.ChainEpoch, max abi.ChainEpoch) {
	return market9.DealDurationBounds(ps)
}

const DealMaxLabelSize = market9.DealMaxLabelSize

func DealPricePerEpochBounds(ps abi.PaddedPieceSize, epoch abi.ChainEpoch) (min abi.TokenAmount, max abi.TokenAmount) {
	return market9.DealPricePerEpochBounds(ps, epoch)
}

func DealProviderCollateralBounds(pieceSize abi.PaddedPieceSize, verified bool, networkRawPower, networkQAPower, baselinePower abi.StoragePower,
	networkCirculatingSupply abi.TokenAmount) (min, max abi.TokenAmount) {
	return market9.DealProviderCollateralBounds(pieceSize, verified, networkRawPower, networkQAPower, baselinePower, networkCirculatingSupply)
}

func DealClientCollateralBounds(ps abi.PaddedPieceSize, epoch abi.ChainEpoch) (min abi.TokenAmount, max abi.TokenAmount) {
	return market9.DealClientCollateralBounds(ps, epoch)
}

// Computes the weight for a deal proposal, which is a function of its size and duration.
func DealWeight(proposal *DealProposal) abi.DealWeight {
	return market9.DealWeight(proposal)
}
