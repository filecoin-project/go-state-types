package init

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/frc0042"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1:                                        {"Constructor", *new(func(*ConstructorParams) *abi.EmptyValue)}, // Constructor
	2:                                        {"Exec", *new(func(*ExecParams) *ExecReturn)},                   // Exec
	frc0042.GenerateExportedMethodNum("Exec"): {"ExecExported", *new(func(*ExecParams) *ExecReturn)}, // ExecExported
	// TODO is this still needed? Are we exporting it?
	3: {"Exec4", *new(func(*Exec4Params) *ExecReturn)}, // Exec4
}
