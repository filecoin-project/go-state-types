package datacap

import (
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
	"github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
)

type StateSummary struct {
	Balances    map[abi.ActorID]abi.TokenAmount
	Allowances  map[abi.ActorID]map[abi.ActorID]abi.TokenAmount
	TotalSupply abi.TokenAmount
}

// Checks internal invariants of verified registry state.
func CheckStateInvariants(st *State, store adt.Store) (*StateSummary, *builtin.MessageAccumulator) {
	acc := &builtin.MessageAccumulator{}
	acc.Require(st.Governor.Protocol() == addr.ID, "governor %v must be ID address", st.Governor)
	checkTokenInvariants(store, st.Token, acc)

	// Check clients
	allBalances := make(map[abi.ActorID]abi.TokenAmount)
	balances, err := adt.AsMap(store, st.Token.Balances, int(st.Token.HamtBitWidth))
	acc.RequireNoError(err, "error getting balances map")

	var balance abi.StoragePower
	err = balances.ForEach(&balance, func(idKey string) error {
		actorId, err := abi.ParseUIntKey(idKey)
		acc.RequireNoError(err, "error parsing actor id to uint")

		allBalances[abi.ActorID(actorId)] = balance.Copy()
		return nil
	})
	acc.RequireNoError(err, "error iterating clients")

	allAllowances := make(map[abi.ActorID]map[abi.ActorID]abi.TokenAmount)
	allowancesMap, err := adt.AsMap(store, st.Token.Allowances, int(st.Token.HamtBitWidth))
	acc.RequireNoError(err, "error getting allowances outer map")

	var innerHamtCid cbg.CborCid
	err = allowancesMap.ForEach(&innerHamtCid, func(idKey string) error {
		owner, err := abi.ParseUIntKey(idKey)
		acc.RequireNoError(err, "error parsing operator id to uint")

		allowances := make(map[abi.ActorID]abi.TokenAmount)
		allowancesInnerMap, err := adt.AsMap(store, cid.Cid(innerHamtCid), int(st.Token.HamtBitWidth))
		acc.RequireNoError(err, "error getting allowances inner map")

		var amount abi.TokenAmount
		err = allowancesInnerMap.ForEach(&amount, func(idKey string) error {
			operator, err := abi.ParseUIntKey(idKey)
			acc.RequireNoError(err, "error parsing operator id to uint")

			allowances[abi.ActorID(operator)] = amount.Copy()
			return nil
		})
		acc.RequireNoError(err, "error iterating over inner allowances map")

		allAllowances[abi.ActorID(owner)] = allowances

		return nil
	})
	acc.RequireNoError(err, "error iterating over outer allowances map")

	return &StateSummary{
		Balances:    allBalances,
		Allowances:  allAllowances,
		TotalSupply: st.Token.Supply.Copy(),
	}, acc
}

// this can be extracted out to check any token contract when we have more than one
func checkTokenInvariants(store adt.Store, tokenState TokenState, acc *builtin.MessageAccumulator) {
	acc.Require(tokenState.Supply.GreaterThanEqual(big.Zero()), "token supply %d cannot be negative", tokenState.Supply)

	// Balances
	balances, err := adt.AsMap(store, tokenState.Balances, int(tokenState.HamtBitWidth))
	acc.RequireNoError(err, "error getting balances map")

	var balanceSum = big.Zero()
	var balance abi.StoragePower
	err = balances.ForEach(&balance, func(idKey string) error {
		actorId, err := abi.ParseUIntKey(idKey)
		acc.RequireNoError(err, "error parsing actor id to uint")

		acc.Require(balance.GreaterThan(big.Zero()), "balance for actor %d is not positive %d", actorId, balance)

		balanceSum = big.Add(balanceSum, balance)

		return nil
	})
	acc.RequireNoError(err, "error iterating clients")

	acc.Require(balanceSum.Equals(tokenState.Supply), "token supply %d does not equal sum of all balances %d", tokenState.Supply, balanceSum)

	// Allowances
	allowancesMap, err := adt.AsMap(store, tokenState.Allowances, int(tokenState.HamtBitWidth))
	acc.RequireNoError(err, "error getting allowances outer map")

	var innerHamtCid cbg.CborCid
	err = allowancesMap.ForEach(&innerHamtCid, func(idKey string) error {
		owner, err := abi.ParseUIntKey(idKey)
		acc.RequireNoError(err, "error parsing operator id to uint")

		allowances := make(map[abi.ActorID]abi.TokenAmount)
		allowancesInnerMap, err := adt.AsMap(store, cid.Cid(innerHamtCid), int(tokenState.HamtBitWidth))
		acc.RequireNoError(err, "error getting allowances inner map")

		var amount abi.TokenAmount
		err = allowancesInnerMap.ForEach(&amount, func(idKey string) error {
			operator, err := abi.ParseUIntKey(idKey)
			acc.RequireNoError(err, "error parsing operator id to uint")

			acc.Require(owner != operator, "owner %d cannot self-store allowance", owner)
			acc.Require(amount.GreaterThan(big.Zero()), "balance %d must be positive", amount)

			allowances[abi.ActorID(operator)] = amount.Copy()
			return nil
		})
		acc.RequireNoError(err, "error iterating over inner allowances map")

		return nil
	})
	acc.RequireNoError(err, "error iterating over outer allowances map")
}
