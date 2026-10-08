package migration

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	reward18 "github.com/filecoin-project/go-state-types/builtin/v18/reward"
	system18 "github.com/filecoin-project/go-state-types/builtin/v18/system"
	reward19 "github.com/filecoin-project/go-state-types/builtin/v19/reward"
	adt19 "github.com/filecoin-project/go-state-types/builtin/v19/util/adt"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/filecoin-project/go-state-types/rt"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"
	cbg "github.com/whyrusleeping/cbor-gen"
)

type testMigrationLogger struct{ t *testing.T }

func (l testMigrationLogger) Log(_ rt.LogLevel, format string, args ...interface{}) {
	l.t.Logf(format, args...)
}

func TestMigrationChecksRewardActorReferences(t *testing.T) {
	// placeActor installs an actor of the given manifest type at an address in the migration input.
	type placeActor func(addr address.Address, actorType string)

	for _, tc := range []struct {
		name      string
		mutate    func(*RewardMigrationConfig, placeActor)
		wantError string
	}{
		{"account references", func(*RewardMigrationConfig, placeActor) {}, ""},
		{
			"system actor references",
			func(config *RewardMigrationConfig, _ placeActor) {
				config.SWAActor = builtin.SystemActorAddr
				config.Streams[1].Distribution.Writer = builtin.SystemActorAddr
			},
			"",
		},
		{
			"payment channel SWA actor",
			func(config *RewardMigrationConfig, place placeActor) {
				place(config.SWAActor, manifest.PaychKey)
			},
			"SWA actor",
		},
		{
			"payment channel distribution writer",
			func(config *RewardMigrationConfig, place placeActor) {
				place(config.Streams[1].Distribution.Writer, manifest.PaychKey)
			},
			"distribution writer",
		},
		{
			"payment channel recipient",
			func(config *RewardMigrationConfig, place placeActor) {
				place(config.Streams[1].Distribution.Shares[0].Recipient, manifest.PaychKey)
			},
			"reward recipient",
		},
		{
			"missing SWA actor",
			func(config *RewardMigrationConfig, _ placeActor) {
				config.SWAActor = migrationIDAddress(t, 103)
			},
			"SWA actor",
		},
		{
			"missing distribution writer",
			func(config *RewardMigrationConfig, _ placeActor) {
				config.Streams[1].Distribution.Writer = migrationIDAddress(t, 103)
			},
			"distribution writer",
		},
		{
			"burn SWA actor",
			func(config *RewardMigrationConfig, _ placeActor) {
				config.SWAActor = builtin.BurntFundsActorAddr
			},
			"SWA actor is the burn actor",
		},
		{
			"burn distribution writer",
			func(config *RewardMigrationConfig, _ placeActor) {
				config.Streams[1].Distribution.Writer = builtin.BurntFundsActorAddr
			},
			"distribution writer is the burn actor",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			store := cbor.NewMemCborStore()
			adtStore := adt19.WrapStore(ctx, store)
			actors, err := builtin.NewTree(adtStore)
			require.NoError(t, err)

			oldCodes := make(map[string]cid.Cid)
			var oldData, newData manifest.ManifestData
			for i, name := range []string{manifest.SystemKey, manifest.RewardKey, manifest.MarketKey, manifest.PaychKey, manifest.AccountKey} {
				oldTag, newTag := cbg.CborInt(i), cbg.CborInt(i+100)
				oldCode, err := store.Put(ctx, &oldTag)
				require.NoError(t, err)
				newCode, err := store.Put(ctx, &newTag)
				require.NoError(t, err)
				oldCodes[name] = oldCode
				oldData.Entries = append(oldData.Entries, manifest.ManifestEntry{Name: name, Code: oldCode})
				newData.Entries = append(newData.Entries, manifest.ManifestEntry{Name: name, Code: newCode})
			}
			oldDataRoot, err := store.Put(ctx, &oldData)
			require.NoError(t, err)
			newDataRoot, err := store.Put(ctx, &newData)
			require.NoError(t, err)
			newManifest, err := store.Put(ctx, &manifest.Manifest{Version: 1, Data: newDataRoot})
			require.NoError(t, err)
			systemHead, err := store.Put(ctx, &system18.State{BuiltinActors: oldDataRoot})
			require.NoError(t, err)
			require.NoError(t, actors.SetActorV5(builtin.SystemActorAddr, &builtin.ActorV5{
				Code: oldCodes[manifest.SystemKey], Head: systemHead, Balance: big.Zero(),
			}))
			rewardHead, err := store.Put(ctx, reward18.ConstructState(big.Zero()))
			require.NoError(t, err)
			require.NoError(t, actors.SetActorV5(builtin.RewardActorAddr, &builtin.ActorV5{
				Code: oldCodes[manifest.RewardKey], Head: rewardHead, Balance: abi.NewTokenAmount(100),
			}))
			config := validRewardMigrationConfig(t)
			for _, addr := range []address.Address{
				config.SWAActor,
				config.Streams[1].Distribution.Writer,
				config.Streams[1].Distribution.Shares[0].Recipient,
			} {
				require.NoError(t, actors.SetActorV5(addr, &builtin.ActorV5{
					Code: oldCodes[manifest.AccountKey], Head: oldDataRoot, Balance: big.Zero(),
				}))
			}
			tc.mutate(&config, func(addr address.Address, actorType string) {
				require.NoError(t, actors.SetActorV5(addr, &builtin.ActorV5{
					Code: oldCodes[actorType], Head: oldDataRoot, Balance: big.Zero(),
				}))
			})
			root, err := actors.Flush()
			require.NoError(t, err)

			migrated, err := MigrateStateTree(ctx, store, newManifest, root, 99, config,
				migration.Config{MaxWorkers: 1}, testMigrationLogger{t}, migration.NewMemMigrationCache())
			if tc.wantError != "" {
				require.ErrorContains(t, err, tc.wantError)
				return
			}
			require.NoError(t, err)
			out, err := builtin.LoadTree(adtStore, migrated)
			require.NoError(t, err)
			rewardActor, found, err := out.GetActorV5(builtin.RewardActorAddr)
			require.NoError(t, err)
			require.True(t, found)
			var state reward19.State
			require.NoError(t, store.Get(ctx, rewardActor.Head, &state))
			streams, err := state.LoadStreams(adtStore)
			require.NoError(t, err)
			require.Equal(t, config.SWAActor, state.SWAActor)
			require.Equal(t, config.Streams[1].Distribution.Writer, streams.Streams[1].Distribution.Writer)
			require.Equal(t, config.Streams[1].Distribution.Shares, streams.Streams[1].Distribution.Shares)
		})
	}
}
