package paych

import (
	"bytes"

	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/crypto"
)

type UpdateChannelStateParams struct {
	Sv     SignedVoucher
	Secret []byte
}

// A voucher is sent by `From` to `To` off-chain in order to enable
// `To` to redeem payments on-chain in the future
type SignedVoucher struct {
	// ChannelAddr is the address of the payment channel this signed voucher is valid for
	ChannelAddr addr.Address
	// TimeLockMin sets a min epoch before which the voucher cannot be redeemed
	TimeLockMin abi.ChainEpoch
	// TimeLockMax sets a max epoch beyond which the voucher cannot be redeemed
	// TimeLockMax set to 0 means no timeout
	TimeLockMax abi.ChainEpoch
	// (optional) The SecretHash is used by `To` to validate
	SecretHash []byte
	// (optional) Extra can be specified by `From` to add a verification method to the voucher.
	Extra *ModVerifyParams
	// Specifies which lane the Voucher merges into (will be created if does not exist)
	Lane uint64
	// Nonce is set by `From` to prevent redemption of stale vouchers on a lane
	Nonce uint64
	// Amount voucher can be redeemed for
	Amount big.Int
	// (optional) MinSettleHeight can extend channel MinSettleHeight if needed
	MinSettleHeight abi.ChainEpoch

	// (optional) Set of lanes to be merged into `Lane`
	Merges []Merge

	// Sender's signature over the voucher
	Signature *crypto.Signature
}

func (t *SignedVoucher) SigningBytes() ([]byte, error) {
	osv := *t
	osv.Signature = nil

	buf := new(bytes.Buffer)
	if err := osv.MarshalCBOR(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type ConstructorParams struct {
	From addr.Address // Payer
	To   addr.Address // Payee
}

// Modular Verification method
type ModVerifyParams struct {
	// Actor on which to invoke the method.
	Actor addr.Address
	// Method to invoke.
	Method abi.MethodNum
	// Pre-serialized method parameters.
	Data []byte
}

//Specifies which `Lane`s to be merged with what `Nonce` on channelUpdate
type Merge struct {
	Lane  uint64
	Nonce uint64
}
