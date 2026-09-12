package migration

import (
	"context"

	power18 "github.com/filecoin-project/go-state-types/builtin/v18/power"
	power19 "github.com/filecoin-project/go-state-types/builtin/v19/power"
	smoothing19 "github.com/filecoin-project/go-state-types/builtin/v19/util/smoothing"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"
)

// powerMigrator drops RampStartEpoch and RampDurationEpochs. The FIP-0081 pledge ramp has
// completed on every network, so the pledge formula is permanently 70% baseline + 30% simple.
type powerMigrator struct {
	OutCodeCID cid.Cid
}

func (p powerMigrator) MigratedCodeCID() cid.Cid {
	return p.OutCodeCID
}

func (p powerMigrator) Deferred() bool {
	return false
}

func (p powerMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (*migration.ActorMigrationResult, error) {
	var inState power18.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load power state for %s: %w", in.Address, err)
	}

	outState := power19.State{
		TotalRawBytePower:         inState.TotalRawBytePower,
		TotalBytesCommitted:       inState.TotalBytesCommitted,
		TotalQualityAdjPower:      inState.TotalQualityAdjPower,
		TotalQABytesCommitted:     inState.TotalQABytesCommitted,
		TotalPledgeCollateral:     inState.TotalPledgeCollateral,
		ThisEpochRawBytePower:     inState.ThisEpochRawBytePower,
		ThisEpochQualityAdjPower:  inState.ThisEpochQualityAdjPower,
		ThisEpochPledgeCollateral: inState.ThisEpochPledgeCollateral,
		ThisEpochQAPowerSmoothed: smoothing19.FilterEstimate{
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
	return &migration.ActorMigrationResult{NewCodeCID: p.OutCodeCID, NewHead: newHead}, nil
}
