package reward

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: builtin.NewMethodMeta("Constructor", *new(func(*abi.StoragePower) *abi.EmptyValue)),
	2: builtin.NewMethodMeta("AwardBlockReward", *new(func(*AwardBlockRewardParams) *abi.EmptyValue)),
	3: builtin.NewMethodMeta("ThisEpochReward", *new(func(*abi.EmptyValue) *ThisEpochRewardReturn)),
	4: builtin.NewMethodMeta("UpdateNetworkKPI", *new(func(*abi.StoragePower) *abi.EmptyValue)),
	builtin.MustGenerateFRCMethodNum("SetWeightRecords"):  builtin.NewMethodMeta("SetWeightRecordsExported", *new(func(*SetWeightRecordsParams) *abi.EmptyValue)),
	builtin.MustGenerateFRCMethodNum("StepWeightRecords"): builtin.NewMethodMeta("StepWeightRecordsExported", *new(func(*StepWeightRecordsParams) *abi.EmptyValue)),
	builtin.MustGenerateFRCMethodNum("RegisterStream"):    builtin.NewMethodMeta("RegisterStreamExported", *new(func(*RegisterStreamParams) *abi.EmptyValue)),
	builtin.MustGenerateFRCMethodNum("RemoveStream"):      builtin.NewMethodMeta("RemoveStreamExported", *new(func(*RemoveStreamParams) *abi.EmptyValue)),
	builtin.MustGenerateFRCMethodNum("SetDistribution"):   builtin.NewMethodMeta("SetDistributionExported", *new(func(*SetDistributionParams) *abi.EmptyValue)),
	builtin.MustGenerateFRCMethodNum("SetShares"):         builtin.NewMethodMeta("SetSharesExported", *new(func(*SetSharesParams) *abi.EmptyValue)),
	builtin.MustGenerateFRCMethodNum("CancelPending"):     builtin.NewMethodMeta("CancelPendingExported", *new(func(*CancelPendingParams) *abi.EmptyValue)),
	builtin.MustGenerateFRCMethodNum("Claim"):             builtin.NewMethodMeta("ClaimExported", *new(func(*ClaimParams) *ClaimReturn)),
}
