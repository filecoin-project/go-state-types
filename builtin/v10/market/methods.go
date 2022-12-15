package market

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/frc0042"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*abi.EmptyValue) *abi.EmptyValue)}, // Constructor
	2: {"AddBalance", *new(func(*address.Address) *abi.EmptyValue)}, // AddBalance
	frc0042.GenerateExportedMethodNum("AddBalance"): {"AddBalanceExported", *new(func(*address.Address) *abi.EmptyValue)}, // AddBalanceExported
	3: {"WithdrawBalance", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)}, // WithdrawBalance
	frc0042.GenerateExportedMethodNum("WithdrawBalance"): {"WithdrawBalanceExported", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)}, // WithdrawBalanceExported
	4: {"PublishStorageDeals", *new(func(*PublishStorageDealsParams) *PublishStorageDealsReturn)}, // PublishStorageDeals
	frc0042.GenerateExportedMethodNum("PublishStorageDeals"): {"PublishStorageDealsExported", *new(func(*PublishStorageDealsParams) *PublishStorageDealsReturn)}, // PublishStorageDealsExported
	5: {"VerifyDealsForActivation", *new(func(*VerifyDealsForActivationParams) *VerifyDealsForActivationReturn)}, // VerifyDealsForActivation
	6: {"ActivateDeals", *new(func(*ActivateDealsParams) *abi.EmptyValue)},                                       // ActivateDeals
	7: {"OnMinerSectorsTerminate", *new(func(*OnMinerSectorsTerminateParams) *abi.EmptyValue)},                   // OnMinerSectorsTerminate
	8: {"ComputeDataCommitment", *new(func(*ComputeDataCommitmentParams) *ComputeDataCommitmentReturn)},          // ComputeDataCommitment
	9: {"CronTick", *new(func(*abi.EmptyValue) *abi.EmptyValue)},                                                 // CronTick
	frc0042.GenerateExportedMethodNum("GetBalance"):                {"GetBalanceExported", *new(func(*address.Address) *GetBalanceReturn)},                                               // GetBalanceExported
	frc0042.GenerateExportedMethodNum("GetDealDataCommitment"):     {"GetDealDataCommitmentExported", *new(func(*GetDealDataCommitmentParams) *GetDealDataCommitmentReturn)},             // GetDealDataCommitmentExported
	frc0042.GenerateExportedMethodNum("GetDealClient"):             {"GetDealClientExported", *new(func(*GetDealClientParams) *GetDealClientReturn)},                                     // GetDealClientExported
	frc0042.GenerateExportedMethodNum("GetDealProvider"):           {"GetDealProviderExported", *new(func(*GetDealProviderParams) *GetDealProviderReturn)},                               // GetDealProviderExported
	frc0042.GenerateExportedMethodNum("GetDealLabel"):              {"GetDealLabelExported", *new(func(*GetDealLabelParams) *GetDealLabelReturn)},                                        // GetDealLabelExported
	frc0042.GenerateExportedMethodNum("GetDealTerm"):               {"GetDealTermExported", *new(func(*GetDealTermParams) *GetDealTermReturn)},                                           // GetDealTermExported
	frc0042.GenerateExportedMethodNum("GetDealTotalPrice"):         {"GetDealTotalPriceExported", *new(func(*GetDealTotalPriceParams) *GetDealTotalPriceReturn)},                         // GetDealTotalPriceExported
	frc0042.GenerateExportedMethodNum("GetDealClientCollateral"):   {"GetDealClientCollateralExported", *new(func(*GetDealClientCollateralParams) *GetDealClientCollateralReturn)},       // GetDealClientCollateralExported
	frc0042.GenerateExportedMethodNum("GetDealProviderCollateral"): {"GetDealProviderCollateralExported", *new(func(*GetDealProviderCollateralParams) *GetDealProviderCollateralReturn)}, // GetDealProviderCollateralExported
	frc0042.GenerateExportedMethodNum("GetDealVerified"):           {"GetDealVerifiedExported", *new(func(*GetDealVerifiedParams) *GetDealVerifiedReturn)},                               // GetDealVerifiedExported
	frc0042.GenerateExportedMethodNum("GetDealActivation"):         {"GetDealActivationExported", *new(func(*GetDealActivationParams) *GetDealActivationReturn)},                         // GetDealActivationExported
}
