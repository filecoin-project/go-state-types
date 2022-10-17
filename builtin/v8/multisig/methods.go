package multisig

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[uint64]builtin.MethodMeta{
	1: {"", *new(func(*ConstructorParams) *abi.EmptyValue)},                 // Constructor
	2: {"", *new(func(*ProposeParams) *ProposeReturn)},                      // Propose
	3: {"", *new(func(*TxnIDParams) *ApproveReturn)},                        // Approve
	4: {"", *new(func(*TxnIDParams) *abi.EmptyValue)},                       // Cancel
	5: {"", *new(func(*AddSignerParams) *abi.EmptyValue)},                   // AddSigner
	6: {"", *new(func(*RemoveSignerParams) *abi.EmptyValue)},                // RemoveSigner
	7: {"", *new(func(*SwapSignerParams) *abi.EmptyValue)},                  // SwapSigner
	8: {"", *new(func(*ChangeNumApprovalsThresholdParams) *abi.EmptyValue)}, // ChangeNumApprovalsThreshold
	9: {"", *new(func(*LockBalanceParams) *abi.EmptyValue)},                 // LockBalance
}
