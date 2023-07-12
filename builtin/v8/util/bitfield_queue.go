package util

import (
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/builtin"
)

// Wrapper for working with an AMT[ChainEpoch]*Bitfield functioning as a queue, bucketed by epoch.
// Keys in the queue are quantized (upwards), modulo some offset, to reduce the cardinality of keys.
type BitfieldQueue struct {
	*adt.Array
	quant builtin.QuantSpec
}

func LoadBitfieldQueue(store adt.Store, root cid.Cid, quant builtin.QuantSpec, bitwidth int) (BitfieldQueue, error) {
	arr, err := adt.AsArray(store, root, bitwidth)
	if err != nil {
		return BitfieldQueue{}, xerrors.Errorf("failed to load epoch queue %v: %w", root, err)
	}
	return BitfieldQueue{arr, quant}, nil
}

// Iterates the queue.
func (q BitfieldQueue) ForEach(cb func(epoch abi.ChainEpoch, bf bitfield.BitField) error) error {
	var bf bitfield.BitField
	return q.Array.ForEach(&bf, func(i int64) error {
		cpy, err := bf.Copy()
		if err != nil {
			return xerrors.Errorf("failed to copy bitfield in queue: %w", err)
		}
		return cb(abi.ChainEpoch(i), cpy)
	})
}
