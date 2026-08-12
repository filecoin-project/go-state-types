package reward

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v19/util/adt"
)

type StateSummary struct {
	StreamCount       int
	TombstoneCount    int
	PendingWriteCount int
}

var FIL = big.NewInt(1e18)
var StorageMiningAllocationCheck = big.Mul(big.NewInt(1_100_000_000), FIL)

// CheckStateInvariants mirrors actors/reward/src/testing.rs::check_state_invariants.
func CheckStateInvariants(st *State, store adt.Store, priorEpoch abi.ChainEpoch, balance abi.TokenAmount) (*StateSummary, *builtin.MessageAccumulator) {
	acc := &builtin.MessageAccumulator{}

	// Anyone can send funds to the reward actor, so this cannot be an equality.
	acc.Require(big.Add(st.TotalMintedReward, balance).GreaterThanEqual(StorageMiningAllocationCheck), "reward minted %v + reward left %v < storage mining allocation %v", st.TotalMintedReward, balance, StorageMiningAllocationCheck)
	acc.Require(st.Epoch == priorEpoch+1, "reward state epoch %d does not match priorEpoch+1 %d", st.Epoch, priorEpoch+1)
	acc.Require(st.EffectiveNetworkTime <= st.Epoch, "effective network time %d greater than state epoch %d", st.EffectiveNetworkTime, st.Epoch)
	acc.Require(st.CumsumRealized.LessThanEqual(st.CumsumBaseline), "cumsum realized %v > cumsum baseline %v", st.CumsumRealized, st.CumsumBaseline)
	acc.Require(st.CumsumRealized.GreaterThanEqual(big.Zero()), "cumsum realized negative (%v)", st.CumsumRealized)
	acc.Require(st.EffectiveBaselinePower.LessThanEqual(st.ThisEpochBaselinePower), "effective baseline power > baseline power")

	for _, total := range []struct {
		name   string
		amount abi.TokenAmount
	}{
		{name: "total minted reward", amount: st.TotalMintedReward},
		{name: "total burn minted", amount: st.TotalBurnMinted},
		{name: "total explicit minted", amount: st.TotalExplicitMinted},
	} {
		acc.Require(total.amount.GreaterThanEqual(big.Zero()), "%s is negative (%v)", total.name, total.amount)
	}
	acc.Require(
		big.Add(st.TotalBurnMinted, st.TotalExplicitMinted).LessThanEqual(st.TotalMintedReward),
		"burn %v + explicit %v exceeds total minted %v",
		st.TotalBurnMinted,
		st.TotalExplicitMinted,
		st.TotalMintedReward,
	)
	acc.Require(st.SWATimelockEpochs >= 0, "SWA timelock is negative (%d)", st.SWATimelockEpochs)
	acc.Require(st.SWAActor.Protocol() == address.ID, "SWA actor %v is not an ID address", st.SWAActor)
	for i, row := range st.Accrued {
		if i > 0 {
			acc.Require(st.Accrued[i-1].ID < row.ID, "explicit-stream accrual rows are not strictly ordered by stream ID")
		}
		acc.Require(row.Amount.GreaterThanEqual(big.Zero()), "explicit-stream accrual for stream %d is negative (%v)", row.ID, row.Amount)
	}

	var streams StreamsState
	if err := store.Get(store.Context(), st.StreamsRoot, &streams); err != nil {
		acc.Addf("error loading streams state: %v", err)
		return &StateSummary{}, acc
	}

	summary := &StateSummary{
		StreamCount:       len(streams.Streams),
		TombstoneCount:    len(streams.Tombstones),
		PendingWriteCount: len(streams.PendingWrites),
	}
	// Mirrors actors/reward/src/streams.rs::validate_streams_state.
	if err := validateStreamsState(&streams, st.Accrued, priorEpoch+1); err != nil {
		acc.Addf("invalid streams state: %v", err)
	}

	streamIDs := make(map[StreamID]struct{}, len(streams.Streams))
	explicitStreamIDs := make(map[StreamID]struct{})
	for i, stream := range streams.Streams {
		if i > 0 {
			acc.Require(streams.Streams[i-1].ID < stream.ID, "streams are not strictly ordered by stream ID")
		}
		streamIDs[stream.ID] = struct{}{}
		if stream.Distribution != nil {
			explicitStreamIDs[stream.ID] = struct{}{}
		}
	}
	tombstoneIDs := make(map[StreamID]struct{}, len(streams.Tombstones))
	for i, tombstone := range streams.Tombstones {
		if i > 0 {
			acc.Require(streams.Tombstones[i-1].ID < tombstone.ID, "tombstones are not strictly ordered by stream ID")
		}
		tombstoneIDs[tombstone.ID] = struct{}{}
		_, live := streamIDs[tombstone.ID]
		acc.Require(!live, "a stream ID is both live and tombstoned")
	}

	accrualIDs := make(map[StreamID]struct{}, len(st.Accrued))
	for _, row := range st.Accrued {
		accrualIDs[row.ID] = struct{}{}
	}
	missing := make([]StreamID, 0)
	unexpected := make([]StreamID, 0)
	for id := range explicitStreamIDs {
		if _, found := accrualIDs[id]; !found {
			missing = append(missing, id)
		}
	}
	for id := range accrualIDs {
		if _, found := explicitStreamIDs[id]; !found {
			unexpected = append(unexpected, id)
		}
	}
	acc.Require(len(missing) == 0 && len(unexpected) == 0, "explicit-stream accrual IDs do not match live explicit streams: missing %v, unexpected %v", missing, unexpected)

	pendingSlots := make(map[pendingSlot]struct{}, len(streams.PendingWrites))
	for i, write := range streams.PendingWrites {
		slot := slotForWrite(write)
		_, duplicate := pendingSlots[slot]
		acc.Require(!duplicate, "duplicate pending slot (%v, %d)", write.ID, write.Op)
		pendingSlots[slot] = struct{}{}
		if i > 0 {
			acc.Require(streams.PendingWrites[i-1].EffectiveEpoch <= write.EffectiveEpoch, "pending writes are not ordered by effective epoch")
		}
	}

	// Mirrors actors/reward/src/streams.rs::compute_service_liability.
	liabilities, err := computeExplicitLiability(&streams, st.Accrued)
	if err != nil {
		acc.Addf("error computing explicit-stream liabilities: %v", err)
	} else {
		acc.Require(liabilities.LessThanEqual(st.TotalExplicitMinted), "explicit-stream liabilities %v exceed total explicit minted %v", liabilities, st.TotalExplicitMinted)
		acc.Require(balance.GreaterThanEqual(liabilities), "reward balance %v does not cover explicit-stream liabilities %v", balance, liabilities)
	}

	return summary, acc
}
