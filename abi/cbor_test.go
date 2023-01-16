package abi_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/go-state-types/abi"
)

func TestCborBytesTransparent(t *testing.T) {
	tBytes := abi.CborBytesTransparent([]byte{0xde, 0xad, 0xbe, 0xef})

	buf := new(bytes.Buffer)
	require.NoError(t, tBytes.MarshalCBOR(buf))

	bytesSer := buf.Bytes()
	require.Equal(t, []byte(tBytes), bytesSer)

	var bytesDeSer abi.CborBytesTransparent
	require.NoError(t, bytesDeSer.UnmarshalCBOR(bytes.NewReader(tBytes)))
	require.Equal(t, tBytes, bytesDeSer)
}
