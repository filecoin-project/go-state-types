package evm

import (
	"io"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	xerrors "golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
)

type Tombstone struct {
	Origin abi.ActorID
	Nonce  uint64
}

type State struct {
	Bytecode      cid.Cid
	BytecodeHash  CodeHash
	ContractState cid.Cid
	Nonce         uint64
	Tombstone     *Tombstone
}

func ConstructState(store adt.Store, bytecode cid.Cid) (*State, error) {
	emptyMapCid, err := adt.StoreEmptyMap(store, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to create empty map: %w", err)
	}

	return &State{
		Bytecode:      bytecode,
		ContractState: emptyMapCid,
		Nonce:         0,
	}, nil
}

// NOTE This struct does not create a valid cbor-encoded byte array. It just passes the bytes through as-is.
type CodeHash [32]byte

// MarshalCBOR Does NOT marshall to a cbor-encoding. This is just syntactic sugar to let us pass bytes transparently through lotus which requires a cbor-marshallable object.
func (t *CodeHash) MarshalCBOR(w io.Writer) error {
	_, err := w.Write((*t)[:])
	return err
}

// UnmarshalCBOR CANNOT read a cbor-encoded byte slice. This will just transparently pass the underlying bytes.
func (t *CodeHash) UnmarshalCBOR(r io.Reader) error {
	read, err := r.Read((*t)[:])
	if err != nil || read != 32 {
		return xerrors.Errorf("failed to read codehash: %w", err)
	}

	return nil
}
