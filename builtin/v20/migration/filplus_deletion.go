package migration

import (
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"golang.org/x/xerrors"
)

func transferBalanceAndDeleteVerifregDataCap(actorsIn *builtin.ActorTree) error {
	verifregKeyAddrs, err := addr.NewIDAddress(80)
	if err != nil {
		return xerrors.Errorf("failed to create verifreg root address: %w", err)
	}

	var actorsToDelete = map[string]addr.Address{
		"Verifreg":      builtin.VerifiedRegistryActorAddr,
		"Datacap":       builtin.DatacapActorAddr,
		"Verifreg Root": verifregKeyAddrs,
	}

	act, found, err := actorsIn.GetActorV5(builtin.BurntFundsActorAddr)
	if err != nil {
		return xerrors.Errorf("failed to get burnt funds actor from store: %w", err)
	}

	if !found {
		return xerrors.Errorf("burnt funds actor not found")
	}

	for name, addrs := range actorsToDelete {
		actr, afound, err := actorsIn.GetActorV5(addrs)
		if err != nil {
			return xerrors.Errorf("failed to get %s actor from store: %w", name, err)
		}

		if afound {
			if err := actorsIn.Map.Delete(abi.AddrKey(addrs)); err != nil {
				return xerrors.Errorf("failed to delete %s actor from store: %w", name, err)
			}
			// Transfer balance to burnt funds actor
			act.Balance = big.Add(act.Balance, actr.Balance)
		} else {
			if addrs != verifregKeyAddrs {
				return xerrors.Errorf("%s actor not found", name)
			}
		}
	}

	return actorsIn.SetActorV5(builtin.BurntFundsActorAddr, act)
}
