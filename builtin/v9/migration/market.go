package migration

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	market8 "github.com/filecoin-project/go-state-types/builtin/v8/market"
	adt8 "github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
	market9 "github.com/filecoin-project/go-state-types/builtin/v9/market"
	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	"github.com/filecoin-project/go-state-types/builtin/v9/verifreg"
	"github.com/ipfs/go-cid"
	typegen "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/xerrors"
)

func migrateMarket(ctx context.Context, adtStore adt8.Store, dealAllocationTuples []DealAllocationTuple, marketStateV8 market8.State, emptyMapCid cid.Cid) (cid.Cid, error) {
	pendingDealAllocationIdsMap, err := adt9.AsMap(adtStore, emptyMapCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to load empty map: %w", err)
	}

	for _, tuple := range dealAllocationTuples {
		cborAllocationID := typegen.CborInt(tuple.Allocation)
		if err = pendingDealAllocationIdsMap.Put(abi.UIntKey(uint64(tuple.Deal)), &cborAllocationID); err != nil {
			return cid.Undef, xerrors.Errorf("failed to populate pending deal allocations map: %w", err)
		}
	}

	pendingDealAllocationIdsMapRoot, err := pendingDealAllocationIdsMap.Root()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to flush pending deal allocations map: %w", err)
	}

	dealStates8, err := adt9.AsArray(adtStore, marketStateV8.States, market8.StatesAmtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to load v8 states array: %w", err)
	}

	emptyStatesArrayCid, err := adt9.StoreEmptyArray(adtStore, market9.StatesAmtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to create empty states array: %w", err)
	}

	dealStates9, err := adt9.AsArray(adtStore, emptyStatesArrayCid, market9.StatesAmtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to load v9 states array: %w", err)
	}

	var dealState8 market8.DealState
	err = dealStates8.ForEach(&dealState8, func(i int64) error {
		return dealStates9.Set(uint64(i), &market9.DealState{
			SectorStartEpoch: dealState8.SectorStartEpoch,
			LastUpdatedEpoch: dealState8.LastUpdatedEpoch,
			SlashEpoch:       dealState8.SlashEpoch,
			VerifiedClaim:    verifreg.NoAllocationID,
		})
	})
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to iterate over v8 states: %w", err)
	}

	dealStates9Root, err := dealStates9.Root()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to flush dealStates 9: %w", err)
	}

	marketStateV9 := market9.State{
		Proposals:                     marketStateV8.Proposals,
		States:                        dealStates9Root,
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
