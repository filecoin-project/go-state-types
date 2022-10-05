package verifreg

import (
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/exitcode"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-varint"
)

// RemoveDataCapProposal A verifier who wants to send/agree to a RemoveDataCapRequest should sign a RemoveDataCapProposal and send the signed proposal to the root key holder.
type RemoveDataCapProposal struct {
	// VerifiedClient is the client address to remove the DataCap from
	// The address must be an ID address
	VerifiedClient addr.Address
	// DataCapAmount is the amount of DataCap to be removed from the VerifiedClient address
	DataCapAmount DataCap
	// RemovalProposalID is the counter of the proposal sent by the Verifier for the VerifiedClient
	RemovalProposalID RmDcProposalID
}

// RemoveDataCapRequest A verifier who wants to submit a request should send their RemoveDataCapRequest to the RKH.
type RemoveDataCapRequest struct {
	// Verifier is the verifier address used for VerifierSignature.
	// The address can be address.SECP256K1 or address.BLS
	Verifier addr.Address
	// VerifierSignature is the Verifier's signature over a RemoveDataCapProposal
	VerifierSignature crypto.Signature
}

type AddVerifierParams struct {
	Address   addr.Address
	Allowance DataCap
}

type AddVerifiedClientParams struct {
	Address   addr.Address
	Allowance DataCap
}

type UseBytesParams struct {
	// Address of verified client.
	Address addr.Address
	// Number of bytes to use.
	DealSize abi.StoragePower
}

type RestoreBytesParams struct {
	Address  addr.Address
	DealSize abi.StoragePower
}

type RemoveDataCapParams struct {
	VerifiedClientToRemove addr.Address
	DataCapAmountToRemove  DataCap
	VerifierRequest1       RemoveDataCapRequest
	VerifierRequest2       RemoveDataCapRequest
}

type RemoveDataCapReturn struct {
	VerifiedClient addr.Address
	DataCapRemoved DataCap
}

type RemoveExpiredAllocationsParams struct {
	Client        abi.ActorID
	AllocationIds []AllocationId
}

type RemoveExpiredAllocationsReturn struct {
	Considered       []AllocationId
	Results          BatchReturn
	DataCapRecovered DataCap
}

type BatchReturn struct {
	SuccessCount uint64
	FailCodes    []FailCode
}

type FailCode struct {
	Idx  uint64
	Code exitcode.ExitCode
}

type AllocationId uint64

func (a AllocationId) Key() string {
	return string(varint.ToUvarint(uint64(a)))
}

type ClaimId uint64

func (a ClaimId) Key() string {
	return string(varint.ToUvarint(uint64(a)))
}

type ClaimAllocationsParams struct {
	Sectors      []SectorAllocationClaim
	AllOrNothing bool
}

type SectorAllocationClaim struct {
	Client       abi.ActorID
	AllocationId AllocationId
	Data         cid.Cid
	Size         abi.PaddedPieceSize
	Sector       abi.SectorNumber
	SectorExpiry abi.ChainEpoch
}

type ClaimAllocationsReturn struct {
	BatchInfo    BatchReturn
	ClaimedSpace big.Int
}

type GetClaimsParams struct {
	Provider abi.ActorID
	ClaimIds []ClaimId
}

type GetClaimsReturn struct {
	BatchInfo BatchReturn
	Claims    []Claim
}

type Claim struct {
	// The provider storing the data (from allocation).
	Provider abi.ActorID
	// The client which allocated the DataCap (from allocation).
	Client abi.ActorID
	// Identifier of the data committed (from allocation).
	Data cid.Cid
	// The (padded) size of data (from allocation).
	Size abi.PaddedPieceSize
	// The min period which the provider must commit to storing data
	TermMin abi.ChainEpoch
	// The max period for which provider can earn QA-power for the data
	TermMax abi.ChainEpoch
	// The epoch at which the (first range of the) piece was committed.
	TermStart abi.ChainEpoch
	// ID of the provider's sector in which the data is committed.
	Sector abi.SectorNumber
}

type Allocation struct {
	// The verified client which allocated the DataCap.
	Client abi.ActorID
	// The provider (miner actor) which may claim the allocation.
	Provider abi.ActorID
	// Identifier of the data to be committed.
	Data cid.Cid
	// The (padded) size of data.
	Size abi.PaddedPieceSize
	// The minimum duration which the provider must commit to storing the piece to avoid
	// early-termination penalties (epochs).
	TermMin abi.ChainEpoch
	// The maximum period for which a provider can earn quality-adjusted power
	// for the piece (epochs).
	TermMax abi.ChainEpoch
	// The latest epoch by which a provider must commit data before the allocation expires.
	Expiration abi.ChainEpoch
}

type UniversalReceiverParams struct {
	Type_   ReceiverType
	Payload []byte
}

type ReceiverType uint64

type AllocationsResponse struct {
	AllocationResults BatchReturn
	ExtensionResults  BatchReturn
	NewAllocations    []AllocationId
}

type ExtendClaimTermsParams struct {
	Terms []ClaimTerm
}

type ClaimTerm struct {
	Provider abi.ActorID
	ClaimId  ClaimId
	TermMax  abi.ChainEpoch
}

type ExtendClaimTermsReturn BatchReturn

type RemoveExpiredClaimsParams struct {
	Provider abi.ActorID
	ClaimIds []ClaimId
}

type RemoveExpiredClaimsReturn struct {
	Considered []AllocationId
	Results    BatchReturn
}

type AllocationRequest struct {
	// The provider (miner actor) which may claim the allocation.
	Provider abi.ActorID
	// Identifier of the data to be committed.
	Data cid.Cid
	// The (padded) size of data.
	Size abi.PaddedPieceSize
	// The minimum duration which the provider must commit to storing the piece to avoid
	// early-termination penalties (epochs).
	TermMin abi.ChainEpoch
	// The maximum period for which a provider can earn quality-adjusted power
	// for the piece (epochs).
	TermMax abi.ChainEpoch
	// The latest epoch by which a provider must commit data before the allocation expires.
	Expiration abi.ChainEpoch
}

type ClaimExtensionRequest struct {
	Provider addr.Address
	Claim    ClaimId
	TermMax  abi.ChainEpoch
}

type AllocationRequests struct {
	Allocations []AllocationRequest
	Extensions  []ClaimExtensionRequest
}
