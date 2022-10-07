package adt

import (
	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"

	"github.com/ipfs/go-cid"
)

var DefaultHamtOptions = adt9.DefaultHamtOptions

type Map = adt9.Map

func AsMap(s Store, root cid.Cid, bitwidth int) (*Map, error) {
	return adt9.AsMap(s, root, bitwidth)
}

func MakeEmptyMap(s Store, bitwidth int) (*Map, error) {
	return adt9.MakeEmptyMap(s, bitwidth)
}

func StoreEmptyMap(s Store, bitwidth int) (cid.Cid, error) {
	return adt9.StoreEmptyMap(s, bitwidth)
}
