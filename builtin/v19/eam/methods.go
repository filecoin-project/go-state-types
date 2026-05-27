package eam

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: builtin.NewMethodMeta("Constructor", *new(func(*abi.EmptyValue) *abi.EmptyValue)),         // Constructor
	2: builtin.NewMethodMeta("Create", *new(func(*CreateParams) *CreateReturn)),                  // Create
	3: builtin.NewMethodMeta("Create2", *new(func(*Create2Params) *Create2Return)),               // Create2
	4: builtin.NewMethodMeta("CreateExternal", *new(func(*abi.CborBytes) *CreateExternalReturn)), // CreateExternal
}
