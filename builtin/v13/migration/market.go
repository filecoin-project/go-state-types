package migration

import (
	"bytes"
	"context"

	"github.com/filecoin-project/go-amt-ipld/v4"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	typegen "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/abi"
	market12 "github.com/filecoin-project/go-state-types/builtin/v12/market"
	market13 "github.com/filecoin-project/go-state-types/builtin/v13/market"
	miner13 "github.com/filecoin-project/go-state-types/builtin/v13/miner"
	adt13 "github.com/filecoin-project/go-state-types/builtin/v13/util/adt"
	"github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	"github.com/filecoin-project/go-state-types/migration"
)

type marketMigrator struct {
	providerSectors *providerSectors
	upgradeEpoch    abi.ChainEpoch
	OutCodeCID      cid.Cid
}

func newMarketMigrator(ctx context.Context, store cbor.IpldStore, outCode cid.Cid, ps *providerSectors, upgradeEpoch abi.ChainEpoch) (*marketMigrator, error) {
	return &marketMigrator{
		providerSectors: ps,
		upgradeEpoch:    upgradeEpoch,
		OutCodeCID:      outCode,
	}, nil
}

func (m *marketMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (result *migration.ActorMigrationResult, err error) {
	var inState market12.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load market state for %s: %w", in.Address, err)
	}

	providerSectors, newStates, err := m.migrateProviderSectorsAndStates(ctx, store, in, inState.States, inState.Proposals)
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

func (m *marketMigrator) migrateProviderSectorsAndStates(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput, states, proposals cid.Cid) (cid.Cid, cid.Cid, error) {
	// providerSectorsRoot: HAMT[ActorID]HAMT[SectorNumber]SectorDealIDs

	okIn, prevInStates, err := in.Cache.Read(migration.MarketPrevDealStatesInKey(in.Address))
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to get previous inRoot from cache: %w", err)
	}

	okInPr, prevInProposals, err := in.Cache.Read(migration.MarketPrevDealProposalsInKey(in.Address))
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to get previous inRoot from cache: %w", err)
	}

	okOut, prevOutStates, err := in.Cache.Read(migration.MarketPrevDealStatesOutKey(in.Address))
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to get previous outRoot from cache: %w", err)
	}

	okOutPs, prevOutProviderSectors, err := in.Cache.Read(migration.MarketPrevProviderSectorsOutKey(in.Address))
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to get previous outRoot from cache: %w", err)
	}

	var providerSectorsRoot, newStateArrayRoot cid.Cid

	if okIn && okInPr && okOut && okOutPs {
		providerSectorsRoot, newStateArrayRoot, err = m.migrateProviderSectorsAndStatesWithDiff(ctx, store, prevInStates, prevOutStates, prevOutProviderSectors, states, prevInProposals, proposals)
		if err != nil {
			return cid.Undef, cid.Undef, xerrors.Errorf("failed to migrate provider sectors (diff): %w", err)
		}
	} else {
		providerSectorsRoot, newStateArrayRoot, err = m.migrateProviderSectorsAndStatesFromScratch(ctx, store, in, states, proposals)
		if err != nil {
			return cid.Undef, cid.Undef, xerrors.Errorf("failed to migrate provider sectors (all): %w", err)
		}
	}

	if err := in.Cache.Write(migration.MarketPrevDealStatesInKey(in.Address), states); err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to write previous inRoot to cache: %w", err)
	}

	if err := in.Cache.Write(migration.MarketPrevDealProposalsInKey(in.Address), proposals); err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to write previous inRoot to cache: %w", err)
	}

	if err := in.Cache.Write(migration.MarketPrevDealStatesOutKey(in.Address), newStateArrayRoot); err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to write previous outRoot to cache: %w", err)
	}

	if err := in.Cache.Write(migration.MarketPrevProviderSectorsOutKey(in.Address), providerSectorsRoot); err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to write previous outRoot to cache: %w", err)
	}

	return providerSectorsRoot, newStateArrayRoot, nil
}

func (m *marketMigrator) migrateProviderSectorsAndStatesWithDiff(ctx context.Context, store cbor.IpldStore, prevInStatesCid, prevOutStatesCid, prevOutProviderSectorsCid, inStatesCid, prevInProposals, inProposals cid.Cid) (cid.Cid, cid.Cid, error) {
	statesDiffs, err := amt.Diff(ctx, store, store, prevInStatesCid, inStatesCid, amt.UseTreeBitWidth(market12.StatesAmtBitwidth))
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to diff old and new deal state AMTs: %w", err)
	}

	ctxStore := adt.WrapStore(ctx, store)

	prevOutStates, err := adt.AsArray(ctxStore, prevOutStatesCid, market13.StatesAmtBitwidth)
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to load prevOutStates array: %w", err)
	}

	prevOutProviderSectors, err := adt.AsMap(ctxStore, prevOutProviderSectorsCid, market13.ProviderSectorsHamtBitwidth)
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to load prevOutProviderSectors map: %w", err)
	}

	for _, change := range statesDiffs {
		deal := abi.DealID(change.Key)

		switch change.Type {
		case amt.Add:
			var oldState market12.DealState
			if err := oldState.UnmarshalCBOR(bytes.NewReader(change.After.Raw)); err != nil {
				return cid.Cid{}, cid.Cid{}, xerrors.Errorf("failed to unmarshal old state: %w", err)
			}

			newState := market13.DealState{
				SectorNumber:     0,
				SectorStartEpoch: oldState.SectorStartEpoch,
				LastUpdatedEpoch: oldState.LastUpdatedEpoch,
				SlashEpoch:       oldState.SlashEpoch,
			}

			// TODO: We should technically only set this if the deal hasn't expired,
			// but if it's a new deal there's no way it's already expired, right?
			if newState.SlashEpoch == -1 {
				dealSectorNumber, ok := m.providerSectors.newDealsToSector[deal]
				if !ok {
					return cid.Undef, cid.Undef, xerrors.Errorf("didn't find new sector number for new deal %d", deal)
				}
				newState.SectorNumber = dealSectorNumber.Number
			}

			if err := prevOutStates.Set(uint64(deal), &newState); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to set new state: %w", err)
			}

		case amt.Remove:
			if err := prevOutStates.Delete(uint64(deal)); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to delete new state: %w", err)
			}

		case amt.Modify:

			var oldState, prevOldState market12.DealState
			var prevNewState market13.DealState
			if err := prevOldState.UnmarshalCBOR(bytes.NewReader(change.Before.Raw)); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to unmarshal old state: %w", err)
			}
			if err := oldState.UnmarshalCBOR(bytes.NewReader(change.After.Raw)); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to unmarshal old state: %w", err)
			}
			ok, err := prevOutStates.Get(uint64(deal), &prevNewState)
			if err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to get previous newstate: %w", err)
			}
			if !ok {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to get previous newstate: not found")
			}

			newState := market13.DealState{
				SectorNumber:     prevNewState.SectorNumber,
				SectorStartEpoch: oldState.SectorStartEpoch,
				LastUpdatedEpoch: oldState.LastUpdatedEpoch,
				SlashEpoch:       oldState.SlashEpoch,
			}

			// if the deal is now slashed, UNSET the sector number
			if newState.SlashEpoch != -1 {
				newState.SectorNumber = 0
			}

			if err := prevOutStates.Set(uint64(deal), &newState); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to set new state: %w", err)
			}
		}
	}

	// now, process updatesToMinerToSectorToDeals, making the appropriate updates

	for miner, sectorUpdates := range m.providerSectors.updatesToMinerToSectorToDeals {
		var actorSectorsMapRoot typegen.CborCid
		found, err := prevOutProviderSectors.Get(abi.UIntKey(uint64(miner)), &actorSectorsMapRoot)
		if err != nil {
			return cid.Undef, cid.Undef, xerrors.Errorf("failed to get actor sectors: %w", err)
		}

		var actorSectors *adt13.Map
		if !found {
			// can happen, this miner just didn't have any sectors before
			actorSectors, err = adt13.MakeEmptyMap(ctxStore, market13.ProviderSectorsHamtBitwidth)
			if err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to make empty actor sectors map: %w", err)
			}
		} else {
			actorSectors, err = adt13.AsMap(ctxStore, cid.Cid(actorSectorsMapRoot), market13.ProviderSectorsHamtBitwidth)
			if err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to load actor sectors map: %w", err)
			}
		}

		var sectorDealIDs market13.SectorDealIDs
		for sector, deals := range sectorUpdates {
			// This signals the unlikely case where a sector with deals got GCed after the premigration, delete the entry for it
			if deals == nil {
				err = actorSectors.Delete(miner13.SectorKey(sector))
				if err != nil {
					return cid.Undef, cid.Undef, xerrors.Errorf("failed to delete actorSectors entry: %w", err)
				}

				continue
			} else {
				sectorDealIDs = deals
				if err := actorSectors.Put(miner13.SectorKey(sector), &sectorDealIDs); err != nil {
					return cid.Undef, cid.Undef, xerrors.Errorf("failed to put sector: %w", err)
				}
			}
		}

		isEmpty, err := actorSectors.IsEmpty()
		if err != nil {
			return cid.Undef, cid.Undef, xerrors.Errorf("failed to check if actorSectors is empty: %w", err)
		}
		if isEmpty {
			if err := prevOutProviderSectors.Delete(abi.UIntKey(uint64(miner))); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to delete actor sectors: %w", err)
			}
		} else {
			if err := prevOutProviderSectors.Put(abi.UIntKey(uint64(miner)), actorSectors); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to put actor sectors: %w", err)
			}
		}
	}

	// flush and return
	outProviderSectors, err := prevOutProviderSectors.Root()
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to get providerSectorsRoot: %w", err)
	}

	outStates, err := prevOutStates.Root()
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to get statesRoot: %w", err)
	}

	return outProviderSectors, outStates, nil
}

func (m *marketMigrator) migrateProviderSectorsAndStatesFromScratch(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput, states cid.Cid, proposals cid.Cid) (cid.Cid, cid.Cid, error) {
	ctxStore := adt.WrapStore(ctx, store)

	oldStateArray, err := adt.AsArray(ctxStore, states, market12.StatesAmtBitwidth)
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to load states array: %w", err)
	}

	newStateArray, err := adt13.MakeEmptyArray(ctxStore, market13.StatesAmtBitwidth)
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to load states array: %w", err)
	}

	proposalsArr, err := adt.AsArray(ctxStore, proposals, market12.ProposalsAmtBitwidth)
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to load proposals array: %w", err)
	}

	proposalsSize := proposalsArr.Length()

	var oldState market12.DealState
	var newState market13.DealState

	err = oldStateArray.ForEach(&oldState, func(i int64) error {
		// Refresh proposals array periodically to avoid holding onto all ~10GB of memory
		// throughout whole migration.
		// This has limited impact on caching speedups because the access pattern is sequential.
		if i%(int64(proposalsSize)/1000) == 0 {
			proposalsArr, err = adt.AsArray(ctxStore, proposals, market12.ProposalsAmtBitwidth)
			if err != nil {
				return xerrors.Errorf("failed to load proposals array: %w", err)
			}
		}
		deal := abi.DealID(i)

		newState.SlashEpoch = oldState.SlashEpoch
		newState.LastUpdatedEpoch = oldState.LastUpdatedEpoch
		newState.SectorStartEpoch = oldState.SectorStartEpoch
		newState.SectorNumber = 0

		var proposal market12.DealProposal
		ok, err := proposalsArr.Get(uint64(i), &proposal)
		if err != nil {
			return xerrors.Errorf("failed to get proposal: %w", err)
		}

		if !ok {
			return xerrors.Errorf("failed to find proposal for deal ID %d", i)
		}

		// FIP: For each unexpired deal state object in the market actor state that has a terminated epoch set to -1:
		if oldState.SlashEpoch == -1 && proposal.EndEpoch >= m.upgradeEpoch {
			// FIP: find the corresponding deal proposal object and extract the provider's actor ID;
			// - we do this by collecting all dealIDs in providerSectors in miner migration

			// in the provider's miner state, find the ID of the sector with the corresponding deal ID in sector metadata;
			sid, ok := m.providerSectors.dealToSector[deal]
			if ok {
				newState.SectorNumber = sid.Number // FIP: set the new deal state object's sector number to the sector ID found;
			}
			//else {
			// TODO: This SHOULD be a fail if we ever get here, but we seem to do so on mainnet.
			// The theory is that because of the "ghost deals" bug, but further investigation is needed.
			//fmt.Println("SUSPECTED GHOST DEAL: ", deal)
			//}
		}

		if err := newStateArray.Set(uint64(deal), &newState); err != nil {
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

	for miner, sectors := range m.providerSectors.minerToSectorToDeals {
		actorSectors, err := adt.MakeEmptyMap(ctxStore, market13.ProviderSectorsHamtBitwidth)
		if err != nil {
			return cid.Undef, cid.Undef, xerrors.Errorf("failed to create empty map: %w", err)
		}

		var sectorDealIDs market13.SectorDealIDs
		for sector, deals := range sectors {
			sectorDealIDs = deals

			if err := actorSectors.Put(miner13.SectorKey(sector), &sectorDealIDs); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to put sector: %w", err)
			}
		}

		if err := outProviderSectors.Put(abi.UIntKey(uint64(miner)), actorSectors); err != nil {
			return cid.Undef, cid.Undef, xerrors.Errorf("failed to put actor sectors: %w", err)
		}
	}

	providerSectorsRoot, err := outProviderSectors.Root()
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to get providerSectorsRoot: %w", err)
	}

	return providerSectorsRoot, newStateArrayRoot, nil
}

func (m *marketMigrator) Deferred() bool {
	return true
}

var _ migration.ActorMigration = (*marketMigrator)(nil)
