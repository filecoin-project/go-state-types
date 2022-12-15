package market

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*abi.EmptyValue) *abi.EmptyValue)}, // Constructor
	2: {"AddBalance", *new(func(*address.Address) *abi.EmptyValue)}, // AddBalance
	builtin.MustGenerateExportedMethodNum("AddBalance"): {"AddBalanceExported", *new(func(*address.Address) *abi.EmptyValue)}, // AddBalanceExported
	3: {"WithdrawBalance", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)}, // WithdrawBalance
	builtin.MustGenerateExportedMethodNum("WithdrawBalance"): {"WithdrawBalanceExported", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)}, // WithdrawBalanceExported
	4: {"PublishStorageDeals", *new(func(*PublishStorageDealsParams) *PublishStorageDealsReturn)}, // PublishStorageDeals
	builtin.MustGenerateExportedMethodNum("PublishStorageDeals"): {"PublishStorageDealsExported", *new(func(*PublishStorageDealsParams) *PublishStorageDealsReturn)}, // PublishStorageDealsExported
	5: {"VerifyDealsForActivation", *new(func(*VerifyDealsForActivationParams) *VerifyDealsForActivationReturn)}, // VerifyDealsForActivation
	6: {"ActivateDeals", *new(func(*ActivateDealsParams) *abi.EmptyValue)},                                       // ActivateDeals
	7: {"OnMinerSectorsTerminate", *new(func(*OnMinerSectorsTerminateParams) *abi.EmptyValue)},                   // OnMinerSectorsTerminate
	8: {"ComputeDataCommitment", *new(func(*ComputeDataCommitmentParams) *ComputeDataCommitmentReturn)},          // ComputeDataCommitment
	9: {"CronTick", *new(func(*abi.EmptyValue) *abi.EmptyValue)},                                                 // CronTick
	builtin.MustGenerateExportedMethodNum("GetBalance"):                {"GetBalanceExported", *new(func(*address.Address) *GetBalanceReturn)},                                               // GetBalanceExported
	builtin.MustGenerateExportedMethodNum("GetDealDataCommitment"):     {"GetDealDataCommitmentExported", *new(func(*GetDealDataCommitmentParams) *GetDealDataCommitmentReturn)},             // GetDealDataCommitmentExported
	builtin.MustGenerateExportedMethodNum("GetDealClient"):             {"GetDealClientExported", *new(func(*GetDealClientParams) *GetDealClientReturn)},                                     // GetDealClientExported
	builtin.MustGenerateExportedMethodNum("GetDealProvider"):           {"GetDealProviderExported", *new(func(*GetDealProviderParams) *GetDealProviderReturn)},                               // GetDealProviderExported
	builtin.MustGenerateExportedMethodNum("GetDealLabel"):              {"GetDealLabelExported", *new(func(*GetDealLabelParams) *GetDealLabelReturn)},                                        // GetDealLabelExported
	builtin.MustGenerateExportedMethodNum("GetDealTerm"):               {"GetDealTermExported", *new(func(*GetDealTermParams) *GetDealTermReturn)},                                           // GetDealTermExported
	builtin.MustGenerateExportedMethodNum("GetDealTotalPrice"):         {"GetDealTotalPriceExported", *new(func(*GetDealTotalPriceParams) *GetDealTotalPriceReturn)},                         // GetDealTotalPriceExported
	builtin.MustGenerateExportedMethodNum("GetDealClientCollateral"):   {"GetDealClientCollateralExported", *new(func(*GetDealClientCollateralParams) *GetDealClientCollateralReturn)},       // GetDealClientCollateralExported
	builtin.MustGenerateExportedMethodNum("GetDealProviderCollateral"): {"GetDealProviderCollateralExported", *new(func(*GetDealProviderCollateralParams) *GetDealProviderCollateralReturn)}, // GetDealProviderCollateralExported
	builtin.MustGenerateExportedMethodNum("GetDealVerified"):           {"GetDealVerifiedExported", *new(func(*GetDealVerifiedParams) *GetDealVerifiedReturn)},                               // GetDealVerifiedExported
	builtin.MustGenerateExportedMethodNum("GetDealActivation"):         {"GetDealActivationExported", *new(func(*GetDealActivationParams) *GetDealActivationReturn)},                         // GetDealActivationExported
}
