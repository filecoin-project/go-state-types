package miner

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v16/util/math"
	"github.com/filecoin-project/go-state-types/builtin/v16/util/smoothing"
)

// Projection period of expected sector block reward for deposit required to pre-commit a sector.
// This deposit is lost if the pre-commitment is not timely followed up by a commitment proof.
var PreCommitDepositFactor = 20
var PreCommitDepositProjectionPeriod = abi.ChainEpoch(PreCommitDepositFactor) * builtin.EpochsInDay

// Projection period of expected sector block rewards for storage pledge required to commit a sector.
// This pledge is lost if a sector is terminated before its full committed lifetime.
var InitialPledgeFactor = 20
var InitialPledgeProjectionPeriod = abi.ChainEpoch(InitialPledgeFactor) * builtin.EpochsInDay

// Cap on initial pledge requirement for sectors.
// The target is 1 FIL (10**18 attoFIL) per 32GiB.
// This does not divide evenly, so the result is fractionally smaller.
var InitialPledgeMaxPerByte = big.Div(big.NewInt(1e18), big.NewInt(32<<30))

// Multiplier of share of circulating money supply for consensus pledge required to commit a sector.
// This pledge is lost if a sector is terminated before its full committed lifetime.
var InitialPledgeLockTarget = builtin.BigFrac{
	Numerator:   big.NewInt(3),
	Denominator: big.NewInt(10),
}

const GammaFixedPointFactor = 1000 // 3 decimal places

// The projected block reward a sector would earn over some period.
// Also known as "BR(t)".
// BR(t) = ProjectedRewardFraction(t) * SectorQualityAdjustedPower
// ProjectedRewardFraction(t) is the sum of estimated reward over estimated total power
// over all epochs in the projection period [t t+projectionDuration]
func ExpectedRewardForPower(rewardEstimate, networkQAPowerEstimate smoothing.FilterEstimate, qaSectorPower abi.StoragePower, projectionDuration abi.ChainEpoch) abi.TokenAmount {
	networkQAPowerSmoothed := smoothing.Estimate(&networkQAPowerEstimate)
	if networkQAPowerSmoothed.IsZero() {
		return smoothing.Estimate(&rewardEstimate)
	}
	expectedRewardForProvingPeriod := smoothing.ExtrapolatedCumSumOfRatio(projectionDuration, 0, rewardEstimate, networkQAPowerEstimate)
	br128 := big.Mul(qaSectorPower, expectedRewardForProvingPeriod) // Q.0 * Q.128 => Q.128
	br := big.Rsh(br128, math.Precision128)

	return big.Max(br, big.Zero())
}

// BR but zero values are clamped at 1 attofil
// Some uses of BR (PCD, IP) require a strictly positive value for BR derived values so
// accounting variables can be used as succinct indicators of miner activity.
func ExpectedRewardForPowerClampedAtAttoFIL(rewardEstimate, networkQAPowerEstimate smoothing.FilterEstimate, qaSectorPower abi.StoragePower, projectionDuration abi.ChainEpoch) abi.TokenAmount {
	br := ExpectedRewardForPower(rewardEstimate, networkQAPowerEstimate, qaSectorPower, projectionDuration)
	if br.LessThanEqual(big.Zero()) {
		br = abi.NewTokenAmount(1)
	}
	return br
}

// Computes the PreCommit deposit given sector qa weight and current network conditions.
// PreCommit Deposit = BR(PreCommitDepositProjectionPeriod)
func PreCommitDepositForPower(rewardEstimate, networkQAPowerEstimate smoothing.FilterEstimate, qaSectorPower abi.StoragePower) abi.TokenAmount {
	return ExpectedRewardForPowerClampedAtAttoFIL(rewardEstimate, networkQAPowerEstimate, qaSectorPower, PreCommitDepositProjectionPeriod)
}

// InitialPledgeForPower computes the pledge requirement for committing new quality-adjusted power
// to the network, given the current network total and baseline power, per-epoch  reward, and
// circulating token supply.
// The pledge comprises two parts:
// - storage pledge, aka IP base: a multiple of the reward expected to be earned by newly-committed power
// - consensus pledge, aka additional IP: a pro-rata fraction of the circulating money supply
//
// IP = IPBase(t) + AdditionalIP(t)
// IPBase(t) = BR(t, InitialPledgeProjectionPeriod)
// AdditionalIP(t) = LockTarget(t)*PledgeShare(t)
// LockTarget = (LockTargetFactorNum / LockTargetFactorDenom) * FILCirculatingSupply(t)
// PledgeShare(t) = sectorQAPower / max(BaselinePower(t), NetworkQAPower(t))
func InitialPledgeForPower(
	qaPower,
	baselinePower abi.StoragePower,
	rewardEstimate,
	networkQAPowerEstimate smoothing.FilterEstimate,
	circulatingSupply abi.TokenAmount,
	epochsSinceRampStart int64,
	rampDurationEpochs uint64,
) abi.TokenAmount {
	ipBase := ExpectedRewardForPowerClampedAtAttoFIL(rewardEstimate, networkQAPowerEstimate, qaPower, InitialPledgeProjectionPeriod)

	lockTargetNum := big.Mul(InitialPledgeLockTarget.Numerator, circulatingSupply)
	lockTargetDenom := InitialPledgeLockTarget.Denominator
	pledgeShareNum := qaPower
	networkQAPower := smoothing.Estimate(&networkQAPowerEstimate)

	// Once FIP-0081 has fully activated, additional pledge will be 70% baseline
	// pledge + 30% simple pledge.
	const fip0081ActivationPermille = 300
	// Gamma/GAMMA_FIXED_POINT_FACTOR is the share of pledge coming from the
	// baseline formulation, with 1-(gamma/GAMMA_FIXED_POINT_FACTOR) coming from
	// simple pledge.
	// gamma = 1000 - 300 * (epochs_since_ramp_start / ramp_duration_epochs).max(0).min(1)
	var skew uint64
	switch {
	case epochsSinceRampStart < 0:
		// No skew before ramp start
		skew = 0
	case rampDurationEpochs == 0 || epochsSinceRampStart >= int64(rampDurationEpochs):
		// 100% skew after ramp end
		skew = fip0081ActivationPermille
	case epochsSinceRampStart > 0:
		skew = (uint64(epochsSinceRampStart*fip0081ActivationPermille) / rampDurationEpochs)
	}
	gamma := big.NewInt(int64(GammaFixedPointFactor - skew))

	additionalIPNum := big.Mul(lockTargetNum, pledgeShareNum)

	pledgeShareDenomBaseline := big.Max(big.Max(networkQAPower, baselinePower), qaPower)
	pledgeShareDenomSimple := big.Max(networkQAPower, qaPower)

	additionalIPDenomBaseline := big.Mul(pledgeShareDenomBaseline, lockTargetDenom)
	additionalIPBaseline := big.Div(big.Mul(gamma, additionalIPNum), big.Mul(additionalIPDenomBaseline, big.NewInt(GammaFixedPointFactor)))
	additionalIPDenomSimple := big.Mul(pledgeShareDenomSimple, lockTargetDenom)
	additionalIPSimple := big.Div(big.Mul(big.Sub(big.NewInt(GammaFixedPointFactor), gamma), additionalIPNum), big.Mul(additionalIPDenomSimple, big.NewInt(GammaFixedPointFactor)))

	// convex combination of simple and baseline pledge
	additionalIP := big.Add(additionalIPBaseline, additionalIPSimple)

	nominalPledge := big.Add(ipBase, additionalIP)
	pledgeCap := big.Mul(InitialPledgeMaxPerByte, qaPower)

	return big.Min(nominalPledge, pledgeCap)
}

// Maximum number of lifetime days penalized when a sector is terminated.
const TerminationLifetimeCap abi.ChainEpoch = 140

// Used to compute termination fees in the base case by multiplying against initial pledge.
var TermFeePledgeMultiple = builtin.BigFrac{
	Numerator:   big.NewInt(85),
	Denominator: big.NewInt(1000),
}

// Used to ensure the termination fee for young sectors is not arbitrarily low.
var TermFeeMinPledgeMultiple = builtin.BigFrac{
	Numerator:   big.NewInt(2),
	Denominator: big.NewInt(100),
}

// Used to compute termination fees when the termination fee of a sector is less than the fault fee for the same sector.
var TermFeeMaxFaultFeeMultiple = builtin.BigFrac{
	Numerator:   big.NewInt(105),
	Denominator: big.NewInt(100),
}

const ContinuedFaultFactorNum = 351
const ContinuedFaultFactorDenom = 100
const ContinuedFaultProjectionPeriod abi.ChainEpoch = (builtin.EpochsInDay * ContinuedFaultFactorNum) / ContinuedFaultFactorDenom

// PledgePenaltyForContinuedFault calculates the penalty for a sector continuing faulty for another
// proving period.
// It is a projection of the expected reward earned by the sector. Also known as "FF(t)"
func PledgePenaltyForContinuedFault(rewardEstimate smoothing.FilterEstimate, networkQaPowerEstimate smoothing.FilterEstimate, qaSectorPower abi.StoragePower) abi.TokenAmount {
	return ExpectedRewardForPower(rewardEstimate, networkQaPowerEstimate, qaSectorPower, ContinuedFaultProjectionPeriod)
}

// PledgePenaltyForTermination Calculates termination fee for a given sector. Normally, it's
// calculated as a fixed percentage of the initial pledge. However, there are some special cases
// outlined in [FIP-0098](https://github.com/filecoin-project/FIPs/blob/master/FIPS/fip-0098.md).
func PledgePenaltyForTermination(
	initialPledge abi.TokenAmount,
	sectorAge abi.ChainEpoch,
	faultFee abi.TokenAmount,
) abi.TokenAmount {
	// Use the Percentage of the initial pledge strategy to determine the termination fee.
	simpleTerminationFee :=
		big.Div(big.Mul(initialPledge, TermFeePledgeMultiple.Numerator), TermFeePledgeMultiple.Denominator)

	durationTerminationFee :=
		big.Div(big.Mul(big.NewInt(int64(sectorAge)), simpleTerminationFee), big.NewInt(int64(TerminationLifetimeCap*builtin.EpochsInDay)))

	// Apply the age adjustment for young sectors to arrive at the base termination fee.
	baseTerminationFee := big.Min(simpleTerminationFee, durationTerminationFee)

	// Calculate the minimum allowed fee (a lower bound on the termination fee) by comparing the absolute minimum termination fee value against the fault fee. Whatever result is Larger sets the lower bound for the termination fee.
	minimumFeeAbs := big.Div(big.Mul(initialPledge, TermFeeMinPledgeMultiple.Numerator), TermFeeMinPledgeMultiple.Denominator)
	minimumFeeFf := big.Div(big.Mul(faultFee, TermFeeMaxFaultFeeMultiple.Numerator), TermFeeMaxFaultFeeMultiple.Denominator)
	minimumFee := big.Max(minimumFeeAbs, minimumFeeFf)

	return big.Max(baseTerminationFee, minimumFee)
}
