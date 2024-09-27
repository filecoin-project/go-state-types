package migration

import (
	"context"
	"testing"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	power14 "github.com/filecoin-project/go-state-types/builtin/v14/power"
	smoothing14 "github.com/filecoin-project/go-state-types/builtin/v14/util/smoothing"
	power15 "github.com/filecoin-project/go-state-types/builtin/v15/power"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"
)

func TestPowerMigration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req := require.New(t)

	cst := cbor.NewMemCborStore()

	pvb := cid.MustParse("bafy2bzacaf2a")
	state14 := power14.State{
		TotalRawBytePower:         abi.NewStoragePower(101),
		TotalBytesCommitted:       abi.NewStoragePower(102),
		TotalQualityAdjPower:      abi.NewStoragePower(103),
		TotalQABytesCommitted:     abi.NewStoragePower(104),
		TotalPledgeCollateral:     abi.NewTokenAmount(105),
		ThisEpochRawBytePower:     abi.NewStoragePower(106),
		ThisEpochQualityAdjPower:  abi.NewStoragePower(107),
		ThisEpochPledgeCollateral: abi.NewTokenAmount(108),
		ThisEpochQAPowerSmoothed:  smoothing14.NewEstimate(big.NewInt(109), big.NewInt(110)),
		MinerCount:                111,
		MinerAboveMinPowerCount:   112,
		CronEventQueue:            cid.MustParse("bafy2bzacafza"),
		FirstCronEpoch:            113,
		Claims:                    cid.MustParse("bafy2bzacafzq"),
		ProofValidationBatch:      &pvb,
	}

	state14Cid, err := cst.Put(ctx, &state14)
	req.NoError(err)

	var rampStartEpoch int64 = 101
	var rampDurationEpochs uint64 = 202

	migrator, err := newPowerMigrator(rampStartEpoch, rampDurationEpochs, cid.MustParse("bafy2bzaca4aaaaaaaaaqk"))
	req.NoError(err)

	result, err := migrator.MigrateState(ctx, cst, migration.ActorMigrationInput{
		Address: address.TestAddress,
		Head:    state14Cid,
		Cache:   nil,
	})
	req.NoError(err)
	req.Equal(cid.MustParse("bafy2bzaca4aaaaaaaaaqk"), result.NewCodeCID)
	req.NotEqual(cid.Undef, result.NewHead)

	newState := power15.State{}
	req.NoError(cst.Get(ctx, result.NewHead, &newState))

	req.Equal(state14.TotalRawBytePower, newState.TotalRawBytePower)
	req.Equal(state14.TotalBytesCommitted, newState.TotalBytesCommitted)
	req.Equal(state14.TotalQualityAdjPower, newState.TotalQualityAdjPower)
	req.Equal(state14.TotalQABytesCommitted, newState.TotalQABytesCommitted)
	req.Equal(state14.TotalPledgeCollateral, newState.TotalPledgeCollateral)
	req.Equal(state14.ThisEpochRawBytePower, newState.ThisEpochRawBytePower)
	req.Equal(state14.ThisEpochQualityAdjPower, newState.ThisEpochQualityAdjPower)
	req.Equal(state14.ThisEpochPledgeCollateral, newState.ThisEpochPledgeCollateral)
	req.Equal(state14.ThisEpochQAPowerSmoothed.PositionEstimate, newState.ThisEpochQAPowerSmoothed.PositionEstimate)
	req.Equal(state14.ThisEpochQAPowerSmoothed.VelocityEstimate, newState.ThisEpochQAPowerSmoothed.VelocityEstimate)
	req.Equal(state14.MinerCount, newState.MinerCount)
	req.Equal(state14.MinerAboveMinPowerCount, newState.MinerAboveMinPowerCount)
	req.Equal(state14.CronEventQueue, newState.CronEventQueue)
	req.Equal(state14.FirstCronEpoch, newState.FirstCronEpoch)
	req.Equal(state14.Claims, newState.Claims)
	req.Equal(state14.ProofValidationBatch, newState.ProofValidationBatch)

	// Ramp parameters
	req.Equal(rampStartEpoch, newState.RampStartEpoch)
	req.Equal(rampDurationEpochs, newState.RampDurationEpochs)
}
