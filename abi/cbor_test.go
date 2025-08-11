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

func TestCborSectorSize(t *testing.T) {
	tSectorSize := abi.SectorSize(20 * 1024)

	buf := new(bytes.Buffer)
	require.NoError(t, tSectorSize.MarshalCBOR(buf))

	sectorSizeSer := buf.Bytes()
	t.Logf("Serialized SectorSize: %x", sectorSizeSer)

	var sectorSizeDeSer abi.SectorSize
	require.NoError(t, sectorSizeDeSer.UnmarshalCBOR(bytes.NewReader(sectorSizeSer)))
	require.Equal(t, tSectorSize, sectorSizeDeSer)

	var nullSectorSize *abi.SectorSize
	buf = new(bytes.Buffer)
	require.NoError(t, nullSectorSize.MarshalCBOR(buf))
	require.Equal(t, []byte{0xf6}, buf.Bytes(), "null SectorSize should serialize to CBOR null")
	require.NoError(t, nullSectorSize.UnmarshalCBOR(bytes.NewReader([]byte{0xf6})))
	require.Nil(t, nullSectorSize, "unmarshaled null SectorSize should be nil")
}

func TestCborSectorNumber(t *testing.T) {
	tSectorNumber := abi.SectorNumber(42)

	buf := new(bytes.Buffer)
	require.NoError(t, tSectorNumber.MarshalCBOR(buf))

	sectorNumberSer := buf.Bytes()
	t.Logf("Serialized SectorNumber: %x", sectorNumberSer)

	var sectorNumberDeSer abi.SectorNumber
	require.NoError(t, sectorNumberDeSer.UnmarshalCBOR(bytes.NewReader(sectorNumberSer)))
	require.Equal(t, tSectorNumber, sectorNumberDeSer)

	var nullSectorNumber *abi.SectorNumber
	buf = new(bytes.Buffer)
	require.NoError(t, nullSectorNumber.MarshalCBOR(buf))
	require.Equal(t, []byte{0xf6}, buf.Bytes(), "null SectorNumber should serialize to CBOR null")
	require.NoError(t, nullSectorNumber.UnmarshalCBOR(bytes.NewReader([]byte{0xf6})))
	require.Nil(t, nullSectorNumber, "unmarshaled null SectorNumber should be nil")
}
