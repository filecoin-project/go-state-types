package cron

import (
	"github.com/filecoin-project/go-state-types/abi"
)

var Methods = []interface{}{
	1: *new(func(interface{}, *ConstructorParams) *abi.EmptyValue), // Constructor
	2: *new(func(interface{}, *abi.EmptyValue) *abi.EmptyValue),    // EpochTick
}
