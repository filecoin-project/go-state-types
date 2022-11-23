package abi

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/network"
)

// AddressValidForNetworkVersion returns true if the address is supported by the given network
// version.
//
// NOTE: It will _also_ return true if the address is "empty", because all versions support empty
// addresses in some places. I.e., it's not a version specific check.
func AddressValidForNetworkVersion(addr address.Address, nv network.Version) bool {
	// We define "undefined" addresses as "supported". The user should check for those
	// separately.
	if addr == address.Undef {
		return true
	}

	switch addr.Protocol() {
	case address.ID, address.SECP256K1, address.Actor, address.BLS:
		return true
	case address.Delegated:
		return nv >= network.Version18
	default:
		return false
	}
}
