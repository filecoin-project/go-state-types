package miner

import (
	"fmt"
	"io"
	"math"
	"sort"

	address "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	verifreg "github.com/filecoin-project/go-state-types/builtin/v16/verifreg"
	proof "github.com/filecoin-project/go-state-types/proof"
	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf
var _ = cid.Undef
var _ = math.E
var _ = sort.Sort

var lengthBufState = []byte{143}

func (t *State) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufState); err != nil {
		return err
	}

	// t.Info (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.Info); err != nil {
		return xerrors.Errorf("failed to write cid field t.Info: %w", err)
	}

	// t.PreCommitDeposits (big.Int) (struct)
	if err := t.PreCommitDeposits.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.LockedFunds (big.Int) (struct)
	if err := t.LockedFunds.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.VestingFunds (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.VestingFunds); err != nil {
		return xerrors.Errorf("failed to write cid field t.VestingFunds: %w", err)
	}

	// t.FeeDebt (big.Int) (struct)
	if err := t.FeeDebt.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.InitialPledge (big.Int) (struct)
	if err := t.InitialPledge.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.PreCommittedSectors (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.PreCommittedSectors); err != nil {
		return xerrors.Errorf("failed to write cid field t.PreCommittedSectors: %w", err)
	}

	// t.PreCommittedSectorsCleanUp (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.PreCommittedSectorsCleanUp); err != nil {
		return xerrors.Errorf("failed to write cid field t.PreCommittedSectorsCleanUp: %w", err)
	}

	// t.AllocatedSectors (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.AllocatedSectors); err != nil {
		return xerrors.Errorf("failed to write cid field t.AllocatedSectors: %w", err)
	}

	// t.Sectors (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.Sectors); err != nil {
		return xerrors.Errorf("failed to write cid field t.Sectors: %w", err)
	}

	// t.ProvingPeriodStart (abi.ChainEpoch) (int64)
	if t.ProvingPeriodStart >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.ProvingPeriodStart)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.ProvingPeriodStart-1)); err != nil {
			return err
		}
	}

	// t.CurrentDeadline (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.CurrentDeadline)); err != nil {
		return err
	}

	// t.Deadlines (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.Deadlines); err != nil {
		return xerrors.Errorf("failed to write cid field t.Deadlines: %w", err)
	}

	// t.EarlyTerminations (bitfield.BitField) (struct)
	if err := t.EarlyTerminations.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.DeadlineCronActive (bool) (bool)
	if err := cbg.WriteBool(w, t.DeadlineCronActive); err != nil {
		return err
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

	if extra != 15 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Info (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Info: %w", err)
		}

		t.Info = c

	}
	// t.PreCommitDeposits (big.Int) (struct)

	{

		if err := t.PreCommitDeposits.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.PreCommitDeposits: %w", err)
		}

	}
	// t.LockedFunds (big.Int) (struct)

	{

		if err := t.LockedFunds.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.LockedFunds: %w", err)
		}

	}
	// t.VestingFunds (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.VestingFunds: %w", err)
		}

		t.VestingFunds = c

	}
	// t.FeeDebt (big.Int) (struct)

	{

		if err := t.FeeDebt.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.FeeDebt: %w", err)
		}

	}
	// t.InitialPledge (big.Int) (struct)

	{

		if err := t.InitialPledge.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.InitialPledge: %w", err)
		}

	}
	// t.PreCommittedSectors (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.PreCommittedSectors: %w", err)
		}

		t.PreCommittedSectors = c

	}
	// t.PreCommittedSectorsCleanUp (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.PreCommittedSectorsCleanUp: %w", err)
		}

		t.PreCommittedSectorsCleanUp = c

	}
	// t.AllocatedSectors (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.AllocatedSectors: %w", err)
		}

		t.AllocatedSectors = c

	}
	// t.Sectors (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Sectors: %w", err)
		}

		t.Sectors = c

	}
	// t.ProvingPeriodStart (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.ProvingPeriodStart = abi.ChainEpoch(extraI)
	}
	// t.CurrentDeadline (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.CurrentDeadline = uint64(extra)

	}
	// t.Deadlines (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Deadlines: %w", err)
		}

		t.Deadlines = c

	}
	// t.EarlyTerminations (bitfield.BitField) (struct)

	{

		if err := t.EarlyTerminations.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.EarlyTerminations: %w", err)
		}

	}
	// t.DeadlineCronActive (bool) (bool)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.DeadlineCronActive = false
	case 21:
		t.DeadlineCronActive = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	return nil
}

var lengthBufMinerInfo = []byte{142}

func (t *MinerInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufMinerInfo); err != nil {
		return err
	}

	// t.Owner (address.Address) (struct)
	if err := t.Owner.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Worker (address.Address) (struct)
	if err := t.Worker.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.ControlAddresses ([]address.Address) (slice)
	if len(t.ControlAddresses) > 8192 {
		return xerrors.Errorf("Slice value in field t.ControlAddresses was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.ControlAddresses))); err != nil {
		return err
	}
	for _, v := range t.ControlAddresses {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}

	// t.PendingWorkerKey (miner.WorkerKeyChange) (struct)
	if err := t.PendingWorkerKey.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.PeerId ([]uint8) (slice)
	if len(t.PeerId) > 2097152 {
		return xerrors.Errorf("Byte array in field t.PeerId was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.PeerId))); err != nil {
		return err
	}

	if _, err := cw.Write(t.PeerId); err != nil {
		return err
	}

	// t.Multiaddrs ([][]uint8) (slice)
	if len(t.Multiaddrs) > 8192 {
		return xerrors.Errorf("Slice value in field t.Multiaddrs was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Multiaddrs))); err != nil {
		return err
	}
	for _, v := range t.Multiaddrs {
		if len(v) > 2097152 {
			return xerrors.Errorf("Byte array in field v was too long")
		}

		if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(v))); err != nil {
			return err
		}

		if _, err := cw.Write(v); err != nil {
			return err
		}

	}

	// t.WindowPoStProofType (abi.RegisteredPoStProof) (int64)
	if t.WindowPoStProofType >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.WindowPoStProofType)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.WindowPoStProofType-1)); err != nil {
			return err
		}
	}

	// t.SectorSize (abi.SectorSize) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SectorSize)); err != nil {
		return err
	}

	// t.WindowPoStPartitionSectors (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.WindowPoStPartitionSectors)); err != nil {
		return err
	}

	// t.ConsensusFaultElapsed (abi.ChainEpoch) (int64)
	if t.ConsensusFaultElapsed >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.ConsensusFaultElapsed)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.ConsensusFaultElapsed-1)); err != nil {
			return err
		}
	}

	// t.PendingOwnerAddress (address.Address) (struct)
	if err := t.PendingOwnerAddress.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Beneficiary (address.Address) (struct)
	if err := t.Beneficiary.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.BeneficiaryTerm (miner.BeneficiaryTerm) (struct)
	if err := t.BeneficiaryTerm.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.PendingBeneficiaryTerm (miner.PendingBeneficiaryChange) (struct)
	if err := t.PendingBeneficiaryTerm.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *MinerInfo) UnmarshalCBOR(r io.Reader) (err error) {
	*t = MinerInfo{}

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

	if extra != 14 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Owner (address.Address) (struct)

	{

		if err := t.Owner.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Owner: %w", err)
		}

	}
	// t.Worker (address.Address) (struct)

	{

		if err := t.Worker.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Worker: %w", err)
		}

	}
	// t.ControlAddresses ([]address.Address) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.ControlAddresses: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.ControlAddresses = make([]address.Address, extra)
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

				if err := t.ControlAddresses[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.ControlAddresses[i]: %w", err)
				}

			}

		}
	}
	// t.PendingWorkerKey (miner.WorkerKeyChange) (struct)

	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}
			t.PendingWorkerKey = new(WorkerKeyChange)
			if err := t.PendingWorkerKey.UnmarshalCBOR(cr); err != nil {
				return xerrors.Errorf("unmarshaling t.PendingWorkerKey pointer: %w", err)
			}
		}

	}
	// t.PeerId ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.PeerId: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.PeerId = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.PeerId); err != nil {
		return err
	}

	// t.Multiaddrs ([][]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Multiaddrs: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Multiaddrs = make([][]uint8, extra)
	}

	for i := 0; i < int(extra); i++ {
		{
			var maj byte
			var extra uint64
			var err error
			_ = maj
			_ = extra
			_ = err

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > 2097152 {
				return fmt.Errorf("t.Multiaddrs[i]: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.Multiaddrs[i] = make([]uint8, extra)
			}

			if _, err := io.ReadFull(cr, t.Multiaddrs[i]); err != nil {
				return err
			}

		}
	}
	// t.WindowPoStProofType (abi.RegisteredPoStProof) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.WindowPoStProofType = abi.RegisteredPoStProof(extraI)
	}
	// t.SectorSize (abi.SectorSize) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.SectorSize = abi.SectorSize(extra)

	}
	// t.WindowPoStPartitionSectors (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.WindowPoStPartitionSectors = uint64(extra)

	}
	// t.ConsensusFaultElapsed (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.ConsensusFaultElapsed = abi.ChainEpoch(extraI)
	}
	// t.PendingOwnerAddress (address.Address) (struct)

	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}
			t.PendingOwnerAddress = new(address.Address)
			if err := t.PendingOwnerAddress.UnmarshalCBOR(cr); err != nil {
				return xerrors.Errorf("unmarshaling t.PendingOwnerAddress pointer: %w", err)
			}
		}

	}
	// t.Beneficiary (address.Address) (struct)

	{

		if err := t.Beneficiary.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Beneficiary: %w", err)
		}

	}
	// t.BeneficiaryTerm (miner.BeneficiaryTerm) (struct)

	{

		if err := t.BeneficiaryTerm.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.BeneficiaryTerm: %w", err)
		}

	}
	// t.PendingBeneficiaryTerm (miner.PendingBeneficiaryChange) (struct)

	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}
			t.PendingBeneficiaryTerm = new(PendingBeneficiaryChange)
			if err := t.PendingBeneficiaryTerm.UnmarshalCBOR(cr); err != nil {
				return xerrors.Errorf("unmarshaling t.PendingBeneficiaryTerm pointer: %w", err)
			}
		}

	}
	return nil
}

var lengthBufDeadlines = []byte{129}

func (t *Deadlines) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufDeadlines); err != nil {
		return err
	}

	// t.Due ([48]cid.Cid) (array)
	if len(t.Due) > 8192 {
		return xerrors.Errorf("Slice value in field t.Due was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Due))); err != nil {
		return err
	}
	for _, v := range t.Due {

		if err := cbg.WriteCid(cw, v); err != nil {
			return xerrors.Errorf("failed to write cid field v: %w", err)
		}

	}
	return nil
}

func (t *Deadlines) UnmarshalCBOR(r io.Reader) (err error) {
	*t = Deadlines{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Due ([48]cid.Cid) (array)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Due: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}
	if extra != 48 {
		return fmt.Errorf("expected array to have 48 elements")
	}

	t.Due = [48]cid.Cid{}
	for i := 0; i < int(extra); i++ {
		{
			var maj byte
			var extra uint64
			var err error
			_ = maj
			_ = extra
			_ = err

			{

				c, err := cbg.ReadCid(cr)
				if err != nil {
					return xerrors.Errorf("failed to read cid field t.Due[i]: %w", err)
				}

				t.Due[i] = c

			}
		}
	}

	return nil
}

var lengthBufDeadline = []byte{139}

func (t *Deadline) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufDeadline); err != nil {
		return err
	}

	// t.Partitions (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.Partitions); err != nil {
		return xerrors.Errorf("failed to write cid field t.Partitions: %w", err)
	}

	// t.ExpirationsEpochs (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.ExpirationsEpochs); err != nil {
		return xerrors.Errorf("failed to write cid field t.ExpirationsEpochs: %w", err)
	}

	// t.PartitionsPoSted (bitfield.BitField) (struct)
	if err := t.PartitionsPoSted.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.EarlyTerminations (bitfield.BitField) (struct)
	if err := t.EarlyTerminations.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.LiveSectors (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.LiveSectors)); err != nil {
		return err
	}

	// t.TotalSectors (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.TotalSectors)); err != nil {
		return err
	}

	// t.FaultyPower (miner.PowerPair) (struct)
	if err := t.FaultyPower.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.OptimisticPoStSubmissions (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.OptimisticPoStSubmissions); err != nil {
		return xerrors.Errorf("failed to write cid field t.OptimisticPoStSubmissions: %w", err)
	}

	// t.SectorsSnapshot (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.SectorsSnapshot); err != nil {
		return xerrors.Errorf("failed to write cid field t.SectorsSnapshot: %w", err)
	}

	// t.PartitionsSnapshot (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.PartitionsSnapshot); err != nil {
		return xerrors.Errorf("failed to write cid field t.PartitionsSnapshot: %w", err)
	}

	// t.OptimisticPoStSubmissionsSnapshot (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.OptimisticPoStSubmissionsSnapshot); err != nil {
		return xerrors.Errorf("failed to write cid field t.OptimisticPoStSubmissionsSnapshot: %w", err)
	}

	return nil
}

func (t *Deadline) UnmarshalCBOR(r io.Reader) (err error) {
	*t = Deadline{}

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

	if extra != 11 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Partitions (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Partitions: %w", err)
		}

		t.Partitions = c

	}
	// t.ExpirationsEpochs (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.ExpirationsEpochs: %w", err)
		}

		t.ExpirationsEpochs = c

	}
	// t.PartitionsPoSted (bitfield.BitField) (struct)

	{

		if err := t.PartitionsPoSted.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.PartitionsPoSted: %w", err)
		}

	}
	// t.EarlyTerminations (bitfield.BitField) (struct)

	{

		if err := t.EarlyTerminations.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.EarlyTerminations: %w", err)
		}

	}
	// t.LiveSectors (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.LiveSectors = uint64(extra)

	}
	// t.TotalSectors (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.TotalSectors = uint64(extra)

	}
	// t.FaultyPower (miner.PowerPair) (struct)

	{

		if err := t.FaultyPower.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.FaultyPower: %w", err)
		}

	}
	// t.OptimisticPoStSubmissions (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.OptimisticPoStSubmissions: %w", err)
		}

		t.OptimisticPoStSubmissions = c

	}
	// t.SectorsSnapshot (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.SectorsSnapshot: %w", err)
		}

		t.SectorsSnapshot = c

	}
	// t.PartitionsSnapshot (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.PartitionsSnapshot: %w", err)
		}

		t.PartitionsSnapshot = c

	}
	// t.OptimisticPoStSubmissionsSnapshot (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.OptimisticPoStSubmissionsSnapshot: %w", err)
		}

		t.OptimisticPoStSubmissionsSnapshot = c

	}
	return nil
}

var lengthBufPartition = []byte{139}

func (t *Partition) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufPartition); err != nil {
		return err
	}

	// t.Sectors (bitfield.BitField) (struct)
	if err := t.Sectors.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Unproven (bitfield.BitField) (struct)
	if err := t.Unproven.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Faults (bitfield.BitField) (struct)
	if err := t.Faults.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Recoveries (bitfield.BitField) (struct)
	if err := t.Recoveries.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Terminated (bitfield.BitField) (struct)
	if err := t.Terminated.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.ExpirationsEpochs (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.ExpirationsEpochs); err != nil {
		return xerrors.Errorf("failed to write cid field t.ExpirationsEpochs: %w", err)
	}

	// t.EarlyTerminated (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.EarlyTerminated); err != nil {
		return xerrors.Errorf("failed to write cid field t.EarlyTerminated: %w", err)
	}

	// t.LivePower (miner.PowerPair) (struct)
	if err := t.LivePower.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.UnprovenPower (miner.PowerPair) (struct)
	if err := t.UnprovenPower.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.FaultyPower (miner.PowerPair) (struct)
	if err := t.FaultyPower.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.RecoveringPower (miner.PowerPair) (struct)
	if err := t.RecoveringPower.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *Partition) UnmarshalCBOR(r io.Reader) (err error) {
	*t = Partition{}

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

	if extra != 11 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Sectors (bitfield.BitField) (struct)

	{

		if err := t.Sectors.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Sectors: %w", err)
		}

	}
	// t.Unproven (bitfield.BitField) (struct)

	{

		if err := t.Unproven.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Unproven: %w", err)
		}

	}
	// t.Faults (bitfield.BitField) (struct)

	{

		if err := t.Faults.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Faults: %w", err)
		}

	}
	// t.Recoveries (bitfield.BitField) (struct)

	{

		if err := t.Recoveries.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Recoveries: %w", err)
		}

	}
	// t.Terminated (bitfield.BitField) (struct)

	{

		if err := t.Terminated.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Terminated: %w", err)
		}

	}
	// t.ExpirationsEpochs (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.ExpirationsEpochs: %w", err)
		}

		t.ExpirationsEpochs = c

	}
	// t.EarlyTerminated (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.EarlyTerminated: %w", err)
		}

		t.EarlyTerminated = c

	}
	// t.LivePower (miner.PowerPair) (struct)

	{

		if err := t.LivePower.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.LivePower: %w", err)
		}

	}
	// t.UnprovenPower (miner.PowerPair) (struct)

	{

		if err := t.UnprovenPower.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.UnprovenPower: %w", err)
		}

	}
	// t.FaultyPower (miner.PowerPair) (struct)

	{

		if err := t.FaultyPower.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.FaultyPower: %w", err)
		}

	}
	// t.RecoveringPower (miner.PowerPair) (struct)

	{

		if err := t.RecoveringPower.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.RecoveringPower: %w", err)
		}

	}
	return nil
}

var lengthBufExpirationSet = []byte{133}

func (t *ExpirationSet) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufExpirationSet); err != nil {
		return err
	}

	// t.OnTimeSectors (bitfield.BitField) (struct)
	if err := t.OnTimeSectors.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.EarlySectors (bitfield.BitField) (struct)
	if err := t.EarlySectors.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.OnTimePledge (big.Int) (struct)
	if err := t.OnTimePledge.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.ActivePower (miner.PowerPair) (struct)
	if err := t.ActivePower.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.FaultyPower (miner.PowerPair) (struct)
	if err := t.FaultyPower.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *ExpirationSet) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ExpirationSet{}

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

	// t.OnTimeSectors (bitfield.BitField) (struct)

	{

		if err := t.OnTimeSectors.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.OnTimeSectors: %w", err)
		}

	}
	// t.EarlySectors (bitfield.BitField) (struct)

	{

		if err := t.EarlySectors.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.EarlySectors: %w", err)
		}

	}
	// t.OnTimePledge (big.Int) (struct)

	{

		if err := t.OnTimePledge.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.OnTimePledge: %w", err)
		}

	}
	// t.ActivePower (miner.PowerPair) (struct)

	{

		if err := t.ActivePower.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.ActivePower: %w", err)
		}

	}
	// t.FaultyPower (miner.PowerPair) (struct)

	{

		if err := t.FaultyPower.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.FaultyPower: %w", err)
		}

	}
	return nil
}

var lengthBufPowerPair = []byte{130}

func (t *PowerPair) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufPowerPair); err != nil {
		return err
	}

	// t.Raw (big.Int) (struct)
	if err := t.Raw.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.QA (big.Int) (struct)
	if err := t.QA.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *PowerPair) UnmarshalCBOR(r io.Reader) (err error) {
	*t = PowerPair{}

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

	// t.Raw (big.Int) (struct)

	{

		if err := t.Raw.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Raw: %w", err)
		}

	}
	// t.QA (big.Int) (struct)

	{

		if err := t.QA.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.QA: %w", err)
		}

	}
	return nil
}

var lengthBufSectorPreCommitOnChainInfo = []byte{131}

func (t *SectorPreCommitOnChainInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufSectorPreCommitOnChainInfo); err != nil {
		return err
	}

	// t.Info (miner.SectorPreCommitInfo) (struct)
	if err := t.Info.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.PreCommitDeposit (big.Int) (struct)
	if err := t.PreCommitDeposit.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.PreCommitEpoch (abi.ChainEpoch) (int64)
	if t.PreCommitEpoch >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.PreCommitEpoch)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.PreCommitEpoch-1)); err != nil {
			return err
		}
	}

	return nil
}

func (t *SectorPreCommitOnChainInfo) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SectorPreCommitOnChainInfo{}

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

	// t.Info (miner.SectorPreCommitInfo) (struct)

	{

		if err := t.Info.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Info: %w", err)
		}

	}
	// t.PreCommitDeposit (big.Int) (struct)

	{

		if err := t.PreCommitDeposit.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.PreCommitDeposit: %w", err)
		}

	}
	// t.PreCommitEpoch (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.PreCommitEpoch = abi.ChainEpoch(extraI)
	}
	return nil
}

var lengthBufSectorPreCommitInfo = []byte{135}

func (t *SectorPreCommitInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufSectorPreCommitInfo); err != nil {
		return err
	}

	// t.SealProof (abi.RegisteredSealProof) (int64)
	if t.SealProof >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SealProof)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.SealProof-1)); err != nil {
			return err
		}
	}

	// t.SectorNumber (abi.SectorNumber) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SectorNumber)); err != nil {
		return err
	}

	// t.SealedCID (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.SealedCID); err != nil {
		return xerrors.Errorf("failed to write cid field t.SealedCID: %w", err)
	}

	// t.SealRandEpoch (abi.ChainEpoch) (int64)
	if t.SealRandEpoch >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SealRandEpoch)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.SealRandEpoch-1)); err != nil {
			return err
		}
	}

	// t.DealIDs ([]abi.DealID) (slice)
	if len(t.DealIDs) > 8192 {
		return xerrors.Errorf("Slice value in field t.DealIDs was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.DealIDs))); err != nil {
		return err
	}
	for _, v := range t.DealIDs {

		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}

	}

	// t.Expiration (abi.ChainEpoch) (int64)
	if t.Expiration >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Expiration)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Expiration-1)); err != nil {
			return err
		}
	}

	// t.UnsealedCid (cid.Cid) (struct)

	if t.UnsealedCid == nil {
		if _, err := cw.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCid(cw, *t.UnsealedCid); err != nil {
			return xerrors.Errorf("failed to write cid field t.UnsealedCid: %w", err)
		}
	}

	return nil
}

func (t *SectorPreCommitInfo) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SectorPreCommitInfo{}

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

	if extra != 7 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.SealProof (abi.RegisteredSealProof) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.SealProof = abi.RegisteredSealProof(extraI)
	}
	// t.SectorNumber (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.SectorNumber = abi.SectorNumber(extra)

	}
	// t.SealedCID (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.SealedCID: %w", err)
		}

		t.SealedCID = c

	}
	// t.SealRandEpoch (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.SealRandEpoch = abi.ChainEpoch(extraI)
	}
	// t.DealIDs ([]abi.DealID) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.DealIDs: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.DealIDs = make([]abi.DealID, extra)
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
				t.DealIDs[i] = abi.DealID(extra)

			}

		}
	}
	// t.Expiration (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.Expiration = abi.ChainEpoch(extraI)
	}
	// t.UnsealedCid (cid.Cid) (struct)

	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}

			c, err := cbg.ReadCid(cr)
			if err != nil {
				return xerrors.Errorf("failed to read cid field t.UnsealedCid: %w", err)
			}

			t.UnsealedCid = &c
		}

	}
	return nil
}

// var lengthBufSectorOnChainInfoOld = []byte{143}
var lengthBufSectorOnChainInfoNew = []byte{144}

func (t *SectorOnChainInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufSectorOnChainInfoNew); err != nil {
		return err
	}

	// t.SectorNumber (abi.SectorNumber) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SectorNumber)); err != nil {
		return err
	}

	// t.SealProof (abi.RegisteredSealProof) (int64)
	if t.SealProof >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SealProof)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.SealProof-1)); err != nil {
			return err
		}
	}

	// t.SealedCID (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.SealedCID); err != nil {
		return xerrors.Errorf("failed to write cid field t.SealedCID: %w", err)
	}

	// t.DealIDs ([]abi.DealID) (slice)
	if len(t.DealIDs) > 8192 {
		return xerrors.Errorf("Slice value in field t.DealIDs was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.DealIDs))); err != nil {
		return err
	}
	for _, v := range t.DealIDs {

		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}

	}

	// t.Activation (abi.ChainEpoch) (int64)
	if t.Activation >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Activation)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Activation-1)); err != nil {
			return err
		}
	}

	// t.Expiration (abi.ChainEpoch) (int64)
	if t.Expiration >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Expiration)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Expiration-1)); err != nil {
			return err
		}
	}

	// t.DealWeight (big.Int) (struct)
	if err := t.DealWeight.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.VerifiedDealWeight (big.Int) (struct)
	if err := t.VerifiedDealWeight.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.InitialPledge (big.Int) (struct)
	if err := t.InitialPledge.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.ExpectedDayReward (big.Int) (struct)
	if err := t.ExpectedDayReward.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.ExpectedStoragePledge (big.Int) (struct)
	if err := t.ExpectedStoragePledge.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.PowerBaseEpoch (abi.ChainEpoch) (int64)
	if t.PowerBaseEpoch >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.PowerBaseEpoch)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.PowerBaseEpoch-1)); err != nil {
			return err
		}
	}

	// t.ReplacedDayReward (big.Int) (struct)
	if err := t.ReplacedDayReward.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.SectorKeyCID (cid.Cid) (struct)

	if t.SectorKeyCID == nil {
		if _, err := cw.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCid(cw, *t.SectorKeyCID); err != nil {
			return xerrors.Errorf("failed to write cid field t.SectorKeyCID: %w", err)
		}
	}

	// t.Flags (miner.SectorOnChainInfoFlags) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Flags)); err != nil {
		return err
	}

	// t.ProvingPeriodFee (big.Int) (struct)
	// Not present for old SectorOnChainInfo objects, always present for new or updated ones

	if err := t.ProvingPeriodFee.MarshalCBOR(cw); err != nil {
		return err
	}

	return nil
}

func (t *SectorOnChainInfo) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SectorOnChainInfo{}

	cr := cbg.NewCborReader(r)

	maj, fieldCount, err := cr.ReadHeader()
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

	if fieldCount < 15 || fieldCount > 16 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.SectorNumber (abi.SectorNumber) (uint64)

	{

		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.SectorNumber = abi.SectorNumber(extra)

	}
	// t.SealProof (abi.RegisteredSealProof) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.SealProof = abi.RegisteredSealProof(extraI)
	}
	// t.SealedCID (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.SealedCID: %w", err)
		}

		t.SealedCID = c

	}
	// t.DealIDs ([]abi.DealID) (slice)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.DealIDs: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.DealIDs = make([]abi.DealID, extra)
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
				t.DealIDs[i] = abi.DealID(extra)

			}

		}
	}
	// t.Activation (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.Activation = abi.ChainEpoch(extraI)
	}
	// t.Expiration (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.Expiration = abi.ChainEpoch(extraI)
	}
	// t.DealWeight (big.Int) (struct)

	{

		if err := t.DealWeight.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.DealWeight: %w", err)
		}

	}
	// t.VerifiedDealWeight (big.Int) (struct)

	{

		if err := t.VerifiedDealWeight.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.VerifiedDealWeight: %w", err)
		}

	}
	// t.InitialPledge (big.Int) (struct)

	{

		if err := t.InitialPledge.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.InitialPledge: %w", err)
		}

	}
	// t.ExpectedDayReward (big.Int) (struct)

	{

		if err := t.ExpectedDayReward.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.ExpectedDayReward: %w", err)
		}

	}
	// t.ExpectedStoragePledge (big.Int) (struct)

	{

		if err := t.ExpectedStoragePledge.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.ExpectedStoragePledge: %w", err)
		}

	}
	// t.PowerBaseEpoch (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.PowerBaseEpoch = abi.ChainEpoch(extraI)
	}
	// t.ReplacedDayReward (big.Int) (struct)

	{

		if err := t.ReplacedDayReward.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.ReplacedDayReward: %w", err)
		}

	}
	// t.SectorKeyCID (cid.Cid) (struct)

	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}

			c, err := cbg.ReadCid(cr)
			if err != nil {
				return xerrors.Errorf("failed to read cid field t.SectorKeyCID: %w", err)
			}

			t.SectorKeyCID = &c
		}

	}
	// t.Flags (miner.SectorOnChainInfoFlags) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Flags = SectorOnChainInfoFlags(extra)

	}

	// t.ProvingPeriodFee (big.Int) (struct)
	// Not present for old SectorOnChainInfo objects, always present for new or updated ones

	if fieldCount == 16 { // new format
		if err := t.ProvingPeriodFee.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.ProvingPeriodFee: %w", err)
		}
	} else {
		t.ProvingPeriodFee = big.Zero()
	}

	return nil
}

var lengthBufWorkerKeyChange = []byte{130}

func (t *WorkerKeyChange) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufWorkerKeyChange); err != nil {
		return err
	}

	// t.NewWorker (address.Address) (struct)
	if err := t.NewWorker.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.EffectiveAt (abi.ChainEpoch) (int64)
	if t.EffectiveAt >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.EffectiveAt)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.EffectiveAt-1)); err != nil {
			return err
		}
	}

	return nil
}

func (t *WorkerKeyChange) UnmarshalCBOR(r io.Reader) (err error) {
	*t = WorkerKeyChange{}

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

	// t.NewWorker (address.Address) (struct)

	{

		if err := t.NewWorker.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.NewWorker: %w", err)
		}

	}
	// t.EffectiveAt (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.EffectiveAt = abi.ChainEpoch(extraI)
	}
	return nil
}

var lengthBufVestingFunds = []byte{129}

func (t *VestingFunds) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufVestingFunds); err != nil {
		return err
	}

	// t.Funds ([]miner.VestingFund) (slice)
	if len(t.Funds) > 8192 {
		return xerrors.Errorf("Slice value in field t.Funds was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Funds))); err != nil {
		return err
	}
	for _, v := range t.Funds {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *VestingFunds) UnmarshalCBOR(r io.Reader) (err error) {
	*t = VestingFunds{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Funds ([]miner.VestingFund) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Funds: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Funds = make([]VestingFund, extra)
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

				if err := t.Funds[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Funds[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufVestingFund = []byte{130}

func (t *VestingFund) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufVestingFund); err != nil {
		return err
	}

	// t.Epoch (abi.ChainEpoch) (int64)
	if t.Epoch >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Epoch)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Epoch-1)); err != nil {
			return err
		}
	}

	// t.Amount (big.Int) (struct)
	if err := t.Amount.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *VestingFund) UnmarshalCBOR(r io.Reader) (err error) {
	*t = VestingFund{}

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

	// t.Epoch (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.Epoch = abi.ChainEpoch(extraI)
	}
	// t.Amount (big.Int) (struct)

	{

		if err := t.Amount.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Amount: %w", err)
		}

	}
	return nil
}

var lengthBufWindowedPoSt = []byte{130}

func (t *WindowedPoSt) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufWindowedPoSt); err != nil {
		return err
	}

	// t.Partitions (bitfield.BitField) (struct)
	if err := t.Partitions.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Proofs ([]proof.PoStProof) (slice)
	if len(t.Proofs) > 8192 {
		return xerrors.Errorf("Slice value in field t.Proofs was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Proofs))); err != nil {
		return err
	}
	for _, v := range t.Proofs {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *WindowedPoSt) UnmarshalCBOR(r io.Reader) (err error) {
	*t = WindowedPoSt{}

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

	// t.Partitions (bitfield.BitField) (struct)

	{

		if err := t.Partitions.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Partitions: %w", err)
		}

	}
	// t.Proofs ([]proof.PoStProof) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Proofs: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Proofs = make([]proof.PoStProof, extra)
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

				if err := t.Proofs[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Proofs[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufActiveBeneficiary = []byte{130}

func (t *ActiveBeneficiary) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufActiveBeneficiary); err != nil {
		return err
	}

	// t.Beneficiary (address.Address) (struct)
	if err := t.Beneficiary.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Term (miner.BeneficiaryTerm) (struct)
	if err := t.Term.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *ActiveBeneficiary) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ActiveBeneficiary{}

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

	// t.Beneficiary (address.Address) (struct)

	{

		if err := t.Beneficiary.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Beneficiary: %w", err)
		}

	}
	// t.Term (miner.BeneficiaryTerm) (struct)

	{

		if err := t.Term.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Term: %w", err)
		}

	}
	return nil
}

var lengthBufBeneficiaryTerm = []byte{131}

func (t *BeneficiaryTerm) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufBeneficiaryTerm); err != nil {
		return err
	}

	// t.Quota (big.Int) (struct)
	if err := t.Quota.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.UsedQuota (big.Int) (struct)
	if err := t.UsedQuota.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Expiration (abi.ChainEpoch) (int64)
	if t.Expiration >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Expiration)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Expiration-1)); err != nil {
			return err
		}
	}

	return nil
}

func (t *BeneficiaryTerm) UnmarshalCBOR(r io.Reader) (err error) {
	*t = BeneficiaryTerm{}

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

	// t.Quota (big.Int) (struct)

	{

		if err := t.Quota.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Quota: %w", err)
		}

	}
	// t.UsedQuota (big.Int) (struct)

	{

		if err := t.UsedQuota.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.UsedQuota: %w", err)
		}

	}
	// t.Expiration (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.Expiration = abi.ChainEpoch(extraI)
	}
	return nil
}

var lengthBufPendingBeneficiaryChange = []byte{133}

func (t *PendingBeneficiaryChange) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufPendingBeneficiaryChange); err != nil {
		return err
	}

	// t.NewBeneficiary (address.Address) (struct)
	if err := t.NewBeneficiary.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.NewQuota (big.Int) (struct)
	if err := t.NewQuota.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.NewExpiration (abi.ChainEpoch) (int64)
	if t.NewExpiration >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.NewExpiration)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.NewExpiration-1)); err != nil {
			return err
		}
	}

	// t.ApprovedByBeneficiary (bool) (bool)
	if err := cbg.WriteBool(w, t.ApprovedByBeneficiary); err != nil {
		return err
	}

	// t.ApprovedByNominee (bool) (bool)
	if err := cbg.WriteBool(w, t.ApprovedByNominee); err != nil {
		return err
	}
	return nil
}

func (t *PendingBeneficiaryChange) UnmarshalCBOR(r io.Reader) (err error) {
	*t = PendingBeneficiaryChange{}

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

	// t.NewBeneficiary (address.Address) (struct)

	{

		if err := t.NewBeneficiary.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.NewBeneficiary: %w", err)
		}

	}
	// t.NewQuota (big.Int) (struct)

	{

		if err := t.NewQuota.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.NewQuota: %w", err)
		}

	}
	// t.NewExpiration (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.NewExpiration = abi.ChainEpoch(extraI)
	}
	// t.ApprovedByBeneficiary (bool) (bool)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.ApprovedByBeneficiary = false
	case 21:
		t.ApprovedByBeneficiary = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	// t.ApprovedByNominee (bool) (bool)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.ApprovedByNominee = false
	case 21:
		t.ApprovedByNominee = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	return nil
}

var lengthBufGetControlAddressesReturn = []byte{131}

func (t *GetControlAddressesReturn) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufGetControlAddressesReturn); err != nil {
		return err
	}

	// t.Owner (address.Address) (struct)
	if err := t.Owner.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Worker (address.Address) (struct)
	if err := t.Worker.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.ControlAddrs ([]address.Address) (slice)
	if len(t.ControlAddrs) > 8192 {
		return xerrors.Errorf("Slice value in field t.ControlAddrs was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.ControlAddrs))); err != nil {
		return err
	}
	for _, v := range t.ControlAddrs {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *GetControlAddressesReturn) UnmarshalCBOR(r io.Reader) (err error) {
	*t = GetControlAddressesReturn{}

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

	// t.Owner (address.Address) (struct)

	{

		if err := t.Owner.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Owner: %w", err)
		}

	}
	// t.Worker (address.Address) (struct)

	{

		if err := t.Worker.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Worker: %w", err)
		}

	}
	// t.ControlAddrs ([]address.Address) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.ControlAddrs: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.ControlAddrs = make([]address.Address, extra)
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

				if err := t.ControlAddrs[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.ControlAddrs[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufChangeWorkerAddressParams = []byte{130}

func (t *ChangeWorkerAddressParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufChangeWorkerAddressParams); err != nil {
		return err
	}

	// t.NewWorker (address.Address) (struct)
	if err := t.NewWorker.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.NewControlAddrs ([]address.Address) (slice)
	if len(t.NewControlAddrs) > 8192 {
		return xerrors.Errorf("Slice value in field t.NewControlAddrs was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.NewControlAddrs))); err != nil {
		return err
	}
	for _, v := range t.NewControlAddrs {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *ChangeWorkerAddressParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ChangeWorkerAddressParams{}

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

	// t.NewWorker (address.Address) (struct)

	{

		if err := t.NewWorker.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.NewWorker: %w", err)
		}

	}
	// t.NewControlAddrs ([]address.Address) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.NewControlAddrs: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.NewControlAddrs = make([]address.Address, extra)
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

				if err := t.NewControlAddrs[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.NewControlAddrs[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufChangePeerIDParams = []byte{129}

func (t *ChangePeerIDParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufChangePeerIDParams); err != nil {
		return err
	}

	// t.NewID ([]uint8) (slice)
	if len(t.NewID) > 2097152 {
		return xerrors.Errorf("Byte array in field t.NewID was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.NewID))); err != nil {
		return err
	}

	if _, err := cw.Write(t.NewID); err != nil {
		return err
	}

	return nil
}

func (t *ChangePeerIDParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ChangePeerIDParams{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.NewID ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.NewID: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.NewID = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.NewID); err != nil {
		return err
	}

	return nil
}

var lengthBufSubmitWindowedPoStParams = []byte{133}

func (t *SubmitWindowedPoStParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufSubmitWindowedPoStParams); err != nil {
		return err
	}

	// t.Deadline (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Deadline)); err != nil {
		return err
	}

	// t.Partitions ([]miner.PoStPartition) (slice)
	if len(t.Partitions) > 8192 {
		return xerrors.Errorf("Slice value in field t.Partitions was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Partitions))); err != nil {
		return err
	}
	for _, v := range t.Partitions {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}

	// t.Proofs ([]proof.PoStProof) (slice)
	if len(t.Proofs) > 8192 {
		return xerrors.Errorf("Slice value in field t.Proofs was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Proofs))); err != nil {
		return err
	}
	for _, v := range t.Proofs {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}

	// t.ChainCommitEpoch (abi.ChainEpoch) (int64)
	if t.ChainCommitEpoch >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.ChainCommitEpoch)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.ChainCommitEpoch-1)); err != nil {
			return err
		}
	}

	// t.ChainCommitRand (abi.Randomness) (slice)
	if len(t.ChainCommitRand) > 2097152 {
		return xerrors.Errorf("Byte array in field t.ChainCommitRand was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.ChainCommitRand))); err != nil {
		return err
	}

	if _, err := cw.Write(t.ChainCommitRand); err != nil {
		return err
	}

	return nil
}

func (t *SubmitWindowedPoStParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SubmitWindowedPoStParams{}

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

	// t.Deadline (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Deadline = uint64(extra)

	}
	// t.Partitions ([]miner.PoStPartition) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Partitions: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Partitions = make([]PoStPartition, extra)
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

				if err := t.Partitions[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Partitions[i]: %w", err)
				}

			}

		}
	}
	// t.Proofs ([]proof.PoStProof) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Proofs: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Proofs = make([]proof.PoStProof, extra)
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

				if err := t.Proofs[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Proofs[i]: %w", err)
				}

			}

		}
	}
	// t.ChainCommitEpoch (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.ChainCommitEpoch = abi.ChainEpoch(extraI)
	}
	// t.ChainCommitRand (abi.Randomness) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.ChainCommitRand: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.ChainCommitRand = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.ChainCommitRand); err != nil {
		return err
	}

	return nil
}

var lengthBufPreCommitSectorParams = []byte{138}

func (t *PreCommitSectorParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufPreCommitSectorParams); err != nil {
		return err
	}

	// t.SealProof (abi.RegisteredSealProof) (int64)
	if t.SealProof >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SealProof)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.SealProof-1)); err != nil {
			return err
		}
	}

	// t.SectorNumber (abi.SectorNumber) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SectorNumber)); err != nil {
		return err
	}

	// t.SealedCID (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.SealedCID); err != nil {
		return xerrors.Errorf("failed to write cid field t.SealedCID: %w", err)
	}

	// t.SealRandEpoch (abi.ChainEpoch) (int64)
	if t.SealRandEpoch >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SealRandEpoch)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.SealRandEpoch-1)); err != nil {
			return err
		}
	}

	// t.DealIDs ([]abi.DealID) (slice)
	if len(t.DealIDs) > 8192 {
		return xerrors.Errorf("Slice value in field t.DealIDs was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.DealIDs))); err != nil {
		return err
	}
	for _, v := range t.DealIDs {

		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}

	}

	// t.Expiration (abi.ChainEpoch) (int64)
	if t.Expiration >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Expiration)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Expiration-1)); err != nil {
			return err
		}
	}

	// t.ReplaceCapacity (bool) (bool)
	if err := cbg.WriteBool(w, t.ReplaceCapacity); err != nil {
		return err
	}

	// t.ReplaceSectorDeadline (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.ReplaceSectorDeadline)); err != nil {
		return err
	}

	// t.ReplaceSectorPartition (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.ReplaceSectorPartition)); err != nil {
		return err
	}

	// t.ReplaceSectorNumber (abi.SectorNumber) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.ReplaceSectorNumber)); err != nil {
		return err
	}

	return nil
}

func (t *PreCommitSectorParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = PreCommitSectorParams{}

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

	if extra != 10 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.SealProof (abi.RegisteredSealProof) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.SealProof = abi.RegisteredSealProof(extraI)
	}
	// t.SectorNumber (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.SectorNumber = abi.SectorNumber(extra)

	}
	// t.SealedCID (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.SealedCID: %w", err)
		}

		t.SealedCID = c

	}
	// t.SealRandEpoch (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.SealRandEpoch = abi.ChainEpoch(extraI)
	}
	// t.DealIDs ([]abi.DealID) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.DealIDs: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.DealIDs = make([]abi.DealID, extra)
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
				t.DealIDs[i] = abi.DealID(extra)

			}

		}
	}
	// t.Expiration (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.Expiration = abi.ChainEpoch(extraI)
	}
	// t.ReplaceCapacity (bool) (bool)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.ReplaceCapacity = false
	case 21:
		t.ReplaceCapacity = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	// t.ReplaceSectorDeadline (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.ReplaceSectorDeadline = uint64(extra)

	}
	// t.ReplaceSectorPartition (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.ReplaceSectorPartition = uint64(extra)

	}
	// t.ReplaceSectorNumber (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.ReplaceSectorNumber = abi.SectorNumber(extra)

	}
	return nil
}

var lengthBufProveCommitSectorParams = []byte{130}

func (t *ProveCommitSectorParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufProveCommitSectorParams); err != nil {
		return err
	}

	// t.SectorNumber (abi.SectorNumber) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SectorNumber)); err != nil {
		return err
	}

	// t.Proof ([]uint8) (slice)
	if len(t.Proof) > 2097152 {
		return xerrors.Errorf("Byte array in field t.Proof was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.Proof))); err != nil {
		return err
	}

	if _, err := cw.Write(t.Proof); err != nil {
		return err
	}

	return nil
}

func (t *ProveCommitSectorParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ProveCommitSectorParams{}

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

	// t.SectorNumber (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.SectorNumber = abi.SectorNumber(extra)

	}
	// t.Proof ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.Proof: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.Proof = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.Proof); err != nil {
		return err
	}

	return nil
}

var lengthBufExtendSectorExpirationParams = []byte{129}

func (t *ExtendSectorExpirationParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufExtendSectorExpirationParams); err != nil {
		return err
	}

	// t.Extensions ([]miner.ExpirationExtension) (slice)
	if len(t.Extensions) > 8192 {
		return xerrors.Errorf("Slice value in field t.Extensions was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Extensions))); err != nil {
		return err
	}
	for _, v := range t.Extensions {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *ExtendSectorExpirationParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ExtendSectorExpirationParams{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Extensions ([]miner.ExpirationExtension) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Extensions: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Extensions = make([]ExpirationExtension, extra)
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

				if err := t.Extensions[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Extensions[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufExtendSectorExpiration2Params = []byte{129}

func (t *ExtendSectorExpiration2Params) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufExtendSectorExpiration2Params); err != nil {
		return err
	}

	// t.Extensions ([]miner.ExpirationExtension2) (slice)
	if len(t.Extensions) > 8192 {
		return xerrors.Errorf("Slice value in field t.Extensions was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Extensions))); err != nil {
		return err
	}
	for _, v := range t.Extensions {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *ExtendSectorExpiration2Params) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ExtendSectorExpiration2Params{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Extensions ([]miner.ExpirationExtension2) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Extensions: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Extensions = make([]ExpirationExtension2, extra)
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

				if err := t.Extensions[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Extensions[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufTerminateSectorsParams = []byte{129}

func (t *TerminateSectorsParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufTerminateSectorsParams); err != nil {
		return err
	}

	// t.Terminations ([]miner.TerminationDeclaration) (slice)
	if len(t.Terminations) > 8192 {
		return xerrors.Errorf("Slice value in field t.Terminations was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Terminations))); err != nil {
		return err
	}
	for _, v := range t.Terminations {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *TerminateSectorsParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = TerminateSectorsParams{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Terminations ([]miner.TerminationDeclaration) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Terminations: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Terminations = make([]TerminationDeclaration, extra)
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

				if err := t.Terminations[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Terminations[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufTerminateSectorsReturn = []byte{129}

func (t *TerminateSectorsReturn) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufTerminateSectorsReturn); err != nil {
		return err
	}

	// t.Done (bool) (bool)
	if err := cbg.WriteBool(w, t.Done); err != nil {
		return err
	}
	return nil
}

func (t *TerminateSectorsReturn) UnmarshalCBOR(r io.Reader) (err error) {
	*t = TerminateSectorsReturn{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Done (bool) (bool)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.Done = false
	case 21:
		t.Done = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	return nil
}

var lengthBufDeclareFaultsParams = []byte{129}

func (t *DeclareFaultsParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufDeclareFaultsParams); err != nil {
		return err
	}

	// t.Faults ([]miner.FaultDeclaration) (slice)
	if len(t.Faults) > 8192 {
		return xerrors.Errorf("Slice value in field t.Faults was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Faults))); err != nil {
		return err
	}
	for _, v := range t.Faults {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *DeclareFaultsParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = DeclareFaultsParams{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Faults ([]miner.FaultDeclaration) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Faults: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Faults = make([]FaultDeclaration, extra)
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

				if err := t.Faults[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Faults[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufDeclareFaultsRecoveredParams = []byte{129}

func (t *DeclareFaultsRecoveredParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufDeclareFaultsRecoveredParams); err != nil {
		return err
	}

	// t.Recoveries ([]miner.RecoveryDeclaration) (slice)
	if len(t.Recoveries) > 8192 {
		return xerrors.Errorf("Slice value in field t.Recoveries was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Recoveries))); err != nil {
		return err
	}
	for _, v := range t.Recoveries {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *DeclareFaultsRecoveredParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = DeclareFaultsRecoveredParams{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Recoveries ([]miner.RecoveryDeclaration) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Recoveries: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Recoveries = make([]RecoveryDeclaration, extra)
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

				if err := t.Recoveries[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Recoveries[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufDeferredCronEventParams = []byte{131}

func (t *DeferredCronEventParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufDeferredCronEventParams); err != nil {
		return err
	}

	// t.EventPayload ([]uint8) (slice)
	if len(t.EventPayload) > 2097152 {
		return xerrors.Errorf("Byte array in field t.EventPayload was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.EventPayload))); err != nil {
		return err
	}

	if _, err := cw.Write(t.EventPayload); err != nil {
		return err
	}

	// t.RewardSmoothed (smoothing.FilterEstimate) (struct)
	if err := t.RewardSmoothed.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.QualityAdjPowerSmoothed (smoothing.FilterEstimate) (struct)
	if err := t.QualityAdjPowerSmoothed.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *DeferredCronEventParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = DeferredCronEventParams{}

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

	// t.EventPayload ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.EventPayload: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.EventPayload = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.EventPayload); err != nil {
		return err
	}

	// t.RewardSmoothed (smoothing.FilterEstimate) (struct)

	{

		if err := t.RewardSmoothed.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.RewardSmoothed: %w", err)
		}

	}
	// t.QualityAdjPowerSmoothed (smoothing.FilterEstimate) (struct)

	{

		if err := t.QualityAdjPowerSmoothed.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.QualityAdjPowerSmoothed: %w", err)
		}

	}
	return nil
}

var lengthBufCheckSectorProvenParams = []byte{129}

func (t *CheckSectorProvenParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufCheckSectorProvenParams); err != nil {
		return err
	}

	// t.SectorNumber (abi.SectorNumber) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SectorNumber)); err != nil {
		return err
	}

	return nil
}

func (t *CheckSectorProvenParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = CheckSectorProvenParams{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.SectorNumber (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.SectorNumber = abi.SectorNumber(extra)

	}
	return nil
}

var lengthBufApplyRewardParams = []byte{130}

func (t *ApplyRewardParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufApplyRewardParams); err != nil {
		return err
	}

	// t.Reward (big.Int) (struct)
	if err := t.Reward.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Penalty (big.Int) (struct)
	if err := t.Penalty.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *ApplyRewardParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ApplyRewardParams{}

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

	// t.Reward (big.Int) (struct)

	{

		if err := t.Reward.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Reward: %w", err)
		}

	}
	// t.Penalty (big.Int) (struct)

	{

		if err := t.Penalty.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Penalty: %w", err)
		}

	}
	return nil
}

var lengthBufReportConsensusFaultParams = []byte{131}

func (t *ReportConsensusFaultParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufReportConsensusFaultParams); err != nil {
		return err
	}

	// t.BlockHeader1 ([]uint8) (slice)
	if len(t.BlockHeader1) > 2097152 {
		return xerrors.Errorf("Byte array in field t.BlockHeader1 was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.BlockHeader1))); err != nil {
		return err
	}

	if _, err := cw.Write(t.BlockHeader1); err != nil {
		return err
	}

	// t.BlockHeader2 ([]uint8) (slice)
	if len(t.BlockHeader2) > 2097152 {
		return xerrors.Errorf("Byte array in field t.BlockHeader2 was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.BlockHeader2))); err != nil {
		return err
	}

	if _, err := cw.Write(t.BlockHeader2); err != nil {
		return err
	}

	// t.BlockHeaderExtra ([]uint8) (slice)
	if len(t.BlockHeaderExtra) > 2097152 {
		return xerrors.Errorf("Byte array in field t.BlockHeaderExtra was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.BlockHeaderExtra))); err != nil {
		return err
	}

	if _, err := cw.Write(t.BlockHeaderExtra); err != nil {
		return err
	}

	return nil
}

func (t *ReportConsensusFaultParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ReportConsensusFaultParams{}

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

	// t.BlockHeader1 ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.BlockHeader1: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.BlockHeader1 = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.BlockHeader1); err != nil {
		return err
	}

	// t.BlockHeader2 ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.BlockHeader2: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.BlockHeader2 = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.BlockHeader2); err != nil {
		return err
	}

	// t.BlockHeaderExtra ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.BlockHeaderExtra: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.BlockHeaderExtra = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.BlockHeaderExtra); err != nil {
		return err
	}

	return nil
}

var lengthBufWithdrawBalanceParams = []byte{129}

func (t *WithdrawBalanceParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufWithdrawBalanceParams); err != nil {
		return err
	}

	// t.AmountRequested (big.Int) (struct)
	if err := t.AmountRequested.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *WithdrawBalanceParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = WithdrawBalanceParams{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.AmountRequested (big.Int) (struct)

	{

		if err := t.AmountRequested.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.AmountRequested: %w", err)
		}

	}
	return nil
}

var lengthBufInternalSectorSetupForPresealParams = []byte{132}

func (t *InternalSectorSetupForPresealParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufInternalSectorSetupForPresealParams); err != nil {
		return err
	}

	// t.Sectors ([]abi.SectorNumber) (slice)
	if len(t.Sectors) > 8192 {
		return xerrors.Errorf("Slice value in field t.Sectors was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Sectors))); err != nil {
		return err
	}
	for _, v := range t.Sectors {

		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}

	}

	// t.RewardSmoothed (smoothing.FilterEstimate) (struct)
	if err := t.RewardSmoothed.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.RewardBaselinePower (big.Int) (struct)
	if err := t.RewardBaselinePower.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.QualityAdjPowerSmoothed (smoothing.FilterEstimate) (struct)
	if err := t.QualityAdjPowerSmoothed.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *InternalSectorSetupForPresealParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = InternalSectorSetupForPresealParams{}

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

	// t.Sectors ([]abi.SectorNumber) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Sectors: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Sectors = make([]abi.SectorNumber, extra)
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
				t.Sectors[i] = abi.SectorNumber(extra)

			}

		}
	}
	// t.RewardSmoothed (smoothing.FilterEstimate) (struct)

	{

		if err := t.RewardSmoothed.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.RewardSmoothed: %w", err)
		}

	}
	// t.RewardBaselinePower (big.Int) (struct)

	{

		if err := t.RewardBaselinePower.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.RewardBaselinePower: %w", err)
		}

	}
	// t.QualityAdjPowerSmoothed (smoothing.FilterEstimate) (struct)

	{

		if err := t.QualityAdjPowerSmoothed.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.QualityAdjPowerSmoothed: %w", err)
		}

	}
	return nil
}

var lengthBufChangeMultiaddrsParams = []byte{129}

func (t *ChangeMultiaddrsParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufChangeMultiaddrsParams); err != nil {
		return err
	}

	// t.NewMultiaddrs ([][]uint8) (slice)
	if len(t.NewMultiaddrs) > 8192 {
		return xerrors.Errorf("Slice value in field t.NewMultiaddrs was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.NewMultiaddrs))); err != nil {
		return err
	}
	for _, v := range t.NewMultiaddrs {
		if len(v) > 2097152 {
			return xerrors.Errorf("Byte array in field v was too long")
		}

		if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(v))); err != nil {
			return err
		}

		if _, err := cw.Write(v); err != nil {
			return err
		}

	}
	return nil
}

func (t *ChangeMultiaddrsParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ChangeMultiaddrsParams{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.NewMultiaddrs ([][]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.NewMultiaddrs: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.NewMultiaddrs = make([][]uint8, extra)
	}

	for i := 0; i < int(extra); i++ {
		{
			var maj byte
			var extra uint64
			var err error
			_ = maj
			_ = extra
			_ = err

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > 2097152 {
				return fmt.Errorf("t.NewMultiaddrs[i]: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.NewMultiaddrs[i] = make([]uint8, extra)
			}

			if _, err := io.ReadFull(cr, t.NewMultiaddrs[i]); err != nil {
				return err
			}

		}
	}
	return nil
}

var lengthBufCompactPartitionsParams = []byte{130}

func (t *CompactPartitionsParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufCompactPartitionsParams); err != nil {
		return err
	}

	// t.Deadline (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Deadline)); err != nil {
		return err
	}

	// t.Partitions (bitfield.BitField) (struct)
	if err := t.Partitions.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *CompactPartitionsParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = CompactPartitionsParams{}

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

	// t.Deadline (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Deadline = uint64(extra)

	}
	// t.Partitions (bitfield.BitField) (struct)

	{

		if err := t.Partitions.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Partitions: %w", err)
		}

	}
	return nil
}

var lengthBufCompactSectorNumbersParams = []byte{129}

func (t *CompactSectorNumbersParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufCompactSectorNumbersParams); err != nil {
		return err
	}

	// t.MaskSectorNumbers (bitfield.BitField) (struct)
	if err := t.MaskSectorNumbers.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *CompactSectorNumbersParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = CompactSectorNumbersParams{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.MaskSectorNumbers (bitfield.BitField) (struct)

	{

		if err := t.MaskSectorNumbers.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.MaskSectorNumbers: %w", err)
		}

	}
	return nil
}

var lengthBufDisputeWindowedPoStParams = []byte{130}

func (t *DisputeWindowedPoStParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufDisputeWindowedPoStParams); err != nil {
		return err
	}

	// t.Deadline (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Deadline)); err != nil {
		return err
	}

	// t.PoStIndex (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.PoStIndex)); err != nil {
		return err
	}

	return nil
}

func (t *DisputeWindowedPoStParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = DisputeWindowedPoStParams{}

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

	// t.Deadline (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Deadline = uint64(extra)

	}
	// t.PoStIndex (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.PoStIndex = uint64(extra)

	}
	return nil
}

var lengthBufPreCommitSectorBatchParams = []byte{129}

func (t *PreCommitSectorBatchParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufPreCommitSectorBatchParams); err != nil {
		return err
	}

	// t.Sectors ([]miner.PreCommitSectorParams) (slice)
	if len(t.Sectors) > 8192 {
		return xerrors.Errorf("Slice value in field t.Sectors was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Sectors))); err != nil {
		return err
	}
	for _, v := range t.Sectors {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *PreCommitSectorBatchParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = PreCommitSectorBatchParams{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Sectors ([]miner.PreCommitSectorParams) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Sectors: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Sectors = make([]PreCommitSectorParams, extra)
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

				if err := t.Sectors[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Sectors[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufProveCommitAggregateParams = []byte{130}

func (t *ProveCommitAggregateParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufProveCommitAggregateParams); err != nil {
		return err
	}

	// t.SectorNumbers (bitfield.BitField) (struct)
	if err := t.SectorNumbers.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.AggregateProof ([]uint8) (slice)
	if len(t.AggregateProof) > 2097152 {
		return xerrors.Errorf("Byte array in field t.AggregateProof was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.AggregateProof))); err != nil {
		return err
	}

	if _, err := cw.Write(t.AggregateProof); err != nil {
		return err
	}

	return nil
}

func (t *ProveCommitAggregateParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ProveCommitAggregateParams{}

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

	// t.SectorNumbers (bitfield.BitField) (struct)

	{

		if err := t.SectorNumbers.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.SectorNumbers: %w", err)
		}

	}
	// t.AggregateProof ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.AggregateProof: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.AggregateProof = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.AggregateProof); err != nil {
		return err
	}

	return nil
}

var lengthBufProveReplicaUpdatesParams = []byte{129}

func (t *ProveReplicaUpdatesParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufProveReplicaUpdatesParams); err != nil {
		return err
	}

	// t.Updates ([]miner.ReplicaUpdate) (slice)
	if len(t.Updates) > 8192 {
		return xerrors.Errorf("Slice value in field t.Updates was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Updates))); err != nil {
		return err
	}
	for _, v := range t.Updates {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *ProveReplicaUpdatesParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ProveReplicaUpdatesParams{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Updates ([]miner.ReplicaUpdate) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Updates: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Updates = make([]ReplicaUpdate, extra)
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

				if err := t.Updates[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Updates[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufCronEventPayload = []byte{129}

func (t *CronEventPayload) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufCronEventPayload); err != nil {
		return err
	}

	// t.EventType (miner.CronEventType) (int64)
	if t.EventType >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.EventType)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.EventType-1)); err != nil {
			return err
		}
	}

	return nil
}

func (t *CronEventPayload) UnmarshalCBOR(r io.Reader) (err error) {
	*t = CronEventPayload{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.EventType (miner.CronEventType) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.EventType = CronEventType(extraI)
	}
	return nil
}

var lengthBufPreCommitSectorBatchParams2 = []byte{129}

func (t *PreCommitSectorBatchParams2) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufPreCommitSectorBatchParams2); err != nil {
		return err
	}

	// t.Sectors ([]miner.SectorPreCommitInfo) (slice)
	if len(t.Sectors) > 8192 {
		return xerrors.Errorf("Slice value in field t.Sectors was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Sectors))); err != nil {
		return err
	}
	for _, v := range t.Sectors {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *PreCommitSectorBatchParams2) UnmarshalCBOR(r io.Reader) (err error) {
	*t = PreCommitSectorBatchParams2{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Sectors ([]miner.SectorPreCommitInfo) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Sectors: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Sectors = make([]SectorPreCommitInfo, extra)
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

				if err := t.Sectors[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Sectors[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufProveReplicaUpdatesParams2 = []byte{129}

func (t *ProveReplicaUpdatesParams2) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufProveReplicaUpdatesParams2); err != nil {
		return err
	}

	// t.Updates ([]miner.ReplicaUpdate2) (slice)
	if len(t.Updates) > 8192 {
		return xerrors.Errorf("Slice value in field t.Updates was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Updates))); err != nil {
		return err
	}
	for _, v := range t.Updates {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *ProveReplicaUpdatesParams2) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ProveReplicaUpdatesParams2{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Updates ([]miner.ReplicaUpdate2) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Updates: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Updates = make([]ReplicaUpdate2, extra)
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

				if err := t.Updates[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Updates[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufChangeBeneficiaryParams = []byte{131}

func (t *ChangeBeneficiaryParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufChangeBeneficiaryParams); err != nil {
		return err
	}

	// t.NewBeneficiary (address.Address) (struct)
	if err := t.NewBeneficiary.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.NewQuota (big.Int) (struct)
	if err := t.NewQuota.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.NewExpiration (abi.ChainEpoch) (int64)
	if t.NewExpiration >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.NewExpiration)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.NewExpiration-1)); err != nil {
			return err
		}
	}

	return nil
}

func (t *ChangeBeneficiaryParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ChangeBeneficiaryParams{}

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

	// t.NewBeneficiary (address.Address) (struct)

	{

		if err := t.NewBeneficiary.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.NewBeneficiary: %w", err)
		}

	}
	// t.NewQuota (big.Int) (struct)

	{

		if err := t.NewQuota.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.NewQuota: %w", err)
		}

	}
	// t.NewExpiration (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.NewExpiration = abi.ChainEpoch(extraI)
	}
	return nil
}

var lengthBufGetBeneficiaryReturn = []byte{130}

func (t *GetBeneficiaryReturn) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufGetBeneficiaryReturn); err != nil {
		return err
	}

	// t.Active (miner.ActiveBeneficiary) (struct)
	if err := t.Active.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Proposed (miner.PendingBeneficiaryChange) (struct)
	if err := t.Proposed.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *GetBeneficiaryReturn) UnmarshalCBOR(r io.Reader) (err error) {
	*t = GetBeneficiaryReturn{}

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

	// t.Active (miner.ActiveBeneficiary) (struct)

	{

		if err := t.Active.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Active: %w", err)
		}

	}
	// t.Proposed (miner.PendingBeneficiaryChange) (struct)

	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}
			t.Proposed = new(PendingBeneficiaryChange)
			if err := t.Proposed.UnmarshalCBOR(cr); err != nil {
				return xerrors.Errorf("unmarshaling t.Proposed pointer: %w", err)
			}
		}

	}
	return nil
}

var lengthBufGetOwnerReturn = []byte{130}

func (t *GetOwnerReturn) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufGetOwnerReturn); err != nil {
		return err
	}

	// t.Owner (address.Address) (struct)
	if err := t.Owner.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Proposed (address.Address) (struct)
	if err := t.Proposed.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *GetOwnerReturn) UnmarshalCBOR(r io.Reader) (err error) {
	*t = GetOwnerReturn{}

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

	// t.Owner (address.Address) (struct)

	{

		if err := t.Owner.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Owner: %w", err)
		}

	}
	// t.Proposed (address.Address) (struct)

	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}
			t.Proposed = new(address.Address)
			if err := t.Proposed.UnmarshalCBOR(cr); err != nil {
				return xerrors.Errorf("unmarshaling t.Proposed pointer: %w", err)
			}
		}

	}
	return nil
}

var lengthBufGetPeerIDReturn = []byte{129}

func (t *GetPeerIDReturn) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufGetPeerIDReturn); err != nil {
		return err
	}

	// t.PeerId ([]uint8) (slice)
	if len(t.PeerId) > 2097152 {
		return xerrors.Errorf("Byte array in field t.PeerId was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.PeerId))); err != nil {
		return err
	}

	if _, err := cw.Write(t.PeerId); err != nil {
		return err
	}

	return nil
}

func (t *GetPeerIDReturn) UnmarshalCBOR(r io.Reader) (err error) {
	*t = GetPeerIDReturn{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.PeerId ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.PeerId: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.PeerId = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.PeerId); err != nil {
		return err
	}

	return nil
}

var lengthBufGetMultiAddrsReturn = []byte{129}

func (t *GetMultiAddrsReturn) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufGetMultiAddrsReturn); err != nil {
		return err
	}

	// t.MultiAddrs ([]uint8) (slice)
	if len(t.MultiAddrs) > 2097152 {
		return xerrors.Errorf("Byte array in field t.MultiAddrs was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.MultiAddrs))); err != nil {
		return err
	}

	if _, err := cw.Write(t.MultiAddrs); err != nil {
		return err
	}

	return nil
}

func (t *GetMultiAddrsReturn) UnmarshalCBOR(r io.Reader) (err error) {
	*t = GetMultiAddrsReturn{}

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

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.MultiAddrs ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.MultiAddrs: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.MultiAddrs = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.MultiAddrs); err != nil {
		return err
	}

	return nil
}

var lengthBufProveCommitSectors3Params = []byte{134}

func (t *ProveCommitSectors3Params) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufProveCommitSectors3Params); err != nil {
		return err
	}

	// t.SectorActivations ([]miner.SectorActivationManifest) (slice)
	if len(t.SectorActivations) > 8192 {
		return xerrors.Errorf("Slice value in field t.SectorActivations was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.SectorActivations))); err != nil {
		return err
	}
	for _, v := range t.SectorActivations {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}

	// t.SectorProofs ([][]uint8) (slice)
	if len(t.SectorProofs) > 8192 {
		return xerrors.Errorf("Slice value in field t.SectorProofs was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.SectorProofs))); err != nil {
		return err
	}
	for _, v := range t.SectorProofs {
		if len(v) > 2097152 {
			return xerrors.Errorf("Byte array in field v was too long")
		}

		if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(v))); err != nil {
			return err
		}

		if _, err := cw.Write(v); err != nil {
			return err
		}

	}

	// t.AggregateProof ([]uint8) (slice)
	if len(t.AggregateProof) > 2097152 {
		return xerrors.Errorf("Byte array in field t.AggregateProof was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.AggregateProof))); err != nil {
		return err
	}

	if _, err := cw.Write(t.AggregateProof); err != nil {
		return err
	}

	// t.AggregateProofType (abi.RegisteredAggregationProof) (int64)
	if t.AggregateProofType == nil {
		if _, err := cw.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if *t.AggregateProofType >= 0 {
			if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(*t.AggregateProofType)); err != nil {
				return err
			}
		} else {
			if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-*t.AggregateProofType-1)); err != nil {
				return err
			}
		}
	}

	// t.RequireActivationSuccess (bool) (bool)
	if err := cbg.WriteBool(w, t.RequireActivationSuccess); err != nil {
		return err
	}

	// t.RequireNotificationSuccess (bool) (bool)
	if err := cbg.WriteBool(w, t.RequireNotificationSuccess); err != nil {
		return err
	}
	return nil
}

func (t *ProveCommitSectors3Params) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ProveCommitSectors3Params{}

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

	if extra != 6 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.SectorActivations ([]miner.SectorActivationManifest) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.SectorActivations: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.SectorActivations = make([]SectorActivationManifest, extra)
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

				if err := t.SectorActivations[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.SectorActivations[i]: %w", err)
				}

			}

		}
	}
	// t.SectorProofs ([][]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.SectorProofs: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.SectorProofs = make([][]uint8, extra)
	}

	for i := 0; i < int(extra); i++ {
		{
			var maj byte
			var extra uint64
			var err error
			_ = maj
			_ = extra
			_ = err

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > 2097152 {
				return fmt.Errorf("t.SectorProofs[i]: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.SectorProofs[i] = make([]uint8, extra)
			}

			if _, err := io.ReadFull(cr, t.SectorProofs[i]); err != nil {
				return err
			}

		}
	}
	// t.AggregateProof ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.AggregateProof: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.AggregateProof = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.AggregateProof); err != nil {
		return err
	}

	// t.AggregateProofType (abi.RegisteredAggregationProof) (int64)
	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}
			maj, extra, err := cr.ReadHeader()
			if err != nil {
				return err
			}
			var extraI int64
			switch maj {
			case cbg.MajUnsignedInt:
				extraI = int64(extra)
				if extraI < 0 {
					return fmt.Errorf("int64 positive overflow")
				}
			case cbg.MajNegativeInt:
				extraI = int64(extra)
				if extraI < 0 {
					return fmt.Errorf("int64 negative overflow")
				}
				extraI = -1 - extraI
			default:
				return fmt.Errorf("wrong type for int64 field: %d", maj)
			}

			t.AggregateProofType = (*abi.RegisteredAggregationProof)(&extraI)
		}
	}
	// t.RequireActivationSuccess (bool) (bool)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.RequireActivationSuccess = false
	case 21:
		t.RequireActivationSuccess = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	// t.RequireNotificationSuccess (bool) (bool)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.RequireNotificationSuccess = false
	case 21:
		t.RequireNotificationSuccess = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	return nil
}

var lengthBufSectorActivationManifest = []byte{130}

func (t *SectorActivationManifest) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufSectorActivationManifest); err != nil {
		return err
	}

	// t.SectorNumber (abi.SectorNumber) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SectorNumber)); err != nil {
		return err
	}

	// t.Pieces ([]miner.PieceActivationManifest) (slice)
	if len(t.Pieces) > 8192 {
		return xerrors.Errorf("Slice value in field t.Pieces was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Pieces))); err != nil {
		return err
	}
	for _, v := range t.Pieces {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *SectorActivationManifest) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SectorActivationManifest{}

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

	// t.SectorNumber (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.SectorNumber = abi.SectorNumber(extra)

	}
	// t.Pieces ([]miner.PieceActivationManifest) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Pieces: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Pieces = make([]PieceActivationManifest, extra)
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

				if err := t.Pieces[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Pieces[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufPieceActivationManifest = []byte{132}

func (t *PieceActivationManifest) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufPieceActivationManifest); err != nil {
		return err
	}

	// t.CID (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.CID); err != nil {
		return xerrors.Errorf("failed to write cid field t.CID: %w", err)
	}

	// t.Size (abi.PaddedPieceSize) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Size)); err != nil {
		return err
	}

	// t.VerifiedAllocationKey (miner.VerifiedAllocationKey) (struct)
	if err := t.VerifiedAllocationKey.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Notify ([]miner.DataActivationNotification) (slice)
	if len(t.Notify) > 8192 {
		return xerrors.Errorf("Slice value in field t.Notify was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Notify))); err != nil {
		return err
	}
	for _, v := range t.Notify {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *PieceActivationManifest) UnmarshalCBOR(r io.Reader) (err error) {
	*t = PieceActivationManifest{}

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

	// t.CID (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.CID: %w", err)
		}

		t.CID = c

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
		t.Size = abi.PaddedPieceSize(extra)

	}
	// t.VerifiedAllocationKey (miner.VerifiedAllocationKey) (struct)

	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}
			t.VerifiedAllocationKey = new(VerifiedAllocationKey)
			if err := t.VerifiedAllocationKey.UnmarshalCBOR(cr); err != nil {
				return xerrors.Errorf("unmarshaling t.VerifiedAllocationKey pointer: %w", err)
			}
		}

	}
	// t.Notify ([]miner.DataActivationNotification) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Notify: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Notify = make([]DataActivationNotification, extra)
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

				if err := t.Notify[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Notify[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufVerifiedAllocationKey = []byte{130}

func (t *VerifiedAllocationKey) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufVerifiedAllocationKey); err != nil {
		return err
	}

	// t.Client (abi.ActorID) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Client)); err != nil {
		return err
	}

	// t.ID (verifreg.AllocationId) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.ID)); err != nil {
		return err
	}

	return nil
}

func (t *VerifiedAllocationKey) UnmarshalCBOR(r io.Reader) (err error) {
	*t = VerifiedAllocationKey{}

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

	// t.Client (abi.ActorID) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Client = abi.ActorID(extra)

	}
	// t.ID (verifreg.AllocationId) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.ID = verifreg.AllocationId(extra)

	}
	return nil
}

var lengthBufDataActivationNotification = []byte{130}

func (t *DataActivationNotification) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufDataActivationNotification); err != nil {
		return err
	}

	// t.Address (address.Address) (struct)
	if err := t.Address.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Payload ([]uint8) (slice)
	if len(t.Payload) > 2097152 {
		return xerrors.Errorf("Byte array in field t.Payload was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.Payload))); err != nil {
		return err
	}

	if _, err := cw.Write(t.Payload); err != nil {
		return err
	}

	return nil
}

func (t *DataActivationNotification) UnmarshalCBOR(r io.Reader) (err error) {
	*t = DataActivationNotification{}

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
	// t.Payload ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.Payload: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.Payload = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.Payload); err != nil {
		return err
	}

	return nil
}

var lengthBufProveReplicaUpdates3Params = []byte{135}

func (t *ProveReplicaUpdates3Params) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufProveReplicaUpdates3Params); err != nil {
		return err
	}

	// t.SectorUpdates ([]miner.SectorUpdateManifest) (slice)
	if len(t.SectorUpdates) > 8192 {
		return xerrors.Errorf("Slice value in field t.SectorUpdates was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.SectorUpdates))); err != nil {
		return err
	}
	for _, v := range t.SectorUpdates {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}

	// t.SectorProofs ([][]uint8) (slice)
	if len(t.SectorProofs) > 8192 {
		return xerrors.Errorf("Slice value in field t.SectorProofs was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.SectorProofs))); err != nil {
		return err
	}
	for _, v := range t.SectorProofs {
		if len(v) > 2097152 {
			return xerrors.Errorf("Byte array in field v was too long")
		}

		if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(v))); err != nil {
			return err
		}

		if _, err := cw.Write(v); err != nil {
			return err
		}

	}

	// t.AggregateProof ([]uint8) (slice)
	if len(t.AggregateProof) > 2097152 {
		return xerrors.Errorf("Byte array in field t.AggregateProof was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.AggregateProof))); err != nil {
		return err
	}

	if _, err := cw.Write(t.AggregateProof); err != nil {
		return err
	}

	// t.UpdateProofsType (abi.RegisteredUpdateProof) (int64)
	if t.UpdateProofsType >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.UpdateProofsType)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.UpdateProofsType-1)); err != nil {
			return err
		}
	}

	// t.AggregateProofType (abi.RegisteredAggregationProof) (int64)
	if t.AggregateProofType == nil {
		if _, err := cw.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if *t.AggregateProofType >= 0 {
			if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(*t.AggregateProofType)); err != nil {
				return err
			}
		} else {
			if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-*t.AggregateProofType-1)); err != nil {
				return err
			}
		}
	}

	// t.RequireActivationSuccess (bool) (bool)
	if err := cbg.WriteBool(w, t.RequireActivationSuccess); err != nil {
		return err
	}

	// t.RequireNotificationSuccess (bool) (bool)
	if err := cbg.WriteBool(w, t.RequireNotificationSuccess); err != nil {
		return err
	}
	return nil
}

func (t *ProveReplicaUpdates3Params) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ProveReplicaUpdates3Params{}

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

	if extra != 7 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.SectorUpdates ([]miner.SectorUpdateManifest) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.SectorUpdates: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.SectorUpdates = make([]SectorUpdateManifest, extra)
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

				if err := t.SectorUpdates[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.SectorUpdates[i]: %w", err)
				}

			}

		}
	}
	// t.SectorProofs ([][]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.SectorProofs: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.SectorProofs = make([][]uint8, extra)
	}

	for i := 0; i < int(extra); i++ {
		{
			var maj byte
			var extra uint64
			var err error
			_ = maj
			_ = extra
			_ = err

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > 2097152 {
				return fmt.Errorf("t.SectorProofs[i]: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.SectorProofs[i] = make([]uint8, extra)
			}

			if _, err := io.ReadFull(cr, t.SectorProofs[i]); err != nil {
				return err
			}

		}
	}
	// t.AggregateProof ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.AggregateProof: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.AggregateProof = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.AggregateProof); err != nil {
		return err
	}

	// t.UpdateProofsType (abi.RegisteredUpdateProof) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.UpdateProofsType = abi.RegisteredUpdateProof(extraI)
	}
	// t.AggregateProofType (abi.RegisteredAggregationProof) (int64)
	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}
			maj, extra, err := cr.ReadHeader()
			if err != nil {
				return err
			}
			var extraI int64
			switch maj {
			case cbg.MajUnsignedInt:
				extraI = int64(extra)
				if extraI < 0 {
					return fmt.Errorf("int64 positive overflow")
				}
			case cbg.MajNegativeInt:
				extraI = int64(extra)
				if extraI < 0 {
					return fmt.Errorf("int64 negative overflow")
				}
				extraI = -1 - extraI
			default:
				return fmt.Errorf("wrong type for int64 field: %d", maj)
			}

			t.AggregateProofType = (*abi.RegisteredAggregationProof)(&extraI)
		}
	}
	// t.RequireActivationSuccess (bool) (bool)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.RequireActivationSuccess = false
	case 21:
		t.RequireActivationSuccess = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	// t.RequireNotificationSuccess (bool) (bool)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.RequireNotificationSuccess = false
	case 21:
		t.RequireNotificationSuccess = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	return nil
}

var lengthBufSectorUpdateManifest = []byte{133}

func (t *SectorUpdateManifest) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufSectorUpdateManifest); err != nil {
		return err
	}

	// t.Sector (abi.SectorNumber) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Sector)); err != nil {
		return err
	}

	// t.Deadline (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Deadline)); err != nil {
		return err
	}

	// t.Partition (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Partition)); err != nil {
		return err
	}

	// t.NewSealedCID (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.NewSealedCID); err != nil {
		return xerrors.Errorf("failed to write cid field t.NewSealedCID: %w", err)
	}

	// t.Pieces ([]miner.PieceActivationManifest) (slice)
	if len(t.Pieces) > 8192 {
		return xerrors.Errorf("Slice value in field t.Pieces was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Pieces))); err != nil {
		return err
	}
	for _, v := range t.Pieces {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *SectorUpdateManifest) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SectorUpdateManifest{}

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

	// t.Sector (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Sector = abi.SectorNumber(extra)

	}
	// t.Deadline (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Deadline = uint64(extra)

	}
	// t.Partition (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Partition = uint64(extra)

	}
	// t.NewSealedCID (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.NewSealedCID: %w", err)
		}

		t.NewSealedCID = c

	}
	// t.Pieces ([]miner.PieceActivationManifest) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Pieces: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Pieces = make([]PieceActivationManifest, extra)
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

				if err := t.Pieces[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Pieces[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufSectorChanges = []byte{131}

func (t *SectorChanges) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufSectorChanges); err != nil {
		return err
	}

	// t.Sector (abi.SectorNumber) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Sector)); err != nil {
		return err
	}

	// t.MinimumCommitmentEpoch (abi.ChainEpoch) (int64)
	if t.MinimumCommitmentEpoch >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.MinimumCommitmentEpoch)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.MinimumCommitmentEpoch-1)); err != nil {
			return err
		}
	}

	// t.Added ([]miner.PieceChange) (slice)
	if len(t.Added) > 8192 {
		return xerrors.Errorf("Slice value in field t.Added was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Added))); err != nil {
		return err
	}
	for _, v := range t.Added {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}
	return nil
}

func (t *SectorChanges) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SectorChanges{}

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

	// t.Sector (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Sector = abi.SectorNumber(extra)

	}
	// t.MinimumCommitmentEpoch (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.MinimumCommitmentEpoch = abi.ChainEpoch(extraI)
	}
	// t.Added ([]miner.PieceChange) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Added: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Added = make([]PieceChange, extra)
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

				if err := t.Added[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Added[i]: %w", err)
				}

			}

		}
	}
	return nil
}

var lengthBufPieceChange = []byte{131}

func (t *PieceChange) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufPieceChange); err != nil {
		return err
	}

	// t.Data (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.Data); err != nil {
		return xerrors.Errorf("failed to write cid field t.Data: %w", err)
	}

	// t.Size (abi.PaddedPieceSize) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Size)); err != nil {
		return err
	}

	// t.Payload ([]uint8) (slice)
	if len(t.Payload) > 2097152 {
		return xerrors.Errorf("Byte array in field t.Payload was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.Payload))); err != nil {
		return err
	}

	if _, err := cw.Write(t.Payload); err != nil {
		return err
	}

	return nil
}

func (t *PieceChange) UnmarshalCBOR(r io.Reader) (err error) {
	*t = PieceChange{}

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

	// t.Data (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Data: %w", err)
		}

		t.Data = c

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
		t.Size = abi.PaddedPieceSize(extra)

	}
	// t.Payload ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.Payload: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.Payload = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.Payload); err != nil {
		return err
	}

	return nil
}

var lengthBufFaultDeclaration = []byte{131}

func (t *FaultDeclaration) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufFaultDeclaration); err != nil {
		return err
	}

	// t.Deadline (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Deadline)); err != nil {
		return err
	}

	// t.Partition (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Partition)); err != nil {
		return err
	}

	// t.Sectors (bitfield.BitField) (struct)
	if err := t.Sectors.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *FaultDeclaration) UnmarshalCBOR(r io.Reader) (err error) {
	*t = FaultDeclaration{}

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

	// t.Deadline (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Deadline = uint64(extra)

	}
	// t.Partition (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Partition = uint64(extra)

	}
	// t.Sectors (bitfield.BitField) (struct)

	{

		if err := t.Sectors.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Sectors: %w", err)
		}

	}
	return nil
}

var lengthBufRecoveryDeclaration = []byte{131}

func (t *RecoveryDeclaration) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufRecoveryDeclaration); err != nil {
		return err
	}

	// t.Deadline (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Deadline)); err != nil {
		return err
	}

	// t.Partition (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Partition)); err != nil {
		return err
	}

	// t.Sectors (bitfield.BitField) (struct)
	if err := t.Sectors.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *RecoveryDeclaration) UnmarshalCBOR(r io.Reader) (err error) {
	*t = RecoveryDeclaration{}

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

	// t.Deadline (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Deadline = uint64(extra)

	}
	// t.Partition (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Partition = uint64(extra)

	}
	// t.Sectors (bitfield.BitField) (struct)

	{

		if err := t.Sectors.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Sectors: %w", err)
		}

	}
	return nil
}

var lengthBufExpirationExtension = []byte{132}

func (t *ExpirationExtension) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufExpirationExtension); err != nil {
		return err
	}

	// t.Deadline (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Deadline)); err != nil {
		return err
	}

	// t.Partition (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Partition)); err != nil {
		return err
	}

	// t.Sectors (bitfield.BitField) (struct)
	if err := t.Sectors.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.NewExpiration (abi.ChainEpoch) (int64)
	if t.NewExpiration >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.NewExpiration)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.NewExpiration-1)); err != nil {
			return err
		}
	}

	return nil
}

func (t *ExpirationExtension) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ExpirationExtension{}

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

	// t.Deadline (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Deadline = uint64(extra)

	}
	// t.Partition (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Partition = uint64(extra)

	}
	// t.Sectors (bitfield.BitField) (struct)

	{

		if err := t.Sectors.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Sectors: %w", err)
		}

	}
	// t.NewExpiration (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.NewExpiration = abi.ChainEpoch(extraI)
	}
	return nil
}

var lengthBufTerminationDeclaration = []byte{131}

func (t *TerminationDeclaration) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufTerminationDeclaration); err != nil {
		return err
	}

	// t.Deadline (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Deadline)); err != nil {
		return err
	}

	// t.Partition (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Partition)); err != nil {
		return err
	}

	// t.Sectors (bitfield.BitField) (struct)
	if err := t.Sectors.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *TerminationDeclaration) UnmarshalCBOR(r io.Reader) (err error) {
	*t = TerminationDeclaration{}

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

	// t.Deadline (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Deadline = uint64(extra)

	}
	// t.Partition (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Partition = uint64(extra)

	}
	// t.Sectors (bitfield.BitField) (struct)

	{

		if err := t.Sectors.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Sectors: %w", err)
		}

	}
	return nil
}

var lengthBufPoStPartition = []byte{130}

func (t *PoStPartition) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufPoStPartition); err != nil {
		return err
	}

	// t.Index (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Index)); err != nil {
		return err
	}

	// t.Skipped (bitfield.BitField) (struct)
	if err := t.Skipped.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *PoStPartition) UnmarshalCBOR(r io.Reader) (err error) {
	*t = PoStPartition{}

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

	// t.Index (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Index = uint64(extra)

	}
	// t.Skipped (bitfield.BitField) (struct)

	{

		if err := t.Skipped.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Skipped: %w", err)
		}

	}
	return nil
}

var lengthBufReplicaUpdate = []byte{135}

func (t *ReplicaUpdate) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufReplicaUpdate); err != nil {
		return err
	}

	// t.SectorID (abi.SectorNumber) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SectorID)); err != nil {
		return err
	}

	// t.Deadline (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Deadline)); err != nil {
		return err
	}

	// t.Partition (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Partition)); err != nil {
		return err
	}

	// t.NewSealedSectorCID (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.NewSealedSectorCID); err != nil {
		return xerrors.Errorf("failed to write cid field t.NewSealedSectorCID: %w", err)
	}

	// t.Deals ([]abi.DealID) (slice)
	if len(t.Deals) > 8192 {
		return xerrors.Errorf("Slice value in field t.Deals was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Deals))); err != nil {
		return err
	}
	for _, v := range t.Deals {

		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}

	}

	// t.UpdateProofType (abi.RegisteredUpdateProof) (int64)
	if t.UpdateProofType >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.UpdateProofType)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.UpdateProofType-1)); err != nil {
			return err
		}
	}

	// t.ReplicaProof ([]uint8) (slice)
	if len(t.ReplicaProof) > 2097152 {
		return xerrors.Errorf("Byte array in field t.ReplicaProof was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.ReplicaProof))); err != nil {
		return err
	}

	if _, err := cw.Write(t.ReplicaProof); err != nil {
		return err
	}

	return nil
}

func (t *ReplicaUpdate) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ReplicaUpdate{}

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

	if extra != 7 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.SectorID (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.SectorID = abi.SectorNumber(extra)

	}
	// t.Deadline (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Deadline = uint64(extra)

	}
	// t.Partition (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Partition = uint64(extra)

	}
	// t.NewSealedSectorCID (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.NewSealedSectorCID: %w", err)
		}

		t.NewSealedSectorCID = c

	}
	// t.Deals ([]abi.DealID) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Deals: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Deals = make([]abi.DealID, extra)
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
				t.Deals[i] = abi.DealID(extra)

			}

		}
	}
	// t.UpdateProofType (abi.RegisteredUpdateProof) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.UpdateProofType = abi.RegisteredUpdateProof(extraI)
	}
	// t.ReplicaProof ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.ReplicaProof: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.ReplicaProof = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.ReplicaProof); err != nil {
		return err
	}

	return nil
}

var lengthBufReplicaUpdate2 = []byte{136}

func (t *ReplicaUpdate2) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufReplicaUpdate2); err != nil {
		return err
	}

	// t.SectorID (abi.SectorNumber) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SectorID)); err != nil {
		return err
	}

	// t.Deadline (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Deadline)); err != nil {
		return err
	}

	// t.Partition (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Partition)); err != nil {
		return err
	}

	// t.NewSealedSectorCID (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.NewSealedSectorCID); err != nil {
		return xerrors.Errorf("failed to write cid field t.NewSealedSectorCID: %w", err)
	}

	// t.NewUnsealedSectorCID (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.NewUnsealedSectorCID); err != nil {
		return xerrors.Errorf("failed to write cid field t.NewUnsealedSectorCID: %w", err)
	}

	// t.Deals ([]abi.DealID) (slice)
	if len(t.Deals) > 8192 {
		return xerrors.Errorf("Slice value in field t.Deals was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Deals))); err != nil {
		return err
	}
	for _, v := range t.Deals {

		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}

	}

	// t.UpdateProofType (abi.RegisteredUpdateProof) (int64)
	if t.UpdateProofType >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.UpdateProofType)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.UpdateProofType-1)); err != nil {
			return err
		}
	}

	// t.ReplicaProof ([]uint8) (slice)
	if len(t.ReplicaProof) > 2097152 {
		return xerrors.Errorf("Byte array in field t.ReplicaProof was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.ReplicaProof))); err != nil {
		return err
	}

	if _, err := cw.Write(t.ReplicaProof); err != nil {
		return err
	}

	return nil
}

func (t *ReplicaUpdate2) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ReplicaUpdate2{}

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

	if extra != 8 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.SectorID (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.SectorID = abi.SectorNumber(extra)

	}
	// t.Deadline (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Deadline = uint64(extra)

	}
	// t.Partition (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Partition = uint64(extra)

	}
	// t.NewSealedSectorCID (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.NewSealedSectorCID: %w", err)
		}

		t.NewSealedSectorCID = c

	}
	// t.NewUnsealedSectorCID (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.NewUnsealedSectorCID: %w", err)
		}

		t.NewUnsealedSectorCID = c

	}
	// t.Deals ([]abi.DealID) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Deals: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Deals = make([]abi.DealID, extra)
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
				t.Deals[i] = abi.DealID(extra)

			}

		}
	}
	// t.UpdateProofType (abi.RegisteredUpdateProof) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.UpdateProofType = abi.RegisteredUpdateProof(extraI)
	}
	// t.ReplicaProof ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.ReplicaProof: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.ReplicaProof = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.ReplicaProof); err != nil {
		return err
	}

	return nil
}

var lengthBufExpirationExtension2 = []byte{133}

func (t *ExpirationExtension2) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufExpirationExtension2); err != nil {
		return err
	}

	// t.Deadline (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Deadline)); err != nil {
		return err
	}

	// t.Partition (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Partition)); err != nil {
		return err
	}

	// t.Sectors (bitfield.BitField) (struct)
	if err := t.Sectors.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.SectorsWithClaims ([]miner.SectorClaim) (slice)
	if len(t.SectorsWithClaims) > 8192 {
		return xerrors.Errorf("Slice value in field t.SectorsWithClaims was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.SectorsWithClaims))); err != nil {
		return err
	}
	for _, v := range t.SectorsWithClaims {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}

	// t.NewExpiration (abi.ChainEpoch) (int64)
	if t.NewExpiration >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.NewExpiration)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.NewExpiration-1)); err != nil {
			return err
		}
	}

	return nil
}

func (t *ExpirationExtension2) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ExpirationExtension2{}

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

	// t.Deadline (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Deadline = uint64(extra)

	}
	// t.Partition (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Partition = uint64(extra)

	}
	// t.Sectors (bitfield.BitField) (struct)

	{

		if err := t.Sectors.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Sectors: %w", err)
		}

	}
	// t.SectorsWithClaims ([]miner.SectorClaim) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.SectorsWithClaims: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.SectorsWithClaims = make([]SectorClaim, extra)
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

				if err := t.SectorsWithClaims[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.SectorsWithClaims[i]: %w", err)
				}

			}

		}
	}
	// t.NewExpiration (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.NewExpiration = abi.ChainEpoch(extraI)
	}
	return nil
}

var lengthBufSectorClaim = []byte{131}

func (t *SectorClaim) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufSectorClaim); err != nil {
		return err
	}

	// t.SectorNumber (abi.SectorNumber) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SectorNumber)); err != nil {
		return err
	}

	// t.MaintainClaims ([]verifreg.ClaimId) (slice)
	if len(t.MaintainClaims) > 8192 {
		return xerrors.Errorf("Slice value in field t.MaintainClaims was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.MaintainClaims))); err != nil {
		return err
	}
	for _, v := range t.MaintainClaims {

		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}

	}

	// t.DropClaims ([]verifreg.ClaimId) (slice)
	if len(t.DropClaims) > 8192 {
		return xerrors.Errorf("Slice value in field t.DropClaims was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.DropClaims))); err != nil {
		return err
	}
	for _, v := range t.DropClaims {

		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}

	}
	return nil
}

func (t *SectorClaim) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SectorClaim{}

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

	// t.SectorNumber (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.SectorNumber = abi.SectorNumber(extra)

	}
	// t.MaintainClaims ([]verifreg.ClaimId) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.MaintainClaims: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.MaintainClaims = make([]verifreg.ClaimId, extra)
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
				t.MaintainClaims[i] = verifreg.ClaimId(extra)

			}

		}
	}
	// t.DropClaims ([]verifreg.ClaimId) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.DropClaims: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.DropClaims = make([]verifreg.ClaimId, extra)
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
				t.DropClaims[i] = verifreg.ClaimId(extra)

			}

		}
	}
	return nil
}

var lengthBufSectorNIActivationInfo = []byte{134}

func (t *SectorNIActivationInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufSectorNIActivationInfo); err != nil {
		return err
	}

	// t.SealingNumber (abi.SectorNumber) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SealingNumber)); err != nil {
		return err
	}

	// t.SealerID (abi.ActorID) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SealerID)); err != nil {
		return err
	}

	// t.SealedCID (cid.Cid) (struct)

	if err := cbg.WriteCid(cw, t.SealedCID); err != nil {
		return xerrors.Errorf("failed to write cid field t.SealedCID: %w", err)
	}

	// t.SectorNumber (abi.SectorNumber) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SectorNumber)); err != nil {
		return err
	}

	// t.SealRandEpoch (abi.ChainEpoch) (int64)
	if t.SealRandEpoch >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SealRandEpoch)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.SealRandEpoch-1)); err != nil {
			return err
		}
	}

	// t.Expiration (abi.ChainEpoch) (int64)
	if t.Expiration >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Expiration)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Expiration-1)); err != nil {
			return err
		}
	}

	return nil
}

func (t *SectorNIActivationInfo) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SectorNIActivationInfo{}

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

	if extra != 6 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.SealingNumber (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.SealingNumber = abi.SectorNumber(extra)

	}
	// t.SealerID (abi.ActorID) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.SealerID = abi.ActorID(extra)

	}
	// t.SealedCID (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(cr)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.SealedCID: %w", err)
		}

		t.SealedCID = c

	}
	// t.SectorNumber (abi.SectorNumber) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.SectorNumber = abi.SectorNumber(extra)

	}
	// t.SealRandEpoch (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.SealRandEpoch = abi.ChainEpoch(extraI)
	}
	// t.Expiration (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.Expiration = abi.ChainEpoch(extraI)
	}
	return nil
}

var lengthBufProveCommitSectorsNIParams = []byte{134}

func (t *ProveCommitSectorsNIParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufProveCommitSectorsNIParams); err != nil {
		return err
	}

	// t.Sectors ([]miner.SectorNIActivationInfo) (slice)
	if len(t.Sectors) > 8192 {
		return xerrors.Errorf("Slice value in field t.Sectors was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Sectors))); err != nil {
		return err
	}
	for _, v := range t.Sectors {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}

	}

	// t.AggregateProof ([]uint8) (slice)
	if len(t.AggregateProof) > 2097152 {
		return xerrors.Errorf("Byte array in field t.AggregateProof was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.AggregateProof))); err != nil {
		return err
	}

	if _, err := cw.Write(t.AggregateProof); err != nil {
		return err
	}

	// t.SealProofType (abi.RegisteredSealProof) (int64)
	if t.SealProofType >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SealProofType)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.SealProofType-1)); err != nil {
			return err
		}
	}

	// t.AggregateProofType (abi.RegisteredAggregationProof) (int64)
	if t.AggregateProofType >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.AggregateProofType)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.AggregateProofType-1)); err != nil {
			return err
		}
	}

	// t.ProvingDeadline (uint64) (uint64)

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.ProvingDeadline)); err != nil {
		return err
	}

	// t.RequireActivationSuccess (bool) (bool)
	if err := cbg.WriteBool(w, t.RequireActivationSuccess); err != nil {
		return err
	}
	return nil
}

func (t *ProveCommitSectorsNIParams) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ProveCommitSectorsNIParams{}

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

	if extra != 6 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Sectors ([]miner.SectorNIActivationInfo) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 8192 {
		return fmt.Errorf("t.Sectors: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Sectors = make([]SectorNIActivationInfo, extra)
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

				if err := t.Sectors[i].UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling t.Sectors[i]: %w", err)
				}

			}

		}
	}
	// t.AggregateProof ([]uint8) (slice)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > 2097152 {
		return fmt.Errorf("t.AggregateProof: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.AggregateProof = make([]uint8, extra)
	}

	if _, err := io.ReadFull(cr, t.AggregateProof); err != nil {
		return err
	}

	// t.SealProofType (abi.RegisteredSealProof) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.SealProofType = abi.RegisteredSealProof(extraI)
	}
	// t.AggregateProofType (abi.RegisteredAggregationProof) (int64)
	{
		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}
		var extraI int64
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative overflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.AggregateProofType = abi.RegisteredAggregationProof(extraI)
	}
	// t.ProvingDeadline (uint64) (uint64)

	{

		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.ProvingDeadline = uint64(extra)

	}
	// t.RequireActivationSuccess (bool) (bool)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.RequireActivationSuccess = false
	case 21:
		t.RequireActivationSuccess = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	return nil
}
