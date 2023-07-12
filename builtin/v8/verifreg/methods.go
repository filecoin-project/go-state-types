package verifreg

import (
	"github.com/filecoin-project/go-address"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*address.Address) *abi.EmptyValue)},                          // Constructor
	2: {"AddVerifier", *new(func(*AddVerifierParams) *abi.EmptyValue)},                        // AddVerifier
	3: {"RemoveVerifier", *new(func(*address.Address) *abi.EmptyValue)},                       // RemoveVerifier
	4: {"AddVerifiedClient", *new(func(*AddVerifiedClientParams) *abi.EmptyValue)},            // AddVerifiedClient
	5: {"UseBytes", *new(func(*UseBytesParams) *abi.EmptyValue)},                              // UseBytes
	6: {"RestoreBytes", *new(func(*RestoreBytesParams) *abi.EmptyValue)},                      // RestoreBytes
	7: {"RemoveVerifiedClientDataCap", *new(func(*RemoveDataCapParams) *RemoveDataCapReturn)}, // RemoveVerifiedClientDataCap
}
