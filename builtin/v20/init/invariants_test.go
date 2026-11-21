package init_test

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v20/account"
	init_ "github.com/filecoin-project/go-state-types/builtin/v20/init"
	"github.com/filecoin-project/go-state-types/builtin/v20/util/adt"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"
	cbg "github.com/whyrusleeping/cbor-gen"
)

func TestCheckStateInvariantsAddressMapIDs(t *testing.T) {
	for _, tc := range []struct {
		name         string
		id           abi.ActorID
		actorPresent bool
		wantMessages []string
	}{
		{name: "retired root key mapping", id: 80},
		{name: "existing ordinary actor", id: 100, actorPresent: true},
		{name: "deleted ordinary actor", id: 100},
		{
			name:         "reserved verified registry ID",
			id:           6,
			wantMessages: []string{"unexpected singleton ID value 6"},
		},
		{
			name:         "reserved datacap ID",
			id:           7,
			wantMessages: []string{"unexpected singleton ID value 7"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			store := adt.WrapStore(ctx, cbor.NewMemCborStore())
			tree, err := builtin.NewTree(store)
			require.NoError(t, err)
			st, err := init_.ConstructState(store, "invariant-test")
			require.NoError(t, err)
			st.NextID = 101

			key, err := address.NewSecp256k1Address([]byte("mapped account key"))
			require.NoError(t, err)
			addressMap, err := adt.AsMap(store, st.AddressMap, builtin.DefaultHamtBitwidth)
			require.NoError(t, err)
			value := cbg.CborInt(tc.id)
			require.NoError(t, addressMap.Put(abi.AddrKey(key), &value))
			st.AddressMap, err = addressMap.Root()
			require.NoError(t, err)

			actorCodes := make(map[string]cid.Cid)
			if tc.actorPresent {
				accountCode, err := abi.CidBuilder.Sum([]byte("account code"))
				require.NoError(t, err)
				actorCodes[manifest.AccountKey] = accountCode
				head, err := store.Put(ctx, &account.State{Address: key})
				require.NoError(t, err)
				idAddress, err := address.NewIDAddress(uint64(tc.id))
				require.NoError(t, err)
				require.NoError(t, tree.SetActorV5(idAddress, &builtin.ActorV5{
					Code:    accountCode,
					Head:    head,
					Balance: big.Zero(),
				}))
			}

			summary, acc := init_.CheckStateInvariants(st, tree, actorCodes)
			require.Equal(t, tc.wantMessages, acc.Messages())
			require.Equal(t, map[address.Address]abi.ActorID{key: tc.id}, summary.AddrIDs)
		})
	}
}
