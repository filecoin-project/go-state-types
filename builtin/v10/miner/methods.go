package miner

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v10/power"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*power.MinerConstructorParams) *abi.EmptyValue)},      // Constructor
	2: {"ControlAddresses", *new(func(*abi.EmptyValue) *GetControlAddressesReturn)},    // ControlAddresses
	3: {"ChangeWorkerAddress", *new(func(*ChangeWorkerAddressParams) *abi.EmptyValue)}, // ChangeWorkerAddress
	builtin.MustGenerateFRCMethodNum("ChangeWorkerAddress"): {"ChangeWorkerAddressExported", *new(func(*ChangeWorkerAddressParams) *abi.EmptyValue)}, // ChangeWorkerAddressExported
	4: {"ChangePeerID", *new(func(*ChangePeerIDParams) *abi.EmptyValue)}, // ChangePeerID
	builtin.MustGenerateFRCMethodNum("ChangePeerID"): {"ChangePeerIDExported", *new(func(*ChangePeerIDParams) *abi.EmptyValue)}, // ChangePeerIDExported
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
	builtin.MustGenerateFRCMethodNum("WithdrawBalance"): {"WithdrawBalanceExported", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)}, // WithdrawBalanceExported
	17: {"ConfirmSectorProofsValid", *new(func(*ConfirmSectorProofsParams) *abi.EmptyValue)}, // ConfirmSectorProofsValid
	18: {"ChangeMultiaddrs", *new(func(*ChangeMultiaddrsParams) *abi.EmptyValue)},            // ChangeMultiaddrs
	builtin.MustGenerateFRCMethodNum("ChangeMultiaddrs"): {"ChangeMultiaddrsExported", *new(func(*ChangeMultiaddrsParams) *abi.EmptyValue)}, // ChangeMultiaddrsExported
	19: {"CompactPartitions", *new(func(*CompactPartitionsParams) *abi.EmptyValue)},       // CompactPartitions
	20: {"CompactSectorNumbers", *new(func(*CompactSectorNumbersParams) *abi.EmptyValue)}, // CompactSectorNumbers
	21: {"ConfirmChangeWorkerAddress", *new(func(*abi.EmptyValue) *abi.EmptyValue)},       // ConfirmChangeWorkerAddress
	builtin.MustGenerateFRCMethodNum("ConfirmChangeWorkerAddress"): {"ConfirmChangeWorkerAddressExported", *new(func(*abi.EmptyValue) *abi.EmptyValue)}, // ConfirmChangeWorkerAddressExported
	22: {"RepayDebt", *new(func(*abi.EmptyValue) *abi.EmptyValue)}, // RepayDebt
	builtin.MustGenerateFRCMethodNum("RepayDebt"): {"RepayDebtExported", *new(func(*abi.EmptyValue) *abi.EmptyValue)}, // RepayDebtExported
	23: {"ChangeOwnerAddress", *new(func(*address.Address) *abi.EmptyValue)}, // ChangeOwnerAddress
	builtin.MustGenerateFRCMethodNum("ChangeOwnerAddress"): {"ChangeOwnerAddressExported", *new(func(*address.Address) *abi.EmptyValue)}, // ChangeOwnerAddressExported
	24: {"DisputeWindowedPoSt", *new(func(*DisputeWindowedPoStParams) *abi.EmptyValue)},    // DisputeWindowedPoSt
	25: {"PreCommitSectorBatch", *new(func(*PreCommitSectorBatchParams) *abi.EmptyValue)},  // PreCommitSectorBatch
	26: {"ProveCommitAggregate", *new(func(*ProveCommitAggregateParams) *abi.EmptyValue)},  // ProveCommitAggregate
	27: {"ProveReplicaUpdates", *new(func(*ProveReplicaUpdatesParams) *bitfield.BitField)}, // ProveReplicaUpdates
	// NB: the name of this method must not change across actor/network versions
	28: {"PreCommitSectorBatch2", *new(func(*PreCommitSectorBatchParams2) *abi.EmptyValue)}, // PreCommitSectorBatch2
	// NB: the name of this method must not change across actor/network versions
	29: {"ProveReplicaUpdates2", *new(func(*ProveReplicaUpdatesParams2) *bitfield.BitField)}, // ProveReplicaUpdates2
	30: {"ChangeBeneficiary", *new(func(*ChangeBeneficiaryParams) *abi.EmptyValue)},          // ChangeBeneficiary
	builtin.MustGenerateFRCMethodNum("ChangeBeneficiary"): {"ChangeBeneficiaryExported", *new(func(*ChangeBeneficiaryParams) *abi.EmptyValue)}, // ChangeBeneficiaryExported
	31: {"GetBeneficiary", *new(func(*abi.EmptyValue) *GetBeneficiaryReturn)}, // GetBeneficiary
	// NB: the name of this method must not change across actor/network versions
	32: {"ExtendSectorExpiration2", *new(func(*ExtendSectorExpiration2Params) *abi.EmptyValue)}, // ExtendSectorExpiration2
	builtin.MustGenerateFRCMethodNum("GetOwner"):             {"GetOwnerExported", *new(func(*abi.EmptyValue) *GetOwnerReturn)},                                            // GetOwnerExported
	builtin.MustGenerateFRCMethodNum("IsControllingAddress"): {"IsControllingAddressExported", *new(func(params *IsControllingAddressParams) *IsControllingAddressReturn)}, // IsControllingAddressExported
	builtin.MustGenerateFRCMethodNum("GetSectorSize"):        {"GetSectorSizeExported", *new(func(*abi.EmptyValue) *GetSectorSizeReturn)},                                  // GetSectorSizeExported
	builtin.MustGenerateFRCMethodNum("GetAvailableBalance"):  {"GetAvailableBalanceExported", *new(func(*abi.EmptyValue) *GetAvailableBalanceReturn)},                      // GetAvailableBalanceExported
	builtin.MustGenerateFRCMethodNum("GetVestingFunds"):      {"GetVestingFundsExported", *new(func(*abi.EmptyValue) *GetVestingFundsReturn)},                              // GetVestingFundsExported
	builtin.MustGenerateFRCMethodNum("GetPeerID"):            {"GetPeerIDExported", *new(func(*abi.EmptyValue) *GetPeerIDReturn)},                                          // GetPeerIDExported
	builtin.MustGenerateFRCMethodNum("GetMultiaddrs"):        {"GetMultiaddrsExported", *new(func(*abi.EmptyValue) *GetMultiAddrsReturn)},                                  // GetMultiaddrsExported
}
