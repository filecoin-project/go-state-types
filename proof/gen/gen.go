package main

import (
	gen "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/go-state-types/proof"
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
