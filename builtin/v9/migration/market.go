package migration

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	market8 "github.com/filecoin-project/go-state-types/builtin/v8/market"
	adt8 "github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
	market9 "github.com/filecoin-project/go-state-types/builtin/v9/market"
	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	verifreg9 "github.com/filecoin-project/go-state-types/builtin/v9/verifreg"
	"github.com/ipfs/go-cid"
	typegen "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/xerrors"
)

func migrateMarket(ctx context.Context, adtStore adt8.Store, dealsToAllocations map[abi.DealID]verifreg9.AllocationId, marketStateV8 market8.State, emptyMapCid cid.Cid) (cid.Cid, error) {
	pendingDealAllocationIdsMap, err := adt9.AsMap(adtStore, emptyMapCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to load empty map: %w", err)
	}

	for dealID, allocationID := range dealsToAllocations {
		cborAllocationID := typegen.CborInt(allocationID)
		if err = pendingDealAllocationIdsMap.Put(abi.UIntKey(uint64(dealID)), &cborAllocationID); err != nil {
			return cid.Undef, xerrors.Errorf("failed to populate pending deal allocations map: %w", err)
		}
	}

	pendingDealAllocationIdsMapRoot, err := pendingDealAllocationIdsMap.Root()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to flush pending deal allocations map: %w", err)
	}

	marketStateV9 := market9.State{
		Proposals:                     marketStateV8.Proposals,
		States:                        marketStateV8.States,
		PendingProposals:              marketStateV8.PendingProposals,
		EscrowTable:                   marketStateV8.EscrowTable,
		LockedTable:                   marketStateV8.LockedTable,
		NextID:                        marketStateV8.NextID,
		DealOpsByEpoch:                marketStateV8.DealOpsByEpoch,
		LastCron:                      marketStateV8.LastCron,
		TotalClientLockedCollateral:   marketStateV8.TotalClientLockedCollateral,
		TotalProviderLockedCollateral: marketStateV8.TotalProviderLockedCollateral,
		TotalClientStorageFee:         marketStateV8.TotalClientStorageFee,
		PendingDealAllocationIds:      pendingDealAllocationIdsMapRoot,
	}

	marketHead, err := adtStore.Put(ctx, &marketStateV9)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to put market state: %w", err)
	}

	return marketHead, nil
}
