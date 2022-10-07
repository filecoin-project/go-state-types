package cron

import (
	cron9 "github.com/filecoin-project/go-state-types/builtin/v9/cron"
)

type State = cron9.State

type Entry = cron9.Entry

func ConstructState(entries []Entry) *State {
	return cron9.ConstructState(entries)
}

// The default entries to install in the cron actor's state at genesis.
func BuiltInEntries() []Entry {
	return cron9.BuiltInEntries()
}
