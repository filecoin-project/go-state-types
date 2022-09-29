package evm

type ConstructorParams struct {
	Bytecode  []byte
	InputData []byte
}

type GetStorageAtParams struct {
	StorageKey []byte
}
