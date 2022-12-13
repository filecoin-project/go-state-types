package multisig

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*ConstructorParams) *abi.EmptyValue)},                                 // Constructor
	2: {"Propose", *new(func(*ProposeParams) *ProposeReturn)},                                          // Propose
	3: {"Approve", *new(func(*TxnIDParams) *ApproveReturn)},                                            // Approve
	4: {"Cancel", *new(func(*TxnIDParams) *abi.EmptyValue)},                                            // Cancel
	5: {"AddSigner", *new(func(*AddSignerParams) *abi.EmptyValue)},                                     // AddSigner
	6: {"RemoveSigner", *new(func(*RemoveSignerParams) *abi.EmptyValue)},                               // RemoveSigner
	7: {"SwapSigner", *new(func(*SwapSignerParams) *abi.EmptyValue)},                                   // SwapSigner
	8: {"ChangeNumApprovalsThreshold", *new(func(*ChangeNumApprovalsThresholdParams) *abi.EmptyValue)}, // ChangeNumApprovalsThreshold
	9: {"LockBalance", *new(func(*LockBalanceParams) *abi.EmptyValue)},                                 // LockBalance
}
