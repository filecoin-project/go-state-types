package main

import (
	"github.com/filecoin-project/go-state-types/builtin/v8/account"
	"github.com/filecoin-project/go-state-types/builtin/v8/cron"
	init_ "github.com/filecoin-project/go-state-types/builtin/v8/init"
	"github.com/filecoin-project/go-state-types/builtin/v8/market"
	"github.com/filecoin-project/go-state-types/builtin/v8/miner"
	"github.com/filecoin-project/go-state-types/builtin/v8/multisig"
	"github.com/filecoin-project/go-state-types/builtin/v8/paych"
	"github.com/filecoin-project/go-state-types/builtin/v8/power"
	"github.com/filecoin-project/go-state-types/builtin/v8/reward"
	"github.com/filecoin-project/go-state-types/builtin/v8/system"
	"github.com/filecoin-project/go-state-types/builtin/v8/util/smoothing"
	"github.com/filecoin-project/go-state-types/builtin/v8/verifreg"
	gen "github.com/whyrusleeping/cbor-gen"
)

func main() {
	// Actors
	if err := gen.WriteTupleEncodersToFile("./builtin/v8/system/cbor_gen.go", "system",
		// actor state
		system.State{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v8/account/cbor_gen.go", "account",
		// actor state
		account.State{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v8/init/cbor_gen.go", "init",
		// actor state
		init_.State{},
		// method params and returns
		init_.ConstructorParams{},
		init_.ExecParams{},
		init_.ExecReturn{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v8/cron/cbor_gen.go", "cron",
		// actor state
		cron.State{},
		cron.Entry{},
		// method params and returns
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v8/reward/cbor_gen.go", "reward",
		// actor state
		reward.State{},
		// method params and returns
		reward.ThisEpochRewardReturn{},
		reward.AwardBlockRewardParams{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v8/multisig/cbor_gen.go", "multisig",
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

	if err := gen.WriteTupleEncodersToFile("./builtin/v8/paych/cbor_gen.go", "paych",
		// actor state
		paych.State{},
		paych.LaneState{},
		// method params and returns
		paych.ConstructorParams{},
		paych.UpdateChannelStateParams{},
		paych.SignedVoucher{},
		paych.ModVerifyParams{},
		// other types
		paych.Merge{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v8/power/cbor_gen.go", "power",
		// actors state
		power.State{},
		power.Claim{},
		// method params and returns
		power.UpdateClaimedPowerParams{},
		power.MinerConstructorParams{},
		power.CreateMinerParams{},
		power.CreateMinerReturn{},
		power.CurrentTotalPowerReturn{},
		power.EnrollCronEventParams{},
		// other types
		power.CronEvent{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v8/market/cbor_gen.go", "market",
		// actor state
		market.State{},
		market.DealState{},
		// method params and returns
		market.WithdrawBalanceParams{},
		market.PublishStorageDealsParams{},
		market.PublishStorageDealsReturn{},
		market.ActivateDealsParams{},
		market.VerifyDealsForActivationParams{},
		market.VerifyDealsForActivationReturn{},
		market.ComputeDataCommitmentParams{},
		market.ComputeDataCommitmentReturn{},
		market.OnMinerSectorsTerminateParams{},
		// other types
		market.DealProposal{},
		market.ClientDealProposal{},
		market.SectorDeals{},
		market.SectorWeights{},
		market.SectorDataSpec{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v8/miner/cbor_gen.go", "miner",
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
		// other types
		miner.FaultDeclaration{},
		miner.RecoveryDeclaration{},
		miner.ExpirationExtension{},
		miner.TerminationDeclaration{},
		miner.PoStPartition{},
		miner.ReplicaUpdate{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v8/verifreg/cbor_gen.go", "verifreg",
		// actor state
		verifreg.State{},

		// method params and returns
		verifreg.AddVerifierParams{},
		verifreg.AddVerifiedClientParams{},
		verifreg.UseBytesParams{},
		verifreg.RestoreBytesParams{},
		verifreg.RemoveDataCapParams{},
		verifreg.RemoveDataCapReturn{},
		// other types
		verifreg.RemoveDataCapRequest{},
		verifreg.RemoveDataCapProposal{},
	); err != nil {
		panic(err)
	}

	if err := gen.WriteTupleEncodersToFile("./builtin/v8/util/smoothing/cbor_gen.go", "smoothing",
		smoothing.FilterEstimate{},
	); err != nil {
		panic(err)
	}
}
