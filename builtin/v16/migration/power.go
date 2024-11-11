package migration

import (
	"context"

	power14 "github.com/filecoin-project/go-state-types/builtin/v14/power"
	power15 "github.com/filecoin-project/go-state-types/builtin/v16/power"
	smoothing15 "github.com/filecoin-project/go-state-types/builtin/v16/util/smoothing"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"
)

type powerMigrator struct {
	rampStartEpoch     int64
	rampDurationEpochs uint64
	outCodeCID         cid.Cid
}

func newPowerMigrator(rampStartEpoch int64, rampDurationEpochs uint64, outCode cid.Cid) (*powerMigrator, error) {
	return &powerMigrator{
		rampStartEpoch:     rampStartEpoch,
		rampDurationEpochs: rampDurationEpochs,
		outCodeCID:         outCode,
	}, nil
}

func (p powerMigrator) MigratedCodeCID() cid.Cid { return p.outCodeCID }

func (p powerMigrator) Deferred() bool { return false }

func (p powerMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (*migration.ActorMigrationResult, error) {
	var inState power14.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load miner state for %s: %w", in.Address, err)
	}

	outState := power15.State{
		TotalRawBytePower:         inState.TotalRawBytePower,
		TotalBytesCommitted:       inState.TotalBytesCommitted,
		TotalQualityAdjPower:      inState.TotalQualityAdjPower,
		TotalQABytesCommitted:     inState.TotalQABytesCommitted,
		TotalPledgeCollateral:     inState.TotalPledgeCollateral,
		ThisEpochRawBytePower:     inState.ThisEpochRawBytePower,
		ThisEpochQualityAdjPower:  inState.ThisEpochQualityAdjPower,
		ThisEpochPledgeCollateral: inState.ThisEpochPledgeCollateral,
		ThisEpochQAPowerSmoothed: smoothing15.FilterEstimate{
			PositionEstimate: inState.ThisEpochQAPowerSmoothed.PositionEstimate,
			VelocityEstimate: inState.ThisEpochQAPowerSmoothed.VelocityEstimate,
		},
		MinerCount:              inState.MinerCount,
		MinerAboveMinPowerCount: inState.MinerAboveMinPowerCount,
		RampStartEpoch:          p.rampStartEpoch,
		RampDurationEpochs:      p.rampDurationEpochs,
		CronEventQueue:          inState.CronEventQueue,
		FirstCronEpoch:          inState.FirstCronEpoch,
		Claims:                  inState.Claims,
		ProofValidationBatch:    inState.ProofValidationBatch,
	}

	newHead, err := store.Put(ctx, &outState)
	if err != nil {
		return nil, xerrors.Errorf("failed to put new state: %w", err)
	}

	return &migration.ActorMigrationResult{
		NewCodeCID: p.MigratedCodeCID(),
		NewHead:    newHead,
	}, nil
}
