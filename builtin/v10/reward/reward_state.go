package reward

import (
	reward9 "github.com/filecoin-project/go-state-types/builtin/v9/reward"

	"github.com/filecoin-project/go-state-types/abi"
)

type Spacetime = reward9.Spacetime

const InitialRewardPositionEstimateStr = reward9.InitialRewardPositionEstimateStr

var InitialRewardPositionEstimate = reward9.InitialRewardPositionEstimate

var InitialRewardVelocityEstimate = reward9.InitialRewardVelocityEstimate

type State = reward9.State

func ConstructState(currRealizedPower abi.StoragePower) *State {
	return reward9.ConstructState(currRealizedPower)
}
