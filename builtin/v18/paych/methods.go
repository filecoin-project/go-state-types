package paych

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: builtin.NewMethodMeta("Constructor", *new(func(*ConstructorParams) *abi.EmptyValue)),               // Constructor
	2: builtin.NewMethodMeta("UpdateChannelState", *new(func(*UpdateChannelStateParams) *abi.EmptyValue)), // UpdateChannelState
	3: builtin.NewMethodMeta("Settle", *new(func(*abi.EmptyValue) *abi.EmptyValue)),                       // Settle
	4: builtin.NewMethodMeta("Collect", *new(func(*abi.EmptyValue) *abi.EmptyValue)),                      // Collect
}
