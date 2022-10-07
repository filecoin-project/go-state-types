package market

import (
	market9 "github.com/filecoin-project/go-state-types/builtin/v9/market"

	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
	"github.com/ipfs/go-cid"
)

const EpochUndefined = market9.EpochUndefined
const ProposalsAmtBitwidth = market9.ProposalsAmtBitwidth
const StatesAmtBitwidth = market9.StatesAmtBitwidth

type State = market9.State

func ConstructState(store adt.Store) (*State, error) {
	return market9.ConstructState(store)
}

type DealArray = market9.DealArray

func AsDealProposalArray(s adt.Store, r cid.Cid) (*DealArray, error) {
	return market9.AsDealProposalArray(s, r)
}

func ValidateDealsForActivation(
	st *State, store adt.Store, dealIDs []abi.DealID, minerAddr addr.Address, sectorExpiry, currEpoch abi.ChainEpoch,
) (big.Int, big.Int, uint64, error) {
	return market9.ValidateDealsForActivation(st, store, dealIDs, minerAddr, sectorExpiry, currEpoch)
}
