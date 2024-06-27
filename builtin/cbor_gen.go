// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package builtin

import (
	"fmt"
	"io"
	"math"
	"sort"

	address "github.com/filecoin-project/go-address"
	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf
var _ = cid.Undef
var _ = math.E
var _ = sort.Sort

var lengthBufActorV4 = []byte{132}

func (t *ActorV4) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufActorV4); err != nil {
		return err
	}

	// t.Code (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.Code); err != nil {
		return xerrors.Errorf("failed to write cid field t.Code: %w", err)
	}

	// t.Head (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.Head); err != nil {
		return xerrors.Errorf("failed to write cid field t.Head: %w", err)
	}

	// t.CallSeqNum (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.CallSeqNum)); err != nil {
		return err
	}

	// t.Balance (big.Int) (struct)
	if err := t.Balance.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *ActorV4) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ActorV4{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 4 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Code (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Code: %w", err)
		}

		t.Code = c

	}
	// t.Head (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Head: %w", err)
		}

		t.Head = c

	}
	// t.CallSeqNum (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.CallSeqNum = uint64(extra)

	}
	// t.Balance (big.Int) (struct)

	{

		if err := t.Balance.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Balance: %w", err)
		}

	}
	return nil
}

var lengthBufActorV5 = []byte{133}

func (t *ActorV5) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufActorV5); err != nil {
		return err
	}

	// t.Code (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.Code); err != nil {
		return xerrors.Errorf("failed to write cid field t.Code: %w", err)
	}

	// t.Head (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.Head); err != nil {
		return xerrors.Errorf("failed to write cid field t.Head: %w", err)
	}

	// t.CallSeqNum (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.CallSeqNum)); err != nil {
		return err
	}

	// t.Balance (big.Int) (struct)
	if err := t.Balance.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.DelegatedAddress (address.Address) (struct)
	if err := t.DelegatedAddress.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *ActorV5) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ActorV5{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 5 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Code (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Code: %w", err)
		}

		t.Code = c

	}
	// t.Head (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Head: %w", err)
		}

		t.Head = c

	}
	// t.CallSeqNum (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.CallSeqNum = uint64(extra)

	}
	// t.Balance (big.Int) (struct)

	{

		if err := t.Balance.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Balance: %w", err)
		}

	}
	// t.DelegatedAddress (address.Address) (struct)

	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}
			t.DelegatedAddress = new(address.Address)
			if err := t.DelegatedAddress.UnmarshalCBOR(cr); err != nil {
				return xerrors.Errorf("unmarshaling t.DelegatedAddress pointer: %w", err)
			}
		}

	}
	return nil
}
