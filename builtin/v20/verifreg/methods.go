package verifreg

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: builtin.NewMethodMeta("Constructor", *new(func(*address.Address) *abi.EmptyValue)),               // Constructor
	2: builtin.NewMethodMeta("AddVerifier", *new(func(*AddVerifierParams) *abi.EmptyValue)),             // AddVerifier
	3: builtin.NewMethodMeta("RemoveVerifier", *new(func(*address.Address) *abi.EmptyValue)),            // RemoveVerifier
	4: builtin.NewMethodMeta("AddVerifiedClient", *new(func(*AddVerifiedClientParams) *abi.EmptyValue)), // AddVerifiedClient
	builtin.MustGenerateFRCMethodNum("AddVerifiedClient"): builtin.NewMethodMeta("AddVerifiedClientExported", *new(func(*AddVerifiedClientParams) *abi.EmptyValue)), // AddVerifiedClientExported
	5: builtin.NewMethodMeta("UseBytes", nil),                                                                                         // deprecated
	6: builtin.NewMethodMeta("RestoreBytes", nil),                                                                                     // deprecated
	7: builtin.NewMethodMeta("RemoveVerifiedClientDataCap", *new(func(*RemoveDataCapParams) *RemoveDataCapReturn)),                    // RemoveVerifiedClientDataCap
	8: builtin.NewMethodMeta("RemoveExpiredAllocations", *new(func(*RemoveExpiredAllocationsParams) *RemoveExpiredAllocationsReturn)), // RemoveExpiredAllocations
	builtin.MustGenerateFRCMethodNum("RemoveExpiredAllocations"): builtin.NewMethodMeta("RemoveExpiredAllocationsExported", *new(func(*RemoveExpiredAllocationsParams) *RemoveExpiredAllocationsReturn)), // RemoveExpiredAllocationsExported
	9:  builtin.NewMethodMeta("ClaimAllocations", *new(func(*ClaimAllocationsParams) *ClaimAllocationsReturn)), // ClaimAllocations
	10: builtin.NewMethodMeta("GetClaims", *new(func(*GetClaimsParams) *GetClaimsReturn)),                      // GetClaims
	builtin.MustGenerateFRCMethodNum("GetClaims"): builtin.NewMethodMeta("GetClaimsExported", *new(func(*GetClaimsParams) *GetClaimsReturn)), // GetClaimsExported
	11: builtin.NewMethodMeta("ExtendClaimTerms", *new(func(*ExtendClaimTermsParams) *ExtendClaimTermsReturn)), // ExtendClaimTerms
	builtin.MustGenerateFRCMethodNum("ExtendClaimTerms"): builtin.NewMethodMeta("ExtendClaimTermsExported", *new(func(*ExtendClaimTermsParams) *ExtendClaimTermsReturn)), // ExtendClaimTermsExported
	12: builtin.NewMethodMeta("RemoveExpiredClaims", *new(func(*RemoveExpiredClaimsParams) *RemoveExpiredClaimsReturn)), // RemoveExpiredClaims
	builtin.MustGenerateFRCMethodNum("RemoveExpiredClaims"): builtin.NewMethodMeta("RemoveExpiredClaimsExported", *new(func(*RemoveExpiredClaimsParams) *RemoveExpiredClaimsReturn)), // RemoveExpiredClaimsExported
	builtin.MustGenerateFRCMethodNum("Receive"):             builtin.NewMethodMeta("UniversalReceiverHook", *new(func(*UniversalReceiverParams) *AllocationsResponse)),               // UniversalReceiverHook
}
