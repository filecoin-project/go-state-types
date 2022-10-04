package verifreg

import (
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	cbg "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

// DataCap is an integer number of bytes.
// We can introduce policy changes and replace this in the future.
type DataCap = abi.StoragePower

var DataCapGranularity = builtin.TokenPrecision

func DataCapToTokens(d DataCap) abi.TokenAmount {
	return big.Mul(d, DataCapGranularity)
}

func TokensToDatacap(t abi.TokenAmount) DataCap {
	return big.Div(t, DataCapGranularity)
}

const SignatureDomainSeparation_RemoveDataCap = "fil_removedatacap:"

type RmDcProposalID struct {
	ProposalID uint64
}

type State struct {
	// Root key holder multisig.
	// Authorize and remove verifiers.
	RootKey address.Address

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

func (st *State) FindAllocation(store adt.Store, addr address.Address, allocationId AllocationId) (*Allocation, bool, error) {
	if addr.Protocol() != address.ID {
		return nil, false, xerrors.Errorf("can only look up ID addresses")
	}

	innerHamtCid, err := GetInnerHamtCid(store, abi.IdAddrKey(addr), st.Allocations)
	if err != nil {
		return nil, false, err
	}

	idToAllocationMap, err := adt.AsMap(store, innerHamtCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, false, xerrors.Errorf("couldn't get inner map: %x", err)
	}

	var allocation Allocation
	if found, err := idToAllocationMap.Get(allocationId, &allocation); err != nil {
		return nil, false, xerrors.Errorf("looking up allocation ID: %d: %w", allocationId, err)
	} else if !found {
		return nil, false, nil
	}

	return &allocation, true, nil
}

func (st *State) FindClaim(store adt.Store, addr address.Address, claimId ClaimId) (*Claim, bool, error) {
	if addr.Protocol() != address.ID {
		return nil, false, xerrors.Errorf("can only look up ID addresses")
	}

	innerHamtCid, err := GetInnerHamtCid(store, abi.IdAddrKey(addr), st.Claims)
	if err != nil {
		return nil, false, err
	}

	idToClaimsMap, err := adt.AsMap(store, innerHamtCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, false, xerrors.Errorf("couldn't get inner map: %x", err)
	}

	var claim Claim
	if found, err := idToClaimsMap.Get(claimId, &claim); err != nil {
		return nil, false, xerrors.Errorf("looking up allocation ID: %d: %w", claimId, err)
	} else if !found {
		return nil, false, nil
	}

	return &claim, true, nil
}

func GetInnerHamtCid(store adt.Store, addr abi.Keyer, mapCid cid.Cid) (cid.Cid, error) {
	actorToHamtMap, err := adt.AsMap(store, mapCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("couldn't get outer map: %x", err)
	}

	var innerHamtCid cbg.CborCid
	if found, err := actorToHamtMap.Get(addr, &innerHamtCid); err != nil {
		return cid.Undef, xerrors.Errorf("looking up key: %s: %w", addr, err)
	} else if !found {
		return cid.Undef, xerrors.Errorf("did not find key: %s", addr)
	}

	return cid.Cid(innerHamtCid), nil
}
