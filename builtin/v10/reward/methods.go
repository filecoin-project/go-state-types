package reward

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: builtin.NewMethodMeta("Constructor", *new(func(*abi.StoragePower) *abi.EmptyValue)),            // Constructor
	2: builtin.NewMethodMeta("AwardBlockReward", *new(func(*AwardBlockRewardParams) *abi.EmptyValue)), // AwardBlockReward
	3: builtin.NewMethodMeta("ThisEpochReward", *new(func(*abi.EmptyValue) *ThisEpochRewardReturn)),   // ThisEpochReward
	4: builtin.NewMethodMeta("UpdateNetworkKPI", *new(func(*abi.StoragePower) *abi.EmptyValue)),       // UpdateNetworkKPI
}
