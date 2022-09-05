package market

import (
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
)

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
	SectorExpiry abi.ChainEpoch
	DealIDs      []abi.DealID
}

// - Array of sectors weights
type VerifyDealsForActivationReturn struct {
	Sectors []SectorDealData
}

type SectorDealData struct {
	CommD *cid.Cid
}

type ActivateDealsParams struct {
	DealIDs      []abi.DealID
	SectorExpiry abi.ChainEpoch
}

type ActivateDealsResult struct {
	Weights DealWeights
}

type SectorDataSpec struct {
	DealIDs    []abi.DealID
	SectorType abi.RegisteredSealProof
}

type DealWeights struct {
	DealSpace          uint64         // Total space in bytes of submitted deals.
	DealWeight         abi.DealWeight // Total space*time of submitted deals.
	VerifiedDealWeight abi.DealWeight // Total space*time of submitted verified deals.
}

type ComputeDataCommitmentParams struct {
	Inputs []*SectorDataSpec
}

type ComputeDataCommitmentReturn struct {
	CommDs []cbg.CborCid
}

type OnMinerSectorsTerminateParams struct {
	Epoch   abi.ChainEpoch
	DealIDs []abi.DealID
}
