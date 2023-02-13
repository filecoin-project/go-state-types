package evm

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*ConstructorParams) *abi.EmptyValue)},
	2: {"Resurrect", *new(func(*ResurrectParams) *abi.EmptyValue)},
	3: {"GetBytecode", *new(func(*abi.EmptyValue) *GetBytecodeReturn)},
	4: {"GetBytecodeHash", *new(func(*abi.EmptyValue) *abi.CborBytes)},
	5: {"GetStorageAt", *new(func(*GetStorageAtParams) *abi.CborBytes)},
	6: {"InvokeContractDelegate", *new(func(params *DelegateCallParams) *abi.CborBytes)},
	builtin.MustGenerateFRCMethodNum("InvokeEVM"): {"InvokeContract", *new(func(bytes *abi.CborBytes) *abi.CborBytes)},
}
