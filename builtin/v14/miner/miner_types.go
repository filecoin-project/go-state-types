package miner

import (
	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/xerrors"

	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v14/power"
	"github.com/filecoin-project/go-state-types/builtin/v14/util/adt"
	"github.com/filecoin-project/go-state-types/builtin/v14/util/smoothing"
	"github.com/filecoin-project/go-state-types/builtin/v14/verifreg"
	xc "github.com/filecoin-project/go-state-types/exitcode"
	"github.com/filecoin-project/go-state-types/proof"
)

type DeclareFaultsRecoveredParams struct {
	Recoveries []RecoveryDeclaration
}

type RecoveryDeclaration struct {
	// The deadline to which the recovered sectors are assigned, in range [0..WPoStPeriodDeadlines)
	Deadline uint64
	// Partition index within the deadline containing the recovered sectors.
	Partition uint64
	// Sectors in the partition being declared recovered.
	Sectors bitfield.BitField
}

type DeclareFaultsParams struct {
	Faults []FaultDeclaration
}

type FaultDeclaration struct {
	// The deadline to which the faulty sectors are assigned, in range [0..WPoStPeriodDeadlines)
	Deadline uint64
	// Partition index within the deadline containing the faulty sectors.
	Partition uint64
	// Sectors in the partition being declared faulty.
	Sectors bitfield.BitField
}

type ReplicaUpdate struct {
	SectorID           abi.SectorNumber
	Deadline           uint64
	Partition          uint64
	NewSealedSectorCID cid.Cid `checked:"true"`
	Deals              []abi.DealID
	UpdateProofType    abi.RegisteredUpdateProof
	ReplicaProof       []byte
}

type ProveReplicaUpdatesParams struct {
	Updates []ReplicaUpdate
}

type ReplicaUpdate2 struct {
	SectorID             abi.SectorNumber
	Deadline             uint64
	Partition            uint64
	NewSealedSectorCID   cid.Cid `checked:"true"`
	NewUnsealedSectorCID cid.Cid `checked:"true"`
	Deals                []abi.DealID
	UpdateProofType      abi.RegisteredUpdateProof
	ReplicaProof         []byte
}

type ProveReplicaUpdatesParams2 struct {
	Updates []ReplicaUpdate2
}

type PoStPartition struct {
	// Partitions are numbered per-deadline, from zero.
	Index uint64
	// Sectors skipped while proving that weren't already declared faulty
	Skipped bitfield.BitField
}

// Information submitted by a miner to provide a Window PoSt.
type SubmitWindowedPoStParams struct {
	// The deadline index which the submission targets.
	Deadline uint64
	// The partitions being proven.
	Partitions []PoStPartition
	// Array of proofs, one per distinct registered proof type present in the sectors being proven.
	// In the usual case of a single proof type, this array will always have a single element (independent of number of partitions).
	Proofs []proof.PoStProof
	// The epoch at which these proofs is being committed to a particular chain.
	ChainCommitEpoch abi.ChainEpoch
	// The ticket randomness on the chain at the ChainCommitEpoch on the chain this post is committed to
	ChainCommitRand abi.Randomness
}

type DisputeWindowedPoStParams struct {
	Deadline  uint64
	PoStIndex uint64 // only one is allowed at a time to avoid loading too many sector infos.
}

type ProveCommitAggregateParams struct {
	SectorNumbers  bitfield.BitField
	AggregateProof []byte
}

type ProveCommitSectorParams struct {
	SectorNumber abi.SectorNumber
	Proof        []byte
}

type MinerConstructorParams = power.MinerConstructorParams

type TerminateSectorsParams struct {
	Terminations []TerminationDeclaration
}

type TerminationDeclaration struct {
	Deadline  uint64
	Partition uint64
	Sectors   bitfield.BitField
}

type TerminateSectorsReturn struct {
	// Set to true if all early termination work has been completed. When
	// false, the miner may choose to repeatedly invoke TerminateSectors
	// with no new sectors to process the remainder of the pending
	// terminations. While pending terminations are outstanding, the miner
	// will not be able to withdraw funds.
	Done bool
}

type ChangePeerIDParams struct {
	NewID abi.PeerID
}

type ChangeMultiaddrsParams struct {
	NewMultiaddrs []abi.Multiaddrs
}

type ChangeWorkerAddressParams struct {
	NewWorker       addr.Address
	NewControlAddrs []addr.Address
}

type ExtendSectorExpirationParams struct {
	Extensions []ExpirationExtension
}

type ExpirationExtension struct {
	Deadline      uint64
	Partition     uint64
	Sectors       bitfield.BitField
	NewExpiration abi.ChainEpoch
}

type ReportConsensusFaultParams struct {
	BlockHeader1     []byte
	BlockHeader2     []byte
	BlockHeaderExtra []byte
}

type GetControlAddressesReturn struct {
	Owner        addr.Address
	Worker       addr.Address
	ControlAddrs []addr.Address
}

type CheckSectorProvenParams struct {
	SectorNumber abi.SectorNumber
}

type WithdrawBalanceParams struct {
	AmountRequested abi.TokenAmount
}

type CompactPartitionsParams struct {
	Deadline   uint64
	Partitions bitfield.BitField
}

type CompactSectorNumbersParams struct {
	MaskSectorNumbers bitfield.BitField
}

type CronEventType int64

const (
	CronEventWorkerKeyChange CronEventType = iota
	CronEventProvingDeadline
	CronEventProcessEarlyTerminations
)

type CronEventPayload struct {
	EventType CronEventType
}

// Identifier for a single partition within a miner.
type PartitionKey struct {
	Deadline  uint64
	Partition uint64
}

type PreCommitSectorParams struct {
	SealProof       abi.RegisteredSealProof
	SectorNumber    abi.SectorNumber
	SealedCID       cid.Cid `checked:"true"` // CommR
	SealRandEpoch   abi.ChainEpoch
	DealIDs         []abi.DealID
	Expiration      abi.ChainEpoch
	ReplaceCapacity bool // DEPRECATED: Whether to replace a "committed capacity" no-deal sector (requires non-empty DealIDs)
	// DEPRECATED: The committed capacity sector to replace, and it's deadline/partition location
	ReplaceSectorDeadline  uint64
	ReplaceSectorPartition uint64
	ReplaceSectorNumber    abi.SectorNumber
}

type PreCommitSectorBatchParams struct {
	Sectors []PreCommitSectorParams
}

type PreCommitSectorBatchParams2 struct {
	Sectors []SectorPreCommitInfo
}

type ChangeBeneficiaryParams struct {
	NewBeneficiary addr.Address
	NewQuota       abi.TokenAmount
	NewExpiration  abi.ChainEpoch
}

type ActiveBeneficiary struct {
	Beneficiary addr.Address
	Term        BeneficiaryTerm
}

type GetBeneficiaryReturn struct {
	Active   ActiveBeneficiary
	Proposed *PendingBeneficiaryChange
}

// ExpirationSet is a collection of sector numbers that are expiring, either due to
// expected "on-time" expiration at the end of their life, or unexpected "early" termination
// due to being faulty for too long consecutively.
// Note that there is not a direct correspondence between on-time sectors and active power;
// a sector may be faulty but expiring on-time if it faults just prior to expected termination.
// Early sectors are always faulty, and active power always represents on-time sectors.
type ExpirationSet struct {
	OnTimeSectors bitfield.BitField // Sectors expiring "on time" at the end of their committed life
	EarlySectors  bitfield.BitField // Sectors expiring "early" due to being faulty for too long
	OnTimePledge  abi.TokenAmount   // Pledge total for the on-time sectors
	ActivePower   PowerPair         // Power that is currently active (not faulty)
	FaultyPower   PowerPair         // Power that is currently faulty
}

// A queue of expiration sets by epoch, representing the on-time or early termination epoch for a collection of sectors.
// Wraps an AMT[ChainEpoch]*ExpirationSet.
// Keys in the queue are quantized (upwards), modulo some offset, to reduce the cardinality of keys.
type ExpirationQueue struct {
	*adt.Array
	quant builtin.QuantSpec
}

// Loads a queue root.
// Epochs provided to subsequent method calls will be quantized upwards to quanta mod offsetSeed before being
// written to/read from queue entries.
func LoadExpirationQueue(store adt.Store, root cid.Cid, quant builtin.QuantSpec, bitwidth int) (ExpirationQueue, error) {
	arr, err := adt.AsArray(store, root, bitwidth)
	if err != nil {
		return ExpirationQueue{}, xerrors.Errorf("failed to load epoch queue %v: %w", root, err)
	}
	return ExpirationQueue{arr, quant}, nil
}
func LoadSectors(store adt.Store, root cid.Cid) (Sectors, error) {
	sectorsArr, err := adt.AsArray(store, root, SectorsAmtBitwidth)
	if err != nil {
		return Sectors{}, err
	}
	return Sectors{sectorsArr}, nil
}

// Sectors is a helper type for accessing/modifying a miner's sectors. It's safe
// to pass this object around as needed.
type Sectors struct {
	*adt.Array
}

func (sa Sectors) Load(sectorNos bitfield.BitField) ([]*SectorOnChainInfo, error) {
	var sectorInfos []*SectorOnChainInfo
	if err := sectorNos.ForEach(func(i uint64) error {
		var sectorOnChain SectorOnChainInfo
		found, err := sa.Array.Get(i, &sectorOnChain)
		if err != nil {
			return xc.ErrIllegalState.Wrapf("failed to load sector %v: %w", abi.SectorNumber(i), err)
		} else if !found {
			return xc.ErrNotFound.Wrapf("can't find sector %d", i)
		}
		sectorInfos = append(sectorInfos, &sectorOnChain)
		return nil
	}); err != nil {
		// Keep the underlying error code, unless the error was from
		// traversing the bitfield. In that case, it's an illegal
		// argument error.
		return nil, xc.Unwrap(err, xc.ErrIllegalArgument).Wrapf("failed to load sectors: %w", err)
	}
	return sectorInfos, nil
}

func (sa Sectors) Get(sectorNumber abi.SectorNumber) (info *SectorOnChainInfo, found bool, err error) {
	var res SectorOnChainInfo
	if found, err := sa.Array.Get(uint64(sectorNumber), &res); err != nil {
		return nil, false, xerrors.Errorf("failed to get sector %d: %w", sectorNumber, err)
	} else if !found {
		return nil, false, nil
	}
	return &res, true, nil
}

// VestingFunds represents the vesting table state for the miner.
// It is a slice of (VestingEpoch, VestingAmount).
// The slice will always be sorted by the VestingEpoch.
type VestingFunds struct {
	Funds []VestingFund
}

// VestingFund represents miner funds that will vest at the given epoch.
type VestingFund struct {
	Epoch  abi.ChainEpoch
	Amount abi.TokenAmount
}

// ConstructVestingFunds constructs empty VestingFunds state.
func ConstructVestingFunds() *VestingFunds {
	v := new(VestingFunds)
	v.Funds = nil
	return v
}

type DeferredCronEventParams struct {
	EventPayload            []byte
	RewardSmoothed          smoothing.FilterEstimate
	QualityAdjPowerSmoothed smoothing.FilterEstimate
}

type ApplyRewardParams struct {
	Reward  abi.TokenAmount
	Penalty abi.TokenAmount
}

type InternalSectorSetupForPresealParams struct {
	Sectors                 []abi.SectorNumber
	RewardSmoothed          smoothing.FilterEstimate
	RewardBaselinePower     abi.StoragePower
	QualityAdjPowerSmoothed smoothing.FilterEstimate
}

type ExtendSectorExpiration2Params struct {
	Extensions []ExpirationExtension2
}

type ExpirationExtension2 struct {
	Deadline          uint64
	Partition         uint64
	Sectors           bitfield.BitField
	SectorsWithClaims []SectorClaim
	NewExpiration     abi.ChainEpoch
}

type SectorClaim struct {
	SectorNumber   abi.SectorNumber
	MaintainClaims []verifreg.ClaimId
	DropClaims     []verifreg.ClaimId
}

type GetOwnerReturn struct {
	Owner    addr.Address
	Proposed *addr.Address
}

type IsControllingAddressParams = addr.Address

type IsControllingAddressReturn = cbg.CborBool

type GetSectorSizeReturn = abi.SectorSize

type GetAvailableBalanceReturn = abi.TokenAmount

type GetVestingFundsReturn = VestingFunds

type GetPeerIDReturn struct {
	PeerId []byte
}

type GetMultiAddrsReturn struct {
	MultiAddrs []byte
}

// ProveCommitSectors3Params represents the parameters for proving committed sectors.
type ProveCommitSectors3Params struct {
	SectorActivations          []SectorActivationManifest
	SectorProofs               [][]byte
	AggregateProof             []byte
	AggregateProofType         *abi.RegisteredAggregationProof
	RequireActivationSuccess   bool
	RequireNotificationSuccess bool
}

// SectorActivationManifest contains data to activate a commitment to one sector and its data.
// All pieces of data must be specified, whether or not not claiming a FIL+ activation or being
// notified to a data consumer.
// An implicit zero piece fills any remaining sector capacity.
type SectorActivationManifest struct {
	SectorNumber abi.SectorNumber
	Pieces       []PieceActivationManifest
}

// PieceActivationManifest represents the manifest for activating a piece.
type PieceActivationManifest struct {
	CID                   cid.Cid
	Size                  abi.PaddedPieceSize
	VerifiedAllocationKey *VerifiedAllocationKey
	Notify                []DataActivationNotification
}

// VerifiedAllocationKey represents the key for a verified allocation.
type VerifiedAllocationKey struct {
	Client abi.ActorID
	ID     verifreg.AllocationId
}

// DataActivationNotification represents a notification for data activation.
type DataActivationNotification struct {
	Address addr.Address
	Payload []byte
}

// ProveCommitSectors3Return represents the return value for the ProveCommit2 function.
type ProveCommitSectors3Return = BatchReturn

type BatchReturn struct {
	SuccessCount uint64
	FailCodes    []FailCode
}

type FailCode struct {
	// Idx represents the index of the operation that failed within the batch.
	Idx  uint64
	Code xc.ExitCode // todo correct?
}

// ProveReplicaUpdates3Params represents the parameters for proving replica updates.
type ProveReplicaUpdates3Params struct {
	SectorUpdates              []SectorUpdateManifest
	SectorProofs               [][]byte
	AggregateProof             []byte
	UpdateProofsType           abi.RegisteredUpdateProof
	AggregateProofType         *abi.RegisteredAggregationProof
	RequireActivationSuccess   bool
	RequireNotificationSuccess bool
}

// SectorUpdateManifest contains data for sector update.
type SectorUpdateManifest struct {
	Sector       abi.SectorNumber
	Deadline     uint64
	Partition    uint64
	NewSealedCID cid.Cid
	Pieces       []PieceActivationManifest
}

// ProveReplicaUpdates3Return represents the return value for the ProveReplicaUpdates3 function.
type ProveReplicaUpdates3Return = BatchReturn

// SectorContentChangedParams represents a notification of change committed to sectors.
type SectorContentChangedParams = []SectorChanges

// SectorChanges describes changes to one sector's content.
type SectorChanges struct {
	Sector                 abi.SectorNumber
	MinimumCommitmentEpoch abi.ChainEpoch
	Added                  []PieceChange
}

// PieceChange describes a piece of data committed to a sector.
type PieceChange struct {
	Data    cid.Cid
	Size    abi.PaddedPieceSize
	Payload []byte
}

// SectorContentChangedReturn represents the return value for the SectorContentChanged function.
type SectorContentChangedReturn = []SectorReturn

// SectorReturn represents a result for each sector that was notified.
type SectorReturn = []PieceReturn

// PieceReturn represents a result for each piece for the sector that was notified.
type PieceReturn = bool // Accepted = true

// SectorNIActivationInfo is the information needed to activate a sector with a "zero" replica.
type SectorNIActivationInfo struct {
	SealingNumber abi.SectorNumber // Sector number used to generate replica id
	SealerID      abi.ActorID      // Must be set to ID of receiving actor for now
	SealedCID     cid.Cid          // CommR
	SectorNumber  abi.SectorNumber // Unique id of sector in actor state
	SealRandEpoch abi.ChainEpoch
	Expiration    abi.ChainEpoch
}

// ProveCommitSectorsNIParams is the parameters for non-interactive prove committing of sectors
// via the miner actor method ProveCommitSectorsNI.
type ProveCommitSectorsNIParams struct {
	Sectors                  []SectorNIActivationInfo       // Information about sealing of each sector
	AggregateProof           []byte                         // Aggregate proof for all sectors
	SealProofType            abi.RegisteredSealProof        // Proof type for each seal (must be an NI-PoRep variant)
	AggregateProofType       abi.RegisteredAggregationProof // Proof type for aggregation
	ProvingDeadline          uint64                         // The Window PoST deadline index at which to schedule the new sectors
	RequireActivationSuccess bool                           // Whether to abort if any sector activation fails
}

type ProveCommitSectorsNIReturn = BatchReturn
