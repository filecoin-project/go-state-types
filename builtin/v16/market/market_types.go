package market

import (
	"github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"

	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/batch"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v16/verifreg"
)

type SectorDealIDs []abi.DealID

type WithdrawBalanceParams struct {
	ProviderOrClientAddress addr.Address
	Amount                  abi.TokenAmount
}

type PublishStorageDealsParams struct {
	Deals []ClientDealProposal
}

type PublishStorageDealsReturn struct {
	IDs        []abi.DealID
	ValidDeals bitfield.BitField
}

// - Array of sectors rather than just one
// - Removed SectorStart (which is unknown at call time)
type VerifyDealsForActivationParams struct {
	Sectors []SectorDeals
}

type SectorDeals struct {
	SectorType   abi.RegisteredSealProof
	SectorExpiry abi.ChainEpoch
	DealIDs      []abi.DealID
}

// Changed to an array of optional CIDs.
type VerifyDealsForActivationReturn struct {
	UnsealedCIDs []*cid.Cid
}

type ActivateDealsParams struct {
	DealIDs      []abi.DealID
	SectorExpiry abi.ChainEpoch
}

type ActivateDealsResult struct {
	NonVerifiedDealSpace big.Int
	VerifiedInfos        []VerifiedDealInfo
}

type VerifiedDealInfo struct {
	Client       abi.ActorID
	AllocationId verifreg.AllocationId
	Data         cid.Cid
	Size         abi.PaddedPieceSize
}

type SectorDataSpec struct {
	DealIDs    []abi.DealID
	SectorType abi.RegisteredSealProof
}

type DealSpaces struct {
	DealSpace         abi.DealWeight // Total space of submitted deals.
	VerifiedDealSpace abi.DealWeight // Total space of submitted verified deals.
}

type OnMinerSectorsTerminateParams struct {
	Epoch   abi.ChainEpoch
	DealIDs []abi.DealID
}

type GetBalanceReturn struct {
	Balance abi.TokenAmount
	Locked  abi.TokenAmount
}

type DealQueryParams = cbg.CborInt // abi.DealID

type GetDealDataCommitmentParams = DealQueryParams

type GetDealDataCommitmentReturn struct {
	Data cid.Cid
	Size abi.PaddedPieceSize
}

type GetDealClientParams = DealQueryParams

type GetDealClientReturn = cbg.CborInt // abi.ActorID

type GetDealProviderParams = DealQueryParams

type GetDealProviderReturn = cbg.CborInt // abi.ActorID

type GetDealLabelParams = DealQueryParams

type GetDealLabelReturn = DealLabel

type GetDealTermParams = DealQueryParams

type GetDealTermReturn struct {
	Start    abi.ChainEpoch
	Duration abi.ChainEpoch
}

type GetDealTotalPriceParams = DealQueryParams

type GetDealTotalPriceReturn = abi.TokenAmount

type GetDealClientCollateralParams = DealQueryParams

type GetDealClientCollateralReturn = abi.TokenAmount

type GetDealProviderCollateralParams = DealQueryParams

type GetDealProviderCollateralReturn = abi.TokenAmount

type GetDealVerifiedParams = DealQueryParams

type GetDealVerifiedReturn = cbg.CborBool

type GetDealActivationParams = DealQueryParams

type GetDealActivationReturn struct {
	// Epoch at which the deal was activated, or -1.
	// This may be before the proposed start epoch.
	Activated abi.ChainEpoch
	// Epoch at which the deal was terminated abnormally, or -1.
	Terminated abi.ChainEpoch
}

type GetDealSectorParams = DealQueryParams

type GetDealSectorReturn = abi.SectorNumber

type SettleDealPaymentsParams = bitfield.BitField

type DealSettlementSummary struct {
	// Incremental amount paid to the provider.
	Payment abi.TokenAmount
	// Whether the deal has settled for the final time.
	Completed bool
}

type SettleDealPaymentsReturn struct {
	// Indicators of success or failure for each deal.
	Results batch.BatchReturn
	// Results for those deals that successfully settled.
	Settlements []DealSettlementSummary
}
