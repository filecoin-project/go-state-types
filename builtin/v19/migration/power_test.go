package migration

import (
	"context"
	"testing"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	power18 "github.com/filecoin-project/go-state-types/builtin/v18/power"
	smoothing18 "github.com/filecoin-project/go-state-types/builtin/v18/util/smoothing"
	power19 "github.com/filecoin-project/go-state-types/builtin/v19/power"
	"github.com/filecoin-project/go-state-types/migration"
	cid "github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"
)

// populatedPowerState builds a v18 power state with a distinct value in every field, so a
// field-shift during migration shows up as a wrong value rather than merely a missing one.
func populatedPowerState() power18.State {
	pvb := cid.MustParse("bafy2bzacaf2a")
	return power18.State{
		TotalRawBytePower:         abi.NewStoragePower(101),
		TotalBytesCommitted:       abi.NewStoragePower(102),
		TotalQualityAdjPower:      abi.NewStoragePower(103),
		TotalQABytesCommitted:     abi.NewStoragePower(104),
		TotalPledgeCollateral:     abi.NewTokenAmount(105),
		ThisEpochRawBytePower:     abi.NewStoragePower(106),
		ThisEpochQualityAdjPower:  abi.NewStoragePower(107),
		ThisEpochPledgeCollateral: abi.NewTokenAmount(108),
		ThisEpochQAPowerSmoothed:  smoothing18.NewEstimate(big.NewInt(109), big.NewInt(110)),
		MinerCount:                111,
		MinerAboveMinPowerCount:   112,
		RampStartEpoch:            113,
		RampDurationEpochs:        114,
		CronEventQueue:            cid.MustParse("bafy2bzacafza"),
		FirstCronEpoch:            115,
		Claims:                    cid.MustParse("bafy2bzacafzq"),
		ProofValidationBatch:      &pvb,
	}
}

// Every surviving field must arrive with its own value: the removed fields sit mid-tuple, so
// a mistake shifts CronEventQueue and everything after it.
func TestPowerMigration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := require.New(t)
	store := cbor.NewMemCborStore()

	inState := populatedPowerState()
	inHead, err := store.Put(ctx, &inState)
	req.NoError(err)

	outCodeCID := cid.MustParse("bafy2bzaca4aaaaaaaaaqk")
	migrator := powerMigrator{OutCodeCID: outCodeCID}
	req.Equal(outCodeCID, migrator.MigratedCodeCID())
	req.False(migrator.Deferred())

	result, err := migrator.MigrateState(ctx, store, migration.ActorMigrationInput{Address: address.TestAddress, Head: inHead})
	req.NoError(err)
	req.Equal(outCodeCID, result.NewCodeCID)

	var outState power19.State
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
	// The fields that follow the removed ones: this is what a shift would corrupt.
	req.Equal(inState.CronEventQueue, outState.CronEventQueue)
	req.Equal(inState.FirstCronEpoch, outState.FirstCronEpoch)
	req.Equal(inState.Claims, outState.Claims)
	req.Equal(inState.ProofValidationBatch, outState.ProofValidationBatch)
}

// The v19 state must round-trip through CBOR without the removed fields, and decode to the
// same values it was written with.
func TestPowerStateRoundTripsWithoutRampFields(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := require.New(t)
	store := cbor.NewMemCborStore()

	inState := populatedPowerState()
	inHead, err := store.Put(ctx, &inState)
	req.NoError(err)
	result, err := powerMigrator{OutCodeCID: cid.MustParse("bafy2bzaca4aaaaaaaaaqk")}.
		MigrateState(ctx, store, migration.ActorMigrationInput{Address: address.TestAddress, Head: inHead})
	req.NoError(err)

	var outState power19.State
	req.NoError(store.Get(ctx, result.NewHead, &outState))
	rewritten, err := store.Put(ctx, &outState)
	req.NoError(err)
	req.Equal(result.NewHead, rewritten, "v19 power state does not round-trip")

	// A v18 decoder must not accept it: the tuple is two fields shorter.
	var asV18 power18.State
	req.Error(store.Get(ctx, result.NewHead, &asV18))
}
