package adt

import (
	"context"

	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"

	ipldcbor "github.com/ipfs/go-ipld-cbor"
)

type Store = adt9.Store

func WrapStore(ctx context.Context, store ipldcbor.IpldStore) Store {
	return adt9.WrapStore(ctx, store)
}
