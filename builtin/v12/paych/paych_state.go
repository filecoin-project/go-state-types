package paych

import (
	"github.com/ipfs/go-cid"

	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
)

// A given payment channel actor is established by From
// to enable off-chain microtransactions to To to be reconciled
// and tallied on chain.
type State struct {
	// Channel owner, who has funded the actor
	From addr.Address
	// Recipient of payouts from channel
	To addr.Address

	// Amount successfully redeemed through the payment channel, paid out on `Collect()`
	ToSend abi.TokenAmount

	// Height at which the channel can be `Collected`
	SettlingAt abi.ChainEpoch
	// Height before which the channel `ToSend` cannot be collected
	MinSettleHeight abi.ChainEpoch

	// Collections of lane states for the channel, maintained in ID order.
	LaneStates cid.Cid // AMT<LaneState>
}

// The Lane state tracks the latest (highest) voucher nonce used to merge the lane
// as well as the amount it has already redeemed.
type LaneState struct {
	Redeemed big.Int
	Nonce    uint64
}

const LaneStatesAmtBitwidth = 3
