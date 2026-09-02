package reward

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v19/util/smoothing"
)

// AwardBlockRewardParams identifies a winning miner and the reward inputs for one block.
type AwardBlockRewardParams struct {
	// Miner receives the block reward and gas reward, less any penalty.
	Miner address.Address
	// Penalty is the non-negative penalty for including invalid messages.
	Penalty abi.TokenAmount
	// GasReward is the non-negative sum of gas rewards from the block.
	GasReward abi.TokenAmount
	// WinCount is the positive number of reward units won.
	WinCount int64
}

// ThisEpochRewardReturn reports the reward estimate and baseline power for the current epoch.
type ThisEpochRewardReturn struct {
	// ThisEpochRewardSmoothed is the smoothed reward per WinCount.
	ThisEpochRewardSmoothed smoothing.FilterEstimate
	// ThisEpochBaselinePower is the baseline power targeted in the current epoch.
	ThisEpochBaselinePower abi.StoragePower
}

// SetWeightRecordsParams queues a cancellable batch of replacement weight records.
type SetWeightRecordsParams struct {
	// Updates names each stream and its complete replacement weight record.
	Updates []WeightRecordUpdate
}

// StepWeightRecordsParams queues an uncancellable, gate-originated batch of replacement weight records.
type StepWeightRecordsParams struct {
	// Updates names each stream and its complete replacement weight record.
	Updates []WeightRecordUpdate
}

// RegisterStreamParams queues a new stream for activation.
type RegisterStreamParams struct {
	// ID is the non-zero stream identifier. It must not be live, tombstoned, or pending registration.
	ID StreamID
	// Weight is the stream's initial block-reward weight.
	Weight WeightRecord
	// Distribution is nil for an implicit stream or initializes an explicit stream.
	Distribution *DistributionInit
	// ActivationEpoch is when registration takes effect and must satisfy the SWA timelock.
	ActivationEpoch abi.ChainEpoch
}

// RemoveStreamParams queues removal of a live stream.
type RemoveStreamParams struct {
	// ID identifies the stream to remove.
	ID StreamID
}

// SetDistributionParams queues replacement of an explicit stream's writer.
type SetDistributionParams struct {
	// ID identifies the explicit stream.
	ID StreamID
	// Writer is the address authorized to install recipient shares after the change applies.
	Writer address.Address
}

// SetSharesParams closes an explicit stream's current period and installs its next share map.
type SetSharesParams struct {
	// ID identifies the explicit stream.
	ID StreamID
	// Shares is the complete next-period recipient allocation and must sum to Denom.
	Shares []RecipientShare
}

// CancelPendingParams identifies one queued operation to cancel after due writes are applied.
type CancelPendingParams struct {
	// ID is nil for schedule-wide weight updates and identifies a stream for per-stream operations.
	ID *StreamID
	// Op identifies the queue slot. PendingWriteOpStepWeightRecords cannot be cancelled.
	Op PendingWriteOp
}

// ClaimParams requests live and carried entitlements for wallets in one explicit stream.
type ClaimParams struct {
	// ID identifies a live explicit stream or its tombstone.
	ID StreamID
	// Wallets lists at most MaxRecipients payout addresses in the order expected in
	// ClaimReturn.Amounts. Unresolvable or ineligible addresses receive a positional zero.
	Wallets []address.Address
}

// ClaimReturn reports one claimed amount for each requested wallet.
type ClaimReturn struct {
	// Amounts preserves ClaimParams.Wallets order and uses zero for wallets with no payout.
	Amounts []abi.TokenAmount
}
