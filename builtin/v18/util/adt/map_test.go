package adt

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/go-state-types/test_util"
	cbor "github.com/ipfs/go-ipld-cbor"

	"github.com/stretchr/testify/require"
)

func TestIsEmpty(t *testing.T) {
	m, err := MakeEmptyMap(WrapStore(context.Background(), cbor.NewCborStore(test_util.NewBlockStoreInMemory())), 5)
	require.NoError(t, err)

	isEmpty, err := m.IsEmpty()
	require.NoError(t, err)
	require.True(t, isEmpty)

	val := abi.CborString("val")
	require.NoError(t, m.Put(abi.IntKey(5), &val))

	isEmpty, err = m.IsEmpty()
	require.NoError(t, err)
	require.False(t, isEmpty)
}
