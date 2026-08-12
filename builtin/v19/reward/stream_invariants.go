package reward

import (
	"bytes"
	"fmt"
	"io"
	"math"
	mathbig "math/big"
	"sort"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
)

// The invariant rules in this file mirror actors/reward/src/streams.rs.

// Mirrors actors/reward/src/streams.rs::compute_weight.
func computeWeight(record WeightRecord, epoch abi.ChainEpoch) uint64 {
	delta := new(mathbig.Int).Sub(mathbig.NewInt(int64(epoch)), mathbig.NewInt(int64(record.TStart)))
	value := new(mathbig.Int).Mul(mathbig.NewInt(record.Slope), delta)
	value.Add(value, new(mathbig.Int).SetUint64(record.VStart))
	floor := new(mathbig.Int).SetUint64(record.Floor)
	cap := new(mathbig.Int).SetUint64(record.Cap)
	if value.Cmp(cap) > 0 {
		value = cap
	}
	if value.Cmp(floor) < 0 {
		value = floor
	}
	return value.Uint64()
}

// Mirrors actors/reward/src/streams.rs::validate_weight_record.
func validateWeightRecord(record WeightRecord) error {
	if record.Floor > record.Cap {
		return fmt.Errorf("weight floor exceeds cap")
	}
	if record.VStart < record.Floor {
		return fmt.Errorf("weight v_start is below floor")
	}
	if record.VStart > record.Cap {
		return fmt.Errorf("weight v_start exceeds cap")
	}
	if record.Cap > Denom {
		return fmt.Errorf("weight cap exceeds DENOM")
	}
	return nil
}

// Mirrors actors/reward/src/streams.rs::validate_weight_updates.
func validateWeightUpdates(updates []WeightRecordUpdate) error {
	if len(updates) == 0 {
		return fmt.Errorf("weight-record update is empty")
	}
	for i, update := range updates {
		if i > 0 {
			if updates[i-1].ID == update.ID {
				return fmt.Errorf("duplicate weight-record stream ID %d", update.ID)
			}
			if updates[i-1].ID > update.ID {
				return fmt.Errorf("weight-record updates are not ordered")
			}
		}
		if err := validateWeightRecord(update.Weight); err != nil {
			return err
		}
	}
	return nil
}

// Mirrors actors/reward/src/streams.rs::weight_breakpoints.
func weightBreakpoints(record WeightRecord, startEpoch abi.ChainEpoch) []abi.ChainEpoch {
	epochs := map[abi.ChainEpoch]struct{}{startEpoch: {}}
	if record.TStart >= startEpoch {
		epochs[record.TStart] = struct{}{}
	}
	if record.Slope != 0 {
		insertCrossing := func(bound uint64) {
			numerator := new(mathbig.Int).Sub(new(mathbig.Int).SetUint64(bound), new(mathbig.Int).SetUint64(record.VStart))
			quotient := new(mathbig.Int).Quo(numerator, mathbig.NewInt(record.Slope))
			for _, offset := range []int64{-1, 0, 1} {
				epoch := new(mathbig.Int).Add(mathbig.NewInt(int64(record.TStart)), quotient)
				epoch.Add(epoch, mathbig.NewInt(offset))
				if epoch.IsInt64() {
					candidate := abi.ChainEpoch(epoch.Int64())
					if candidate >= startEpoch {
						epochs[candidate] = struct{}{}
					}
				}
			}
		}
		if record.Slope > 0 {
			insertCrossing(record.Cap)
			if startEpoch < record.TStart {
				insertCrossing(record.Floor)
			}
		} else {
			insertCrossing(record.Floor)
			if startEpoch < record.TStart {
				insertCrossing(record.Cap)
			}
		}
	}

	var last abi.ChainEpoch = math.MinInt64
	for epoch := range epochs {
		if epoch > last {
			last = epoch
		}
	}
	if last < math.MaxInt64 {
		epochs[last+1] = struct{}{}
	}
	epochs[math.MaxInt64] = struct{}{}

	out := make([]abi.ChainEpoch, 0, len(epochs))
	for epoch := range epochs {
		out = append(out, epoch)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// Mirrors actors/reward/src/streams.rs::validate_weight_schedule.
func validateWeightSchedule(streams []Stream, startEpoch abi.ChainEpoch) error {
	epochs := map[abi.ChainEpoch]struct{}{startEpoch: {}, math.MaxInt64: {}}
	for _, stream := range streams {
		if err := validateWeightRecord(stream.Weight); err != nil {
			return err
		}
		for _, epoch := range weightBreakpoints(stream.Weight, startEpoch) {
			epochs[epoch] = struct{}{}
		}
	}
	for epoch := range epochs {
		var sum uint64
		for _, stream := range streams {
			sum += computeWeight(stream.Weight, epoch)
		}
		if sum > Denom {
			return fmt.Errorf("stream weights exceed DENOM at epoch %d: %d", epoch, sum)
		}
	}
	return nil
}

// Mirrors actors/reward/src/streams.rs::validate_id_address.
func validateIDAddress(addr address.Address, label string) error {
	if addr.Protocol() != address.ID {
		return fmt.Errorf("%s %s is not an ID address", label, addr)
	}
	return nil
}

// Mirrors actors/reward/src/streams.rs::validate_shares.
func validateShares(shares []RecipientShare) error {
	if len(shares) > MaxRecipients {
		return fmt.Errorf("recipient count %d exceeds maximum %d", len(shares), MaxRecipients)
	}
	recipients := make(map[address.Address]struct{}, len(shares))
	var total uint64
	for _, row := range shares {
		if err := validateIDAddress(row.Recipient, "share recipient"); err != nil {
			return err
		}
		if row.Share == 0 {
			return fmt.Errorf("share for recipient %s is zero", row.Recipient)
		}
		if _, found := recipients[row.Recipient]; found {
			return fmt.Errorf("duplicate share recipient %s", row.Recipient)
		}
		recipients[row.Recipient] = struct{}{}
		if row.Share > Denom-total {
			return fmt.Errorf("shares exceed DENOM")
		}
		total += row.Share
	}
	if total != Denom {
		return fmt.Errorf("shares sum to %d, expected %d", total, Denom)
	}
	return nil
}

func addressRowsStrictlyOrdered[T any](rows []T, addressOf func(T) address.Address) bool {
	for i := 1; i < len(rows); i++ {
		if bytes.Compare(addressOf(rows[i-1]).Bytes(), addressOf(rows[i]).Bytes()) >= 0 {
			return false
		}
	}
	return true
}

// Mirrors actors/reward/src/streams.rs::validate_amount_rows.
func validateAmountRows(rows []RecipientAmount, label string) error {
	if !addressRowsStrictlyOrdered(rows, func(row RecipientAmount) address.Address { return row.Recipient }) {
		return fmt.Errorf("%s recipients are not ordered", label)
	}
	for _, row := range rows {
		if err := validateIDAddress(row.Recipient, label); err != nil {
			return err
		}
		if !row.Amount.GreaterThan(big.Zero()) {
			return fmt.Errorf("%s amount is not positive", label)
		}
	}
	return nil
}

// Mirrors actors/reward/src/streams.rs::validate_period_claims.
func validatePeriodClaims(distribution *ExplicitDistribution, pool abi.TokenAmount) error {
	if err := validateAmountRows(distribution.Payable, "payable"); err != nil {
		return err
	}
	if err := validateAmountRows(distribution.ClaimedPeriod, "claimed-period"); err != nil {
		return err
	}
	for _, claimed := range distribution.ClaimedPeriod {
		var share *RecipientShare
		for i := range distribution.Shares {
			if distribution.Shares[i].Recipient == claimed.Recipient {
				share = &distribution.Shares[i]
				break
			}
		}
		if share == nil {
			return fmt.Errorf("claimed-period recipient is absent from shares")
		}
		earned := big.Div(big.Mul(pool, big.NewIntUnsigned(share.Share)), big.NewIntUnsigned(Denom))
		if claimed.Amount.GreaterThan(earned) {
			return fmt.Errorf("claimed amount exceeds earnings for recipient %s", claimed.Recipient)
		}
	}
	return nil
}

// Mirrors actors/reward/src/streams.rs::recipient_union_len.
func recipientUnionLen(payable []RecipientAmount, shares []RecipientShare) int {
	recipients := make(map[address.Address]struct{}, len(payable)+len(shares))
	for _, row := range payable {
		recipients[row.Recipient] = struct{}{}
	}
	for _, row := range shares {
		recipients[row.Recipient] = struct{}{}
	}
	return len(recipients)
}

// Mirrors actors/reward/src/streams.rs::validate_stream_configuration_without_weights.
func validateStreamConfigurationWithoutWeights(streams []Stream) error {
	if len(streams) > MaxStreams {
		return fmt.Errorf("stream count exceeds maximum %d", MaxStreams)
	}
	implicit := 0
	for i, stream := range streams {
		if stream.ID == 0 {
			return fmt.Errorf("stream ID 0 is reserved")
		}
		if i > 0 && streams[i-1].ID >= stream.ID {
			return fmt.Errorf("stream IDs are not ordered")
		}
		if stream.Distribution == nil {
			implicit++
			continue
		}
		distribution := stream.Distribution
		if err := validateIDAddress(distribution.Writer, "distribution writer"); err != nil {
			return err
		}
		if !addressRowsStrictlyOrdered(distribution.Shares, func(row RecipientShare) address.Address { return row.Recipient }) {
			return fmt.Errorf("share recipients are not ordered")
		}
		if err := validateShares(distribution.Shares); err != nil {
			return err
		}
		if err := validateAmountRows(distribution.Payable, "payable"); err != nil {
			return err
		}
		if err := validateAmountRows(distribution.ClaimedPeriod, "claimed-period"); err != nil {
			return err
		}
		reservedRows := recipientUnionLen(distribution.Payable, distribution.Shares)
		if reservedRows > MaxPayableRowsPerStream {
			return fmt.Errorf("stream %d payable row reservation %d exceeds maximum %d", stream.ID, reservedRows, MaxPayableRowsPerStream)
		}
		if len(distribution.ClaimedPeriod) > MaxRecipients {
			return fmt.Errorf("stream %d claimed-period row count %d exceeds maximum %d", stream.ID, len(distribution.ClaimedPeriod), MaxRecipients)
		}
	}
	if implicit > 1 {
		return fmt.Errorf("multiple implicit streams")
	}
	return nil
}

// Mirrors actors/reward/src/streams.rs::validate_stream_configuration.
func validateStreamConfiguration(streams []Stream) error {
	if err := validateStreamConfigurationWithoutWeights(streams); err != nil {
		return err
	}
	for _, stream := range streams {
		if err := validateWeightRecord(stream.Weight); err != nil {
			return err
		}
	}
	return nil
}

type cborUnmarshaler interface {
	UnmarshalCBOR(io.Reader) error
}

func unmarshalPendingPayload(payload []byte, out cborUnmarshaler) error {
	reader := bytes.NewReader(payload)
	if err := out.UnmarshalCBOR(reader); err != nil {
		return err
	}
	if reader.Len() != 0 {
		return fmt.Errorf("trailing CBOR data")
	}
	return nil
}

// Mirrors actors/reward/src/streams.rs::validate_distribution_init.
func validateDistributionInit(distribution *DistributionInit) error {
	if distribution == nil {
		return nil
	}
	if err := validateIDAddress(distribution.Writer, "distribution writer"); err != nil {
		return err
	}
	return validateShares(distribution.Shares)
}

// Mirrors actors/reward/src/streams.rs::validate_pending_payload.
func validatePendingPayload(write PendingWrite) error {
	switch write.Op {
	case PendingWriteOpSetWeightRecords, PendingWriteOpStepWeightRecords:
		var payload SetWeightRecordsParams
		if err := unmarshalPendingPayload(write.Payload, &payload); err != nil {
			return err
		}
		return validateWeightUpdates(payload.Updates)
	case PendingWriteOpRegisterStream:
		var payload RegisterStreamPayload
		if err := unmarshalPendingPayload(write.Payload, &payload); err != nil {
			return err
		}
		if err := validateWeightRecord(payload.Weight); err != nil {
			return err
		}
		return validateDistributionInit(payload.Distribution)
	case PendingWriteOpRemoveStream:
		if !bytes.Equal(write.Payload, []byte{0x80}) {
			return fmt.Errorf("RemoveStream payload is not an empty tuple")
		}
		return nil
	case PendingWriteOpSetDistribution:
		var payload SetDistributionPayload
		if err := unmarshalPendingPayload(write.Payload, &payload); err != nil {
			return err
		}
		return validateIDAddress(payload.Writer, "distribution writer")
	default:
		return fmt.Errorf("unknown pending operation %d", write.Op)
	}
}

type pendingSlot struct {
	id    StreamID
	hasID bool
	op    PendingWriteOp
}

func slotForWrite(write PendingWrite) pendingSlot {
	slot := pendingSlot{op: write.Op, hasID: write.ID != nil}
	if write.ID != nil {
		slot.id = *write.ID
	}
	return slot
}

// Mirrors actors/reward/src/streams.rs::validate_pending_queue.
func validatePendingQueue(writes []PendingWrite) error {
	if len(writes) > MaxPendingWrites {
		return fmt.Errorf("pending write count %d exceeds maximum %d", len(writes), MaxPendingWrites)
	}
	slots := make(map[pendingSlot]struct{}, len(writes))
	for i, write := range writes {
		if i > 0 && writes[i-1].EffectiveEpoch > write.EffectiveEpoch {
			return fmt.Errorf("pending writes are not ordered")
		}
		isSchedule := write.Op == PendingWriteOpSetWeightRecords || write.Op == PendingWriteOpStepWeightRecords
		if isSchedule != (write.ID == nil) {
			return fmt.Errorf("pending call (%v, %d) has a non-canonical stream ID", write.ID, write.Op)
		}
		if err := validatePendingPayload(write); err != nil {
			return err
		}
		slot := slotForWrite(write)
		if _, found := slots[slot]; found {
			return fmt.Errorf("duplicate pending slot (%v, %d)", write.ID, write.Op)
		}
		slots[slot] = struct{}{}
	}
	return nil
}

// Mirrors actors/reward/src/streams.rs::validate_tombstone_capacity.
func validateTombstoneCapacity(streams *StreamsState) error {
	rows := 0
	for _, tombstone := range streams.Tombstones {
		rows += len(tombstone.Payable)
	}
	for _, write := range streams.PendingWrites {
		if write.Op != PendingWriteOpRemoveStream {
			continue
		}
		reserved := MaxRecipients
		if write.ID != nil {
			for _, stream := range streams.Streams {
				if stream.ID == *write.ID && stream.Distribution != nil {
					reserved = max(reserved, recipientUnionLen(stream.Distribution.Payable, stream.Distribution.Shares))
					break
				}
			}
		}
		rows += reserved
	}
	if rows > MaxTombstoneRows {
		return fmt.Errorf("tombstone row reservation %d exceeds maximum %d", rows, MaxTombstoneRows)
	}
	return nil
}

// Mirrors actors/reward/src/streams.rs::validate_award_state_structure.
func validateAwardStateStructure(streams *StreamsState) error {
	if err := validatePendingQueue(streams.PendingWrites); err != nil {
		return err
	}
	if err := validateStreamConfigurationWithoutWeights(streams.Streams); err != nil {
		return err
	}
	liveIDs := make(map[StreamID]struct{}, len(streams.Streams))
	for _, stream := range streams.Streams {
		liveIDs[stream.ID] = struct{}{}
	}
	tombstoneIDs := make(map[StreamID]struct{}, len(streams.Tombstones))
	for i, tombstone := range streams.Tombstones {
		if tombstone.ID == 0 {
			return fmt.Errorf("stream ID 0 is reserved")
		}
		if i > 0 && streams.Tombstones[i-1].ID >= tombstone.ID {
			return fmt.Errorf("tombstones are not ordered")
		}
		if _, found := liveIDs[tombstone.ID]; found {
			return fmt.Errorf("a stream ID is live and tombstoned")
		}
		if len(tombstone.Payable) == 0 {
			return fmt.Errorf("tombstone %d is empty", tombstone.ID)
		}
		if err := validateAmountRows(tombstone.Payable, "tombstone payable"); err != nil {
			return err
		}
		tombstoneIDs[tombstone.ID] = struct{}{}
	}
	for _, write := range streams.PendingWrites {
		if write.Op != PendingWriteOpRegisterStream || write.ID == nil {
			continue
		}
		if _, found := liveIDs[*write.ID]; found {
			return fmt.Errorf("pending registration reuses stream ID %d", *write.ID)
		}
		if _, found := tombstoneIDs[*write.ID]; found {
			return fmt.Errorf("pending registration reuses stream ID %d", *write.ID)
		}
	}
	return validateTombstoneCapacity(streams)
}

func cloneStreamsState(streams *StreamsState) StreamsState {
	clone := StreamsState{
		Streams:       make([]Stream, len(streams.Streams)),
		Tombstones:    append([]Tombstone(nil), streams.Tombstones...),
		PendingWrites: append([]PendingWrite(nil), streams.PendingWrites...),
	}
	for i, stream := range streams.Streams {
		clone.Streams[i] = stream
		if stream.Distribution != nil {
			distribution := *stream.Distribution
			distribution.Shares = append([]RecipientShare(nil), distribution.Shares...)
			distribution.Payable = append([]RecipientAmount(nil), distribution.Payable...)
			distribution.ClaimedPeriod = append([]RecipientAmount(nil), distribution.ClaimedPeriod...)
			clone.Streams[i].Distribution = &distribution
		}
	}
	return clone
}

// Mirrors the projection-only path through
// actors/reward/src/streams.rs::apply_pending_transition.
func applyPendingTransition(streams *StreamsState, write PendingWrite) error {
	switch write.Op {
	case PendingWriteOpSetWeightRecords, PendingWriteOpStepWeightRecords:
		var payload SetWeightRecordsParams
		if err := unmarshalPendingPayload(write.Payload, &payload); err != nil {
			return err
		}
		if err := validateWeightUpdates(payload.Updates); err != nil {
			return err
		}
		for _, update := range payload.Updates {
			found := false
			for i := range streams.Streams {
				if streams.Streams[i].ID == update.ID {
					streams.Streams[i].Weight = update.Weight
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("stream %d not found", update.ID)
			}
		}
	case PendingWriteOpRegisterStream:
		if write.ID == nil {
			return fmt.Errorf("RegisterStream call has no stream ID")
		}
		if *write.ID == 0 {
			return fmt.Errorf("stream ID 0 is reserved")
		}
		for _, stream := range streams.Streams {
			if stream.ID == *write.ID {
				return fmt.Errorf("stream ID %d is already registered", *write.ID)
			}
		}
		for _, tombstone := range streams.Tombstones {
			if tombstone.ID == *write.ID {
				return fmt.Errorf("stream ID %d is tombstoned", *write.ID)
			}
		}
		var payload RegisterStreamPayload
		if err := unmarshalPendingPayload(write.Payload, &payload); err != nil {
			return err
		}
		if err := validateWeightRecord(payload.Weight); err != nil {
			return err
		}
		if err := validateDistributionInit(payload.Distribution); err != nil {
			return err
		}
		var distribution *ExplicitDistribution
		if payload.Distribution != nil {
			distribution = &ExplicitDistribution{
				Writer: payload.Distribution.Writer,
				Shares: append([]RecipientShare(nil), payload.Distribution.Shares...),
			}
		}
		streams.Streams = append(streams.Streams, Stream{ID: *write.ID, Weight: payload.Weight, Distribution: distribution})
		sort.Slice(streams.Streams, func(i, j int) bool { return streams.Streams[i].ID < streams.Streams[j].ID })
	case PendingWriteOpRemoveStream:
		if write.ID == nil {
			return fmt.Errorf("RemoveStream call has no stream ID")
		}
		found := false
		for i, stream := range streams.Streams {
			if stream.ID == *write.ID {
				streams.Streams = append(streams.Streams[:i], streams.Streams[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("stream %d not found", *write.ID)
		}
	case PendingWriteOpSetDistribution:
		if write.ID == nil {
			return fmt.Errorf("SetDistribution call has no stream ID")
		}
		var payload SetDistributionPayload
		if err := unmarshalPendingPayload(write.Payload, &payload); err != nil {
			return err
		}
		for i := range streams.Streams {
			if streams.Streams[i].ID != *write.ID {
				continue
			}
			if streams.Streams[i].Distribution == nil {
				return fmt.Errorf("stream %d is implicit", *write.ID)
			}
			streams.Streams[i].Distribution.Writer = payload.Writer
			return nil
		}
		return fmt.Errorf("stream %d not found", *write.ID)
	default:
		return fmt.Errorf("unknown pending operation %d", write.Op)
	}
	return nil
}

// Mirrors actors/reward/src/streams.rs::validate_projected_queue.
func validateProjectedQueue(streams *StreamsState, currentEpoch abi.ChainEpoch) error {
	if err := validatePendingQueue(streams.PendingWrites); err != nil {
		return err
	}
	if err := validateStreamConfiguration(streams.Streams); err != nil {
		return err
	}
	if err := validateWeightSchedule(streams.Streams, currentEpoch); err != nil {
		return err
	}
	projected := cloneStreamsState(streams)
	for _, write := range streams.PendingWrites {
		candidate := cloneStreamsState(&projected)
		if err := applyPendingTransition(&candidate, write); err != nil {
			continue
		}
		if err := validateStreamConfiguration(candidate.Streams); err != nil {
			continue
		}
		if err := validateWeightSchedule(candidate.Streams, write.EffectiveEpoch); err != nil {
			continue
		}
		if err := validateTombstoneCapacity(&candidate); err != nil {
			continue
		}
		projected = candidate
	}
	return nil
}

// Mirrors actors/reward/src/streams.rs::validate_streams_state.
func validateStreamsState(streams *StreamsState, accruals []StreamAccrual, currentEpoch abi.ChainEpoch) error {
	if err := validateAwardStateStructure(streams); err != nil {
		return err
	}
	for i := 1; i < len(accruals); i++ {
		if accruals[i-1].ID >= accruals[i].ID {
			return fmt.Errorf("explicit-stream accruals are not ordered")
		}
	}
	explicit := make(map[StreamID]*ExplicitDistribution)
	for _, stream := range streams.Streams {
		if stream.Distribution != nil {
			explicit[stream.ID] = stream.Distribution
		}
	}
	if len(explicit) != len(accruals) {
		return fmt.Errorf("explicit-stream accrual IDs do not match live explicit streams")
	}
	for _, accrual := range accruals {
		distribution, found := explicit[accrual.ID]
		if !found {
			return fmt.Errorf("explicit-stream accrual IDs do not match live explicit streams")
		}
		if accrual.Amount.LessThan(big.Zero()) {
			return fmt.Errorf("explicit-stream accrual %d is negative", accrual.ID)
		}
		if err := validatePeriodClaims(distribution, accrual.Amount); err != nil {
			return err
		}
	}
	for _, stream := range streams.Streams {
		if err := validateWeightRecord(stream.Weight); err != nil {
			return err
		}
	}
	return validateProjectedQueue(streams, currentEpoch)
}

// Mirrors actors/reward/src/streams.rs::compute_service_liability.
func computeExplicitLiability(streams *StreamsState, accruals []StreamAccrual) (abi.TokenAmount, error) {
	total := big.Zero()
	accrualIndex := 0
	for _, stream := range streams.Streams {
		if stream.Distribution == nil {
			continue
		}
		if accrualIndex >= len(accruals) {
			return big.Zero(), fmt.Errorf("missing accrual for stream %d", stream.ID)
		}
		accrual := accruals[accrualIndex]
		accrualIndex++
		if accrual.ID != stream.ID {
			return big.Zero(), fmt.Errorf("explicit-stream accrual %d does not match explicit stream %d", accrual.ID, stream.ID)
		}
		if accrual.Amount.LessThan(big.Zero()) {
			return big.Zero(), fmt.Errorf("explicit-stream accrual for stream %d is negative", stream.ID)
		}
		if err := validatePeriodClaims(stream.Distribution, accrual.Amount); err != nil {
			return big.Zero(), err
		}
		claimed := big.Zero()
		for _, row := range stream.Distribution.ClaimedPeriod {
			claimed = big.Add(claimed, row.Amount)
		}
		total = big.Add(total, big.Sub(accrual.Amount, claimed))
		for _, row := range stream.Distribution.Payable {
			total = big.Add(total, row.Amount)
		}
	}
	if accrualIndex < len(accruals) {
		return big.Zero(), fmt.Errorf("explicit-stream accrual %d has no matching explicit stream", accruals[accrualIndex].ID)
	}
	for _, tombstone := range streams.Tombstones {
		if err := validateAmountRows(tombstone.Payable, "tombstone payable"); err != nil {
			return big.Zero(), err
		}
		for _, row := range tombstone.Payable {
			total = big.Add(total, row.Amount)
		}
	}
	return total, nil
}

// ValidateMigrationStreams validates and constructs the two activation streams.
func ValidateMigrationStreams(params []RegisterStreamParams, activationEpoch abi.ChainEpoch) (*StreamsState, []StreamAccrual, error) {
	if len(params) != 2 {
		return nil, nil, fmt.Errorf("bootstrap requires exactly two streams")
	}
	if params[0].ID != 1 || params[1].ID != 2 {
		return nil, nil, fmt.Errorf("bootstrap stream IDs must be 1 and 2")
	}
	if params[0].Distribution != nil || params[1].Distribution == nil {
		return nil, nil, fmt.Errorf("bootstrap distribution forms are invalid")
	}
	for _, param := range params {
		if param.ActivationEpoch != activationEpoch {
			return nil, nil, fmt.Errorf("stream %d activation epoch %d does not match upgrade epoch %d", param.ID, param.ActivationEpoch, activationEpoch)
		}
		if param.Weight.TStart != activationEpoch {
			return nil, nil, fmt.Errorf("stream %d weight start %d does not match upgrade epoch %d", param.ID, param.Weight.TStart, activationEpoch)
		}
	}

	pct := Denom / 100
	consensus := params[0].Weight
	explicit := params[1].Weight
	if consensus.VStart != 95*pct || consensus.Floor != 50*pct || consensus.Cap != 95*pct {
		return nil, nil, fmt.Errorf("consensus bootstrap weight is invalid")
	}
	if explicit.VStart != 5*pct || explicit.Floor != 5*pct || explicit.Cap != 10*pct {
		return nil, nil, fmt.Errorf("explicit bootstrap weight is invalid")
	}
	if consensus.Slope >= 0 || explicit.Slope <= 0 || consensus.Slope != -explicit.Slope {
		return nil, nil, fmt.Errorf("bootstrap weight slopes are invalid")
	}
	if len(params[1].Distribution.Shares) != 1 || params[1].Distribution.Shares[0].Share != Denom {
		return nil, nil, fmt.Errorf("explicit bootstrap requires one full-share recipient")
	}

	streams := &StreamsState{Streams: make([]Stream, 0, len(params))}
	accruals := make([]StreamAccrual, 0, 1)
	for _, param := range params {
		var distribution *ExplicitDistribution
		if param.Distribution != nil {
			distribution = &ExplicitDistribution{
				Writer: param.Distribution.Writer,
				Shares: append([]RecipientShare(nil), param.Distribution.Shares...),
			}
			accruals = append(accruals, StreamAccrual{ID: param.ID, Amount: big.Zero()})
		}
		streams.Streams = append(streams.Streams, Stream{
			ID:           param.ID,
			Weight:       param.Weight,
			Distribution: distribution,
		})
	}
	if err := validateStreamsState(streams, accruals, activationEpoch); err != nil {
		return nil, nil, err
	}
	return streams, accruals, nil
}
