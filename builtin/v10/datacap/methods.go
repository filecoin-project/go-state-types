package datacap

import (
	"github.com/filecoin-project/go-address"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1:                                        {"Constructor", *new(func(*address.Address) *abi.EmptyValue)}, // Constructor
	builtin.MustGenerateExportedMethodNum("Mint"):              {"MintExported", *new(func(*MintParams) *MintReturn)},                                // MintExported
	builtin.MustGenerateExportedMethodNum("Destroy"):           {"DestroyExported", *new(func(*DestroyParams) *BurnReturn)},                          // DestroyExported
	builtin.MustGenerateExportedMethodNum("Name"):              {"NameExported", *new(func(*abi.EmptyValue) *abi.CborString)},                        // NameExported
	builtin.MustGenerateExportedMethodNum("Symbol"):            {"SymbolExported", *new(func(*abi.EmptyValue) *abi.CborString)},                      // SymbolExported
	builtin.MustGenerateExportedMethodNum("TotalSupply"):       {"TotalSupplyExported", *new(func(*abi.EmptyValue) *abi.TokenAmount)},                // TotalSupplyExported
	builtin.MustGenerateExportedMethodNum("Balance"):           {"BalanceExported", *new(func(*address.Address) *abi.TokenAmount)},                   // BalanceExported
	builtin.MustGenerateExportedMethodNum("Transfer"):          {"TransferExported", *new(func(*TransferParams) *TransferReturn)},                    // TransferExported
	builtin.MustGenerateExportedMethodNum("TransferFrom"):      {"TransferFromExported", *new(func(*TransferFromParams) *TransferFromReturn)},        // TransferFromExported
	builtin.MustGenerateExportedMethodNum("IncreaseAllowance"): {"IncreaseAllowanceExported", *new(func(*IncreaseAllowanceParams) *abi.TokenAmount)}, // IncreaseAllowanceExported
	builtin.MustGenerateExportedMethodNum("DecreaseAllowance"): {"DecreaseAllowanceExported", *new(func(*DecreaseAllowanceParams) *abi.TokenAmount)}, // DecreaseAllowanceExported
	builtin.MustGenerateExportedMethodNum("RevokeAllowance"):   {"RevokeAllowanceExported", *new(func(*RevokeAllowanceParams) *abi.TokenAmount)},     // RevokeAllowanceExported
	builtin.MustGenerateExportedMethodNum("Burn"):              {"BurnExported", *new(func(*BurnParams) *BurnReturn)},                                // BurnExported
	builtin.MustGenerateExportedMethodNum("BurnFrom"):          {"BurnFromExported", *new(func(*BurnFromParams) *BurnFromReturn)},                    // BurnFromExported
	builtin.MustGenerateExportedMethodNum("Allowance"):         {"AllowanceExported", *new(func(*GetAllowanceParams) *abi.TokenAmount)},              // AllowanceExported
	builtin.MustGenerateExportedMethodNum("Granularity"):       {"GranularityExported", *new(func(value *abi.EmptyValue) *GranularityReturn)},        // GranularityExported
}
