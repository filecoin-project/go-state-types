package evm

import (
	"github.com/ipfs/go-cid"
	xerrors "golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v16/util/adt"
)

// TransientDataLifespan represents the lifespan of transient data.
// It includes the origin actor ID and a nonce to uniquely identify the transaction.
type TransientDataLifespan struct {
	Origin abi.ActorID // The origin actor ID associated with the transient data.
	Nonce  uint64      // A unique nonce identifying the transaction.
}

// TransientData holds the transient state of an EVM contract.
// It includes a CID representing the transient data state and its lifespan.
type TransientData struct {
	TransientDataState    cid.Cid               // The contract transient data state CID.
	TransientDataLifespan TransientDataLifespan // The lifespan of the transient data.
}

// Tombstone represents a self-destructed contract.
// It stores the origin actor ID and nonce from the message that triggered the self-destruction.
type Tombstone struct {
	Origin abi.ActorID // The message origin when this actor was self-destructed.
	Nonce  uint64      // The message nonce when this actor was self-destructed.
}

// State defines the storage structure for an EVM contract.
// It contains contract bytecode, contract state, nonce, transient data, and a possible tombstone.
type State struct {
	Bytecode      cid.Cid        // The EVM contract bytecode CID.
	BytecodeHash  [32]byte       // The Keccak256 hash of the contract bytecode.
	ContractState cid.Cid        // The CID representing the contract's persistent state.
	TransientData *TransientData // The transient state of the contract, if any.
	Nonce         uint64         // The contract nonce used for CREATE/CREATE2.
	Tombstone     *Tombstone     // A tombstone indicating self-destruction, if applicable.
}

// ConstructState initializes a new contract state with the provided bytecode.
// It creates an empty state map and returns the constructed state.
//
// Parameters:
// - store: The ADT store to manage state storage.
// - bytecode: The CID of the contract's bytecode.
//
// Returns:
// - A pointer to the newly constructed State instance.
// - An error if the empty map creation fails.
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
