package abi

import (
	"fmt"
	"io"

	cbg "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/xerrors"
)

type CborBytes []byte

func (t *CborBytes) MarshalCBOR(w io.Writer) error {
	if len(*t) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("byte array was too long")
	}

	if err := cbg.WriteMajorTypeHeader(w, cbg.MajByteString, uint64(len(*t))); err != nil {
		return err
	}

	_, err := w.Write((*t)[:])
	return err
}

func (t *CborBytes) UnmarshalCBOR(r io.Reader) error {

	br := cbg.GetPeeker(r)
	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}

	if extra > cbg.ByteArrayMaxLen {
		return fmt.Errorf("byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		ret := make([]byte, extra)
		if _, err := io.ReadFull(br, ret[:]); err != nil {
			return err
		}

		*t = ret
	}

	return nil
}
