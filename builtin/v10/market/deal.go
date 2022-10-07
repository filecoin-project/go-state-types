package market

import (
	market9 "github.com/filecoin-project/go-state-types/builtin/v9/market"
)

var PieceCIDPrefix = market9.PieceCIDPrefix

type DealState = market9.DealState
type DealLabel = market9.DealLabel

// Zero value of DealLabel is canonical EmptyDealLabel
var EmptyDealLabel = market9.EmptyDealLabel

func NewLabelFromString(s string) (DealLabel, error) {
	return market9.NewLabelFromString(s)
}

func NewLabelFromBytes(b []byte) (DealLabel, error) {
	return market9.NewLabelFromBytes(b)
}

type DealProposal = market9.DealProposal
type ClientDealProposal = market9.ClientDealProposal
