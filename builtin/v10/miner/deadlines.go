package miner

import (
	miner9 "github.com/filecoin-project/go-state-types/builtin/v9/miner"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
	"github.com/filecoin-project/go-state-types/dline"
)

func NewDeadlineInfo(periodStart abi.ChainEpoch, deadlineIdx uint64, currEpoch abi.ChainEpoch) *dline.Info {
	return miner9.NewDeadlineInfo(periodStart, deadlineIdx, currEpoch)
}

func QuantSpecForDeadline(di *dline.Info) builtin.QuantSpec {
	return miner9.QuantSpecForDeadline(di)
}

func FindSector(store adt.Store, deadlines *Deadlines, sectorNum abi.SectorNumber) (uint64, uint64, error) {
	return miner9.FindSector(store, deadlines, sectorNum)
}
