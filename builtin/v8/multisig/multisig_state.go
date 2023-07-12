package multisig

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs/go-cid"
)

type State struct {
	Signers               []address.Address // Signers must be canonical ID-addresses.
	NumApprovalsThreshold uint64
	NextTxnID             TxnID

	// Linear unlock
	InitialBalance abi.TokenAmount
	StartEpoch     abi.ChainEpoch
	UnlockDuration abi.ChainEpoch

	PendingTxns cid.Cid // HAMT[TxnID]Transaction
}

func (st *State) AmountLocked(elapsedEpoch abi.ChainEpoch) abi.TokenAmount {
	if elapsedEpoch >= st.UnlockDuration {
		return abi.NewTokenAmount(0)
	}
	if elapsedEpoch <= 0 {
		return st.InitialBalance
	}

	unlockDuration := big.NewInt(int64(st.UnlockDuration))
	remainingLockDuration := big.Sub(unlockDuration, big.NewInt(int64(elapsedEpoch)))

	// locked = ceil(InitialBalance * remainingLockDuration / UnlockDuration)
	numerator := big.Mul(st.InitialBalance, remainingLockDuration)
	denominator := unlockDuration
	quot := big.Div(numerator, denominator)
	rem := big.Mod(numerator, denominator)

	locked := quot
	if !rem.IsZero() {
		locked = big.Add(locked, big.NewInt(1))
	}
	return locked
}
