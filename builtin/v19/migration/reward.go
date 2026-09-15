package migration

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	reward18 "github.com/filecoin-project/go-state-types/builtin/v18/reward"
	reward19 "github.com/filecoin-project/go-state-types/builtin/v19/reward"
	smoothing19 "github.com/filecoin-project/go-state-types/builtin/v19/util/smoothing"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"
)

type rewardMigrator struct {
	outCodeCID        cid.Cid
	streams           *reward19.StreamsState
	accruals          []reward19.StreamAccrual
	swaTimelockEpochs abi.ChainEpoch
	swaActor          address.Address
}

// ValidateRewardMigrationConfig checks the bootstrap parameters at activationEpoch.
// MigrateStateTree additionally checks recipient existence and actor types in the input tree.
func ValidateRewardMigrationConfig(config RewardMigrationConfig, activationEpoch abi.ChainEpoch) error {
	_, _, err := validateRewardMigrationConfig(config, activationEpoch)
	return err
}

func validateRewardMigrationConfig(config RewardMigrationConfig, activationEpoch abi.ChainEpoch) (*reward19.StreamsState, []reward19.StreamAccrual, error) {
	streamParams := make([]reward19.RegisterStreamParams, len(config.Streams))
	for i, stream := range config.Streams {
		streamParams[i] = reward19.RegisterStreamParams{
			ID: stream.ID,
			Weight: reward19.WeightRecord{
				VStart: stream.Weight.VStart,
				Slope:  stream.Weight.Slope,
				TStart: activationEpoch,
				Floor:  stream.Weight.Floor,
				Cap:    stream.Weight.Cap,
			},
			Distribution:    stream.Distribution,
			ActivationEpoch: activationEpoch,
		}
	}
	streams, accruals, err := reward19.ValidateMigrationStreams(streamParams, activationEpoch)
	if err != nil {
		return nil, nil, xerrors.Errorf("invalid reward migration streams: %w", err)
	}
	if config.SWATimelockEpochs < 0 {
		return nil, nil, xerrors.Errorf("SWA timelock is negative")
	}
	if config.SWAActor.Protocol() != address.ID {
		return nil, nil, xerrors.Errorf("SWA actor is not an ID address")
	}
	return streams, accruals, nil
}

func newRewardMigrator(config RewardMigrationConfig, activationEpoch abi.ChainEpoch, outCodeCID cid.Cid) (*rewardMigrator, error) {
	streams, accruals, err := validateRewardMigrationConfig(config, activationEpoch)
	if err != nil {
		return nil, err
	}
	return &rewardMigrator{
		outCodeCID:        outCodeCID,
		streams:           streams,
		accruals:          accruals,
		swaTimelockEpochs: config.SWATimelockEpochs,
		swaActor:          config.SWAActor,
	}, nil
}

func (m rewardMigrator) validateRecipients(actors *builtin.ActorTree, paychCode cid.Cid) error {
	for _, stream := range m.streams.Streams {
		if stream.Distribution == nil {
			continue
		}
		if !paychCode.Defined() {
			return xerrors.Errorf("code cid for payment channel actor not found in old manifest")
		}
		for _, share := range stream.Distribution.Shares {
			actor, found, err := actors.GetActorV5(share.Recipient)
			if err != nil {
				return xerrors.Errorf("failed to load reward recipient %s: %w", share.Recipient, err)
			}
			if !found {
				return xerrors.Errorf("reward recipient %s does not exist", share.Recipient)
			}
			if actor.Code == paychCode {
				return xerrors.Errorf("reward recipient %s is a payment channel", share.Recipient)
			}
		}
	}
	return nil
}

func (m rewardMigrator) MigratedCodeCID() cid.Cid {
	return m.outCodeCID
}

func (m rewardMigrator) Deferred() bool {
	return false
}

func (m rewardMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (*migration.ActorMigrationResult, error) {
	var inState reward18.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load reward state for %s: %w", in.Address, err)
	}

	streamsRoot, err := store.Put(ctx, m.streams)
	if err != nil {
		return nil, xerrors.Errorf("failed to put reward streams state: %w", err)
	}
	outState := reward19.State{
		CumsumBaseline:          inState.CumsumBaseline,
		CumsumRealized:          inState.CumsumRealized,
		EffectiveNetworkTime:    inState.EffectiveNetworkTime,
		EffectiveBaselinePower:  inState.EffectiveBaselinePower,
		ThisEpochReward:         inState.ThisEpochReward,
		ThisEpochRewardSmoothed: smoothing19.FilterEstimate{PositionEstimate: inState.ThisEpochRewardSmoothed.PositionEstimate, VelocityEstimate: inState.ThisEpochRewardSmoothed.VelocityEstimate},
		ThisEpochBaselinePower:  inState.ThisEpochBaselinePower,
		Epoch:                   inState.Epoch,
		TotalMintedReward:       inState.TotalStoragePowerReward,
		TotalBurnMinted:         big.Zero(),
		TotalExplicitMinted:     big.Zero(),
		Accrued:                 append([]reward19.StreamAccrual(nil), m.accruals...),
		SWATimelockEpochs:       m.swaTimelockEpochs,
		SWAActor:                m.swaActor,
		StreamsRoot:             streamsRoot,
	}
	newHead, err := store.Put(ctx, &outState)
	if err != nil {
		return nil, xerrors.Errorf("failed to put new reward state: %w", err)
	}
	return &migration.ActorMigrationResult{NewCodeCID: m.outCodeCID, NewHead: newHead}, nil
}
