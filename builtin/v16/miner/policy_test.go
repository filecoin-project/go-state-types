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
