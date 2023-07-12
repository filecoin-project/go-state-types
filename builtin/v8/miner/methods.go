package miner

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v8/power"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1:  {"Constructor", *new(func(*power.MinerConstructorParams) *abi.EmptyValue)},            // Constructor
	2:  {"ControlAddresses", *new(func(*abi.EmptyValue) *GetControlAddressesReturn)},          // ControlAddresses
	3:  {"ChangeWorkerAddress", *new(func(*ChangeWorkerAddressParams) *abi.EmptyValue)},       // ChangeWorkerAddress
	4:  {"ChangePeerID", *new(func(*ChangePeerIDParams) *abi.EmptyValue)},                     // ChangePeerID
	5:  {"SubmitWindowedPoSt", *new(func(*SubmitWindowedPoStParams) *abi.EmptyValue)},         // SubmitWindowedPoSt
	6:  {"PreCommitSector", *new(func(*PreCommitSectorParams) *abi.EmptyValue)},               // PreCommitSector
	7:  {"ProveCommitSector", *new(func(*ProveCommitSectorParams) *abi.EmptyValue)},           // ProveCommitSector
	8:  {"ExtendSectorExpiration", *new(func(*ExtendSectorExpirationParams) *abi.EmptyValue)}, // ExtendSectorExpiration
	9:  {"TerminateSectors", *new(func(*TerminateSectorsParams) *TerminateSectorsReturn)},     // TerminateSectors
	10: {"DeclareFaults", *new(func(*DeclareFaultsParams) *abi.EmptyValue)},                   // DeclareFaults
	11: {"DeclareFaultsRecovered", *new(func(*DeclareFaultsRecoveredParams) *abi.EmptyValue)}, // DeclareFaultsRecovered
	12: {"OnDeferredCronEvent", *new(func(*DeferredCronEventParams) *abi.EmptyValue)},         // OnDeferredCronEvent
	13: {"CheckSectorProven", *new(func(*CheckSectorProvenParams) *abi.EmptyValue)},           // CheckSectorProven
	14: {"ApplyRewards", *new(func(*ApplyRewardParams) *abi.EmptyValue)},                      // ApplyRewards
	15: {"ReportConsensusFault", *new(func(*ReportConsensusFaultParams) *abi.EmptyValue)},     // ReportConsensusFault
	16: {"WithdrawBalance", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)},              // WithdrawBalance
	17: {"ConfirmSectorProofsValid", *new(func(*ConfirmSectorProofsParams) *abi.EmptyValue)},  // ConfirmSectorProofsValid
	18: {"ChangeMultiaddrs", *new(func(*ChangeMultiaddrsParams) *abi.EmptyValue)},             // ChangeMultiaddrs
	19: {"CompactPartitions", *new(func(*CompactPartitionsParams) *abi.EmptyValue)},           // CompactPartitions
	20: {"CompactSectorNumbers", *new(func(*CompactSectorNumbersParams) *abi.EmptyValue)},     // CompactSectorNumbers
	21: {"ConfirmUpdateWorkerKey", *new(func(*abi.EmptyValue) *abi.EmptyValue)},               // ConfirmUpdateWorkerKey
	22: {"RepayDebt", *new(func(*abi.EmptyValue) *abi.EmptyValue)},                            // RepayDebt
	23: {"ChangeOwnerAddress", *new(func(*address.Address) *abi.EmptyValue)},                  // ChangeOwnerAddress
	24: {"DisputeWindowedPoSt", *new(func(*DisputeWindowedPoStParams) *abi.EmptyValue)},       // DisputeWindowedPoSt
	25: {"PreCommitSectorBatch", *new(func(*PreCommitSectorBatchParams) *abi.EmptyValue)},     // PreCommitSectorBatch
	26: {"ProveCommitAggregate", *new(func(*ProveCommitAggregateParams) *abi.EmptyValue)},     // ProveCommitAggregate
	27: {"ProveReplicaUpdates", *new(func(*ProveReplicaUpdatesParams) *bitfield.BitField)},    // ProveReplicaUpdates
}
