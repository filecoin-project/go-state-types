package miner

import (
	miner9 "github.com/filecoin-project/go-state-types/builtin/v9/miner"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v10/util/smoothing"
)

var PreCommitDepositFactor = miner9.PreCommitDepositFactor
var PreCommitDepositProjectionPeriod = miner9.PreCommitDepositProjectionPeriod
var InitialPledgeFactor = miner9.InitialPledgeFactor
var InitialPledgeProjectionPeriod = miner9.InitialPledgeProjectionPeriod
var InitialPledgeMaxPerByte = miner9.InitialPledgeMaxPerByte
var InitialPledgeLockTarget = miner9.InitialPledgeLockTarget

func ExpectedRewardForPower(rewardEstimate, networkQAPowerEstimate smoothing.FilterEstimate, qaSectorPower abi.StoragePower, projectionDuration abi.ChainEpoch) abi.TokenAmount {
	return miner9.ExpectedRewardForPower(rewardEstimate, networkQAPowerEstimate, qaSectorPower, projectionDuration)
}

func ExpectedRewardForPowerClampedAtAttoFIL(rewardEstimate, networkQAPowerEstimate smoothing.FilterEstimate, qaSectorPower abi.StoragePower, projectionDuration abi.ChainEpoch) abi.TokenAmount {
	return miner9.ExpectedRewardForPowerClampedAtAttoFIL(rewardEstimate, networkQAPowerEstimate, qaSectorPower, projectionDuration)
}

func PreCommitDepositForPower(rewardEstimate, networkQAPowerEstimate smoothing.FilterEstimate, qaSectorPower abi.StoragePower) abi.TokenAmount {
	return miner9.PreCommitDepositForPower(rewardEstimate, networkQAPowerEstimate, qaSectorPower)
}

func InitialPledgeForPower(qaPower, baselinePower abi.StoragePower, rewardEstimate, networkQAPowerEstimate smoothing.FilterEstimate, circulatingSupply abi.TokenAmount) abi.TokenAmount {
	return miner9.InitialPledgeForPower(qaPower, baselinePower, rewardEstimate, networkQAPowerEstimate, circulatingSupply)
}

var EstimatedSingleProveCommitGasUsage = miner9.EstimatedSingleProveCommitGasUsage
var EstimatedSinglePreCommitGasUsage = miner9.EstimatedSinglePreCommitGasUsage
var BatchDiscount = miner9.BatchDiscount
var BatchBalancer = miner9.BatchBalancer

func AggregateProveCommitNetworkFee(aggregateSize int, baseFee abi.TokenAmount) abi.TokenAmount {
	return miner9.AggregateProveCommitNetworkFee(aggregateSize, baseFee)
}

func AggregatePreCommitNetworkFee(aggregateSize int, baseFee abi.TokenAmount) abi.TokenAmount {
	return miner9.AggregatePreCommitNetworkFee(aggregateSize, baseFee)
}
