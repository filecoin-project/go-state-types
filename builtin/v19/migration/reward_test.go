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

func validRewardMigrationConfig(t *testing.T, activationEpoch abi.ChainEpoch) RewardMigrationConfig {
	t.Helper()
	pct := reward19.Denom / 100
	return RewardMigrationConfig{
		SWATimelockEpochs: 20_160,
		SWAActor:          migrationIDAddress(t, 100),
		Streams: []reward19.RegisterStreamParams{
			{
				ID:              1,
				Weight:          reward19.WeightRecord{VStart: 95 * pct, Slope: -1, TStart: activationEpoch, Floor: 50 * pct, Cap: 95 * pct},
				ActivationEpoch: activationEpoch,
			},
			{
				ID:     2,
				Weight: reward19.WeightRecord{VStart: 5 * pct, Slope: 1, TStart: activationEpoch, Floor: 5 * pct, Cap: 10 * pct},
				Distribution: &reward19.DistributionInit{
					Writer: migrationIDAddress(t, 101),
					Shares: []reward19.RecipientShare{{Recipient: migrationIDAddress(t, 102), Share: reward19.Denom}},
				},
				ActivationEpoch: activationEpoch,
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
	config := validRewardMigrationConfig(t, activationEpoch)
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

	var streams reward19.StreamsState
	req.NoError(store.Get(ctx, outState.StreamsRoot, &streams))
	req.Len(streams.Streams, 2)
	req.Equal(reward19.StreamID(1), streams.Streams[0].ID)
	req.Nil(streams.Streams[0].Distribution)
	req.Equal(reward19.StreamID(2), streams.Streams[1].ID)
	req.Equal(config.Streams[1].Distribution.Writer, streams.Streams[1].Distribution.Writer)
	req.Equal(config.Streams[1].Distribution.Shares, streams.Streams[1].Distribution.Shares)
	req.Empty(streams.Tombstones)
	req.Empty(streams.PendingWrites)
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
			name: "activation epoch mismatch",
			mutate: func(config *RewardMigrationConfig) {
				config.Streams[0].ActivationEpoch++
			},
			expected: "activation epoch 101 does not match upgrade epoch 100",
		},
		{
			name: "weight start mismatch",
			mutate: func(config *RewardMigrationConfig) {
				config.Streams[1].Weight.TStart++
			},
			expected: "weight start 101 does not match upgrade epoch 100",
		},
		{
			name: "wrong stream count",
			mutate: func(config *RewardMigrationConfig) {
				config.Streams = config.Streams[:1]
			},
			expected: "requires exactly two streams",
		},
		{
			name: "wrong bootstrap weight",
			mutate: func(config *RewardMigrationConfig) {
				config.Streams[0].Weight.VStart--
			},
			expected: "consensus bootstrap weight is invalid",
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := validRewardMigrationConfig(t, activationEpoch)
			tc.mutate(&config)
			_, err := newRewardMigrator(config, activationEpoch, outCodeCID)
			require.ErrorContains(t, err, tc.expected)
		})
	}
}
