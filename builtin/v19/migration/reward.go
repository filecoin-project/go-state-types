package migration

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
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

func newRewardMigrator(config RewardMigrationConfig, activationEpoch abi.ChainEpoch, outCodeCID cid.Cid) (*rewardMigrator, error) {
	streams, accruals, err := reward19.ValidateMigrationStreams(config.Streams, activationEpoch)
	if err != nil {
		return nil, xerrors.Errorf("invalid reward migration streams: %w", err)
	}
	if config.SWATimelockEpochs < 0 {
		return nil, xerrors.Errorf("SWA timelock is negative")
	}
	if config.SWAActor.Protocol() != address.ID {
		return nil, xerrors.Errorf("SWA actor is not an ID address")
	}
	return &rewardMigrator{
		outCodeCID:        outCodeCID,
		streams:           streams,
		accruals:          accruals,
		swaTimelockEpochs: config.SWATimelockEpochs,
		swaActor:          config.SWAActor,
	}, nil
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
