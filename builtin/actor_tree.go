package builtin

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

// Value type of the top level of the state tree.
// Represents the on-chain state of a single actor.
// This is the actor state for state tree version up to 4
type ActorV4 struct {
	Code       cid.Cid // CID representing the code associated with the actor
	Head       cid.Cid // CID of the head state object for the actor
	CallSeqNum uint64  // CallSeqNum for the next message to be received by the actor (non-zero for accounts only)
	Balance    big.Int // Token balance of the actor
}

// As above, but this is the actor state for state tree version 5 and above.
type ActorV5 struct {
	Code       cid.Cid          // CID representing the code associated with the actor
	Head       cid.Cid          // CID of the head state object for the actor
	CallSeqNum uint64           // CallSeqNum for the next message to be received by the actor (non-zero for accounts only)
	Balance    big.Int          // Token balance of the actor
	Address    *address.Address // Delegated (f4) actor address
}

// A specialization of a map of ID-addresses to actor heads.
type ActorTree struct {
	Map   *adt.Map
	Store adt.Store
}

// Initializes a new, empty state tree backed by a store.
func NewTree(store adt.Store) (*ActorTree, error) {
	emptyMap, err := adt.MakeEmptyMap(store, DefaultHamtBitwidth)
	if err != nil {
		return nil, err
	}
	return &ActorTree{
		Map:   emptyMap,
		Store: store,
	}, nil
}

// Loads a tree from a root CID and store.
func LoadTree(s adt.Store, r cid.Cid) (*ActorTree, error) {
	m, err := adt.AsMap(s, r, DefaultHamtBitwidth)
	if err != nil {
		return nil, err
	}
	return &ActorTree{
		Map:   m,
		Store: s,
	}, nil
}

// Writes the tree root node to the store, and returns its CID.
func (t *ActorTree) Flush() (cid.Cid, error) {
	return t.Map.Root()
}

// Loads the state associated with an address.
func (t *ActorTree) GetActorV4(addr address.Address) (*ActorV4, bool, error) {
	if addr.Protocol() != address.ID {
		return nil, false, xerrors.Errorf("non-ID address %v invalid as actor key", addr)
	}
	var actor ActorV4
	found, err := t.Map.Get(abi.AddrKey(addr), &actor)
	return &actor, found, err
}

func (t *ActorTree) GetActorV5(addr address.Address) (*ActorV5, bool, error) {
	if addr.Protocol() != address.ID {
		return nil, false, xerrors.Errorf("non-ID address %v invalid as actor key", addr)
	}
	var actor ActorV5
	found, err := t.Map.Get(abi.AddrKey(addr), &actor)
	return &actor, found, err
}

// Sets the state associated with an address, overwriting if it already present.
func (t *ActorTree) SetActorV4(addr address.Address, actor *ActorV4) error {
	if addr.Protocol() != address.ID {
		return xerrors.Errorf("non-ID address %v invalid as actor key", addr)
	}
	return t.Map.Put(abi.AddrKey(addr), actor)
}

func (t *ActorTree) SetActorV5(addr address.Address, actor *ActorV5) error {
	if addr.Protocol() != address.ID {
		return xerrors.Errorf("non-ID address %v invalid as actor key", addr)
	}
	return t.Map.Put(abi.AddrKey(addr), actor)
}

// Traverses all entries in the tree.
func (t *ActorTree) ForEachV4(fn func(addr address.Address, actor *ActorV4) error) error {
	var val ActorV4
	return t.Map.ForEach(&val, func(key string) error {
		addr, err := address.NewFromBytes([]byte(key))
		if err != nil {
			return err
		}
		return fn(addr, &val)
	})
}

func (t *ActorTree) ForEachV5(fn func(addr address.Address, actor *ActorV5) error) error {
	var val ActorV5
	return t.Map.ForEach(&val, func(key string) error {
		addr, err := address.NewFromBytes([]byte(key))
		if err != nil {
			return err
		}
		return fn(addr, &val)
	})
}

// Traverses all keys in the tree, without decoding the actor states.
func (t *ActorTree) ForEachKey(fn func(addr address.Address) error) error {
	return t.Map.ForEach(nil, func(key string) error {
		addr, err := address.NewFromBytes([]byte(key))
		if err != nil {
			return err
		}
		return fn(addr)
	})
}
