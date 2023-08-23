package migration

import (
	"context"

	miner11 "github.com/filecoin-project/go-state-types/builtin/v11/miner"
	adt11 "github.com/filecoin-project/go-state-types/builtin/v11/util/adt"
	miner12 "github.com/filecoin-project/go-state-types/builtin/v12/miner"
	adt12 "github.com/filecoin-project/go-state-types/builtin/v12/util/adt"
	"github.com/filecoin-project/go-state-types/migration"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-amt-ipld/v4"
	"golang.org/x/xerrors"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"

	"github.com/filecoin-project/go-state-types/abi"
)

// The minerMigrator performs the following migrations:
// - FIP-0061: Updates all miner info PoSt proof types to V1_1 types
// - #914 fix: Set ActivationEpoch to when the sector was FIRST activated, and PowerBaseEpoch to latest update epoch

type minerMigrator struct {
	emptyDeadlineV11  cid.Cid
	emptyDeadlinesV11 cid.Cid
	emptyDeadlineV12  cid.Cid
	emptyDeadlinesV12 cid.Cid
	OutCodeCID        cid.Cid
}

func newMinerMigrator(ctx context.Context, store cbor.IpldStore, outCode cid.Cid) (*minerMigrator, error) {
	ctxStore := adt11.WrapStore(ctx, store)

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

	edv12, err := miner12.ConstructDeadline(ctxStore)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct empty deadline v12: %w", err)
	}

	edv12cid, err := store.Put(ctx, edv12)
	if err != nil {
		return nil, xerrors.Errorf("failed to put empty deadline v12: %w", err)
	}

	edsv12 := miner12.ConstructDeadlines(edv12cid)
	edsv12cid, err := store.Put(ctx, edsv12)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct empty deadlines v12: %w", err)

	}

	return &minerMigrator{
		emptyDeadlineV11:  edv11cid,
		emptyDeadlinesV11: edsv11cid,
		emptyDeadlineV12:  edv12cid,
		emptyDeadlinesV12: edsv12cid,
		OutCodeCID:        outCode,
	}, nil
}

func (m minerMigrator) MigratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m minerMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (*migration.ActorMigrationResult, error) {
	var inState miner11.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load miner state for %s: %w", in.Address, err)
	}

	var inInfo miner11.MinerInfo
	if err := store.Get(ctx, inState.Info, &inInfo); err != nil {
		return nil, xerrors.Errorf("failed to load miner info for %s: %w", in.Address, err)
	}

	outProof, err := inInfo.WindowPoStProofType.ToV1_1PostProof()
	if err != nil {
		return nil, xerrors.Errorf("failed to convert to v1_1 proof: %w", err)
	}

	outInfo := miner12.MinerInfo{
		Owner:                      inInfo.Owner,
		Worker:                     inInfo.Worker,
		ControlAddresses:           inInfo.ControlAddresses,
		PendingWorkerKey:           (*miner12.WorkerKeyChange)(inInfo.PendingWorkerKey),
		PeerId:                     inInfo.PeerId,
		Multiaddrs:                 inInfo.Multiaddrs,
		WindowPoStProofType:        outProof,
		SectorSize:                 inInfo.SectorSize,
		WindowPoStPartitionSectors: inInfo.WindowPoStPartitionSectors,
		ConsensusFaultElapsed:      inInfo.ConsensusFaultElapsed,
		PendingOwnerAddress:        inInfo.PendingOwnerAddress,
		Beneficiary:                inInfo.Beneficiary,
		BeneficiaryTerm:            miner12.BeneficiaryTerm(inInfo.BeneficiaryTerm),
		PendingBeneficiaryTerm:     (*miner12.PendingBeneficiaryChange)(inInfo.PendingBeneficiaryTerm),
	}

	outInfoCid, err := store.Put(ctx, &outInfo)
	if err != nil {
		return nil, xerrors.Errorf("failed to write new miner info: %w", err)
	}

	wrappedStore := adt11.WrapStore(ctx, store)

	newSectors, err := migrateSectorsWithCache(ctx, wrappedStore, in.Cache, in.Address, inState.Sectors)
	if err != nil {
		return nil, xerrors.Errorf("failed to migrate sectors for miner: %s: %w", in.Address, err)
	}

	newDeadlines, err := m.migrateDeadlines(ctx, wrappedStore, in.Cache, inState.Deadlines)
	if err != nil {
		return nil, xerrors.Errorf("failed to migrate deadlines: %w", err)
	}

	outState := miner12.State{
		Info:                       outInfoCid,
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
	if err != nil {
		return nil, xerrors.Errorf("failed to put new state: %w", err)
	}

	return &migration.ActorMigrationResult{
		NewCodeCID: m.MigratedCodeCID(),
		NewHead:    newHead,
	}, nil
}

func migrateSectorsWithCache(ctx context.Context, store adt11.Store, cache migration.MigrationCache, minerAddr address.Address, inRoot cid.Cid) (cid.Cid, error) {
	return cache.Load(migration.SectorsAmtKey(inRoot), func() (cid.Cid, error) {
		inArray, err := adt11.AsArray(store, inRoot, miner11.SectorsAmtBitwidth)
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

		var outArray *adt12.Array
		if okIn && okOut {
			// we have previous work, but the AMT has changed -- diff them
			diffs, err := amt.Diff(ctx, store, store, prevInRoot, inRoot, amt.UseTreeBitWidth(miner12.SectorsAmtBitwidth))
			if err != nil {
				return cid.Undef, xerrors.Errorf("failed to diff old and new Sector AMTs: %w", err)
			}

			inSectors, err := miner11.LoadSectors(store, inRoot)
			if err != nil {
				return cid.Undef, xerrors.Errorf("failed to load inSectors: %w", err)
			}

			prevOutSectors, err := miner12.LoadSectors(store, prevOutRoot)
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

func migrateSectorsFromScratch(ctx context.Context, store adt11.Store, inArray *adt11.Array) (*adt12.Array, error) {
	outArray, err := adt12.MakeEmptyArray(store, miner12.SectorsAmtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct new sectors array: %w", err)
	}

	var sectorInfo miner11.SectorOnChainInfo
	if err = inArray.ForEach(&sectorInfo, func(k int64) error {
		return outArray.Set(uint64(k), migrateSectorInfo(sectorInfo))
	}); err != nil {
		return nil, err
	}

	return outArray, err
}

func (m minerMigrator) migrateDeadlines(ctx context.Context, store adt11.Store, cache migration.MigrationCache, deadlines cid.Cid) (cid.Cid, error) {
	if deadlines == m.emptyDeadlinesV11 {
		return m.emptyDeadlinesV12, nil
	}

	var inDeadlines miner11.Deadlines
	err := store.Get(store.Context(), deadlines, &inDeadlines)
	if err != nil {
		return cid.Undef, err
	}

	var outDeadlines miner12.Deadlines
	for i, c := range inDeadlines.Due {
		if c == m.emptyDeadlineV11 {
			outDeadlines.Due[i] = m.emptyDeadlineV12
		} else {
			var inDeadline miner11.Deadline
			if err = store.Get(ctx, c, &inDeadline); err != nil {
				return cid.Undef, err
			}

			outSectorsSnapshotCid, err := cache.Load(migration.SectorsAmtKey(inDeadline.SectorsSnapshot), func() (cid.Cid, error) {
				inSectorsSnapshot, err := adt11.AsArray(store, inDeadline.SectorsSnapshot, miner11.SectorsAmtBitwidth)
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

			outDeadline := miner12.Deadline{
				Partitions:                        inDeadline.Partitions,
				ExpirationsEpochs:                 inDeadline.ExpirationsEpochs,
				PartitionsPoSted:                  inDeadline.PartitionsPoSted,
				EarlyTerminations:                 inDeadline.EarlyTerminations,
				LiveSectors:                       inDeadline.LiveSectors,
				TotalSectors:                      inDeadline.TotalSectors,
				FaultyPower:                       miner12.PowerPair(inDeadline.FaultyPower),
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

func migrateSectorInfo(sectorInfo miner11.SectorOnChainInfo) *miner12.SectorOnChainInfo {
	// For a sector that has not been updated: the Activation is correct and ReplacedSectorAge is zero.
	// For a sector that has been updated through SnapDeals: Activation is the epoch at which it was upgraded, and ReplacedSectorAge is delta since the true activation.
	// For a sector that has been updated through the old CC path: Activation is correct
	// Thus, we want to set:
	//
	// PowerBaseEpoch := Activation (in all cases)
	// Activation := Activation (for non-upgraded sectors and sectors upgraded through old CC path)
	// Activation := OldActivation - ReplacedSectorAge (for sectors updated through SnapDeals)

	powerBaseEpoch := sectorInfo.Activation
	activationEpoch := sectorInfo.Activation
	if sectorInfo.SectorKeyCID != nil {
		activationEpoch = sectorInfo.Activation - sectorInfo.ReplacedSectorAge
	}

	return &miner12.SectorOnChainInfo{
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
