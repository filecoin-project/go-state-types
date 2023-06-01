package migration

import (
	"context"

	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/builtin"
)

type migrationJob struct {
	address.Address
	// This assumes ActorV5, which is relatively stable.
	builtin.ActorV5
	ActorMigration
	cache MigrationCache
}

type migrationJobResult struct {
	address.Address
	builtin.ActorV5
}

func (job *migrationJob) run(ctx context.Context, store cbor.IpldStore) (*migrationJobResult, error) {
	result, err := job.MigrateState(ctx, store, ActorMigrationInput{
		Address: job.Address,
		Head:    job.ActorV5.Head,
		Cache:   job.cache,
	})
	if err != nil {
		return nil, xerrors.Errorf("state migration failed for actor code %s, addr %s: %w",
			job.ActorV5.Code, job.Address, err)
	}

	// Set up new actor record with the migrated state.
	return &migrationJobResult{
		job.Address, // Unchanged
		builtin.ActorV5{
			Code:       result.NewCodeCID,
			Head:       result.NewHead,
			CallSeqNum: job.ActorV5.CallSeqNum, // Unchanged
			Balance:    job.ActorV5.Balance,    // Unchanged
			Address:    job.ActorV5.Address,    // Unchanged
		},
	}, nil
}
