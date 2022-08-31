/*
Package ipld provides tools for working with Filecoin types as IPLD nodes.

These functions allow you to wrap Filecoin types into a go-ipld-prime node
interface using the bindnode package.
*/
package ipld

import (
	"fmt"

	"github.com/ipld/go-ipld-prime/node/bindnode"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/crypto"
)

// go type converter functions for bindnode for common Filecoin data types

// BigIntBindnodeOption converts a big.Int type to and from a Bytes field in a
// schema
var BigIntBindnodeOption = bindnode.TypedBytesConverter(&big.Int{}, bigIntFromBytes, bigIntToBytes)

// TokenAmountBindnodeOption converts a filecoin abi.TokenAmount type to and
// from a Bytes field in a schema
var TokenAmountBindnodeOption = bindnode.TypedBytesConverter(&abi.TokenAmount{}, tokenAmountFromBytes, tokenAmountToBytes)

// AddressBindnodeOption converts a filecoin Address type to and from a Bytes
// field in a schema
var AddressBindnodeOption = bindnode.TypedBytesConverter(&address.Address{}, addressFromBytes, addressToBytes)

// SignatureBindnodeOption converts a filecoin Signature type to and from a
// Bytes field in a schema
var SignatureBindnodeOption = bindnode.TypedBytesConverter(&crypto.Signature{}, signatureFromBytes, signatureToBytes)

func tokenAmountFromBytes(b []byte) (interface{}, error) {
	return bigIntFromBytes(b)
}

func bigIntFromBytes(b []byte) (interface{}, error) {
	if len(b) == 0 {
		return big.NewInt(0), nil
	}
	return big.FromBytes(b)
}

func tokenAmountToBytes(iface interface{}) ([]byte, error) {
	return bigIntToBytes(iface)
}

func bigIntToBytes(iface interface{}) ([]byte, error) {
	bi, ok := iface.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("expected *big.Int value")
	}
	if bi == nil || bi.Int == nil {
		*bi = big.Zero()
	}
	return bi.Bytes()
}

func addressFromBytes(b []byte) (interface{}, error) {
	return address.NewFromBytes(b)
}

func addressToBytes(iface interface{}) ([]byte, error) {
	addr, ok := iface.(*address.Address)
	if !ok {
		return nil, fmt.Errorf("expected *Address value")
	}
	return addr.Bytes(), nil
}

// Signature is a byteprefix union
func signatureFromBytes(b []byte) (interface{}, error) {
	if len(b) > crypto.SignatureMaxLength {
		return nil, fmt.Errorf("string too long")
	}
	if len(b) == 0 {
		return nil, fmt.Errorf("string empty")
	}
	var s crypto.Signature
	switch crypto.SigType(b[0]) {
	default:
		return nil, fmt.Errorf("invalid signature type in cbor input: %d", b[0])
	case crypto.SigTypeSecp256k1:
		s.Type = crypto.SigTypeSecp256k1
	case crypto.SigTypeBLS:
		s.Type = crypto.SigTypeBLS
	}
	s.Data = b[1:]
	return &s, nil
}

func signatureToBytes(iface interface{}) ([]byte, error) {
	s, ok := iface.(*crypto.Signature)
	if !ok {
		return nil, fmt.Errorf("expected *Signature value")
	}
	ba := append([]byte{byte(s.Type)}, s.Data...)
	return ba, nil
}
