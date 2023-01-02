package evm

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	typegen "github.com/whyrusleeping/cbor-gen"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*ConstructorParams) *abi.EmptyValue)},
	2: {"InvokeContract", *new(func(bytes *abi.CborBytes) *abi.CborBytes)},
	3: {"GetBytecode", *new(func(*abi.EmptyValue) *typegen.CborCid)},
	4: {"GetStorageAt", *new(func(*GetStorageAtParams) *big.Int)},
	5: {"InvokeContractDelegate", *new(func(params *DelegateCallParams) *abi.CborBytes)},
	6: {"GetBytecodeHash", *new(func(*abi.EmptyValue) *typegen.CborCid)},
}
