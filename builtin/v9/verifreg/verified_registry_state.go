package verifreg

import (
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v9/util/adt"

	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

// DataCap is an integer number of bytes.
// We can introduce policy changes and replace this in the future.
type DataCap = abi.StoragePower

var DatacapGranularity = builtin.TokenPrecision

func DataCapToTokens(d DataCap) abi.TokenAmount {
	return big.Mul(d, DatacapGranularity)
}

func TokensToDatacap(t abi.TokenAmount) DataCap {
	return big.Div(t, DatacapGranularity)
}

const SignatureDomainSeparation_RemoveDataCap = "fil_removedatacap:"

type RmDcProposalID struct {
	ProposalID uint64
}

type State struct {
	// Root key holder multisig.
	// Authorize and remove verifiers.
	RootKey addr.Address

	// Verifiers authorize VerifiedClients.
	// Verifiers delegate their DataCap.
	Verifiers cid.Cid // HAMT[addr.Address]DataCap

	// RemoveDataCapProposalIDs keeps the counters of the datacap removal proposal a verifier has submitted for a
	//specific client. Unique proposal ids ensure that removal proposals cannot be replayed.√
	// AddrPairKey is constructed as <verifier address, client address>, both using ID addresses.
	RemoveDataCapProposalIDs cid.Cid // HAMT[AddrPairKey]RmDcProposalID

	// Maps client IDs to allocations made by that client.
	Allocations cid.Cid // HAMT[ActorID]HAMT[AllocationID]Allocation

	// Next allocation identifier to use.
	// The value 0 is reserved to mean "no allocation".
	NextAllocationId uint64

	// Maps provider IDs to allocations claimed by that provider.
	Claims cid.Cid // HAMT[ActorID]HAMT[ClaimID]Claim
}

var MinVerifiedDealSize = abi.NewStoragePower(1 << 20)

// rootKeyAddress comes from genesis.
func ConstructState(store adt.Store, rootKeyAddress addr.Address) (*State, error) {
	emptyMapCid, err := adt.StoreEmptyMap(store, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to create empty map: %w", err)
	}

	return &State{
		RootKey:                  rootKeyAddress,
		Verifiers:                emptyMapCid,
		RemoveDataCapProposalIDs: emptyMapCid,
		Allocations:              emptyMapCid,
		NextAllocationId:         1,
		Claims:                   emptyMapCid,
	}, nil
}
