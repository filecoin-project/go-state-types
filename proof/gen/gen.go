package main

import (
	"github.com/filecoin-project/go-state-types/proof"
	gen "github.com/whyrusleeping/cbor-gen"
)

func main() {
	// Actor manifest
	if err := gen.WriteTupleEncodersToFile("./proof/cbor_gen.go", "proof",
		// actor manifest
		proof.PoStProof{},
		proof.ExtendedSectorInfo{},
		proof.SealVerifyInfo{},
		proof.WindowPoStVerifyInfo{},
		proof.WinningPoStVerifyInfo{},
		proof.SectorInfo{},
	); err != nil {
		panic(err)
	}
}
