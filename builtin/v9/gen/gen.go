package main

import (
	"github.com/filecoin-project/go-state-types/builtin/v9/account"
	"github.com/filecoin-project/go-state-types/builtin/v9/cron"
	"github.com/filecoin-project/go-state-types/builtin/v9/datacap"
	init_ "github.com/filecoin-project/go-state-types/builtin/v9/init"
	"github.com/filecoin-project/go-state-types/builtin/v9/market"
	"github.com/filecoin-project/go-state-types/builtin/v9/migration"
	"github.com/filecoin-project/go-state-types/builtin/v9/miner"
	"github.com/filecoin-project/go-state-types/builtin/v9/multisig"
	"github.com/filecoin-project/go-state-types/builtin/v9/paych"
	"github.com/filecoin-project/go-state-types/builtin/v9/power"
	"github.com/filecoin-project/go-state-types/builtin/v9/reward"
	"github.com/filecoin-project/go-state-types/builtin/v9/system"
	"github.com/filecoin-project/go-state-types/builtin/v9/util/smoothing"
	"github.com/filecoin-project/go-state-types/builtin/v9/verifreg"
	gen "github.com/whyrusleeping/cbor-gen"
)

func main() {
	// Actors
	if err := gen.WriteTupleEncodersToFile("./builtin/v9/system/cbor_gen.go", "system",
		// actor state
		system.State{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v9/account/cbor_gen.go", "account",
		// actor state
		account.State{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v9/init/cbor_gen.go", "init",
		// actor state
		init_.State{},
		// method params and returns
		init_.ConstructorParams{},
		init_.ExecParams{},
		init_.ExecReturn{},
		init_.InstallParams{},
		init_.InstallReturn{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v9/cron/cbor_gen.go", "cron",
		// actor state
		cron.State{},
		cron.Entry{},
		// method params and returns
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v9/reward/cbor_gen.go", "reward",
		// actor state
		reward.State{},
		// method params and returns
		reward.ThisEpochRewardReturn{},
		reward.AwardBlockRewardParams{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v9/multisig/cbor_gen.go", "multisig",
		// actor state
		multisig.State{},
		multisig.Transaction{},
		multisig.ProposalHashData{},
		// method params and returns
		multisig.ConstructorParams{},
		multisig.ProposeParams{},
		multisig.ProposeReturn{},
		multisig.AddSignerParams{},
		multisig.RemoveSignerParams{},
		multisig.TxnIDParams{},
		multisig.ApproveReturn{},
		multisig.ChangeNumApprovalsThresholdParams{},
		multisig.SwapSignerParams{},
		multisig.LockBalanceParams{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v9/paych/cbor_gen.go", "paych",
		// actor state
		paych.State{},
		paych.LaneState{},
		//method params and returns
		paych.ConstructorParams{},
		paych.UpdateChannelStateParams{},
		paych.SignedVoucher{},
		paych.ModVerifyParams{},
		//other types
		paych.Merge{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v9/power/cbor_gen.go", "power",
		// actors state
		power.State{},
		power.Claim{},
		// method params and returns
		power.UpdateClaimedPowerParams{},
		power.MinerConstructorParams{},
		power.CreateMinerReturn{},
		power.CurrentTotalPowerReturn{},
		power.EnrollCronEventParams{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v9/market/cbor_gen.go", "market",
		// actor state
		market.State{},
		market.DealState{},
		// method params and returns
		market.WithdrawBalanceParams{},
		market.PublishStorageDealsParams{},
		market.PublishStorageDealsReturn{},
		market.ActivateDealsParams{},
		market.ActivateDealsResult{},
		market.VerifyDealsForActivationParams{},
		market.VerifyDealsForActivationReturn{},
		market.ComputeDataCommitmentParams{},
		market.ComputeDataCommitmentReturn{},
		market.OnMinerSectorsTerminateParams{},
		// other types
		market.DealProposal{},
		market.ClientDealProposal{},
		market.SectorDeals{},
		market.SectorDealData{},
		market.DealWeights{},
		market.SectorDataSpec{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v9/miner/cbor_gen.go", "miner",
		// actor state
		miner.State{},
		miner.MinerInfo{},
		miner.Deadlines{},
		miner.Deadline{},
		miner.Partition{},
		miner.ExpirationSet{},
		miner.PowerPair{},
		miner.SectorPreCommitOnChainInfo{},
		miner.SectorPreCommitInfo{},
		miner.SectorOnChainInfo{},
		miner.WorkerKeyChange{},
		miner.VestingFunds{},
		miner.VestingFund{},
		miner.WindowedPoSt{},
		miner.ActiveBeneficiary{},
		miner.BeneficiaryTerm{},
		miner.PendingBeneficiaryChange{},
		// method params and returns

		miner.GetControlAddressesReturn{},
		miner.ChangeWorkerAddressParams{},
		miner.ChangePeerIDParams{},
		miner.SubmitWindowedPoStParams{},
		miner.PreCommitSectorParams{},
		miner.ProveCommitSectorParams{},
		miner.ExtendSectorExpirationParams{},
		miner.TerminateSectorsParams{},
		miner.TerminateSectorsReturn{},
		miner.DeclareFaultsParams{},
		miner.DeclareFaultsRecoveredParams{},
		miner.DeferredCronEventParams{},
		miner.CheckSectorProvenParams{},
		miner.ApplyRewardParams{},
		miner.ReportConsensusFaultParams{},
		miner.WithdrawBalanceParams{},
		miner.ConfirmSectorProofsParams{},
		miner.ChangeMultiaddrsParams{},
		miner.CompactPartitionsParams{},
		miner.CompactSectorNumbersParams{},
		miner.DisputeWindowedPoStParams{},
		miner.PreCommitSectorBatchParams{},
		miner.ProveCommitAggregateParams{},
		miner.ProveReplicaUpdatesParams{},
		miner.CronEventPayload{},
		miner.PreCommitSectorBatchParams2{},
		miner.ProveReplicaUpdatesParams2{},
		miner.ChangeBeneficiaryParams{},
		miner.GetBeneficiaryReturn{},
		// other types
		miner.FaultDeclaration{},
		miner.RecoveryDeclaration{},
		miner.ExpirationExtension{},
		miner.TerminationDeclaration{},
		miner.PoStPartition{},
		miner.ReplicaUpdate{},
		miner.ReplicaUpdate2{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v9/verifreg/cbor_gen.go", "verifreg",
		// actor state
		verifreg.State{},

		// method params and returns
		verifreg.AddVerifierParams{},
		verifreg.AddVerifiedClientParams{},
		verifreg.UseBytesParams{},
		verifreg.RestoreBytesParams{},
		verifreg.RemoveDataCapParams{},
		verifreg.RemoveDataCapReturn{},
		verifreg.RemoveExpiredAllocationsParams{},
		verifreg.RemoveExpiredAllocationsReturn{},
		verifreg.BatchReturn{},
		verifreg.ClaimAllocationsParams{},
		verifreg.ClaimAllocationsReturn{},
		verifreg.GetClaimsParams{},
		verifreg.GetClaimsReturn{},
		verifreg.UniversalReceiverParams{},
		verifreg.AllocationsResponse{},
		// other types
		verifreg.RemoveDataCapRequest{},
		verifreg.RemoveDataCapProposal{},
		verifreg.RmDcProposalID{},
		verifreg.FailCode{},
		verifreg.SectorAllocationClaim{},
		verifreg.Claim{},
		verifreg.Allocation{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v9/datacap/cbor_gen.go", "datacap",
		// actor state
		datacap.State{},
		datacap.TokenState{},

		// method params and returns
		datacap.MintParams{},
		datacap.MintReturn{},
		datacap.DestroyParams{},
		datacap.TransferParams{},
		datacap.TransferReturn{},
		datacap.TransferFromParams{},
		datacap.TransferFromReturn{},
		datacap.IncreaseAllowanceParams{},
		datacap.DecreaseAllowanceParams{},
		datacap.RevokeAllowanceParams{},
		datacap.GetAllowanceParams{},
		datacap.BurnParams{},
		datacap.BurnReturn{},
		datacap.BurnFromParams{},
		datacap.BurnFromReturn{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v9/util/smoothing/cbor_gen.go", "smoothing",
		smoothing.FilterEstimate{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v9/migration/cbor_gen.go", "migration",
		migration.Actor{},
	); err != nil {
		panic(err)
	}
}
