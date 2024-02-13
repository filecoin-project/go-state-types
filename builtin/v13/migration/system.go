package migration

import (
	"context"

	system13 "github.com/filecoin-project/go-state-types/builtin/v13/system"

	"github.com/filecoin-project/go-state-types/migration"

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
	state := system13.State{BuiltinActors: m.ManifestData}
	stateHead, err := store.Put(ctx, &state)
	if err != nil {
		return nil, err
	}

	return &migration.ActorMigrationResult{
		NewCodeCID: m.OutCodeCID,
		NewHead:    stateHead,
	}, nil
}

func (m systemActorMigrator) Deferred() bool {
	return false
}
