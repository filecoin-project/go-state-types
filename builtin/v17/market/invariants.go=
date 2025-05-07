package market

import (
	"bytes"
	"encoding/binary"

	"github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v16/util/adt"
	"github.com/filecoin-project/go-state-types/builtin/v16/verifreg"
)

type DealSummary struct {
	Provider         address.Address
	StartEpoch       abi.ChainEpoch
	EndEpoch         abi.ChainEpoch
	SectorStartEpoch abi.ChainEpoch
	LastUpdatedEpoch abi.ChainEpoch
	SlashEpoch       abi.ChainEpoch
	// Populated from market's States, not the miner state
	SectorNumber abi.SectorNumber
	PieceCid     cid.Cid
}

type StateSummary struct {
	Deals                    map[abi.DealID]*DealSummary
	PendingDealAllocationIds map[abi.DealID]verifreg.AllocationId
	ClaimIdToDealId          map[verifreg.ClaimId]abi.DealID
	AllocIdToDealId          map[verifreg.AllocationId]abi.DealID
	ProviderSectors          map[abi.SectorID][]abi.DealID
	PendingProposalCount     uint64
	DealStateCount           uint64
	LockTableCount           uint64
	DealOpEpochCount         uint64
	DealOpCount              uint64
}

// Checks internal invariants of market state.
func CheckStateInvariants(st *State, store adt.Store, balance abi.TokenAmount, currEpoch abi.ChainEpoch) (*StateSummary, *builtin.MessageAccumulator) {
	acc := &builtin.MessageAccumulator{}

	acc.Require(
		st.TotalClientLockedCollateral.GreaterThanEqual(big.Zero()),
		"negative total client locked collateral: %v", st.TotalClientLockedCollateral)

	acc.Require(
		st.TotalProviderLockedCollateral.GreaterThanEqual(big.Zero()),
		"negative total provider locked collateral: %v", st.TotalClientLockedCollateral)

	acc.Require(
		st.TotalClientStorageFee.GreaterThanEqual(big.Zero()),
		"negative total client storage fee: %v", st.TotalClientLockedCollateral)

	//
	// Proposals
	//

	proposalCids := make(map[cid.Cid]struct{})
	maxDealID := int64(-1)
	proposalStats := make(map[abi.DealID]*DealSummary)
	expectedDealOps := make(map[abi.DealID]struct{})
	totalProposalCollateral := abi.NewTokenAmount(0)

	if proposals, err := adt.AsArray(store, st.Proposals, ProposalsAmtBitwidth); err != nil {
		acc.Addf("error loading proposals: %v", err)
	} else {
		var proposal DealProposal
		err = proposals.ForEach(&proposal, func(dealID int64) error {
			pcid, err := proposal.Cid()
			acc.RequireNoError(err, "error getting cid from proposal")

			if proposal.StartEpoch >= currEpoch {
				expectedDealOps[abi.DealID(dealID)] = struct{}{}
			}

			// keep some state
			proposalCids[pcid] = struct{}{}
			if dealID > maxDealID {
				maxDealID = dealID
			}
			proposalStats[abi.DealID(dealID)] = &DealSummary{
				Provider:         proposal.Provider,
				StartEpoch:       proposal.StartEpoch,
				EndEpoch:         proposal.EndEpoch,
				SectorStartEpoch: abi.ChainEpoch(-1),
				LastUpdatedEpoch: abi.ChainEpoch(-1),
				SlashEpoch:       abi.ChainEpoch(-1),
				PieceCid:         proposal.PieceCID,
			}

			totalProposalCollateral = big.Sum(totalProposalCollateral, proposal.ClientCollateral, proposal.ProviderCollateral)

			acc.Require(proposal.Client.Protocol() == address.ID, "client address for deal %d is not an ID address", dealID)
			acc.Require(proposal.Provider.Protocol() == address.ID, "provider address for deal %d is not an ID address", dealID)
			return nil
		})
		acc.RequireNoError(err, "error iterating proposals")
	}

	// next id should be higher than any existing deal
	acc.Require(int64(st.NextID) > maxDealID, "next id, %d, is not greater than highest id in proposals, %d", st.NextID, maxDealID)

	//
	// Deal States
	//

	pendingDealAllocationIds, err := st.GetPendingDealAllocationIds(store)
	acc.RequireNoError(err, "error loading pending deal proposal Ids")

	allocationIdToDealId := make(map[verifreg.AllocationId]abi.DealID)
	for dealId, allocationId := range pendingDealAllocationIds {
		_, found := proposalStats[dealId]
		acc.Require(found, "pending deal allocation %d not found in proposals", dealId)

		allocationIdToDealId[allocationId] = dealId
	}

	dealStateCount := uint64(0)
	claimIdToDealId := make(map[verifreg.ClaimId]abi.DealID)
	expectedProviderSectors := make(map[abi.DealID]struct{})
	if dealStates, err := adt.AsArray(store, st.States, StatesAmtBitwidth); err != nil {
		acc.Addf("error loading deal states: %v", err)
	} else {
		var dealState DealState
		err = dealStates.ForEach(&dealState, func(dealID int64) error {
			acc.Require(
				dealState.SectorStartEpoch >= 0,
				"deal %d state start epoch undefined: %v", dealID, dealState)

			acc.Require(
				dealState.LastUpdatedEpoch == EpochUndefined || dealState.LastUpdatedEpoch >= dealState.SectorStartEpoch,
				"deal %d state last updated before sector start: %v", dealID, dealState)

			acc.Require(
				dealState.LastUpdatedEpoch == EpochUndefined || dealState.LastUpdatedEpoch <= currEpoch,
				"deal %d last updated epoch %d after current %d", dealID, dealState.LastUpdatedEpoch, currEpoch)

			acc.Require(
				dealState.SlashEpoch == EpochUndefined || dealState.SlashEpoch >= dealState.SectorStartEpoch,
				"deal %d state slashed before sector start: %v", dealID, dealState)

			acc.Require(
				dealState.SlashEpoch == EpochUndefined || dealState.SlashEpoch <= currEpoch,
				"deal %d state slashed after current epoch %d: %v", dealID, currEpoch, dealState)

			stats, found := proposalStats[abi.DealID(dealID)]
			if !found {
				acc.Addf("no deal proposal for deal state %d", dealID)
			} else {
				stats.SectorStartEpoch = dealState.SectorStartEpoch
				stats.LastUpdatedEpoch = dealState.LastUpdatedEpoch
				stats.SlashEpoch = dealState.SlashEpoch
				stats.SectorNumber = dealState.SectorNumber
			}
			_, found = pendingDealAllocationIds[abi.DealID(dealID)]
			acc.Require(!found, "deal %d has pending allocation", dealID)

			if dealState.SlashEpoch == EpochUndefined && dealState.SectorStartEpoch != EpochUndefined && stats.EndEpoch > currEpoch {
				expectedProviderSectors[abi.DealID(dealID)] = struct{}{}
			}

			dealStateCount++

			return nil
		})
		acc.RequireNoError(err, "error iterating deal states")
	}

	//
	// Pending Proposals
	//

	pendingProposalCount := uint64(0)
	if pendingProposals, err := adt.AsMap(store, st.PendingProposals, builtin.DefaultHamtBitwidth); err != nil {
		acc.Addf("error loading pending proposals: %v", err)
	} else {
		err = pendingProposals.ForEach(nil, func(key string) error {
			proposalCID, err := cid.Parse([]byte(key))
			acc.RequireNoError(err, "error getting cid from proposal")

			_, found := proposalCids[proposalCID]
			acc.Require(found, "pending proposal with cid %v not found within proposals %v", proposalCID, pendingProposals)

			pendingProposalCount++
			return nil
		})
		acc.RequireNoError(err, "error iterating pending proposals")
	}

	//
	// Escrow Table and Locked Table
	//

	lockTableCount := uint64(0)
	escrowTable, err := adt.AsBalanceTable(store, st.EscrowTable)
	acc.RequireNoError(err, "error loading escrow table")
	lockTable, err := adt.AsBalanceTable(store, st.LockedTable)
	acc.RequireNoError(err, "error loading locked table")
	if escrowTable != nil && lockTable != nil {
		var lockedAmount abi.TokenAmount
		lockedTotal := abi.NewTokenAmount(0)
		err = (*adt.Map)(lockTable).ForEach(&lockedAmount, func(key string) error {
			addr, err := address.NewFromBytes([]byte(key))
			acc.RequireNoError(err, "error getting address from bytes")
			lockedTotal = big.Add(lockedTotal, lockedAmount)

			// every entry in locked table should have a corresponding entry in escrow table that is at least as high
			escrowAmount, err := escrowTable.Get(addr)
			acc.RequireNoError(err, "error escrow amount from table for %s", addr)
			acc.Require(escrowAmount.GreaterThanEqual(lockedAmount),
				"locked funds for %s, %s, greater than escrow amount, %s", addr, lockedAmount, escrowAmount)

			lockTableCount++
			return nil
		})
		acc.RequireNoError(err, "error iterating locked table")

		// lockTable total should be sum of client and provider locked plus client storage fee
		expectedLockTotal := big.Sum(st.TotalProviderLockedCollateral, st.TotalClientLockedCollateral, st.TotalClientStorageFee)
		acc.Require(lockedTotal.Equals(expectedLockTotal),
			"locked total, %s, does not sum to provider locked, %s, client locked, %s, and client storage fee, %s",
			lockedTotal, st.TotalProviderLockedCollateral, st.TotalClientLockedCollateral, st.TotalClientStorageFee)

		// assert escrow <= actor balance
		// lockTable item <= escrow item and escrowTotal <= balance implies lockTable total <= balance
		escrowTotal, err := escrowTable.Total()
		acc.RequireNoError(err, "error calculating escrow total")
		acc.Require(escrowTotal.LessThanEqual(balance), "escrow total, %v, greater than actor balance, %v", escrowTotal, balance)
		acc.Require(escrowTotal.GreaterThanEqual(totalProposalCollateral), "escrow total, %v, less than sum of proposal collateral, %v", escrowTotal, totalProposalCollateral)
	}

	//
	// Deal Ops by Epoch
	//

	dealOpEpochCount := uint64(0)
	dealOpCount := uint64(0)
	if dealOps, err := AsSetMultimap(store, st.DealOpsByEpoch, builtin.DefaultHamtBitwidth, builtin.DefaultHamtBitwidth); err != nil {
		acc.Addf("error loading deal ops: %v", err)
	} else {
		// get into internals just to iterate through full data structure
		var setRoot cbg.CborCid
		err = dealOps.mp.ForEach(&setRoot, func(key string) error {
			epoch, err := binary.ReadUvarint(bytes.NewReader([]byte(key)))
			acc.RequireNoError(err, "error epoch from bytes")

			dealOpEpochCount++
			return dealOps.ForEach(abi.ChainEpoch(epoch), func(id abi.DealID) error {
				_, found := proposalStats[id]
				acc.Require(found, "deal op found for deal id %d with missing proposal at epoch %d", id, epoch)
				delete(expectedDealOps, id)
				dealOpCount++
				return nil
			})
		})
		acc.RequireNoError(err, "error iterating deal ops")
	}

	acc.Require(len(expectedDealOps) == 0, "missing deal ops for proposals: %v", expectedDealOps)

	//
	// Provider Sectors
	// Provider->sector->deal mapping
	// A deal may have expired but remain in the mapping until settlement.

	providerSectors := make(map[abi.SectorID][]abi.DealID)
	if sectorDeals, err := adt.AsMap(store, st.ProviderSectors, ProviderSectorsHamtBitwidth); err != nil {
		acc.Addf("error loading st.ProviderSectors: %v", err)
	} else {
		var sectorMapRoot cbg.CborCid
		err = sectorDeals.ForEach(&sectorMapRoot, func(providerID string) error {
			provider, err := abi.ParseUIntKey(providerID)
			if err != nil {
				acc.RequireNoError(err, "error getting address from bytes")
				return nil
			}

			sectorMap, err := adt.AsMap(store, cid.Cid(sectorMapRoot), ProviderSectorsHamtBitwidth)
			if err != nil {
				acc.RequireNoError(err, "error loading sector map for provider %s", provider)
				return nil
			}

			var dealIDs SectorDealIDs
			err = sectorMap.ForEach(&dealIDs, func(sectorIDBytes string) error {
				sectorNumber, err := abi.ParseUIntKey(sectorIDBytes)
				acc.RequireNoError(err, "error sector number from bytes")

				acc.Require(len(dealIDs) > 0, "no deal ids in sector %d", sectorNumber)

				dealIDsCopy := make([]abi.DealID, len(dealIDs))
				copy(dealIDsCopy, dealIDs)

				providerSectors[abi.SectorID{Miner: abi.ActorID(provider), Number: abi.SectorNumber(sectorNumber)}] = dealIDsCopy

				// check against proposalStats
				for _, dealID := range dealIDsCopy {
					st, found := proposalStats[dealID]
					if !found || st.SlashEpoch != EpochUndefined || st.EndEpoch < currEpoch {
						continue
					}

					acc.Require(st.SectorNumber == abi.SectorNumber(sectorNumber), "deal id %d sector number %d does not match sector id %d", dealID, st.SectorNumber, sectorNumber)

					_, ok := expectedProviderSectors[dealID]
					acc.Require(ok, "found unexpected deal in ProviderSectors")
					delete(expectedProviderSectors, dealID)

					provID, err := address.IDFromAddress(st.Provider)
					if err != nil {
						acc.Addf("error creating ID address: %v", err)
						continue
					}

					acc.Require(provider == provID, "deal %d has provider %v, expected %v", dealID, provID, provider)
				}

				return nil
			})
			acc.RequireNoError(err, "error iterating sector deals")
			return nil
		})
		acc.RequireNoError(err, "error iterating sector deals")
	}

	acc.Require(len(expectedProviderSectors) == 0, "missing %d providersectors entries for deals", len(expectedProviderSectors))

	return &StateSummary{
		Deals:                    proposalStats,
		PendingDealAllocationIds: pendingDealAllocationIds,
		PendingProposalCount:     pendingProposalCount,
		DealStateCount:           dealStateCount,
		LockTableCount:           lockTableCount,
		DealOpEpochCount:         dealOpEpochCount,
		DealOpCount:              dealOpCount,
		ClaimIdToDealId:          claimIdToDealId,
		AllocIdToDealId:          allocationIdToDealId,
		ProviderSectors:          providerSectors,
	}, acc
}
