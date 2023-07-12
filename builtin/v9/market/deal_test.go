package market_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/filecoin-project/go-state-types/builtin/v9/market"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs/go-cid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDealLabel(t *testing.T) {

	// empty label serialize and deserialize to string label of empty string
	label1 := market.EmptyDealLabel
	buf := bytes.Buffer{}
	require.NoError(t, label1.MarshalCBOR(&buf), "failed to marshal empty deal label")
	ser := buf.Bytes()
	assert.Equal(t, []byte{0x60}, ser) // cbor empty str (maj type 3)
	label2 := &market.DealLabel{}
	require.NoError(t, label2.UnmarshalCBOR(bytes.NewReader(ser)))
	assert.True(t, label2.IsString())
	assert.False(t, label2.IsBytes())
	str, err := label2.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	// bytes label
	label1, err = market.NewLabelFromBytes([]byte{0xca, 0xfe, 0xb0, 0x0a})
	assert.NoError(t, err)
	buf = bytes.Buffer{}
	require.NoError(t, label1.MarshalCBOR(&buf), "failed to marshal bytes deal label")
	ser = buf.Bytes()
	label2 = &market.DealLabel{}
	require.NoError(t, label2.UnmarshalCBOR(bytes.NewReader(ser)))
	assert.True(t, label2.IsBytes())
	assert.False(t, label2.IsString())
	bs, err := label2.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, []byte{0xca, 0xfe, 0xb0, 0x0a}, bs)

	// string label
	label1, err = market.NewLabelFromString("i am a label, turn me into cbor maj typ 3 plz")
	assert.NoError(t, err)
	buf = bytes.Buffer{}
	require.NoError(t, label1.MarshalCBOR(&buf), "failed to marshal string deal label")
	ser = buf.Bytes()
	label2 = &market.DealLabel{}
	require.NoError(t, label2.UnmarshalCBOR(bytes.NewReader(ser)))
	assert.True(t, label2.IsString())
	assert.False(t, label2.IsBytes())
	str, err = label2.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "i am a label, turn me into cbor maj typ 3 plz", str)

	// invalid utf8 string
	_, err = market.NewLabelFromString(string([]byte{0xde, 0xad, 0xbe, 0xef}))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid utf8")

	// nil label marshals
	labelPtr := (*market.DealLabel)(nil)
	buf = bytes.Buffer{}
	require.NoError(t, labelPtr.MarshalCBOR(&buf), "failed to marshal empty deal label")
	ser = buf.Bytes()
	assert.Equal(t, []byte{0x60}, ser) // cbor empty str (maj type 3)

	// nil label unmarshal fails nicely
	labelPtr = (*market.DealLabel)(nil)
	err = labelPtr.UnmarshalCBOR(bytes.NewReader(ser))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot unmarshal into nil pointer")

}
func TestDealLabelJSON(t *testing.T) {

	// Non-empty string
	label1, err := market.NewLabelFromString("i am a label, json me correctly plz")
	require.NoError(t, err, "failed to create label from string")
	label1JSON, err := json.Marshal(&label1)
	require.NoError(t, err, "failed to JSON marshal string label")
	label2 := &market.DealLabel{}
	require.NoError(t, label2.UnmarshalJSON(label1JSON))
	assert.True(t, label2.IsString())
	assert.False(t, label2.IsBytes())
	str, err := label2.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "i am a label, json me correctly plz", str)

	dp := &market.DealProposal{
		PieceCID:             cid.Undef,
		PieceSize:            0,
		VerifiedDeal:         false,
		Client:               address.Undef,
		Provider:             address.Undef,
		Label:                label1,
		StoragePricePerEpoch: big.Zero(),
		ProviderCollateral:   big.Zero(),
		ClientCollateral:     big.Zero(),
	}

	dpJSON, err := json.Marshal(dp)
	require.NoError(t, err, "failed to JSON marshal deal proposal")
	dp2 := market.DealProposal{}
	require.NoError(t, json.Unmarshal(dpJSON, &dp2))
	assert.True(t, dp2.Label.IsString())
	assert.False(t, dp2.Label.IsBytes())
	str, err = dp2.Label.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "i am a label, json me correctly plz", str)

	// empty label
	label1, err = market.NewLabelFromString("")
	require.NoError(t, err, "failed to create label from string")
	label1JSON, err = json.Marshal(&label1)
	require.NoError(t, err, "failed to JSON marshal string label")
	label2 = &market.DealLabel{}
	require.NoError(t, label2.UnmarshalJSON(label1JSON))
	assert.True(t, label2.IsString())
	assert.False(t, label2.IsBytes())
	str, err = label2.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	dp = &market.DealProposal{
		PieceCID:             cid.Undef,
		PieceSize:            0,
		VerifiedDeal:         false,
		Client:               address.Undef,
		Provider:             address.Undef,
		Label:                label1,
		StoragePricePerEpoch: big.Zero(),
		ProviderCollateral:   big.Zero(),
		ClientCollateral:     big.Zero(),
	}

	dpJSON, err = json.Marshal(dp)
	require.NoError(t, err, "failed to JSON marshal deal proposal")
	dp2 = market.DealProposal{}
	require.NoError(t, json.Unmarshal(dpJSON, &dp2))
	assert.True(t, dp2.Label.IsString())
	assert.False(t, dp2.Label.IsBytes())
	str, err = dp2.Label.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	// unset label
	dp = &market.DealProposal{
		PieceCID:             cid.Undef,
		PieceSize:            0,
		VerifiedDeal:         false,
		Client:               address.Undef,
		Provider:             address.Undef,
		StoragePricePerEpoch: big.Zero(),
		ProviderCollateral:   big.Zero(),
		ClientCollateral:     big.Zero(),
	}

	dpJSON, err = json.Marshal(dp)
	require.NoError(t, err, "failed to JSON marshal deal proposal")
	dp2 = market.DealProposal{}
	require.NoError(t, json.Unmarshal(dpJSON, &dp2))
	assert.True(t, dp2.Label.IsString())
	assert.False(t, dp2.Label.IsBytes())
	str, err = dp2.Label.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	// bytes label
	label1, err = market.NewLabelFromBytes([]byte{0xde, 0xad, 0xbe, 0xef})
	assert.NoError(t, err)
	assert.True(t, label1.IsBytes())
	assert.False(t, label1.IsString())

	dp = &market.DealProposal{
		PieceCID:             cid.Undef,
		PieceSize:            0,
		VerifiedDeal:         false,
		Client:               address.Undef,
		Provider:             address.Undef,
		Label:                label1,
		StoragePricePerEpoch: big.Zero(),
		ProviderCollateral:   big.Zero(),
		ClientCollateral:     big.Zero(),
	}

	dpJSON, err = json.Marshal(dp)
	require.NoError(t, err, "failed to JSON marshal deal proposal")
	dp2 = market.DealProposal{}
	require.NoError(t, json.Unmarshal(dpJSON, &dp2))
	// a JSON-unmarshaled label is always string type
	assert.False(t, dp2.Label.IsBytes())
	assert.True(t, dp2.Label.IsString())

	// it can be converted ToString, but it's empty
	str, err = dp2.Label.ToString()
	assert.NoError(t, err)
	assert.Equal(t, str, "")
}

func TestDealLabelFromCBOR(t *testing.T) {
	// empty string
	// b011_00000
	emptyCBORText := []byte{0x60}
	var label1 market.DealLabel
	require.NoError(t, label1.UnmarshalCBOR(bytes.NewReader(emptyCBORText)))
	assert.True(t, label1.IsString())
	str, err := label1.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	// valid utf8 string
	// b011_01000 "deadbeef"
	cborTestValid := append([]byte{0x68}, []byte("deadbeef")...)
	var label2 market.DealLabel
	require.NoError(t, label2.UnmarshalCBOR(bytes.NewReader(cborTestValid)))
	assert.True(t, label2.IsString())
	str, err = label2.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "deadbeef", str)

	// invalid utf8 string
	// b011_00100 0xde 0xad 0xbe 0xef
	cborText := []byte{0x64, 0xde, 0xad, 0xbe, 0xef}
	var label3 market.DealLabel
	err = label3.UnmarshalCBOR(bytes.NewReader(cborText))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not valid utf8")

	// empty bytes
	// b010_00000
	emptyCBORBytes := []byte{0x40}
	var label4 market.DealLabel
	require.NoError(t, label4.UnmarshalCBOR(bytes.NewReader(emptyCBORBytes)))
	assert.True(t, label4.IsBytes())
	bs, err := label4.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, []byte{}, bs)

	// bytes
	// b010_00100 0xde 0xad 0xbe 0xef
	cborBytes := []byte{0x44, 0xde, 0xad, 0xbe, 0xef}
	var label5 market.DealLabel
	require.NoError(t, label5.UnmarshalCBOR(bytes.NewReader(cborBytes)))
	assert.True(t, label5.IsBytes())
	bs, err = label5.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, []byte{0xde, 0xad, 0xbe, 0xef}, bs)

	// bad major type
	// array of empty array b100_00001 b100_00000
	arrayBytes := []byte{0x81, 0x80}
	var label6 market.DealLabel
	err = label6.UnmarshalCBOR(bytes.NewReader(arrayBytes))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected major tag")
}
