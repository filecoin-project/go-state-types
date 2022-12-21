package migration

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	init10 "github.com/filecoin-project/go-state-types/builtin/v10/init"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	system9 "github.com/filecoin-project/go-state-types/builtin/v9/system"
	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/filecoin-project/go-state-types/rt"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

// MigrateStateTree Migrates the filecoin state tree starting from the global state tree and upgrading all actor state.
// The store must support concurrent writes (even if the configured worker count is 1).
func MigrateStateTree(ctx context.Context, store cbor.IpldStore, newManifestCID cid.Cid, actorsRootIn cid.Cid, priorEpoch abi.ChainEpoch, cfg migration.Config, log migration.Logger, cache migration.MigrationCache) (cid.Cid, error) {
	if cfg.MaxWorkers <= 0 {
		return cid.Undef, xerrors.Errorf("invalid migration config with %d workers", cfg.MaxWorkers)
	}

	adtStore := adt9.WrapStore(ctx, store)

	// Load input and output state trees
	actorsIn, err := builtin.LoadTree(adtStore, actorsRootIn)
	if err != nil {
		return cid.Undef, xerrors.Errorf("loading state tree: %w", err)
	}
	actorsOut, err := builtin.NewTree(adtStore)
	if err != nil {
		return cid.Undef, xerrors.Errorf("creating new state tree: %w", err)
	}

	// load old manifest data
	systemActor, ok, err := actorsIn.GetActorV4(builtin.SystemActorAddr)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to get system actor: %w", err)
	}

	if !ok {
		return cid.Undef, xerrors.New("didn't find system actor")
	}

	var systemState system9.State
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
	migrations := make(map[cid.Cid]migration.ActorMigration)
	// Set of prior version code CIDs for actors to defer during iteration, for explicit migration afterwards.
	deferredCodeIDs := make(map[cid.Cid]struct{})

	oldInitCodeCID := cid.Undef
	for _, oldEntry := range oldManifestData.Entries {
		newCodeCID, ok := newManifest.Get(oldEntry.Name)
		if !ok {
			return cid.Undef, xerrors.Errorf("code cid for %s actor not found in new manifest", oldEntry.Name)
		}

		migrations[oldEntry.Code] = migration.CodeMigrator{OutCodeCID: newCodeCID}

		if oldEntry.Name == manifest.InitKey {
			oldInitCodeCID = oldEntry.Code
		}
	}

	if !oldInitCodeCID.Defined() {
		return cid.Undef, xerrors.New("didn't find init actor in old manifest")
	}
	// migrations that migrate both code and state, override entries in `migrations`

	// The System Actor
	newSystemCodeCID, ok := newManifest.Get(manifest.SystemKey)
	if !ok {
		return cid.Undef, xerrors.Errorf("code cid for system actor not found in new manifest")
	}

	migrations[systemActor.Code] = systemActorMigrator{OutCodeCID: newSystemCodeCID, ManifestData: newManifest.Data}

	// The Init Actor
	newInitCodeCID, ok := newManifest.Get(manifest.InitKey)
	if !ok {
		return cid.Undef, xerrors.Errorf("code cid for init actor not found in new manifest")
	}

	ethZeroAddr, err := makeEthZeroAddress()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to make eth zero address: %w", err)
	}

	migrations[oldInitCodeCID] = initActorMigrator{OutCodeCID: newInitCodeCID, EthZeroAddress: ethZeroAddr}
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
		if err = actorsIn.ForEachV4(func(addr address.Address, actorIn *builtin.ActorV4) error {
			if _, ok := deferredCodeIDs[actorIn.Code]; ok {
				return nil
			}

			migration, ok := migrations[actorIn.Code]
			if !ok {
				return xerrors.Errorf("actor with code %s has no registered migration function", actorIn.Code)
			}

			nextInput := &migrationJob{
				Address:        addr,
				ActorV4:        *actorIn, // Must take a copy, the pointer is not stable.
				cache:          cache,
				ActorMigration: migration,
			}

			select {
			case jobCh <- nextInput:
			case <-ctx.Done():
				return ctx.Err()
			}
			atomic.AddUint32(&jobCount, 1)
			return nil
		}); err != nil {
			return xerrors.Errorf("error iterating v4 actors: %w", err)
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
					return xerrors.Errorf("running job: %w", err)
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
			if err := actorsOut.SetActorV5(result.Address, &result.ActorV5); err != nil {
				return xerrors.Errorf("error setting actor %s: %w", result.Address, err)
			}
			resultCount++
		}
		log.Log(rt.INFO, "Result writer wrote %d results to state tree after %v", resultCount, time.Since(startTime))
		return nil
	})

	if err := grp.Wait(); err != nil {
		return cid.Undef, xerrors.Errorf("migration group error: %w", err)
	}

	// Make the empty state

	emptyObj, err := makeEmptyState(store)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to make empty obj: %w", err)
	}

	// Create the EAM Actor

	eamAct, err := CreateEAMActor(&newManifest, emptyObj)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to create the EAM Actor: %w", err)
	}

	if err := actorsOut.SetActorV5(builtin.EthereumAddressManagerActorAddr, eamAct); err != nil {
		return cid.Undef, xerrors.Errorf("failed to set EAM Actor: %w", err)
	}

	// Create the EthZeroAddress as an EthAccount
	initAct, ok, err := actorsOut.GetActorV5(builtin.InitActorAddr)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to get init actor: %w", err)
	}
	if !ok {
		return cid.Undef, xerrors.New("didn't find init actor")
	}

	var initSt init10.State
	if err := store.Get(ctx, initAct.Head, &initSt); err != nil {
		return cid.Undef, xerrors.Errorf("failed to get init state: %w", err)
	}

	ethZeroAddrId, ok, err := initSt.ResolveAddress(adt9.WrapStore(ctx, store), ethZeroAddr)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to get ethZero actor: %w", err)
	}
	if !ok {
		return cid.Undef, xerrors.New("didn't find ethZero actor")
	}

	ethAccountCID, ok := newManifest.Get(manifest.EthAccountKey)
	if !ok {
		return cid.Undef, xerrors.New("didn't find ethAccount code CID")
	}

	if err := actorsOut.SetActorV5(ethZeroAddrId, &builtin.ActorV5{
		Code:       ethAccountCID,
		Head:       emptyObj,
		CallSeqNum: 0,
		Balance:    abi.NewTokenAmount(0),
		Address:    &ethZeroAddr,
	}); err != nil {
		return cid.Undef, xerrors.Errorf("failed to set ethZeroActor: %w", err)
	}

	elapsed := time.Since(startTime)
	rate := float64(doneCount) / elapsed.Seconds()
	log.Log(rt.INFO, "All %d done after %v (%.0f/s), flushing state root.", doneCount, elapsed, rate)

	outCid, err := actorsOut.Flush()
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to flush actorsOut: %w", err)
	}

	return outCid, nil
}

type migrationJob struct {
	address.Address
	builtin.ActorV4
	migration.ActorMigration
	cache migration.MigrationCache
}

type migrationJobResult struct {
	address.Address
	builtin.ActorV5
}

func (job *migrationJob) run(ctx context.Context, store cbor.IpldStore, priorEpoch abi.ChainEpoch) (*migrationJobResult, error) {
	result, err := job.MigrateState(ctx, store, migration.ActorMigrationInput{
		Address:    job.Address,
		Head:       job.ActorV4.Head,
		PriorEpoch: priorEpoch,
		Cache:      job.cache,
	})
	if err != nil {
		return nil, xerrors.Errorf("state migration failed for actor code %s, addr %s: %w",
			job.ActorV4.Code, job.Address, err)
	}

	// Set up new actor record with the migrated state.
	return &migrationJobResult{
		job.Address, // Unchanged
		builtin.ActorV5{
			Code:       result.NewCodeCID,
			Head:       result.NewHead,
			CallSeqNum: job.ActorV4.CallSeqNum, // Unchanged
			Balance:    job.ActorV4.Balance,    // Unchanged
			Address:    nil,                    // TODO check that this does not need to be populated
		},
	}, nil
}

func makeEthZeroAddress() (address.Address, error) {
	return address.NewDelegatedAddress(builtin.EthereumAddressManagerActorID, make([]byte, 20))
}

func makeEmptyState(store cbor.IpldStore) (cid.Cid, error) {
	emptyObject, err := store.Put(context.TODO(), []struct{}{})
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to make empty object: %w", err)
	}

	return emptyObject, nil
}
