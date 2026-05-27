package verifreg

import (
	"github.com/filecoin-project/go-state-types/builtin"
)

const EndOfLifeClaimDropPeriod = 30 * builtin.EpochsInDay

const MaximumVerifiedAllocationExpiration = 60 * builtin.EpochsInDay

const MinimumVerifiedAllocationTerm = 180 * builtin.EpochsInDay

const MaximumVerifiedAllocationTerm = 5 * builtin.EpochsInYear

const NoAllocationID = AllocationId(0)

const MinimumVerifiedAllocationSize = 1 << 20
