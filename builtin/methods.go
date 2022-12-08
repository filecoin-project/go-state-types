package builtin

import (
	"encoding/binary"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/minio/blake2b-simd"
	"golang.org/x/xerrors"
	"unicode"
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
}{MethodConstructor, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, UniversalReceiverHookMethodNum}

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
	Allowance         abi.MethodNum
}{MethodConstructor, 2, 3, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21}

var MethodsEVM = struct {
	Constructor            abi.MethodNum
	InvokeContract         abi.MethodNum
	GetBytecode            abi.MethodNum
	GetStorageAt           abi.MethodNum
	InvokeContractReadOnly abi.MethodNum
	InvokeContractDelegate abi.MethodNum
}{MethodConstructor, 2, 3, 4, 5, 6}

var MethodsEAM = struct {
	Constructor abi.MethodNum
	Create      abi.MethodNum
	Create2     abi.MethodNum
}{MethodConstructor, 2, 3}

// Generates a standard FRC-42 compliant method number
// Reference: https://github.com/filecoin-project/FIPs/blob/master/FRCs/frc-0042.md
func GenerateMethodNum(name string) (abi.MethodNum, error) {
	err := validateMethodName(name)
	if err != nil {
		return 0, err
	}

	if name == "Constructor" {
		return MethodConstructor, nil
	}

	digest := blake2b.Sum512([]byte("1|" + name))

	for i := 0; i < 64; i += 4 {
		methodId := binary.BigEndian.Uint32(digest[i : i+4])
		if methodId >= (1 << 24) {
			return abi.MethodNum(methodId), nil
		}
	}

	return abi.MethodNum(0), xerrors.Errorf("Could not generate method num from method name :", name)
}

func validateMethodName(name string) error {
	if name == "" {
		return xerrors.Errorf("empty name string")
	}

	if !unicode.IsUpper(rune(name[0])) {
		return xerrors.Errorf("Method name first letter must be uppercase, method name: ", name)
	}

	for _, c := range name {
		if !(unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_') {
			return xerrors.Errorf("method name has illegal characters, method name: ", name)
		}
	}

	return nil
}
