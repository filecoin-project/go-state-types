package datacap

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[uint64]builtin.MethodMeta{
	1: {"", *new(func(*address.Address) *abi.EmptyValue)}, // Constructor
	2: {"", *new(func(*MintParams) *MintReturn)},          // Mint
	3: {"", *new(func(*DestroyParams) *BurnReturn)},       // Destroy
	// Reserved
	10: {"", *new(func(*abi.EmptyValue) *abi.CborString)},           // Name
	11: {"", *new(func(*abi.EmptyValue) *abi.CborString)},           // Symbol
	12: {"", *new(func(*abi.EmptyValue) *abi.TokenAmount)},          // TotalSupply
	13: {"", *new(func(*address.Address) *abi.TokenAmount)},         // BalanceOf
	14: {"", *new(func(*TransferParams) *TransferReturn)},           // Transfer
	15: {"", *new(func(*TransferFromParams) *TransferFromReturn)},   // TransferFrom
	16: {"", *new(func(*IncreaseAllowanceParams) *abi.TokenAmount)}, // IncreaseAllowance
	17: {"", *new(func(*DecreaseAllowanceParams) *abi.TokenAmount)}, // DecreaseAllowance
	18: {"", *new(func(*RevokeAllowanceParams) *abi.TokenAmount)},   // RevokeAllowance
	19: {"", *new(func(*BurnParams) *BurnReturn)},                   // Burn
	20: {"", *new(func(*BurnFromParams) *BurnFromReturn)},           // BurnFrom
}
