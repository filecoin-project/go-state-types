package market

import (
	"github.com/filecoin-project/go-address"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: builtin.NewMethodMeta("Constructor", *new(func(*abi.EmptyValue) *abi.EmptyValue)),                                              // Constructor
	2: builtin.NewMethodMeta("AddBalance", *new(func(*address.Address) *abi.EmptyValue)),                                              // AddBalance
	3: builtin.NewMethodMeta("WithdrawBalance", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)),                                  // WithdrawBalance
	4: builtin.NewMethodMeta("PublishStorageDeals", *new(func(*PublishStorageDealsParams) *PublishStorageDealsReturn)),                // PublishStorageDeals
	5: builtin.NewMethodMeta("VerifyDealsForActivation", *new(func(*VerifyDealsForActivationParams) *VerifyDealsForActivationReturn)), // VerifyDealsForActivation
	6: builtin.NewMethodMeta("ActivateDeals", *new(func(*ActivateDealsParams) *abi.EmptyValue)),                                       // ActivateDeals
	7: builtin.NewMethodMeta("OnMinerSectorsTerminate", *new(func(*OnMinerSectorsTerminateParams) *abi.EmptyValue)),                   // OnMinerSectorsTerminate
	8: builtin.NewMethodMeta("ComputeDataCommitment", *new(func(*ComputeDataCommitmentParams) *ComputeDataCommitmentReturn)),          // ComputeDataCommitment
	9: builtin.NewMethodMeta("CronTick", *new(func(*abi.EmptyValue) *abi.EmptyValue)),                                                 // CronTick
}
