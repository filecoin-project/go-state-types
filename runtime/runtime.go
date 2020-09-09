package runtime

// Enumeration of network upgrades where actor behaviour can change (without necessarily
// vendoring and versioning the whole actor codebase).
type NetworkVersion uint // FIXME move to types

const (
	NetworkVersion0 = NetworkVersion(iota) // specs-actors v0.9.3
	NetworkVersion1                        // specs-actors v0.9.7
	NetworkVersion2                        // specs-actors v2.0.?

	NetworkVersionLatest = NetworkVersion2
)
