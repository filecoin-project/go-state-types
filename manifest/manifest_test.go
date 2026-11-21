package manifest_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/go-state-types/actors"
	"github.com/filecoin-project/go-state-types/manifest"
)

func TestGetBuiltinActorsKeys(t *testing.T) {
	for _, tc := range []struct {
		name    string
		version actors.Version
		want    []string
	}{
		{
			name:    "v8",
			version: actors.Version8,
			want: []string{
				"account", "cron", "init", "storagemarket", "storageminer", "multisig",
				"paymentchannel", "storagepower", "reward", "system", "verifiedregistry",
			},
		},
		{
			name:    "v9",
			version: actors.Version9,
			want: []string{
				"account", "cron", "init", "storagemarket", "storageminer", "multisig",
				"paymentchannel", "storagepower", "reward", "system", "verifiedregistry", "datacap",
			},
		},
		{
			name:    "v10",
			version: actors.Version10,
			want: []string{
				"account", "cron", "init", "storagemarket", "storageminer", "multisig",
				"paymentchannel", "storagepower", "reward", "system", "verifiedregistry", "datacap",
				"evm", "eam", "placeholder", "ethaccount",
			},
		},
		{
			name:    "v19",
			version: actors.Version19,
			want: []string{
				"account", "cron", "init", "storagemarket", "storageminer", "multisig",
				"paymentchannel", "storagepower", "reward", "system", "verifiedregistry", "datacap",
				"evm", "eam", "placeholder", "ethaccount",
			},
		},
		{
			name:    "v20",
			version: actors.Version20,
			want: []string{
				"account", "cron", "init", "storagemarket", "storageminer", "multisig",
				"paymentchannel", "storagepower", "reward", "system",
				"evm", "eam", "placeholder", "ethaccount",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, manifest.GetBuiltinActorsKeys(tc.version))
		})
	}
}
