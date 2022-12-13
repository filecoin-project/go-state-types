package datacap

import (
	"github.com/filecoin-project/go-address"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1:                                        {"Constructor", *new(func(*address.Address) *abi.EmptyValue)}, // Constructor
	builtin.MustGenerateFRCMethodNum("Mint"): {"MintExported", *new(func(*MintParams) *MintReturn)},         // MintExported
	builtin.MustGenerateFRCMethodNum("Destroy"):           {"DestroyExported", *new(func(*DestroyParams) *BurnReturn)},                          // DestroyExported
	builtin.MustGenerateFRCMethodNum("Name"):              {"NameExported", *new(func(*abi.EmptyValue) *abi.CborString)},                        // NameExported
	builtin.MustGenerateFRCMethodNum("Symbol"):            {"SymbolExported", *new(func(*abi.EmptyValue) *abi.CborString)},                      // SymbolExported
	builtin.MustGenerateFRCMethodNum("TotalSupply"):       {"TotalSupplyExported", *new(func(*abi.EmptyValue) *abi.TokenAmount)},                // TotalSupplyExported
	builtin.MustGenerateFRCMethodNum("Balance"):           {"BalanceExported", *new(func(*address.Address) *abi.TokenAmount)},                   // BalanceExported
	builtin.MustGenerateFRCMethodNum("Transfer"):          {"TransferExported", *new(func(*TransferParams) *TransferReturn)},                    // TransferExported
	builtin.MustGenerateFRCMethodNum("TransferFrom"):      {"TransferFromExported", *new(func(*TransferFromParams) *TransferFromReturn)},        // TransferFromExported
	builtin.MustGenerateFRCMethodNum("IncreaseAllowance"): {"IncreaseAllowanceExported", *new(func(*IncreaseAllowanceParams) *abi.TokenAmount)}, // IncreaseAllowanceExported
	builtin.MustGenerateFRCMethodNum("DecreaseAllowance"): {"DecreaseAllowanceExported", *new(func(*DecreaseAllowanceParams) *abi.TokenAmount)}, // DecreaseAllowanceExported
	builtin.MustGenerateFRCMethodNum("RevokeAllowance"):   {"RevokeAllowanceExported", *new(func(*RevokeAllowanceParams) *abi.TokenAmount)},     // RevokeAllowanceExported
	builtin.MustGenerateFRCMethodNum("Burn"):              {"BurnExported", *new(func(*BurnParams) *BurnReturn)},                                // BurnExported
	builtin.MustGenerateFRCMethodNum("BurnFrom"):          {"BurnFromExported", *new(func(*BurnFromParams) *BurnFromReturn)},                    // BurnFromExported
	builtin.MustGenerateFRCMethodNum("Allowance"):         {"AllowanceExported", *new(func(*GetAllowanceParams) *abi.TokenAmount)},              // AllowanceExported
	builtin.MustGenerateFRCMethodNum("Granularity"):       {"GranularityExported", *new(func(value *abi.EmptyValue) *GranularityReturn)},        // GranularityExported
}
