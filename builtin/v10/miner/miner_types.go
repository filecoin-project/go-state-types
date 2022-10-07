package miner

import (
	power "github.com/filecoin-project/go-state-types/builtin/v10/power"
	miner9 "github.com/filecoin-project/go-state-types/builtin/v9/miner"

	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
	cid "github.com/ipfs/go-cid"
)

type DeclareFaultsRecoveredParams = miner9.DeclareFaultsRecoveredParams
type RecoveryDeclaration = miner9.RecoveryDeclaration
type DeclareFaultsParams = miner9.DeclareFaultsParams
type FaultDeclaration = miner9.FaultDeclaration
type ReplicaUpdate = miner9.ReplicaUpdate
type ProveReplicaUpdatesParams = miner9.ProveReplicaUpdatesParams
type ReplicaUpdate2 = miner9.ReplicaUpdate2
type ProveReplicaUpdatesParams2 = miner9.ProveReplicaUpdatesParams2
type PoStPartition = miner9.PoStPartition
type SubmitWindowedPoStParams = miner9.SubmitWindowedPoStParams
type DisputeWindowedPoStParams = miner9.DisputeWindowedPoStParams
type ProveCommitAggregateParams = miner9.ProveCommitAggregateParams
type ProveCommitSectorParams = miner9.ProveCommitSectorParams
type MinerConstructorParams = power.MinerConstructorParams
type TerminateSectorsParams = miner9.TerminateSectorsParams
type TerminationDeclaration = miner9.TerminationDeclaration
type TerminateSectorsReturn = miner9.TerminateSectorsReturn
type ChangePeerIDParams = miner9.ChangePeerIDParams
type ChangeMultiaddrsParams = miner9.ChangeMultiaddrsParams
type ChangeWorkerAddressParams = miner9.ChangeWorkerAddressParams
type ExtendSectorExpirationParams = miner9.ExtendSectorExpirationParams
type ExpirationExtension = miner9.ExpirationExtension
type ReportConsensusFaultParams = miner9.ReportConsensusFaultParams
type GetControlAddressesReturn = miner9.GetControlAddressesReturn
type CheckSectorProvenParams = miner9.CheckSectorProvenParams
type WithdrawBalanceParams = miner9.WithdrawBalanceParams
type CompactPartitionsParams = miner9.CompactPartitionsParams
type CompactSectorNumbersParams = miner9.CompactSectorNumbersParams
type CronEventType = miner9.CronEventType

const (
	CronEventWorkerKeyChange          = miner9.CronEventWorkerKeyChange
	CronEventProvingDeadline          = miner9.CronEventProvingDeadline
	CronEventProcessEarlyTerminations = miner9.CronEventProcessEarlyTerminations
)

type CronEventPayload = miner9.CronEventPayload
type PartitionKey = miner9.PartitionKey
type PreCommitSectorParams = miner9.PreCommitSectorParams
type PreCommitSectorBatchParams = miner9.PreCommitSectorBatchParams
type PreCommitSectorBatchParams2 = miner9.PreCommitSectorBatchParams2
type ChangeBeneficiaryParams = miner9.ChangeBeneficiaryParams
type ActiveBeneficiary = miner9.ActiveBeneficiary
type GetBeneficiaryReturn = miner9.GetBeneficiaryReturn
type ExpirationSet = miner9.ExpirationSet
type ExpirationQueue = miner9.ExpirationQueue

func LoadExpirationQueue(store adt.Store, root cid.Cid, quant builtin.QuantSpec, bitwidth int) (ExpirationQueue, error) {
	return miner9.LoadExpirationQueue(store, root, quant, bitwidth)
}

type Sectors = miner9.Sectors
type VestingFunds = miner9.VestingFunds
type VestingFund = miner9.VestingFund

func ConstructVestingFunds() *VestingFunds {
	return miner9.ConstructVestingFunds()
}

type DeferredCronEventParams = miner9.DeferredCronEventParams
type ApplyRewardParams = miner9.ApplyRewardParams
type ConfirmSectorProofsParams = miner9.ConfirmSectorProofsParams
type ExtendSectorExpiration2Params = miner9.ExtendSectorExpiration2Params
type ExpirationExtension2 = miner9.ExpirationExtension2
type SectorClaim = miner9.SectorClaim
