package migration

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	market12 "github.com/filecoin-project/go-state-types/builtin/v12/market"
	market13 "github.com/filecoin-project/go-state-types/builtin/v13/market"
	miner13 "github.com/filecoin-project/go-state-types/builtin/v13/miner"
	"github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"
)

type marketMigrator struct {
	providerSectors *providerSectors
	OutCodeCID      cid.Cid
}

func newMarketMigrator(ctx context.Context, store cbor.IpldStore, outCode cid.Cid, ps *providerSectors) (*marketMigrator, error) {
	return &marketMigrator{
		providerSectors: ps,
		OutCodeCID:      outCode,
	}, nil
}

func (m *marketMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (result *migration.ActorMigrationResult, err error) {
	var inState market12.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load market state for %s: %w", in.Address, err)
	}

	providerSectors, newStates, err := m.migrateProviderSectorsAndStates(ctx, store, inState.States)
	if err != nil {
		return nil, xerrors.Errorf("failed to migrate provider sectors: %w", err)
	}

	outState := market13.State{
		Proposals:                     inState.Proposals,
		States:                        newStates,
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
		ProviderSectors:               providerSectors,
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

func (m *marketMigrator) MigratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m *marketMigrator) migrateProviderSectorsAndStates(ctx context.Context, store cbor.IpldStore, states cid.Cid) (cid.Cid, cid.Cid, error) {

	// out HAMT[Address]HAMT[SectorNumber]SectorDealIDs
	ctxStore := adt.WrapStore(ctx, store)

	oldStateArray, err := adt.AsArray(ctxStore, states, market12.StatesAmtBitwidth)
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to load states array: %w", err)
	}

	newStateArray, err := adt.AsArray(ctxStore, states, market13.StatesAmtBitwidth)
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to load states array: %w", err)
	}

	providerSectorsMem := map[abi.ActorID]map[abi.SectorNumber][]abi.DealID{}

	var oldState market12.DealState
	var newState market13.DealState
	err = oldStateArray.ForEach(&oldState, func(i int64) error {
		deal := abi.DealID(i)

		newState.SlashEpoch = oldState.SlashEpoch
		newState.LastUpdatedEpoch = oldState.LastUpdatedEpoch
		newState.SectorStartEpoch = oldState.SectorStartEpoch
		newState.SectorNumber = 0 // terminated / not found (?)

		if oldState.SlashEpoch == -1 {
			sid, ok := m.providerSectors.dealToSector[deal]
			if ok {
				newState.SectorNumber = sid.Number
			}

			if _, ok := providerSectorsMem[sid.Miner]; !ok {
				providerSectorsMem[sid.Miner] = make(map[abi.SectorNumber][]abi.DealID)
			}
			providerSectorsMem[sid.Miner][sid.Number] = append(providerSectorsMem[sid.Miner][sid.Number], deal)
		}

		if err := newStateArray.AppendContinuous(&newState); err != nil {
			return xerrors.Errorf("failed to append new state: %w", err)
		}

		return nil
	})
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to iterate states: %w", err)
	}

	newStateArrayRoot, err := newStateArray.Root()
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to get newStateArrayRoot: %w", err)
	}

	outProviderSectors, err := adt.MakeEmptyMap(ctxStore, market13.ProviderSectorsHamtBitwidth)
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to create empty map: %w", err)
	}

	for miner, sectors := range providerSectorsMem {
		actorSectors, err := adt.MakeEmptyMap(ctxStore, market13.ProviderSectorsHamtBitwidth)
		if err != nil {
			return cid.Undef, cid.Undef, xerrors.Errorf("failed to create empty map: %w", err)
		}

		var sectorDealIDs market13.SectorDealIDs
		for sector, deals := range sectors {
			sectorDealIDs.Deals = deals

			if err := actorSectors.Put(miner13.SectorKey(sector), &sectorDealIDs); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to put sector: %w", err)
			}
		}

		maddr, err := address.NewIDAddress(uint64(miner))
		if err != nil {
			return cid.Undef, cid.Undef, xerrors.Errorf("failed to convert miner ID to address: %w", err)
		}

		if err := outProviderSectors.Put(abi.AddrKey(maddr), actorSectors); err != nil {
			return cid.Undef, cid.Undef, xerrors.Errorf("failed to put actor sectors: %w", err)
		}
	}

	providerSectorsRoot, err := outProviderSectors.Root()
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to get providerSectorsRoot: %w", err)
	}

	return providerSectorsRoot, newStateArrayRoot, nil
}

var _ migration.ActorMigration = (*marketMigrator)(nil)
