package datacap

import (
	"github.com/filecoin-project/go-address"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1:                                        builtin.NewMethodMeta("Constructor", *new(func(*address.Address) *abi.EmptyValue)), // Constructor
	builtin.MustGenerateFRCMethodNum("Mint"): builtin.NewMethodMeta("MintExported", *new(func(*MintParams) *MintReturn)),         // MintExported
	builtin.MustGenerateFRCMethodNum("Destroy"):           builtin.NewMethodMeta("DestroyExported", *new(func(*DestroyParams) *BurnReturn)),                          // DestroyExported
	builtin.MustGenerateFRCMethodNum("Name"):              builtin.NewMethodMeta("NameExported", *new(func(*abi.EmptyValue) *abi.CborString)),                        // NameExported
	builtin.MustGenerateFRCMethodNum("Symbol"):            builtin.NewMethodMeta("SymbolExported", *new(func(*abi.EmptyValue) *abi.CborString)),                      // SymbolExported
	builtin.MustGenerateFRCMethodNum("TotalSupply"):       builtin.NewMethodMeta("TotalSupplyExported", *new(func(*abi.EmptyValue) *abi.TokenAmount)),                // TotalSupplyExported
	builtin.MustGenerateFRCMethodNum("Balance"):           builtin.NewMethodMeta("BalanceExported", *new(func(*address.Address) *abi.TokenAmount)),                   // BalanceExported
	builtin.MustGenerateFRCMethodNum("Transfer"):          builtin.NewMethodMeta("TransferExported", *new(func(*TransferParams) *TransferReturn)),                    // TransferExported
	builtin.MustGenerateFRCMethodNum("TransferFrom"):      builtin.NewMethodMeta("TransferFromExported", *new(func(*TransferFromParams) *TransferFromReturn)),        // TransferFromExported
	builtin.MustGenerateFRCMethodNum("IncreaseAllowance"): builtin.NewMethodMeta("IncreaseAllowanceExported", *new(func(*IncreaseAllowanceParams) *abi.TokenAmount)), // IncreaseAllowanceExported
	builtin.MustGenerateFRCMethodNum("DecreaseAllowance"): builtin.NewMethodMeta("DecreaseAllowanceExported", *new(func(*DecreaseAllowanceParams) *abi.TokenAmount)), // DecreaseAllowanceExported
	builtin.MustGenerateFRCMethodNum("RevokeAllowance"):   builtin.NewMethodMeta("RevokeAllowanceExported", *new(func(*RevokeAllowanceParams) *abi.TokenAmount)),     // RevokeAllowanceExported
	builtin.MustGenerateFRCMethodNum("Burn"):              builtin.NewMethodMeta("BurnExported", *new(func(*BurnParams) *BurnReturn)),                                // BurnExported
	builtin.MustGenerateFRCMethodNum("BurnFrom"):          builtin.NewMethodMeta("BurnFromExported", *new(func(*BurnFromParams) *BurnFromReturn)),                    // BurnFromExported
	builtin.MustGenerateFRCMethodNum("Allowance"):         builtin.NewMethodMeta("AllowanceExported", *new(func(*GetAllowanceParams) *abi.TokenAmount)),              // AllowanceExported
	builtin.MustGenerateFRCMethodNum("Granularity"):       builtin.NewMethodMeta("GranularityExported", *new(func(value *abi.EmptyValue) *GranularityReturn)),        // GranularityExported
}
