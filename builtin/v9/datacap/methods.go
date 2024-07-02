package datacap

import (
	"github.com/filecoin-project/go-address"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: builtin.NewMethodMeta("Constructor", *new(func(*address.Address) *abi.EmptyValue)), // Constructor
	2: builtin.NewMethodMeta("Mint", *new(func(*MintParams) *MintReturn)),                 // Mint
	3: builtin.NewMethodMeta("Destroy", *new(func(*DestroyParams) *BurnReturn)),           // Destroy
	// Reserved
	10: builtin.NewMethodMeta("Name", *new(func(*abi.EmptyValue) *abi.CborString)),                        // Name
	11: builtin.NewMethodMeta("Symbol", *new(func(*abi.EmptyValue) *abi.CborString)),                      // Symbol
	12: builtin.NewMethodMeta("TotalSupply", *new(func(*abi.EmptyValue) *abi.TokenAmount)),                // TotalSupply
	13: builtin.NewMethodMeta("BalanceOf", *new(func(*address.Address) *abi.TokenAmount)),                 // BalanceOf
	14: builtin.NewMethodMeta("Transfer", *new(func(*TransferParams) *TransferReturn)),                    // Transfer
	15: builtin.NewMethodMeta("TransferFrom", *new(func(*TransferFromParams) *TransferFromReturn)),        // TransferFrom
	16: builtin.NewMethodMeta("IncreaseAllowance", *new(func(*IncreaseAllowanceParams) *abi.TokenAmount)), // IncreaseAllowance
	17: builtin.NewMethodMeta("DecreaseAllowance", *new(func(*DecreaseAllowanceParams) *abi.TokenAmount)), // DecreaseAllowance
	18: builtin.NewMethodMeta("RevokeAllowance", *new(func(*RevokeAllowanceParams) *abi.TokenAmount)),     // RevokeAllowance
	19: builtin.NewMethodMeta("Burn", *new(func(*BurnParams) *BurnReturn)),                                // Burn
	20: builtin.NewMethodMeta("BurnFrom", *new(func(*BurnFromParams) *BurnFromReturn)),                    // BurnFrom
	21: builtin.NewMethodMeta("Allowance", *new(func(*GetAllowanceParams) *abi.TokenAmount)),              // Allowance
}
