package evm

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
)

type ConstructorParams struct {
	Creator  [20]byte
	Initcode []byte
}

type ResurrectParams = ConstructorParams

type GetStorageAtParams struct {
	StorageKey [32]byte
}

type DelegateCallParams struct {
	Code   cid.Cid
	Input  []byte
	Caller [20]byte
	Value  abi.TokenAmount
}
