package miner

import (
	miner9 "github.com/filecoin-project/go-state-types/builtin/v9/miner"

	"github.com/filecoin-project/go-state-types/abi"
)

type Partition = miner9.Partition

const PartitionExpirationAmtBitwidth = miner9.PartitionExpirationAmtBitwidth
const PartitionEarlyTerminationArrayAmtBitwidth = miner9.PartitionEarlyTerminationArrayAmtBitwidth

type PowerPair = miner9.PowerPair

func NewPowerPairZero() PowerPair {
	return miner9.NewPowerPairZero()
}

func NewPowerPair(raw, qa abi.StoragePower) PowerPair {
	return miner9.NewPowerPair(raw, qa)
}
