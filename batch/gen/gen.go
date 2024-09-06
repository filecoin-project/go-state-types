package main

import (
	"github.com/filecoin-project/go-state-types/batch"
	gen "github.com/whyrusleeping/cbor-gen"
)

func main() {
	// Actor manifest
	if err := gen.WriteTupleEncodersToFile("./batch/cbor_gen.go", "batch",
		// actor manifest
		batch.BatchReturn{},
		batch.FailCode{},
	); err != nil {
		panic(err)
	}
}
