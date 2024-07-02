package verifreg

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1:  builtin.NewMethodMeta("Constructor", *new(func(*address.Address) *abi.EmptyValue)),                                             // Constructor
	2:  builtin.NewMethodMeta("AddVerifier", *new(func(*AddVerifierParams) *abi.EmptyValue)),                                           // AddVerifier
	3:  builtin.NewMethodMeta("RemoveVerifier", *new(func(*address.Address) *abi.EmptyValue)),                                          // RemoveVerifier
	4:  builtin.NewMethodMeta("AddVerifiedClient", *new(func(*AddVerifiedClientParams) *abi.EmptyValue)),                               // AddVerifiedClient
	5:  builtin.NewMethodMeta("UseBytes", nil),                                                                                         // deprecated
	6:  builtin.NewMethodMeta("RestoreBytes", nil),                                                                                     // deprecated
	7:  builtin.NewMethodMeta("RemoveVerifiedClientDataCap", *new(func(*RemoveDataCapParams) *RemoveDataCapReturn)),                    // RemoveVerifiedClientDataCap
	8:  builtin.NewMethodMeta("RemoveExpiredAllocations", *new(func(*RemoveExpiredAllocationsParams) *RemoveExpiredAllocationsReturn)), // RemoveExpiredAllocations
	9:  builtin.NewMethodMeta("ClaimAllocations", *new(func(*ClaimAllocationsParams) *ClaimAllocationsReturn)),                         // ClaimAllocations
	10: builtin.NewMethodMeta("GetClaims", *new(func(*GetClaimsParams) *GetClaimsReturn)),                                              // GetClaims
	11: builtin.NewMethodMeta("ExtendClaimTerms", *new(func(*ExtendClaimTermsParams) *ExtendClaimTermsReturn)),                         // ExtendClaimTerms
	12: builtin.NewMethodMeta("RemoveExpiredClaims", *new(func(*RemoveExpiredClaimsParams) *RemoveExpiredClaimsReturn)),                // RemoveExpiredClaims
	builtin.MustGenerateFRCMethodNum("Receive"): builtin.NewMethodMeta("UniversalReceiverHook", *new(func(*UniversalReceiverParams) *AllocationsResponse)), // UniversalReceiverHook
}
