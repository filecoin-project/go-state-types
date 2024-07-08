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
		proof.SectorInfo{},
		proof.ExtendedSectorInfo{},
		proof.WinningPoStVerifyInfo{},
		proof.WindowPoStVerifyInfo{},
		proof.SealVerifyInfo{},
		proof.AggregateSealVerifyInfo{},
		proof.AggregateSealVerifyProofAndInfos{},
		proof.ReplicaUpdateInfo{},
	); err != nil {
		panic(err)
	}
}
