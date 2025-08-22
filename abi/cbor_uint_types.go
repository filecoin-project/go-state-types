package abi

import (
	"fmt"
	"io"

	cbg "github.com/whyrusleeping/cbor-gen"
)

// marshalNullableUint64 is a helper for marshaling nullable uint64 types
func marshalNullableUint64(w io.Writer, v *uint64) error {
	if v == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	if err := cbg.WriteMajorTypeHeader(w, cbg.MajUnsignedInt, *v); err != nil {
		return err
	}
	return nil
}

// unmarshalNullableUint64 is a helper for unmarshaling nullable uint64 types
func unmarshalNullableUint64(r io.Reader, v *uint64, typeName string) error {
	cr := cbg.NewCborReader(r)
	b, err := cr.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to read byte for %s: %w", typeName, err)
	}
	if b == cbg.CborNull[0] {
		// Null value - don't modify v as it might be nil
		// Just return as we've successfully read a null
		return nil
	}

	cr.UnreadByte()
	maj, extra, err := cbg.CborReadHeader(cr)
	if err != nil {
		return fmt.Errorf("failed to parse CBOR header for %s: %w", typeName, err)
	}

	if maj != cbg.MajUnsignedInt {
		return fmt.Errorf("wrong type for %s field: %d", typeName, maj)
	}

	*v = extra
	return nil
}

func (t *SectorSize) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	return marshalNullableUint64(w, (*uint64)(t))
}

func (t *SectorSize) UnmarshalCBOR(r io.Reader) error {
	if t == nil {
		// Just consume the value without storing it
		var dummy uint64
		return unmarshalNullableUint64(r, &dummy, "SectorSize")
	}

	value := uint64(*t)
	if err := unmarshalNullableUint64(r, &value, "SectorSize"); err != nil {
		return err
	}
	*t = SectorSize(value)
	return nil
}

func (t *SectorNumber) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	return marshalNullableUint64(w, (*uint64)(t))
}

func (t *SectorNumber) UnmarshalCBOR(r io.Reader) error {
	if t == nil {
		// Just consume the value without storing it
		var dummy uint64
		return unmarshalNullableUint64(r, &dummy, "SectorNumber")
	}

	value := uint64(*t)
	if err := unmarshalNullableUint64(r, &value, "SectorNumber"); err != nil {
		return err
	}
	*t = SectorNumber(value)
	return nil
}
