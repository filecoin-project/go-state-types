package datacap

import (
	datacap9 "github.com/filecoin-project/go-state-types/builtin/v9/datacap"
)

var InfiniteAllowance = datacap9.InfiniteAllowance

type MintParams = datacap9.MintParams
type DestroyParams = datacap9.DestroyParams
type MintReturn = datacap9.MintReturn
type TransferParams = datacap9.TransferParams
type TransferReturn = datacap9.TransferReturn
type TransferFromParams = datacap9.TransferFromParams
type TransferFromReturn = datacap9.TransferFromReturn
type IncreaseAllowanceParams = datacap9.IncreaseAllowanceParams
type DecreaseAllowanceParams = datacap9.DecreaseAllowanceParams
type RevokeAllowanceParams = datacap9.RevokeAllowanceParams
type GetAllowanceParams = datacap9.GetAllowanceParams
type BurnParams = datacap9.BurnParams
type BurnReturn = datacap9.BurnReturn
type BurnFromParams = datacap9.BurnFromParams
type BurnFromReturn = datacap9.BurnFromReturn
