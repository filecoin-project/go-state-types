package account

import (
	"github.com/filecoin-project/go-address"
	typegen "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*address.Address) *abi.EmptyValue)},   // Constructor
	2: {"PubkeyAddress", *new(func(*abi.EmptyValue) *address.Address)}, // PubkeyAddress
	builtin.MustGenerateFRCMethodNum("AuthenticateMessage"): {"AuthenticateMessage", *new(func(*AuthenticateMessageParams) *typegen.CborBool)}, // AuthenticateMessage
}
