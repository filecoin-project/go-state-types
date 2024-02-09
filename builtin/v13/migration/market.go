package migration

import (
	"bytes"
	"context"
	"errors"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-amt-ipld/v4"
	"github.com/filecoin-project/go-state-types/abi"
	market12 "github.com/filecoin-project/go-state-types/builtin/v12/market"
	market13 "github.com/filecoin-project/go-state-types/builtin/v13/market"
	miner13 "github.com/filecoin-project/go-state-types/builtin/v13/miner"
	adt13 "github.com/filecoin-project/go-state-types/builtin/v13/util/adt"
	"github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	typegen "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/xerrors"
)

var errItemFound = errors.New("item found")

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
		providerSectorsRoot, newStateArrayRoot, err = m.migrateProviderSectorsAndStatesWithDiff(ctx, store, prevInStates, prevOutStates, prevOutProviderSectors, states, prevInProposals)
		if err != nil {
			return cid.Undef, cid.Undef, xerrors.Errorf("failed to migrate provider sectors (diff): %w", err)
		}
	} else {
		providerSectorsRoot, newStateArrayRoot, err = m.migrateProviderSectorsAndStatesFromScratch(ctx, store, in, states)
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

func (m *marketMigrator) migrateProviderSectorsAndStatesWithDiff(ctx context.Context, store cbor.IpldStore, prevInStatesCid, prevOutStatesCid, prevOutProviderSectorsCid, inStatesCid, proposals cid.Cid) (cid.Cid, cid.Cid, error) {
	diffs, err := amt.Diff(ctx, store, store, prevInStatesCid, inStatesCid, amt.UseTreeBitWidth(market12.StatesAmtBitwidth))
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

	proposalsArr, err := adt.AsArray(ctxStore, proposals, market12.ProposalsAmtBitwidth)
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to load proposals map: %w", err)
	}

	// in-memory maps with changesets to be applied to prevOutProviderSectors
	providerSectorsMem := map[abi.ActorID]map[abi.SectorNumber][]abi.DealID{}        // added
	providerSectorsMemRemoved := map[abi.ActorID]map[abi.SectorNumber][]abi.DealID{} // removed

	addProviderSectorEntry := func(deal abi.DealID) (abi.SectorNumber, error) {
		sid, ok := m.providerSectors.dealToSector[deal]
		if !ok {
			return 0, xerrors.Errorf("deal %d not found in providerSectors", deal) // todo is this normal and possible??
		}

		if _, ok := providerSectorsMem[sid.Miner]; !ok {
			providerSectorsMem[sid.Miner] = make(map[abi.SectorNumber][]abi.DealID)
		}
		providerSectorsMem[sid.Miner][sid.Number] = append(providerSectorsMem[sid.Miner][sid.Number], deal)

		return sid.Number, nil
	}

	removeProviderSectorEntry := func(deal abi.DealID, newStatePrevState *market13.DealState) error {
		snum := newStatePrevState.SectorNumber

		var proposals market12.DealProposal
		found, err := proposalsArr.Get(uint64(deal), &proposals)
		if err != nil {
			return xerrors.Errorf("failed to get proposal for removed deal: %w", err)
		}
		if !found {
			return xerrors.Errorf("proposal not found")
		}

		midd, err := address.IDFromAddress(proposals.Provider)
		if err != nil {
			return xerrors.Errorf("failed to get miner ID: %w", err)
		}
		mid := abi.ActorID(midd)

		newStatePrevState.SectorNumber = 0
		if _, ok := providerSectorsMemRemoved[mid]; !ok {
			providerSectorsMemRemoved[mid] = make(map[abi.SectorNumber][]abi.DealID)
		}
		providerSectorsMemRemoved[mid][snum] = append(providerSectorsMemRemoved[mid][snum], deal)

		return nil
	}

	//fmt.Printf("state diffs: %d\n", len(diffs))
	//fmt.Printf("dealToSector: %d\n", len(m.providerSectors.dealToSector))
	//fmt.Printf("removedDealToSector: %d\n", len(m.providerSectors.removedDealToSector))

	for _, change := range diffs {
		deal := abi.DealID(change.Key)

		switch change.Type {
		case amt.Add:
			var oldState market12.DealState
			if err := oldState.UnmarshalCBOR(bytes.NewReader(change.After.Raw)); err != nil {
				return cid.Cid{}, cid.Cid{}, xerrors.Errorf("failed to unmarshal old state: %w", err)
			}

			var newState market13.DealState
			newState.SlashEpoch = oldState.SlashEpoch
			newState.LastUpdatedEpoch = oldState.LastUpdatedEpoch
			newState.SectorStartEpoch = oldState.SectorStartEpoch
			newState.SectorNumber = 0 // terminated / not found (?)

			if oldState.SlashEpoch == -1 {
				si, err := addProviderSectorEntry(deal)
				if err != nil {
					return cid.Cid{}, cid.Cid{}, xerrors.Errorf("failed to add provider sector entry: %w", err)
				}
				newState.SectorNumber = si
			}

			//fmt.Printf("add deal %d to sector %d\n", deal, newState.SectorNumber)

			if err := prevOutStates.Set(uint64(deal), &newState); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to set new state: %w", err)
			}

		case amt.Remove:
			//fmt.Printf("remove deal %d\n", deal)

			var prevOutState market13.DealState // note: this says "prevOut", not "prevOld"
			ok, err := prevOutStates.Get(uint64(deal), &prevOutState)
			if err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to get previous newstate: %w", err)
			}
			if !ok {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to get previous newstate: not found")
			}

			if prevOutState.SlashEpoch == -1 {
				// if the previous OUT state was not slashed then it has a provider sector entry that needs to be removed

				if err := removeProviderSectorEntry(deal, &prevOutState); err != nil {
					return cid.Undef, cid.Undef, xerrors.Errorf("failed to remove provider sector entry: %w", err)
				}
			}

			if err := prevOutStates.Delete(uint64(deal)); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to delete new state: %w", err)
			}

		case amt.Modify:

			var oldState, prevOldState market12.DealState
			var newState market13.DealState
			if err := prevOldState.UnmarshalCBOR(bytes.NewReader(change.Before.Raw)); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to unmarshal old state: %w", err)
			}
			if err := oldState.UnmarshalCBOR(bytes.NewReader(change.After.Raw)); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to unmarshal old state: %w", err)
			}
			ok, err := prevOutStates.Get(uint64(deal), &newState)
			if err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to get previous newstate: %w", err)
			}
			if !ok {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to get previous newstate: not found")
			}

			newState.SlashEpoch = oldState.SlashEpoch
			newState.LastUpdatedEpoch = oldState.LastUpdatedEpoch
			newState.SectorStartEpoch = oldState.SectorStartEpoch

			// if nowOld.Slash == -1, then 'now' is not slashed, so we should try to find the sector
			// we probably don't care about prevOldSlash?? beyond it changing from newSlash?

			//fmt.Printf("deal %d slash %d -> %d, update %d -> %d (prev sec: %d)\n", deal, prevOldState.SlashEpoch, oldState.SlashEpoch, prevOldState.LastUpdatedEpoch, oldState.LastUpdatedEpoch, newState.SectorNumber)

			if prevOldState.SlashEpoch == -1 && oldState.SlashEpoch != -1 {
				// not slashed -> slashed
				//fmt.Printf("deal %d slash -1 -> %d\n", deal, oldState.SlashEpoch)

				if err := removeProviderSectorEntry(deal, &newState); err != nil {
					return cid.Cid{}, cid.Cid{}, xerrors.Errorf("failed to remove provider sector entry: %w", err)
				}
			}

			if err := prevOutStates.Set(uint64(deal), &newState); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to set new state: %w", err)
			}
		}
	}

	// process prevOutProviderSectors, first removes, then adds

	var actorSectorsMapRoot typegen.CborCid
	var sectorDeals market13.SectorDealIDs

	for miner, sectors := range providerSectorsMemRemoved {
		found, err := prevOutProviderSectors.Get(abi.UIntKey(uint64(miner)), &actorSectorsMapRoot)
		if err != nil {
			return cid.Undef, cid.Undef, xerrors.Errorf("failed to get actor sectors: %w", err)
		}
		if !found {
			// this is fine, all sectors of this miner were already not present
			// in ProviderSectors. Sadly because the default value of a non-present
			// sector number in deal state is 0, we can't tell if a sector was
			// removed or if it was never there to begin with, which is why we
			// may occasionally end up here.

			//fmt.Printf("no actor sectors for miner %d\n", miner)

			continue
		}

		actorSectors, err := adt.AsMap(ctxStore, cid.Cid(actorSectorsMapRoot), market13.ProviderSectorsHamtBitwidth)
		if err != nil {
			return cid.Undef, cid.Undef, xerrors.Errorf("failed to load actor sectors map: %w", err)
		}

		for sector, deals := range sectors {
			var sectorDeals market13.SectorDealIDs // []abi.DealID
			found, err := actorSectors.Get(miner13.SectorKey(sector), &sectorDeals)
			if err != nil {
				return cid.Cid{}, cid.Cid{}, xerrors.Errorf("failed to get sector: %w", err)
			}
			if !found {
				continue
			}

			// remove deals from sectorDeals
			for _, deal := range deals {
				for j, d := range sectorDeals {
					if d == deal {
						// delete j-th element
						sectorDeals = append(sectorDeals[:j], sectorDeals[j+1:]...)
						break
					}
				}
			}

			if len(sectorDeals) == 0 {
				//fmt.Printf("delete providersector %d, deals %v\n", sector, deals)
				if err := actorSectors.Delete(miner13.SectorKey(sector)); err != nil {
					return cid.Cid{}, cid.Cid{}, xerrors.Errorf("failed to delete sector: %w", err)
				}
			} else {
				//fmt.Printf("update providersector %d, deals %v\n", sector, sectorDeals)
				if err := actorSectors.Put(miner13.SectorKey(sector), &sectorDeals); err != nil {
					return cid.Cid{}, cid.Cid{}, xerrors.Errorf("failed to put sector: %w", err)
				}
			}
		}

		// check if actorSectors are empty
		err = actorSectors.ForEach(nil, func(k string) error {
			return errItemFound
		})
		var nonEmpty bool
		if err == errItemFound {
			nonEmpty = true
			err = nil
		}
		if err != nil {
			return cid.Undef, cid.Undef, xerrors.Errorf("failed to iterate actor sectors: %w", err)
		}

		if nonEmpty {
			newActorSectorsMapRoot, err := actorSectors.Root()
			if err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to get actor sectors root: %w", err)
			}
			actorSectorsMapRoot = typegen.CborCid(newActorSectorsMapRoot)

			if err := prevOutProviderSectors.Put(abi.UIntKey(uint64(miner)), &actorSectorsMapRoot); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to put actor sectors: %w", err)
			}
		} else {
			if err := prevOutProviderSectors.Delete(abi.UIntKey(uint64(miner))); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to delete actor sectors: %w", err)
			}
		}
	}

	for miner, sectors := range providerSectorsMem {
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

		for sector, deals := range sectors {
			sectorDeals = deals

			if err := actorSectors.Put(miner13.SectorKey(sector), &sectorDeals); err != nil {
				return cid.Undef, cid.Undef, xerrors.Errorf("failed to put sector: %w", err)
			}
		}

		newActorSectorsMapRoot, err := actorSectors.Root()
		if err != nil {
			return cid.Undef, cid.Undef, xerrors.Errorf("failed to get actor sectors root: %w", err)
		}
		actorSectorsMapRoot = typegen.CborCid(newActorSectorsMapRoot)

		if err := prevOutProviderSectors.Put(abi.UIntKey(uint64(miner)), &actorSectorsMapRoot); err != nil {
			return cid.Undef, cid.Undef, xerrors.Errorf("failed to put actor sectors: %w", err)
		}
	}

	// get roots
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

func (m *marketMigrator) migrateProviderSectorsAndStatesFromScratch(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput, states cid.Cid) (cid.Cid, cid.Cid, error) {
	ctxStore := adt.WrapStore(ctx, store)

	oldStateArray, err := adt.AsArray(ctxStore, states, market12.StatesAmtBitwidth)
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to load states array: %w", err)
	}

	newStateArray, err := adt13.MakeEmptyArray(ctxStore, market13.StatesAmtBitwidth)
	if err != nil {
		return cid.Undef, cid.Undef, xerrors.Errorf("failed to load states array: %w", err)
	}

	providerSectorsMem := map[abi.ActorID]map[abi.SectorNumber][]abi.DealID{}

	var oldState market12.DealState
	var newState market13.DealState

	// FIP: For each deal state object in the market actor state that has a terminated epoch set to -1:
	err = oldStateArray.ForEach(&oldState, func(i int64) error {
		deal := abi.DealID(i)

		newState.SlashEpoch = oldState.SlashEpoch
		newState.LastUpdatedEpoch = oldState.LastUpdatedEpoch
		newState.SectorStartEpoch = oldState.SectorStartEpoch
		newState.SectorNumber = 0 // non- -1 slashed epoch

		if oldState.SlashEpoch == -1 {
			// FIP: find the corresponding deal proposal object and extract the provider's actor ID;
			// - we do this by collecting all dealIDs in providerSectors in miner migration

			// in the provider's miner state, find the ID of the sector with the corresponding deal ID in sector metadata;
			sid, ok := m.providerSectors.dealToSector[deal]
			if ok {
				newState.SectorNumber = sid.Number // FIP: set the new deal state object's sector number to the sector ID found;

				// FIP: add the deal ID to the ProviderSectors mapping for the provider's actor ID and sector number.
				if _, ok := providerSectorsMem[sid.Miner]; !ok {
					providerSectorsMem[sid.Miner] = make(map[abi.SectorNumber][]abi.DealID)
				}
				providerSectorsMem[sid.Miner][sid.Number] = append(providerSectorsMem[sid.Miner][sid.Number], deal)
			} else {
				//fmt.Printf("deal %d not found in providerSectors: %v\n", deal, oldState)

				newState.SectorNumber = 0 // FIP: if such a sector cannot be found, assert that the deal's end epoch has passed and use sector ID 0
			}
		} /*else {
			fmt.Printf("deal %d slashed, not inserting ProviderSectors record: %v\n", deal, oldState)
		}*/

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

	for miner, sectors := range providerSectorsMem {
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
