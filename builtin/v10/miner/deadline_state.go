package miner

import (
	miner9 "github.com/filecoin-project/go-state-types/builtin/v9/miner"

	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
	"github.com/ipfs/go-cid"
)

type Deadlines = miner9.Deadlines
type Deadline = miner9.Deadline
type WindowedPoSt = miner9.WindowedPoSt

const DeadlinePartitionsAmtBitwidth = miner9.DeadlinePartitionsAmtBitwidth
const DeadlineExpirationAmtBitwidth = miner9.DeadlineExpirationAmtBitwidth
const DeadlineOptimisticPoStSubmissionsAmtBitwidth = miner9.DeadlineOptimisticPoStSubmissionsAmtBitwidth

func ConstructDeadline(store adt.Store) (*Deadline, error) {
	return miner9.ConstructDeadline(store)
}
func ConstructDeadlines(emptyDeadlineCid cid.Cid) *Deadlines {
	return miner9.ConstructDeadlines(emptyDeadlineCid)
}
