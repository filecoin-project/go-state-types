package evm

type ConstructorParams struct {
	Bytecode  []byte
	InputData []byte
}

type InvokeParams struct {
	InputData []byte
}
