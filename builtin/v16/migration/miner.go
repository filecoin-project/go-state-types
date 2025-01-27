package migration

import (
	"context"
	"fmt"
	"sync/atomic"

	miner15 "github.com/filecoin-project/go-state-types/builtin/v15/miner"
	miner16 "github.com/filecoin-project/go-state-types/builtin/v16/miner"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"
)

type minerMigrator struct {
	OutCodeCID                      cid.Cid
	EmptyPreCommittedSectorsHamtCid cid.Cid
	EmptyPrecommitCleanUpAmtCid     cid.Cid
}

func newMinerMigrator(
	_ context.Context,
	_ cbor.IpldStore,
	outCode,
	emptyPreCommittedSectorsHamtCid,
	emptyPrecommitCleanUpAmtCid cid.Cid,
) (*minerMigrator, error) {
	return &minerMigrator{
		OutCodeCID:                      outCode,
		EmptyPreCommittedSectorsHamtCid: emptyPreCommittedSectorsHamtCid,
		EmptyPrecommitCleanUpAmtCid:     emptyPrecommitCleanUpAmtCid,
	}, nil
}

// TOOD: remove these, they're just for debugging
var minerCount, withEmptyPrecommitCleanUpAmtCidCount, withEmptyPreCommittedSectorsHamtCidCount, migratedCount uint64

func (m *minerMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (result *migration.ActorMigrationResult, err error) {
	var inState miner15.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load miner state for %s: %w", in.Address, err)
	}

	newHead := in.Head

	var epccuacc, epcshcc uint64
	var preCommittedSectors, preCommittedSectorsCleanUp *cid.Cid
	if !m.EmptyPrecommitCleanUpAmtCid.Equals(inState.PreCommittedSectorsCleanUp) {
		preCommittedSectorsCleanUp = &inState.PreCommittedSectorsCleanUp
	} else {
		epccuacc = atomic.AddUint64(&withEmptyPrecommitCleanUpAmtCidCount, 1)
	}
	if !m.EmptyPreCommittedSectorsHamtCid.Equals(inState.PreCommittedSectors) {
		preCommittedSectors = &inState.PreCommittedSectors
	} else {
		epcshcc = atomic.AddUint64(&withEmptyPreCommittedSectorsHamtCidCount, 1)
	}

	// we only need to update the state if one of these need to be nullable, otherwise the existing
	// state, with a valid CID, will decode correctly even though they're now pointer fields
	if preCommittedSectors == nil || preCommittedSectorsCleanUp == nil {
		outState := miner16.State{
			Info:                       inState.Info,
			PreCommitDeposits:          inState.PreCommitDeposits,
			LockedFunds:                inState.LockedFunds,
			VestingFunds:               inState.VestingFunds,
			FeeDebt:                    inState.FeeDebt,
			InitialPledge:              inState.InitialPledge,
			PreCommittedSectors:        preCommittedSectors,
			PreCommittedSectorsCleanUp: preCommittedSectorsCleanUp,
			AllocatedSectors:           inState.AllocatedSectors,
			Sectors:                    inState.Sectors,
			ProvingPeriodStart:         inState.ProvingPeriodStart,
			CurrentDeadline:            inState.CurrentDeadline,
			Deadlines:                  inState.Deadlines,
			EarlyTerminations:          inState.EarlyTerminations,
			DeadlineCronActive:         inState.DeadlineCronActive,
		}

		if newHead, err = store.Put(ctx, &outState); err != nil {
			return nil, xerrors.Errorf("failed to put new state: %w", err)
		}

		atomic.AddUint64(&migratedCount, 1)
	}

	if nc := atomic.AddUint64(&minerCount, 1); nc%10000 == 0 {
		fmt.Printf("Checked %d miners, %d with empty PreCommittedSectors, %d with empty PreCommittedSectorsCleanUp, %d migrated\n", nc, epcshcc, epccuacc, atomic.LoadUint64(&migratedCount))
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
