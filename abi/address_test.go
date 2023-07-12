package abi

import (
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/stretchr/testify/require"
)

func TestAddressValidForNetworkVersion(t *testing.T) {
	id, _ := address.NewIDAddress(1)
	bls, _ := address.NewBLSAddress(make([]byte, address.BlsPublicKeyBytes))
	secp, _ := address.NewSecp256k1Address(make([]byte, address.PayloadHashLength))
	actor, _ := address.NewActorAddress(make([]byte, address.PayloadHashLength))
	for _, addr := range []address.Address{id, bls, secp, actor} {
		require.True(t, AddressValidForNetworkVersion(addr, network.Version17))
	}
}
