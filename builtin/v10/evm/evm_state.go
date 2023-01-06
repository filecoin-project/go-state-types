package evm

import (
	"github.com/ipfs/go-cid"
	xerrors "golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
)

type State struct {
	Bytecode      cid.Cid
	ContractState cid.Cid
	Nonce         uint64
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
