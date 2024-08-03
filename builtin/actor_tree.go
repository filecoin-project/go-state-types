package builtin

import (
	"crypto/sha256"

	"github.com/filecoin-project/go-address"
	hamt "github.com/filecoin-project/go-hamt-ipld/v4"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v14/util/adt"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

type LegacyActorTree interface {
	GetStore() adt.Store
	GetMap() *adt.Map
	Flush() (cid.Cid, error)
	GetActorV4(addr address.Address) (*ActorV4, bool, error)
	GetActorV5(addr address.Address) (*ActorV5, bool, error)
	SetActorV4(addr address.Address, actor *ActorV4) error
	SetActorV5(addr address.Address, actor *ActorV5) error
	ForEachV4(fn func(addr address.Address, actor *ActorV4) error) error
	ForEachV5(fn func(addr address.Address, actor *ActorV5) error) error
	ForEachKey(fn func(addr address.Address) error) error
}

type ActorTree interface {
	GetStore() adt.Store
	Flush() (cid.Cid, error)
	GetActorV5(addr address.Address) (*ActorV5, bool, error)
	SetActorV5(addr address.Address, actor *ActorV5) error
	ForEachV5(fn func(addr address.Address, actor *ActorV5) error) error
	ForEachKey(fn func(addr address.Address) error) error
}

var _ LegacyActorTree = (*legacyActorTree)(nil)
var _ ActorTree = (*legacyActorTree)(nil)

// var _ ActorTree = (*actorTree)(nil)

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
	Code             cid.Cid          // CID representing the code associated with the actor
	Head             cid.Cid          // CID of the head state object for the actor
	CallSeqNum       uint64           // CallSeqNum for the next message to be received by the actor (non-zero for accounts only)
	Balance          big.Int          // Token balance of the actor
	DelegatedAddress *address.Address // Delegated (f4) actor address
}

func (a *ActorV5) Equals(o *ActorV5) bool {
	if a == nil && o == nil {
		return true
	}
	if a == nil || o == nil {
		return false
	}
	return a.Code == o.Code &&
		a.Head == o.Head &&
		a.CallSeqNum == o.CallSeqNum &&
		a.Balance.Equals(o.Balance) &&
		((a.DelegatedAddress == nil && o.DelegatedAddress == nil) || (a.DelegatedAddress != nil && o.DelegatedAddress != nil && *a.DelegatedAddress == *o.DelegatedAddress))
}

func (a *ActorV5) New() *ActorV5 {
	return new(ActorV5)
}

// A specialization of a map of ID-addresses to actor heads.
type actorTree struct {
	Map      *hamt.Node[*ActorV5]
	LastRoot cid.Cid
	Store    adt.Store
}

// Initializes a new, empty state tree backed by a store.
func NewTree(store adt.Store) (ActorTree, error) {
	nd, err := hamt.NewNode[*ActorV5](
		store,
		hamt.UseHashFunction(func(input []byte) []byte {
			res := sha256.Sum256(input)
			return res[:]
		}),
		hamt.UseTreeBitWidth(DefaultHamtBitwidth),
	)
	if err != nil {
		return nil, err
	}
	return &actorTree{
		Map:      nd,
		LastRoot: cid.Undef,
		Store:    store,
	}, nil
}

// Loads a tree from a root CID and store.
func LoadTree(store adt.Store, root cid.Cid) (ActorTree, error) {
	nd, err := hamt.LoadNode[*ActorV5](
		store.Context(),
		store,
		root,
		hamt.UseHashFunction(func(input []byte) []byte {
			res := sha256.Sum256(input)
			return res[:]
		}),
		hamt.UseTreeBitWidth(DefaultHamtBitwidth),
	)
	if err != nil {
		return nil, err
	}
	return &actorTree{
		Map:      nd,
		LastRoot: root,
		Store:    store,
	}, nil
}

func (t *actorTree) GetStore() adt.Store {
	return t.Store
}

// Writes the tree root node to the store, and returns its CID.
func (t *actorTree) Flush() (cid.Cid, error) {
	if err := t.Map.Flush(t.Store.Context()); err != nil {
		return cid.Undef, xerrors.Errorf("failed to flush map root: %w", err)
	}
	c, err := t.Store.Put(t.Store.Context(), t.Map)
	if err != nil {
		return cid.Undef, xerrors.Errorf("writing map root object: %w", err)
	}
	t.LastRoot = c
	return c, nil
}

func (t *actorTree) GetActorV5(addr address.Address) (*ActorV5, bool, error) {
	if addr.Protocol() != address.ID {
		return nil, false, xerrors.Errorf("non-ID address %v invalid as actor key", addr)
	}
	if addr.Protocol() != address.ID {
		return nil, false, xerrors.Errorf("non-ID address %v invalid as actor key", addr)
	}
	return t.Map.Find(t.Store.Context(), abi.AddrKey(addr).Key())
}

func (t *actorTree) SetActorV5(addr address.Address, actor *ActorV5) error {
	if addr.Protocol() != address.ID {
		return xerrors.Errorf("non-ID address %v invalid as actor key", addr)
	}
	return t.Map.Set(t.Store.Context(), abi.AddrKey(addr).Key(), actor)
}

func (t *actorTree) ForEachV5(fn func(addr address.Address, actor *ActorV5) error) error {
	return t.Map.ForEach(t.Store.Context(), func(key string, val *ActorV5) error {
		addr, err := address.NewFromBytes([]byte(key))
		if err != nil {
			return err
		}
		return fn(addr, val)
	})
}

// Traverses all keys in the tree, without decoding the actor states.
func (t *actorTree) ForEachKey(fn func(addr address.Address) error) error {
	return t.Map.ForEach(t.Store.Context(), func(key string, _ *ActorV5) error {
		addr, err := address.NewFromBytes([]byte(key))
		if err != nil {
			return err
		}
		return fn(addr)
	})
}

// A specialization of a map of ID-addresses to actor heads.
type legacyActorTree struct {
	Map   *adt.Map
	Store adt.Store
}

// Initializes a new, empty state tree backed by a store.
func NewLegacyTree(store adt.Store) (LegacyActorTree, error) {
	emptyMap, err := adt.MakeEmptyMap(store, DefaultHamtBitwidth)
	if err != nil {
		return nil, err
	}
	return &legacyActorTree{
		Map:   emptyMap,
		Store: store,
	}, nil
}

// Loads a tree from a root CID and store.
func LoadLegacyTree(s adt.Store, r cid.Cid) (LegacyActorTree, error) {
	m, err := adt.AsMap(s, r, DefaultHamtBitwidth)
	if err != nil {
		return nil, err
	}
	return &legacyActorTree{
		Map:   m,
		Store: s,
	}, nil
}

func (t *legacyActorTree) GetStore() adt.Store {
	return t.Store
}

func (t *legacyActorTree) GetMap() *adt.Map {
	return t.Map
}

// Writes the tree root node to the store, and returns its CID.
func (t *legacyActorTree) Flush() (cid.Cid, error) {
	return t.Map.Root()
}

// Loads the state associated with an address.
func (t *legacyActorTree) GetActorV4(addr address.Address) (*ActorV4, bool, error) {
	if addr.Protocol() != address.ID {
		return nil, false, xerrors.Errorf("non-ID address %v invalid as actor key", addr)
	}
	var actor ActorV4
	found, err := t.Map.Get(abi.AddrKey(addr), &actor)
	return &actor, found, err
}

func (t *legacyActorTree) GetActorV5(addr address.Address) (*ActorV5, bool, error) {
	if addr.Protocol() != address.ID {
		return nil, false, xerrors.Errorf("non-ID address %v invalid as actor key", addr)
	}
	var actor ActorV5
	found, err := t.Map.Get(abi.AddrKey(addr), &actor)
	return &actor, found, err
}

// Sets the state associated with an address, overwriting if it already present.
func (t *legacyActorTree) SetActorV4(addr address.Address, actor *ActorV4) error {
	if addr.Protocol() != address.ID {
		return xerrors.Errorf("non-ID address %v invalid as actor key", addr)
	}
	return t.Map.Put(abi.AddrKey(addr), actor)
}

func (t *legacyActorTree) SetActorV5(addr address.Address, actor *ActorV5) error {
	if addr.Protocol() != address.ID {
		return xerrors.Errorf("non-ID address %v invalid as actor key", addr)
	}
	return t.Map.Put(abi.AddrKey(addr), actor)
}

// Traverses all entries in the tree.
func (t *legacyActorTree) ForEachV4(fn func(addr address.Address, actor *ActorV4) error) error {
	var val ActorV4
	return t.Map.ForEach(&val, func(key string) error {
		addr, err := address.NewFromBytes([]byte(key))
		if err != nil {
			return err
		}
		return fn(addr, &val)
	})
}

func (t *legacyActorTree) ForEachV5(fn func(addr address.Address, actor *ActorV5) error) error {
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
func (t *legacyActorTree) ForEachKey(fn func(addr address.Address) error) error {
	return t.Map.ForEach(nil, func(key string) error {
		addr, err := address.NewFromBytes([]byte(key))
		if err != nil {
			return err
		}
		return fn(addr)
	})
}
