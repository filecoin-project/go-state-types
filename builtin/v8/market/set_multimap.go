package market

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
	cbg "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/xerrors"

	cid "github.com/ipfs/go-cid"
)

type SetMultimap struct {
	mp            *adt.Map
	store         adt.Store
	innerBitwidth int
}

// Interprets a store as a HAMT-based map of HAMT-based sets with root `r`.
// Both inner and outer HAMTs are interpreted with branching factor 2^bitwidth.
func AsSetMultimap(s adt.Store, r cid.Cid, outerBitwidth, innerBitwidth int) (*SetMultimap, error) {
	m, err := adt.AsMap(s, r, outerBitwidth)
	if err != nil {
		return nil, err
	}
	return &SetMultimap{mp: m, store: s, innerBitwidth: innerBitwidth}, nil
}

// Creates a new map backed by an empty HAMT and flushes it to the store.
// Both inner and outer HAMTs have branching factor 2^bitwidth.
func MakeEmptySetMultimap(s adt.Store, bitwidth int) (*SetMultimap, error) {
	m, err := adt.MakeEmptyMap(s, bitwidth)
	if err != nil {
		return nil, err
	}
	return &SetMultimap{mp: m, store: s, innerBitwidth: bitwidth}, nil
}

// Writes a new empty map to the store and returns its CID.
func StoreEmptySetMultimap(s adt.Store, bitwidth int) (cid.Cid, error) {
	mm, err := MakeEmptySetMultimap(s, bitwidth)
	if err != nil {
		return cid.Undef, err
	}
	return mm.Root()
}

// Returns the root cid of the underlying HAMT.
func (mm *SetMultimap) Root() (cid.Cid, error) {
	return mm.mp.Root()
}

func parseDealKey(s string) (abi.DealID, error) {
	key, err := abi.ParseUIntKey(s)
	return abi.DealID(key), err
}

func (mm *SetMultimap) get(key abi.Keyer) (*adt.Set, bool, error) {
	var setRoot cbg.CborCid
	found, err := mm.mp.Get(key, &setRoot)
	if err != nil {
		return nil, false, xerrors.Errorf("failed to load set key: %v: %w", key, err)
	}
	var set *adt.Set
	if found {
		set, err = adt.AsSet(mm.store, cid.Cid(setRoot), mm.innerBitwidth)
		if err != nil {
			return nil, false, err
		}
	}
	return set, found, nil
}

// Iterates all entries for a key, iteration halts if the function returns an error.
func (mm *SetMultimap) ForEach(epoch abi.ChainEpoch, fn func(id abi.DealID) error) error {
	set, found, err := mm.get(abi.UIntKey(uint64(epoch)))
	if err != nil {
		return err
	}
	if found {
		return set.ForEach(func(k string) error {
			v, err := parseDealKey(k)
			if err != nil {
				return err
			}
			return fn(v)
		})
	}
	return nil
}
