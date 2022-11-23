package verifreg

import (
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
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
	NextAllocationId AllocationId

	// Maps provider IDs to allocations claimed by that provider.
	Claims cid.Cid // HAMT[ActorID]HAMT[ClaimID]Claim
}

var MinVerifiedDealSize = abi.NewStoragePower(1 << 20)

// rootKeyAddress comes from genesis.
func ConstructState(store adt.Store, rootKeyAddress address.Address) (*State, error) {
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

func (st *State) FindAllocation(store adt.Store, clientIdAddr address.Address, allocationId AllocationId) (*Allocation, bool, error) {
	if clientIdAddr.Protocol() != address.ID {
		return nil, false, xerrors.Errorf("can only look up ID addresses")
	}

	innerHamtCid, err := getInnerHamtCid(store, abi.IdAddrKey(clientIdAddr), st.Allocations)
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

	clientId, err := address.IDFromAddress(clientIdAddr)
	if err != nil {
		return nil, false, xerrors.Errorf("couldn't get ID from clientIdAddr: %s", clientIdAddr)
	}

	if uint64(allocation.Client) != clientId {
		return nil, false, xerrors.Errorf("clientId: %d did not match client in allocation: %d", clientId, allocation.Client)
	}

	return &allocation, true, nil
}

func (st *State) FindClaim(store adt.Store, providerIdAddr address.Address, claimId ClaimId) (*Claim, bool, error) {
	if providerIdAddr.Protocol() != address.ID {
		return nil, false, xerrors.Errorf("can only look up ID addresses")
	}

	innerHamtCid, err := getInnerHamtCid(store, abi.IdAddrKey(providerIdAddr), st.Claims)
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

	providerId, err := address.IDFromAddress(providerIdAddr)
	if err != nil {
		return nil, false, xerrors.Errorf("couldn't get ID from providerIdAddr: %s", providerIdAddr)
	}

	if uint64(claim.Provider) != providerId {
		return nil, false, xerrors.Errorf("providerId: %d did not match provider in claim: %d", providerId, claim.Provider)
	}

	return &claim, true, nil
}

func getInnerHamtCid(store adt.Store, key abi.Keyer, mapCid cid.Cid) (cid.Cid, error) {
	actorToHamtMap, err := adt.AsMap(store, mapCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("couldn't get outer map: %x", err)
	}

	var innerHamtCid cbg.CborCid
	if found, err := actorToHamtMap.Get(key, &innerHamtCid); err != nil {
		return cid.Undef, xerrors.Errorf("looking up key: %s: %w", key, err)
	} else if !found {
		return cid.Undef, xerrors.Errorf("did not find key: %s", key)
	}

	return cid.Cid(innerHamtCid), nil
}

func (st *State) LoadAllocationsToMap(store adt.Store, clientIdAddr address.Address) (map[AllocationId]Allocation, error) {
	if clientIdAddr.Protocol() != address.ID {
		return nil, xerrors.Errorf("can only look up ID addresses")
	}

	innerHamtCid, err := getInnerHamtCid(store, abi.IdAddrKey(clientIdAddr), st.Allocations)
	if err != nil {
		return nil, err
	}

	adtMap, err := adt.AsMap(store, innerHamtCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("couldn't get map: %x", err)
	}

	var allocIdToAlloc = make(map[AllocationId]Allocation)
	var out Allocation
	err = adtMap.ForEach(&out, func(key string) error {
		uintKey, err := abi.ParseUIntKey(key)
		if err != nil {
			return xerrors.Errorf("couldn't parse key to uint: %x", err)
		}
		allocIdToAlloc[AllocationId(uintKey)] = out
		return nil
	})
	if err != nil {
		return nil, err
	}

	return allocIdToAlloc, nil
}

func (st *State) LoadClaimsToMap(store adt.Store, providerIdAddr address.Address) (map[ClaimId]Claim, error) {
	if providerIdAddr.Protocol() != address.ID {
		return nil, xerrors.Errorf("can only look up ID addresses")
	}

	innerHamtCid, err := getInnerHamtCid(store, abi.IdAddrKey(providerIdAddr), st.Claims)
	if err != nil {
		return nil, err
	}

	adtMap, err := adt.AsMap(store, innerHamtCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("couldn't get map: %x", err)
	}

	var claimIdToClaim = make(map[ClaimId]Claim)
	var out Claim
	err = adtMap.ForEach(&out, func(key string) error {
		uintKey, err := abi.ParseUIntKey(key)
		if err != nil {
			return xerrors.Errorf("couldn't parse key to uint: %w", err)
		}
		claimIdToClaim[ClaimId(uintKey)] = out
		return nil
	})
	if err != nil {
		return nil, err
	}

	return claimIdToClaim, nil
}

func (st *State) GetAllClaims(store adt.Store) (map[ClaimId]Claim, error) {
	allClaims := make(map[ClaimId]Claim)

	actorToHamtMap, err := adt.AsMap(store, st.Claims, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("couldn't get outer map: %x", err)
	}

	var innerHamtCid cbg.CborCid
	err = actorToHamtMap.ForEach(&innerHamtCid, func(idKey string) error {
		innerMap, err := adt.AsMap(store, cid.Cid(innerHamtCid), builtin.DefaultHamtBitwidth)
		if err != nil {
			return xerrors.Errorf("couldn't get inner map: %x", err)
		}

		var out Claim
		err = innerMap.ForEach(&out, func(key string) error {
			uintKey, err := abi.ParseUIntKey(key)
			if err != nil {
				return xerrors.Errorf("couldn't parse idKey to uint: %w", err)
			}
			allClaims[ClaimId(uintKey)] = out
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return allClaims, nil
}

func (st *State) GetAllAllocations(store adt.Store) (map[AllocationId]Allocation, error) {
	allAllocations := make(map[AllocationId]Allocation)

	actorToHamtMap, err := adt.AsMap(store, st.Allocations, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("couldn't get outer map: %x", err)
	}

	var innerHamtCid cbg.CborCid
	err = actorToHamtMap.ForEach(&innerHamtCid, func(idKey string) error {
		innerMap, err := adt.AsMap(store, cid.Cid(innerHamtCid), builtin.DefaultHamtBitwidth)
		if err != nil {
			return xerrors.Errorf("couldn't get inner map: %x", err)
		}

		var out Allocation
		err = innerMap.ForEach(&out, func(key string) error {
			uintKey, err := abi.ParseUIntKey(key)
			if err != nil {
				return xerrors.Errorf("couldn't parse idKey to uint: %w", err)
			}
			allAllocations[AllocationId(uintKey)] = out
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return allAllocations, nil
}
