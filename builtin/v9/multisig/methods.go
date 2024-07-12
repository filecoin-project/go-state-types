package multisig

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: builtin.NewMethodMeta("Constructor", *new(func(*ConstructorParams) *abi.EmptyValue)),                                 // Constructor
	2: builtin.NewMethodMeta("Propose", *new(func(*ProposeParams) *ProposeReturn)),                                          // Propose
	3: builtin.NewMethodMeta("Approve", *new(func(*TxnIDParams) *ApproveReturn)),                                            // Approve
	4: builtin.NewMethodMeta("Cancel", *new(func(*TxnIDParams) *abi.EmptyValue)),                                            // Cancel
	5: builtin.NewMethodMeta("AddSigner", *new(func(*AddSignerParams) *abi.EmptyValue)),                                     // AddSigner
	6: builtin.NewMethodMeta("RemoveSigner", *new(func(*RemoveSignerParams) *abi.EmptyValue)),                               // RemoveSigner
	7: builtin.NewMethodMeta("SwapSigner", *new(func(*SwapSignerParams) *abi.EmptyValue)),                                   // SwapSigner
	8: builtin.NewMethodMeta("ChangeNumApprovalsThreshold", *new(func(*ChangeNumApprovalsThresholdParams) *abi.EmptyValue)), // ChangeNumApprovalsThreshold
	9: builtin.NewMethodMeta("LockBalance", *new(func(*LockBalanceParams) *abi.EmptyValue)),                                 // LockBalance
	builtin.MustGenerateFRCMethodNum("Receive"): builtin.NewMethodMeta("UniversalReceiverHook", *new(func(*abi.CborBytesTransparent) *abi.EmptyValue)), // UniversalReceiverHook
}
