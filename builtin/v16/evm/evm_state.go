package evm

import (
	"github.com/ipfs/go-cid"
	xerrors "golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v16/util/adt"
)

// TransientDataLifespan is a structure representing the transient data lifespan.
type TransientDataLifespan struct {
	Origin abi.ActorID // The origin actor ID associated with the transient data.
	Nonce  uint64 // A unique nonce identifying the transaction.
}

// TransientData is a structure representing Transient Data
type TransientData struct {
	TransientDataState    cid.Cid // The contract transient data state dictionary. Transient Data State is a map of U256 -> U256 values. KAMT<U256, U256>
	TransientDataLifespan TransientDataLifespan // The data representing the transient data lifespan
}

type Tombstone struct {
	Origin abi.ActorID
	Nonce  uint64
}

type State struct {
	Bytecode      cid.Cid
	BytecodeHash  [32]byte
	ContractState cid.Cid
	TransientData *TransientData
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
