package datacap

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
)

var Methods = map[uint64]interface{}{
	1: *new(func(interface{}, *address.Address) *abi.EmptyValue), // Constructor
	2: *new(func(interface{}, *MintParams) *MintReturn),          // Mint
	3: *new(func(interface{}, *DestroyParams) *BurnReturn),       // Destroy
	// Reserved
	10: *new(func(interface{}, *abi.EmptyValue) *abi.CborString),           // Name
	11: *new(func(interface{}, *abi.EmptyValue) *abi.CborString),           // Symbol
	12: *new(func(interface{}, *abi.EmptyValue) *abi.TokenAmount),          // TotalSupply
	13: *new(func(interface{}, *address.Address) *abi.TokenAmount),         // BalanceOf
	14: *new(func(interface{}, *TransferParams) *TransferReturn),           // Transfer
	15: *new(func(interface{}, *TransferFromParams) *TransferFromReturn),   // TransferFrom
	16: *new(func(interface{}, *IncreaseAllowanceParams) *abi.TokenAmount), // IncreaseAllowance
	17: *new(func(interface{}, *DecreaseAllowanceParams) *abi.TokenAmount), // DecreaseAllowance
	18: *new(func(interface{}, *RevokeAllowanceParams) *abi.TokenAmount),   // RevokeAllowance
	19: *new(func(interface{}, *BurnParams) *BurnReturn),                   // Burn
	20: *new(func(interface{}, *BurnFromParams) *BurnFromReturn),           // BurnFrom
}
