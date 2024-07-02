package verifreg

import (
	"github.com/filecoin-project/go-address"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: builtin.NewMethodMeta("Constructor", *new(func(*address.Address) *abi.EmptyValue)),                          // Constructor
	2: builtin.NewMethodMeta("AddVerifier", *new(func(*AddVerifierParams) *abi.EmptyValue)),                        // AddVerifier
	3: builtin.NewMethodMeta("RemoveVerifier", *new(func(*address.Address) *abi.EmptyValue)),                       // RemoveVerifier
	4: builtin.NewMethodMeta("AddVerifiedClient", *new(func(*AddVerifiedClientParams) *abi.EmptyValue)),            // AddVerifiedClient
	5: builtin.NewMethodMeta("UseBytes", *new(func(*UseBytesParams) *abi.EmptyValue)),                              // UseBytes
	6: builtin.NewMethodMeta("RestoreBytes", *new(func(*RestoreBytesParams) *abi.EmptyValue)),                      // RestoreBytes
	7: builtin.NewMethodMeta("RemoveVerifiedClientDataCap", *new(func(*RemoveDataCapParams) *RemoveDataCapReturn)), // RemoveVerifiedClientDataCap
}
