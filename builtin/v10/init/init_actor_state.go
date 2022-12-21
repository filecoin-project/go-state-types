package init

import (
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
	"github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/xerrors"
)

type State struct {
	AddressMap  cid.Cid // HAMT[addr.Address]abi.ActorID
	NextID      abi.ActorID
	NetworkName string
}

func ConstructState(store adt.Store, networkName string) (*State, error) {
	emptyAddressMapCid, err := adt.StoreEmptyMap(store, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to create empty map: %w", err)
	}

	return &State{
		AddressMap:  emptyAddressMapCid,
		NextID:      abi.ActorID(builtin.FirstNonSingletonActorId),
		NetworkName: networkName,
	}, nil
}

// ResolveAddress resolves an address to an ID-address, if possible.
// If the provided address is an ID address, it is returned as-is.
// This means that mapped ID-addresses (which should only appear as values, not keys) and
// singleton actor addresses (which are not in the map) pass through unchanged.
//
// Returns an ID-address and `true` if the address was already an ID-address or was resolved in the mapping.
// Returns an undefined address and `false` if the address was not an ID-address and not found in the mapping.
// Returns an error only if state was inconsistent.
func (s *State) ResolveAddress(store adt.Store, address addr.Address) (addr.Address, bool, error) {
	// Short-circuit ID address resolution.
	if address.Protocol() == addr.ID {
		return address, true, nil
	}

	// Lookup address.
	m, err := adt.AsMap(store, s.AddressMap, builtin.DefaultHamtBitwidth)
	if err != nil {
		return addr.Undef, false, xerrors.Errorf("failed to load address map: %w", err)
	}

	var actorID cbg.CborInt
	if found, err := m.Get(abi.AddrKey(address), &actorID); err != nil {
		return addr.Undef, false, xerrors.Errorf("failed to get from address map: %w", err)
	} else if found {
		// Reconstruct address from the ActorID.
		idAddr, err := addr.NewIDAddress(uint64(actorID))
		return idAddr, true, err
	} else {
		return addr.Undef, false, nil
	}
}

// Maps argument addresses to to a new or existing actor ID.
// With no delegated address, or if the delegated address is not already mapped,
// allocates a new ID address and maps both to it.
// If the delegated address is already present, maps the robust address to that actor ID.
// Fails if the robust address is already mapped, providing tombstone.
//
//0 Returns the actor ID and a boolean indicating whether or not the actor already exists.
func (s *State) MapAddressToNewID(store adt.Store, robustAddress addr.Address, delegatedAddress *addr.Address) (addr.Address, bool, error) {
	m, err := adt.AsMap(store, s.AddressMap, builtin.DefaultHamtBitwidth)
	if err != nil {
		return addr.Undef, false, xerrors.Errorf("failed to load address map: %w", err)
	}

	var actorID cbg.CborInt
	existing := false
	if delegatedAddress != nil {
		// If there's a delegated address, either recall the already-mapped actor ID or
		// create and map a new one.
		existing, err = m.Get(abi.AddrKey(*delegatedAddress), &actorID)
		if err != nil {
			return addr.Undef, false, xerrors.Errorf("failed to check existing delegated addr: %w", err)
		}
		if !existing {
			actorID = cbg.CborInt(s.NextID)
			s.NextID++
			if err := m.Put(abi.AddrKey(*delegatedAddress), &actorID); err != nil {
				return addr.Undef, false, xerrors.Errorf("failed to put new delegated addr: %w", err)
			}
		}
	} else {
		// With no delegated address, always create a new actor ID.
		actorID = cbg.CborInt(s.NextID)
		s.NextID++
	}

	isNew, err := m.PutIfAbsent(abi.AddrKey(robustAddress), &actorID)
	if err != nil {
		return addr.Undef, false, xerrors.Errorf("map address failed to store entry: %w", err)
	}
	if !isNew {
		return addr.Undef, false, xerrors.Errorf("robust address %s is already allocated in the address map", robustAddress)
	}

	amr, err := m.Root()
	if err != nil {
		return addr.Undef, false, xerrors.Errorf("failed to get address map root: %w", err)
	}
	s.AddressMap = amr

	idAddr, err := addr.NewIDAddress(uint64(actorID))
	if err != nil {
		return addr.Undef, false, xerrors.Errorf("failed to convert actorID to address: %w", err)
	}

	return idAddr, existing, nil
}
