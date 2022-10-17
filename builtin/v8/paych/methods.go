package paych

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[uint64]builtin.MethodMeta{
	1: {"", *new(func(*ConstructorParams) *abi.EmptyValue)},        // Constructor
	2: {"", *new(func(*UpdateChannelStateParams) *abi.EmptyValue)}, // UpdateChannelState
	3: {"", *new(func(*abi.EmptyValue) *abi.EmptyValue)},           // Settle
	4: {"", *new(func(*abi.EmptyValue) *abi.EmptyValue)},           // Collect
}
