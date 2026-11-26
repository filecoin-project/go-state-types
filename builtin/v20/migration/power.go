package migration

import (
	"context"

	power19 "github.com/filecoin-project/go-state-types/builtin/v19/power"
	power20 "github.com/filecoin-project/go-state-types/builtin/v20/power"
	smoothing20 "github.com/filecoin-project/go-state-types/builtin/v20/util/smoothing"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"
)

type powerMigrator struct {
	OutCodeCID cid.Cid
}

func (m powerMigrator) MigratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m powerMigrator) Deferred() bool {
	return false
}

func (m powerMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (*migration.ActorMigrationResult, error) {
	var inState power19.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load power state for %s: %w", in.Address, err)
	}

	// Removing PowerRamp parameters from FIP81 as ramp up is over
	// RampStartEpoch int64
	// RampDurationEpochs uint64
	outState := power20.State{
		TotalRawBytePower:         inState.TotalRawBytePower,
		TotalBytesCommitted:       inState.TotalBytesCommitted,
		TotalQualityAdjPower:      inState.TotalQualityAdjPower,
		TotalQABytesCommitted:     inState.TotalQABytesCommitted,
		TotalPledgeCollateral:     inState.TotalPledgeCollateral,
		ThisEpochRawBytePower:     inState.ThisEpochRawBytePower,
		ThisEpochQualityAdjPower:  inState.ThisEpochQualityAdjPower,
		ThisEpochPledgeCollateral: inState.ThisEpochPledgeCollateral,
		ThisEpochQAPowerSmoothed: smoothing20.FilterEstimate{
			PositionEstimate: inState.ThisEpochQAPowerSmoothed.PositionEstimate,
			VelocityEstimate: inState.ThisEpochQAPowerSmoothed.VelocityEstimate,
		},
		MinerCount:              inState.MinerCount,
		MinerAboveMinPowerCount: inState.MinerAboveMinPowerCount,
		CronEventQueue:          inState.CronEventQueue,
		FirstCronEpoch:          inState.FirstCronEpoch,
		Claims:                  inState.Claims,
		ProofValidationBatch:    inState.ProofValidationBatch,
	}

	newHead, err := store.Put(ctx, &outState)
	if err != nil {
		return nil, xerrors.Errorf("failed to put new power state: %w", err)
	}
	return &migration.ActorMigrationResult{NewCodeCID: m.OutCodeCID, NewHead: newHead}, nil
}
