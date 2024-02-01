package migration

import (
	"context"

	init10 "github.com/filecoin-project/go-state-types/builtin/v10/init"

	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	cbg "github.com/whyrusleeping/cbor-gen"

	"golang.org/x/xerrors"

	init9 "github.com/filecoin-project/go-state-types/builtin/v9/init"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
)

// Init Actor migrator
type initActorMigrator struct {
	OutCodeCID     cid.Cid
	EthZeroAddress addr.Address
}

func (m initActorMigrator) MigratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m initActorMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (*migration.ActorMigrationResult, error) {
	var inState init9.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to get init actor state: %w", err)
	}

	outState := init10.State{
		NetworkName: inState.NetworkName,
	}

	ctxStore := adt9.WrapStore(ctx, store)

	inAddrMap, err := adt9.AsMap(ctxStore, inState.AddressMap, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to load address map: %w", err)
	}

	actorID := cbg.CborInt(inState.NextID)
	outState.NextID = inState.NextID + 1
	if err := inAddrMap.Put(abi.AddrKey(m.EthZeroAddress), &actorID); err != nil {
		return nil, xerrors.Errorf("failed to put new delegated addr: %w", err)
	}

	outAddrMapHead, err := inAddrMap.Root()
	if err != nil {
		return nil, xerrors.Errorf("failed to flush init addr map: %w", err)
	}
	outState.AddressMap = outAddrMapHead

	outHead, err := store.Put(ctx, &outState)
	if err != nil {
		return nil, xerrors.Errorf("failed to put new init state: %w", err)
	}

	return &migration.ActorMigrationResult{
		NewCodeCID: m.OutCodeCID,
		NewHead:    outHead,
	}, nil
}

func (m initActorMigrator) Deferred() bool {
	return false
}
