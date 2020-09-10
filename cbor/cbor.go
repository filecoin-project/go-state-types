package cbor

import "io"

// These interfaces are intended to match those from whyrusleeping/cbor-gen, such that code generated from that
// system is automatically usable here (but not mandatory).
type CBORMarshaler interface {
	MarshalCBOR(w io.Writer) error
}

type CBORUnmarshaler interface {
	UnmarshalCBOR(r io.Reader) error
}

type CBORer interface {
	CBORMarshaler
	CBORUnmarshaler
}
