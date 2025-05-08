package market

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v17/miner"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: builtin.NewMethodMeta("Constructor", *new(func(*abi.EmptyValue) *abi.EmptyValue)), // Constructor
	2: builtin.NewMethodMeta("AddBalance", *new(func(*address.Address) *abi.EmptyValue)), // AddBalance
	builtin.MustGenerateFRCMethodNum("AddBalance"): builtin.NewMethodMeta("AddBalanceExported", *new(func(*address.Address) *abi.EmptyValue)), // AddBalanceExported
	3: builtin.NewMethodMeta("WithdrawBalance", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)), // WithdrawBalance
	builtin.MustGenerateFRCMethodNum("WithdrawBalance"): builtin.NewMethodMeta("WithdrawBalanceExported", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)), // WithdrawBalanceExported
	4: builtin.NewMethodMeta("PublishStorageDeals", *new(func(*PublishStorageDealsParams) *PublishStorageDealsReturn)), // PublishStorageDeals
	builtin.MustGenerateFRCMethodNum("PublishStorageDeals"): builtin.NewMethodMeta("PublishStorageDealsExported", *new(func(*PublishStorageDealsParams) *PublishStorageDealsReturn)), // PublishStorageDealsExported
	5: builtin.NewMethodMeta("VerifyDealsForActivation", *new(func(*VerifyDealsForActivationParams) *VerifyDealsForActivationReturn)), // VerifyDealsForActivation
	6: builtin.NewMethodMeta("ActivateDeals", *new(func(*ActivateDealsParams) *abi.EmptyValue)),                                       // ActivateDeals
	7: builtin.NewMethodMeta("OnMinerSectorsTerminate", *new(func(*OnMinerSectorsTerminateParams) *abi.EmptyValue)),                   // OnMinerSectorsTerminate
	8: builtin.NewMethodMeta("ComputeDataCommitment", nil),                                                                            // deprecated
	9: builtin.NewMethodMeta("CronTick", *new(func(*abi.EmptyValue) *abi.EmptyValue)),                                                 // CronTick
	builtin.MustGenerateFRCMethodNum("GetBalance"):                builtin.NewMethodMeta("GetBalanceExported", *new(func(*address.Address) *GetBalanceReturn)),                                               // GetBalanceExported
	builtin.MustGenerateFRCMethodNum("GetDealDataCommitment"):     builtin.NewMethodMeta("GetDealDataCommitmentExported", *new(func(*GetDealDataCommitmentParams) *GetDealDataCommitmentReturn)),             // GetDealDataCommitmentExported
	builtin.MustGenerateFRCMethodNum("GetDealClient"):             builtin.NewMethodMeta("GetDealClientExported", *new(func(*GetDealClientParams) *GetDealClientReturn)),                                     // GetDealClientExported
	builtin.MustGenerateFRCMethodNum("GetDealProvider"):           builtin.NewMethodMeta("GetDealProviderExported", *new(func(*GetDealProviderParams) *GetDealProviderReturn)),                               // GetDealProviderExported
	builtin.MustGenerateFRCMethodNum("GetDealLabel"):              builtin.NewMethodMeta("GetDealLabelExported", *new(func(*GetDealLabelParams) *GetDealLabelReturn)),                                        // GetDealLabelExported
	builtin.MustGenerateFRCMethodNum("GetDealTerm"):               builtin.NewMethodMeta("GetDealTermExported", *new(func(*GetDealTermParams) *GetDealTermReturn)),                                           // GetDealTermExported
	builtin.MustGenerateFRCMethodNum("GetDealTotalPrice"):         builtin.NewMethodMeta("GetDealTotalPriceExported", *new(func(*GetDealTotalPriceParams) *GetDealTotalPriceReturn)),                         // GetDealTotalPriceExported
	builtin.MustGenerateFRCMethodNum("GetDealClientCollateral"):   builtin.NewMethodMeta("GetDealClientCollateralExported", *new(func(*GetDealClientCollateralParams) *GetDealClientCollateralReturn)),       // GetDealClientCollateralExported
	builtin.MustGenerateFRCMethodNum("GetDealProviderCollateral"): builtin.NewMethodMeta("GetDealProviderCollateralExported", *new(func(*GetDealProviderCollateralParams) *GetDealProviderCollateralReturn)), // GetDealProviderCollateralExported
	builtin.MustGenerateFRCMethodNum("GetDealVerified"):           builtin.NewMethodMeta("GetDealVerifiedExported", *new(func(*GetDealVerifiedParams) *GetDealVerifiedReturn)),                               // GetDealVerifiedExported
	builtin.MustGenerateFRCMethodNum("GetDealActivation"):         builtin.NewMethodMeta("GetDealActivationExported", *new(func(*GetDealActivationParams) *GetDealActivationReturn)),                         // GetDealActivationExported
	builtin.MustGenerateFRCMethodNum("GetDealSector"):             builtin.NewMethodMeta("GetDealSectorExported", *new(func(*GetDealSectorParams) *GetDealSectorReturn)),                                     // GetDealSectorExported
	builtin.MethodSectorContentChanged:                            builtin.NewMethodMeta("SectorContentChanged", *new(func(*miner.SectorContentChangedParams) *miner.SectorContentChangedReturn)),            // SectorContentChanged
	builtin.MustGenerateFRCMethodNum("SettleDealPayments"):        builtin.NewMethodMeta("SettleDealPaymentsExported", *new(func(*SettleDealPaymentsParams) *SettleDealPaymentsReturn)),                      // SettleDealPaymentsExported
}
