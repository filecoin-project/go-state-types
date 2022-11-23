package verifreg

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
	Verifiers   map[addr.Address]DataCap
	Claims      map[ClaimId]Claim
	Allocations map[AllocationId]Allocation
}

// Checks internal invariants of verified registry state.
func CheckStateInvariants(st *State, store adt.Store, priorEpoch abi.ChainEpoch) (*StateSummary, *builtin.MessageAccumulator) {
	acc := &builtin.MessageAccumulator{}
	acc.Require(st.RootKey.Protocol() == addr.ID, "root key %v should have ID protocol", st.RootKey)

	// Check verifiers
	allVerifiers := map[addr.Address]DataCap{}
	if verifiers, err := adt.AsMap(store, st.Verifiers, builtin.DefaultHamtBitwidth); err != nil {
		acc.Addf("error loading verifiers: %v", err)
	} else {
		var vcap abi.StoragePower
		err = verifiers.ForEach(&vcap, func(key string) error {
			verifier, err := addr.NewFromBytes([]byte(key))
			acc.RequireNoError(err, "error getting verifier from bytes")

			acc.Require(verifier.Protocol() == addr.ID, "verifier %v should have ID protocol", verifier)
			acc.Require(vcap.GreaterThanEqual(big.Zero()), "verifier %v cap %v is negative", verifier, vcap)
			allVerifiers[verifier] = vcap.Copy()
			return nil
		})
		acc.RequireNoError(err, "error iterating verifiers")
	}

	// Check Claims
	allClaims := make(map[ClaimId]Claim)
	outerClaimsMap, err := adt.AsMap(store, st.Claims, builtin.DefaultHamtBitwidth)
	acc.RequireNoError(err, "error getting outer claims map")

	var innerClaimHamtCid cbg.CborCid
	err = outerClaimsMap.ForEach(&innerClaimHamtCid, func(idKey string) error {
		providerId, err := abi.ParseUIntKey(idKey)
		acc.RequireNoError(err, "error parsing provider id to uint")
		innerMap, err := adt.AsMap(store, cid.Cid(innerClaimHamtCid), builtin.DefaultHamtBitwidth)
		acc.RequireNoError(err, "error getting inner claims map")

		var out Claim
		err = innerMap.ForEach(&out, func(claimIdStr string) error {
			claimId, err := abi.ParseUIntKey(claimIdStr)
			acc.RequireNoError(err, "error parsing claim id to uint")

			checkClaimState(ClaimId(claimId), &out, abi.ActorID(providerId), st.NextAllocationId, priorEpoch, acc)
			allClaims[ClaimId(claimId)] = out
			return nil
		})
		acc.RequireNoError(err, "error iterating over inner claims map")
		return nil
	})
	acc.RequireNoError(err, "error iterating over claims")

	// Check Allocations
	allAllocations := make(map[AllocationId]Allocation)
	outerAllocMap, err := adt.AsMap(store, st.Allocations, builtin.DefaultHamtBitwidth)
	acc.RequireNoError(err, "error getting outer allocations map")

	var innerAllocHamtCid cbg.CborCid
	err = outerAllocMap.ForEach(&innerAllocHamtCid, func(idKey string) error {
		clientId, err := abi.ParseUIntKey(idKey)
		acc.RequireNoError(err, "error parsing client id to uint")
		innerMap, err := adt.AsMap(store, cid.Cid(innerAllocHamtCid), builtin.DefaultHamtBitwidth)
		acc.RequireNoError(err, "error getting inner allocations map")

		var alloc Allocation
		err = innerMap.ForEach(&alloc, func(claimIdStr string) error {
			allocId, err := abi.ParseUIntKey(claimIdStr)
			acc.RequireNoError(err, "error parsing allocation id to uint")

			checkAllocationState(AllocationId(allocId), &alloc, abi.ActorID(clientId), st.NextAllocationId, priorEpoch, acc)
			allAllocations[AllocationId(allocId)] = alloc
			return nil
		})
		acc.RequireNoError(err, "error iterating over inner allocations map")
		return nil
	})
	acc.RequireNoError(err, "error iterating over allocations")

	return &StateSummary{
		Verifiers:   allVerifiers,
		Claims:      allClaims,
		Allocations: allAllocations,
	}, acc
}

func checkAllocationState(id AllocationId, alloc *Allocation, client abi.ActorID, nextAllocId AllocationId, priorEpoch abi.ChainEpoch, acc *builtin.MessageAccumulator) {
	acc.Require(id < nextAllocId, "allocation id %d exceeds next %d", id, nextAllocId)

	acc.Require(alloc.Client == client, "allocation %d client %d doesn't match key %d", id, alloc.Client, client)

	acc.Require(alloc.Size >= MinimumVerifiedAllocationSize, "allocation %d size %d too small", id, alloc.Size)

	acc.Require(alloc.TermMin >= MinimumVerifiedAllocationTerm, "allocation %d term min %d too small", id, alloc.TermMin)

	acc.Require(alloc.TermMax <= MaximumVerifiedAllocationTerm, "allocation %d term max %d too large", id, alloc.TermMax)

	acc.Require(alloc.TermMin <= alloc.TermMax, "allocation %d term min %d exceeds max %d", id, alloc.TermMin, alloc.TermMin)

	acc.Require(alloc.Expiration <= priorEpoch+MaximumVerifiedAllocationExpiration, "allocation %d expiration %d too far from now %d", id, alloc.Expiration, priorEpoch)
}

func checkClaimState(id ClaimId, claim *Claim, provider abi.ActorID, nextAllocId AllocationId, priorEpoch abi.ChainEpoch, acc *builtin.MessageAccumulator) {
	acc.Require(AllocationId(id) < nextAllocId, "claim id %d exceeds next %d", id, nextAllocId)

	acc.Require(claim.Provider == provider, "claim %d provider %d doesn't match key %d", id, claim.Provider, provider)

	acc.Require(claim.Size >= MinimumVerifiedAllocationSize, "claim %d size %d too small", id, claim.Size)

	acc.Require(claim.TermMin >= MinimumVerifiedAllocationTerm, "claim %d term min %d too small", id, claim.TermMin)

	acc.Require(claim.TermMin <= claim.TermMax, "claim %d term min %d exceeds max %d", id, claim.TermMin, claim.TermMin)

	acc.Require(claim.TermStart <= priorEpoch, "claim %d term start %d after now %d", id, claim.TermStart, priorEpoch)
}
