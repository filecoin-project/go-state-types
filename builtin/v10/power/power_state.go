package power

import (
	power9 "github.com/filecoin-project/go-state-types/builtin/v9/power"

	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
)

var InitialQAPowerEstimatePosition = power9.InitialQAPowerEstimatePosition
var InitialQAPowerEstimateVelocity = power9.InitialQAPowerEstimateVelocity

const CronQueueHamtBitwidth = power9.CronQueueHamtBitwidth
const CronQueueAmtBitwidth = power9.CronQueueAmtBitwidth
const ProofValidationBatchAmtBitwidth = power9.ProofValidationBatchAmtBitwidth
const ConsensusMinerMinMiners = power9.ConsensusMinerMinMiners

type State = power9.State

func ConstructState(store adt.Store) (*State, error) {
	return power9.ConstructState(store)
}

type Claim = power9.Claim
