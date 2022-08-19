package evm

import (
	"fmt"
	"io"

	cbg "github.com/whyrusleeping/cbor-gen"
)

type ConstructorParams struct {
	Bytecode  []byte
	InputData []byte
}

type InvokeParams struct {
	InputData []byte
}

type InvokeReturn struct {
	OutputData []byte
}

func (t *InvokeReturn) UnmarshalCBOR(r io.Reader) error {
	*t = InvokeReturn{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.ByteArrayMaxLen {
		return fmt.Errorf("t.InputData: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.OutputData = make([]uint8, extra)
	}

	if _, err := io.ReadFull(br, t.OutputData[:]); err != nil {
		return err
	}
	return nil
}
