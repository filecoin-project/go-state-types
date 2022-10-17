package market

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[uint64]builtin.MethodMeta{
	1: {"", *new(func(*abi.EmptyValue) *abi.EmptyValue)},                                 // Constructor
	2: {"", *new(func(*address.Address) *abi.EmptyValue)},                                // AddBalance
	3: {"", *new(func(*WithdrawBalanceParams) *abi.TokenAmount)},                         // WithdrawBalance
	4: {"", *new(func(*PublishStorageDealsParams) *PublishStorageDealsReturn)},           // PublishStorageDeals
	5: {"", *new(func(*VerifyDealsForActivationParams) *VerifyDealsForActivationReturn)}, // VerifyDealsForActivation
	6: {"", *new(func(*ActivateDealsParams) *abi.EmptyValue)},                            // ActivateDeals
	7: {"", *new(func(*OnMinerSectorsTerminateParams) *abi.EmptyValue)},                  // OnMinerSectorsTerminate
	8: {"", *new(func(*ComputeDataCommitmentParams) *ComputeDataCommitmentReturn)},       // ComputeDataCommitment
	9: {"", *new(func(*abi.EmptyValue) *abi.EmptyValue)},                                 // CronTick
}
