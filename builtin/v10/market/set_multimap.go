package market

import (
	market9 "github.com/filecoin-project/go-state-types/builtin/v9/market"

	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
	cid "github.com/ipfs/go-cid"
)

type SetMultimap = market9.SetMultimap

func MakeEmptySetMultimap(s adt.Store, bitwidth int) (*SetMultimap, error) {
	return market9.MakeEmptySetMultimap(s, bitwidth)
}

// Writes a new empty map to the store and returns its CID.
func StoreEmptySetMultimap(s adt.Store, bitwidth int) (cid.Cid, error) {
	return market9.StoreEmptySetMultimap(s, bitwidth)
}
