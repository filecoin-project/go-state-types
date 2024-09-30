package miner_test

import (
	"testing"

	abi "github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v15/miner"
	"github.com/filecoin-project/go-state-types/builtin/v15/util/smoothing"
)

// See filecoin-project/builtin-actors actors/miner/tests/fip0081_initial_pledge.rs
func TestInitialPledgeForPowerFip0081(t *testing.T) {
	filPrecision := big.NewInt(1_000_000_000_000_000_000)

	epochTargetReward := abi.TokenAmount(big.Zero())
	qaSectorPower := abi.StoragePower(big.NewInt(1 << 36))
	networkQAPower := abi.StoragePower(big.NewInt(1 << 10))
	powerRateOfChange := abi.StoragePower(big.NewInt(1 << 10))
	rewardEstimate := smoothing.FilterEstimate{
		PositionEstimate: epochTargetReward,
		VelocityEstimate: big.Zero(),
	}
	powerEstimate := smoothing.FilterEstimate{
		PositionEstimate: networkQAPower,
		VelocityEstimate: powerRateOfChange,
	}
	circulatingSupply := abi.TokenAmount(filPrecision)

	testCases := []struct {
		name                  string
		qaSectorPower         abi.StoragePower
		rewardEstimate        smoothing.FilterEstimate
		powerEstimate         smoothing.FilterEstimate
		circulatingSupply     abi.TokenAmount
		epochsSinceRampStart  int64
		rampDurationEpochs    uint64
		expectedInitialPledge abi.TokenAmount
	}{
		{
			name:                  "pre-ramp where 'baseline power' dominates (negative epochsSinceRampStart)",
			qaSectorPower:         qaSectorPower,
			rewardEstimate:        rewardEstimate,
			powerEstimate:         powerEstimate,
			circulatingSupply:     circulatingSupply,
			epochsSinceRampStart:  -100,
			rampDurationEpochs:    100,
			expectedInitialPledge: big.Add(big.Div(big.Mul(big.NewInt(1500), filPrecision), big.NewInt(10000)), big.NewInt(1)),
		},
		{
			name:                  "zero ramp duration",
			qaSectorPower:         qaSectorPower,
			rewardEstimate:        rewardEstimate,
			powerEstimate:         powerEstimate,
			circulatingSupply:     circulatingSupply,
			epochsSinceRampStart:  0,
			rampDurationEpochs:    0,
			expectedInitialPledge: big.Add(big.Div(big.Mul(big.NewInt(1950), filPrecision), big.NewInt(10000)), big.NewInt(1)),
		},
		{
			name:                  "zero ramp duration (10 epochs since)",
			qaSectorPower:         qaSectorPower,
			rewardEstimate:        rewardEstimate,
			powerEstimate:         powerEstimate,
			circulatingSupply:     circulatingSupply,
			epochsSinceRampStart:  10,
			rampDurationEpochs:    0,
			expectedInitialPledge: big.Add(big.Div(big.Mul(big.NewInt(1950), filPrecision), big.NewInt(10000)), big.NewInt(1)),
		},
		{
			name:                  "pre-ramp where 'baseline power' dominates",
			qaSectorPower:         qaSectorPower,
			rewardEstimate:        rewardEstimate,
			powerEstimate:         powerEstimate,
			circulatingSupply:     circulatingSupply,
			epochsSinceRampStart:  0,
			rampDurationEpochs:    100,
			expectedInitialPledge: big.Add(big.Div(big.Mul(big.NewInt(1500), filPrecision), big.NewInt(10000)), big.NewInt(1)),
		},
		{
			name:                  "on-ramp where 'baseline power' is at 85% and `simple power` is at 15%",
			qaSectorPower:         qaSectorPower,
			rewardEstimate:        rewardEstimate,
			powerEstimate:         powerEstimate,
			circulatingSupply:     circulatingSupply,
			epochsSinceRampStart:  50,
			rampDurationEpochs:    100,
			expectedInitialPledge: big.Add(big.Div(big.Mul(big.NewInt(1725), filPrecision), big.NewInt(10000)), big.NewInt(1)),
		},
		{
			name:                  "after-ramp where 'baseline power' is at 70% and `simple power` is at 30%",
			qaSectorPower:         qaSectorPower,
			rewardEstimate:        rewardEstimate,
			powerEstimate:         powerEstimate,
			circulatingSupply:     circulatingSupply,
			epochsSinceRampStart:  150,
			rampDurationEpochs:    100,
			expectedInitialPledge: big.Add(big.Div(big.Mul(big.NewInt(1950), filPrecision), big.NewInt(10000)), big.NewInt(1)),
		},
		{
			name:                  "on-ramp where 'baseline power' has reduced effect (97%)",
			qaSectorPower:         qaSectorPower,
			rewardEstimate:        rewardEstimate,
			powerEstimate:         powerEstimate,
			circulatingSupply:     circulatingSupply,
			epochsSinceRampStart:  10,
			rampDurationEpochs:    100,
			expectedInitialPledge: big.Add(big.Div(big.Mul(big.NewInt(1545), filPrecision), big.NewInt(10000)), big.NewInt(1)),
		},
		{
			name:                  "on-ramp, first epoch, pledge should be 97% 'baseline' + 3% simple",
			qaSectorPower:         qaSectorPower,
			rewardEstimate:        rewardEstimate,
			powerEstimate:         powerEstimate,
			circulatingSupply:     circulatingSupply,
			epochsSinceRampStart:  1,
			rampDurationEpochs:    10,
			expectedInitialPledge: big.Add(big.Div(big.Mul(big.NewInt(1545), filPrecision), big.NewInt(10000)), big.NewInt(1)),
		},
		{
			name:                  "validate pledges 1 epoch before and after ramp start: before ramp start",
			qaSectorPower:         qaSectorPower,
			rewardEstimate:        rewardEstimate,
			powerEstimate:         powerEstimate,
			circulatingSupply:     circulatingSupply,
			epochsSinceRampStart:  -1,
			rampDurationEpochs:    10,
			expectedInitialPledge: big.Add(big.Div(big.Mul(big.NewInt(1500), filPrecision), big.NewInt(10000)), big.NewInt(1)),
		},
		{
			name:                  "validate pledges 1 epoch before and after ramp start: at ramp start",
			qaSectorPower:         qaSectorPower,
			rewardEstimate:        rewardEstimate,
			powerEstimate:         powerEstimate,
			circulatingSupply:     circulatingSupply,
			epochsSinceRampStart:  0,
			rampDurationEpochs:    10,
			expectedInitialPledge: big.Add(big.Div(big.Mul(big.NewInt(1500), filPrecision), big.NewInt(10000)), big.NewInt(1)),
		},
		{
			name:                  "validate pledges 1 epoch before and after ramp start: on ramp start",
			qaSectorPower:         qaSectorPower,
			rewardEstimate:        rewardEstimate,
			powerEstimate:         powerEstimate,
			circulatingSupply:     circulatingSupply,
			epochsSinceRampStart:  1,
			rampDurationEpochs:    10,
			expectedInitialPledge: big.Add(big.Div(big.Mul(big.NewInt(1545), filPrecision), big.NewInt(10000)), big.NewInt(1)),
		},
		{
			name:                  "post-ramp where 'baseline power' has reduced effect (70%)",
			qaSectorPower:         qaSectorPower,
			rewardEstimate:        rewardEstimate,
			powerEstimate:         powerEstimate,
			circulatingSupply:     circulatingSupply,
			epochsSinceRampStart:  500,
			rampDurationEpochs:    100,
			expectedInitialPledge: big.Add(big.Div(big.Mul(big.NewInt(1950), filPrecision), big.NewInt(10000)), big.NewInt(1)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initialPledge := miner.InitialPledgeForPower(
				tc.qaSectorPower,
				abi.StoragePower(big.NewInt(1<<37)),
				tc.rewardEstimate,
				tc.powerEstimate,
				tc.circulatingSupply,
				tc.epochsSinceRampStart,
				tc.rampDurationEpochs,
			)
			if !initialPledge.Equals(tc.expectedInitialPledge) {
				t.Fatalf("expected initial pledge %v, got %v", tc.expectedInitialPledge, initialPledge)
			}
		})
	}
}
