package reward

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*abi.StoragePower) *abi.EmptyValue)},            // Constructor
	2: {"AwardBlockReward", *new(func(*AwardBlockRewardParams) *abi.EmptyValue)}, // AwardBlockReward
	3: {"ThisEpochReward", *new(func(*abi.EmptyValue) *ThisEpochRewardReturn)},   // ThisEpochReward
	4: {"UpdateNetworkKPI", *new(func(*abi.StoragePower) *abi.EmptyValue)},       // UpdateNetworkKPI
}
