package multisig

import (
	"bytes"

	"github.com/filecoin-project/go-state-types/exitcode"

	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
)

type TxnID int64

type Transaction struct {
	To     addr.Address
	Value  abi.TokenAmount
	Method abi.MethodNum
	Params []byte

	// This address at index 0 is the transaction proposer, order of this slice must be preserved.
	Approved []addr.Address
}

//Data for a BLAKE2B-256 to be attached to methods referencing proposals via TXIDs.
//Ensures the existence of a cryptographic reference to the original proposal. Useful
//for offline signers and for protection when reorgs change a multisig TXID.
//
//Requester - The requesting multisig wallet member.
//All other fields - From the "Transaction" struct.
type ProposalHashData struct {
	Requester addr.Address
	To        addr.Address
	Value     abi.TokenAmount
	Method    abi.MethodNum
	Params    []byte
}

type ConstructorParams struct {
	Signers               []addr.Address
	NumApprovalsThreshold uint64
	UnlockDuration        abi.ChainEpoch
	StartEpoch            abi.ChainEpoch
}

type ProposeParams struct {
	To     addr.Address
	Value  abi.TokenAmount
	Method abi.MethodNum
	Params []byte
}

type ProposeReturn struct {
	// TxnID is the ID of the proposed transaction
	TxnID TxnID
	// Applied indicates if the transaction was applied as opposed to proposed but not applied due to lack of approvals
	Applied bool
	// Code is the exitcode of the transaction, if Applied is false this field should be ignored.
	Code exitcode.ExitCode
	// Ret is the return vale of the transaction, if Applied is false this field should be ignored.
	Ret []byte
}

type TxnIDParams struct {
	ID TxnID
	// Optional hash of proposal to ensure an operation can only apply to a
	// specific proposal.
	ProposalHash []byte
}

type ApproveReturn struct {
	// Applied indicates if the transaction was applied as opposed to proposed but not applied due to lack of approvals
	Applied bool
	// Code is the exitcode of the transaction, if Applied is false this field should be ignored.
	Code exitcode.ExitCode
	// Ret is the return vale of the transaction, if Applied is false this field should be ignored.
	Ret []byte
}

type AddSignerParams struct {
	Signer   addr.Address
	Increase bool
}

type RemoveSignerParams struct {
	Signer   addr.Address
	Decrease bool
}

type SwapSignerParams struct {
	From addr.Address
	To   addr.Address
}

type ChangeNumApprovalsThresholdParams struct {
	NewThreshold uint64
}

type LockBalanceParams struct {
	StartEpoch     abi.ChainEpoch
	UnlockDuration abi.ChainEpoch
	Amount         abi.TokenAmount
}

func (phd *ProposalHashData) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := phd.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
