package migration

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"
	miner15 "github.com/filecoin-project/go-state-types/builtin/v15/miner"
	miner16 "github.com/filecoin-project/go-state-types/builtin/v16/miner"
	"github.com/filecoin-project/go-state-types/builtin/v16/util/adt"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"
)

type minerMigrator struct {
	OutCodeCID cid.Cid
}

func newMinerMigrator(_ context.Context, _ cbor.IpldStore, outCode cid.Cid) (*minerMigrator, error) {
	return &minerMigrator{
		OutCodeCID: outCode,
	}, nil
}

func (m *minerMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (result *migration.ActorMigrationResult, err error) {
	var inState miner15.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load miner state for %s: %w", in.Address, err)
	}

	ctxStore := adt.WrapStore(ctx, store)

	inSectors, err := miner15.LoadSectors(ctxStore, inState.Sectors)
	if err != nil {
		return nil, xerrors.Errorf("failed to load sectors array: %w", err)
	}

	arr, err := adt.MakeEmptyArray(ctxStore, miner16.SectorsAmtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to create sectors array: %w", err)
	}
	outSectors := miner16.Sectors{Array: arr, Store: ctxStore}

	inSectors.ForEach(func(sn abi.SectorNumber, soci *miner15.SectorOnChainInfo) error {
		return outSectors.Set(sn, &miner16.SectorOnChainInfo{
			SectorNumber:          soci.SectorNumber,
			SealProof:             soci.SealProof,
			SealedCID:             soci.SealedCID,
			DealIDs:               soci.DealIDs,
			Activation:            soci.Activation,
			Expiration:            soci.Expiration,
			DealWeight:            soci.DealWeight,
			VerifiedDealWeight:    soci.VerifiedDealWeight,
			InitialPledge:         soci.InitialPledge,
			ExpectedDayReward:     soci.ExpectedDayReward,
			ExpectedStoragePledge: soci.ExpectedStoragePledge,
			PowerBaseEpoch:        soci.PowerBaseEpoch,
			ReplacedDayReward:     soci.ReplacedDayReward,
			SectorKeyCID:          soci.SectorKeyCID,
			Flags:                 miner16.SectorOnChainInfoFlags(soci.Flags),
		})
	})

	outSectorsRoot, err := outSectors.Root()
	if err != nil {
		return nil, xerrors.Errorf("failed to flush sectors: %w", err)
	}

	// TODO: implement cached migrator with diff, see v13 for example

	outState := miner16.State{
		Info:                       inState.Info,
		PreCommitDeposits:          inState.PreCommitDeposits,
		LockedFunds:                inState.LockedFunds,
		VestingFunds:               inState.VestingFunds,
		FeeDebt:                    inState.FeeDebt,
		InitialPledge:              inState.InitialPledge,
		PreCommittedSectors:        inState.PreCommittedSectors,
		PreCommittedSectorsCleanUp: inState.PreCommittedSectorsCleanUp,
		AllocatedSectors:           inState.AllocatedSectors,
		Sectors:                    outSectorsRoot,
		ProvingPeriodStart:         inState.ProvingPeriodStart,
		CurrentDeadline:            inState.CurrentDeadline,
		Deadlines:                  inState.Deadlines,
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

func (m *minerMigrator) MigratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m *minerMigrator) Deferred() bool {
	return false
}

var _ migration.ActorMigration = (*minerMigrator)(nil)
