package evm

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: builtin.NewMethodMeta("Constructor", *new(func(*ConstructorParams) *abi.EmptyValue)),
	2: builtin.NewMethodMeta("Resurrect", *new(func(*ResurrectParams) *abi.EmptyValue)),
	3: builtin.NewMethodMeta("GetBytecode", *new(func(*abi.EmptyValue) *GetBytecodeReturn)),
	4: builtin.NewMethodMeta("GetBytecodeHash", *new(func(*abi.EmptyValue) *abi.CborBytes)),
	5: builtin.NewMethodMeta("GetStorageAt", *new(func(*GetStorageAtParams) *abi.CborBytes)),
	6: builtin.NewMethodMeta("InvokeContractDelegate", *new(func(params *DelegateCallParams) *abi.CborBytes)),
	builtin.MustGenerateFRCMethodNum("InvokeEVM"): builtin.NewMethodMeta("InvokeContract", *new(func(bytes *abi.CborBytes) *abi.CborBytes)),
}
