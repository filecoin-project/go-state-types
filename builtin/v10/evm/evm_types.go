package evm

import (
	"io"

	"github.com/ipfs/go-cid"
	xerrors "golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/abi"
	cbg "github.com/whyrusleeping/cbor-gen"
)

type ConstructorParams struct {
	Creator  [20]byte
	Initcode []byte
}

type ResurrectParams = ConstructorParams

type GetStorageAtParams struct {
	StorageKey [32]byte
}

type DelegateCallParams struct {
	Code   cid.Cid
	Input  []byte
	Caller [20]byte
	Value  abi.TokenAmount
}

type GetBytecodeReturn struct {
	Cid *cid.Cid
}

func (bc *GetBytecodeReturn) UnmarshalCBOR(r io.Reader) error {
	if bc == nil {
		return xerrors.Errorf("cannot unmarshal into nil pointer")
	}

	br := cbg.GetPeeker(r)

	// reset fields
	bc.Cid = nil

	// Check if it's null
	if byte, err := br.ReadByte(); err != nil {
		return err
	} else if byte == cbg.CborNull[0] {
		return nil
	} else if err := br.UnreadByte(); err != nil {
		return err
	}

	c, err := cbg.ReadCid(br)
	if err != nil {
		return xerrors.Errorf("failed to read cid field t.OldPtr: %w", err)
	}
	bc.Cid = &c
	return nil
}

func (bc *GetBytecodeReturn) MarshalCBOR(w io.Writer) error {
	if bc == nil || bc.Cid == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	return cbg.WriteCid(w, *bc.Cid)
}
