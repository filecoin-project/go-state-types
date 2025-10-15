package actors

import (
	"fmt"

	"github.com/filecoin-project/go-state-types/network"
)

type Version int

const (
	Version0  Version = 0
	Version2  Version = 2
	Version3  Version = 3
	Version4  Version = 4
	Version5  Version = 5
	Version6  Version = 6
	Version7  Version = 7
	Version8  Version = 8
	Version9  Version = 9
	Version10 Version = 10
	Version11 Version = 11
	Version12 Version = 12
	Version13 Version = 13
	Version14 Version = 14
	Version15 Version = 15
	Version16 Version = 16
	Version17 Version = 17
	Version18 Version = 18
)

// Converts a network version into an actors adt version.
func VersionForNetwork(version network.Version) (Version, error) {
	switch version {
	case network.Version0, network.Version1, network.Version2, network.Version3:
		return Version0, nil
	case network.Version4, network.Version5, network.Version6, network.Version7, network.Version8, network.Version9:
		return Version2, nil
	case network.Version10, network.Version11:
		return Version3, nil
	case network.Version12:
		return Version4, nil
	case network.Version13:
		return Version5, nil
	case network.Version14:
		return Version6, nil
	case network.Version15:
		return Version7, nil
	case network.Version16:
		return Version8, nil
	case network.Version17:
		return Version9, nil
	case network.Version18:
		return Version10, nil
	case network.Version19, network.Version20:
		return Version11, nil
	case network.Version21:
		return Version12, nil
	case network.Version22:
		return Version13, nil
	case network.Version23:
		return Version14, nil
	case network.Version24:
		return Version15, nil
	case network.Version25, network.Version26:
		return Version16, nil
	case network.Version27:
		return Version17, nil
	case network.Version28:
		return Version18, nil
	default:
		return -1, fmt.Errorf("unsupported network version %d", version)
	}
}
