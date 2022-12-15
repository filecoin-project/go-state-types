package builtin

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/frc0042"
)

const (
	MethodSend        = abi.MethodNum(0)
	MethodConstructor = abi.MethodNum(1)
)

var MethodsAccount = struct {
	Constructor           abi.MethodNum
	PubkeyAddress         abi.MethodNum
	AuthenticateMessage   abi.MethodNum
	UniversalReceiverHook abi.MethodNum
}{
	MethodConstructor,
	2,
	frc0042.GenerateExportedMethodNum("AuthenticateMessage"),
	frc0042.GenerateExportedMethodNum("Receive"),
}

var MethodsInit = struct {
	Constructor  abi.MethodNum
	Exec         abi.MethodNum
	ExecExported abi.MethodNum
	// TODO Exec4?
}{
	MethodConstructor,
	2,
	frc0042.GenerateExportedMethodNum("Exec"),
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
	frc0042.GenerateExportedMethodNum("Propose"),
	3,
	frc0042.GenerateExportedMethodNum("Approve"),
	4,
	frc0042.GenerateExportedMethodNum("Cancel"),
	5,
	frc0042.GenerateExportedMethodNum("AddSigner"),
	6,
	frc0042.GenerateExportedMethodNum("RemoveSigner"),
	7,
	frc0042.GenerateExportedMethodNum("SwapSigner"),
	8,
	frc0042.GenerateExportedMethodNum("ChangeNumApprovalsThreshold"),
	9,
	frc0042.GenerateExportedMethodNum("LockBalance"),
	frc0042.GenerateExportedMethodNum("Receive"),
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
	frc0042.GenerateExportedMethodNum("AddBalance"),
	3,
	frc0042.GenerateExportedMethodNum("WithdrawBalance"),
	4,
	frc0042.GenerateExportedMethodNum("PublishStorageDeals"),
	5,
	6,
	7,
	8,
	9,
	frc0042.GenerateExportedMethodNum("GetBalance"),
	frc0042.GenerateExportedMethodNum("GetDealDataCommitment"),
	frc0042.GenerateExportedMethodNum("GetDealClient"),
	frc0042.GenerateExportedMethodNum("GetDealProvider"),
	frc0042.GenerateExportedMethodNum("GetDealLabel"),
	frc0042.GenerateExportedMethodNum("GetDealTerm"),
	frc0042.GenerateExportedMethodNum("GetDealTotalPrice"),
	frc0042.GenerateExportedMethodNum("GetDealClientCollateral"),
	frc0042.GenerateExportedMethodNum("GetDealProviderCollateral"),
	frc0042.GenerateExportedMethodNum("GetDealVerified"),
	frc0042.GenerateExportedMethodNum("GetDealActivation"),
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
	frc0042.GenerateExportedMethodNum("CreateMiner"),
	3,
	4,
	5,
	6,
	7,
	8,
	9,
	frc0042.GenerateExportedMethodNum("NetworkRawPower"),
	frc0042.GenerateExportedMethodNum("MinerRawPower"),
	frc0042.GenerateExportedMethodNum("MinerCount"),
	frc0042.GenerateExportedMethodNum("MinerConsensusCount"),
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
	frc0042.GenerateExportedMethodNum("ChangeWorkerAddress"),
	4,
	frc0042.GenerateExportedMethodNum("ChangePeerID"),
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
	frc0042.GenerateExportedMethodNum("WithdrawBalance"),
	17,
	18,
	frc0042.GenerateExportedMethodNum("ChangeMultiaddrs"),
	19,
	20,
	21,
	frc0042.GenerateExportedMethodNum("ConfirmChangeWorkerAddress"),
	22,
	frc0042.GenerateExportedMethodNum("RepayDebt"),
	23,
	frc0042.GenerateExportedMethodNum("ChangeOwnerAddress"),
	24,
	25,
	26,
	27,
	28,
	29,
	30,
	frc0042.GenerateExportedMethodNum("ChangeBeneficiary"),
	31,
	32,
	frc0042.GenerateExportedMethodNum("GetOwner"),
	frc0042.GenerateExportedMethodNum("IsControllingAddress"),
	frc0042.GenerateExportedMethodNum("GetSectorSize"),
	frc0042.GenerateExportedMethodNum("GetAvailableBalance"),
	frc0042.GenerateExportedMethodNum("GetVestingFunds"),
	frc0042.GenerateExportedMethodNum("GetPeerID"),
	frc0042.GenerateExportedMethodNum("GetMultiaddrs"),
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
	frc0042.GenerateExportedMethodNum("AddVerifiedClient"),
	5,
	6,
	7,
	8,
	frc0042.GenerateExportedMethodNum("RemoveExpiredAllocations"),
	9,
	10,
	frc0042.GenerateExportedMethodNum("GetClaims"),
	11,
	frc0042.GenerateExportedMethodNum("ExtendClaimTerms"),
	12,
	frc0042.GenerateExportedMethodNum("RemoveExpiredClaims"),
	frc0042.GenerateExportedMethodNum("Receive"),
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
	frc0042.GenerateExportedMethodNum("Mint"),
	frc0042.GenerateExportedMethodNum("Destroy"),
	frc0042.GenerateExportedMethodNum("Name"),
	frc0042.GenerateExportedMethodNum("Symbol"),
	frc0042.GenerateExportedMethodNum("TotalSupply"),
	frc0042.GenerateExportedMethodNum("Balance"),
	frc0042.GenerateExportedMethodNum("Transfer"),
	frc0042.GenerateExportedMethodNum("TransferFrom"),
	frc0042.GenerateExportedMethodNum("IncreaseAllowance"),
	frc0042.GenerateExportedMethodNum("DecreaseAllowance"),
	frc0042.GenerateExportedMethodNum("RevokeAllowance"),
	frc0042.GenerateExportedMethodNum("Burn"),
	frc0042.GenerateExportedMethodNum("BurnFrom"),
	frc0042.GenerateExportedMethodNum("Allowance"),
	frc0042.GenerateExportedMethodNum("Granularity"),
}

var MethodsEVM = struct {
	Constructor            abi.MethodNum
	InvokeContract         abi.MethodNum
	GetBytecode            abi.MethodNum
	GetStorageAt           abi.MethodNum
	InvokeContractReadOnly abi.MethodNum
	InvokeContractDelegate abi.MethodNum
}{
	MethodConstructor,
	2,
	3,
	4,
	5,
	6,
}

var MethodsEAM = struct {
	Constructor abi.MethodNum
	Create      abi.MethodNum
	Create2     abi.MethodNum
}{
	MethodConstructor,
	2,
	3,
}
