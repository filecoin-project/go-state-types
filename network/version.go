package network

import "math"

// Enumeration of network upgrades where actor behaviour can change (without necessarily
// vendoring and versioning the whole actor codebase).
type Version uint

const (
	Version0 = Version(iota) // genesis   (specs-actors v0.9.3)
	Version1                 // breeze    (specs-actors v0.9.7)
	Version2                 // smoke     (specs-actors v0.9.8)
	Version3                 // ignition  (specs-actors v0.9.11)
	Version4                 // actors v2 (specs-actors v2.0.x (future))
	Version5                 // sometime?

	// VersionMax is the maximum version number
	VersionMax = Version(math.MaxUint32)
)
