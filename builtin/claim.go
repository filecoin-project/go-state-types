package builtin

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
)

type OwnedClaim struct {
	// Address of the miner.
	Address address.Address

	// Sum of raw byte power for a miner's sectors.
	RawBytePower abi.StoragePower
	// Sum of quality adjusted power for a miner's sectors.
	QualityAdjPower abi.StoragePower
}
