package migration

import (
	"context"

	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	xerrors "golang.org/x/xerrors"

	evm16 "github.com/filecoin-project/go-state-types/builtin/v16/evm"
	evm17 "github.com/filecoin-project/go-state-types/builtin/v17/evm"
)

// evmMigrator performs the migration for the EVM contract state,
// adding an empty TransientData field.
type evmMigrator struct {
	OutCodeCID cid.Cid
}

func (m evmMigrator) MigratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m evmMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (*migration.ActorMigrationResult, error) {
	// Load the existing state (v16).
	var inState evm16.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load evm state for %s: %w", in.Address, err)
	}

	// Create the new state (v17) with TransientData set to nil.
	outState := evm17.State{
		Bytecode:      inState.Bytecode,
		BytecodeHash:  inState.BytecodeHash,
		ContractState: inState.ContractState,
		Nonce:         inState.Nonce,
		TransientData: nil, // Add empty TransientData as nil
		Tombstone:     inState.Tombstone,
	}

	// Store the new state.
	newHead, err := store.Put(ctx, &outState)
	if err != nil {
		return nil, xerrors.Errorf("failed to put new evm state: %w", err)
	}

	return &migration.ActorMigrationResult{
		NewCodeCID: m.MigratedCodeCID(),
		NewHead:    newHead,
	}, nil
}

func (m evmMigrator) Deferred() bool {
	return false
}
