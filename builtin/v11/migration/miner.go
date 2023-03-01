package migration

import (
	"context"

	miner10 "github.com/filecoin-project/go-state-types/builtin/v10/miner"
	adt10 "github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
	miner11 "github.com/filecoin-project/go-state-types/builtin/v11/miner"
	adt11 "github.com/filecoin-project/go-state-types/builtin/v11/util/adt"
	"github.com/filecoin-project/go-state-types/migration"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-amt-ipld/v4"
	"golang.org/x/xerrors"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"

	"github.com/filecoin-project/go-state-types/abi"
)

// The minerMigrator performs the following migrations:
// #914 fix: Set ActivationEpoch to when the sector was FIRST activated, and PowerBaseEpoch to latest update epoch
type minerMigrator struct {
	emptyDeadlineV10  cid.Cid
	emptyDeadlinesV10 cid.Cid
	emptyDeadlineV11  cid.Cid
	emptyDeadlinesV11 cid.Cid
	OutCodeCID        cid.Cid
}

func newMinerMigrator(ctx context.Context, store cbor.IpldStore, outCode cid.Cid) (*minerMigrator, error) {
	ctxStore := adt10.WrapStore(ctx, store)

	edv10, err := miner10.ConstructDeadline(ctxStore)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct empty deadline v10: %w", err)
	}

	edv10cid, err := store.Put(ctx, edv10)
	if err != nil {
		return nil, xerrors.Errorf("failed to put empty deadline v10: %w", err)
	}

	edsv10 := miner10.ConstructDeadlines(edv10cid)
	edsv10cid, err := store.Put(ctx, edsv10)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct empty deadlines v10: %w", err)
	}

	edv11, err := miner11.ConstructDeadline(ctxStore)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct empty deadline v11: %w", err)
	}

	edv11cid, err := store.Put(ctx, edv11)
	if err != nil {
		return nil, xerrors.Errorf("failed to put empty deadline v11: %w", err)
	}

	edsv11 := miner11.ConstructDeadlines(edv11cid)
	edsv11cid, err := store.Put(ctx, edsv11)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct empty deadlines v11: %w", err)

	}

	return &minerMigrator{
		emptyDeadlineV10:  edv10cid,
		emptyDeadlinesV10: edsv10cid,
		emptyDeadlineV11:  edv11cid,
		emptyDeadlinesV11: edsv11cid,
		OutCodeCID:        outCode,
	}, nil
}

func (m minerMigrator) MigratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m minerMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (*migration.ActorMigrationResult, error) {
	var inState miner10.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, err
	}
	var inInfo miner10.MinerInfo
	if err := store.Get(ctx, inState.Info, &inInfo); err != nil {
		return nil, err
	}
	wrappedStore := adt10.WrapStore(ctx, store)

	newSectors, err := migrateSectorsWithCache(ctx, wrappedStore, in.Cache, in.Address, inState.Sectors)
	if err != nil {
		return nil, xerrors.Errorf("failed to migrate sectors for miner: %s: %w", in.Address, err)
	}

	newDeadlines, err := m.migrateDeadlines(ctx, wrappedStore, in.Cache, inState.Deadlines)
	if err != nil {
		return nil, xerrors.Errorf("failed to migrate deadlines: %w", err)
	}

	outState := miner11.State{
		Info:                       inState.Info,
		PreCommitDeposits:          inState.PreCommitDeposits,
		LockedFunds:                inState.LockedFunds,
		VestingFunds:               inState.VestingFunds,
		FeeDebt:                    inState.FeeDebt,
		InitialPledge:              inState.InitialPledge,
		PreCommittedSectors:        inState.PreCommittedSectors,
		PreCommittedSectorsCleanUp: inState.PreCommittedSectorsCleanUp,
		AllocatedSectors:           inState.AllocatedSectors,
		Sectors:                    newSectors,
		ProvingPeriodStart:         inState.ProvingPeriodStart,
		CurrentDeadline:            inState.CurrentDeadline,
		Deadlines:                  newDeadlines,
		EarlyTerminations:          inState.EarlyTerminations,
		DeadlineCronActive:         inState.DeadlineCronActive,
	}

	newHead, err := store.Put(ctx, &outState)
	return &migration.ActorMigrationResult{
		NewCodeCID: m.MigratedCodeCID(),
		NewHead:    newHead,
	}, err
}

func migrateSectorsWithCache(ctx context.Context, store adt10.Store, cache migration.MigrationCache, minerAddr address.Address, inRoot cid.Cid) (cid.Cid, error) {
	return cache.Load(migration.SectorsAmtKey(inRoot), func() (cid.Cid, error) {
		inArray, err := adt10.AsArray(store, inRoot, miner10.SectorsAmtBitwidth)
		if err != nil {
			return cid.Undef, xerrors.Errorf("failed to read sectors array: %w", err)
		}

		okIn, prevInRoot, err := cache.Read(migration.MinerPrevSectorsInKey(minerAddr))
		if err != nil {
			return cid.Undef, xerrors.Errorf("failed to get previous inRoot from cache: %w", err)
		}

		okOut, prevOutRoot, err := cache.Read(migration.MinerPrevSectorsOutKey(minerAddr))
		if err != nil {
			return cid.Undef, xerrors.Errorf("failed to get previous outRoot from cache: %w", err)
		}

		var outArray *adt11.Array
		if okIn && okOut {
			// we have previous work, but the AMT has changed -- diff them
			diffs, err := amt.Diff(ctx, store, store, prevInRoot, inRoot, amt.UseTreeBitWidth(miner11.SectorsAmtBitwidth))
			if err != nil {
				return cid.Undef, xerrors.Errorf("failed to diff old and new Sector AMTs: %w", err)
			}

			inSectors, err := miner10.LoadSectors(store, inRoot)
			if err != nil {
				return cid.Undef, xerrors.Errorf("failed to load inSectors: %w", err)
			}

			prevOutSectors, err := miner11.LoadSectors(store, prevOutRoot)
			if err != nil {
				return cid.Undef, xerrors.Errorf("failed to load prevOutSectors: %w", err)
			}

			for _, change := range diffs {
				switch change.Type {
				case amt.Remove:
					if err := prevOutSectors.Delete(change.Key); err != nil {
						return cid.Undef, xerrors.Errorf("failed to delete sector from prevOutSectors: %w", err)
					}
				case amt.Add:
					fallthrough
				case amt.Modify:
					sectorNo := abi.SectorNumber(change.Key)
					info, found, err := inSectors.Get(sectorNo)
					if err != nil {
						return cid.Undef, xerrors.Errorf("failed to get sector %d in inSectors: %w", sectorNo, err)
					}

					if !found {
						return cid.Undef, xerrors.Errorf("didn't find sector %d in inSectors", sectorNo)
					}

					if err := prevOutSectors.Set(change.Key, migrateSectorInfo(*info)); err != nil {
						return cid.Undef, xerrors.Errorf("failed to set migrated sector %d in prevOutSectors", sectorNo)
					}
				}
			}

			outArray = prevOutSectors.Array
		} else {
			// first time we're doing this, do all the work
			outArray, err = migrateSectorsFromScratch(ctx, store, inArray)
			if err != nil {
				return cid.Undef, xerrors.Errorf("failed to migrate sectors from scratch: %w", err)
			}
		}

		outRoot, err := outArray.Root()
		if err != nil {
			return cid.Undef, xerrors.Errorf("error writing new sectors AMT: %w", err)
		}

		if err = cache.Write(migration.MinerPrevSectorsInKey(minerAddr), inRoot); err != nil {
			return cid.Undef, xerrors.Errorf("failed to write inkey to cache: %w", err)
		}

		if err = cache.Write(migration.MinerPrevSectorsOutKey(minerAddr), outRoot); err != nil {
			return cid.Undef, xerrors.Errorf("failed to write inkey to cache: %w", err)
		}

		return outRoot, nil
	})
}

func migrateSectorsFromScratch(ctx context.Context, store adt10.Store, inArray *adt10.Array) (*adt11.Array, error) {
	outArray, err := adt11.MakeEmptyArray(store, miner11.SectorsAmtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct new sectors array: %w", err)
	}

	var sectorInfo miner10.SectorOnChainInfo
	if err = inArray.ForEach(&sectorInfo, func(k int64) error {
		return outArray.Set(uint64(k), migrateSectorInfo(sectorInfo))
	}); err != nil {
		return nil, err
	}

	return outArray, err
}

func (m minerMigrator) migrateDeadlines(ctx context.Context, store adt10.Store, cache migration.MigrationCache, deadlines cid.Cid) (cid.Cid, error) {
	if deadlines == m.emptyDeadlinesV10 {
		return m.emptyDeadlinesV11, nil
	}

	var inDeadlines miner10.Deadlines
	err := store.Get(store.Context(), deadlines, &inDeadlines)
	if err != nil {
		return cid.Undef, err
	}

	var outDeadlines miner11.Deadlines
	for i, c := range inDeadlines.Due {
		if c == m.emptyDeadlineV10 {
			outDeadlines.Due[i] = m.emptyDeadlineV11
		} else {
			var inDeadline miner10.Deadline
			if err = store.Get(ctx, c, &inDeadline); err != nil {
				return cid.Undef, err
			}

			outSectorsSnapshotCid, err := cache.Load(migration.SectorsAmtKey(inDeadline.SectorsSnapshot), func() (cid.Cid, error) {
				inSectorsSnapshot, err := adt10.AsArray(store, inDeadline.SectorsSnapshot, miner10.SectorsAmtBitwidth)
				if err != nil {
					return cid.Undef, err
				}

				outSectorsSnapshot, err := migrateSectorsFromScratch(ctx, store, inSectorsSnapshot)
				if err != nil {
					return cid.Undef, xerrors.Errorf("failed to migrate sectors: %w", err)
				}

				return outSectorsSnapshot.Root()
			})

			if err != nil {
				return cid.Undef, xerrors.Errorf("failed to migrate sectors snapshot: %w", err)
			}

			outDeadline := miner11.Deadline{
				Partitions:                        inDeadline.Partitions,
				ExpirationsEpochs:                 inDeadline.ExpirationsEpochs,
				PartitionsPoSted:                  inDeadline.PartitionsPoSted,
				EarlyTerminations:                 inDeadline.EarlyTerminations,
				LiveSectors:                       inDeadline.LiveSectors,
				TotalSectors:                      inDeadline.TotalSectors,
				FaultyPower:                       miner11.PowerPair(inDeadline.FaultyPower),
				OptimisticPoStSubmissions:         inDeadline.OptimisticPoStSubmissions,
				SectorsSnapshot:                   outSectorsSnapshotCid,
				PartitionsSnapshot:                inDeadline.PartitionsSnapshot,
				OptimisticPoStSubmissionsSnapshot: inDeadline.OptimisticPoStSubmissionsSnapshot,
			}

			outDlCid, err := store.Put(ctx, &outDeadline)
			if err != nil {
				return cid.Undef, err
			}

			outDeadlines.Due[i] = outDlCid
		}
	}

	return store.Put(ctx, &outDeadlines)
}

func migrateSectorInfo(sectorInfo miner10.SectorOnChainInfo) *miner11.SectorOnChainInfo {
	// For a sector that has not been updated: the Activation is correct and ReplacedSectorAge is zero.
	// For a sector that has been updated: Activation is the epoch at which it was upgraded, and ReplacedSectorAge is delta since the true activation.
	//
	// Thus we want to set:
	//
	// Activation := OldActivation - ReplacedSectorAge (a no-op for non-upgraded sectors)
	// PowerBaseEpoch := Activation (in both upgrade and not-upgraded cases)

	powerBaseEpoch := sectorInfo.Activation
	activationEpoch := sectorInfo.Activation - sectorInfo.ReplacedSectorAge

	return &miner11.SectorOnChainInfo{
		SectorNumber:          sectorInfo.SectorNumber,
		SealProof:             sectorInfo.SealProof,
		SealedCID:             sectorInfo.SealedCID,
		DealIDs:               sectorInfo.DealIDs,
		Activation:            activationEpoch,
		Expiration:            sectorInfo.Expiration,
		DealWeight:            sectorInfo.DealWeight,
		VerifiedDealWeight:    sectorInfo.VerifiedDealWeight,
		InitialPledge:         sectorInfo.InitialPledge,
		ExpectedDayReward:     sectorInfo.ExpectedDayReward,
		ExpectedStoragePledge: sectorInfo.ExpectedStoragePledge,
		PowerBaseEpoch:        powerBaseEpoch,
		ReplacedDayReward:     sectorInfo.ReplacedDayReward,
		SectorKeyCID:          sectorInfo.SectorKeyCID,
		SimpleQAPower:         sectorInfo.SimpleQAPower,
	}
}
