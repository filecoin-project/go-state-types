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

func TestReplicaID(t *testing.T) {
	actor := abi.ActorID(123)
	sector := abi.SectorNumber(1)
	ticket := abi.SealRandomness{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2}
	commd := []byte{
		0xfc, 0x7e, 0x92, 0x82, 0x96, 0xe5, 0x16, 0xfa,
		0xad, 0xe9, 0x86, 0xb2, 0x8f, 0x92, 0xd4, 0x4a,
		0x4f, 0x24, 0xb9, 0x35, 0x48, 0x52, 0x23, 0x37,
		0x6a, 0x79, 0x90, 0x27, 0xbc, 0x18, 0xf8, 0x33,
	}

	spt := abi.RegisteredSealProof_StackedDrg2KiBV1

	out, err := spt.ReplicaId(actor, sector, ticket, commd)
	require.NoError(t, err)

	// extracted from rust-fil-proofs from a PreCommit1 call
	expect := [32]byte{
		0x1e, 0x2e, 0x06, 0x39, 0x1a, 0x0d, 0x1b, 0x99,
		0xdc, 0x89, 0x6a, 0xa6, 0xf7, 0xfa, 0x11, 0x3f,
		0xf8, 0xe9, 0x93, 0xcb, 0xdc, 0xbf, 0xb5, 0x03,
		0xe9, 0x95, 0x71, 0x60, 0x3b, 0x61, 0xa0, 0x16,
	}

	assert.EqualValues(t, expect, out)
}
