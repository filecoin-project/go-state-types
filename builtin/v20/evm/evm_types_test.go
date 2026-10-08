package evm

import (
	"bytes"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
	cbg "github.com/whyrusleeping/cbor-gen"
)

func TestGetBytecodeReturn(t *testing.T) {
	randomCid, err := cid.Decode("bafy2bzacecu7n7wbtogznrtuuvf73dsz7wasgyneqasksdblxupnyovmtwxxu")
	require.NoError(t, err)

	var cidbuf bytes.Buffer
	require.NoError(t, cbg.WriteCid(&cidbuf, randomCid))

	in := &GetBytecodeReturn{
		Cid: &randomCid,
	}
	var buf bytes.Buffer
	require.NoError(t, in.MarshalCBOR(&buf))
	require.Equal(t, cidbuf.Bytes(), buf.Bytes())

	var out GetBytecodeReturn
	require.NoError(t, out.UnmarshalCBOR(&buf))
	require.Equal(t, &randomCid, out.Cid)

	in.Cid = nil
	buf.Reset()
	require.NoError(t, in.MarshalCBOR(&buf))
	require.Equal(t, cbg.CborNull, buf.Bytes())

	require.NoError(t, out.UnmarshalCBOR(&buf))
	require.Nil(t, out.Cid)
}
