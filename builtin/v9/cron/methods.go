package cron

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

var Methods = map[abi.MethodNum]builtin.MethodMeta{
	1: {"Constructor", *new(func(*ConstructorParams) *abi.EmptyValue)}, // Constructor
	2: {"EpochTick", *new(func(*abi.EmptyValue) *abi.EmptyValue)},      // EpochTick
}
