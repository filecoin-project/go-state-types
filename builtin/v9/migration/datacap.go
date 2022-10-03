package migration

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	adt8 "github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
	verifreg8 "github.com/filecoin-project/go-state-types/builtin/v8/verifreg"
	datacap9 "github.com/filecoin-project/go-state-types/builtin/v9/datacap"
	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	verifreg9 "github.com/filecoin-project/go-state-types/builtin/v9/verifreg"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

func createDatacap(ctx context.Context, adtStore adt8.Store, verifregStateV8 verifreg8.State, emptyMapCid cid.Cid) (cid.Cid, error) {

	verifiedClients, err := adt8.AsMap(adtStore, verifregStateV8.VerifiedClients, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to get verified clients: %w", err)
	}

	tokenSupply := big.Zero()

	balancesMap, err := adt9.AsMap(adtStore, emptyMapCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to load empty map: %w", err)
	}

	allowancesMap, err := adt9.AsMap(adtStore, emptyMapCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to load empty map: %w", err)
	}

	var dcap abi.StoragePower
	if err = verifiedClients.ForEach(&dcap, func(key string) error {
		a, err := address.NewFromBytes([]byte(key))
		if err != nil {
			return err
		}

		tokenAmount := verifreg9.DataCapToTokens(dcap)
		tokenSupply = big.Add(tokenSupply, tokenAmount)
		if err = balancesMap.Put(abi.IdAddrKey(a), &tokenAmount); err != nil {
			return xerrors.Errorf("failed to put new balancesMap entry: %w", err)
		}

		allowancesMapEntry, err := adt9.AsMap(adtStore, emptyMapCid, builtin.DefaultHamtBitwidth)
		if err != nil {
			return xerrors.Errorf("failed to load empty map: %w", err)
		}

		if err = allowancesMapEntry.Put(abi.IdAddrKey(builtin.StorageMarketActorAddr), &datacap9.InfiniteAllowance); err != nil {
			return xerrors.Errorf("failed to populate allowance map: %w", err)
		}

		return allowancesMap.Put(abi.IdAddrKey(a), allowancesMapEntry)
	}); err != nil {
		return cid.Undef, xerrors.Errorf("failed to loop over verified clients: %w", err)
	}

	balancesMapRoot, err := balancesMap.Root()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to flush balances map: %w", err)
	}

	allowancesMapRoot, err := allowancesMap.Root()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to flush allowances map: %w", err)
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
		return cid.Undef, xerrors.Errorf("failed to put data cap state: %w", err)
	}

	return dataCapHead, nil

}
