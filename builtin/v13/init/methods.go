package init

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*ConstructorParams) *abi.EmptyValue)}, // Constructor
	2: {"Exec", *new(func(*ExecParams) *ExecReturn)},                   // Exec
	// TODO Are we exporting Exec4
	3: {"Exec4", *new(func(*Exec4Params) *ExecReturn)}, // Exec4
}
