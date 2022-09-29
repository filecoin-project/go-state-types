package account

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[uint64]interface{}{
	1: *new(func(interface{}, *address.Address) *abi.EmptyValue),           // Constructor
	2: *new(func(interface{}, *abi.EmptyValue) *address.Address),           // PubkeyAddress
	3: *new(func(interface{}, *AuthenticateMessageParams) *abi.EmptyValue), // AuthenticateMessage
	uint64(builtin.UniversalReceiverHookMethodNum): *new(func(interface{}, *[]byte) *abi.EmptyValue), // UniversalReceiverHook
}
