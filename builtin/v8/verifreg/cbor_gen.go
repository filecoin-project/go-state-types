// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package verifreg

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

var lengthBufState = []byte{132}

func (t *State) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufState); err != nil {
		return err
	}

	// t.RootKey (address.Address) (struct)
	if err := t.RootKey.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Verifiers (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.Verifiers); err != nil {
		return xerrors.Errorf("failed to write cid field t.Verifiers: %w", err)
	}

	// t.VerifiedClients (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.VerifiedClients); err != nil {
		return xerrors.Errorf("failed to write cid field t.VerifiedClients: %w", err)
	}

	// t.RemoveDataCapProposalIDs (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.RemoveDataCapProposalIDs); err != nil {
		return xerrors.Errorf("failed to write cid field t.RemoveDataCapProposalIDs: %w", err)
	}

	return nil
}

func (t *State) UnmarshalCBOR(r io.Reader) (err error) {
	*t = State{}

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

	// t.RootKey (address.Address) (struct)

	{

		if err := t.RootKey.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.RootKey: %w", err)
		}

	}
	// t.Verifiers (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Verifiers: %w", err)
		}

		t.Verifiers = c

	}
	// t.VerifiedClients (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.VerifiedClients: %w", err)
		}

		t.VerifiedClients = c

	}
	// t.RemoveDataCapProposalIDs (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.RemoveDataCapProposalIDs: %w", err)
		}

		t.RemoveDataCapProposalIDs = c

	}
	return nil
}

var lengthBufAddVerifierParams = []byte{130}

func (t *AddVerifierParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufAddVerifierParams); err != nil {
		return err
	}

	// t.Address (address.Address) (struct)
	if err := t.Address.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Allowance (big.Int) (struct)
	if err := t.Allowance.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *AddVerifierParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = AddVerifierParams{}

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

	// t.Address (address.Address) (struct)

	{

		if err := t.Address.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Address: %w", err)
		}

	}
	// t.Allowance (big.Int) (struct)

	{

		if err := t.Allowance.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Allowance: %w", err)
		}

	}
	return nil
}

var lengthBufAddVerifiedClientParams = []byte{130}

func (t *AddVerifiedClientParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufAddVerifiedClientParams); err != nil {
		return err
	}

	// t.Address (address.Address) (struct)
	if err := t.Address.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Allowance (big.Int) (struct)
	if err := t.Allowance.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *AddVerifiedClientParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = AddVerifiedClientParams{}

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

	// t.Address (address.Address) (struct)

	{

		if err := t.Address.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Address: %w", err)
		}

	}
	// t.Allowance (big.Int) (struct)

	{

		if err := t.Allowance.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Allowance: %w", err)
		}

	}
	return nil
}

var lengthBufUseBytesParams = []byte{130}

func (t *UseBytesParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufUseBytesParams); err != nil {
		return err
	}

	// t.Address (address.Address) (struct)
	if err := t.Address.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.DealSize (big.Int) (struct)
	if err := t.DealSize.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *UseBytesParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = UseBytesParams{}

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

	// t.Address (address.Address) (struct)

	{

		if err := t.Address.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Address: %w", err)
		}

	}
	// t.DealSize (big.Int) (struct)

	{

		if err := t.DealSize.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.DealSize: %w", err)
		}

	}
	return nil
}

var lengthBufRestoreBytesParams = []byte{130}

func (t *RestoreBytesParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufRestoreBytesParams); err != nil {
		return err
	}

	// t.Address (address.Address) (struct)
	if err := t.Address.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.DealSize (big.Int) (struct)
	if err := t.DealSize.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *RestoreBytesParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = RestoreBytesParams{}

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

	// t.Address (address.Address) (struct)

	{

		if err := t.Address.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Address: %w", err)
		}

	}
	// t.DealSize (big.Int) (struct)

	{

		if err := t.DealSize.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.DealSize: %w", err)
		}

	}
	return nil
}

var lengthBufRemoveDataCapParams = []byte{132}

func (t *RemoveDataCapParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufRemoveDataCapParams); err != nil {
		return err
	}

	// t.VerifiedClientToRemove (address.Address) (struct)
	if err := t.VerifiedClientToRemove.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.DataCapAmountToRemove (big.Int) (struct)
	if err := t.DataCapAmountToRemove.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.VerifierRequest1 (verifreg.RemoveDataCapRequest) (struct)
	if err := t.VerifierRequest1.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.VerifierRequest2 (verifreg.RemoveDataCapRequest) (struct)
	if err := t.VerifierRequest2.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *RemoveDataCapParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = RemoveDataCapParams{}

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

	// t.VerifiedClientToRemove (address.Address) (struct)

	{

		if err := t.VerifiedClientToRemove.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.VerifiedClientToRemove: %w", err)
		}

	}
	// t.DataCapAmountToRemove (big.Int) (struct)

	{

		if err := t.DataCapAmountToRemove.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.DataCapAmountToRemove: %w", err)
		}

	}
	// t.VerifierRequest1 (verifreg.RemoveDataCapRequest) (struct)

	{

		if err := t.VerifierRequest1.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.VerifierRequest1: %w", err)
		}

	}
	// t.VerifierRequest2 (verifreg.RemoveDataCapRequest) (struct)

	{

		if err := t.VerifierRequest2.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.VerifierRequest2: %w", err)
		}

	}
	return nil
}

var lengthBufRemoveDataCapReturn = []byte{130}

func (t *RemoveDataCapReturn) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufRemoveDataCapReturn); err != nil {
		return err
	}

	// t.VerifiedClient (address.Address) (struct)
	if err := t.VerifiedClient.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.DataCapRemoved (big.Int) (struct)
	if err := t.DataCapRemoved.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *RemoveDataCapReturn) UnmarshalCBOR(r io.Reader) (err error) {
	*t = RemoveDataCapReturn{}

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

	// t.VerifiedClient (address.Address) (struct)

	{

		if err := t.VerifiedClient.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.VerifiedClient: %w", err)
		}

	}
	// t.DataCapRemoved (big.Int) (struct)

	{

		if err := t.DataCapRemoved.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.DataCapRemoved: %w", err)
		}

	}
	return nil
}

var lengthBufRemoveDataCapRequest = []byte{130}

func (t *RemoveDataCapRequest) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufRemoveDataCapRequest); err != nil {
		return err
	}

	// t.Verifier (address.Address) (struct)
	if err := t.Verifier.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.VerifierSignature (crypto.Signature) (struct)
	if err := t.VerifierSignature.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *RemoveDataCapRequest) UnmarshalCBOR(r io.Reader) (err error) {
	*t = RemoveDataCapRequest{}

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

	// t.Verifier (address.Address) (struct)

	{

		if err := t.Verifier.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Verifier: %w", err)
		}

	}
	// t.VerifierSignature (crypto.Signature) (struct)

	{

		if err := t.VerifierSignature.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.VerifierSignature: %w", err)
		}

	}
	return nil
}

var lengthBufRemoveDataCapProposal = []byte{131}

func (t *RemoveDataCapProposal) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufRemoveDataCapProposal); err != nil {
		return err
	}

	// t.VerifiedClient (address.Address) (struct)
	if err := t.VerifiedClient.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.DataCapAmount (big.Int) (struct)
	if err := t.DataCapAmount.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.RemovalProposalID (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.RemovalProposalID)); err != nil {
		return err
	}

	return nil
}

func (t *RemoveDataCapProposal) UnmarshalCBOR(r io.Reader) (err error) {
	*t = RemoveDataCapProposal{}

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

	if extra != 3 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.VerifiedClient (address.Address) (struct)

	{

		if err := t.VerifiedClient.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.VerifiedClient: %w", err)
		}

	}
	// t.DataCapAmount (big.Int) (struct)

	{

		if err := t.DataCapAmount.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.DataCapAmount: %w", err)
		}

	}
	// t.RemovalProposalID (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.RemovalProposalID = uint64(extra)

	}
	return nil
}
