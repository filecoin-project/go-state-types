package account

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[uint64]builtin.MethodMeta{
	1: {"", *new(func(*address.Address) *abi.EmptyValue)}, // Constructor
	2: {"", *new(func(*abi.EmptyValue) *address.Address)}, // PubkeyAddress
}
