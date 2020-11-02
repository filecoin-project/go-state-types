package abi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/go-state-types/abi"
)

func TestSectorSizeString(t *testing.T) {
	assert.Equal(t, "0", abi.SectorSize(0).String())
	assert.Equal(t, "1", abi.SectorSize(1).String())
	assert.Equal(t, "1024", abi.SectorSize(1024).String())
	assert.Equal(t, "1234", abi.SectorSize(1234).String())
	assert.Equal(t, "1125899906842624", abi.SectorSize(1125899906842624).String())
}

func TestSectorSizeShortString(t *testing.T) {
	kib := uint64(1024)
	pib := uint64(1125899906842624)

	assert.Equal(t, "0B", abi.SectorSize(0).ShortString())
	assert.Equal(t, "1B", abi.SectorSize(1).ShortString())
	assert.Equal(t, "1023B", abi.SectorSize(1023).ShortString())
	assert.Equal(t, "1KiB", abi.SectorSize(kib).ShortString())
	assert.Equal(t, "1KiB", abi.SectorSize(kib+1).ShortString())   // truncated
	assert.Equal(t, "1KiB", abi.SectorSize(kib*2-1).ShortString()) // truncated
	assert.Equal(t, "2KiB", abi.SectorSize(kib*2).ShortString())
	assert.Equal(t, "2KiB", abi.SectorSize(kib*2+1).ShortString()) // truncated
	assert.Equal(t, "1023KiB", abi.SectorSize(kib*1023).ShortString())
	assert.Equal(t, "1MiB", abi.SectorSize(1048576).ShortString())
	assert.Equal(t, "1GiB", abi.SectorSize(1073741824).ShortString())
	assert.Equal(t, "1TiB", abi.SectorSize(1099511627776).ShortString())
	assert.Equal(t, "1PiB", abi.SectorSize(pib).ShortString())
	assert.Equal(t, "1EiB", abi.SectorSize(pib*kib).ShortString())
	assert.Equal(t, "10EiB", abi.SectorSize(pib*kib*10).ShortString())
}

func TestV1_1SealProofEquivalence(t *testing.T) {
	for v1 := abi.RegisteredSealProof_StackedDrg2KiBV1; v1 < abi.RegisteredSealProof_StackedDrg2KiBV1_1; v1++ {
		v1Size, err := v1.SectorSize()
		require.NoError(t, err)
		v1Wpost, err := v1.RegisteredWindowPoStProof()
		require.NoError(t, err)
		v1Gpost, err := v1.RegisteredWinningPoStProof()
		require.NoError(t, err)

		v1_1 := v1 + 5 // RegisteredSealProof_StackedDrgXxxV1_1
		v11Size, err := v1_1.SectorSize()
		require.NoError(t, err)
		v11Wpost, err := v1_1.RegisteredWindowPoStProof()
		require.NoError(t, err)
		v11Gpost, err := v1_1.RegisteredWinningPoStProof()
		require.NoError(t, err)

		assert.Equal(t, v1Size, v11Size)
		assert.Equal(t, v1Wpost, v11Wpost)
		assert.Equal(t, v1Gpost, v11Gpost)

		gpost := abi.RegisteredPoStProof(v1) // RegisteredPoStProof_StackedDrgWinningXxxV1
		gpostSize, err := gpost.SectorSize()
		require.NoError(t, err)
		assert.Equal(t, v1Size, gpostSize)

		wpost := gpost + 5 // RegisteredPoStProof_StackedDrgWindowXxxV1
		wpostSize, err := wpost.SectorSize()
		require.NoError(t, err)
		assert.Equal(t, v1Size, wpostSize)
	}
}
