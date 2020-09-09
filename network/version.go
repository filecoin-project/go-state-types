package network

import "math"

// Enumeration of network upgrades where actor behaviour can change (without necessarily
// vendoring and versioning the whole actor codebase).
type Version uint

const (
	Version0 = Version(iota) // specs-actors v0.9.3
	Version1                 // specs-actors v0.9.7
	Version2                 // specs-actors v0.9.?
	Version3                 // Coming soon
	Version4                 // Who knows?

	VersionMax = Version(math.MaxUint32)
)
