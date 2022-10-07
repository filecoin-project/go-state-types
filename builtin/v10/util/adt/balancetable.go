package adt

import (
	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"

	"github.com/ipfs/go-cid"
)

const BalanceTableBitwidth = adt9.BalanceTableBitwidth

type BalanceTable = adt9.BalanceTable

func AsBalanceTable(s Store, r cid.Cid) (*BalanceTable, error) {
	return adt9.AsBalanceTable(s, r)
}
