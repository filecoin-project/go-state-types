package migration

import (
	"context"
	"sync"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/rt"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/multiformats/go-multibase"

	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

type MemMigrationCache struct {
	MigrationMap sync.Map
}

func NewMemMigrationCache() *MemMigrationCache {
	return new(MemMigrationCache)
}

func (m *MemMigrationCache) Write(key string, c cid.Cid) error {
	m.MigrationMap.Store(key, c)
	return nil
}

func (m *MemMigrationCache) Read(key string) (bool, cid.Cid, error) {
	val, found := m.MigrationMap.Load(key)
	if !found {
		return false, cid.Undef, nil
	}
	c, ok := val.(cid.Cid)
	if !ok {
		return false, cid.Undef, xerrors.Errorf("non cid value in cache")
	}

	return true, c, nil
}

func (m *MemMigrationCache) Load(key string, loadFunc func() (cid.Cid, error)) (cid.Cid, error) {
	found, c, err := m.Read(key)
	if err != nil {
		return cid.Undef, err
	}
	if found {
		return c, nil
	}
	c, err = loadFunc()
	if err != nil {
		return cid.Undef, err
	}
	m.MigrationMap.Store(key, c)
	return c, nil
}

func (m *MemMigrationCache) Clone() *MemMigrationCache {
	newCache := NewMemMigrationCache()
	newCache.Update(m)
	return newCache
}

func (m *MemMigrationCache) Update(other *MemMigrationCache) {
	other.MigrationMap.Range(func(key, value interface{}) bool {
		m.MigrationMap.Store(key, value)
		return true
	})
}

// Config parameterizes a state tree migration
type Config struct {
	// Number of migration worker goroutines to run.
	// More workers enables higher CPU utilization doing migration computations (including state encoding)
	MaxWorkers uint
	// Capacity of the queue of jobs available to workers (zero for unbuffered).
	// A queue length of hundreds to thousands improves throughput at the cost of memory.
	JobQueueSize uint
	// Capacity of the queue receiving migration results from workers, for persisting (zero for unbuffered).
	// A queue length of tens to hundreds improves throughput at the cost of memory.
	ResultQueueSize uint
	// Time between progress logs to emit.
	// Zero (the default) results in no progress logs.
	ProgressLogPeriod time.Duration
}

type Logger interface {
	// This is the same logging interface provided by the Runtime
	Log(level rt.LogLevel, msg string, args ...interface{})
}

// MigrationCache stores and loads cached data. Its implementation must be thread-safe
type MigrationCache interface {
	Write(key string, newCid cid.Cid) error
	Read(key string) (bool, cid.Cid, error)
	Load(key string, loadFunc func() (cid.Cid, error)) (cid.Cid, error)
}

func ActorHeadKey(addr address.Address, head cid.Cid) string {
	headKey, err := head.StringOfBase(multibase.Base32)
	if err != nil {
		panic(err)
	}

	return addr.String() + "-head-" + headKey
}

func SectorsAmtKey(sectorsAmt cid.Cid) string {
	sectorsAmtKey, err := sectorsAmt.StringOfBase(multibase.Base32)
	if err != nil {
		panic(err)
	}

	return "sectorsAmt-" + sectorsAmtKey
}

func PartitionsAmtKey(partitionsAmt cid.Cid) string {
	partitionsAmtKey, err := partitionsAmt.StringOfBase(multibase.Base32)
	if err != nil {
		panic(err)
	}

	return "partitionsAmt-" + partitionsAmtKey
}

func ExpirationsAmtKey(expirationsAmt cid.Cid) string {
	expirationsAmtKey, err := expirationsAmt.StringOfBase(multibase.Base32)
	if err != nil {
		panic(err)
	}

	return "partitionsAmt-" + expirationsAmtKey
}

func MinerPrevSectorsInKey(m address.Address) string {
	return "prevSectorsIn-" + m.String()
}

func MinerPrevSectorsOutKey(m address.Address) string {
	return "prevSectorsOut-" + m.String()
}

func MarketPrevDealStatesInKey(m address.Address) string {
	return "prevDealStatesIn-" + m.String()
}

func MarketPrevDealProposalsInKey(m address.Address) string {
	return "prevDealProposalsIn-" + m.String()
}

func MarketPrevDealStatesOutKey(m address.Address) string {
	return "prevDealStatesOut-" + m.String()
}

func MarketPrevProviderSectorsOutKey(m address.Address) string {
	return "prevProviderSectorsOut-" + m.String()
}

type ActorMigrationInput struct {
	Address address.Address // actor's address
	Head    cid.Cid
	Cache   MigrationCache // cache of existing cid -> cid migrations for this actor
}

type ActorMigrationResult struct {
	NewCodeCID cid.Cid
	NewHead    cid.Cid
}

type ActorMigration interface {
	// MigrateState Loads an actor's state from an input store and writes new state to an output store.
	// Returns the new state head CID.
	MigrateState(ctx context.Context, store cbor.IpldStore, input ActorMigrationInput) (result *ActorMigrationResult, err error)
	MigratedCodeCID() cid.Cid

	// Deferred returns true if this migration should be run after all non-deferred migrations have completed.
	Deferred() bool
}

// Migrator which preserves the head CID and provides a fixed result code CID.
type CodeMigrator struct {
	OutCodeCID cid.Cid
}

func (n CodeMigrator) MigrateState(_ context.Context, _ cbor.IpldStore, in ActorMigrationInput) (*ActorMigrationResult, error) {
	return &ActorMigrationResult{
		NewCodeCID: n.OutCodeCID,
		NewHead:    in.Head,
	}, nil
}

func (n CodeMigrator) MigratedCodeCID() cid.Cid {
	return n.OutCodeCID
}

func (n CodeMigrator) Deferred() bool {
	return false
}

// Migrator that uses cached transformation if it exists
type CachedMigrator struct {
	cache MigrationCache
	ActorMigration
}

func (c CachedMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in ActorMigrationInput) (*ActorMigrationResult, error) {
	newHead, err := c.cache.Load(ActorHeadKey(in.Address, in.Head), func() (cid.Cid, error) {
		result, err := c.ActorMigration.MigrateState(ctx, store, in)
		if err != nil {
			return cid.Undef, xerrors.Errorf("migrating state: %w", err)
		}
		return result.NewHead, nil
	})
	if err != nil {
		return nil, xerrors.Errorf("using cache: %w", err)
	}
	return &ActorMigrationResult{
		NewCodeCID: c.MigratedCodeCID(),
		NewHead:    newHead,
	}, nil
}

func (c CachedMigrator) Deferred() bool {
	return c.ActorMigration.Deferred()
}

func CachedMigration(cache MigrationCache, m ActorMigration) ActorMigration {
	return CachedMigrator{
		ActorMigration: m,
		cache:          cache,
	}
}
