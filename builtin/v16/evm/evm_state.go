package evm

import (
	"github.com/ipfs/go-cid"
	xerrors "golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v16/util/adt"
)

// TransientDataLifespan represents the lifespan of transient data.
// It includes an origin ActorID and a unique nonce to track the transaction.
type TransientDataLifespan struct {
	Origin abi.ActorID // The origin actor ID associated with the transient data.
	Nonce  uint64      // A unique nonce identifying the transaction.
}

// TransientData represents transient storage data in an EVM contract.
// It consists of a state dictionary (KAMT<U256, U256>) and a lifespan tracker.
type TransientData struct {
	TransientDataState    cid.Cid               // CID of the transient data state dictionary.
	TransientDataLifespan TransientDataLifespan // Data representing the transient data lifespan.
}

// Tombstone represents a marker for contracts that have been self-destructed.
type Tombstone struct {
	Origin abi.ActorID // The ActorID that initiated the self-destruct.
	Nonce  uint64      // The transaction nonce at the time of self-destruction.
}

// State represents the on-chain state of an EVM contract.
// It includes the contract bytecode, its hash, storage state, transient data, and nonce tracking.
type State struct {
	Bytecode      cid.Cid        // CID of the EVM contract bytecode.
	BytecodeHash  [32]byte       // Keccak256 hash of the contract bytecode.
	ContractState cid.Cid        // CID of the contract's persistent state dictionary.
	TransientData *TransientData // Optional transient storage data associated with the contract.
	Nonce         uint64         // Tracks how many times CREATE or CREATE2 have been invoked.
	Tombstone     *Tombstone     // Optional marker indicating if the contract has been self-destructed.
}

// ConstructState initializes and returns a new contract state with an empty state dictionary.
//
// Parameters:
//   - store: The ADT store used for state management.
//   - bytecode: CID representing the contract bytecode.
//
// Returns:
//   - *State: A pointer to the newly constructed State object.
//   - error: An error if state initialization fails.
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
