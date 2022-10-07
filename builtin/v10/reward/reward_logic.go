package reward

import (
	reward9 "github.com/filecoin-project/go-state-types/builtin/v9/reward"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
)

var BaselineExponent = reward9.BaselineExponent
var BaselineInitialValue = reward9.BaselineInitialValue

func InitBaselinePower() abi.StoragePower {
	return reward9.InitBaselinePower()
}

func BaselinePowerFromPrev(prevEpochBaselinePower abi.StoragePower) abi.StoragePower {
	return reward9.BaselinePowerFromPrev(prevEpochBaselinePower)
}

var DefaultSimpleTotal = reward9.DefaultSimpleTotal
var DefaultBaselineTotal = reward9.DefaultBaselineTotal

func ComputeRTheta(effectiveNetworkTime abi.ChainEpoch, baselinePowerAtEffectiveNetworkTime, cumsumRealized, cumsumBaseline big.Int) big.Int {
	return reward9.ComputeRTheta(effectiveNetworkTime, baselinePowerAtEffectiveNetworkTime, cumsumRealized, cumsumBaseline)
}

var (
	Lambda       = reward9.Lambda
	ExpLamSubOne = reward9.ExpLamSubOne
)
