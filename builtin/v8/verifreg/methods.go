package verifreg

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[uint64]builtin.MethodMeta{
	1: {"", *new(func(*address.Address) *abi.EmptyValue)},          // Constructor
	2: {"", *new(func(*AddVerifierParams) *abi.EmptyValue)},        // AddVerifier
	3: {"", *new(func(*address.Address) *abi.EmptyValue)},          // RemoveVerifier
	4: {"", *new(func(*AddVerifiedClientParams) *abi.EmptyValue)},  // AddVerifiedClient
	5: {"", *new(func(*UseBytesParams) *abi.EmptyValue)},           // UseBytes
	6: {"", *new(func(*RestoreBytesParams) *abi.EmptyValue)},       // RestoreBytes
	7: {"", *new(func(*RemoveDataCapParams) *RemoveDataCapReturn)}, // RemoveVerifiedClientDataCap
}
