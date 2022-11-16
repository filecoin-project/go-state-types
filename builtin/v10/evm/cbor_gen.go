// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package evm

import (
	"fmt"
	"io"
	"sort"

	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf
var _ = cid.Undef
var _ = sort.Sort

var lengthBufState = []byte{130}

func (t *State) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufState); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Bytecode (cid.Cid) (struct)

	if err := cbg.WriteCidBuf(scratch, w, t.Bytecode); err != nil {
		return xerrors.Errorf("failed to write cid field t.Bytecode: %w", err)
	}

	// t.ContractState (cid.Cid) (struct)

	if err := cbg.WriteCidBuf(scratch, w, t.ContractState); err != nil {
		return xerrors.Errorf("failed to write cid field t.ContractState: %w", err)
	}

	return nil
}

func (t *State) UnmarshalCBOR(r io.Reader) error {
	*t = State{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Bytecode (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Bytecode: %w", err)
		}

		t.Bytecode = c

	}
	// t.ContractState (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.ContractState: %w", err)
		}

		t.ContractState = c

	}
	return nil
}

var lengthBufConstructorParams = []byte{130}

func (t *ConstructorParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufConstructorParams); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Creator ([]uint8) (slice)
	if len(t.Creator) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.Creator was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.Creator))); err != nil {
		return err
	}

	if _, err := w.Write(t.Creator[:]); err != nil {
		return err
	}

	// t.Initcode ([]uint8) (slice)
	if len(t.Initcode) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.Initcode was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.Initcode))); err != nil {
		return err
	}

	if _, err := w.Write(t.Initcode[:]); err != nil {
		return err
	}
	return nil
}

func (t *ConstructorParams) UnmarshalCBOR(r io.Reader) error {
	*t = ConstructorParams{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Creator ([]uint8) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.ByteArrayMaxLen {
		return fmt.Errorf("t.Creator: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.Creator = make([]uint8, extra)
	}

	if _, err := io.ReadFull(br, t.Creator[:]); err != nil {
		return err
	}
	// t.Initcode ([]uint8) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.ByteArrayMaxLen {
		return fmt.Errorf("t.Initcode: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.Initcode = make([]uint8, extra)
	}

	if _, err := io.ReadFull(br, t.Initcode[:]); err != nil {
		return err
	}
	return nil
}

var lengthBufGetStorageAtParams = []byte{129}

func (t *GetStorageAtParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufGetStorageAtParams); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.StorageKey ([]uint8) (slice)
	if len(t.StorageKey) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.StorageKey was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.StorageKey))); err != nil {
		return err
	}

	if _, err := w.Write(t.StorageKey[:]); err != nil {
		return err
	}
	return nil
}

func (t *GetStorageAtParams) UnmarshalCBOR(r io.Reader) error {
	*t = GetStorageAtParams{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.StorageKey ([]uint8) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.ByteArrayMaxLen {
		return fmt.Errorf("t.StorageKey: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.StorageKey = make([]uint8, extra)
	}

	if _, err := io.ReadFull(br, t.StorageKey[:]); err != nil {
		return err
	}
	return nil
}
