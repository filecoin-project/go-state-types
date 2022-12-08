package migration

import (
	"context"

	"github.com/filecoin-project/go-state-types/migration"

	"github.com/filecoin-project/go-state-types/builtin/v9/system"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
)

// System Actor migrator
type systemActorMigrator struct {
	OutCodeCID   cid.Cid
	ManifestData cid.Cid
}

func (m systemActorMigrator) MigratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m systemActorMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (*migration.ActorMigrationResult, error) {
	// The ManifestData itself is already in the blockstore
	state := system.State{BuiltinActors: m.ManifestData}
	stateHead, err := store.Put(ctx, &state)
	if err != nil {
		return nil, err
	}

	return &migration.ActorMigrationResult{
		NewCodeCID: m.OutCodeCID,
		NewHead:    stateHead,
	}, nil
}
