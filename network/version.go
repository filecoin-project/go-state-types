package network

import "math"

// Enumeration of network upgrades where actor behaviour can change (without necessarily
// vendoring and versioning the whole actor codebase).
type Version uint

const (
	Version0         = Version(0)    // genesis    (specs-actors v0.9.3)
	Version1         = Version(100)  // breeze     (specs-actors v0.9.7)
	Version2         = Version(200)  // smoke      (specs-actors v0.9.8)
	Version3         = Version(300)  // ignition   (specs-actors v0.9.11)
	Version4         = Version(400)  // actors v2  (specs-actors v2.0.3)
	Version5         = Version(500)  // tape       (specs-actors v2.1.0)
	Version6         = Version(600)  // kumquat    (specs-actors v2.2.0)
	Version6AndAHalf = Version(650)  // pre-calico (specs-actors v2.2.0)
	Version7         = Version(700)  // calico     (specs-actors v2.3.2)
	Version8         = Version(800)  // persian    (post-2.3.2 behaviour transition)
	Version9         = Version(900)  // orange     (post-2.3.2 behaviour transition)
	Version10        = Version(1000) // trust      (specs-actors v3.0.1)
	Version11        = Version(1100) // norwegian  (specs-actors v3.1.0)
	Version12        = Version(1200) // turbo      (specs-actors v4.0.0)
	Version13        = Version(1300) // hyperdrive (specs-actors v5.0.1)

	// VersionMax is the maximum version number
	VersionMax = Version(math.MaxUint32)
)
