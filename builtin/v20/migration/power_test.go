package migration

import (
	"context"
	"testing"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	power19 "github.com/filecoin-project/go-state-types/builtin/v19/power"
	smoothing19 "github.com/filecoin-project/go-state-types/builtin/v19/util/smoothing"
	power20 "github.com/filecoin-project/go-state-types/builtin/v20/power"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"
)

func TestPowerMigration(t *testing.T) {
	pvb := cid.MustParse("bafy2bzacaf2a")
	for _, tc := range []struct {
		name                 string
		proofValidationBatch *cid.Cid
	}{
		{name: "without proof validation batch"},
		{name: "with proof validation batch", proofValidationBatch: &pvb},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			req := require.New(t)
			store := cbor.NewMemCborStore()

			inState := power19.State{
				TotalRawBytePower:         abi.NewStoragePower(101),
				TotalBytesCommitted:       abi.NewStoragePower(102),
				TotalQualityAdjPower:      abi.NewStoragePower(103),
				TotalQABytesCommitted:     abi.NewStoragePower(104),
				TotalPledgeCollateral:     abi.NewTokenAmount(105),
				ThisEpochRawBytePower:     abi.NewStoragePower(106),
				ThisEpochQualityAdjPower:  abi.NewStoragePower(107),
				ThisEpochPledgeCollateral: abi.NewTokenAmount(108),
				ThisEpochQAPowerSmoothed:  smoothing19.NewEstimate(big.NewInt(109), big.NewInt(110)),
				MinerCount:                112,
				MinerAboveMinPowerCount:   111,
				RampStartEpoch:            113,
				RampDurationEpochs:        114,
				CronEventQueue:            cid.MustParse("bafy2bzacafza"),
				FirstCronEpoch:            115,
				Claims:                    cid.MustParse("bafy2bzacafzq"),
				ProofValidationBatch:      tc.proofValidationBatch,
			}
			inHead, err := store.Put(ctx, &inState)
			req.NoError(err)

			outCodeCID := cid.MustParse("bafy2bzaca4aaaaaaaaaqk")
			migrator := powerMigrator{OutCodeCID: outCodeCID}
			req.Equal(outCodeCID, migrator.MigratedCodeCID())
			req.False(migrator.Deferred())

			result, err := migrator.MigrateState(ctx, store, migration.ActorMigrationInput{
				Address: address.TestAddress,
				Head:    inHead,
			})
			req.NoError(err)
			req.Equal(outCodeCID, result.NewCodeCID)
			req.NotEqual(cid.Undef, result.NewHead)

			var outState power20.State
			req.NoError(store.Get(ctx, result.NewHead, &outState))

			req.Equal(inState.TotalRawBytePower, outState.TotalRawBytePower)
			req.Equal(inState.TotalBytesCommitted, outState.TotalBytesCommitted)
			req.Equal(inState.TotalQualityAdjPower, outState.TotalQualityAdjPower)
			req.Equal(inState.TotalQABytesCommitted, outState.TotalQABytesCommitted)
			req.Equal(inState.TotalPledgeCollateral, outState.TotalPledgeCollateral)
			req.Equal(inState.ThisEpochRawBytePower, outState.ThisEpochRawBytePower)
			req.Equal(inState.ThisEpochQualityAdjPower, outState.ThisEpochQualityAdjPower)
			req.Equal(inState.ThisEpochPledgeCollateral, outState.ThisEpochPledgeCollateral)
			req.Equal(inState.ThisEpochQAPowerSmoothed.PositionEstimate, outState.ThisEpochQAPowerSmoothed.PositionEstimate)
			req.Equal(inState.ThisEpochQAPowerSmoothed.VelocityEstimate, outState.ThisEpochQAPowerSmoothed.VelocityEstimate)
			req.Equal(inState.MinerCount, outState.MinerCount)
			req.Equal(inState.MinerAboveMinPowerCount, outState.MinerAboveMinPowerCount)
			req.Equal(inState.CronEventQueue, outState.CronEventQueue)
			req.Equal(inState.FirstCronEpoch, outState.FirstCronEpoch)
			req.Equal(inState.Claims, outState.Claims)
			req.Equal(inState.ProofValidationBatch, outState.ProofValidationBatch)

			// Removed ramp parameters must not affect the migrated state.
			inState.RampStartEpoch++
			inState.RampDurationEpochs++
			alternateHead, err := store.Put(ctx, &inState)
			req.NoError(err)
			req.NotEqual(inHead, alternateHead)
			alternateResult, err := migrator.MigrateState(ctx, store, migration.ActorMigrationInput{
				Address: address.TestAddress,
				Head:    alternateHead,
			})
			req.NoError(err)
			req.Equal(result.NewHead, alternateResult.NewHead)
		})
	}
}
