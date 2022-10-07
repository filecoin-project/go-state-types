package evm

import (
	"github.com/ipfs/go-cid"
)

type State struct {
	Bytecode      cid.Cid
	ContractState cid.Cid
}
