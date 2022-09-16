package datacap

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

type State struct {
	Governor address.Address
	Token    TokenState
}

type TokenState struct {
	Supply       abi.TokenAmount
	Balances     cid.Cid // HAMT address.Address[abi.TokenAmount]
	Allowances   cid.Cid // HAMT address.Address[address.Address[abi.TokenAmount]]
	HamtBitWidth uint32
}

func ConstructState(store adt.Store, governor address.Address, bitWidth uint32) (*State, error) {
	emptyMapCid, err := adt.StoreEmptyMap(store, int(bitWidth))
	if err != nil {
		return nil, xerrors.Errorf("failed to create empty map: %w", err)
	}

	return &State{
		Governor: governor,
		Token: TokenState{
			Supply:       big.Zero(),
			Balances:     emptyMapCid,
			Allowances:   emptyMapCid,
			HamtBitWidth: bitWidth,
		},
	}, nil
}
