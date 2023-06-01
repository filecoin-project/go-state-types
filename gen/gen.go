package main

import (
	gen "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
)

func main() {
	// Common types
	if err := gen.WriteTupleEncodersToFile("./abi/cbor_gen.go", "abi",
		abi.PieceInfo{},
		abi.SectorID{},
		abi.AddrPairKey{},
	); err != nil {
		panic(err)
	}
	if err := gen.WriteTupleEncodersToFile("./builtin/cbor_gen.go", "builtin",
		builtin.ActorV4{},
		builtin.ActorV5{},
	); err != nil {
		panic(err)
	}
}
