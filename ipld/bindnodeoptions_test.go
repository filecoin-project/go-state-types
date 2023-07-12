package ipld_test

import (
	"bytes"
	_ "embed"
	"math/rand"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v8/paych"
	"github.com/filecoin-project/go-state-types/crypto"
	. "github.com/filecoin-project/go-state-types/ipld"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/stretchr/testify/assert"
)

var ipldSchema string = `
type SignedVoucher struct {
	ChannelAddr Bytes # addr.Address
	TimeLockMin Int # abi.ChainEpoch
	TimeLockMax Int # abi.ChainEpoch
	SecretHash Bytes
	Extra nullable ModVerifyParams
	Lane Int
	Nonce Int
	Amount Bytes # big.Int
	MinSettleHeight Int # abi.ChainEpoch
	Merges [Merge]
	Signature nullable Bytes # crypto.Signature
} representation tuple

type ModVerifyParams struct {
	Actor Bytes # addr.Address
	Method Int # abi.MethodNum
	Data Bytes
} representation tuple

type Merge struct {
	Lane Int
	Nonce Int
} representation tuple
`

// This test verifies compatibility between cbor-gen and ipld-prime serialization
// using the payment channel SignedVoucher as a sample type
func TestCborGenPrimeRoundTrip(t *testing.T) {

	// make a randomly generated payment channel voucher
	testVoucher := paych.SignedVoucher{
		ChannelAddr: address.TestAddress,
		TimeLockMin: abi.ChainEpoch(rand.Uint32()),
		TimeLockMax: abi.ChainEpoch(rand.Uint32()),
		SecretHash:  []byte("a secret"),
		Extra: &paych.ModVerifyParams{
			Actor:  address.TestAddress2,
			Method: abi.MethodNum(rand.Uint64()),
			Data:   []byte("some extra"),
		},
		Lane:            rand.Uint64(),
		Nonce:           rand.Uint64(),
		Amount:          big.NewInt(rand.Int63()),
		MinSettleHeight: abi.ChainEpoch(rand.Uint32()),
		Merges: []paych.Merge{{
			Lane:  rand.Uint64(),
			Nonce: rand.Uint64(),
		}},
		Signature: &crypto.Signature{
			Type: crypto.SigTypeSecp256k1,
			Data: []byte{1, 2, 3, 4},
		},
	}

	// serialize the payment channel voucher using cbor-gen
	buffer := new(bytes.Buffer)
	err := testVoucher.MarshalCBOR(buffer)
	assert.NoError(t, err, "should encode to CBOR with cbor-gen")

	// deserialize using bindnode in ipld prime
	typeSystem, err := ipld.LoadSchemaBytes([]byte(ipldSchema))
	assert.NoError(t, err, "should load a properly encoded schema")
	schemaType := typeSystem.TypeByName("SignedVoucher")
	assert.NotNil(t, schemaType, "should have the expected type")
	proto := bindnode.Prototype((*paych.SignedVoucher)(nil), schemaType, BigIntBindnodeOption, TokenAmountBindnodeOption, AddressBindnodeOption, SignatureBindnodeOption)

	nd, err := ipld.DecodeStreamingUsingPrototype(buffer, dagcbor.Decode, proto)
	assert.NoError(t, err, "should decode successfully from CBOR using bindnode")

	// unwrap deserialized node
	primeVoucher, ok := bindnode.Unwrap(nd).(*paych.SignedVoucher)
	assert.True(t, ok, "should unwrap to the correct go type")

	// verify equality with original voucher
	assert.Equal(t, &testVoucher, primeVoucher, "original go value should be unchanged after encoding with cbor-gen and decoding with bindnode")

	// write back to bytes with ipld-prime encoders
	err = ipld.EncodeStreaming(buffer, nd, dagcbor.Encode)
	assert.NoError(t, err, "should encode to CBOR with bindnode")

	// read back out with cbor
	var finalCborGenVoucher paych.SignedVoucher
	err = finalCborGenVoucher.UnmarshalCBOR(buffer)
	assert.NoError(t, err, "should decode successfully from CBOR using cbor-gen")

	// verify we still have the same as the original voucher
	assert.Equal(t, testVoucher, finalCborGenVoucher, "original go value should be unchanged after encoding with bindnode and decoding with cbor-gen")
}
