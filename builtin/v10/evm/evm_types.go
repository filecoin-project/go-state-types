package evm

import (
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-state-types/abi"
)

type ConstructorParams struct {
	Creator  []byte
	Initcode []byte
}

type ResurrectParams = ConstructorParams

type GetStorageAtParams struct {
	StorageKey []byte
}

type DelegateCallParams struct {
	Code   cid.Cid
	Input  []byte
	Caller [20]byte
	Value  abi.TokenAmount
}
