package abi

import (
	"fmt"
	"math"
	"strconv"

	"github.com/filecoin-project/go-state-types/network"

	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/big"
)

// SectorNumber is a numeric identifier for a sector. It is usually relative to a miner.
type SectorNumber uint64

func (s SectorNumber) String() string {
	return strconv.FormatUint(uint64(s), 10)
}

// The maximum assignable sector number.
// Raising this would require modifying our AMT implementation.
const MaxSectorNumber = math.MaxInt64

// SectorSize indicates one of a set of possible sizes in the network.
// Ideally, SectorSize would be an enum
//
//	type SectorSize enum {
//	  1KiB = 1024
//	  1MiB = 1048576
//	  1GiB = 1073741824
//	  1TiB = 1099511627776
//	  1PiB = 1125899906842624
//	  1EiB = 1152921504606846976
//	  max  = 18446744073709551615
//	}
type SectorSize uint64

// Formats the size as a decimal string.
func (s SectorSize) String() string {
	return strconv.FormatUint(uint64(s), 10)
}

// Abbreviates the size as a human-scale number.
// This approximates (truncates) the size unless it is a power of 1024.
func (s SectorSize) ShortString() string {
	var biUnits = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	unit := 0
	for s >= 1024 && unit < len(biUnits)-1 {
		s /= 1024
		unit++
	}
	return fmt.Sprintf("%d%s", s, biUnits[unit])
}

type SectorID struct {
	Miner  ActorID
	Number SectorNumber
}

// The unit of storage power (measured in bytes)
type StoragePower = big.Int

type SectorQuality = big.Int

func NewStoragePower(n int64) StoragePower {
	return big.NewInt(n)
}

// These enumerations must match the proofs library and never change.
type RegisteredSealProof int64

const (
	RegisteredSealProof_StackedDrg2KiBV1   = RegisteredSealProof(0)
	RegisteredSealProof_StackedDrg8MiBV1   = RegisteredSealProof(1)
	RegisteredSealProof_StackedDrg512MiBV1 = RegisteredSealProof(2)
	RegisteredSealProof_StackedDrg32GiBV1  = RegisteredSealProof(3)
	RegisteredSealProof_StackedDrg64GiBV1  = RegisteredSealProof(4)

	RegisteredSealProof_StackedDrg2KiBV1_1   = RegisteredSealProof(5)
	RegisteredSealProof_StackedDrg8MiBV1_1   = RegisteredSealProof(6)
	RegisteredSealProof_StackedDrg512MiBV1_1 = RegisteredSealProof(7)
	RegisteredSealProof_StackedDrg32GiBV1_1  = RegisteredSealProof(8)
	RegisteredSealProof_StackedDrg64GiBV1_1  = RegisteredSealProof(9)

	RegisteredSealProof_StackedDrg2KiBV1_1_Feat_SyntheticPoRep   = RegisteredSealProof(10)
	RegisteredSealProof_StackedDrg8MiBV1_1_Feat_SyntheticPoRep   = RegisteredSealProof(11)
	RegisteredSealProof_StackedDrg512MiBV1_1_Feat_SyntheticPoRep = RegisteredSealProof(12)
	RegisteredSealProof_StackedDrg32GiBV1_1_Feat_SyntheticPoRep  = RegisteredSealProof(13)
	RegisteredSealProof_StackedDrg64GiBV1_1_Feat_SyntheticPoRep  = RegisteredSealProof(14)
)

var Synthetic = map[RegisteredSealProof]bool{
	RegisteredSealProof_StackedDrg2KiBV1_1_Feat_SyntheticPoRep:   true,
	RegisteredSealProof_StackedDrg8MiBV1_1_Feat_SyntheticPoRep:   true,
	RegisteredSealProof_StackedDrg512MiBV1_1_Feat_SyntheticPoRep: true,
	RegisteredSealProof_StackedDrg32GiBV1_1_Feat_SyntheticPoRep:  true,
	RegisteredSealProof_StackedDrg64GiBV1_1_Feat_SyntheticPoRep:  true,
}

type RegisteredPoStProof int64

const (
	RegisteredPoStProof_StackedDrgWinning2KiBV1   = RegisteredPoStProof(0)
	RegisteredPoStProof_StackedDrgWinning8MiBV1   = RegisteredPoStProof(1)
	RegisteredPoStProof_StackedDrgWinning512MiBV1 = RegisteredPoStProof(2)
	RegisteredPoStProof_StackedDrgWinning32GiBV1  = RegisteredPoStProof(3)
	RegisteredPoStProof_StackedDrgWinning64GiBV1  = RegisteredPoStProof(4)

	RegisteredPoStProof_StackedDrgWindow2KiBV1   = RegisteredPoStProof(5)
	RegisteredPoStProof_StackedDrgWindow8MiBV1   = RegisteredPoStProof(6)
	RegisteredPoStProof_StackedDrgWindow512MiBV1 = RegisteredPoStProof(7)
	RegisteredPoStProof_StackedDrgWindow32GiBV1  = RegisteredPoStProof(8)
	RegisteredPoStProof_StackedDrgWindow64GiBV1  = RegisteredPoStProof(9)

	RegisteredPoStProof_StackedDrgWindow2KiBV1_1   = RegisteredPoStProof(10)
	RegisteredPoStProof_StackedDrgWindow8MiBV1_1   = RegisteredPoStProof(11)
	RegisteredPoStProof_StackedDrgWindow512MiBV1_1 = RegisteredPoStProof(12)
	RegisteredPoStProof_StackedDrgWindow32GiBV1_1  = RegisteredPoStProof(13)
	RegisteredPoStProof_StackedDrgWindow64GiBV1_1  = RegisteredPoStProof(14)
)

func (r RegisteredPoStProof) ToV1_1PostProof() (RegisteredPoStProof, error) {
	switch r {
	case RegisteredPoStProof_StackedDrgWindow2KiBV1, RegisteredPoStProof_StackedDrgWindow2KiBV1_1:
		return RegisteredPoStProof_StackedDrgWindow2KiBV1_1, nil
	case RegisteredPoStProof_StackedDrgWindow8MiBV1, RegisteredPoStProof_StackedDrgWindow8MiBV1_1:
		return RegisteredPoStProof_StackedDrgWindow8MiBV1_1, nil
	case RegisteredPoStProof_StackedDrgWindow512MiBV1, RegisteredPoStProof_StackedDrgWindow512MiBV1_1:
		return RegisteredPoStProof_StackedDrgWindow512MiBV1_1, nil
	case RegisteredPoStProof_StackedDrgWindow32GiBV1, RegisteredPoStProof_StackedDrgWindow32GiBV1_1:
		return RegisteredPoStProof_StackedDrgWindow32GiBV1_1, nil
	case RegisteredPoStProof_StackedDrgWindow64GiBV1, RegisteredPoStProof_StackedDrgWindow64GiBV1_1:
		return RegisteredPoStProof_StackedDrgWindow64GiBV1_1, nil
	}

	return -1, xerrors.Errorf("input %d is not a V1 PostProof", r)
}

type RegisteredAggregationProof int64

const (
	RegisteredAggregationProof_SnarkPackV1 = RegisteredAggregationProof(0)
	RegisteredAggregationProof_SnarkPackV2 = RegisteredAggregationProof(1)
)

type RegisteredUpdateProof int64

const (
	RegisteredUpdateProof_StackedDrg2KiBV1   = RegisteredUpdateProof(0)
	RegisteredUpdateProof_StackedDrg8MiBV1   = RegisteredUpdateProof(1)
	RegisteredUpdateProof_StackedDrg512MiBV1 = RegisteredUpdateProof(2)
	RegisteredUpdateProof_StackedDrg32GiBV1  = RegisteredUpdateProof(3)
	RegisteredUpdateProof_StackedDrg64GiBV1  = RegisteredUpdateProof(4)
)

// Metadata about a seal proof type.
type SealProofInfo struct {
	// The proof sizes are 192 * the number of "porep" partitions.
	// https://github.com/filecoin-project/rust-fil-proofs/blob/64390b6fcedb04dd1fdbe43c82b1e91c1439cea2/filecoin-proofs/src/constants.rs#L68-L80
	ProofSize        uint64
	SectorSize       SectorSize
	WinningPoStProof RegisteredPoStProof
	WindowPoStProof  RegisteredPoStProof
	UpdateProof      RegisteredUpdateProof
}

const (
	ss2KiB   = 2 << 10
	ss8MiB   = 8 << 20
	ss512MiB = 512 << 20
	ss32GiB  = 32 << 30
	ss64GiB  = 64 << 30
)

var SealProofInfos = map[RegisteredSealProof]*SealProofInfo{
	RegisteredSealProof_StackedDrg2KiBV1: {
		ProofSize:        192,
		SectorSize:       ss2KiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning2KiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow2KiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg2KiBV1,
	},

	RegisteredSealProof_StackedDrg8MiBV1: {
		ProofSize:        192,
		SectorSize:       ss8MiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning8MiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow8MiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg8MiBV1,
	},
	RegisteredSealProof_StackedDrg512MiBV1: {
		ProofSize:        192,
		SectorSize:       ss512MiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning512MiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow512MiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg512MiBV1,
	},
	RegisteredSealProof_StackedDrg32GiBV1: {
		ProofSize:        1920,
		SectorSize:       ss32GiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning32GiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow32GiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg32GiBV1,
	},
	RegisteredSealProof_StackedDrg64GiBV1: {
		ProofSize:        1920,
		SectorSize:       ss64GiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning64GiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow64GiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg64GiBV1,
	},

	RegisteredSealProof_StackedDrg2KiBV1_1: {
		ProofSize:        192,
		SectorSize:       ss2KiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning2KiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow2KiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg2KiBV1,
	},
	RegisteredSealProof_StackedDrg8MiBV1_1: {
		ProofSize:        192,
		SectorSize:       ss8MiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning8MiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow8MiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg8MiBV1,
	},
	RegisteredSealProof_StackedDrg512MiBV1_1: {
		ProofSize:        192,
		SectorSize:       ss512MiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning512MiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow512MiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg512MiBV1,
	},
	RegisteredSealProof_StackedDrg32GiBV1_1: {
		ProofSize:        1920,
		SectorSize:       ss32GiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning32GiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow32GiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg32GiBV1,
	},
	RegisteredSealProof_StackedDrg64GiBV1_1: {
		ProofSize:        1920,
		SectorSize:       ss64GiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning64GiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow64GiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg64GiBV1,
	},
	RegisteredSealProof_StackedDrg2KiBV1_1_Feat_SyntheticPoRep: {
		ProofSize:        192,
		SectorSize:       ss2KiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning2KiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow2KiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg2KiBV1,
	},
	RegisteredSealProof_StackedDrg8MiBV1_1_Feat_SyntheticPoRep: {
		ProofSize:        192,
		SectorSize:       ss8MiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning8MiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow8MiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg8MiBV1,
	},
	RegisteredSealProof_StackedDrg512MiBV1_1_Feat_SyntheticPoRep: {
		ProofSize:        192,
		SectorSize:       ss512MiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning512MiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow512MiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg512MiBV1,
	},
	RegisteredSealProof_StackedDrg32GiBV1_1_Feat_SyntheticPoRep: {
		ProofSize:        1920,
		SectorSize:       ss32GiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning32GiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow32GiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg32GiBV1,
	},
	RegisteredSealProof_StackedDrg64GiBV1_1_Feat_SyntheticPoRep: {
		ProofSize:        1920,
		SectorSize:       ss64GiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning64GiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow64GiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg64GiBV1,
	},
}

// ProofSize returns the size of seal proofs for the given sector type.
func (p RegisteredSealProof) ProofSize() (uint64, error) {
	info, ok := SealProofInfos[p]
	if !ok {
		return 0, xerrors.Errorf("unsupported proof type: %v", p)
	}
	return info.ProofSize, nil
}

func (p RegisteredSealProof) SectorSize() (SectorSize, error) {
	info, ok := SealProofInfos[p]
	if !ok {
		return 0, xerrors.Errorf("unsupported proof type: %v", p)
	}
	return info.SectorSize, nil
}

// RegisteredWinningPoStProof produces the PoSt-specific RegisteredProof corresponding
// to the receiving RegisteredProof.
func (p RegisteredSealProof) RegisteredWinningPoStProof() (RegisteredPoStProof, error) {
	info, ok := SealProofInfos[p]
	if !ok {
		return 0, xerrors.Errorf("unsupported proof type: %v", p)
	}
	return info.WinningPoStProof, nil
}

// RegisteredWindowPoStProof produces the V1 PoSt-specific RegisteredPoStProof corresponding
// to the receiving RegisteredSealProof.
func (p RegisteredSealProof) RegisteredWindowPoStProof() (RegisteredPoStProof, error) {
	info, ok := SealProofInfos[p]
	if !ok {
		return 0, xerrors.Errorf("unsupported proof type: %v", p)
	}

	return info.WindowPoStProof, nil
}

// RegisteredWindowPoStProofByNetworkVersion produces the V1 PoSt-specific RegisteredPoStProof corresponding
// to the receiving RegisteredSealProof.
// Before nv19, the V1 Proof is returned, from nv19 onwards the v1_1 proof is returned.
func (p RegisteredSealProof) RegisteredWindowPoStProofByNetworkVersion(nv network.Version) (RegisteredPoStProof, error) {
	info, ok := SealProofInfos[p]
	if !ok {
		return 0, xerrors.Errorf("unsupported proof type: %v", p)
	}

	if nv <= network.Version18 {
		return info.WindowPoStProof, nil
	}

	return info.WindowPoStProof.ToV1_1PostProof()
}

// RegisteredUpdateProof produces the Update-specific RegisteredProof corresponding
// to the receiving RegisteredProof.
func (p RegisteredSealProof) RegisteredUpdateProof() (RegisteredUpdateProof, error) {
	info, ok := SealProofInfos[p]
	if !ok {
		return 0, xerrors.Errorf("unsupported proof type: %v", p)
	}
	return info.UpdateProof, nil
}

// Metadata about a PoSt proof type.
type PoStProofInfo struct {
	SectorSize SectorSize

	// Size of a single proof.
	ProofSize uint64
}

var PoStProofInfos = map[RegisteredPoStProof]*PoStProofInfo{
	RegisteredPoStProof_StackedDrgWinning2KiBV1: {
		SectorSize: ss2KiB,
		ProofSize:  192,
	},
	RegisteredPoStProof_StackedDrgWinning8MiBV1: {
		SectorSize: ss8MiB,
		ProofSize:  192,
	},
	RegisteredPoStProof_StackedDrgWinning512MiBV1: {
		SectorSize: ss512MiB,
		ProofSize:  192,
	},
	RegisteredPoStProof_StackedDrgWinning32GiBV1: {
		SectorSize: ss32GiB,
		ProofSize:  192,
	},
	RegisteredPoStProof_StackedDrgWinning64GiBV1: {
		SectorSize: ss64GiB,
		ProofSize:  192,
	},
	RegisteredPoStProof_StackedDrgWindow2KiBV1: {
		SectorSize: ss2KiB,
		ProofSize:  192,
	},
	RegisteredPoStProof_StackedDrgWindow8MiBV1: {
		SectorSize: ss8MiB,
		ProofSize:  192,
	},
	RegisteredPoStProof_StackedDrgWindow512MiBV1: {
		SectorSize: ss512MiB,
		ProofSize:  192,
	},
	RegisteredPoStProof_StackedDrgWindow32GiBV1: {
		SectorSize: ss32GiB,
		ProofSize:  192,
	},
	RegisteredPoStProof_StackedDrgWindow64GiBV1: {
		SectorSize: ss64GiB,
		ProofSize:  192,
	},
	RegisteredPoStProof_StackedDrgWindow2KiBV1_1: {
		SectorSize: ss2KiB,
		ProofSize:  192,
	},
	RegisteredPoStProof_StackedDrgWindow8MiBV1_1: {
		SectorSize: ss8MiB,
		ProofSize:  192,
	},
	RegisteredPoStProof_StackedDrgWindow512MiBV1_1: {
		SectorSize: ss512MiB,
		ProofSize:  192,
	},
	RegisteredPoStProof_StackedDrgWindow32GiBV1_1: {
		SectorSize: ss32GiB,
		ProofSize:  192,
	},
	RegisteredPoStProof_StackedDrgWindow64GiBV1_1: {
		SectorSize: ss64GiB,
		ProofSize:  192,
	},
}

func (p RegisteredPoStProof) SectorSize() (SectorSize, error) {
	info, ok := PoStProofInfos[p]
	if !ok {
		return 0, xerrors.Errorf("unsupported proof type: %v", p)
	}
	return info.SectorSize, nil
}

// ProofSize returns the size of window post proofs for the given sector type.
func (p RegisteredPoStProof) ProofSize() (uint64, error) {
	info, ok := PoStProofInfos[p]
	if !ok {
		return 0, xerrors.Errorf("unsupported proof type: %v", p)
	}
	return info.ProofSize, nil
}

type SealRandomness Randomness
type InteractiveSealRandomness Randomness
type PoStRandomness Randomness
