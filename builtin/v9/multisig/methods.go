package multisig

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[uint64]interface{}{
	1: *new(func(interface{}, *ConstructorParams) *abi.EmptyValue),                 // Constructor
	2: *new(func(interface{}, *ProposeParams) *ProposeReturn),                      // Propose
	3: *new(func(interface{}, *TxnIDParams) *ApproveReturn),                        // Approve
	4: *new(func(interface{}, *TxnIDParams) *abi.EmptyValue),                       // Cancel
	5: *new(func(interface{}, *AddSignerParams) *abi.EmptyValue),                   // AddSigner
	6: *new(func(interface{}, *RemoveSignerParams) *abi.EmptyValue),                // RemoveSigner
	7: *new(func(interface{}, *SwapSignerParams) *abi.EmptyValue),                  // SwapSigner
	8: *new(func(interface{}, *ChangeNumApprovalsThresholdParams) *abi.EmptyValue), // ChangeNumApprovalsThreshold
	9: *new(func(interface{}, *LockBalanceParams) *abi.EmptyValue),                 // LockBalance
	uint64(builtin.UniversalReceiverHookMethodNum): *new(func(interface{}, *[]byte) *abi.EmptyValue), // UniversalReceiverHook
}
