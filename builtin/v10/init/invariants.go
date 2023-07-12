package init

import (
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
)

type StateSummary struct {
	AddrIDs map[addr.Address]abi.ActorID
	NextID  abi.ActorID
}

// Checks internal invariants of init state.
func CheckStateInvariants(st *State, tree *builtin.ActorTree, actorCodes map[string]cid.Cid) (*StateSummary, *builtin.MessageAccumulator) {
	acc := &builtin.MessageAccumulator{}
	store := tree.Store

	acc.Require(len(st.NetworkName) > 0, "network name is empty")
	acc.Require(st.NextID >= builtin.FirstNonSingletonActorId, "next id %d is too low", st.NextID)

	initSummary := &StateSummary{
		AddrIDs: nil,
		NextID:  st.NextID,
	}

	lut, err := adt.AsMap(store, st.AddressMap, builtin.DefaultHamtBitwidth)
	if err != nil {
		acc.Addf("error loading address map: %v", err)
		// Stop here, it's hard to make other useful checks.
		return initSummary, acc
	}

	initSummary.AddrIDs = map[addr.Address]abi.ActorID{}
	reverse := map[abi.ActorID]addr.Address{}
	var value cbg.CborInt
	err = lut.ForEach(&value, func(key string) error {
		actorId := abi.ActorID(value)
		keyAddr, err := addr.NewFromBytes([]byte(key))
		if err != nil {
			return err
		}

		acc.Require(keyAddr.Protocol() != addr.ID, "key %v is an ID address", keyAddr)
		acc.Require(keyAddr.Protocol() <= addr.Delegated, "unknown address protocol for key %v", keyAddr)
		acc.Require(actorId >= builtin.FirstNonSingletonActorId, "unexpected singleton ID value %v", actorId)

		foundAddr, found := reverse[actorId]
		isPair := (keyAddr.Protocol() == addr.Actor && foundAddr.Protocol() == addr.Delegated) ||
			(keyAddr.Protocol() == addr.Delegated && foundAddr.Protocol() == addr.Actor)
		dup := found && !isPair
		acc.Require(!dup, "duplicate mapping to ID %v: %v, %v", actorId, keyAddr, foundAddr)
		reverse[actorId] = keyAddr

		initSummary.AddrIDs[keyAddr] = actorId

		idaddr, err := addr.NewIDAddress(uint64(actorId))
		acc.RequireNoError(err, "unable to convert actorId %v to id address", actorId)
		actor, found, err := tree.GetActorV5(idaddr)
		acc.RequireNoError(err, "unable to retrieve actor with idaddr %v", idaddr)
		acc.Require(found, "actor not found idaddr %v", idaddr)

		if keyAddr.Protocol() == addr.Delegated {
			acc.Require(canHaveDelegatedAddress(actor, actorCodes), "actor %v not supposed to have a delegated address", idaddr)
		}

		// we expect the address field to be populated for the below actors
		if (actor.Code == actorCodes[manifest.EthAccountKey] ||
			actor.Code == actorCodes[manifest.EvmKey] ||
			actor.Code == actorCodes[manifest.PlaceholderKey]) &&
			keyAddr.Protocol() != addr.Actor {
			acc.Require(keyAddr == *actor.Address, "address field in actor state differs from addr available in init actor map: actor=%v, init=%v", *actor.Address, keyAddr)
		}

		return nil
	})
	acc.RequireNoError(err, "error iterating address map")
	return initSummary, acc
}

func canHaveDelegatedAddress(actor *builtin.ActorV5, actorCodes map[string]cid.Cid) bool {
	return actor.Code == actorCodes[manifest.EthAccountKey] ||
		actor.Code == actorCodes[manifest.EvmKey] ||
		actor.Code == actorCodes[manifest.PlaceholderKey]
}
