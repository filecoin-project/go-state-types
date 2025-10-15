package miner_test

import (
	"fmt"
	"testing"

	abi "github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v18/miner"
	"github.com/filecoin-project/go-state-types/builtin/v18/util/smoothing"
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

func TestNegativeBRClamp(t *testing.T) {
	epochTargetReward := big.NewInt(1 << 50)
	qaSectorPower := abi.StoragePower(big.NewInt(1 << 36))
	networkQAPower := abi.StoragePower(big.NewInt(1 << 10))
	powerRateOfChange := big.NewInt(1 << 10).Neg()
	rewardEstimate := smoothing.FilterEstimate{
		PositionEstimate: abi.TokenAmount(epochTargetReward),
		VelocityEstimate: big.Zero(),
	}
	powerEstimate := smoothing.FilterEstimate{
		PositionEstimate: networkQAPower,
		VelocityEstimate: powerRateOfChange,
	}

	if big.Add(powerEstimate.PositionEstimate, big.Mul(powerEstimate.VelocityEstimate, big.NewInt(4))).GreaterThan(networkQAPower) {
		t.Fatalf("power estimate extrapolated incorrectly")
	}

	fourBR := miner.ExpectedRewardForPower(rewardEstimate, powerEstimate, qaSectorPower, 4)
	if !fourBR.IsZero() {
		t.Fatalf("expected zero BR, got %v", fourBR)
	}
}

func TestZeroPowerMeansZeroFaultPenalty(t *testing.T) {
	epochTargetReward := big.NewInt(1 << 50)
	zeroQAPower := abi.StoragePower(big.Zero())
	networkQAPower := abi.StoragePower(big.NewInt(1 << 10))
	powerRateOfChange := big.NewInt(1 << 10)
	rewardEstimate := smoothing.FilterEstimate{
		PositionEstimate: abi.TokenAmount(epochTargetReward),
		VelocityEstimate: big.Zero(),
	}
	powerEstimate := smoothing.FilterEstimate{
		PositionEstimate: networkQAPower,
		VelocityEstimate: powerRateOfChange,
	}

	penaltyForZeroPowerFaulted := miner.PledgePenaltyForContinuedFault(rewardEstimate, powerEstimate, zeroQAPower)
	if !penaltyForZeroPowerFaulted.IsZero() {
		t.Fatalf("expected zero penalty, got %v", penaltyForZeroPowerFaulted)
	}
}

func TestAggregatePowerPledgePenaltyForContinuedFault(t *testing.T) {
	epochTargetReward := big.NewInt(1 << 50)
	networkQAPower := abi.StoragePower(big.NewInt(1 << 10))
	powerRateOfChange := big.NewInt(1 << 10)
	rewardEstimate := smoothing.NewEstimate(abi.TokenAmount(epochTargetReward), big.Zero())
	powerEstimate := smoothing.NewEstimate(networkQAPower, powerRateOfChange)

	testCases := []struct {
		sectorMultiple int64
		qaPower        abi.StoragePower
	}{
		{10, abi.StoragePower(big.NewInt(1 << 6))},
		{10, abi.StoragePower(big.NewInt(1 << 36))},
		{10, abi.StoragePower(big.NewInt(1 << 50))},
		{1000, abi.StoragePower(big.NewInt(1 << 6))},
		{1000, abi.StoragePower(big.NewInt(1 << 36))},
		{1000, abi.StoragePower(big.NewInt(1 << 50))},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d sectors, %s qap", tc.sectorMultiple, tc.qaPower), func(t *testing.T) {
			sectorMultiple := tc.sectorMultiple
			qaPower := tc.qaPower

			aggregatePenalty := miner.PledgePenaltyForContinuedFault(
				rewardEstimate,
				powerEstimate,
				big.Mul(qaPower, big.NewInt(sectorMultiple)),
			)

			individualPenalties := big.Zero()
			for i := int64(0); i < sectorMultiple; i++ {
				individualPenalty := miner.PledgePenaltyForContinuedFault(rewardEstimate, powerEstimate, qaPower)
				individualPenalties = big.Add(individualPenalties, individualPenalty)
			}
			if aggregatePenalty.LessThanEqual(big.Zero()) {
				t.Fatalf("aggregate penalty is not positive: %s", aggregatePenalty)
			}

			diff := big.Sub(aggregatePenalty, individualPenalties).Abs()
			allowedAttoDifference := big.NewInt(sectorMultiple)
			if diff.GreaterThan(allowedAttoDifference) {
				t.Fatalf("aggregate_penalty: %v, individual_penalties: %v, diff: %v", aggregatePenalty, individualPenalties, diff)
			}
		})
	}
}

func TestPledgePenaltyForTermination(t *testing.T) {
	t.Run("when sector age exceeds cap returns percentage of initial pledge", func(t *testing.T) {
		sectorAgeInDays := miner.TerminationLifetimeCap + 1
		sectorAge := sectorAgeInDays * builtin.EpochsInDay

		initialPledge := abi.NewTokenAmount(1 << 10)
		faultFee := abi.NewTokenAmount(0)
		fee := miner.PledgePenaltyForTermination(initialPledge, sectorAge, faultFee)

		expectedFee := big.Div(big.Mul(initialPledge, miner.TermFeePledgeMultiple.Numerator), miner.TermFeePledgeMultiple.Denominator)
		if !fee.Equals(abi.TokenAmount(expectedFee)) {
			t.Fatalf("expected fee %v, got %v", expectedFee, fee)
		}
	})

	t.Run("when sector age below cap returns percentage of initial pledge percentage", func(t *testing.T) {
		sectorAgeInDays := miner.TerminationLifetimeCap / 2
		sectorAge := sectorAgeInDays * builtin.EpochsInDay

		initialPledge := abi.NewTokenAmount(1 << 10)
		faultFee := abi.NewTokenAmount(0)
		fee := miner.PledgePenaltyForTermination(initialPledge, sectorAge, faultFee)

		simpleTerminationFee := big.Div(big.Mul(initialPledge, miner.TermFeePledgeMultiple.Numerator), miner.TermFeePledgeMultiple.Denominator)
		expectedFee := big.Div(big.Mul(simpleTerminationFee, big.NewInt(int64(sectorAgeInDays))), big.NewInt(int64(miner.TerminationLifetimeCap)))

		if !fee.Equals(abi.TokenAmount(expectedFee)) {
			t.Fatalf("expected fee %v, got %v", expectedFee, fee)
		}
	})

	t.Run("when termination fee less than fault fee returns multiple of fault fee", func(t *testing.T) {
		sectorAgeInDays := miner.TerminationLifetimeCap + 1
		sectorAge := sectorAgeInDays * builtin.EpochsInDay

		initialPledge := abi.NewTokenAmount(1 << 10)
		faultFee := abi.NewTokenAmount(1 << 10)
		fee := miner.PledgePenaltyForTermination(initialPledge, sectorAge, faultFee)

		expectedFee := big.Div(big.Mul(faultFee, miner.TermFeeMaxFaultFeeMultiple.Numerator), miner.TermFeeMaxFaultFeeMultiple.Denominator)
		if !fee.Equals(abi.TokenAmount(expectedFee)) {
			t.Fatalf("expected fee %v, got %v", expectedFee, fee)
		}
	})

	t.Run("when termination fee less than minimum returns minimum", func(t *testing.T) {
		sectorAge := abi.ChainEpoch(0)

		initialPledge := abi.NewTokenAmount(1 << 10)
		faultFee := abi.NewTokenAmount(0)
		fee := miner.PledgePenaltyForTermination(initialPledge, sectorAge, faultFee)

		expectedFee := big.Div(big.Mul(initialPledge, miner.TermFeeMinPledgeMultiple.Numerator), miner.TermFeeMinPledgeMultiple.Denominator)
		if !fee.Equals(abi.TokenAmount(expectedFee)) {
			t.Fatalf("expected fee %v, got %v", expectedFee, fee)
		}
	})
}
