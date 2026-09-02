package reward

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
)

// StreamID identifies a reward stream. Zero is reserved; an ID becomes reusable after its
// retired stream's tombstone drains.
type StreamID uint64

const (
	// Denom is the fixed-point denominator for weights and recipient shares.
	Denom uint64 = 1_000_000_000_000_000_000
	// MaxStreams bounds live streams in actor state.
	MaxStreams = 8
	// MaxRecipients bounds one explicit distribution's share map and one claim request.
	MaxRecipients = 64
	// MaxPayableRowsPerStream bounds distinct recipients across current shares and carried balances.
	MaxPayableRowsPerStream = 2 * MaxRecipients
	// MaxTombstoneRows bounds existing and reserved carried-balance rows across retired streams.
	MaxTombstoneRows = 256
	// MaxPendingWrites bounds the complete deferred-write queue.
	MaxPendingWrites = MaxStreams*3 + 2
)

// PendingWriteOp is a stable operation tag persisted in PendingWrite and exposed by CancelPending.
type PendingWriteOp uint8

const (
	// PendingWriteOpSetWeightRecords identifies a cancellable, schedule-wide weight update.
	PendingWriteOpSetWeightRecords PendingWriteOp = iota
	// PendingWriteOpStepWeightRecords identifies an uncancellable, gate-originated weight update.
	PendingWriteOpStepWeightRecords
	// PendingWriteOpRegisterStream identifies a per-stream registration.
	PendingWriteOpRegisterStream
	// PendingWriteOpRemoveStream identifies a per-stream removal.
	PendingWriteOpRemoveStream
	// PendingWriteOpSetDistribution identifies a per-stream writer replacement.
	PendingWriteOpSetDistribution
)

// WeightRecord is a clamped linear weight in Denom fixed point.
//
// Persisted records satisfy Floor <= VStart <= Cap <= Denom.
type WeightRecord struct {
	// VStart is the weight at TStart, in Denom fixed point.
	VStart uint64
	// Slope is the signed weight change per epoch, in Denom fixed point.
	Slope int64
	// TStart is the epoch at which VStart applies.
	TStart abi.ChainEpoch
	// Floor is the inclusive lower clamp, in Denom fixed point.
	Floor uint64
	// Cap is the inclusive upper clamp, in Denom fixed point.
	Cap uint64
}

// WeightRecordUpdate is one stream update in SetWeightRecords and StepWeightRecords parameters.
type WeightRecordUpdate struct {
	// ID identifies the stream to update.
	ID StreamID
	// Weight completely replaces the stream's current weight record.
	Weight WeightRecord
}

// RecipientShare is one recipient entry in a share-map message and persisted distribution state.
type RecipientShare struct {
	// Recipient is the payout address, persisted in ID-address form.
	Recipient address.Address
	// Share is this recipient's portion of the stream, in Denom fixed point.
	Share uint64
}

// ExplicitDistribution is persisted allocation state for an explicit stream.
//
// The accounting rows are actor-owned state, not caller-supplied share-map fields.
type ExplicitDistribution struct {
	// Writer is the ID address authorized to install the next share map.
	Writer address.Address
	// Shares is the current period's complete recipient allocation, ordered by recipient.
	Shares []RecipientShare
	// Payable carries settled but unclaimed amounts from prior periods, ordered by recipient.
	Payable []RecipientAmount
	// ClaimedPeriod tracks current-period withdrawals, ordered by recipient.
	ClaimedPeriod []RecipientAmount
}

// RecipientAmount is a persisted recipient balance in a live distribution or tombstone.
type RecipientAmount struct {
	// Recipient is the payout address in ID-address form.
	Recipient address.Address
	// Amount is the recipient's positive token amount.
	Amount abi.TokenAmount
}

// Stream is a live stream persisted in StreamsState.
type Stream struct {
	// ID is the stream's stable, non-zero identifier.
	ID StreamID
	// Weight determines the stream's portion of each block reward.
	Weight WeightRecord
	// Distribution is nil for the implicit stream and non-nil for an explicit stream.
	Distribution *ExplicitDistribution
}

// Tombstone holds persisted liabilities for a removed stream.
type Tombstone struct {
	// ID is the retired stream identifier.
	ID StreamID
	// Payable carries its remaining unclaimed amounts, ordered by recipient.
	Payable []RecipientAmount
}

// PendingWrite is a deferred SWA operation persisted in StreamsState.
type PendingWrite struct {
	// ID is nil for schedule-wide operations and identifies the target for per-stream operations.
	ID *StreamID
	// Op identifies the deferred operation and its queue slot.
	Op PendingWriteOp
	// Payload is the operation's canonical CBOR tuple.
	Payload []byte
	// EffectiveEpoch is the first epoch at which the write is due.
	EffectiveEpoch abi.ChainEpoch
}

// DistributionInit is the caller-supplied subset of a new explicit distribution.
type DistributionInit struct {
	// Writer is the address authorized to install subsequent share maps.
	Writer address.Address
	// Shares is the initial complete recipient allocation and must sum to Denom.
	Shares []RecipientShare
}

// RegisterStreamPayload is the private tuple stored in PendingWrite.Payload for RegisterStream.
//
// It is not the public actor-method parameter tuple.
type RegisterStreamPayload struct {
	// Weight is the stream's initial weight record.
	Weight WeightRecord
	// Distribution is nil for an implicit stream or initializes an explicit stream.
	Distribution *DistributionInit
}

// SetDistributionPayload is the private tuple stored in PendingWrite.Payload for SetDistribution.
//
// It is not the public actor-method parameter tuple.
type SetDistributionPayload struct {
	// Writer is the replacement distribution writer.
	Writer address.Address
}

// StreamsState is persisted as the block referenced by State.StreamsRoot.
type StreamsState struct {
	// Streams contains live streams ordered by ID.
	Streams []Stream
	// Tombstones contains retired streams with liabilities, ordered by ID.
	Tombstones []Tombstone
	// PendingWrites is stably ordered by EffectiveEpoch.
	PendingWrites []PendingWrite
}

// StreamAccrual is the current-period gross accrual persisted inline in State.
//
// It stays outside StreamsState because it changes on every award.
type StreamAccrual struct {
	// ID identifies a live explicit stream.
	ID StreamID
	// Amount is its non-negative gross accrual for the current period.
	Amount abi.TokenAmount
}
