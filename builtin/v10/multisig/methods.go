package multisig

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/frc0042"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*ConstructorParams) *abi.EmptyValue)}, // Constructor
	2: {"Propose", *new(func(*ProposeParams) *ProposeReturn)},          // Propose
	frc0042.GenerateExportedMethodNum("Propose"): {"ProposeExported", *new(func(*ProposeParams) *ProposeReturn)}, // ProposeExported
	3: {"Approve", *new(func(*TxnIDParams) *ApproveReturn)}, // Approve
	frc0042.GenerateExportedMethodNum("Approve"): {"ApproveExported", *new(func(*TxnIDParams) *ApproveReturn)}, // ApproveExported
	4: {"Cancel", *new(func(*TxnIDParams) *abi.EmptyValue)}, // Cancel
	frc0042.GenerateExportedMethodNum("Cancel"): {"CancelExported", *new(func(*TxnIDParams) *abi.EmptyValue)}, // CancelExported
	5: {"AddSigner", *new(func(*AddSignerParams) *abi.EmptyValue)}, // AddSigner
	frc0042.GenerateExportedMethodNum("AddSigner"): {"AddSignerExported", *new(func(*AddSignerParams) *abi.EmptyValue)}, // AddSignerExported
	6: {"RemoveSigner", *new(func(*RemoveSignerParams) *abi.EmptyValue)}, // RemoveSigner
	frc0042.GenerateExportedMethodNum("RemoveSigner"): {"RemoveSignerExported", *new(func(*RemoveSignerParams) *abi.EmptyValue)}, // RemoveSignerExported
	7: {"SwapSigner", *new(func(*SwapSignerParams) *abi.EmptyValue)}, // SwapSigner
	frc0042.GenerateExportedMethodNum("SwapSigner"): {"SwapSignerExported", *new(func(*SwapSignerParams) *abi.EmptyValue)}, // SwapSignerExported
	8: {"ChangeNumApprovalsThreshold", *new(func(*ChangeNumApprovalsThresholdParams) *abi.EmptyValue)}, // ChangeNumApprovalsThreshold
	frc0042.GenerateExportedMethodNum("ChangeNumApprovalsThreshold"): {"ChangeNumApprovalsThresholdExported", *new(func(*ChangeNumApprovalsThresholdParams) *abi.EmptyValue)}, // ChangeNumApprovalsThresholdExported
	9: {"LockBalance", *new(func(*LockBalanceParams) *abi.EmptyValue)}, // LockBalance
	frc0042.GenerateExportedMethodNum("LockBalance"): {"LockBalanceExported", *new(func(*LockBalanceParams) *abi.EmptyValue)},          // LockBalanceExported
	frc0042.GenerateExportedMethodNum("Receive"):     {"UniversalReceiverHook", *new(func(*abi.CborBytesTransparent) *abi.EmptyValue)}, // UniversalReceiverHook
}
