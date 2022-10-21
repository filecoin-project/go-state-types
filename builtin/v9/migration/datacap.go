package migration

import (
	"context"

	verifreg8 "github.com/filecoin-project/go-state-types/builtin/v8/verifreg"

	cbor "github.com/ipfs/go-ipld-cbor"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	adt8 "github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
	datacap9 "github.com/filecoin-project/go-state-types/builtin/v9/datacap"
	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	verifreg9 "github.com/filecoin-project/go-state-types/builtin/v9/verifreg"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

type datacapMigrator struct {
	emptyMapCid             cid.Cid
	verifregStateV8         verifreg8.State
	OutCodeCID              cid.Cid
	pendingVerifiedDealSize uint64
}

func (d *datacapMigrator) migratedCodeCID() cid.Cid {
	return d.OutCodeCID

}
func (d *datacapMigrator) migrateState(ctx context.Context, store cbor.IpldStore, input actorMigrationInput) (result *actorMigrationResult, err error) {
	adtStore := adt9.WrapStore(ctx, store)
	verifiedClients, err := adt8.AsMap(adtStore, d.verifregStateV8.VerifiedClients, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to get verified clients: %w", err)
	}

	tokenSupply := big.Zero()

	balancesMap, err := adt9.AsMap(adtStore, d.emptyMapCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to load empty map: %w", err)
	}

	allowancesMap, err := adt9.AsMap(adtStore, d.emptyMapCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to load empty map: %w", err)
	}

	var dcap abi.StoragePower
	if err = verifiedClients.ForEach(&dcap, func(key string) error {
		a, err := address.NewFromBytes([]byte(key))
		if err != nil {
			return err
		}

		tokenAmount := big.Mul(dcap, verifreg9.DataCapGranularity)
		tokenSupply = big.Add(tokenSupply, tokenAmount)
		if err = balancesMap.Put(abi.IdAddrKey(a), &tokenAmount); err != nil {
			return xerrors.Errorf("failed to put new balancesMap entry: %w", err)
		}

		allowancesMapEntry, err := adt9.AsMap(adtStore, d.emptyMapCid, builtin.DefaultHamtBitwidth)
		if err != nil {
			return xerrors.Errorf("failed to load empty map: %w", err)
		}

		if err = allowancesMapEntry.Put(abi.IdAddrKey(builtin.StorageMarketActorAddr), &datacap9.InfiniteAllowance); err != nil {
			return xerrors.Errorf("failed to populate allowance map: %w", err)
		}

		return allowancesMap.Put(abi.IdAddrKey(a), allowancesMapEntry)
	}); err != nil {
		return nil, xerrors.Errorf("failed to loop over verified clients: %w", err)
	}
	verifregBalance := big.Mul(big.NewIntUnsigned(d.pendingVerifiedDealSize), verifreg9.DataCapGranularity)
	tokenSupply = big.Add(tokenSupply, verifregBalance)
	if err = balancesMap.Put(abi.IdAddrKey(builtin.VerifiedRegistryActorAddr), &verifregBalance); err != nil {
		return nil, xerrors.Errorf("failed to put verifreg balance in balancesMap: %w", err)
	}

	balancesMapRoot, err := balancesMap.Root()
	if err != nil {
		return nil, xerrors.Errorf("failed to flush balances map: %w", err)
	}

	allowancesMapRoot, err := allowancesMap.Root()
	if err != nil {
		return nil, xerrors.Errorf("failed to flush allowances map: %w", err)
	}

	dataCapState := datacap9.State{
		Governor: builtin.VerifiedRegistryActorAddr,
		Token: datacap9.TokenState{
			Supply:       tokenSupply,
			Balances:     balancesMapRoot,
			Allowances:   allowancesMapRoot,
			HamtBitWidth: builtin.DefaultHamtBitwidth,
		},
	}

	dataCapHead, err := adtStore.Put(ctx, &dataCapState)
	if err != nil {
		return nil, xerrors.Errorf("failed to put data cap state: %w", err)
	}

	return &actorMigrationResult{
		newCodeCID: d.OutCodeCID,
		newHead:    dataCapHead,
	}, nil

}
