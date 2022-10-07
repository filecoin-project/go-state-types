package migration

import (
	"context"

	"github.com/filecoin-project/go-state-types/exitcode"

	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-amt-ipld/v4"
	"golang.org/x/xerrors"

	commp "github.com/filecoin-project/go-commp-utils/nonffi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v8/market"
	miner8 "github.com/filecoin-project/go-state-types/builtin/v8/miner"
	adt8 "github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
	miner9 "github.com/filecoin-project/go-state-types/builtin/v9/miner"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"

	"github.com/filecoin-project/go-state-types/abi"
)

// The minerMigrator performs the following migrations:
// FIP-0029: Sets the Beneficary to the Owner, and sets empty values for BeneficiaryTerm and PendingBeneficiaryTerm
// FIP-0034: For each SectorPreCommitOnChainInfo in PreCommitedSectors, calculates the unsealed CID (assuming there are deals)
// FIP-0045: For each SectorOnChainInfo in Sectors, set SimpleQAPower = (DealWeight == 0 && VerifiedDealWeight == 0)
// FIP-0045: For each Deadline in Deadlines: for each SectorOnChainInfo in SectorsSnapshot, set SimpleQAPower = (DealWeight == 0 && VerifiedDealWeight == 0)

type minerMigrator struct {
	emptyPrecommitOnChainInfosV9 cid.Cid
	emptyDeadlineV8              cid.Cid
	emptyDeadlinesV8             cid.Cid
	emptyDeadlineV9              cid.Cid
	emptyDeadlinesV9             cid.Cid
	proposals                    *market.DealArray
	OutCodeCID                   cid.Cid
}

func newMinerMigrator(ctx context.Context, store cbor.IpldStore, marketProposals *market.DealArray, outCode cid.Cid) (*minerMigrator, error) {
	ctxStore := adt8.WrapStore(ctx, store)

	emptyPrecommitMapCidV9, err := adt9.StoreEmptyMap(ctxStore, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct empty precommit map v9: %w", err)
	}

	edv8, err := miner8.ConstructDeadline(ctxStore)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct empty deadline v8: %w", err)
	}

	edv8cid, err := store.Put(ctx, edv8)
	if err != nil {
		return nil, xerrors.Errorf("failed to put empty deadline v8: %w", err)
	}

	edsv8 := miner8.ConstructDeadlines(edv8cid)
	edsv8cid, err := store.Put(ctx, edsv8)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct empty deadlines v8: %w", err)
	}

	edv9, err := miner9.ConstructDeadline(ctxStore)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct empty deadline v9: %w", err)
	}

	edv9cid, err := store.Put(ctx, edv9)
	if err != nil {
		return nil, xerrors.Errorf("failed to put empty deadline v9: %w", err)
	}

	edsv9 := miner9.ConstructDeadlines(edv9cid)
	edsv9cid, err := store.Put(ctx, edsv9)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct empty deadlines v9: %w", err)

	}

	return &minerMigrator{
		emptyPrecommitOnChainInfosV9: emptyPrecommitMapCidV9,
		emptyDeadlineV8:              edv8cid,
		emptyDeadlinesV8:             edsv8cid,
		emptyDeadlineV9:              edv9cid,
		emptyDeadlinesV9:             edsv9cid,
		proposals:                    marketProposals,
		OutCodeCID:                   outCode,
	}, nil
}

func (m minerMigrator) migratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m minerMigrator) migrateState(ctx context.Context, store cbor.IpldStore, in actorMigrationInput) (*actorMigrationResult, error) {
	var inState miner8.State
	if err := store.Get(ctx, in.head, &inState); err != nil {
		return nil, err
	}
	var inInfo miner8.MinerInfo
	if err := store.Get(ctx, inState.Info, &inInfo); err != nil {
		return nil, err
	}
	wrappedStore := adt8.WrapStore(ctx, store)

	newPrecommits, err := m.migratePrecommits(ctx, wrappedStore, inState.PreCommittedSectors)
	if err != nil {
		return nil, xerrors.Errorf("failed to migrate precommits for miner: %s: %w", in.address, err)
	}

	newSectors, err := migrateSectorsWithCache(ctx, wrappedStore, in.cache, in.address, inState.Sectors)
	if err != nil {
		return nil, xerrors.Errorf("failed to migrate sectors for miner: %s: %w", in.address, err)
	}

	newDeadlines, err := m.migrateDeadlines(ctx, wrappedStore, in.cache, inState.Deadlines)
	if err != nil {
		return nil, xerrors.Errorf("failed to migrate deadlines: %w", err)
	}

	var newPendingWorkerKey *miner9.WorkerKeyChange
	if inInfo.PendingWorkerKey != nil {
		newPendingWorkerKey = &miner9.WorkerKeyChange{
			NewWorker:   inInfo.PendingWorkerKey.NewWorker,
			EffectiveAt: inInfo.PendingWorkerKey.EffectiveAt,
		}
	}

	outInfo := miner9.MinerInfo{
		Owner:       inInfo.Owner,
		Worker:      inInfo.Worker,
		Beneficiary: inInfo.Owner,
		BeneficiaryTerm: miner9.BeneficiaryTerm{
			Quota:      abi.NewTokenAmount(0),
			UsedQuota:  abi.NewTokenAmount(0),
			Expiration: 0,
		},
		PendingBeneficiaryTerm:     nil,
		ControlAddresses:           inInfo.ControlAddresses,
		PendingWorkerKey:           newPendingWorkerKey,
		PeerId:                     inInfo.PeerId,
		Multiaddrs:                 inInfo.Multiaddrs,
		WindowPoStProofType:        inInfo.WindowPoStProofType,
		SectorSize:                 inInfo.SectorSize,
		WindowPoStPartitionSectors: inInfo.WindowPoStPartitionSectors,
		ConsensusFaultElapsed:      inInfo.ConsensusFaultElapsed,
		PendingOwnerAddress:        inInfo.PendingOwnerAddress,
	}
	newInfoCid, err := store.Put(ctx, &outInfo)
	if err != nil {
		return nil, xerrors.Errorf("failed to flush new miner info: %w", err)
	}

	outState := miner9.State{
		Info:                       newInfoCid,
		PreCommitDeposits:          inState.PreCommitDeposits,
		LockedFunds:                inState.LockedFunds,
		VestingFunds:               inState.VestingFunds,
		FeeDebt:                    inState.FeeDebt,
		InitialPledge:              inState.InitialPledge,
		PreCommittedSectors:        newPrecommits,
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
	return &actorMigrationResult{
		newCodeCID: m.migratedCodeCID(),
		newHead:    newHead,
	}, err
}

func (m minerMigrator) migratePrecommits(ctx context.Context, wrappedStore adt8.Store, inRoot cid.Cid) (cid.Cid, error) {
	oldPrecommitOnChainInfos, err := adt8.AsMap(wrappedStore, inRoot, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to load old precommit onchain infos: %w", err)
	}

	newPrecommitOnChainInfos, err := adt9.AsMap(wrappedStore, m.emptyPrecommitOnChainInfosV9, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to load empty map: %w", err)
	}

	var info miner8.SectorPreCommitOnChainInfo
	err = oldPrecommitOnChainInfos.ForEach(&info, func(key string) error {
		var unsealedCid *cid.Cid
		var pieces []abi.PieceInfo
		for _, dealID := range info.Info.DealIDs {
			deal, err := m.proposals.GetDealProposal(dealID)
			if err != nil {
				// Possible for the proposal to be missing if it's expired (but the deal is still in a precommit that's yet to be cleaned up)
				// Just continue in this case, the sector is unProveCommitable anyway, will just fail later
				if exitcode.Unwrap(err, exitcode.ErrIllegalState) != exitcode.ErrNotFound {
					return xerrors.Errorf("error getting deal proposal for sector: %d: %w", info.Info.SectorNumber, err)
				}

				continue
			}

			pieces = append(pieces, abi.PieceInfo{
				PieceCID: deal.PieceCID,
				Size:     deal.PieceSize,
			})
		}

		if len(pieces) != 0 {
			commd, err := commp.GenerateUnsealedCID(info.Info.SealProof, pieces)
			if err != nil {
				return xerrors.Errorf("failed to generate unsealed CID: %w", err)
			}

			unsealedCid = &commd
		}

		err = newPrecommitOnChainInfos.Put(miner9.SectorKey(info.Info.SectorNumber), &miner9.SectorPreCommitOnChainInfo{
			Info: miner9.SectorPreCommitInfo{
				SealProof:     info.Info.SealProof,
				SectorNumber:  info.Info.SectorNumber,
				SealedCID:     info.Info.SealedCID,
				SealRandEpoch: info.Info.SealRandEpoch,
				DealIDs:       info.Info.DealIDs,
				Expiration:    info.Info.Expiration,
				UnsealedCid:   unsealedCid,
			},
			PreCommitDeposit: info.PreCommitDeposit,
			PreCommitEpoch:   info.PreCommitEpoch,
		})

		if err != nil {
			return xerrors.Errorf("failed to write new precommitinfo: %w", err)
		}

		return nil
	})

	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to iterate over precommitinfos: %w", err)
	}

	newPrecommits, err := newPrecommitOnChainInfos.Root()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to flush new precommits: %w", err)
	}

	return newPrecommits, nil
}

func migrateSectorsWithCache(ctx context.Context, store adt8.Store, cache MigrationCache, minerAddr address.Address, inRoot cid.Cid) (cid.Cid, error) {
	return cache.Load(SectorsAmtKey(inRoot), func() (cid.Cid, error) {
		inArray, err := adt8.AsArray(store, inRoot, miner8.SectorsAmtBitwidth)
		if err != nil {
			return cid.Undef, xerrors.Errorf("failed to read sectors array: %w", err)
		}

		okIn, prevInRoot, err := cache.Read(MinerPrevSectorsInKey(minerAddr))
		if err != nil {
			return cid.Undef, xerrors.Errorf("failed to get previous inRoot from cache: %w", err)
		}

		okOut, prevOutRoot, err := cache.Read(MinerPrevSectorsOutKey(minerAddr))
		if err != nil {
			return cid.Undef, xerrors.Errorf("failed to get previous outRoot from cache: %w", err)
		}

		var outArray *adt9.Array
		if okIn && okOut {
			// we have previous work, but the AMT has changed -- diff them
			diffs, err := amt.Diff(ctx, store, store, prevInRoot, inRoot, amt.UseTreeBitWidth(miner9.SectorsAmtBitwidth))
			if err != nil {
				return cid.Undef, xerrors.Errorf("failed to diff old and new Sector AMTs: %w", err)
			}

			inSectors, err := miner8.LoadSectors(store, inRoot)
			if err != nil {
				return cid.Undef, xerrors.Errorf("failed to load inSectors: %w", err)
			}

			prevOutSectors, err := miner9.LoadSectors(store, prevOutRoot)
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

		if err = cache.Write(MinerPrevSectorsInKey(minerAddr), inRoot); err != nil {
			return cid.Undef, xerrors.Errorf("failed to write inkey to cache: %w", err)
		}

		if err = cache.Write(MinerPrevSectorsOutKey(minerAddr), outRoot); err != nil {
			return cid.Undef, xerrors.Errorf("failed to write inkey to cache: %w", err)
		}

		return outRoot, nil
	})
}

func migrateSectorsFromScratch(ctx context.Context, store adt8.Store, inArray *adt8.Array) (*adt9.Array, error) {
	outArray, err := adt9.MakeEmptyArray(store, miner9.SectorsAmtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct new sectors array: %w", err)
	}

	var sectorInfo miner8.SectorOnChainInfo
	if err = inArray.ForEach(&sectorInfo, func(k int64) error {
		return outArray.Set(uint64(k), migrateSectorInfo(sectorInfo))
	}); err != nil {
		return nil, err
	}

	return outArray, err
}

func (m minerMigrator) migrateDeadlines(ctx context.Context, store adt8.Store, cache MigrationCache, deadlines cid.Cid) (cid.Cid, error) {
	if deadlines == m.emptyDeadlinesV8 {
		return m.emptyDeadlinesV9, nil
	}

	var inDeadlines miner8.Deadlines
	err := store.Get(store.Context(), deadlines, &inDeadlines)
	if err != nil {
		return cid.Undef, err
	}

	var outDeadlines miner9.Deadlines
	for i, c := range inDeadlines.Due {
		if c == m.emptyDeadlineV8 {
			outDeadlines.Due[i] = m.emptyDeadlineV9
		} else {
			var inDeadline miner8.Deadline
			if err = store.Get(ctx, c, &inDeadline); err != nil {
				return cid.Undef, err
			}

			outSectorsSnapshotCid, err := cache.Load(SectorsAmtKey(inDeadline.SectorsSnapshot), func() (cid.Cid, error) {
				inSectorsSnapshot, err := adt8.AsArray(store, inDeadline.SectorsSnapshot, miner8.SectorsAmtBitwidth)
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

			outDeadline := miner9.Deadline{
				Partitions:                        inDeadline.Partitions,
				ExpirationsEpochs:                 inDeadline.ExpirationsEpochs,
				PartitionsPoSted:                  inDeadline.PartitionsPoSted,
				EarlyTerminations:                 inDeadline.EarlyTerminations,
				LiveSectors:                       inDeadline.LiveSectors,
				TotalSectors:                      inDeadline.TotalSectors,
				FaultyPower:                       miner9.PowerPair(inDeadline.FaultyPower),
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

func migrateSectorInfo(sectorInfo miner8.SectorOnChainInfo) *miner9.SectorOnChainInfo {
	return &miner9.SectorOnChainInfo{
		SectorNumber:          sectorInfo.SectorNumber,
		SealProof:             sectorInfo.SealProof,
		SealedCID:             sectorInfo.SealedCID,
		DealIDs:               sectorInfo.DealIDs,
		Activation:            sectorInfo.Activation,
		Expiration:            sectorInfo.Expiration,
		DealWeight:            sectorInfo.DealWeight,
		VerifiedDealWeight:    sectorInfo.VerifiedDealWeight,
		InitialPledge:         sectorInfo.InitialPledge,
		ExpectedDayReward:     sectorInfo.ExpectedDayReward,
		ExpectedStoragePledge: sectorInfo.ExpectedStoragePledge,
		ReplacedSectorAge:     sectorInfo.ReplacedSectorAge,
		ReplacedDayReward:     sectorInfo.ReplacedDayReward,
		SectorKeyCID:          sectorInfo.SectorKeyCID,
		SimpleQAPower:         sectorInfo.DealWeight.IsZero() && sectorInfo.VerifiedDealWeight.IsZero(),
	}
}
