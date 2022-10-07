package miner

import (
	miner9 "github.com/filecoin-project/go-state-types/builtin/v9/miner"

	"github.com/filecoin-project/go-state-types/abi"
)

var WPoStProvingPeriod = miner9.WPoStProvingPeriod
var WPoStChallengeWindow = miner9.WPoStChallengeWindow
var WPoStDisputeWindow = miner9.WPoStDisputeWindow

const WPoStPeriodDeadlines = miner9.WPoStPeriodDeadlines
const MaxPartitionsPerDeadline = miner9.MaxPartitionsPerDeadline
const AddressedPartitionsMax = miner9.AddressedPartitionsMax
const DeclarationsMax = miner9.DeclarationsMax
const AddressedSectorsMax = miner9.AddressedSectorsMax
const ChainFinality = miner9.ChainFinality

var SealedCIDPrefix = miner9.SealedCIDPrefix
var WindowPoStProofTypes = miner9.WindowPoStProofTypes

var MaxProveCommitDuration = miner9.MaxProveCommitDuration

const PreCommitSectorBatchMaxSize = miner9.PreCommitSectorBatchMaxSize
const ProveReplicaUpdatesMaxSize = miner9.ProveReplicaUpdatesMaxSize

var MaxPreCommitRandomnessLookback = miner9.MaxPreCommitRandomnessLookback
var PreCommitChallengeDelay = miner9.PreCommitChallengeDelay

const WPoStChallengeLookback = miner9.WPoStChallengeLookback
const FaultDeclarationCutoff = miner9.FaultDeclarationCutoff

var FaultMaxAge = miner9.FaultMaxAge

const WorkerKeyChangeDelay = miner9.WorkerKeyChangeDelay
const MinSectorExpiration = miner9.MinSectorExpiration
const MaxSectorExpirationExtension = miner9.MaxSectorExpirationExtension

func QualityForWeight(size abi.SectorSize, duration abi.ChainEpoch, dealWeight, verifiedWeight abi.DealWeight) abi.SectorQuality {
	return miner9.QualityForWeight(size, duration, dealWeight, verifiedWeight)
}

func QAPowerForWeight(size abi.SectorSize, duration abi.ChainEpoch, dealWeight, verifiedWeight abi.DealWeight) abi.StoragePower {
	return miner9.QAPowerForWeight(size, duration, dealWeight, verifiedWeight)
}

const MaxAggregatedSectors = miner9.MaxAggregatedSectors
const MinAggregatedSectors = miner9.MinAggregatedSectors
const MaxAggregateProofSize = miner9.MaxAggregateProofSize

type VestSpec = miner9.VestSpec

func QAPowerMax(size abi.SectorSize) abi.StoragePower {
	return miner9.QAPowerMax(size)
}
