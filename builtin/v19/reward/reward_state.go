package reward

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v19/util/adt"
	"github.com/filecoin-project/go-state-types/builtin/v19/util/smoothing"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

// A quantity of space * time (in byte-epochs) representing power committed to the network for some duration.
type Spacetime = big.Int

// 36.266260308195979333 FIL
// https://www.wolframalpha.com/input/?i=IntegerPart%5B330%2C000%2C000+*+%281+-+Exp%5B-Log%5B2%5D+%2F+%286+*+%281+year+%2F+30+seconds%29%29%5D%29+*+10%5E18%5D
const InitialRewardPositionEstimateStr = "36266260308195979333"

var InitialRewardPositionEstimate = big.MustFromString(InitialRewardPositionEstimateStr)

// -1.0982489*10^-7 FIL per epoch. Change of simple minted tokens between epochs 0 and 1.
// https://www.wolframalpha.com/input/?i=IntegerPart%5B%28Exp%5B-Log%5B2%5D+%2F+%286+*+%281+year+%2F+30+seconds%29%29%5D+-+1%29+*+10%5E18%5D
var InitialRewardVelocityEstimate = abi.NewTokenAmount(-109897758509)

// State is the reward actor state.
//
// Changed since v18:
// - TotalStoragePowerReward is renamed TotalMintedReward and counts all block-reward minting.
// - SimpleTotal and BaselineTotal move from state to code constants.
// - Stream accounting, SWA configuration, and offboarded stream state are added.
type State struct {
	// CumsumBaseline is the target CumsumRealized must reach for EffectiveNetworkTime to increase.
	// It is expressed in byte-epochs.
	CumsumBaseline Spacetime
	// CumsumRealized is the cumulative network power capped by BaselinePower(epoch).
	// It is expressed in byte-epochs.
	CumsumRealized Spacetime
	// EffectiveNetworkTime is the ceiling of real effective network time theta, based on
	// CumsumBaselinePower(theta) == CumsumRealizedPower. Theta captures how far the
	// network has progressed against its baseline and in advancing network time.
	EffectiveNetworkTime abi.ChainEpoch
	// EffectiveBaselinePower is the baseline power at EffectiveNetworkTime.
	EffectiveBaselinePower abi.StoragePower
	// ThisEpochReward is the reward per WinCount paid to block producers.
	// The actual total depends on the number of winners in the round.
	// It is recomputed every non-null epoch and used in the next non-null epoch.
	ThisEpochReward abi.TokenAmount
	// ThisEpochRewardSmoothed is the smoothed ThisEpochReward.
	ThisEpochRewardSmoothed smoothing.FilterEstimate
	// ThisEpochBaselinePower is the baseline power targeted at Epoch.
	ThisEpochBaselinePower abi.StoragePower
	// Epoch identifies the epoch for which the reward was computed.
	Epoch abi.ChainEpoch
	// TotalMintedReward is the total FIL minted through block rewards.
	TotalMintedReward abi.TokenAmount
	// TotalBurnMinted is the cumulative block-reward residual sent to the burnt funds actor.
	TotalBurnMinted abi.TokenAmount
	// TotalExplicitMinted is the cumulative block reward accrued to explicit streams.
	TotalExplicitMinted abi.TokenAmount
	// Accrued holds current-period accrual for each explicit stream, ordered by stream ID.
	Accrued []StreamAccrual
	// SWATimelockEpochs is the hold applied to SWA writes. Construction leaves zero;
	// the activation migration sets the operational value.
	SWATimelockEpochs abi.ChainEpoch
	// SWAActor manages stream configuration. Construction uses f00;
	// the activation migration sets the operational address.
	SWAActor address.Address
	// StreamsRoot references offboarded stream, tombstone, and queued-write state.
	StreamsRoot cid.Cid
}

func ConstructState(store adt.Store, currRealizedPower abi.StoragePower) (*State, error) {
	streamsRoot, err := store.Put(store.Context(), &StreamsState{})
	if err != nil {
		return nil, xerrors.Errorf("failed to store empty streams state: %w", err)
	}
	swaActor, err := address.NewIDAddress(0)
	if err != nil {
		return nil, xerrors.Errorf("failed to construct default SWA address: %w", err)
	}
	st := &State{
		CumsumBaseline:          big.Zero(),
		CumsumRealized:          big.Zero(),
		EffectiveBaselinePower:  BaselineInitialValue,
		ThisEpochReward:         big.Zero(),
		ThisEpochRewardSmoothed: smoothing.NewEstimate(InitialRewardPositionEstimate, InitialRewardVelocityEstimate),
		ThisEpochBaselinePower:  InitBaselinePower(),
		Epoch:                   -1,
		TotalMintedReward:       big.Zero(),
		TotalBurnMinted:         big.Zero(),
		TotalExplicitMinted:     big.Zero(),
		Accrued:                 []StreamAccrual{},
		SWAActor:                swaActor,
		StreamsRoot:             streamsRoot,
	}
	st.updateToNextEpochWithReward(currRealizedPower)
	return st, nil
}

// updateToNextEpoch updates internal state for the current realized power.
// It is used during null-round processing.
func (st *State) updateToNextEpoch(currRealizedPower abi.StoragePower) {
	st.Epoch++
	st.ThisEpochBaselinePower = BaselinePowerFromPrev(st.ThisEpochBaselinePower)
	cappedRealizedPower := big.Min(st.ThisEpochBaselinePower, currRealizedPower)
	st.CumsumRealized = big.Add(st.CumsumRealized, cappedRealizedPower)
	for st.CumsumRealized.GreaterThan(st.CumsumBaseline) {
		st.EffectiveNetworkTime++
		st.EffectiveBaselinePower = BaselinePowerFromPrev(st.EffectiveBaselinePower)
		st.CumsumBaseline = big.Add(st.CumsumBaseline, st.EffectiveBaselinePower)
	}
}

// updateToNextEpochWithReward advances internal state and computes the next epoch's reward.
func (st *State) updateToNextEpochWithReward(currRealizedPower abi.StoragePower) {
	prevRewardTheta := ComputeRTheta(st.EffectiveNetworkTime, st.EffectiveBaselinePower, st.CumsumRealized, st.CumsumBaseline)
	st.updateToNextEpoch(currRealizedPower)
	currRewardTheta := ComputeRTheta(st.EffectiveNetworkTime, st.EffectiveBaselinePower, st.CumsumRealized, st.CumsumBaseline)
	st.ThisEpochReward = computeReward(st.Epoch, prevRewardTheta, currRewardTheta)
}
