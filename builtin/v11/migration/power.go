package migration

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v11/util/smoothing"

	"github.com/filecoin-project/go-address"
	power11 "github.com/filecoin-project/go-state-types/builtin/v11/power"

	adt10 "github.com/filecoin-project/go-state-types/builtin/v10/util/adt"

	"github.com/filecoin-project/go-state-types/builtin"
	power10 "github.com/filecoin-project/go-state-types/builtin/v10/power"
	adt11 "github.com/filecoin-project/go-state-types/builtin/v11/util/adt"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"
)

// The powerMigrator performs the following migrations:
// - FIP-0061: Updates all claims in the power table to based on V1_1 proof types
type powerMigrator struct {
	OutCodeCID cid.Cid
}

func (m powerMigrator) MigratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m powerMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (*migration.ActorMigrationResult, error) {
	var inState power10.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load power state: %w", err)
	}

	ctxStore := adt11.WrapStore(ctx, store)
	inClaims, err := adt10.AsMap(ctxStore, inState.Claims, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to load claims: %w", err)
	}

	emptyClaimsMapCid, err := adt11.StoreEmptyMap(ctxStore, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct empty claims: %w", err)
	}

	outClaims, err := adt11.AsMap(ctxStore, emptyClaimsMapCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to load empty claims: %w", err)
	}

	var inClaim power10.Claim
	if err := inClaims.ForEach(&inClaim, func(key string) error {
		addr, err := address.NewFromBytes([]byte(key))
		if err != nil {
			return xerrors.Errorf("failed to get addr: %w", err)
		}

		newProofType, err := inClaim.WindowPoStProofType.ToV1_1PostProof()
		if err != nil {
			return xerrors.Errorf("failed to convert proof type %d for miner %s: %w", inClaim.WindowPoStProofType, addr, err)
		}

		outClaim := power11.Claim{
			WindowPoStProofType: newProofType,
			RawBytePower:        inClaim.RawBytePower,
			QualityAdjPower:     inClaim.QualityAdjPower,
		}

		if err := outClaims.Put(abi.AddrKey(addr), &outClaim); err != nil {
			return xerrors.Errorf("failed to put new claim for miner %s: %w", addr, err)
		}

		return nil
	}); err != nil {
		return nil, xerrors.Errorf("failed to iterate over inClaims: %w", err)
	}

	outClaimsRoot, err := outClaims.Root()
	if err != nil {
		return nil, xerrors.Errorf("failed to flush outClaims: %w", err)
	}

	outState := power11.State{
		TotalRawBytePower:         inState.TotalRawBytePower,
		TotalBytesCommitted:       inState.TotalBytesCommitted,
		TotalQualityAdjPower:      inState.TotalQualityAdjPower,
		TotalQABytesCommitted:     inState.TotalQABytesCommitted,
		TotalPledgeCollateral:     inState.TotalPledgeCollateral,
		ThisEpochRawBytePower:     inState.ThisEpochRawBytePower,
		ThisEpochQualityAdjPower:  inState.ThisEpochQualityAdjPower,
		ThisEpochPledgeCollateral: inState.ThisEpochPledgeCollateral,
		ThisEpochQAPowerSmoothed:  smoothing.FilterEstimate(inState.ThisEpochQAPowerSmoothed),
		MinerCount:                inState.MinerCount,
		MinerAboveMinPowerCount:   inState.MinerAboveMinPowerCount,
		CronEventQueue:            inState.CronEventQueue,
		FirstCronEpoch:            inState.FirstCronEpoch,
		Claims:                    outClaimsRoot,
		ProofValidationBatch:      inState.ProofValidationBatch,
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

func (m powerMigrator) Deferred() bool {
	return false
}
