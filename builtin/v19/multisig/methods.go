package multisig

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: builtin.NewMethodMeta("Constructor", *new(func(*ConstructorParams) *abi.EmptyValue)), // Constructor
	2: builtin.NewMethodMeta("Propose", *new(func(*ProposeParams) *ProposeReturn)),          // Propose
	builtin.MustGenerateFRCMethodNum("Propose"): builtin.NewMethodMeta("ProposeExported", *new(func(*ProposeParams) *ProposeReturn)), // ProposeExported
	3: builtin.NewMethodMeta("Approve", *new(func(*TxnIDParams) *ApproveReturn)), // Approve
	builtin.MustGenerateFRCMethodNum("Approve"): builtin.NewMethodMeta("ApproveExported", *new(func(*TxnIDParams) *ApproveReturn)), // ApproveExported
	4: builtin.NewMethodMeta("Cancel", *new(func(*TxnIDParams) *abi.EmptyValue)), // Cancel
	builtin.MustGenerateFRCMethodNum("Cancel"): builtin.NewMethodMeta("CancelExported", *new(func(*TxnIDParams) *abi.EmptyValue)), // CancelExported
	5: builtin.NewMethodMeta("AddSigner", *new(func(*AddSignerParams) *abi.EmptyValue)), // AddSigner
	builtin.MustGenerateFRCMethodNum("AddSigner"): builtin.NewMethodMeta("AddSignerExported", *new(func(*AddSignerParams) *abi.EmptyValue)), // AddSignerExported
	6: builtin.NewMethodMeta("RemoveSigner", *new(func(*RemoveSignerParams) *abi.EmptyValue)), // RemoveSigner
	builtin.MustGenerateFRCMethodNum("RemoveSigner"): builtin.NewMethodMeta("RemoveSignerExported", *new(func(*RemoveSignerParams) *abi.EmptyValue)), // RemoveSignerExported
	7: builtin.NewMethodMeta("SwapSigner", *new(func(*SwapSignerParams) *abi.EmptyValue)), // SwapSigner
	builtin.MustGenerateFRCMethodNum("SwapSigner"): builtin.NewMethodMeta("SwapSignerExported", *new(func(*SwapSignerParams) *abi.EmptyValue)), // SwapSignerExported
	8: builtin.NewMethodMeta("ChangeNumApprovalsThreshold", *new(func(*ChangeNumApprovalsThresholdParams) *abi.EmptyValue)), // ChangeNumApprovalsThreshold
	builtin.MustGenerateFRCMethodNum("ChangeNumApprovalsThreshold"): builtin.NewMethodMeta("ChangeNumApprovalsThresholdExported", *new(func(*ChangeNumApprovalsThresholdParams) *abi.EmptyValue)), // ChangeNumApprovalsThresholdExported
	9: builtin.NewMethodMeta("LockBalance", *new(func(*LockBalanceParams) *abi.EmptyValue)), // LockBalance
	builtin.MustGenerateFRCMethodNum("LockBalance"): builtin.NewMethodMeta("LockBalanceExported", *new(func(*LockBalanceParams) *abi.EmptyValue)),          // LockBalanceExported
	builtin.MustGenerateFRCMethodNum("Receive"):     builtin.NewMethodMeta("UniversalReceiverHook", *new(func(*abi.CborBytesTransparent) *abi.EmptyValue)), // UniversalReceiverHook
}
