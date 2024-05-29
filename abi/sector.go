package abi

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/minio/sha256-simd"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/network"
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

	RegisteredSealProof_StackedDrg2KiBV1_2_Feat_NiPoRep   = RegisteredSealProof(15)
	RegisteredSealProof_StackedDrg8MiBV1_2_Feat_NiPoRep   = RegisteredSealProof(16)
	RegisteredSealProof_StackedDrg512MiBV1_2_Feat_NiPoRep = RegisteredSealProof(17)
	RegisteredSealProof_StackedDrg32GiBV1_2_Feat_NiPoRep  = RegisteredSealProof(18)
	RegisteredSealProof_StackedDrg64GiBV1_2_Feat_NiPoRep  = RegisteredSealProof(19)
)

var Synthetic = map[RegisteredSealProof]bool{
	RegisteredSealProof_StackedDrg2KiBV1_1_Feat_SyntheticPoRep:   true,
	RegisteredSealProof_StackedDrg8MiBV1_1_Feat_SyntheticPoRep:   true,
	RegisteredSealProof_StackedDrg512MiBV1_1_Feat_SyntheticPoRep: true,
	RegisteredSealProof_StackedDrg32GiBV1_1_Feat_SyntheticPoRep:  true,
	RegisteredSealProof_StackedDrg64GiBV1_1_Feat_SyntheticPoRep:  true,
}

var NonInteractive = map[RegisteredSealProof]bool{
	RegisteredSealProof_StackedDrg2KiBV1_2_Feat_NiPoRep:   true,
	RegisteredSealProof_StackedDrg8MiBV1_2_Feat_NiPoRep:   true,
	RegisteredSealProof_StackedDrg512MiBV1_2_Feat_NiPoRep: true,
	RegisteredSealProof_StackedDrg32GiBV1_2_Feat_NiPoRep:  true,
	RegisteredSealProof_StackedDrg64GiBV1_2_Feat_NiPoRep:  true,
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
	RegisteredSealProof_StackedDrg2KiBV1_2_Feat_NiPoRep: {
		ProofSize:        192,
		SectorSize:       ss2KiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning2KiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow2KiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg2KiBV1,
	},
	RegisteredSealProof_StackedDrg8MiBV1_2_Feat_NiPoRep: {
		ProofSize:        192,
		SectorSize:       ss8MiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning8MiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow8MiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg8MiBV1,
	},
	RegisteredSealProof_StackedDrg512MiBV1_2_Feat_NiPoRep: {
		ProofSize:        192,
		SectorSize:       ss512MiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning512MiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow512MiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg512MiBV1,
	},
	RegisteredSealProof_StackedDrg32GiBV1_2_Feat_NiPoRep: {
		ProofSize:        1920,
		SectorSize:       ss32GiB,
		WinningPoStProof: RegisteredPoStProof_StackedDrgWinning32GiBV1,
		WindowPoStProof:  RegisteredPoStProof_StackedDrgWindow32GiBV1,
		UpdateProof:      RegisteredUpdateProof_StackedDrg32GiBV1,
	},
	RegisteredSealProof_StackedDrg64GiBV1_2_Feat_NiPoRep: {
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

// mapping to porep_id values
// https://github.com/filecoin-project/rust-filecoin-proofs-api/blob/9b580c2791f028b7ce3aaef0cf9d68956c50170d/src/registry.rs#L32
var registeredProofIds = map[RegisteredSealProof]uint64{
	RegisteredSealProof_StackedDrg2KiBV1:   0,
	RegisteredSealProof_StackedDrg8MiBV1:   1,
	RegisteredSealProof_StackedDrg512MiBV1: 2,
	RegisteredSealProof_StackedDrg32GiBV1:  3,
	RegisteredSealProof_StackedDrg64GiBV1:  4,

	RegisteredSealProof_StackedDrg2KiBV1_1:   5,
	RegisteredSealProof_StackedDrg8MiBV1_1:   6,
	RegisteredSealProof_StackedDrg512MiBV1_1: 7,
	RegisteredSealProof_StackedDrg32GiBV1_1:  8,
	RegisteredSealProof_StackedDrg64GiBV1_1:  9,

	RegisteredSealProof_StackedDrg2KiBV1_1_Feat_SyntheticPoRep:   10,
	RegisteredSealProof_StackedDrg8MiBV1_1_Feat_SyntheticPoRep:   11,
	RegisteredSealProof_StackedDrg512MiBV1_1_Feat_SyntheticPoRep: 12,
	RegisteredSealProof_StackedDrg32GiBV1_1_Feat_SyntheticPoRep:  13,
	RegisteredSealProof_StackedDrg64GiBV1_1_Feat_SyntheticPoRep:  14,
}

func (p RegisteredSealProof) porepNonce() uint64 {
	// https://github.com/filecoin-project/rust-filecoin-proofs-api/blob/9b580c2791f028b7ce3aaef0cf9d68956c50170d/src/registry.rs#L166
	return 0
}

// PoRepID produces the porep_id for this RegisteredSealProof. Mainly used for
// computing replica_id.
func (p RegisteredSealProof) PoRepID() ([32]byte, error) {
	// https://github.com/filecoin-project/rust-filecoin-proofs-api/blob/9b580c2791f028b7ce3aaef0cf9d68956c50170d/src/registry.rs#L174

	var id [32]byte

	proofId, ok := registeredProofIds[p]
	if !ok {
		return id, xerrors.Errorf("unsupported proof type: %v", p)
	}

	// Convert the proofId and nonce to little endian byte slices
	proofIdBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(proofIdBytes, proofId)

	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, p.porepNonce())

	// Copy these byte slices into the PoRep ID
	copy(id[0:8], proofIdBytes)
	copy(id[8:16], nonceBytes)

	return id, nil
}

// ReplicaId produces the replica_id for this RegisteredSealProof. This is used
// as the main input for computing SDR
func (p RegisteredSealProof) ReplicaId(prover ActorID, sector SectorNumber, ticket []byte, commd []byte) ([32]byte, error) {
	// https://github.com/filecoin-project/rust-fil-proofs/blob/5b46d4ac88e19003416bb110e2b2871523cc2892/storage-proofs-porep/src/stacked/vanilla/params.rs#L758-L775

	pi, err := MakeProverID(prover)
	if err != nil {
		return [32]byte{}, err
	}
	porepID, err := p.PoRepID()
	if err != nil {
		return [32]byte{}, err
	}

	if len(ticket) != 32 {
		return [32]byte{}, xerrors.Errorf("invalid ticket length %d", len(ticket))
	}
	if len(commd) != 32 {
		return [32]byte{}, xerrors.Errorf("invalid commd length %d", len(commd))
	}

	var sectorID [8]byte
	binary.BigEndian.PutUint64(sectorID[:], uint64(sector))

	s := sha256.New()

	// sha256 writes never error
	_, _ = s.Write(pi[:])
	_, _ = s.Write(sectorID[:])
	_, _ = s.Write(ticket)
	_, _ = s.Write(commd)
	_, _ = s.Write(porepID[:])

	return bytesIntoFr32Safe(s.Sum(nil)), nil
}

func (p RegisteredSealProof) IsSynthetic() bool {
	_, ok := Synthetic[p]
	return ok
}

func (p RegisteredSealProof) IsNonInteractive() bool {
	_, ok := NonInteractive[p]
	return ok
}

type ProverID [32]byte

// ProverID returns a 32 byte proverID used when computing ReplicaID
func MakeProverID(e ActorID) (ProverID, error) {
	maddr, err := address.NewIDAddress(uint64(e))
	if err != nil {
		return ProverID{}, xerrors.Errorf("failed to convert ActorID to prover id ([32]byte): %w", err)
	}

	var proverID ProverID
	copy(proverID[:], maddr.Payload())
	return proverID, nil
}

func bytesIntoFr32Safe(in []byte) [32]byte {
	var out [32]byte
	copy(out[:], in)

	out[31] &= 0b0011_1111

	return out
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
