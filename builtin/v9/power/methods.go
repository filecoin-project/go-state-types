package power

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/proof"
)

var Methods = map[uint64]builtin.MethodMeta{
	1: {"", *new(func(*abi.EmptyValue) *abi.EmptyValue)},           // Constructor
	2: {"", *new(func(*CreateMinerParams) *CreateMinerReturn)},     // CreateMiner
	3: {"", *new(func(*UpdateClaimedPowerParams) *abi.EmptyValue)}, // UpdateClaimedPower
	4: {"", *new(func(*EnrollCronEventParams) *abi.EmptyValue)},    // EnrollCronEvent
	5: {"", *new(func(*abi.EmptyValue) *abi.EmptyValue)},           // CronTick
	6: {"", *new(func(*abi.TokenAmount) *abi.EmptyValue)},          // UpdatePledgeTotal
	7: {"deprecated", nil},
	8: {"", *new(func(*proof.SealVerifyInfo) *abi.EmptyValue)},    // SubmitPoRepForBulkVerify
	9: {"", *new(func(*abi.EmptyValue) *CurrentTotalPowerReturn)}, // CurrentTotalPower
}
