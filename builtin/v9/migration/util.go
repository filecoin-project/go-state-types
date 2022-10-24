package migration

import (
	"context"
	"sync"

	adt9 "github.com/filecoin-project/go-state-types/builtin/v9/util/adt"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	market8 "github.com/filecoin-project/go-state-types/builtin/v8/market"
	adt8 "github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
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

func getPendingVerifiedDealsAndTotalSize(ctx context.Context, adtStore adt8.Store, marketStateV8 market8.State) ([]abi.DealID, uint64, error) {
	pendingProposals, err := adt8.AsSet(adtStore, marketStateV8.PendingProposals, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, 0, xerrors.Errorf("failed to load pending proposals: %w", err)
	}

	proposals, err := market8.AsDealProposalArray(adtStore, marketStateV8.Proposals)
	if err != nil {
		return nil, 0, xerrors.Errorf("failed to get proposals: %w", err)
	}

	// We only want those pending deals that haven't been activated -- an activated deal has an entry in dealStates8
	dealStates8, err := adt9.AsArray(adtStore, marketStateV8.States, market8.StatesAmtBitwidth)
	if err != nil {
		return nil, 0, xerrors.Errorf("failed to load v8 states array: %w", err)
	}

	var pendingVerifiedDeals []abi.DealID
	pendingSize := uint64(0)
	var proposal market8.DealProposal
	if err = proposals.ForEach(&proposal, func(dealID int64) error {
		// Nothing to do for unverified deals
		if !proposal.VerifiedDeal {
			return nil
		}

		pcid, err := proposal.Cid()
		if err != nil {
			return err
		}

		isPending, err := pendingProposals.Has(abi.CidKey(pcid))
		if err != nil {
			return xerrors.Errorf("failed to check pending: %w", err)
		}

		// Nothing to do for not-pending deals
		if !isPending {
			return nil
		}

		var _dealState8 market8.DealState
		found, err := dealStates8.Get(uint64(dealID), &_dealState8)
		if err != nil {
			return xerrors.Errorf("failed to lookup deal state: %w", err)
		}

		// the deal has an entry in deal states, which means it's already been allocated, nothing to do
		if found {
			return nil
		}

		pendingVerifiedDeals = append(pendingVerifiedDeals, abi.DealID(dealID))
		pendingSize += uint64(proposal.PieceSize)
		return nil
	}); err != nil {
		return nil, 0, xerrors.Errorf("failed to iterate over proposals: %w", err)
	}
	return pendingVerifiedDeals, pendingSize, nil
}
