package migration

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	miner12 "github.com/filecoin-project/go-state-types/builtin/v12/miner"
	"github.com/filecoin-project/go-state-types/builtin/v13/miner"
	"github.com/filecoin-project/go-state-types/builtin/v13/util/adt"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"
	"sync"
)

type providerSectors struct {
	lk sync.Mutex

	dealToSector map[abi.DealID]abi.SectorID
}

// minerMigration is technically a no-op, but it collects a cache for market migration
type minerMigrator struct {
	providerSectors *providerSectors

	OutCodeCID cid.Cid
}

func newMinerMigrator(ctx context.Context, store cbor.IpldStore, outCode cid.Cid, ps *providerSectors) (*minerMigrator, error) {
	return &minerMigrator{
		providerSectors: ps,

		OutCodeCID: outCode,
	}, nil
}

func (m *minerMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (result *migration.ActorMigrationResult, err error) {
	var inState miner12.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load miner state for %s: %w", in.Address, err)
	}

	ctxStore := adt.WrapStore(ctx, store)

	sa, err := adt.AsArray(ctxStore, inState.Sectors, miner.SectorsAmtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to load sectors array: %w", err)
	}

	var sector miner.SectorOnChainInfo

	mid, err := address.IDFromAddress(in.Address)
	if err != nil {
		return nil, xerrors.Errorf("failed to get miner ID: %w", err)
	}

	err = sa.ForEach(&sector, func(i int64) error {
		if len(sector.DealIDs) == 0 {
			return nil
		}

		m.providerSectors.lk.Lock()
		for _, dealID := range sector.DealIDs {
			m.providerSectors.dealToSector[dealID] = abi.SectorID{
				Miner:  abi.ActorID(mid),
				Number: abi.SectorNumber(i),
			}
		}

		dealIDsCopy := make([]abi.DealID, len(sector.DealIDs))
		copy(dealIDsCopy, sector.DealIDs)

		m.providerSectors.lk.Unlock()

		return nil
	})
	if err != nil {
		return nil, xerrors.Errorf("failed to iterate sectors: %w", err)
	}

	return &migration.ActorMigrationResult{
		NewCodeCID: m.MigratedCodeCID(),
		NewHead:    in.Head,
	}, nil
}

func (m *minerMigrator) MigratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

var _ migration.ActorMigration = (*minerMigrator)(nil)
