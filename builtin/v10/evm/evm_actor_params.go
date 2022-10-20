package evm

type ConstructorParams struct {
	Initcode []byte
}

type GetStorageAtParams struct {
	StorageKey []byte
}
