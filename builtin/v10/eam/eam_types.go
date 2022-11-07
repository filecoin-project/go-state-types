package eam

import "github.com/filecoin-project/go-address"

type CreateParams struct {
	Initcode []byte
	Nonce    uint64
}

type Create2Params struct {
	Initcode []byte
	Salt     [32]byte
}

type CreateAccount struct {
	// Pubkey is the secp256k1 public key
	Pubkey [65]byte
}

type Return struct {
	ActorID       uint64
	RobustAddress address.Address
	EthAddress    [20]byte
}

type CreateReturn Return
type Create2Return Return
