package main

import (
	"github.com/filecoin-project/go-state-types/manifest"
	gen "github.com/whyrusleeping/cbor-gen"
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
