package paych

import (
	"github.com/filecoin-project/go-state-types/abi"
)

var Methods = map[uint64]interface{}{
	1: *new(func(interface{}, *ConstructorParams) *abi.EmptyValue),        // Constructor
	2: *new(func(interface{}, *UpdateChannelStateParams) *abi.EmptyValue), // UpdateChannelState
	3: *new(func(interface{}, *abi.EmptyValue) *abi.EmptyValue),           // Settle
	4: *new(func(interface{}, *abi.EmptyValue) *abi.EmptyValue),           // Collect
}
