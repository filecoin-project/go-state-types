package migration

import (
	"context"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	xerrors "golang.org/x/xerrors"

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

	// Create the new state (v16) with VestingFunds set to nil.
	outState := miner16.State{
		Info:                       inState.Info,
		PreCommitDeposits:          inState.PreCommitDeposits,
		LockedFunds:                inState.LockedFunds,
		VestingFunds:               nil,
		FeeDebt:                    inState.FeeDebt,
		InitialPledge:              inState.InitialPledge,
		PreCommittedSectors:        inState.PreCommittedSectors,
		PreCommittedSectorsCleanUp: inState.PreCommittedSectorsCleanUp,
		AllocatedSectors:           inState.AllocatedSectors,
		Sectors:                    inState.Sectors,
		ProvingPeriodStart:         inState.ProvingPeriodStart,
		CurrentDeadline:            inState.CurrentDeadline,
		Deadlines:                  inState.Deadlines,
		EarlyTerminations:          inState.EarlyTerminations,
		DeadlineCronActive:         inState.DeadlineCronActive,
	}

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
