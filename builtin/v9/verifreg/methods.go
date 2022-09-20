package verifreg

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
)

const UniversalReceiverHookMethodNum = 302636343

var Methods = map[uint64]interface{}{
	1:                              *new(func(interface{}, *address.Address) *abi.EmptyValue),         // Constructor
	2:                              *new(func(interface{}, *AddVerifierParams) *abi.EmptyValue),       // AddVerifier
	3:                              *new(func(interface{}, *address.Address) *abi.EmptyValue),         // RemoveVerifier
	4:                              *new(func(interface{}, *AddVerifiedClientParams) *abi.EmptyValue), // AddVerifiedClient
	5:                              nil,
	6:                              nil,
	7:                              *new(func(interface{}, *RemoveDataCapParams) *RemoveDataCapReturn),                       // RemoveVerifiedClientDataCap
	8:                              *new(func(interface{}, *RemoveExpiredAllocationsParams) *RemoveExpiredAllocationsReturn), // RemoveExpiredAllocations
	9:                              *new(func(interface{}, *ClaimAllocationsParams) *ClaimAllocationsReturn),                 // ClaimAllocations
	10:                             *new(func(interface{}, *GetClaimsParams) *GetClaimsReturn),                               // GetClaims
	UniversalReceiverHookMethodNum: *new(func(interface{}, *UniversalReceiverParams) *AllocationsResponse),                   // UniversalReceiverHook
}
