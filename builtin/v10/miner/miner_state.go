package miner

import (
	miner9 "github.com/filecoin-project/go-state-types/builtin/v9/miner"

	"github.com/filecoin-project/go-state-types/abi"
)

type State = miner9.State

const PrecommitCleanUpAmtBitwidth = miner9.PrecommitCleanUpAmtBitwidth
const SectorsAmtBitwidth = miner9.SectorsAmtBitwidth

type MinerInfo = miner9.MinerInfo
type WorkerKeyChange = miner9.WorkerKeyChange
type SectorPreCommitInfo = miner9.SectorPreCommitInfo
type SectorPreCommitOnChainInfo = miner9.SectorPreCommitOnChainInfo
type SectorOnChainInfo = miner9.SectorOnChainInfo
type BeneficiaryTerm = miner9.BeneficiaryTerm
type PendingBeneficiaryChange = miner9.PendingBeneficiaryChange

func SectorKey(e abi.SectorNumber) abi.Keyer {
	return miner9.SectorKey(e)
}
