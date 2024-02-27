package test_util

import (
	"context"
	"fmt"
	"sync"

	block "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
)

type BlockStoreInMemory struct {
	data map[cid.Cid]block.Block
}

func NewBlockStoreInMemory() *BlockStoreInMemory {
	return &BlockStoreInMemory{make(map[cid.Cid]block.Block)}
}

func (mb *BlockStoreInMemory) Get(ctx context.Context, c cid.Cid) (block.Block, error) {
	d, ok := mb.data[c]
	if ok {
		return d, nil
	}
	return nil, fmt.Errorf("not found")
}

func (mb *BlockStoreInMemory) Put(ctx context.Context, b block.Block) error {
	mb.data[b.Cid()] = b
	return nil
}

type SyncBlockStoreInMemory struct {
	bs *BlockStoreInMemory
	mu sync.Mutex
}

func (ss *SyncBlockStoreInMemory) Context() context.Context {
	return context.Background()
}

func NewSyncBlockStoreInMemory() *SyncBlockStoreInMemory {
	return &SyncBlockStoreInMemory{
		bs: NewBlockStoreInMemory(),
	}
}

func (ss *SyncBlockStoreInMemory) Get(ctx context.Context, c cid.Cid) (block.Block, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.bs.Get(ctx, c)
}

func (ss *SyncBlockStoreInMemory) Put(ctx context.Context, b block.Block) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.bs.Put(ctx, b)
}
