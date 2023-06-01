package main

import (
	gen "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/go-state-types/manifest"
)

func main() {
	// Actor manifest
	if err := gen.WriteTupleEncodersToFile("./manifest/cbor_gen.go", "manifest",
		// actor manifest
		manifest.Manifest{},
		manifest.ManifestEntry{},
	); err != nil {
		panic(err)
	}
}
