package miner_test

import (
	"testing"

	abi "github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	v14miner "github.com/filecoin-project/go-state-types/builtin/v14/miner"
	"github.com/filecoin-project/go-state-types/builtin/v16/miner"
	"github.com/stretchr/testify/require"
)

func TestQualityForWeight(t *testing.T) {
	emptyQuality := big.NewInt(1 << builtin.SectorQualityPrecision)
	verifiedQuality := big.Mul(emptyQuality, big.Div(builtin.VerifiedDealWeightMultiplier, builtin.QualityBaseMultiplier))
	halfVerifiedQuality := big.Add(big.Div(emptyQuality, big.NewInt(2)), big.Div(verifiedQuality, big.NewInt(2)))

	sizeRange := []abi.SectorSize{
		abi.SectorSize(2 << 10),
		abi.SectorSize(8 << 20),
		abi.SectorSize(512 << 20),
		abi.SectorSize(32 << 30),
		abi.SectorSize(64 << 30),
	}
	durationRange := []abi.ChainEpoch{
		abi.ChainEpoch(1),
		abi.ChainEpoch(10),
		abi.ChainEpoch(1000),
		1000 * builtin.EpochsInDay,
	}

	for _, size := range sizeRange {
		for _, duration := range durationRange {
			fullWeight := big.NewInt(int64(size) * int64(duration))
			halfWeight := big.Div(fullWeight, big.NewInt(2))

			require.Equal(t, emptyQuality, miner.QualityForWeight(size, duration, big.Zero()))
			require.Equal(t, verifiedQuality, miner.QualityForWeight(size, duration, fullWeight))
			require.Equal(t, halfVerifiedQuality, miner.QualityForWeight(size, duration, halfWeight))

			// test against old form that takes a dealWeight argument
			require.Equal(t, emptyQuality, v14miner.QualityForWeight(size, duration, big.Zero(), big.Zero()))
			require.Equal(t, emptyQuality, v14miner.QualityForWeight(size, duration, halfWeight, big.Zero()))
			require.Equal(t, emptyQuality, v14miner.QualityForWeight(size, duration, fullWeight, big.Zero()))
			require.Equal(t, verifiedQuality, v14miner.QualityForWeight(size, duration, big.Zero(), fullWeight))
			require.Equal(t, verifiedQuality, v14miner.QualityForWeight(size, duration, fullWeight, fullWeight))
			require.Equal(t, halfVerifiedQuality, v14miner.QualityForWeight(size, duration, big.Zero(), halfWeight))
			require.Equal(t, halfVerifiedQuality, v14miner.QualityForWeight(size, duration, halfWeight, halfWeight))
		}
	}
}

// matches builtin-actors/actors/miner/tests/policy_test.rs
func TestDailyProofFeeCalc(t *testing.T) {
	// Given a CS of 680M FIL, 32GiB QAP, a fee multiplier of 7.4e-15 per 32GiB QAP, the daily proof
	// fee should be 5032 nanoFIL.
	//   680M * 7.4e-15 = 0.000005032 FIL
	//   0.000005032 * 1e9 = 5032 nanoFIL
	//   0.000005032 * 1e18 = 5032000000000 attoFIL
	// As a per-byte multiplier we use 2.1536e-25, a close approximation of 7.4e-15 / 32GiB.
	//   680M * 32GiB * 2.1536e-25 = 0.000005031805013354 FIL
	//   0.000005031805013354 * 1e18 = 5031805013354 attoFIL
	circulatingSupply := big.Mul(big.NewInt(680_000_000), builtin.TokenPrecision)
	ref32GibFee := big.NewInt(3780793052776)

	for _, tc := range []struct {
		size     uint64
		expected big.Int
	}{
		{32, ref32GibFee},
		{64, big.Mul(ref32GibFee, big.NewInt(2))},
		{32 * 10, big.Mul(ref32GibFee, big.NewInt(10))},
		{32 * 5, big.Mul(ref32GibFee, big.NewInt(5))},
		{64 * 10, big.Mul(ref32GibFee, big.NewInt(20))},
	} {
		power := big.NewInt(int64(tc.size << 30)) // 32GiB raw QAP
		fee := miner.DailyProofFee(circulatingSupply, power)
		delta := big.Sub(fee, tc.expected).Abs()
		require.LessOrEqual(t, delta.Uint64(), uint64(10),
			"size: %s, fee: %s, expected_fee: %s (±10)", tc.size, fee, tc.expected)
	}
}
