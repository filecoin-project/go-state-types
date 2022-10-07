package datacap

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
	datacap9 "github.com/filecoin-project/go-state-types/builtin/v9/datacap"
)

type State = datacap9.State
type TokenState = datacap9.TokenState

func ConstructState(store adt.Store, governor address.Address, bitWidth uint64) (*State, error) {
	return datacap9.ConstructState(store, governor, bitWidth)
}
