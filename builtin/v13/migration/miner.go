package migration

import (
	"context"
	"sync"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-amt-ipld/v4"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v13/miner"
	"github.com/filecoin-project/go-state-types/builtin/v13/util/adt"

	miner12 "github.com/filecoin-project/go-state-types/builtin/v12/miner"
	"github.com/filecoin-project/go-state-types/migration"
)

type providerSectors struct {
	lk sync.Mutex

	// Populated ONLY on the first run-through, when there's no cache to diff against
	dealToSector         map[abi.DealID]abi.SectorID
	minerToSectorToDeals map[abi.ActorID]map[abi.SectorNumber][]abi.DealID

	// Populated ONLY on future run-throughs if there is a cache to diff against
	newDealsToSector              map[abi.DealID]abi.SectorID
	updatesToMinerToSectorToDeals map[abi.ActorID]map[abi.SectorNumber][]abi.DealID
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

	mid, err := address.IDFromAddress(in.Address)
	if err != nil {
		return nil, xerrors.Errorf("failed to get miner ID: %w", err)
	}

	inSectors, err := adt.AsArray(ctxStore, inState.Sectors, miner12.SectorsAmtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to load sectors array: %w", err)
	}

	hasCached, prevSectors, err := in.Cache.Read(migration.MinerPrevSectorsInKey(in.Address))
	if err != nil {
		return nil, xerrors.Errorf("failed to read prev sectors from cache: %w", err)
	}

	var sector miner.SectorOnChainInfo
	if !hasCached {
		// no cached migration, so we simply iterate all sectors and collect deal IDs

		err = inSectors.ForEach(&sector, func(i int64) error {
			if len(sector.DealIDs) == 0 || sector.Expiration < UpgradeHeight {
				return nil
			}

			m.providerSectors.lk.Lock()
			for _, dealID := range sector.DealIDs {
				m.providerSectors.dealToSector[dealID] = abi.SectorID{
					Miner:  abi.ActorID(mid),
					Number: abi.SectorNumber(i),
				}
			}

			sectorDeals, ok := m.providerSectors.minerToSectorToDeals[abi.ActorID(mid)]
			if !ok {
				sectorDeals = make(map[abi.SectorNumber][]abi.DealID)
				m.providerSectors.minerToSectorToDeals[abi.ActorID(mid)] = sectorDeals
			}
			// There can't already be an entry for this sector, so just set (not append)
			// TODO: golang: I need to copy this, since sector is a reference that's constantly updated? Right? Steb?
			sectorDealIDs := sector.DealIDs
			sectorDeals[abi.SectorNumber(i)] = sectorDealIDs

			m.providerSectors.lk.Unlock()

			return nil
		})
		if err != nil {
			return nil, xerrors.Errorf("failed to iterate sectors: %w", err)
		}
	} else {
		prevInSectors, err := adt.AsArray(ctxStore, prevSectors, miner12.SectorsAmtBitwidth)
		if err != nil {
			return nil, xerrors.Errorf("failed to load previous input sectors array: %w", err)
		}

		diffs, err := amt.Diff(ctx, store, store, prevSectors, inState.Sectors, amt.UseTreeBitWidth(miner12.SectorsAmtBitwidth))
		if err != nil {
			return nil, xerrors.Errorf("failed to diff old and new Sector AMTs: %w", err)
		}

		for _, change := range diffs {
			sectorNo := abi.SectorNumber(change.Key)

			switch change.Type {
			case amt.Add:
				fallthrough
			case amt.Modify:
				found, err := inSectors.Get(change.Key, &sector)
				if err != nil {
					return nil, xerrors.Errorf("failed to get sector %d in inSectors: %w", sectorNo, err)
				}

				if !found {
					return nil, xerrors.Errorf("didn't find sector %d in inSectors", sectorNo)
				}

				if len(sector.DealIDs) == 0 {
					// if no deals don't even bother taking the lock
					continue
				}

				m.providerSectors.lk.Lock()
				for _, dealID := range sector.DealIDs {
					m.providerSectors.newDealsToSector[dealID] = abi.SectorID{
						Miner:  abi.ActorID(mid),
						Number: sectorNo,
					}
				}

				sectorDeals, ok := m.providerSectors.updatesToMinerToSectorToDeals[abi.ActorID(mid)]
				if !ok {
					sectorDeals = make(map[abi.SectorNumber][]abi.DealID)
					m.providerSectors.updatesToMinerToSectorToDeals[abi.ActorID(mid)] = sectorDeals
				}
				// There can't already be an entry for this sector, so just set (not append)
				// TODO: golang: I need to copy this, since sector is a reference that's constantly updated? Right? Steb?
				sectorDealIDs := sector.DealIDs
				sectorDeals[sectorNo] = sectorDealIDs

				m.providerSectors.lk.Unlock()
			case amt.Remove:
				// In this unlikely case, we need to load the deleted SectorOnChainInfo, and mark it as not having any deals
				var oldSector miner.SectorOnChainInfo
				found, err := prevInSectors.Get(change.Key, &oldSector)
				if err != nil {
					return nil, xerrors.Errorf("failed to get previous sector info: %w", err)
				}
				if !found {
					return nil, xerrors.Errorf("didn't find previous sector info for %d", sectorNo)
				}

				if len(oldSector.DealIDs) != 0 {
					sectorDeals, ok := m.providerSectors.updatesToMinerToSectorToDeals[abi.ActorID(mid)]
					if !ok {
						sectorDeals = make(map[abi.SectorNumber][]abi.DealID)
						m.providerSectors.updatesToMinerToSectorToDeals[abi.ActorID(mid)] = sectorDeals
					}
					// TODO: Is it better to use nil to communicate deals to be removed, or just a separate map (of maps)?
					sectorDeals[sectorNo] = nil
				}
			}
		}
	}

	err = in.Cache.Write(migration.MinerPrevSectorsInKey(in.Address), inState.Sectors)
	if err != nil {
		return nil, xerrors.Errorf("failed to write prev sectors to cache: %w", err)
	}

	return &migration.ActorMigrationResult{
		NewCodeCID: m.MigratedCodeCID(),
		NewHead:    in.Head,
	}, nil
}

func (m *minerMigrator) MigratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m *minerMigrator) Deferred() bool {
	return false
}

var _ migration.ActorMigration = (*minerMigrator)(nil)
