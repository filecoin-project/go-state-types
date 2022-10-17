package verifreg

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[uint64]builtin.MethodMeta{
	1:  {"", *new(func(*address.Address) *abi.EmptyValue)},         // Constructor
	2:  {"", *new(func(*AddVerifierParams) *abi.EmptyValue)},       // AddVerifier
	3:  {"", *new(func(*address.Address) *abi.EmptyValue)},         // RemoveVerifier
	4:  {"", *new(func(*AddVerifiedClientParams) *abi.EmptyValue)}, // AddVerifiedClient
	5:  {"deprecated", nil},
	6:  {"deprecated", nil},
	7:  {"", *new(func(*RemoveDataCapParams) *RemoveDataCapReturn)},                       // RemoveVerifiedClientDataCap
	8:  {"", *new(func(*RemoveExpiredAllocationsParams) *RemoveExpiredAllocationsReturn)}, // RemoveExpiredAllocations
	9:  {"", *new(func(*ClaimAllocationsParams) *ClaimAllocationsReturn)},                 // ClaimAllocations
	10: {"", *new(func(*GetClaimsParams) *GetClaimsReturn)},                               // GetClaims
	11: {"", *new(func(*ExtendClaimTermsParams) *ExtendClaimTermsReturn)},                 // ExtendClaimTerms
	12: {"", *new(func(*RemoveExpiredClaimsParams) *RemoveExpiredClaimsReturn)},           // RemoveExpiredClaims
	uint64(builtin.UniversalReceiverHookMethodNum): {"", *new(func(*UniversalReceiverParams) *AllocationsResponse)}, // UniversalReceiverHook
}
