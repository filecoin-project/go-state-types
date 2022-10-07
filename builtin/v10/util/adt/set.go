package adt

import (
	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"

	"github.com/ipfs/go-cid"
)

type Set = adt9.Set

func AsSet(s Store, r cid.Cid, bitwidth int) (*Set, error) {
	return adt9.AsSet(s, r, bitwidth)
}

func MakeEmptySet(s Store, bitwidth int) (*Set, error) {
	return adt9.MakeEmptySet(s, bitwidth)
}
