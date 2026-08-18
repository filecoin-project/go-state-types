package migration

import (
	"context"
	"testing"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	reward18 "github.com/filecoin-project/go-state-types/builtin/v18/reward"
	smoothing18 "github.com/filecoin-project/go-state-types/builtin/v18/util/smoothing"
	reward19 "github.com/filecoin-project/go-state-types/builtin/v19/reward"
	adt19 "github.com/filecoin-project/go-state-types/builtin/v19/util/adt"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"
)

func migrationIDAddress(t *testing.T, id uint64) address.Address {
	t.Helper()
	addr, err := address.NewIDAddress(id)
	require.NoError(t, err)
	return addr
}

func validRewardMigrationConfig(t *testing.T) RewardMigrationConfig {
	t.Helper()
	pct := reward19.Denom / 100
	return RewardMigrationConfig{
		SWATimelockEpochs: 20_160,
		SWAActor:          migrationIDAddress(t, 100),
		Streams: []RewardMigrationStream{
			{
				ID:     1,
				Weight: RewardMigrationWeight{VStart: 95 * pct, Slope: -1, Floor: 50 * pct, Cap: 95 * pct},
			},
			{
				ID:     2,
				Weight: RewardMigrationWeight{VStart: 5 * pct, Slope: 1, Floor: 5 * pct, Cap: 10 * pct},
				Distribution: &reward19.DistributionInit{
					Writer: migrationIDAddress(t, 101),
					Shares: []reward19.RecipientShare{{Recipient: migrationIDAddress(t, 102), Share: reward19.Denom}},
				},
			},
		},
	}
}

func TestRewardMigration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := require.New(t)
	store := cbor.NewMemCborStore()

	inState := reward18.State{
		CumsumBaseline:          big.NewInt(1),
		CumsumRealized:          big.NewInt(2),
		EffectiveNetworkTime:    3,
		EffectiveBaselinePower:  big.NewInt(4),
		ThisEpochReward:         abi.NewTokenAmount(5),
		ThisEpochRewardSmoothed: smoothing18.NewEstimate(big.NewInt(6), big.NewInt(7)),
		ThisEpochBaselinePower:  big.NewInt(8),
		Epoch:                   9,
		TotalStoragePowerReward: abi.NewTokenAmount(10),
		SimpleTotal:             abi.NewTokenAmount(11),
		BaselineTotal:           abi.NewTokenAmount(12),
	}
	inHead, err := store.Put(ctx, &inState)
	req.NoError(err)

	activationEpoch := abi.ChainEpoch(100)
	config := validRewardMigrationConfig(t)
	outCodeCID := cid.MustParse("bafy2bzaca4aaaaaaaaaqk")
	migrator, err := newRewardMigrator(config, activationEpoch, outCodeCID)
	req.NoError(err)
	result, err := migrator.MigrateState(ctx, store, migration.ActorMigrationInput{Address: address.TestAddress, Head: inHead})
	req.NoError(err)
	req.Equal(outCodeCID, result.NewCodeCID)

	var outState reward19.State
	req.NoError(store.Get(ctx, result.NewHead, &outState))
	req.Equal(inState.CumsumBaseline, outState.CumsumBaseline)
	req.Equal(inState.CumsumRealized, outState.CumsumRealized)
	req.Equal(inState.EffectiveNetworkTime, outState.EffectiveNetworkTime)
	req.Equal(inState.EffectiveBaselinePower, outState.EffectiveBaselinePower)
	req.Equal(inState.ThisEpochReward, outState.ThisEpochReward)
	req.Equal(inState.ThisEpochRewardSmoothed.PositionEstimate, outState.ThisEpochRewardSmoothed.PositionEstimate)
	req.Equal(inState.ThisEpochRewardSmoothed.VelocityEstimate, outState.ThisEpochRewardSmoothed.VelocityEstimate)
	req.Equal(inState.ThisEpochBaselinePower, outState.ThisEpochBaselinePower)
	req.Equal(inState.Epoch, outState.Epoch)
	req.Equal(inState.TotalStoragePowerReward, outState.TotalMintedReward)
	req.Equal(big.Zero(), outState.TotalBurnMinted)
	req.Equal(big.Zero(), outState.TotalExplicitMinted)
	req.Equal([]reward19.StreamAccrual{{ID: 2, Amount: big.Zero()}}, outState.Accrued)
	req.Equal(config.SWATimelockEpochs, outState.SWATimelockEpochs)
	req.Equal(config.SWAActor, outState.SWAActor)

	streams, err := outState.LoadStreams(adt19.WrapStore(ctx, store))
	req.NoError(err)
	req.Len(streams.Streams, 2)
	req.Equal(reward19.StreamID(1), streams.Streams[0].ID)
	req.Equal(reward19.WeightRecord{
		VStart: config.Streams[0].Weight.VStart,
		Slope:  config.Streams[0].Weight.Slope,
		TStart: activationEpoch,
		Floor:  config.Streams[0].Weight.Floor,
		Cap:    config.Streams[0].Weight.Cap,
	}, streams.Streams[0].Weight)
	req.Nil(streams.Streams[0].Distribution)
	req.Equal(reward19.StreamID(2), streams.Streams[1].ID)
	req.Equal(reward19.WeightRecord{
		VStart: config.Streams[1].Weight.VStart,
		Slope:  config.Streams[1].Weight.Slope,
		TStart: activationEpoch,
		Floor:  config.Streams[1].Weight.Floor,
		Cap:    config.Streams[1].Weight.Cap,
	}, streams.Streams[1].Weight)
	req.Equal(config.Streams[1].Distribution.Writer, streams.Streams[1].Distribution.Writer)
	req.Equal(config.Streams[1].Distribution.Shares, streams.Streams[1].Distribution.Shares)
	req.Empty(streams.Tombstones)
	req.Empty(streams.PendingWrites)
}

func TestRewardMigrationDropsStoredRewardTotals(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := require.New(t)
	store := cbor.NewMemCborStore()
	activationEpoch := abi.ChainEpoch(100)
	migrator, err := newRewardMigrator(
		validRewardMigrationConfig(t),
		activationEpoch,
		cid.MustParse("bafy2bzaca4aaaaaaaaaqk"),
	)
	req.NoError(err)

	base := reward18.State{
		CumsumBaseline:          big.NewInt(1),
		CumsumRealized:          big.NewInt(2),
		EffectiveNetworkTime:    3,
		EffectiveBaselinePower:  big.NewInt(4),
		ThisEpochReward:         abi.NewTokenAmount(5),
		ThisEpochRewardSmoothed: smoothing18.NewEstimate(big.NewInt(6), big.NewInt(7)),
		ThisEpochBaselinePower:  big.NewInt(8),
		Epoch:                   9,
		TotalStoragePowerReward: abi.NewTokenAmount(10),
		SimpleTotal:             abi.NewTokenAmount(11),
		BaselineTotal:           abi.NewTokenAmount(12),
	}
	alternate := base
	alternate.SimpleTotal = abi.NewTokenAmount(13)
	alternate.BaselineTotal = abi.NewTokenAmount(14)
	req.NotEqual(base.SimpleTotal, alternate.SimpleTotal)
	req.NotEqual(base.BaselineTotal, alternate.BaselineTotal)

	migrate := func(input reward18.State) reward19.State {
		head, err := store.Put(ctx, &input)
		req.NoError(err)
		result, err := migrator.MigrateState(ctx, store, migration.ActorMigrationInput{
			Address: address.TestAddress,
			Head:    head,
		})
		req.NoError(err)
		var output reward19.State
		req.NoError(store.Get(ctx, result.NewHead, &output))
		return output
	}

	req.Equal(migrate(base), migrate(alternate))
}

func TestNewRewardMigratorAcceptsAlternativeBootstrapWeights(t *testing.T) {
	activationEpoch := abi.ChainEpoch(100)
	pct := reward19.Denom / 100
	config := validRewardMigrationConfig(t)
	config.Streams[0].Weight = RewardMigrationWeight{
		VStart: 80 * pct,
		Slope:  -1,
		Floor:  60 * pct,
		Cap:    80 * pct,
	}
	config.Streams[1].Weight = RewardMigrationWeight{
		VStart: 20 * pct,
		Slope:  1,
		Floor:  10 * pct,
		Cap:    20 * pct,
	}

	_, err := newRewardMigrator(config, activationEpoch, cid.MustParse("bafy2bzaca4aaaaaaaaaqk"))
	require.NoError(t, err)
}

func TestNewRewardMigratorRejectsInvalidConfig(t *testing.T) {
	activationEpoch := abi.ChainEpoch(100)
	outCodeCID := cid.MustParse("bafy2bzaca4aaaaaaaaaqk")
	testCases := []struct {
		name     string
		mutate   func(*RewardMigrationConfig)
		expected string
	}{
		{
			name: "wrong stream count",
			mutate: func(config *RewardMigrationConfig) {
				config.Streams = config.Streams[:1]
			},
			expected: "requires exactly two streams",
		},
		{
			name: "starting weights under-sum",
			mutate: func(config *RewardMigrationConfig) {
				config.Streams[1].Weight.VStart--
			},
			expected: "bootstrap starting weights must sum to denominator",
		},
		{
			name: "unequal slopes",
			mutate: func(config *RewardMigrationConfig) {
				config.Streams[1].Weight.Slope++
			},
			expected: "bootstrap weight slopes are invalid",
		},
		{
			name: "invalid initial distribution",
			mutate: func(config *RewardMigrationConfig) {
				config.Streams[1].Distribution.Shares[0].Share--
			},
			expected: "one full-share recipient",
		},
		{
			name: "negative timelock",
			mutate: func(config *RewardMigrationConfig) {
				config.SWATimelockEpochs = -1
			},
			expected: "SWA timelock is negative",
		},
		{
			name: "non-ID SWA actor",
			mutate: func(config *RewardMigrationConfig) {
				addr, err := address.NewDelegatedAddress(10, []byte{1})
				require.NoError(t, err)
				config.SWAActor = addr
			},
			expected: "SWA actor is not an ID address",
		},
		{
			name: "non-ID SRA actor",
			mutate: func(config *RewardMigrationConfig) {
				addr, err := address.NewDelegatedAddress(10, []byte{1})
				require.NoError(t, err)
				config.Streams[1].Distribution.Writer = addr
			},
			expected: "distribution writer",
		},
		{
			name: "non-ID initial orchestrator",
			mutate: func(config *RewardMigrationConfig) {
				addr, err := address.NewDelegatedAddress(10, []byte{1})
				require.NoError(t, err)
				config.Streams[1].Distribution.Shares[0].Recipient = addr
			},
			expected: "share recipient",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := validRewardMigrationConfig(t)
			tc.mutate(&config)
			_, err := newRewardMigrator(config, activationEpoch, outCodeCID)
			require.ErrorContains(t, err, tc.expected)
		})
	}
}
