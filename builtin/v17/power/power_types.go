package power

import (
	cbg "github.com/whyrusleeping/cbor-gen"

	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v17/util/smoothing"
)

// Storage miner actor constructor params are defined here so the power actor can send them to the init actor
// to instantiate miners.
type MinerConstructorParams struct {
	OwnerAddr           addr.Address
	WorkerAddr          addr.Address
	ControlAddrs        []addr.Address
	WindowPoStProofType abi.RegisteredPoStProof
	PeerId              abi.PeerID
	Multiaddrs          []abi.Multiaddrs
}

type CreateMinerParams struct {
	Owner               addr.Address
	Worker              addr.Address
	WindowPoStProofType abi.RegisteredPoStProof
	Peer                abi.PeerID
	Multiaddrs          []abi.Multiaddrs
}

type CreateMinerReturn struct {
	IDAddress     addr.Address // The canonical ID-based address for the actor.
	RobustAddress addr.Address // A more expensive but re-org-safe address for the newly created actor.
}

type UpdateClaimedPowerParams struct {
	RawByteDelta         abi.StoragePower
	QualityAdjustedDelta abi.StoragePower
}

type EnrollCronEventParams struct {
	EventEpoch abi.ChainEpoch
	Payload    []byte
}

type CurrentTotalPowerReturn struct {
	RawBytePower            abi.StoragePower
	QualityAdjPower         abi.StoragePower
	PledgeCollateral        abi.TokenAmount
	QualityAdjPowerSmoothed smoothing.FilterEstimate
	RampStartEpoch          int64
	RampDurationEpochs      uint64
}

type NetworkRawPowerReturn = abi.StoragePower

type MinerRawPowerParams = cbg.CborInt // abi.ActorID

type MinerRawPowerReturn struct {
	RawBytePower          abi.StoragePower
	MeetsConsensusMinimum bool
}

type MinerCountReturn = cbg.CborInt

type MinerConsensusCountReturn = cbg.CborInt

type MinerPowerParams = cbg.CborInt // abi.ActorID

type MinerPowerReturn struct {
	RawBytePower    abi.StoragePower
	QualityAdjPower abi.StoragePower
}
