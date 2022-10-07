package verifreg

import (
	verifreg9 "github.com/filecoin-project/go-state-types/builtin/v9/verifreg"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
)

type DataCap = verifreg9.DataCap

var DataCapGranularity = verifreg9.DataCapGranularity

const SignatureDomainSeparation_RemoveDataCap = verifreg9.SignatureDomainSeparation_RemoveDataCap

type RmDcProposalID = verifreg9.RmDcProposalID
type State = verifreg9.State

var MinVerifiedDealSize = verifreg9.MinVerifiedDealSize

// rootKeyAddress comes from genesis.
func ConstructState(store adt.Store, rootKeyAddress address.Address) (*State, error) {
	return verifreg9.ConstructState(store, rootKeyAddress)
}
