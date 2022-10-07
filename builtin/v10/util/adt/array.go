package adt

import (
	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"

	"github.com/ipfs/go-cid"
)

var DefaultAmtOptions = adt9.DefaultAmtOptions

type Array = adt9.Array

func AsArray(s Store, r cid.Cid, bitwidth int) (*Array, error) {
	return adt9.AsArray(s, r, bitwidth)
}

func MakeEmptyArray(s Store, bitwidth int) (*Array, error) {
	return adt9.MakeEmptyArray(s, bitwidth)
}

func StoreEmptyArray(s Store, bitwidth int) (cid.Cid, error) {
	return adt9.StoreEmptyArray(s, bitwidth)
}
