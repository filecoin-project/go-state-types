package abi

import (
	"io"
)

// CborBytesTransparent NOTE This struct does not create a valid cbor-encoded byte slice. It just passes the bytes through as-is.
type CborBytesTransparent []byte

// MarshalCBOR Does NOT marshall to a cbor-encoding. This is just syntactic sugar to let us pass bytes transparently through lotus which requires a cbor-marshallable object.
func (t *CborBytesTransparent) MarshalCBOR(w io.Writer) error {
	_, err := w.Write(*t)
	return err
}

// UnmarshalCBOR CANNOT read a cbor-encoded byte slice. This will just transparently pass the underlying bytes.
func (t *CborBytesTransparent) UnmarshalCBOR(r io.Reader) error {
	_, err := r.Read(*t)
	return err
}
