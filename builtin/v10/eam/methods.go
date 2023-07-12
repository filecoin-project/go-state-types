package eam

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*abi.EmptyValue) *abi.EmptyValue)},  // Constructor
	2: {"Create", *new(func(*CreateParams) *CreateReturn)},           // Create
	3: {"Create2", *new(func(*Create2Params) *Create2Return)},        // Create2
	4: {"CreateExternal", *new(func(*abi.CborBytes) *CreateExternalReturn)}, // CreateExternal
}
