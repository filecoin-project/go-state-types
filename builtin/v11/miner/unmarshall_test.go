package miner

import (
	"bytes"
	"testing"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"

	"github.com/stretchr/testify/require"
)

func TestCborVestingFunds(t *testing.T) {
	var vf = VestingFunds{
		Funds: []VestingFund{
			{
				Epoch:  1,
				Amount: big.NewInt(3),
			},
			{
				Epoch:  2,
				Amount: big.NewInt(4),
			},
		},
	}
	buf := new(bytes.Buffer)
	err := vf.MarshalCBOR(buf)
	require.NoError(t, err)

	// Value taken from builtin-actors test.
	b := []byte{129, 130, 130, 1, 66, 0, 3, 130, 2, 66, 0, 4}

	require.Equal(t, b, buf.Bytes())

	var vf2 VestingFunds
	err = vf2.UnmarshalCBOR(bytes.NewReader(b))
	require.NoError(t, err)
	require.Equal(t, abi.ChainEpoch(1), vf2.Funds[0].Epoch)
	require.Equal(t, abi.ChainEpoch(2), vf2.Funds[1].Epoch)
	require.Equal(t, big.NewInt(3), vf2.Funds[0].Amount)
	require.Equal(t, big.NewInt(4), vf2.Funds[1].Amount)
}
