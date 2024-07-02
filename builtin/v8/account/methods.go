package account

import (
	"github.com/filecoin-project/go-address"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: builtin.NewMethodMeta("Constructor", *new(func(*address.Address) *abi.EmptyValue)),   // Constructor
	2: builtin.NewMethodMeta("PubkeyAddress", *new(func(*abi.EmptyValue) *address.Address)), // PubkeyAddress
}
