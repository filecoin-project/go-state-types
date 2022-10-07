package multisig

import (
	multisig9 "github.com/filecoin-project/go-state-types/builtin/v9/multisig"
)

type TxnID = multisig9.TxnID

type Transaction = multisig9.Transaction
type ProposalHashData = multisig9.ProposalHashData
type ConstructorParams = multisig9.ConstructorParams
type ProposeParams = multisig9.ProposeParams
type ProposeReturn = multisig9.ProposeReturn
type TxnIDParams = multisig9.TxnIDParams
type ApproveReturn = multisig9.ApproveReturn
type AddSignerParams = multisig9.AddSignerParams
type RemoveSignerParams = multisig9.RemoveSignerParams
type SwapSignerParams = multisig9.SwapSignerParams
type ChangeNumApprovalsThresholdParams = multisig9.ChangeNumApprovalsThresholdParams
type LockBalanceParams = multisig9.LockBalanceParams
