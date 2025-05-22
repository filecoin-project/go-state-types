package power

import (
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"

	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-hamt-ipld/v3"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v17/util/adt"
	"github.com/filecoin-project/go-state-types/builtin/v17/util/smoothing"
)

// genesis power in bytes = 750,000 GiB
var InitialQAPowerEstimatePosition = big.Mul(big.NewInt(750_000), big.NewInt(1<<30))

// max chain throughput in bytes per epoch = 120 ProveCommits / epoch = 3,840 GiB
var InitialQAPowerEstimateVelocity = big.Mul(big.NewInt(3_840), big.NewInt(1<<30))

// Bitwidth of CronEventQueue HAMT determined empirically from mutation
// patterns and projections of mainnet data.
const CronQueueHamtBitwidth = 6

// Bitwidth of CronEventQueue AMT determined empirically from mutation
// patterns and projections of mainnet data.
const CronQueueAmtBitwidth = 6

// Bitwidth of ProofValidationBatch AMT determined empirically from mutation
// patterns and projections of mainnet data.
const ProofValidationBatchAmtBitwidth = 4

// The number of miners that must meet the consensus minimum miner power before that minimum power is enforced
// as a condition of leader election.
// This ensures a network still functions before any miners reach that threshold.
const ConsensusMinerMinMiners = 4 // PARAM_SPEC

// PARAM_SPEC// Maximum number of prove-commits each miner can submit in one epoch.
//
// This limits the number of proof partitions we may need to load in the cron call path.
// Onboarding 1EiB/year requires at least 32 prove-commits per epoch.
const MaxMinerProveCommitsPerEpoch = 200 // PARAM_SPEC

type State struct {
	TotalRawBytePower abi.StoragePower
	// TotalBytesCommitted includes claims from miners below min power threshold
	TotalBytesCommitted  abi.StoragePower
	TotalQualityAdjPower abi.StoragePower
	// TotalQABytesCommitted includes claims from miners below min power threshold
	TotalQABytesCommitted abi.StoragePower
	TotalPledgeCollateral abi.TokenAmount

	// These fields are set once per epoch in the previous cron tick and used
	// for consistent values across a single epoch's state transition.
	ThisEpochRawBytePower     abi.StoragePower
	ThisEpochQualityAdjPower  abi.StoragePower
	ThisEpochPledgeCollateral abi.TokenAmount
	ThisEpochQAPowerSmoothed  smoothing.FilterEstimate

	MinerCount int64
	// Number of miners having proven the minimum consensus power.
	MinerAboveMinPowerCount int64

	// FIP0081 changed pledge calculations, moving from ruleset A to ruleset B.
	// This change is spread over several epochs to avoid sharp jumps in pledge
	// amounts. At `RampStartEpoch`, we use the old ruleset. At
	// `RampStartEpoch + RampDurationEpochs`, we use 70% old rules + 30%
	// new rules. See FIP0081 for more details.
	RampStartEpoch int64
	// Number of epochs over which the new pledge calculation is ramped up.
	RampDurationEpochs uint64

	// A queue of events to be triggered by cron, indexed by epoch.
	CronEventQueue cid.Cid // Multimap, (HAMT[ChainEpoch]AMT[CronEvent])

	// First epoch in which a cron task may be stored.
	// Cron will iterate every epoch between this and the current epoch inclusively to find tasks to execute.
	FirstCronEpoch abi.ChainEpoch

	// Claimed power for each miner.
	Claims cid.Cid // Map, HAMT[address]Claim

	ProofValidationBatch *cid.Cid // Multimap, (HAMT[Address]AMT[SealVerifyInfo])
}

func ConstructState(store adt.Store) (*State, error) {
	emptyClaimsMapCid, err := adt.StoreEmptyMap(store, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to create empty map: %w", err)
	}
	emptyCronQueueMMapCid, err := adt.StoreEmptyMultimap(store, CronQueueHamtBitwidth, CronQueueAmtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to create empty multimap: %w", err)
	}

	return &State{
		TotalRawBytePower:         abi.NewStoragePower(0),
		TotalBytesCommitted:       abi.NewStoragePower(0),
		TotalQualityAdjPower:      abi.NewStoragePower(0),
		TotalQABytesCommitted:     abi.NewStoragePower(0),
		TotalPledgeCollateral:     abi.NewTokenAmount(0),
		ThisEpochRawBytePower:     abi.NewStoragePower(0),
		ThisEpochQualityAdjPower:  abi.NewStoragePower(0),
		ThisEpochPledgeCollateral: abi.NewTokenAmount(0),
		ThisEpochQAPowerSmoothed:  smoothing.NewEstimate(InitialQAPowerEstimatePosition, InitialQAPowerEstimateVelocity),
		FirstCronEpoch:            0,
		CronEventQueue:            emptyCronQueueMMapCid,
		Claims:                    emptyClaimsMapCid,
		MinerCount:                0,
		MinerAboveMinPowerCount:   0,
		RampStartEpoch:            0,
		RampDurationEpochs:        0,
	}, nil
}

type Claim struct {
	// Miner's proof type used to determine minimum miner size
	WindowPoStProofType abi.RegisteredPoStProof

	// Sum of raw byte power for a miner's sectors.
	RawBytePower abi.StoragePower

	// Sum of quality adjusted power for a miner's sectors.
	QualityAdjPower abi.StoragePower
}

type CronEvent struct {
	MinerAddr       addr.Address
	CallbackPayload []byte
}

// ClaimMeetsConsensusMinimums checks if given claim meets the minimums set by the network for mining.
func (st *State) ClaimMeetsConsensusMinimums(claim *Claim) (bool, error) {
	minerNominalPower := claim.RawBytePower
	minerMinPower, err := builtin.ConsensusMinerMinPower(claim.WindowPoStProofType)
	if err != nil {
		return false, xerrors.Errorf("could not get miner min power from proof type: %w", err)
	}

	// if miner is larger than min power requirement, we're set
	if minerNominalPower.GreaterThanEqual(minerMinPower) {
		return true, nil
	}

	// otherwise, if ConsensusMinerMinMiners miners meet min power requirement, return false
	if st.MinerAboveMinPowerCount >= ConsensusMinerMinMiners {
		return false, nil
	}

	// If fewer than ConsensusMinerMinMiners over threshold miner can win a block with non-zero power
	return minerNominalPower.GreaterThan(abi.NewStoragePower(0)), nil
}

type powerMapReduceCache struct {
	cmr *hamt.CachedMapReduce[Claim, *Claim, []builtin.OwnedClaim]
}

func (st *State) CollectEligibleClaims(s adt.Store, cacheInOut *builtin.MapReduceCache) ([]builtin.OwnedClaim, error) {
	if st.MinerAboveMinPowerCount < ConsensusMinerMinMiners {
		// simple collect all claims,
		var res []builtin.OwnedClaim
		claims, err := adt.AsMap(s, st.Claims, builtin.DefaultHamtBitwidth)
		if err != nil {
			return nil, xerrors.Errorf("failed to load claims: %w", err)
		}
		var out Claim
		err = claims.ForEach(&out, func(k string) error {
			if !out.RawBytePower.GreaterThan(abi.NewStoragePower(0)) {
				return nil
			}
			addr, err := addr.NewFromBytes([]byte(k))
			if err != nil {
				return xerrors.Errorf("parsing address from bytes: %w", err)
			}
			res = append(res, builtin.OwnedClaim{
				Address:         addr,
				RawBytePower:    out.RawBytePower,
				QualityAdjPower: out.QualityAdjPower,
			})
			return nil
		})
		return res, err
	}
	cache, ok := (*cacheInOut).(powerMapReduceCache)
	if !ok {
		mapper := func(k string, claim Claim) ([]builtin.OwnedClaim, error) {
			minerMinPower, err := builtin.ConsensusMinerMinPower(claim.WindowPoStProofType)
			if err != nil {
				return nil, xerrors.Errorf("could not get miner min power from proof type: %w", err)
			}
			if !claim.RawBytePower.GreaterThanEqual(minerMinPower) {
				return nil, nil
			}
			addr, err := addr.NewFromBytes([]byte(k))
			if err != nil {
				return nil, err
			}
			return []builtin.OwnedClaim{
				{
					Address:         addr,
					RawBytePower:    claim.RawBytePower,
					QualityAdjPower: claim.QualityAdjPower,
				},
			}, nil
		}
		reducer := func(in [][]builtin.OwnedClaim) ([]builtin.OwnedClaim, error) {
			var out []builtin.OwnedClaim
			for _, v := range in {
				out = append(out, v...)
			}
			return out, nil
		}
		// cache size of 2000 is arbitrary, but seems to work well mith 600k claims with room for it to grow
		cmr, err := hamt.NewCachedMapReduce[Claim, *Claim, []builtin.OwnedClaim](mapper, reducer, 2000)
		if err != nil {
			return nil, err
		}

		cache = powerMapReduceCache{
			cmr: cmr,
		}
		(*cacheInOut) = cache
	}

	claims, err := cache.cmr.MapReduce(s.Context(), s, st.Claims,
		hamt.UseTreeBitWidth(builtin.DefaultHamtBitwidth))
	if err != nil {
		return nil, xerrors.Errorf("failed to map reduce claims: %w", err)
	}

	return claims, nil
}

// MinerNominalPowerMeetsConsensusMinimum is used to validate Election PoSt
// winners outside the chain state. If the miner has over a threshold of power
// the miner meets the minimum.  If the network is a below a threshold of
// miners and has power > zero the miner meets the minimum.
func (st *State) MinerNominalPowerMeetsConsensusMinimum(s adt.Store, miner addr.Address) (bool, error) { //nolint:deadcode,unused
	claims, err := adt.AsMap(s, st.Claims, builtin.DefaultHamtBitwidth)
	if err != nil {
		return false, xerrors.Errorf("failed to load claims: %w", err)
	}

	claim, ok, err := getClaim(claims, miner)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, xerrors.Errorf("no claim for actor %w", miner)
	}

	return st.ClaimMeetsConsensusMinimums(claim)
}

func (st *State) GetClaim(s adt.Store, a addr.Address) (*Claim, bool, error) {
	claims, err := adt.AsMap(s, st.Claims, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, false, xerrors.Errorf("failed to load claims: %w", err)
	}
	return getClaim(claims, a)
}

func getClaim(claims *adt.Map, a addr.Address) (*Claim, bool, error) {
	var out Claim
	found, err := claims.Get(abi.AddrKey(a), &out)
	if err != nil {
		return nil, false, xerrors.Errorf("failed to get claim for address %v: %w", a, err)
	}
	if !found {
		return nil, false, nil
	}
	return &out, true, nil
}
