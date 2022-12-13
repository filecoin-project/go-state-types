package power

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/proof"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*abi.EmptyValue) *abi.EmptyValue)},                    // Constructor
	2: {"CreateMiner", *new(func(*CreateMinerParams) *CreateMinerReturn)},              // CreateMiner
	3: {"UpdateClaimedPower", *new(func(*UpdateClaimedPowerParams) *abi.EmptyValue)},   // UpdateClaimedPower
	4: {"EnrollCronEvent", *new(func(*EnrollCronEventParams) *abi.EmptyValue)},         // EnrollCronEvent
	5: {"CronTick", *new(func(*abi.EmptyValue) *abi.EmptyValue)},                       // CronTick
	6: {"UpdatePledgeTotal", *new(func(*abi.TokenAmount) *abi.EmptyValue)},             // UpdatePledgeTotal
	7: {"OnConsensusFault", nil},                                                       // deprecated
	8: {"SubmitPoRepForBulkVerify", *new(func(*proof.SealVerifyInfo) *abi.EmptyValue)}, // SubmitPoRepForBulkVerify
	9: {"CurrentTotalPower", *new(func(*abi.EmptyValue) *CurrentTotalPowerReturn)},     // CurrentTotalPower
}
