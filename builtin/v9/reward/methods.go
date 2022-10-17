package reward

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[uint64]builtin.MethodMeta{
	1: {"", *new(func(*abi.StoragePower) *abi.EmptyValue)}, // Constructor
	2: {"", *new(func(*AwardBlockRewardParams) *abi.EmptyValue)},            // AwardBlockReward
	3: {"", *new(func(*abi.EmptyValue) *ThisEpochRewardReturn)},             // ThisEpochReward
	4: {"", *new(func(*abi.StoragePower) *abi.EmptyValue)},                  // UpdateNetworkKPI
}
