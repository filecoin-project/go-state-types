package verifreg

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
)

var Methods = []interface{}{
	1: *new(func(interface{}, *address.Address) *abi.EmptyValue),          // Constructor
	2: *new(func(interface{}, *AddVerifierParams) *abi.EmptyValue),        // AddVerifier
	3: *new(func(interface{}, *address.Address) *abi.EmptyValue),          // RemoveVerifier
	4: *new(func(interface{}, *AddVerifiedClientParams) *abi.EmptyValue),  // AddVerifiedClient
	5: *new(func(interface{}, *UseBytesParams) *abi.EmptyValue),           // UseBytes
	6: *new(func(interface{}, *RestoreBytesParams) *abi.EmptyValue),       // RestoreBytes
	7: *new(func(interface{}, *RemoveDataCapParams) *RemoveDataCapReturn), // RemoveVerifiedClientDataCap
}
