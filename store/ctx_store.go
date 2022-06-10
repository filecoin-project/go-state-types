package store

import (
	"context"

	ipldcbor "github.com/ipfs/go-ipld-cbor"
)

type Store interface {
	Context() context.Context
	ipldcbor.IpldStore
}

// WrapStore Adapts a vanilla IPLD store as an ADT store.
func WrapStore(ctx context.Context, store ipldcbor.IpldStore) Store {
	return &wstore{
		ctx:       ctx,
		IpldStore: store,
	}
}

// WrapBlockStore Adapts a block store as an ADT store.
func WrapBlockStore(ctx context.Context, bs ipldcbor.IpldBlockstore) Store {
	return WrapStore(ctx, ipldcbor.NewCborStore(bs))
}

type wstore struct {
	ctx context.Context
	ipldcbor.IpldStore
}

var _ Store = &wstore{}

func (s *wstore) Context() context.Context {
	return s.ctx
}
