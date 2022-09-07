package account

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
)

var Methods = []interface{}{
	1: *new(func(interface{}, *address.Address) *abi.EmptyValue), // Constructor
	2: *new(func(interface{}, *abi.EmptyValue) *address.Address), // PubkeyAddress
}
