package migration

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/actors"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	paych19 "github.com/filecoin-project/go-state-types/builtin/v19/paych"
	system19 "github.com/filecoin-project/go-state-types/builtin/v19/system"
	adt19 "github.com/filecoin-project/go-state-types/builtin/v19/util/adt"
	paych20 "github.com/filecoin-project/go-state-types/builtin/v20/paych"
	adt20 "github.com/filecoin-project/go-state-types/builtin/v20/util/adt"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/filecoin-project/go-state-types/rt"
	"github.com/filecoin-project/go-state-types/test_util"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"
	cbg "github.com/whyrusleeping/cbor-gen"
)

type paychMigrationLogger struct{ t *testing.T }

func (l paychMigrationLogger) Log(_ rt.LogLevel, format string, args ...interface{}) {
	l.t.Logf(format, args...)
}

func TestMigrationPreservesPaymentChannels(t *testing.T) {
	// FIP-0124 disables creation only. Existing channels must retain their funds,
	// settlement schedule, and vouchers' lane nonces across the bundle upgrade.
	for _, tc := range []struct {
		name            string
		settlingAt      abi.ChainEpoch
		minSettleHeight abi.ChainEpoch
	}{
		{name: "active"},
		{name: "settling", settlingAt: 500, minSettleHeight: 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			store := cbor.NewCborStore(test_util.NewSyncBlockStoreInMemory())
			oldStore := adt19.WrapStore(ctx, store)
			newStore := adt20.WrapStore(ctx, store)
			put := func(value interface{}) cid.Cid {
				t.Helper()
				root, err := store.Put(ctx, value)
				require.NoError(t, err)
				return root
			}
			idAddress := func(id uint64) address.Address {
				t.Helper()
				addr, err := address.NewIDAddress(id)
				require.NoError(t, err)
				return addr
			}
			manifestData := func(version actors.Version, offset int) (cid.Cid, map[string]cid.Cid) {
				t.Helper()
				codes := make(map[string]cid.Cid)
				var data manifest.ManifestData
				for i, name := range manifest.GetBuiltinActorsKeys(version) {
					tag := cbg.CborInt(offset + i)
					codes[name] = put(&tag)
					data.Entries = append(data.Entries, manifest.ManifestEntry{Name: name, Code: codes[name]})
				}
				return put(&data), codes
			}
			oldData, oldCodes := manifestData(actors.Version19, 1000)
			newData, newCodes := manifestData(actors.Version20, 2000)
			newManifest := put(&manifest.Manifest{Version: 1, Data: newData})

			lanes, err := adt19.MakeEmptyArray(oldStore, paych19.LaneStatesAmtBitwidth)
			require.NoError(t, err)
			lane := paych19.LaneState{Redeemed: abi.NewTokenAmount(30), Nonce: 7}
			require.NoError(t, lanes.Set(42, &lane))
			laneRoot, err := lanes.Root()
			require.NoError(t, err)
			state := paych19.State{
				From: idAddress(100), To: idAddress(101), ToSend: lane.Redeemed,
				SettlingAt: tc.settlingAt, MinSettleHeight: tc.minSettleHeight, LaneStates: laneRoot,
			}
			channelAddr := idAddress(102)
			channel := builtin.ActorV5{
				Code: oldCodes[manifest.PaychKey], Head: put(&state), Balance: abi.NewTokenAmount(100),
			}
			tree, err := builtin.NewTree(oldStore)
			require.NoError(t, err)
			require.NoError(t, tree.SetActorV5(builtin.SystemActorAddr, &builtin.ActorV5{
				Code: oldCodes[manifest.SystemKey], Head: put(&system19.State{BuiltinActors: oldData}), Balance: big.Zero(),
			}))
			require.NoError(t, tree.SetActorV5(channelAddr, &channel))
			root, err := tree.Flush()
			require.NoError(t, err)

			migrated, err := MigrateStateTree(ctx, store, newManifest, root, 99,
				migration.Config{MaxWorkers: 2}, paychMigrationLogger{t}, migration.NewMemMigrationCache())
			require.NoError(t, err)
			out, err := builtin.LoadTree(newStore, migrated)
			require.NoError(t, err)
			actual, found, err := out.GetActorV5(channelAddr)
			require.NoError(t, err)
			require.True(t, found, "existing payment channel was removed")
			expected := channel
			expected.Code = newCodes[manifest.PaychKey]
			require.Equal(t, expected, *actual)

			var decoded paych20.State
			require.NoError(t, store.Get(ctx, actual.Head, &decoded))
			require.Equal(t, paych20.State(state), decoded)
			newLanes, err := adt20.AsArray(newStore, decoded.LaneStates, paych20.LaneStatesAmtBitwidth)
			require.NoError(t, err)
			var decodedLane paych20.LaneState
			found, err = newLanes.Get(42, &decodedLane)
			require.NoError(t, err)
			require.True(t, found)
			require.Equal(t, paych20.LaneState(lane), decodedLane)
			_, messages := paych20.CheckStateInvariants(&decoded, newStore, actual.Balance)
			require.Empty(t, messages.Messages())
		})
	}
}
