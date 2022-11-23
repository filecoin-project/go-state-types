package datacap

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
)

var InfiniteAllowance = big.Mul(big.MustFromString("1000000000000000000000"), builtin.TokenPrecision)

type MintParams struct {
	To        address.Address
	Amount    abi.TokenAmount
	Operators []address.Address
}

type DestroyParams struct {
	Owner  address.Address
	Amount abi.TokenAmount
}

type MintReturn struct {
	Balance       abi.TokenAmount
	Supply        abi.TokenAmount
	RecipientData []byte
}

type TransferParams struct {
	To           address.Address
	Amount       abi.TokenAmount
	OperatorData []byte
}

type TransferReturn struct {
	FromBalance   abi.TokenAmount
	ToBalance     abi.TokenAmount
	RecipientData []byte
}

type TransferFromParams struct {
	From         address.Address
	To           address.Address
	Amount       abi.TokenAmount
	OperatorData []byte
}

type TransferFromReturn struct {
	FromBalance   abi.TokenAmount
	ToBalance     abi.TokenAmount
	Allowance     abi.TokenAmount
	RecipientData []byte
}

type IncreaseAllowanceParams struct {
	Operator address.Address
	Increase abi.TokenAmount
}

type DecreaseAllowanceParams struct {
	Operator address.Address
	Decrease abi.TokenAmount
}

type RevokeAllowanceParams struct {
	Operator address.Address
}

type GetAllowanceParams struct {
	Owner    address.Address
	Operator address.Address
}

type BurnParams struct {
	Amount abi.TokenAmount
}

type BurnReturn struct {
	Balance abi.TokenAmount
}

type BurnFromParams struct {
	Owner  address.Address
	Amount abi.TokenAmount
}

type BurnFromReturn struct {
	Balance   abi.TokenAmount
	Allowance abi.TokenAmount
}
