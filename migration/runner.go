package migration

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/builtin"
	adt10 "github.com/filecoin-project/go-state-types/builtin/v10/util/adt"
	"github.com/filecoin-project/go-state-types/rt"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

func RunMigration(ctx context.Context, cfg Config, cache MigrationCache, store cbor.IpldStore, log Logger, actorsIn *builtin.ActorTree, migrations map[cid.Cid]ActorMigration) (*builtin.ActorTree, error) {
	startTime := time.Now()

	deferred := map[address.Address]struct{}{}
	var deferredLk sync.Mutex

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
		log.Log(rt.INFO, "Creating migration jobs")
		if err := actorsIn.ForEachV5(func(addr address.Address, actorIn *builtin.ActorV5) error {
			actorMigration, ok := migrations[actorIn.Code]
			if !ok {
				return xerrors.Errorf("actor with code %s has no registered migration function", actorIn.Code)
			}

			nextInput := &migrationJob{
				Address:        addr,
				ActorV5:        *actorIn, // Must take a copy, the pointer is not stable.
				cache:          cache,
				ActorMigration: actorMigration,
			}

			select {
			case jobCh <- nextInput:
			case <-ctx.Done():
				return ctx.Err()
			}
			atomic.AddUint32(&jobCount, 1)
			return nil
		}); err != nil {
			return xerrors.Errorf("error iterating actors: %w", err)
		}
		log.Log(rt.INFO, "Done creating %d migration jobs after %v", jobCount, time.Since(startTime))
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
				if job.ActorMigration.Deferred() {
					deferredLk.Lock()
					deferred[job.Address] = struct{}{}
					deferredLk.Unlock()

					atomic.AddUint32(&doneCount, 1)
					continue
				}

				result, err := job.run(ctx, store)
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

	// Setup the new actors tree

	actorsOut, err := builtin.NewTree(adt10.WrapStore(ctx, store))
	if err != nil {
		return nil, xerrors.Errorf("creating new state tree: %w", err)
	}

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
		return nil, xerrors.Errorf("migration group error: %w", err)
	}

	if len(deferred) > 0 {
		// deferred round
		// NOTE this is not parralelized for now as this was only ever needed for singleton actor migrations

		log.Log(rt.INFO, "Running deferred migrations")

		for addr := range deferred {
			log.Log(rt.INFO, "Running deferred migration for %s", addr)

			actorIn, found, err := actorsIn.GetActorV5(addr)
			if err != nil {
				return nil, xerrors.Errorf("failed to get actor %s: %w", addr, err)
			}

			if !found {
				return nil, xerrors.Errorf("failed to find actor %s", addr)
			}

			actorMigration, ok := migrations[actorIn.Code]
			if !ok {
				return nil, xerrors.Errorf("actor with code %s has no registered migration function", actorIn.Code)
			}

			res, err := (&migrationJob{
				Address:        addr,
				ActorV5:        *actorIn,
				ActorMigration: actorMigration,
				cache:          cache,
			}).run(ctx, store)
			if err != nil {
				return nil, xerrors.Errorf("running deferred job: %w", err)
			}

			if err := actorsOut.SetActorV5(res.Address, &res.ActorV5); err != nil {
				return nil, xerrors.Errorf("error setting actor %s: %w", res.Address, err)
			}
		}
	}

	elapsed := time.Since(startTime)
	rate := float64(doneCount) / elapsed.Seconds()
	log.Log(rt.INFO, "All %d done after %v (%.0f/s)", doneCount, elapsed, rate)

	return actorsOut, nil
}
