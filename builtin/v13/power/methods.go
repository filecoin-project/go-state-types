package power

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/proof"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: builtin.NewMethodMeta("Constructor", *new(func(*abi.EmptyValue) *abi.EmptyValue)),       // Constructor
	2: builtin.NewMethodMeta("CreateMiner", *new(func(*CreateMinerParams) *CreateMinerReturn)), // CreateMiner
	builtin.MustGenerateFRCMethodNum("CreateMiner"): builtin.NewMethodMeta("CreateMinerExported", *new(func(*CreateMinerParams) *CreateMinerReturn)), // CreateMinerExported
	3: builtin.NewMethodMeta("UpdateClaimedPower", *new(func(*UpdateClaimedPowerParams) *abi.EmptyValue)),   // UpdateClaimedPower
	4: builtin.NewMethodMeta("EnrollCronEvent", *new(func(*EnrollCronEventParams) *abi.EmptyValue)),         // EnrollCronEvent
	5: builtin.NewMethodMeta("CronTick", *new(func(*abi.EmptyValue) *abi.EmptyValue)),                       // CronTick
	6: builtin.NewMethodMeta("UpdatePledgeTotal", *new(func(*abi.TokenAmount) *abi.EmptyValue)),             // UpdatePledgeTotal
	7: builtin.NewMethodMeta("OnConsensusFault", nil),                                                       // deprecated
	8: builtin.NewMethodMeta("SubmitPoRepForBulkVerify", *new(func(*proof.SealVerifyInfo) *abi.EmptyValue)), // SubmitPoRepForBulkVerify
	9: builtin.NewMethodMeta("CurrentTotalPower", *new(func(*abi.EmptyValue) *CurrentTotalPowerReturn)),     // CurrentTotalPower
	builtin.MustGenerateFRCMethodNum("NetworkRawPower"):     builtin.NewMethodMeta("NetworkRawPowerExported", *new(func(*abi.EmptyValue) *NetworkRawPowerReturn)),         // NetworkRawPowerExported
	builtin.MustGenerateFRCMethodNum("MinerRawPower"):       builtin.NewMethodMeta("MinerRawPowerExported", *new(func(*MinerRawPowerParams) *MinerRawPowerReturn)),        // MinerRawPowerExported
	builtin.MustGenerateFRCMethodNum("MinerCount"):          builtin.NewMethodMeta("MinerCountExported", *new(func(*abi.EmptyValue) *MinerCountReturn)),                   // MinerCountExported
	builtin.MustGenerateFRCMethodNum("MinerConsensusCount"): builtin.NewMethodMeta("MinerConsensusCountExported", *new(func(*abi.EmptyValue) *MinerConsensusCountReturn)), // MinerConsensusCountExported
}
