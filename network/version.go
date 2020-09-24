package network

import "math"

// Enumeration of network upgrades where actor behaviour can change (without necessarily
// vendoring and versioning the whole actor codebase).
type Version uint

const (
	Version0 = Version(iota) // specs-actors v0.9.3
	Version1                 // specs-actors v0.9.7
	Version2                 // specs-actors v0.9.8
	Version3                 // specs-actors v0.9.11
	Version4                 // Who knows?

	VersionMax = Version(math.MaxUint32)
)
