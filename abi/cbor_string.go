package abi

import (
	"io"

	"golang.org/x/xerrors"

	cbg "github.com/whyrusleeping/cbor-gen"
)

type CborString string

func (t *CborString) MarshalCBOR(w io.Writer) error {
	scratch := make([]byte, 8)

	if len(*t) > cbg.MaxLength {
		return xerrors.Errorf("Value in t was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(*t))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(*t)); err != nil {
		return err
	}

	return nil
}

func (t *CborString) UnmarshalCBOR(r io.Reader) error {
	*t = ""

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	{
		sval, err := cbg.ReadStringBuf(br, scratch)
		if err != nil {
			return err
		}

		*t = CborString(sval)
	}

	return nil
}
