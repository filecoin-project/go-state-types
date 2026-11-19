package miner_test

import (
	"fmt"
	"testing"

	abi "github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v20/miner"
	"github.com/filecoin-project/go-state-types/builtin/v20/util/smoothing"
)

// FIP-0081 is fully activated: additional pledge is 70% baseline and 30% simple pledge.
func TestInitialPledgeForPowerFip0081(t *testing.T) {
	const sectorPower int64 = 1 << 36
	filPrecision := big.NewInt(1_000_000_000_000_000_000)

	testCases := []struct {
		name                  string
		baselinePower         int64
		networkQAPower        int64
		epochReward           int64
		circulatingSupply     abi.TokenAmount
		expectedInitialPledge abi.TokenAmount
	}{
		{
			name:                  "70 percent baseline and 30 percent simple pledge",
			baselinePower:         2 * sectorPower,
			networkQAPower:        sectorPower,
			circulatingSupply:     filPrecision,
			expectedInitialPledge: big.NewInt(195_000_000_000_000_001),
		},
		{
			name:                  "network power exceeds baseline",
			baselinePower:         2 * sectorPower,
			networkQAPower:        4 * sectorPower,
			circulatingSupply:     filPrecision,
			expectedInitialPledge: big.NewInt(75_000_000_000_000_001),
		},
		{
			name:                  "sector power bounds both denominators",
			baselinePower:         sectorPower / 2,
			networkQAPower:        sectorPower / 4,
			circulatingSupply:     filPrecision,
			expectedInitialPledge: big.NewInt(300_000_000_000_000_001),
		},
		{
			name:              "weighted pledge terms are rounded separately",
			baselinePower:     2 * sectorPower,
			networkQAPower:    sectorPower,
			circulatingSupply: big.NewInt(16),
			// floor(1.68) + floor(1.44) + 1 attoFIL minimum reward.
			expectedInitialPledge: big.NewInt(3),
		},
		{
			name:              "projected reward is added to consensus pledge",
			baselinePower:     2 * sectorPower,
			networkQAPower:    sectorPower,
			epochReward:       1,
			circulatingSupply: filPrecision,
			// At full network power, 20 days of reward adds 57,600 attoFIL.
			expectedInitialPledge: big.NewInt(195_000_000_000_057_600),
		},
		{
			name:                  "pledge is capped",
			baselinePower:         2 * sectorPower,
			networkQAPower:        sectorPower,
			circulatingSupply:     big.Mul(filPrecision, big.NewInt(100)),
			expectedInitialPledge: big.Mul(miner.InitialPledgeMaxPerByte, big.NewInt(sectorPower)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initialPledge := miner.InitialPledgeForPower(
				abi.NewStoragePower(sectorPower),
				abi.NewStoragePower(tc.baselinePower),
				smoothing.NewEstimate(big.NewInt(tc.epochReward), big.Zero()),
				smoothing.NewEstimate(big.NewInt(tc.networkQAPower), big.Zero()),
				tc.circulatingSupply,
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
