package miner

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v9/power"
)

var Methods = map[uint64]builtin.MethodMeta{
	1:  {"", *new(func(*power.MinerConstructorParams) *abi.EmptyValue)},   // Constructor
	2:  {"", *new(func(*abi.EmptyValue) *GetControlAddressesReturn)},      // ControlAddresses
	3:  {"", *new(func(*ChangeWorkerAddressParams) *abi.EmptyValue)},      // ChangeWorkerAddress
	4:  {"", *new(func(*ChangePeerIDParams) *abi.EmptyValue)},             // ChangePeerID
	5:  {"", *new(func(*SubmitWindowedPoStParams) *abi.EmptyValue)},       // SubmitWindowedPoSt
	6:  {"", *new(func(*PreCommitSectorParams) *abi.EmptyValue)},          // PreCommitSector
	7:  {"", *new(func(*ProveCommitSectorParams) *abi.EmptyValue)},        // ProveCommitSector
	8:  {"", *new(func(*ExtendSectorExpirationParams) *abi.EmptyValue)},   // ExtendSectorExpiration
	9:  {"", *new(func(*TerminateSectorsParams) *TerminateSectorsReturn)}, // TerminateSectors
	10: {"", *new(func(*DeclareFaultsParams) *abi.EmptyValue)},            // DeclareFaults
	11: {"", *new(func(*DeclareFaultsRecoveredParams) *abi.EmptyValue)},   // DeclareFaultsRecovered
	12: {"", *new(func(*DeferredCronEventParams) *abi.EmptyValue)},        // OnDeferredCronEvent
	13: {"", *new(func(*CheckSectorProvenParams) *abi.EmptyValue)},        // CheckSectorProven
	14: {"", *new(func(*ApplyRewardParams) *abi.EmptyValue)},              // ApplyRewards
	15: {"", *new(func(*ReportConsensusFaultParams) *abi.EmptyValue)},     // ReportConsensusFault
	16: {"", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)},         // WithdrawBalance
	17: {"", *new(func(*ConfirmSectorProofsParams) *abi.EmptyValue)},      // ConfirmSectorProofsValid
	18: {"", *new(func(*ChangeMultiaddrsParams) *abi.EmptyValue)},         // ChangeMultiaddrs
	19: {"", *new(func(*CompactPartitionsParams) *abi.EmptyValue)},        // CompactPartitions
	20: {"", *new(func(*CompactSectorNumbersParams) *abi.EmptyValue)},     // CompactSectorNumbers
	21: {"", *new(func(*abi.EmptyValue) *abi.EmptyValue)},                 // ConfirmUpdateWorkerKey
	22: {"", *new(func(*abi.EmptyValue) *abi.EmptyValue)},                 // RepayDebt
	23: {"", *new(func(*address.Address) *abi.EmptyValue)},                // ChangeOwnerAddress
	24: {"", *new(func(*DisputeWindowedPoStParams) *abi.EmptyValue)},      // DisputeWindowedPoSt
	25: {"", *new(func(*PreCommitSectorBatchParams) *abi.EmptyValue)},     // PreCommitSectorBatch
	26: {"", *new(func(*ProveCommitAggregateParams) *abi.EmptyValue)},     // ProveCommitAggregate
	27: {"", *new(func(*ProveReplicaUpdatesParams) *bitfield.BitField)},   // ProveReplicaUpdates
	28: {"", *new(func(*PreCommitSectorBatchParams2) *abi.EmptyValue)},    // PreCommitSectorBatch2
	29: {"", *new(func(*ProveReplicaUpdatesParams2) *bitfield.BitField)},  // ProveReplicaUpdates2
	30: {"", *new(func(*ChangeBeneficiaryParams) *abi.EmptyValue)},        // ChangeBeneficiary
	31: {"", *new(func(*abi.EmptyValue) *GetBeneficiaryReturn)},           // GetBeneficiary
	32: {"", *new(func(*ExtendSectorExpiration2Params) *abi.EmptyValue)},  // ExtendSectorExpiration2
}
