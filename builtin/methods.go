package builtin

import (
	"github.com/filecoin-project/go-state-types/abi"
)

const (
	MethodSend                     = abi.MethodNum(0)
	MethodConstructor              = abi.MethodNum(1)
	UniversalReceiverHookMethodNum = abi.MethodNum(3726118371)
)

var MethodsAccount = struct {
	Constructor           abi.MethodNum
	PubkeyAddress         abi.MethodNum
	AuthenticateMessage   abi.MethodNum
	UniversalReceiverHook abi.MethodNum
}{MethodConstructor, 2, 3, UniversalReceiverHookMethodNum}

var MethodsInit = struct {
	Constructor abi.MethodNum
	Exec        abi.MethodNum
}{MethodConstructor, 2}

var MethodsCron = struct {
	Constructor abi.MethodNum
	EpochTick   abi.MethodNum
}{MethodConstructor, 2}

var MethodsReward = struct {
	Constructor      abi.MethodNum
	AwardBlockReward abi.MethodNum
	ThisEpochReward  abi.MethodNum
	UpdateNetworkKPI abi.MethodNum
}{MethodConstructor, 2, 3, 4}

var MethodsMultisig = struct {
	Constructor                 abi.MethodNum
	Propose                     abi.MethodNum
	Approve                     abi.MethodNum
	Cancel                      abi.MethodNum
	AddSigner                   abi.MethodNum
	RemoveSigner                abi.MethodNum
	SwapSigner                  abi.MethodNum
	ChangeNumApprovalsThreshold abi.MethodNum
	LockBalance                 abi.MethodNum
	UniversalReceiverHook       abi.MethodNum
}{MethodConstructor, 2, 3, 4, 5, 6, 7, 8, 9, UniversalReceiverHookMethodNum}

var MethodsPaych = struct {
	Constructor        abi.MethodNum
	UpdateChannelState abi.MethodNum
	Settle             abi.MethodNum
	Collect            abi.MethodNum
}{MethodConstructor, 2, 3, 4}

var MethodsMarket = struct {
	Constructor              abi.MethodNum
	AddBalance               abi.MethodNum
	WithdrawBalance          abi.MethodNum
	PublishStorageDeals      abi.MethodNum
	VerifyDealsForActivation abi.MethodNum
	ActivateDeals            abi.MethodNum
	OnMinerSectorsTerminate  abi.MethodNum
	ComputeDataCommitment    abi.MethodNum
	CronTick                 abi.MethodNum
}{MethodConstructor, 2, 3, 4, 5, 6, 7, 8, 9}

var MethodsPower = struct {
	Constructor              abi.MethodNum
	CreateMiner              abi.MethodNum
	UpdateClaimedPower       abi.MethodNum
	EnrollCronEvent          abi.MethodNum
	CronTick                 abi.MethodNum
	UpdatePledgeTotal        abi.MethodNum
	Deprecated1              abi.MethodNum
	SubmitPoRepForBulkVerify abi.MethodNum
	CurrentTotalPower        abi.MethodNum
}{MethodConstructor, 2, 3, 4, 5, 6, 7, 8, 9}

var MethodsMiner = struct {
	Constructor              abi.MethodNum
	ControlAddresses         abi.MethodNum
	ChangeWorkerAddress      abi.MethodNum
	ChangePeerID             abi.MethodNum
	SubmitWindowedPoSt       abi.MethodNum
	PreCommitSector          abi.MethodNum
	ProveCommitSector        abi.MethodNum
	ExtendSectorExpiration   abi.MethodNum
	TerminateSectors         abi.MethodNum
	DeclareFaults            abi.MethodNum
	DeclareFaultsRecovered   abi.MethodNum
	OnDeferredCronEvent      abi.MethodNum
	CheckSectorProven        abi.MethodNum
	ApplyRewards             abi.MethodNum
	ReportConsensusFault     abi.MethodNum
	WithdrawBalance          abi.MethodNum
	ConfirmSectorProofsValid abi.MethodNum
	ChangeMultiaddrs         abi.MethodNum
	CompactPartitions        abi.MethodNum
	CompactSectorNumbers     abi.MethodNum
	ConfirmUpdateWorkerKey   abi.MethodNum
	RepayDebt                abi.MethodNum
	ChangeOwnerAddress       abi.MethodNum
	DisputeWindowedPoSt      abi.MethodNum
	PreCommitSectorBatch     abi.MethodNum
	ProveCommitAggregate     abi.MethodNum
	ProveReplicaUpdates      abi.MethodNum
	PreCommitSectorBatch2    abi.MethodNum
	ProveReplicaUpdates2     abi.MethodNum
	ChangeBeneficiary        abi.MethodNum
	GetBeneficiary           abi.MethodNum
	ExtendSectorExpiration2  abi.MethodNum
}{MethodConstructor, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

var MethodsVerifiedRegistry = struct {
	Constructor                 abi.MethodNum
	AddVerifier                 abi.MethodNum
	RemoveVerifier              abi.MethodNum
	AddVerifiedClient           abi.MethodNum
	Deprecated1                 abi.MethodNum
	Deprecated2                 abi.MethodNum
	RemoveVerifiedClientDataCap abi.MethodNum
	RemoveExpiredAllocations    abi.MethodNum
	ClaimAllocations            abi.MethodNum
	GetClaims                   abi.MethodNum
	ExtendClaimTerms            abi.MethodNum
	RemoveExpiredClaims         abi.MethodNum
	UniversalReceiverHook       abi.MethodNum
}{MethodConstructor, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 123 /* FIXME */}

var MethodsDatacap = struct {
	Constructor       abi.MethodNum
	Mint              abi.MethodNum
	Destroy           abi.MethodNum
	Name              abi.MethodNum
	Symbol            abi.MethodNum
	TotalSupply       abi.MethodNum
	BalanceOf         abi.MethodNum
	Transfer          abi.MethodNum
	TransferFrom      abi.MethodNum
	IncreaseAllowance abi.MethodNum
	DecreaseAllowance abi.MethodNum
	RevokeAllowance   abi.MethodNum
	Burn              abi.MethodNum
	BurnFrom          abi.MethodNum
}{MethodConstructor, 2, 3, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
