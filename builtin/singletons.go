package builtin

import (
	addr "github.com/filecoin-project/go-address"
)

// Addresses for singleton system actors.
var (
	// Distinguished AccountActor that is the source of system implicit messages.
	SystemActorAddr           = mustMakeAddress(0)
	InitActorAddr             = mustMakeAddress(1)
	RewardActorAddr           = mustMakeAddress(2)
	CronActorAddr             = mustMakeAddress(3)
	StoragePowerActorAddr     = mustMakeAddress(4)
	StorageMarketActorAddr    = mustMakeAddress(5)
	VerifiedRegistryActorAddr = mustMakeAddress(6)
	DatacapActorAddr          = mustMakeAddress(7)
	// Distinguished AccountActor that is the destination of all burnt funds.
	BurntFundsActorAddr = mustMakeAddress(99)

	// EthereumAddressManagerActorID is the actor ID of the Ethereum Address Manager singleton.
	EthereumAddressManagerActorID   = uint64(10)
	EthereumAddressManagerActorAddr = mustMakeAddress(EthereumAddressManagerActorID)
)

const FirstNonSingletonActorId = 100

func mustMakeAddress(id uint64) addr.Address {
	address, err := addr.NewIDAddress(id)
	if err != nil {
		panic(err)
	}
	return address
}
