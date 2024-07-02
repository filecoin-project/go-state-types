package miner

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v14/power"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: builtin.NewMethodMeta("Constructor", *new(func(*power.MinerConstructorParams) *abi.EmptyValue)),      // Constructor
	2: builtin.NewMethodMeta("ControlAddresses", *new(func(*abi.EmptyValue) *GetControlAddressesReturn)),    // ControlAddresses
	3: builtin.NewMethodMeta("ChangeWorkerAddress", *new(func(*ChangeWorkerAddressParams) *abi.EmptyValue)), // ChangeWorkerAddress
	builtin.MustGenerateFRCMethodNum("ChangeWorkerAddress"): builtin.NewMethodMeta("ChangeWorkerAddressExported", *new(func(*ChangeWorkerAddressParams) *abi.EmptyValue)), // ChangeWorkerAddressExported
	4: builtin.NewMethodMeta("ChangePeerID", *new(func(*ChangePeerIDParams) *abi.EmptyValue)), // ChangePeerID
	builtin.MustGenerateFRCMethodNum("ChangePeerID"): builtin.NewMethodMeta("ChangePeerIDExported", *new(func(*ChangePeerIDParams) *abi.EmptyValue)), // ChangePeerIDExported
	5:  builtin.NewMethodMeta("SubmitWindowedPoSt", *new(func(*SubmitWindowedPoStParams) *abi.EmptyValue)),         // SubmitWindowedPoSt
	6:  builtin.NewMethodMeta("PreCommitSector", *new(func(*PreCommitSectorParams) *abi.EmptyValue)),               // PreCommitSector
	7:  builtin.NewMethodMeta("ProveCommitSector", *new(func(*ProveCommitSectorParams) *abi.EmptyValue)),           // ProveCommitSector
	8:  builtin.NewMethodMeta("ExtendSectorExpiration", *new(func(*ExtendSectorExpirationParams) *abi.EmptyValue)), // ExtendSectorExpiration
	9:  builtin.NewMethodMeta("TerminateSectors", *new(func(*TerminateSectorsParams) *TerminateSectorsReturn)),     // TerminateSectors
	10: builtin.NewMethodMeta("DeclareFaults", *new(func(*DeclareFaultsParams) *abi.EmptyValue)),                   // DeclareFaults
	11: builtin.NewMethodMeta("DeclareFaultsRecovered", *new(func(*DeclareFaultsRecoveredParams) *abi.EmptyValue)), // DeclareFaultsRecovered
	12: builtin.NewMethodMeta("OnDeferredCronEvent", *new(func(*DeferredCronEventParams) *abi.EmptyValue)),         // OnDeferredCronEvent
	13: builtin.NewMethodMeta("CheckSectorProven", *new(func(*CheckSectorProvenParams) *abi.EmptyValue)),           // CheckSectorProven
	14: builtin.NewMethodMeta("ApplyRewards", *new(func(*ApplyRewardParams) *abi.EmptyValue)),                      // ApplyRewards
	15: builtin.NewMethodMeta("ReportConsensusFault", *new(func(*ReportConsensusFaultParams) *abi.EmptyValue)),     // ReportConsensusFault
	16: builtin.NewMethodMeta("WithdrawBalance", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)),              // WithdrawBalance
	builtin.MustGenerateFRCMethodNum("WithdrawBalance"): builtin.NewMethodMeta("WithdrawBalanceExported", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)), // WithdrawBalanceExported
	17: builtin.NewMethodMeta("InternalSectorSetupForPreseal", *new(func(*InternalSectorSetupForPresealParams) *abi.EmptyValue)), // InternalSectorSetupForPreseal
	18: builtin.NewMethodMeta("ChangeMultiaddrs", *new(func(*ChangeMultiaddrsParams) *abi.EmptyValue)),                           // ChangeMultiaddrs
	builtin.MustGenerateFRCMethodNum("ChangeMultiaddrs"): builtin.NewMethodMeta("ChangeMultiaddrsExported", *new(func(*ChangeMultiaddrsParams) *abi.EmptyValue)), // ChangeMultiaddrsExported
	19: builtin.NewMethodMeta("CompactPartitions", *new(func(*CompactPartitionsParams) *abi.EmptyValue)),       // CompactPartitions
	20: builtin.NewMethodMeta("CompactSectorNumbers", *new(func(*CompactSectorNumbersParams) *abi.EmptyValue)), // CompactSectorNumbers
	21: builtin.NewMethodMeta("ConfirmChangeWorkerAddress", *new(func(*abi.EmptyValue) *abi.EmptyValue)),       // ConfirmChangeWorkerAddress
	builtin.MustGenerateFRCMethodNum("ConfirmChangeWorkerAddress"): builtin.NewMethodMeta("ConfirmChangeWorkerAddressExported", *new(func(*abi.EmptyValue) *abi.EmptyValue)), // ConfirmChangeWorkerAddressExported
	22: builtin.NewMethodMeta("RepayDebt", *new(func(*abi.EmptyValue) *abi.EmptyValue)), // RepayDebt
	builtin.MustGenerateFRCMethodNum("RepayDebt"): builtin.NewMethodMeta("RepayDebtExported", *new(func(*abi.EmptyValue) *abi.EmptyValue)), // RepayDebtExported
	23: builtin.NewMethodMeta("ChangeOwnerAddress", *new(func(*address.Address) *abi.EmptyValue)), // ChangeOwnerAddress
	builtin.MustGenerateFRCMethodNum("ChangeOwnerAddress"): builtin.NewMethodMeta("ChangeOwnerAddressExported", *new(func(*address.Address) *abi.EmptyValue)), // ChangeOwnerAddressExported
	24: builtin.NewMethodMeta("DisputeWindowedPoSt", *new(func(*DisputeWindowedPoStParams) *abi.EmptyValue)),    // DisputeWindowedPoSt
	25: builtin.NewMethodMeta("PreCommitSectorBatch", *new(func(*PreCommitSectorBatchParams) *abi.EmptyValue)),  // PreCommitSectorBatch
	26: builtin.NewMethodMeta("ProveCommitAggregate", *new(func(*ProveCommitAggregateParams) *abi.EmptyValue)),  // ProveCommitAggregate
	27: builtin.NewMethodMeta("ProveReplicaUpdates", *new(func(*ProveReplicaUpdatesParams) *bitfield.BitField)), // ProveReplicaUpdates
	// NB: the name of this method must not change across actor/network versions
	28: builtin.NewMethodMeta("PreCommitSectorBatch2", *new(func(*PreCommitSectorBatchParams2) *abi.EmptyValue)), // PreCommitSectorBatch2
	// NB: the name of this method must not change across actor/network versions
	29: builtin.NewMethodMeta("ProveReplicaUpdates2", *new(func(*ProveReplicaUpdatesParams2) *bitfield.BitField)), // ProveReplicaUpdates2
	30: builtin.NewMethodMeta("ChangeBeneficiary", *new(func(*ChangeBeneficiaryParams) *abi.EmptyValue)),          // ChangeBeneficiary
	builtin.MustGenerateFRCMethodNum("ChangeBeneficiary"): builtin.NewMethodMeta("ChangeBeneficiaryExported", *new(func(*ChangeBeneficiaryParams) *abi.EmptyValue)), // ChangeBeneficiaryExported
	31: builtin.NewMethodMeta("GetBeneficiary", *new(func(*abi.EmptyValue) *GetBeneficiaryReturn)), // GetBeneficiary
	// NB: the name of this method must not change across actor/network versions
	32: builtin.NewMethodMeta("ExtendSectorExpiration2", *new(func(*ExtendSectorExpiration2Params) *abi.EmptyValue)), // ExtendSectorExpiration2
	builtin.MustGenerateFRCMethodNum("GetOwner"):             builtin.NewMethodMeta("GetOwnerExported", *new(func(*abi.EmptyValue) *GetOwnerReturn)),                                            // GetOwnerExported
	builtin.MustGenerateFRCMethodNum("IsControllingAddress"): builtin.NewMethodMeta("IsControllingAddressExported", *new(func(params *IsControllingAddressParams) *IsControllingAddressReturn)), // IsControllingAddressExported
	builtin.MustGenerateFRCMethodNum("GetSectorSize"):        builtin.NewMethodMeta("GetSectorSizeExported", *new(func(*abi.EmptyValue) *GetSectorSizeReturn)),                                  // GetSectorSizeExported
	builtin.MustGenerateFRCMethodNum("GetAvailableBalance"):  builtin.NewMethodMeta("GetAvailableBalanceExported", *new(func(*abi.EmptyValue) *GetAvailableBalanceReturn)),                      // GetAvailableBalanceExported
	builtin.MustGenerateFRCMethodNum("GetVestingFunds"):      builtin.NewMethodMeta("GetVestingFundsExported", *new(func(*abi.EmptyValue) *GetVestingFundsReturn)),                              // GetVestingFundsExported
	builtin.MustGenerateFRCMethodNum("GetPeerID"):            builtin.NewMethodMeta("GetPeerIDExported", *new(func(*abi.EmptyValue) *GetPeerIDReturn)),                                          // GetPeerIDExported
	builtin.MustGenerateFRCMethodNum("GetMultiaddrs"):        builtin.NewMethodMeta("GetMultiaddrsExported", *new(func(*abi.EmptyValue) *GetMultiAddrsReturn)),                                  // GetMultiaddrsExported
	// 33 MovePartitions
	34: builtin.NewMethodMeta("ProveCommitSectors3", *new(func(*ProveCommitSectors3Params) *ProveCommitSectors3Return)),    // ProveCommitSectors3
	35: builtin.NewMethodMeta("ProveReplicaUpdates3", *new(func(*ProveReplicaUpdates3Params) *ProveReplicaUpdates3Return)), // ProveReplicaUpdates3
	36: builtin.NewMethodMeta("ProveCommitSectorsNI", *new(func(*ProveCommitSectorsNIParams) *abi.EmptyValue)),             // ProveCommitSectorsNI
}
