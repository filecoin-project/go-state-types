package multisig

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*ConstructorParams) *abi.EmptyValue)}, // Constructor
	2: {"Propose", *new(func(*ProposeParams) *ProposeReturn)},          // Propose
	builtin.MustGenerateFRCMethodNum("Propose"): {"ProposeExported", *new(func(*ProposeParams) *ProposeReturn)}, // ProposeExported
	3: {"Approve", *new(func(*TxnIDParams) *ApproveReturn)}, // Approve
	builtin.MustGenerateFRCMethodNum("Approve"): {"ApproveExported", *new(func(*TxnIDParams) *ApproveReturn)}, // ApproveExported
	4: {"Cancel", *new(func(*TxnIDParams) *abi.EmptyValue)}, // Cancel
	builtin.MustGenerateFRCMethodNum("Cancel"): {"CancelExported", *new(func(*TxnIDParams) *abi.EmptyValue)}, // CancelExported
	5: {"AddSigner", *new(func(*AddSignerParams) *abi.EmptyValue)}, // AddSigner
	builtin.MustGenerateFRCMethodNum("AddSigner"): {"AddSignerExported", *new(func(*AddSignerParams) *abi.EmptyValue)}, // AddSignerExported
	6: {"RemoveSigner", *new(func(*RemoveSignerParams) *abi.EmptyValue)}, // RemoveSigner
	builtin.MustGenerateFRCMethodNum("RemoveSigner"): {"RemoveSignerExported", *new(func(*RemoveSignerParams) *abi.EmptyValue)}, // RemoveSignerExported
	7: {"SwapSigner", *new(func(*SwapSignerParams) *abi.EmptyValue)}, // SwapSigner
	builtin.MustGenerateFRCMethodNum("SwapSigner"): {"SwapSignerExported", *new(func(*SwapSignerParams) *abi.EmptyValue)}, // SwapSignerExported
	8: {"ChangeNumApprovalsThreshold", *new(func(*ChangeNumApprovalsThresholdParams) *abi.EmptyValue)}, // ChangeNumApprovalsThreshold
	builtin.MustGenerateFRCMethodNum("ChangeNumApprovalsThreshold"): {"ChangeNumApprovalsThresholdExported", *new(func(*ChangeNumApprovalsThresholdParams) *abi.EmptyValue)}, // ChangeNumApprovalsThresholdExported
	9: {"LockBalance", *new(func(*LockBalanceParams) *abi.EmptyValue)}, // LockBalance
	builtin.MustGenerateFRCMethodNum("LockBalance"): {"LockBalanceExported", *new(func(*LockBalanceParams) *abi.EmptyValue)},          // LockBalanceExported
	builtin.MustGenerateFRCMethodNum("Receive"):     {"UniversalReceiverHook", *new(func(*abi.CborBytesTransparent) *abi.EmptyValue)}, // UniversalReceiverHook
}
