// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package abi

import (
	"fmt"
	"io"
	"math"
	"sort"

	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf
var _ = cid.Undef
var _ = math.E
var _ = sort.Sort

var lengthBufPieceInfo = []byte{130}

func (t *PieceInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufPieceInfo); err != nil {
		return err
	}

	// t.Size (abi.PaddedPieceSize) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Size)); err != nil {
		return err
	}

	// t.PieceCID (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.PieceCID); err != nil {
		return xerrors.Errorf("failed to write cid field t.PieceCID: %w", err)
	}

	return nil
}

func (t *PieceInfo) UnmarshalCBOR(r io.Reader) (err error) {
	*t = PieceInfo{}

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

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Size (abi.PaddedPieceSize) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Size = PaddedPieceSize(extra)

	}
	// t.PieceCID (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.PieceCID: %w", err)
		}

		t.PieceCID = c

	}
	return nil
}

var lengthBufSectorID = []byte{130}

func (t *SectorID) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufSectorID); err != nil {
		return err
	}

	// t.Miner (abi.ActorID) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Miner)); err != nil {
		return err
	}

	// t.Number (abi.SectorNumber) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Number)); err != nil {
		return err
	}

	return nil
}

func (t *SectorID) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SectorID{}

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

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Miner (abi.ActorID) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Miner = ActorID(extra)

	}
	// t.Number (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Number = SectorNumber(extra)

	}
	return nil
}

var lengthBufAddrPairKey = []byte{130}

func (t *AddrPairKey) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufAddrPairKey); err != nil {
		return err
	}

	// t.First (address.Address) (struct)
	if err := t.First.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Second (address.Address) (struct)
	if err := t.Second.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *AddrPairKey) UnmarshalCBOR(r io.Reader) (err error) {
	*t = AddrPairKey{}

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

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.First (address.Address) (struct)

	{

		if err := t.First.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.First: %w", err)
		}

	}
	// t.Second (address.Address) (struct)

	{

		if err := t.Second.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Second: %w", err)
		}

	}
	return nil
}

func (t *DealIDList) MarshalCBOR(w io.Writer) error {
	cw := cbg.NewCborWriter(w)

	// (*t) (abi.DealIDList) (slice)
	if len((*t)) > 8192 {
		return xerrors.Errorf("Slice value in field (*t) was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len((*t)))); err != nil {
		return err
	}
	for _, v := range *t {

		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}

	}
	return nil
}

func (t *DealIDList) UnmarshalCBOR(r io.Reader) (err error) {
	*t = DealIDList{}

	cr := cbg.NewCborReader(r)
	var maj byte
	var extra uint64
	_ = maj
	_ = extra
	// (*t) (abi.DealIDList) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("(*t): array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		(*t) = make([]DealID, extra)
	}

	for i := 0; i < int(extra); i++ {
		{
			var maj byte
			var extra uint64
			var err error
			_ = maj
			_ = extra
			_ = err

			{

				maj, extra, err = cr.ReadHeader()
				if err != nil {
					return err
				}
				if maj != cbg.MajUnsignedInt {
					return fmt.Errorf("wrong type for uint64 field")
				}
				(*t)[i] = DealID(extra)

			}

		}
	}
	return nil
}
