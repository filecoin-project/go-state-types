package cron

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[uint64]builtin.MethodMeta{
	1: {"", *new(func(*ConstructorParams) *abi.EmptyValue)}, // Constructor
	2: {"", *new(func(*abi.EmptyValue) *abi.EmptyValue)},    // EpochTick
}
