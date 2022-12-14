package verifreg

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1:  {"Constructor", *new(func(*address.Address) *abi.EmptyValue)},                                             // Constructor
	2:  {"AddVerifier", *new(func(*AddVerifierParams) *abi.EmptyValue)},                                           // AddVerifier
	3:  {"RemoveVerifier", *new(func(*address.Address) *abi.EmptyValue)},                                          // RemoveVerifier
	4:  {"AddVerifiedClient", *new(func(*AddVerifiedClientParams) *abi.EmptyValue)},                               // AddVerifiedClient
	5:  {"UseBytes", nil},                                                                                         // deprecated
	6:  {"RestoreBytes", nil},                                                                                     // deprecated
	7:  {"RemoveVerifiedClientDataCap", *new(func(*RemoveDataCapParams) *RemoveDataCapReturn)},                    // RemoveVerifiedClientDataCap
	8:  {"RemoveExpiredAllocations", *new(func(*RemoveExpiredAllocationsParams) *RemoveExpiredAllocationsReturn)}, // RemoveExpiredAllocations
	9:  {"ClaimAllocations", *new(func(*ClaimAllocationsParams) *ClaimAllocationsReturn)},                         // ClaimAllocations
	10: {"GetClaims", *new(func(*GetClaimsParams) *GetClaimsReturn)},                                              // GetClaims
	11: {"ExtendClaimTerms", *new(func(*ExtendClaimTermsParams) *ExtendClaimTermsReturn)},                         // ExtendClaimTerms
	12: {"RemoveExpiredClaims", *new(func(*RemoveExpiredClaimsParams) *RemoveExpiredClaimsReturn)},                // RemoveExpiredClaims
	builtin.MustGenerateFRCMethodNum("Receive"): {"UniversalReceiverHook", *new(func(*UniversalReceiverParams) *AllocationsResponse)}, // UniversalReceiverHook
}
