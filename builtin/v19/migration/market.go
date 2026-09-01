package migration

import (
	"context"

	market18 "github.com/filecoin-project/go-state-types/builtin/v18/market"
	market19 "github.com/filecoin-project/go-state-types/builtin/v19/market"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"
)

// marketMigrator drops PendingDealAllocationIds, which FIP-0118 removes along with the
// verified registry's role in deal activation.
type marketMigrator struct {
	OutCodeCID cid.Cid
}

func (m marketMigrator) MigratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m marketMigrator) Deferred() bool {
	return false
}

func (m marketMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (*migration.ActorMigrationResult, error) {
	var inState market18.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load market state for %s: %w", in.Address, err)
	}

	// FIP-0118 ends the verified registry's role in deal activation, so nothing claims
	// an allocation and nothing reads these ids again so `PublishStorageDeals` becomes useless.
	outState := market19.State{
		Proposals:                     inState.Proposals,
		States:                        inState.States,
		PendingProposals:              inState.PendingProposals,
		EscrowTable:                   inState.EscrowTable,
		LockedTable:                   inState.LockedTable,
		NextID:                        inState.NextID,
		DealOpsByEpoch:                inState.DealOpsByEpoch,
		LastCron:                      inState.LastCron,
		TotalClientLockedCollateral:   inState.TotalClientLockedCollateral,
		TotalProviderLockedCollateral: inState.TotalProviderLockedCollateral,
		TotalClientStorageFee:         inState.TotalClientStorageFee,
		ProviderSectors:               inState.ProviderSectors,
	}

	newHead, err := store.Put(ctx, &outState)
	if err != nil {
		return nil, xerrors.Errorf("failed to put new market state: %w", err)
	}
	return &migration.ActorMigrationResult{NewCodeCID: m.OutCodeCID, NewHead: newHead}, nil
}
