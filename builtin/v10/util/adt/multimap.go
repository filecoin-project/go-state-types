package adt

import (
	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"

	"github.com/ipfs/go-cid"
)

type Multimap = adt9.Multimap

func MakeEmptyMultimap(s Store, outerBitwidth, innerBitwidth int) (*Multimap, error) {
	return adt9.MakeEmptyMultimap(s, outerBitwidth, innerBitwidth)
}

func StoreEmptyMultimap(store Store, outerBitwidth, innerBitwidth int) (cid.Cid, error) {
	return adt9.StoreEmptyMultimap(store, outerBitwidth, innerBitwidth)
}
