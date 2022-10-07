package main

import (
	evm "github.com/filecoin-project/go-state-types/builtin/v10/evm"
	init_ "github.com/filecoin-project/go-state-types/builtin/v10/init"
	gen "github.com/whyrusleeping/cbor-gen"
)

func main() {
	if err := gen.WriteTupleEncodersToFile("./builtin/v10/init/cbor_gen.go", "init",
		// actor state
		init_.State{},
		// method params and returns
		init_.ConstructorParams{},
		init_.ExecParams{},
		init_.ExecReturn{},
		init_.Exec4Params{},
		init_.InstallParams{},
		init_.InstallReturn{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v10/evm/cbor_gen.go", "evm",
		// actor state
		evm.State{},
		// method params and returns
		evm.ConstructorParams{},
		evm.GetStorageAtParams{},
	); err != nil {
		panic(err)
	}
}
