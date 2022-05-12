package power

import (
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
)

// Storage miner actor constructor params are defined here so the power actor can send them to the init actor
// to instantiate miners.
// Changed since v2:
// - Seal proof type replaced with PoSt proof type
type MinerConstructorParams struct {
	OwnerAddr           addr.Address
	WorkerAddr          addr.Address
	ControlAddrs        []addr.Address
	WindowPoStProofType abi.RegisteredPoStProof
	PeerId              abi.PeerID
	Multiaddrs          []abi.Multiaddrs
}
