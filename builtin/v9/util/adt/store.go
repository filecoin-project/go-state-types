package adt

import (
	"context"

	ipldcbor "github.com/ipfs/go-ipld-cbor"
)

// Store defines an interface required to back the ADTs in this package.
type Store interface {
	Context() context.Context
	ipldcbor.IpldStore
}
