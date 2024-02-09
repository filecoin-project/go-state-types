package migration

import (
	"bytes"
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

	var sector miner.SectorOnChainInfo

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

	if !hasCached {
		// no cached migration, so we simply iterate all sectors and collect deal IDs

		err = inSectors.ForEach(&sector, func(i int64) error {
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
			m.providerSectors.lk.Unlock()

			return nil
		})
		if err != nil {
			return nil, xerrors.Errorf("failed to iterate sectors: %w", err)
		}
	} else {
		diffs, err := amt.Diff(ctx, store, store, prevSectors, inState.Sectors, amt.UseTreeBitWidth(miner12.SectorsAmtBitwidth))
		if err != nil {
			return nil, xerrors.Errorf("failed to diff old and new Sector AMTs: %w", err)
		}

		for i, change := range diffs {
			sectorNo := abi.SectorNumber(change.Key)

			switch change.Type {
			case amt.Add:

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

				//fmt.Printf("prov dealsector ADD %d: %v\n", sectorNo, sector.DealIDs)

				m.providerSectors.lk.Lock()
				for _, dealID := range sector.DealIDs {
					m.providerSectors.dealToSector[dealID] = abi.SectorID{
						Miner:  abi.ActorID(mid),
						Number: sectorNo,
					}
				}
				m.providerSectors.lk.Unlock()
			case amt.Modify:
				// oh snap deals??

				var sectorBefore miner.SectorOnChainInfo
				if err := sectorBefore.UnmarshalCBOR(bytes.NewReader(change.Before.Raw)); err != nil {
					return nil, xerrors.Errorf("failed to unmarshal sector %d before: %w", sectorNo, err)
				}

				var sectorAfter miner.SectorOnChainInfo

				if err := sectorAfter.UnmarshalCBOR(bytes.NewReader(change.After.Raw)); err != nil {
					return nil, xerrors.Errorf("failed to unmarshal sector %d after: %w", sectorNo, err)
				}

				if len(sectorBefore.DealIDs) != len(sectorAfter.DealIDs) {
					if len(sectorBefore.DealIDs) != 0 {
						return nil, xerrors.Errorf("WHAT?! sector %d modified, this not supported and not supposed to happen", i) // todo: is it? Can't happen w/o a deep, deep reorg, and even then we wouldn't use the cache??
					}
					// snap

					//fmt.Printf("prov dealsector MOD %d: %v -> %v\n", sectorNo, sectorBefore.DealIDs, sectorAfter.DealIDs)

					m.providerSectors.lk.Lock()
					for _, dealID := range sectorAfter.DealIDs {
						m.providerSectors.dealToSector[dealID] = abi.SectorID{
							Miner:  abi.ActorID(mid),
							Number: sectorNo,
						}
					}
					m.providerSectors.lk.Unlock()
				}

				// extensions, etc. here; we don't care about those
			case amt.Remove:
				// nothing to do here, market removes deals based on activation/slash status, and can tell what
				// mappings to remove because non-slashed deals already had them

				//fmt.Printf("prov dealsector REM %d: %v\n", sectorNo, sector.DealIDs)
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
