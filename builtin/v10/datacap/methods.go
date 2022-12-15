package datacap

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/builtin/frc0042"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1:                                        {"Constructor", *new(func(*address.Address) *abi.EmptyValue)}, // Constructor
	frc0042.GenerateExportedMethodNum("Mint"):              {"MintExported", *new(func(*MintParams) *MintReturn)},                                // MintExported
	frc0042.GenerateExportedMethodNum("Destroy"):           {"DestroyExported", *new(func(*DestroyParams) *BurnReturn)},                          // DestroyExported
	frc0042.GenerateExportedMethodNum("Name"):              {"NameExported", *new(func(*abi.EmptyValue) *abi.CborString)},                        // NameExported
	frc0042.GenerateExportedMethodNum("Symbol"):            {"SymbolExported", *new(func(*abi.EmptyValue) *abi.CborString)},                      // SymbolExported
	frc0042.GenerateExportedMethodNum("TotalSupply"):       {"TotalSupplyExported", *new(func(*abi.EmptyValue) *abi.TokenAmount)},                // TotalSupplyExported
	frc0042.GenerateExportedMethodNum("Balance"):           {"BalanceExported", *new(func(*address.Address) *abi.TokenAmount)},                   // BalanceExported
	frc0042.GenerateExportedMethodNum("Transfer"):          {"TransferExported", *new(func(*TransferParams) *TransferReturn)},                    // TransferExported
	frc0042.GenerateExportedMethodNum("TransferFrom"):      {"TransferFromExported", *new(func(*TransferFromParams) *TransferFromReturn)},        // TransferFromExported
	frc0042.GenerateExportedMethodNum("IncreaseAllowance"): {"IncreaseAllowanceExported", *new(func(*IncreaseAllowanceParams) *abi.TokenAmount)}, // IncreaseAllowanceExported
	frc0042.GenerateExportedMethodNum("DecreaseAllowance"): {"DecreaseAllowanceExported", *new(func(*DecreaseAllowanceParams) *abi.TokenAmount)}, // DecreaseAllowanceExported
	frc0042.GenerateExportedMethodNum("RevokeAllowance"):   {"RevokeAllowanceExported", *new(func(*RevokeAllowanceParams) *abi.TokenAmount)},     // RevokeAllowanceExported
	frc0042.GenerateExportedMethodNum("Burn"):              {"BurnExported", *new(func(*BurnParams) *BurnReturn)},                                // BurnExported
	frc0042.GenerateExportedMethodNum("BurnFrom"):          {"BurnFromExported", *new(func(*BurnFromParams) *BurnFromReturn)},                    // BurnFromExported
	frc0042.GenerateExportedMethodNum("Allowance"):         {"AllowanceExported", *new(func(*GetAllowanceParams) *abi.TokenAmount)},              // AllowanceExported
	frc0042.GenerateExportedMethodNum("Granularity"):       {"GranularityExported", *new(func(value *abi.EmptyValue) *GranularityReturn)},        // GranularityExported
}
