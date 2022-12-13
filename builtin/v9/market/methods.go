package market

import (
	"github.com/filecoin-project/go-address"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*abi.EmptyValue) *abi.EmptyValue)},                                              // Constructor
	2: {"AddBalance", *new(func(*address.Address) *abi.EmptyValue)},                                              // AddBalance
	3: {"WithdrawBalance", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)},                                  // WithdrawBalance
	4: {"PublishStorageDeals", *new(func(*PublishStorageDealsParams) *PublishStorageDealsReturn)},                // PublishStorageDeals
	5: {"VerifyDealsForActivation", *new(func(*VerifyDealsForActivationParams) *VerifyDealsForActivationReturn)}, // VerifyDealsForActivation
	6: {"ActivateDeals", *new(func(*ActivateDealsParams) *abi.EmptyValue)},                                       // ActivateDeals
	7: {"OnMinerSectorsTerminate", *new(func(*OnMinerSectorsTerminateParams) *abi.EmptyValue)},                   // OnMinerSectorsTerminate
	8: {"ComputeDataCommitment", *new(func(*ComputeDataCommitmentParams) *ComputeDataCommitmentReturn)},          // ComputeDataCommitment
	9: {"CronTick", *new(func(*abi.EmptyValue) *abi.EmptyValue)},                                                 // CronTick
}
