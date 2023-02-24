package builtin

import (
	"github.com/filecoin-project/go-state-types/abi"
)

const (
	MethodSend        = abi.MethodNum(0)
	MethodConstructor = abi.MethodNum(1)
)

var MethodsAccount = struct {
	Constructor         abi.MethodNum
	PubkeyAddress       abi.MethodNum
	AuthenticateMessage abi.MethodNum
}{
	MethodConstructor,
	2,
	MustGenerateFRCMethodNum("AuthenticateMessage"),
}

var MethodsInit = struct {
	Constructor abi.MethodNum
	Exec        abi.MethodNum
	Exec4       abi.MethodNum
}{
	MethodConstructor,
	2,
	3,
}

var MethodsCron = struct {
	Constructor abi.MethodNum
	EpochTick   abi.MethodNum
}{
	MethodConstructor,
	2,
}

var MethodsReward = struct {
	Constructor      abi.MethodNum
	AwardBlockReward abi.MethodNum
	ThisEpochReward  abi.MethodNum
	UpdateNetworkKPI abi.MethodNum
}{
	MethodConstructor,
	2,
	3,
	4,
}

var MethodsMultisig = struct {
	Constructor                         abi.MethodNum
	Propose                             abi.MethodNum
	ProposeExported                     abi.MethodNum
	Approve                             abi.MethodNum
	ApproveExported                     abi.MethodNum
	Cancel                              abi.MethodNum
	CancelExported                      abi.MethodNum
	AddSigner                           abi.MethodNum
	AddSignerExported                   abi.MethodNum
	RemoveSigner                        abi.MethodNum
	RemoveSignerExported                abi.MethodNum
	SwapSigner                          abi.MethodNum
	SwapSignerExported                  abi.MethodNum
	ChangeNumApprovalsThreshold         abi.MethodNum
	ChangeNumApprovalsThresholdExported abi.MethodNum
	LockBalance                         abi.MethodNum
	LockBalanceExported                 abi.MethodNum
	UniversalReceiverHook               abi.MethodNum
}{
	MethodConstructor,
	2,
	MustGenerateFRCMethodNum("Propose"),
	3,
	MustGenerateFRCMethodNum("Approve"),
	4,
	MustGenerateFRCMethodNum("Cancel"),
	5,
	MustGenerateFRCMethodNum("AddSigner"),
	6,
	MustGenerateFRCMethodNum("RemoveSigner"),
	7,
	MustGenerateFRCMethodNum("SwapSigner"),
	8,
	MustGenerateFRCMethodNum("ChangeNumApprovalsThreshold"),
	9,
	MustGenerateFRCMethodNum("LockBalance"),
	MustGenerateFRCMethodNum("Receive"),
}

var MethodsPaych = struct {
	Constructor        abi.MethodNum
	UpdateChannelState abi.MethodNum
	Settle             abi.MethodNum
	Collect            abi.MethodNum
}{
	MethodConstructor,
	2,
	3,
	4,
}

var MethodsMarket = struct {
	Constructor                       abi.MethodNum
	AddBalance                        abi.MethodNum
	AddBalanceExported                abi.MethodNum
	WithdrawBalance                   abi.MethodNum
	WithdrawBalanceExported           abi.MethodNum
	PublishStorageDeals               abi.MethodNum
	PublishStorageDealsExported       abi.MethodNum
	VerifyDealsForActivation          abi.MethodNum
	ActivateDeals                     abi.MethodNum
	OnMinerSectorsTerminate           abi.MethodNum
	ComputeDataCommitment             abi.MethodNum
	CronTick                          abi.MethodNum
	GetBalanceExported                abi.MethodNum
	GetDealDataCommitmentExported     abi.MethodNum
	GetDealClientExported             abi.MethodNum
	GetDealProviderExported           abi.MethodNum
	GetDealLabelExported              abi.MethodNum
	GetDealTermExported               abi.MethodNum
	GetDealTotalPriceExported         abi.MethodNum
	GetDealClientCollateralExported   abi.MethodNum
	GetDealProviderCollateralExported abi.MethodNum
	GetDealVerifiedExported           abi.MethodNum
	GetDealActivationExported         abi.MethodNum
}{
	MethodConstructor,
	2,
	MustGenerateFRCMethodNum("AddBalance"),
	3,
	MustGenerateFRCMethodNum("WithdrawBalance"),
	4,
	MustGenerateFRCMethodNum("PublishStorageDeals"),
	5,
	6,
	7,
	8,
	9,
	MustGenerateFRCMethodNum("GetBalance"),
	MustGenerateFRCMethodNum("GetDealDataCommitment"),
	MustGenerateFRCMethodNum("GetDealClient"),
	MustGenerateFRCMethodNum("GetDealProvider"),
	MustGenerateFRCMethodNum("GetDealLabel"),
	MustGenerateFRCMethodNum("GetDealTerm"),
	MustGenerateFRCMethodNum("GetDealTotalPrice"),
	MustGenerateFRCMethodNum("GetDealClientCollateral"),
	MustGenerateFRCMethodNum("GetDealProviderCollateral"),
	MustGenerateFRCMethodNum("GetDealVerified"),
	MustGenerateFRCMethodNum("GetDealActivation"),
}

var MethodsPower = struct {
	Constructor                                  abi.MethodNum
	CreateMiner                                  abi.MethodNum
	CreateMinerExported                          abi.MethodNum
	UpdateClaimedPower                           abi.MethodNum
	EnrollCronEvent                              abi.MethodNum
	CronTick                                     abi.MethodNum
	UpdatePledgeTotal                            abi.MethodNum
	Deprecated1                                  abi.MethodNum
	SubmitPoRepForBulkVerify                     abi.MethodNum
	CurrentTotalPower                            abi.MethodNum
	CurrentTotalPowerNetworkRawPowerExported     abi.MethodNum
	CurrentTotalPowerMinerRawPowerExported       abi.MethodNum
	CurrentTotalPowerMinerCountExported          abi.MethodNum
	CurrentTotalPowerMinerConsensusCountExported abi.MethodNum
}{
	MethodConstructor,
	2,
	MustGenerateFRCMethodNum("CreateMiner"),
	3,
	4,
	5,
	6,
	7,
	8,
	9,
	MustGenerateFRCMethodNum("NetworkRawPower"),
	MustGenerateFRCMethodNum("MinerRawPower"),
	MustGenerateFRCMethodNum("MinerCount"),
	MustGenerateFRCMethodNum("MinerConsensusCount"),
}

var MethodsMiner = struct {
	Constructor                        abi.MethodNum
	ControlAddresses                   abi.MethodNum
	ChangeWorkerAddress                abi.MethodNum
	ChangeWorkerAddressExported        abi.MethodNum
	ChangePeerID                       abi.MethodNum
	ChangePeerIDExported               abi.MethodNum
	SubmitWindowedPoSt                 abi.MethodNum
	PreCommitSector                    abi.MethodNum
	ProveCommitSector                  abi.MethodNum
	ExtendSectorExpiration             abi.MethodNum
	TerminateSectors                   abi.MethodNum
	DeclareFaults                      abi.MethodNum
	DeclareFaultsRecovered             abi.MethodNum
	OnDeferredCronEvent                abi.MethodNum
	CheckSectorProven                  abi.MethodNum
	ApplyRewards                       abi.MethodNum
	ReportConsensusFault               abi.MethodNum
	WithdrawBalance                    abi.MethodNum
	WithdrawBalanceExported            abi.MethodNum
	ConfirmSectorProofsValid           abi.MethodNum
	ChangeMultiaddrs                   abi.MethodNum
	ChangeMultiaddrsExported           abi.MethodNum
	CompactPartitions                  abi.MethodNum
	CompactSectorNumbers               abi.MethodNum
	ConfirmChangeWorkerAddress         abi.MethodNum
	ConfirmChangeWorkerAddressExported abi.MethodNum
	RepayDebt                          abi.MethodNum
	RepayDebtExported                  abi.MethodNum
	ChangeOwnerAddress                 abi.MethodNum
	ChangeOwnerAddressExported         abi.MethodNum
	DisputeWindowedPoSt                abi.MethodNum
	PreCommitSectorBatch               abi.MethodNum
	ProveCommitAggregate               abi.MethodNum
	ProveReplicaUpdates                abi.MethodNum
	PreCommitSectorBatch2              abi.MethodNum
	ProveReplicaUpdates2               abi.MethodNum
	ChangeBeneficiary                  abi.MethodNum
	ChangeBeneficiaryExported          abi.MethodNum
	GetBeneficiary                     abi.MethodNum
	ExtendSectorExpiration2            abi.MethodNum
	GetOwnerExported                   abi.MethodNum
	IsControllingAddressExported       abi.MethodNum
	GetSectorSizeExported              abi.MethodNum
	GetAvailableBalanceExported        abi.MethodNum
	GetVestingFundsExported            abi.MethodNum
	GetPeerIDExported                  abi.MethodNum
	GetMultiaddrsExported              abi.MethodNum
}{
	MethodConstructor,
	2,
	3,
	MustGenerateFRCMethodNum("ChangeWorkerAddress"),
	4,
	MustGenerateFRCMethodNum("ChangePeerID"),
	5,
	6,
	7,
	8,
	9,
	10,
	11,
	12,
	13,
	14,
	15,
	16,
	MustGenerateFRCMethodNum("WithdrawBalance"),
	17,
	18,
	MustGenerateFRCMethodNum("ChangeMultiaddrs"),
	19,
	20,
	21,
	MustGenerateFRCMethodNum("ConfirmChangeWorkerAddress"),
	22,
	MustGenerateFRCMethodNum("RepayDebt"),
	23,
	MustGenerateFRCMethodNum("ChangeOwnerAddress"),
	24,
	25,
	26,
	27,
	28,
	29,
	30,
	MustGenerateFRCMethodNum("ChangeBeneficiary"),
	31,
	32,
	MustGenerateFRCMethodNum("GetOwner"),
	MustGenerateFRCMethodNum("IsControllingAddress"),
	MustGenerateFRCMethodNum("GetSectorSize"),
	MustGenerateFRCMethodNum("GetAvailableBalance"),
	MustGenerateFRCMethodNum("GetVestingFunds"),
	MustGenerateFRCMethodNum("GetPeerID"),
	MustGenerateFRCMethodNum("GetMultiaddrs"),
}

var MethodsVerifiedRegistry = struct {
	Constructor                      abi.MethodNum
	AddVerifier                      abi.MethodNum
	RemoveVerifier                   abi.MethodNum
	AddVerifiedClient                abi.MethodNum
	AddVerifiedClientExported        abi.MethodNum
	Deprecated1                      abi.MethodNum
	Deprecated2                      abi.MethodNum
	RemoveVerifiedClientDataCap      abi.MethodNum
	RemoveExpiredAllocations         abi.MethodNum
	RemoveExpiredAllocationsExported abi.MethodNum
	ClaimAllocations                 abi.MethodNum
	GetClaims                        abi.MethodNum
	GetClaimsExported                abi.MethodNum
	ExtendClaimTerms                 abi.MethodNum
	ExtendClaimTermsExported         abi.MethodNum
	RemoveExpiredClaims              abi.MethodNum
	RemoveExpiredClaimsExported      abi.MethodNum
	UniversalReceiverHook            abi.MethodNum
}{
	MethodConstructor,
	2,
	3,
	4,
	MustGenerateFRCMethodNum("AddVerifiedClient"),
	5,
	6,
	7,
	8,
	MustGenerateFRCMethodNum("RemoveExpiredAllocations"),
	9,
	10,
	MustGenerateFRCMethodNum("GetClaims"),
	11,
	MustGenerateFRCMethodNum("ExtendClaimTerms"),
	12,
	MustGenerateFRCMethodNum("RemoveExpiredClaims"),
	MustGenerateFRCMethodNum("Receive"),
}

var MethodsDatacap = struct {
	Constructor               abi.MethodNum
	MintExported              abi.MethodNum
	DestroyExported           abi.MethodNum
	NameExported              abi.MethodNum
	SymbolExported            abi.MethodNum
	TotalSupplyExported       abi.MethodNum
	BalanceExported           abi.MethodNum
	TransferExported          abi.MethodNum
	TransferFromExported      abi.MethodNum
	IncreaseAllowanceExported abi.MethodNum
	DecreaseAllowanceExported abi.MethodNum
	RevokeAllowanceExported   abi.MethodNum
	BurnExported              abi.MethodNum
	BurnFromExported          abi.MethodNum
	AllowanceExported         abi.MethodNum
	GranularityExported       abi.MethodNum
}{
	MethodConstructor,
	MustGenerateFRCMethodNum("Mint"),
	MustGenerateFRCMethodNum("Destroy"),
	MustGenerateFRCMethodNum("Name"),
	MustGenerateFRCMethodNum("Symbol"),
	MustGenerateFRCMethodNum("TotalSupply"),
	MustGenerateFRCMethodNum("Balance"),
	MustGenerateFRCMethodNum("Transfer"),
	MustGenerateFRCMethodNum("TransferFrom"),
	MustGenerateFRCMethodNum("IncreaseAllowance"),
	MustGenerateFRCMethodNum("DecreaseAllowance"),
	MustGenerateFRCMethodNum("RevokeAllowance"),
	MustGenerateFRCMethodNum("Burn"),
	MustGenerateFRCMethodNum("BurnFrom"),
	MustGenerateFRCMethodNum("Allowance"),
	MustGenerateFRCMethodNum("Granularity"),
}

var MethodsEVM = struct {
	Constructor            abi.MethodNum
	Resurrect              abi.MethodNum
	GetBytecode            abi.MethodNum
	GetBytecodeHash        abi.MethodNum
	GetStorageAt           abi.MethodNum
	InvokeContractDelegate abi.MethodNum
	InvokeContract         abi.MethodNum
}{
	MethodConstructor,
	2,
	3,
	4,
	5,
	6,
	MustGenerateFRCMethodNum("InvokeEVM"),
}

var MethodsEAM = struct {
	Constructor    abi.MethodNum
	Create         abi.MethodNum
	Create2        abi.MethodNum
	CreateExternal abi.MethodNum
}{
	MethodConstructor,
	2,
	3,
	4,
}

var MethodsPlaceholder = struct {
	Constructor abi.MethodNum
}{
	MethodConstructor,
}

var MethodsEthAccount = struct {
	Constructor abi.MethodNum
}{
	MethodConstructor,
}
