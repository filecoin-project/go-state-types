package system

import (
	system9 "github.com/filecoin-project/go-state-types/builtin/v9/system"

	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
)

type State = system9.State

func ConstructState(store adt.Store) (*State, error) {
	return system9.ConstructState(store)
}
