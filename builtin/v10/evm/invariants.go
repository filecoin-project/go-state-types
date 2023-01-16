package evm

import (
	"bytes"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
	"github.com/ipfs/go-cid"
	"golang.org/x/crypto/sha3"
	"golang.org/x/xerrors"
)

// Checks internal invariants of evm state.
func CheckStateInvariants(st *State, store adt.Store) *builtin.MessageAccumulator {
	acc := &builtin.MessageAccumulator{}

	acc.Require(st.Nonce > 0, "EVM actor state nonce needs to be greater than 0")

	byteCode, err := getBytecode(st.Bytecode, store)
	acc.RequireNoError(err, "Unable to retrieve bytecode")

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(byteCode)
	byteCodeHash := hasher.Sum(nil)

	acc.Require(bytes.Equal(byteCodeHash, st.BytecodeHash[:]), "Bytecode hash doesn't match bytecode cid, bytecode_hash: %x hash from bytecode cid: %x", st.BytecodeHash, byteCodeHash)

	return acc
}

func getBytecode(byteCodeCid cid.Cid, store adt.Store) ([]byte, error) {
	var bytecode abi.CborBytesTransparent
	if err := store.Get(store.Context(), byteCodeCid, &bytecode); err != nil {
		return nil, xerrors.Errorf("failed to get bytecode %w", err)
	}
	return bytecode, nil
}
