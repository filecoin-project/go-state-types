package eam

import (
	"fmt"
	"io"

	"github.com/filecoin-project/go-address"

	cbg "github.com/whyrusleeping/cbor-gen"
)

type CreateParams struct {
	Initcode []byte
	Nonce    uint64
}

type Create2Params struct {
	Initcode []byte
	Salt     [32]byte
}

// this is transparent
type CreateExternalParams struct {
	Initcode []byte
}

type Return struct {
	ActorID       uint64
	RobustAddress *address.Address
	EthAddress    [20]byte
}

type CreateReturn Return
type Create2Return Return
type Create3Return Return

// cbg doesn't support transparent structs  so we do this by hand
func (t *CreateExternalParams) MarshalCBOR(w io.Writer) error {
	if len(t.Initcode) > cbg.ByteArrayMaxLen {
		return fmt.Errorf("Byte array in field t.Initcode was too long")
	}

	scratch := make([]byte, 9)

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.Initcode))); err != nil {
		return err
	}

	if _, err := w.Write(t.Initcode[:]); err != nil {
		return err
	}

	return nil
}
