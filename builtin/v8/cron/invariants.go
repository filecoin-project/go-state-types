package cron

import (
	"github.com/filecoin-project/go-address"

	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
)

type StateSummary struct {
	EntryCount int
}

// Checks internal invariants of cron state.
func CheckStateInvariants(st *State, _ adt.Store) (*StateSummary, *builtin.MessageAccumulator) {
	acc := &builtin.MessageAccumulator{}
	cronSummary := &StateSummary{
		EntryCount: len(st.Entries),
	}
	for i, e := range st.Entries {
		acc.Require(e.Receiver.Protocol() == address.ID, "entry %d receiver address %v must be ID protocol", i, e.Receiver)
		acc.Require(e.MethodNum > 0, "entry %d has invalid method number %d", i, e.MethodNum)
	}
	return cronSummary, acc
}
