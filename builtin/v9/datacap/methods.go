package datacap

import (
	"github.com/filecoin-project/go-address"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*address.Address) *abi.EmptyValue)}, // Constructor
	2: {"Mint", *new(func(*MintParams) *MintReturn)},                 // Mint
	3: {"Destroy", *new(func(*DestroyParams) *BurnReturn)},           // Destroy
	// Reserved
	10: {"Name", *new(func(*abi.EmptyValue) *abi.CborString)},                        // Name
	11: {"Symbol", *new(func(*abi.EmptyValue) *abi.CborString)},                      // Symbol
	12: {"TotalSupply", *new(func(*abi.EmptyValue) *abi.TokenAmount)},                // TotalSupply
	13: {"BalanceOf", *new(func(*address.Address) *abi.TokenAmount)},                 // BalanceOf
	14: {"Transfer", *new(func(*TransferParams) *TransferReturn)},                    // Transfer
	15: {"TransferFrom", *new(func(*TransferFromParams) *TransferFromReturn)},        // TransferFrom
	16: {"IncreaseAllowance", *new(func(*IncreaseAllowanceParams) *abi.TokenAmount)}, // IncreaseAllowance
	17: {"DecreaseAllowance", *new(func(*DecreaseAllowanceParams) *abi.TokenAmount)}, // DecreaseAllowance
	18: {"RevokeAllowance", *new(func(*RevokeAllowanceParams) *abi.TokenAmount)},     // RevokeAllowance
	19: {"Burn", *new(func(*BurnParams) *BurnReturn)},                                // Burn
	20: {"BurnFrom", *new(func(*BurnFromParams) *BurnFromReturn)},                    // BurnFrom
	21: {"Allowance", *new(func(*GetAllowanceParams) *abi.TokenAmount)},              // Allowance
}
