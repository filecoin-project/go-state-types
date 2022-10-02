package migration

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	market9 "github.com/filecoin-project/go-state-types/builtin/v9/market"

	init9 "github.com/filecoin-project/go-state-types/builtin/v9/init"

	datacap9 "github.com/filecoin-project/go-state-types/builtin/v9/datacap"

	verifreg9 "github.com/filecoin-project/go-state-types/builtin/v9/verifreg"

	"github.com/filecoin-project/go-state-types/big"
	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"

	verifreg8 "github.com/filecoin-project/go-state-types/builtin/v8/verifreg"

	market8 "github.com/filecoin-project/go-state-types/builtin/v8/market"

	system8 "github.com/filecoin-project/go-state-types/builtin/v8/system"

	"github.com/filecoin-project/go-state-types/builtin"

	"github.com/multiformats/go-multibase"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	adt8 "github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/filecoin-project/go-state-types/rt"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

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

// MigrationCache stores and loads cached data. Its implementation must be threadsafe
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

func MinerPrevSectorsInKey(m address.Address) string {
	return "prevSectorsIn-" + m.String()
}

func MinerPrevSectorsOutKey(m address.Address) string {
	return "prevSectorsOut-" + m.String()
}

// Migrates the filecoin state tree starting from the global state tree and upgrading all actor state.
// The store must support concurrent writes (even if the configured worker count is 1).
func MigrateStateTree(ctx context.Context, store cbor.IpldStore, newManifestCID cid.Cid, actorsRootIn cid.Cid, priorEpoch abi.ChainEpoch, cfg Config, log Logger, cache MigrationCache) (cid.Cid, error) {
	if cfg.MaxWorkers <= 0 {
		return cid.Undef, xerrors.Errorf("invalid migration config with %d workers", cfg.MaxWorkers)
	}

	adtStore := adt8.WrapStore(ctx, store)

	// Load input and output state trees
	actorsIn, err := LoadTree(adtStore, actorsRootIn)
	if err != nil {
		return cid.Undef, err
	}
	actorsOut, err := NewTree(adtStore)
	if err != nil {
		return cid.Undef, err
	}

	// load old manifest data
	systemActor, ok, err := actorsIn.GetActor(builtin.SystemActorAddr)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to get system actor: %w", err)
	}

	if !ok {
		return cid.Undef, xerrors.New("didn't find system actor")
	}

	var systemState system8.State
	if err := store.Get(ctx, systemActor.Head, &systemState); err != nil {
		return cid.Undef, xerrors.Errorf("failed to get system actor state: %w", err)
	}

	var oldManifestData manifest.ManifestData
	if err := store.Get(ctx, systemState.BuiltinActors, &oldManifestData); err != nil {
		return cid.Undef, xerrors.Errorf("failed to get old manifest data: %w", err)
	}

	// load new manifest
	var newManifest manifest.Manifest
	if err := adtStore.Get(ctx, newManifestCID, &newManifest); err != nil {
		return cid.Undef, xerrors.Errorf("error reading actor manifest: %w", err)
	}

	if err := newManifest.Load(ctx, adtStore); err != nil {
		return cid.Undef, xerrors.Errorf("error loading actor manifest: %w", err)
	}

	// Maps prior version code CIDs to migration functions.
	migrations := make(map[cid.Cid]actorMigration)
	// Set of prior version code CIDs for actors to defer during iteration, for explicit migration afterwards.
	deferredCodeIDs := make(map[cid.Cid]struct{})

	// Populated from oldManifestData
	oldCodeIDMap := make(map[string]cid.Cid, len(oldManifestData.Entries))

	miner8Cid := cid.Undef
	for _, entry := range oldManifestData.Entries {
		oldCodeIDMap[entry.Name] = entry.Code
		if entry.Name == manifest.MinerKey {
			miner8Cid = entry.Code
		}
	}

	if miner8Cid == cid.Undef {
		return cid.Undef, xerrors.Errorf("didn't find miner in old manifest entries")
	}

	for name, oldCodeCID := range oldCodeIDMap { //nolint:nomaprange
		if name == manifest.MarketKey || name == manifest.VerifregKey {
			deferredCodeIDs[oldCodeCID] = struct{}{}
		} else {
			newCodeCID, ok := newManifest.Get(name)
			if !ok {
				return cid.Undef, xerrors.Errorf("code cid for %s actor not found in new manifest", name)
			}

			migrations[oldCodeCID] = codeMigrator{newCodeCID}
		}
	}

	// migrations that migrate both code and state, override entries in `migrations`

	// The System Actor

	newSystemCodeCID, ok := newManifest.Get("system")
	if !ok {
		return cid.Undef, xerrors.Errorf("code cid for system actor not found in manifest")
	}

	migrations[systemActor.Code] = systemActorMigrator{newSystemCodeCID, newManifest.Data}

	// The Miner Actor

	// load market proposals
	marketActorV8, ok, err := actorsIn.GetActor(builtin.StorageMarketActorAddr)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to get market actor: %w", err)
	}

	if !ok {
		return cid.Undef, xerrors.New("didn't find market actor")
	}

	var marketStateV8 market8.State
	if err := store.Get(ctx, marketActorV8.Head, &marketStateV8); err != nil {
		return cid.Undef, xerrors.Errorf("failed to get market actor state: %w", err)
	}

	proposals, err := market8.AsDealProposalArray(adtStore, marketStateV8.Proposals)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to get proposals: %w", err)
	}

	miner9Cid, ok := newManifest.Get(manifest.MinerKey)
	if !ok {
		return cid.Undef, xerrors.Errorf("code cid for miner actor not found in new manifest")
	}

	mm, err := newMinerMigrator(ctx, store, proposals, miner9Cid)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to create miner migrator: %w", err)
	}

	migrations[miner8Cid] = cachedMigration(cache, *mm)

	if len(migrations)+len(deferredCodeIDs) != len(oldManifestData.Entries) {
		return cid.Undef, xerrors.Errorf("incomplete migration specification with %d code CIDs, need %d", len(migrations), len(oldManifestData.Entries))
	}
	startTime := time.Now()

	// Setup synchronization
	grp, ctx := errgroup.WithContext(ctx)
	// Input and output queues for workers.
	jobCh := make(chan *migrationJob, cfg.JobQueueSize)
	jobResultCh := make(chan *migrationJobResult, cfg.ResultQueueSize)
	// Atomically-modified counters for logging progress
	var jobCount uint32
	var doneCount uint32

	// Iterate all actors in old state root to create migration jobs for each non-deferred actor.
	grp.Go(func() error {
		defer close(jobCh)
		log.Log(rt.INFO, "Creating migration jobs for tree %s", actorsRootIn)
		if err = actorsIn.ForEach(func(addr address.Address, actorIn *Actor) error {
			if _, ok := deferredCodeIDs[actorIn.Code]; ok {
				return nil
			}

			migration, ok := migrations[actorIn.Code]
			if !ok {
				return xerrors.Errorf("actor with code %s has no registered migration function", actorIn.Code)
			}

			nextInput := &migrationJob{
				Address:        addr,
				Actor:          *actorIn, // Must take a copy, the pointer is not stable.
				cache:          cache,
				actorMigration: migration,
			}

			select {
			case jobCh <- nextInput:
			case <-ctx.Done():
				return ctx.Err()
			}
			atomic.AddUint32(&jobCount, 1)
			return nil
		}); err != nil {
			return err
		}
		log.Log(rt.INFO, "Done creating %d migration jobs for tree %s after %v", jobCount, actorsRootIn, time.Since(startTime))
		return nil
	})

	// Worker threads run jobs.
	var workerWg sync.WaitGroup
	for i := uint(0); i < cfg.MaxWorkers; i++ {
		workerWg.Add(1)
		workerId := i
		grp.Go(func() error {
			defer workerWg.Done()
			for job := range jobCh {
				result, err := job.run(ctx, store, priorEpoch)
				if err != nil {
					return err
				}
				select {
				case jobResultCh <- result:
				case <-ctx.Done():
					return ctx.Err()
				}
				atomic.AddUint32(&doneCount, 1)
			}
			log.Log(rt.INFO, "Worker %d done", workerId)
			return nil
		})
	}
	log.Log(rt.INFO, "Started %d workers", cfg.MaxWorkers)

	// Monitor the job queue. This non-critical goroutine is outside the errgroup and exits when
	// workersFinished is closed, or the context done.
	workersFinished := make(chan struct{}) // Closed when waitgroup is emptied.
	if cfg.ProgressLogPeriod > 0 {
		go func() {
			defer log.Log(rt.DEBUG, "Job queue monitor done")
			for {
				select {
				case <-time.After(cfg.ProgressLogPeriod):
					jobsNow := jobCount // Snapshot values to avoid incorrect-looking arithmetic if they change.
					doneNow := doneCount
					pendingNow := jobsNow - doneNow
					elapsed := time.Since(startTime)
					rate := float64(doneNow) / elapsed.Seconds()
					log.Log(rt.INFO, "%d jobs created, %d done, %d pending after %v (%.0f/s)",
						jobsNow, doneNow, pendingNow, elapsed, rate)
				case <-workersFinished:
					return
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// Close result channel when workers are done sending to it.
	grp.Go(func() error {
		workerWg.Wait()
		close(jobResultCh)
		close(workersFinished)
		log.Log(rt.INFO, "All workers done after %v", time.Since(startTime))
		return nil
	})

	// Insert migrated records in output state tree and accumulators.
	grp.Go(func() error {
		log.Log(rt.INFO, "Result writer started")
		resultCount := 0
		for result := range jobResultCh {
			if err := actorsOut.SetActor(result.Address, &result.Actor); err != nil {
				return err
			}
			resultCount++
		}
		log.Log(rt.INFO, "Result writer wrote %d results to state tree after %v", resultCount, time.Since(startTime))
		return nil
	})

	if err := grp.Wait(); err != nil {
		return cid.Undef, err
	}

	elapsed := time.Since(startTime)
	rate := float64(doneCount) / elapsed.Seconds()
	log.Log(rt.INFO, "All %d done after %v (%.0f/s). Starting deferred migrations.", doneCount, elapsed, rate)

	// Create the Datacap actor

	log.Log(rt.INFO, "Migrating datacap actor")

	verifregActorV8, ok, err := actorsIn.GetActor(builtin.VerifiedRegistryActorAddr)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to get verifreg actor: %w", err)
	}

	if !ok {
		return cid.Undef, xerrors.New("didn't find verifreg actor")
	}

	var verifregStateV8 verifreg8.State
	if err = store.Get(ctx, verifregActorV8.Head, &verifregStateV8); err != nil {
		return cid.Undef, xerrors.Errorf("failed to get verifreg actor state: %w", err)
	}

	verifiedClients, err := adt8.AsMap(adtStore, verifregStateV8.VerifiedClients, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to get verified clients: %w", err)
	}

	tokenSupply := big.Zero()

	emptyMapCid, err := adt9.StoreEmptyMap(adtStore, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to create empty map: %w", err)
	}

	balancesMap, err := adt9.AsMap(adtStore, emptyMapCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to load empty map: %w", err)
	}

	allowancesMap, err := adt9.AsMap(adtStore, emptyMapCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to load empty map: %w", err)
	}

	var dcap abi.StoragePower
	if err = verifiedClients.ForEach(&dcap, func(key string) error {
		a, err := address.NewFromBytes([]byte(key))
		if err != nil {
			return err
		}

		tokenAmount := verifreg9.DataCapToTokens(dcap)
		tokenSupply = big.Add(tokenSupply, tokenAmount)
		if err = balancesMap.Put(abi.IdAddrKey(a), &tokenAmount); err != nil {
			return xerrors.Errorf("failed to put new balancesMap entry: %w", err)
		}

		allowancesMapEntry, err := adt9.AsMap(adtStore, emptyMapCid, builtin.DefaultHamtBitwidth)
		if err != nil {
			return xerrors.Errorf("failed to load empty map: %w", err)
		}

		if err = allowancesMapEntry.Put(abi.IdAddrKey(builtin.StorageMarketActorAddr), &datacap9.InfiniteAllowance); err != nil {
			return xerrors.Errorf("failed to populate allowance map: %w", err)
		}

		return allowancesMap.Put(abi.IdAddrKey(a), allowancesMapEntry)
	}); err != nil {
		return cid.Undef, xerrors.Errorf("failed to loop over verified clients: %w", err)
	}

	balancesMapRoot, err := balancesMap.Root()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to flush balances map: %w", err)
	}

	allowancesMapRoot, err := allowancesMap.Root()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to flush allowances map: %w", err)
	}

	dataCapState := datacap9.State{
		Governor: builtin.VerifiedRegistryActorAddr,
		Token: datacap9.TokenState{
			Supply:       tokenSupply,
			Balances:     balancesMapRoot,
			Allowances:   allowancesMapRoot,
			HamtBitWidth: builtin.DefaultHamtBitwidth,
		},
	}

	dataCapHead, err := adtStore.Put(ctx, &dataCapState)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to put data cap state: %w", err)
	}

	dataCapCode, ok := newManifest.Get(manifest.DataCapKey)
	if !ok {
		return cid.Undef, xerrors.Errorf("failed to find datacap code ID: %w", err)
	}

	if err = actorsOut.SetActor(builtin.DatacapActorAddr, &Actor{
		Code:       dataCapCode,
		Head:       dataCapHead,
		CallSeqNum: 0,
		Balance:    big.Zero(),
	}); err != nil {
		return cid.Undef, xerrors.Errorf("failed to set datacap actor: %w", err)
	}

	// Migrate the Verified Registry Actor

	log.Log(rt.INFO, "Migrating the verified registry actor")

	initActorV9, ok, err := actorsOut.GetActor(builtin.InitActorAddr)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to load init actor: %w", err)
	}

	if !ok {
		return cid.Undef, xerrors.New("failed to find init actor")
	}

	var initStateV9 init9.State
	if err = adtStore.Get(ctx, initActorV9.Head, &initStateV9); err != nil {
		return cid.Undef, xerrors.Errorf("failed to load init state: %w", err)
	}

	pendingProposals, err := adt8.AsSet(adtStore, marketStateV8.PendingProposals, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to load pending proposals: %w", err)
	}

	pendingMap := make(map[abi.DealID]market8.DealProposal)
	var proposal market8.DealProposal
	if err = proposals.ForEach(&proposal, func(dealID int64) error {
		pcid, err := proposal.Cid()
		if err != nil {
			return err
		}

		isPending, err := pendingProposals.Has(abi.CidKey(pcid))
		if err != nil {
			return xerrors.Errorf("failed to check pending: %w", err)
		}

		if isPending {
			pendingMap[abi.DealID(dealID)] = proposal
		}

		return nil
	}); err != nil {
		return cid.Undef, xerrors.Errorf("failed to iterate over proposals: %w", err)
	}

	nextAllocationID := verifreg9.AllocationId(1)
	allocationsMapMap := make(map[address.Address]*adt9.Map)
	dealsToAllocations := make(map[abi.DealID]verifreg9.AllocationId)
	for dealID, proposal := range pendingMap {
		clientIDAddress, ok, err := initStateV9.ResolveAddress(adtStore, proposal.Client)
		if err != nil {
			return cid.Undef, xerrors.Errorf("failed to resolve client %s: %w", proposal.Client, err)
		}

		if !ok {
			return cid.Undef, xerrors.New("failed to find client in init actor map")
		}

		clientIDu64, err := address.IDFromAddress(clientIDAddress)
		if err != nil {
			return cid.Undef, err
		}

		providerIDAddress, ok, err := initStateV9.ResolveAddress(adtStore, proposal.Provider)
		if err != nil {
			return cid.Undef, xerrors.Errorf("failed to resolve provider %s: %w", proposal.Provider, err)
		}

		if !ok {
			return cid.Undef, xerrors.New("failed to find provider in init actor map")
		}

		providerIDu64, err := address.IDFromAddress(providerIDAddress)
		if err != nil {
			return cid.Undef, err
		}

		clientAllocationMap, ok := allocationsMapMap[clientIDAddress]
		if !ok {
			clientAllocationMap, err = adt9.AsMap(adtStore, emptyMapCid, builtin.DefaultHamtBitwidth)
			if err != nil {
				return cid.Undef, xerrors.Errorf("failed to load empty map: %w", err)
			}
		}

		if err = clientAllocationMap.Put(nextAllocationID, &verifreg9.Allocation{
			Client:   abi.ActorID(clientIDu64),
			Provider: abi.ActorID(providerIDu64),
			Data:     proposal.PieceCID,
			Size:     proposal.PieceSize,
			TermMin:  proposal.Duration(),
			TermMax:  market9.DealMaxDuration,
			// TODO: priorEpoch + 1???
			Expiration: verifreg9.MaximumVerifiedAllocationExpiration + priorEpoch,
		}); err != nil {
			return cid.Undef, xerrors.Errorf("failed to put new allocation obj: %w", err)
		}

		dealsToAllocations[dealID] = nextAllocationID
		nextAllocationID++
	}

	allocationsMap, err := adt9.AsMap(adtStore, emptyMapCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to load empty map: %w", err)
	}

	for clientID, clientAllocationsMap := range allocationsMapMap {
		if err = allocationsMap.Put(abi.IdAddrKey(clientID), clientAllocationsMap); err != nil {
			return cid.Undef, xerrors.Errorf("failed to populate allocationsMap: %w", err)
		}
	}

	allocationsMapRoot, err := allocationsMap.Root()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to flush allocations map: %w", err)
	}

	verifregStateV9 := verifreg9.State{
		RootKey:                  verifregStateV8.RootKey,
		Verifiers:                verifregStateV8.Verifiers,
		RemoveDataCapProposalIDs: verifregStateV8.RemoveDataCapProposalIDs,
		Allocations:              allocationsMapRoot,
		NextAllocationId:         1,
		Claims:                   emptyMapCid,
	}

	verifregHead, err := adtStore.Put(ctx, &verifregStateV9)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to put verifreg9 state: %w", err)
	}

	verifregCode, ok := newManifest.Get(manifest.VerifregKey)
	if !ok {
		return cid.Undef, xerrors.Errorf("failed to find verifreg code ID: %w", err)
	}

	if err = actorsOut.SetActor(builtin.VerifiedRegistryActorAddr, &Actor{
		Code:       verifregCode,
		Head:       verifregHead,
		CallSeqNum: verifregActorV8.CallSeqNum,
		Balance:    verifregActorV8.Balance,
	}); err != nil {
		return cid.Undef, xerrors.Errorf("failed to set verifreg actor: %w", err)
	}

	// Migrate the Market Actor

	log.Log(rt.INFO, "Migrating the market actor")

	pendingDealAllocationIdsMap, err := adt9.AsMap(adtStore, emptyMapCid, builtin.DefaultHamtBitwidth)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to load empty map: %w", err)
	}

	for dealID, allocationID := range dealsToAllocations {
		if err = pendingDealAllocationIdsMap.Put(abi.UIntKey(uint64(dealID)), allocationID); err != nil {
			return cid.Undef, xerrors.Errorf("failed to populate pending deal allocations map: %w", err)
		}
	}

	pendingDealAllocationIdsMapRoot, err := pendingDealAllocationIdsMap.Root()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to flush pending deal allocations map: %w", err)
	}

	marketStateV9 := market9.State{
		Proposals:                     marketStateV8.Proposals,
		States:                        marketStateV8.States,
		PendingProposals:              marketStateV8.PendingProposals,
		EscrowTable:                   marketStateV8.EscrowTable,
		LockedTable:                   marketStateV8.LockedTable,
		NextID:                        marketStateV8.NextID,
		DealOpsByEpoch:                marketStateV8.DealOpsByEpoch,
		LastCron:                      marketStateV8.LastCron,
		TotalClientLockedCollateral:   marketStateV8.TotalClientLockedCollateral,
		TotalProviderLockedCollateral: marketStateV8.TotalProviderLockedCollateral,
		TotalClientStorageFee:         marketStateV8.TotalClientStorageFee,
		PendingDealAllocationIds:      pendingDealAllocationIdsMapRoot,
	}

	marketHead, err := adtStore.Put(ctx, &marketStateV9)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to put market state: %w", err)
	}

	marketCode, ok := newManifest.Get(manifest.MarketKey)
	if !ok {
		return cid.Undef, xerrors.Errorf("failed to find market code ID: %w", err)
	}

	if err = actorsOut.SetActor(builtin.StorageMarketActorAddr, &Actor{
		Code:       marketCode,
		Head:       marketHead,
		CallSeqNum: marketActorV8.CallSeqNum,
		Balance:    marketActorV8.Balance,
	}); err != nil {
		return cid.Undef, xerrors.Errorf("failed to set market actor: %w", err)
	}

	log.Log(rt.INFO, "Done all migrations, flushing state root")

	return actorsOut.Flush()
}

type actorMigrationInput struct {
	address    address.Address // actor's address
	head       cid.Cid
	priorEpoch abi.ChainEpoch // epoch of last state transition prior to migration
	cache      MigrationCache // cache of existing cid -> cid migrations for this actor
}

type actorMigrationResult struct {
	newCodeCID cid.Cid
	newHead    cid.Cid
}

type actorMigration interface {
	// Loads an actor's state from an input store and writes new state to an output store.
	// Returns the new state head CID.
	migrateState(ctx context.Context, store cbor.IpldStore, input actorMigrationInput) (result *actorMigrationResult, err error)
	migratedCodeCID() cid.Cid
}

type migrationJob struct {
	address.Address
	Actor
	actorMigration
	cache MigrationCache
}

type migrationJobResult struct {
	address.Address
	Actor
}

func (job *migrationJob) run(ctx context.Context, store cbor.IpldStore, priorEpoch abi.ChainEpoch) (*migrationJobResult, error) {
	result, err := job.migrateState(ctx, store, actorMigrationInput{
		address:    job.Address,
		head:       job.Actor.Head,
		priorEpoch: priorEpoch,
		cache:      job.cache,
	})
	if err != nil {
		return nil, xerrors.Errorf("state migration failed for actor code %s, addr %s: %w",
			job.Actor.Code, job.Address, err)
	}

	// Set up new actor record with the migrated state.
	return &migrationJobResult{
		job.Address, // Unchanged
		Actor{
			Code:       result.newCodeCID,
			Head:       result.newHead,
			CallSeqNum: job.Actor.CallSeqNum, // Unchanged
			Balance:    job.Actor.Balance,    // Unchanged
		},
	}, nil
}

// Migrator which preserves the head CID and provides a fixed result code CID.
type codeMigrator struct {
	OutCodeCID cid.Cid
}

func (n codeMigrator) migrateState(_ context.Context, _ cbor.IpldStore, in actorMigrationInput) (*actorMigrationResult, error) {
	return &actorMigrationResult{
		newCodeCID: n.OutCodeCID,
		newHead:    in.head,
	}, nil
}

func (n codeMigrator) migratedCodeCID() cid.Cid {
	return n.OutCodeCID
}

// Migrator that uses cached transformation if it exists
type cachedMigrator struct {
	cache MigrationCache
	actorMigration
}

func (c cachedMigrator) migrateState(ctx context.Context, store cbor.IpldStore, in actorMigrationInput) (*actorMigrationResult, error) {
	newHead, err := c.cache.Load(ActorHeadKey(in.address, in.head), func() (cid.Cid, error) {
		result, err := c.actorMigration.migrateState(ctx, store, in)
		if err != nil {
			return cid.Undef, err
		}
		return result.newHead, nil
	})
	if err != nil {
		return nil, err
	}
	return &actorMigrationResult{
		newCodeCID: c.migratedCodeCID(),
		newHead:    newHead,
	}, nil
}

func cachedMigration(cache MigrationCache, m actorMigration) actorMigration {
	return cachedMigrator{
		actorMigration: m,
		cache:          cache,
	}
}
