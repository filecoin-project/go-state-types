package reward

import (
	"github.com/filecoin-project/go-state-types/abi"
)

var Methods = []interface{}{
	1: *new(func(interface{}, *abi.StoragePower) *abi.EmptyValue),       // Constructor
	2: *new(func(interface{}, *AwardBlockRewardParams) *abi.EmptyValue), // AwardBlockReward
	3: *new(func(interface{}, *abi.EmptyValue) *ThisEpochRewardReturn),  // ThisEpochReward
	4: *new(func(interface{}, *abi.StoragePower) *abi.EmptyValue),       // UpdateNetworkKPI
}
