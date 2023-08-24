package migration

import (
	"context"

	market11 "github.com/filecoin-project/go-state-types/builtin/v11/market"
	market12 "github.com/filecoin-project/go-state-types/builtin/v12/market"

	"github.com/filecoin-project/go-state-types/migration"

	"golang.org/x/xerrors"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
)

// The marketMigrator performs the following migrations:
// TOOD

func newMarketMigrator(ctx context.Context, store cbor.IpldStore, outCode cid.Cid, x *minerMigrator) (*marketMigrator, error) {
	sectorToDealIdHamtCid, err := x.sectorDeals.Map.Root()
	if err != nil {
		return nil, err
	}

	return &marketMigrator{
		OutCodeCID:    outCode,
		sectorDealIDs: sectorToDealIdHamtCid,
	}, nil
}

type marketMigrator struct {
	sectorDealIDs cid.Cid

	OutCodeCID cid.Cid
}

func (m marketMigrator) MigratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m marketMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (*migration.ActorMigrationResult, error) {
	var inState market11.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load miner state for %s: %w", in.Address, err)
	}

	newSectorDeals := m.sectorDealIDs

	// Create the new state
	outState := market12.State{
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
		PendingDealAllocationIds:      inState.PendingDealAllocationIds,
		SectorDeals:                   newSectorDeals, // Updated value
	}

	newHead, err := store.Put(ctx, &outState)
	if err != nil {
		return nil, xerrors.Errorf("failed to put new state: %w", err)
	}

	return &migration.ActorMigrationResult{
		NewCodeCID: m.MigratedCodeCID(),
		NewHead:    newHead,
	}, nil
}
