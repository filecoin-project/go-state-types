package market

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*abi.EmptyValue) *abi.EmptyValue)}, // Constructor
	2: {"AddBalance", *new(func(*address.Address) *abi.EmptyValue)}, // AddBalance
	builtin.MustGenerateFRCMethodNum("AddBalance"): {"AddBalanceExported", *new(func(*address.Address) *abi.EmptyValue)}, // AddBalanceExported
	3: {"WithdrawBalance", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)}, // WithdrawBalance
	builtin.MustGenerateFRCMethodNum("WithdrawBalance"): {"WithdrawBalanceExported", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)}, // WithdrawBalanceExported
	4: {"PublishStorageDeals", *new(func(*PublishStorageDealsParams) *PublishStorageDealsReturn)}, // PublishStorageDeals
	builtin.MustGenerateFRCMethodNum("PublishStorageDeals"): {"PublishStorageDealsExported", *new(func(*PublishStorageDealsParams) *PublishStorageDealsReturn)}, // PublishStorageDealsExported
	5: {"VerifyDealsForActivation", *new(func(*VerifyDealsForActivationParams) *VerifyDealsForActivationReturn)}, // VerifyDealsForActivation
	6: {"ActivateDeals", *new(func(*ActivateDealsParams) *abi.EmptyValue)},                                       // ActivateDeals
	7: {"OnMinerSectorsTerminate", *new(func(*OnMinerSectorsTerminateParams) *abi.EmptyValue)},                   // OnMinerSectorsTerminate
	8: {"ComputeDataCommitment", *new(func(*ComputeDataCommitmentParams) *ComputeDataCommitmentReturn)},          // ComputeDataCommitment
	9: {"CronTick", *new(func(*abi.EmptyValue) *abi.EmptyValue)},                                                 // CronTick
	builtin.MustGenerateFRCMethodNum("GetBalance"):                {"GetBalanceExported", *new(func(*address.Address) *GetBalanceReturn)},                                               // GetBalanceExported
	builtin.MustGenerateFRCMethodNum("GetDealDataCommitment"):     {"GetDealDataCommitmentExported", *new(func(*GetDealDataCommitmentParams) *GetDealDataCommitmentReturn)},             // GetDealDataCommitmentExported
	builtin.MustGenerateFRCMethodNum("GetDealClient"):             {"GetDealClientExported", *new(func(*GetDealClientParams) *GetDealClientReturn)},                                     // GetDealClientExported
	builtin.MustGenerateFRCMethodNum("GetDealProvider"):           {"GetDealProviderExported", *new(func(*GetDealProviderParams) *GetDealProviderReturn)},                               // GetDealProviderExported
	builtin.MustGenerateFRCMethodNum("GetDealLabel"):              {"GetDealLabelExported", *new(func(*GetDealLabelParams) *GetDealLabelReturn)},                                        // GetDealLabelExported
	builtin.MustGenerateFRCMethodNum("GetDealTerm"):               {"GetDealTermExported", *new(func(*GetDealTermParams) *GetDealTermReturn)},                                           // GetDealTermExported
	builtin.MustGenerateFRCMethodNum("GetDealTotalPrice"):         {"GetDealTotalPriceExported", *new(func(*GetDealTotalPriceParams) *GetDealTotalPriceReturn)},                         // GetDealTotalPriceExported
	builtin.MustGenerateFRCMethodNum("GetDealClientCollateral"):   {"GetDealClientCollateralExported", *new(func(*GetDealClientCollateralParams) *GetDealClientCollateralReturn)},       // GetDealClientCollateralExported
	builtin.MustGenerateFRCMethodNum("GetDealProviderCollateral"): {"GetDealProviderCollateralExported", *new(func(*GetDealProviderCollateralParams) *GetDealProviderCollateralReturn)}, // GetDealProviderCollateralExported
	builtin.MustGenerateFRCMethodNum("GetDealVerified"):           {"GetDealVerifiedExported", *new(func(*GetDealVerifiedParams) *GetDealVerifiedReturn)},                               // GetDealVerifiedExported
	builtin.MustGenerateFRCMethodNum("GetDealActivation"):         {"GetDealActivationExported", *new(func(*GetDealActivationParams) *GetDealActivationReturn)},                         // GetDealActivationExported
}
