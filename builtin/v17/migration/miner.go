package migration

import (
	"context"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	xerrors "golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/abi"
	miner15 "github.com/filecoin-project/go-state-types/builtin/v15/miner"
	miner16 "github.com/filecoin-project/go-state-types/builtin/v16/miner"
	adt16 "github.com/filecoin-project/go-state-types/builtin/v16/util/adt"
	"github.com/filecoin-project/go-state-types/migration"
)

// minerMigrator performs the migration for the Miner contract state.
type minerMigrator struct {
	OutCodeCID cid.Cid
}

func (m minerMigrator) MigratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m minerMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (*migration.ActorMigrationResult, error) {
	// Load the existing state (v15).
	var inState miner15.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load miner state for %s: %w", in.Address, err)
	}

	adtStore := adt16.WrapStore(ctx, store)

	// Create the new state (v16) with the same values as the old state but VestedFunds and Deadlines
	// set to zero values (to be filled in below).
	outState := miner16.State{
		Info:                       inState.Info,
		PreCommitDeposits:          inState.PreCommitDeposits,
		LockedFunds:                inState.LockedFunds,
		VestingFunds:               nil, // set below
		FeeDebt:                    inState.FeeDebt,
		InitialPledge:              inState.InitialPledge,
		PreCommittedSectors:        inState.PreCommittedSectors,
		PreCommittedSectorsCleanUp: inState.PreCommittedSectorsCleanUp,
		AllocatedSectors:           inState.AllocatedSectors,
		Sectors:                    inState.Sectors,
		ProvingPeriodStart:         inState.ProvingPeriodStart,
		CurrentDeadline:            inState.CurrentDeadline,
		Deadlines:                  cid.Undef, // set below
		EarlyTerminations:          inState.EarlyTerminations,
		DeadlineCronActive:         inState.DeadlineCronActive,
	}

	// Copy over the vesting funds into the new structure - head (first VestingFund) in the miner
	// state, tail (the rest) in a separate block.

	oldVestingFunds, err := inState.LoadVestingFunds(adtStore)
	if err != nil {
		return nil, xerrors.Errorf("failed to load vesting funds for %s: %w", in.Address, err)
	}

	newVestingFunds := make([]miner16.VestingFund, len(oldVestingFunds))
	for i, oldVestingFund := range oldVestingFunds {
		newVestingFunds[i] = miner16.VestingFund{
			Epoch:  oldVestingFund.Epoch,
			Amount: oldVestingFund.Amount,
		}
	}

	if len(newVestingFunds) > 0 {
		head := newVestingFunds[0]
		tail := newVestingFunds[1:]
		tailCid, err := store.Put(ctx, &miner16.VestingFundsTail{Funds: tail})
		if err != nil {
			return nil, xerrors.Errorf("failed to persist vesting funds tail for %s: %w", in.Address, err)
		}
		outState.VestingFunds = &miner16.VestingFunds{
			Head: head,
			Tail: tailCid,
		}
	}

	// Copy over the deadlines into the new structure, but add the new LivePower and DailyFee fields.
	// LivePower can be calculated from the partitions, and DailyFee starts at zero for existing
	// miners.

	deadlines, err := inState.LoadDeadlines(adtStore)
	if err != nil {
		return nil, xerrors.Errorf("failed to load deadlines for %s: %w", in.Address, err)
	}
	var newDeadlines miner16.Deadlines
	for dlIdx := range deadlines.Due {
		oldDeadline, err := deadlines.LoadDeadline(adtStore, uint64(dlIdx))
		if err != nil {
			return nil, xerrors.Errorf("failed to load deadline %d for %s: %w", dlIdx, in.Address, err)
		}

		// Copy old deadline
		newDeadline := miner16.Deadline{
			Partitions:                        oldDeadline.Partitions,
			ExpirationsEpochs:                 oldDeadline.ExpirationsEpochs,
			PartitionsPoSted:                  oldDeadline.PartitionsPoSted,
			EarlyTerminations:                 oldDeadline.EarlyTerminations,
			LiveSectors:                       oldDeadline.LiveSectors,
			TotalSectors:                      oldDeadline.TotalSectors,
			FaultyPower:                       miner16.NewPowerPair(oldDeadline.FaultyPower.Raw, oldDeadline.FaultyPower.QA),
			OptimisticPoStSubmissions:         oldDeadline.OptimisticPoStSubmissions,
			SectorsSnapshot:                   oldDeadline.SectorsSnapshot,
			PartitionsSnapshot:                oldDeadline.PartitionsSnapshot,
			OptimisticPoStSubmissionsSnapshot: oldDeadline.OptimisticPoStSubmissionsSnapshot,
			LivePower:                         miner16.NewPowerPairZero(),
			DailyFee:                          abi.NewTokenAmount(0),
		}

		// Sum up live power in the partitions of this deadline
		partitions, err := oldDeadline.PartitionsArray(adtStore)
		if err != nil {
			return nil, xerrors.Errorf("failed to load partitions for deadline %d for %s: %w", dlIdx, in.Address, err)
		}
		var partition miner16.Partition
		err = partitions.ForEach(&partition, func(partIdx int64) error {
			newDeadline.LivePower = newDeadline.LivePower.Add(partition.LivePower)
			return nil
		})
		if err != nil {
			return nil, xerrors.Errorf("failed to sum up live power for deadline %d for %s: %w", dlIdx, in.Address, err)
		}

		// Store the new deadline
		dlcid, err := store.Put(ctx, &newDeadline)
		if err != nil {
			return nil, xerrors.Errorf("failed to persist deadline %d for %s: %w", dlIdx, in.Address, err)
		}
		newDeadlines.Due[dlIdx] = dlcid
	}

	// Store the new deadlines array
	dlscid, err := store.Put(ctx, &newDeadlines)
	if err != nil {
		return nil, xerrors.Errorf("failed to persist deadlines for %s: %w", in.Address, err)
	}
	outState.Deadlines = dlscid

	// Store the new state.
	newHead, err := store.Put(ctx, &outState)
	if err != nil {
		return nil, xerrors.Errorf("failed to put new evm state: %w", err)
	}

	return &migration.ActorMigrationResult{
		NewCodeCID: m.MigratedCodeCID(),
		NewHead:    newHead,
	}, nil
}

func (m minerMigrator) Deferred() bool {
	return false
}
