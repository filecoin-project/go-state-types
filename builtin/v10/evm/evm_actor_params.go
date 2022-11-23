package evm

type ConstructorParams struct {
	Creator  []byte
	Initcode []byte
}

type GetStorageAtParams struct {
	StorageKey []byte
}
