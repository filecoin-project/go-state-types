package migration

import (
	"context"

	init8 "github.com/filecoin-project/go-state-types/builtin/v8/init"

	verifreg8 "github.com/filecoin-project/go-state-types/builtin/v8/verifreg"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	market8 "github.com/filecoin-project/go-state-types/builtin/v8/market"
	adt8 "github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
	market9 "github.com/filecoin-project/go-state-types/builtin/v9/market"
	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	verifreg9 "github.com/filecoin-project/go-state-types/builtin/v9/verifreg"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

func migrateVerifreg(ctx context.Context, adtStore adt8.Store, priorEpoch abi.ChainEpoch, initStateV9 init8.State, marketStateV8 market8.State, verifregStateV8 verifreg8.State, emptyMapCid cid.Cid) (cid.Cid, map[abi.DealID]verifreg9.AllocationId, error) {
	pendingProposals, err := adt8.AsSet(adtStore, marketStateV8.PendingProposals, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, nil, xerrors.Errorf("failed to load pending proposals: %w", err)
	}

	proposals, err := market8.AsDealProposalArray(adtStore, marketStateV8.Proposals)
	if err != nil {
		return cid.Undef, nil, xerrors.Errorf("failed to get proposals: %w", err)
	}

	pendingMap := make(map[abi.DealID]market8.DealProposal)
	var proposal market8.DealProposal
	if err = proposals.ForEach(&proposal, func(dealID int64) error {
		// Nothing to do for unverified deals
		if !proposal.VerifiedDeal {
			return nil
		}

		pcid, err := proposal.Cid()
		if err != nil {
			return err
		}

		isPending, err := pendingProposals.Has(abi.CidKey(pcid))
		if err != nil {
			return xerrors.Errorf("failed to check pending: %w", err)
		}

		if isPending {
			pendingMap[abi.DealID(dealID)] = proposal
		}

		return nil
	}); err != nil {
		return cid.Undef, nil, xerrors.Errorf("failed to iterate over proposals: %w", err)
	}

	nextAllocationID := verifreg9.AllocationId(1)
	allocationsMapMap := make(map[address.Address]*adt9.Map)
	dealsToAllocations := make(map[abi.DealID]verifreg9.AllocationId)
	for dealID, proposal := range pendingMap {
		clientIDAddress, ok, err := initStateV9.ResolveAddress(adtStore, proposal.Client)
		if err != nil {
			return cid.Undef, nil, xerrors.Errorf("failed to resolve client %s: %w", proposal.Client, err)
		}

		if !ok {
			return cid.Undef, nil, xerrors.New("failed to find client in init actor map")
		}

		clientIDu64, err := address.IDFromAddress(clientIDAddress)
		if err != nil {
			return cid.Undef, nil, err
		}

		providerIDAddress, ok, err := initStateV9.ResolveAddress(adtStore, proposal.Provider)
		if err != nil {
			return cid.Undef, nil, xerrors.Errorf("failed to resolve provider %s: %w", proposal.Provider, err)
		}

		if !ok {
			return cid.Undef, nil, xerrors.New("failed to find provider in init actor map")
		}

		providerIDu64, err := address.IDFromAddress(providerIDAddress)
		if err != nil {
			return cid.Undef, nil, err
		}

		clientAllocationMap, ok := allocationsMapMap[clientIDAddress]
		if !ok {
			clientAllocationMap, err = adt9.AsMap(adtStore, emptyMapCid, builtin.DefaultHamtBitwidth)
			if err != nil {
				return cid.Undef, nil, xerrors.Errorf("failed to load empty map: %w", err)
			}

			allocationsMapMap[clientIDAddress] = clientAllocationMap
		}

		if err = clientAllocationMap.Put(nextAllocationID, &verifreg9.Allocation{
			Client:     abi.ActorID(clientIDu64),
			Provider:   abi.ActorID(providerIDu64),
			Data:       proposal.PieceCID,
			Size:       proposal.PieceSize,
			TermMin:    proposal.Duration(),
			TermMax:    market9.DealMaxDuration,
			Expiration: verifreg9.MaximumVerifiedAllocationExpiration + priorEpoch,
		}); err != nil {
			return cid.Undef, nil, xerrors.Errorf("failed to put new allocation obj: %w", err)
		}

		dealsToAllocations[dealID] = nextAllocationID
		nextAllocationID++
	}

	allocationsMap, err := adt9.AsMap(adtStore, emptyMapCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, nil, xerrors.Errorf("failed to load empty map: %w", err)
	}

	for clientID, clientAllocationsMap := range allocationsMapMap {
		if err = allocationsMap.Put(abi.IdAddrKey(clientID), clientAllocationsMap); err != nil {
			return cid.Undef, nil, xerrors.Errorf("failed to populate allocationsMap: %w", err)
		}
	}

	allocationsMapRoot, err := allocationsMap.Root()
	if err != nil {
		return cid.Undef, nil, xerrors.Errorf("failed to flush allocations map: %w", err)
	}

	verifregStateV9 := verifreg9.State{
		RootKey:                  verifregStateV8.RootKey,
		Verifiers:                verifregStateV8.Verifiers,
		RemoveDataCapProposalIDs: verifregStateV8.RemoveDataCapProposalIDs,
		Allocations:              allocationsMapRoot,
		NextAllocationId:         1,
		Claims:                   emptyMapCid,
	}

	verifregHead, err := adtStore.Put(ctx, &verifregStateV9)
	if err != nil {
		return cid.Undef, nil, xerrors.Errorf("failed to put verifreg9 state: %w", err)
	}

	return verifregHead, dealsToAllocations, nil
}
