package migration

import (
	"context"

	"golang.org/x/xerrors"

	miner10 "github.com/filecoin-project/go-state-types/builtin/v10/miner"
	miner11 "github.com/filecoin-project/go-state-types/builtin/v11/miner"
	"github.com/filecoin-project/go-state-types/migration"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
)

// The minerMigrator performs the following migrations:
// - FIP-0061: Updates all miner info PoSt proof types to V1_1 types
type minerMigrator struct {
	OutCodeCID cid.Cid
}

func (m minerMigrator) MigratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m minerMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, in migration.ActorMigrationInput) (*migration.ActorMigrationResult, error) {
	var inState miner10.State
	if err := store.Get(ctx, in.Head, &inState); err != nil {
		return nil, xerrors.Errorf("failed to load miner state for %s: %w", in.Address, err)
	}
	var inInfo miner10.MinerInfo
	if err := store.Get(ctx, inState.Info, &inInfo); err != nil {
		return nil, xerrors.Errorf("failed to load miner info for %s: %w", in.Address, err)
	}

	outProof, err := inInfo.WindowPoStProofType.ToV1_1PostProof()
	if err != nil {
		return nil, xerrors.Errorf("failed to convert to v1_1 proof: %w", err)
	}

	outInfo := miner11.MinerInfo{
		Owner:                      inInfo.Owner,
		Worker:                     inInfo.Worker,
		ControlAddresses:           inInfo.ControlAddresses,
		PendingWorkerKey:           (*miner11.WorkerKeyChange)(inInfo.PendingWorkerKey),
		PeerId:                     inInfo.PeerId,
		Multiaddrs:                 inInfo.Multiaddrs,
		WindowPoStProofType:        outProof,
		SectorSize:                 inInfo.SectorSize,
		WindowPoStPartitionSectors: inInfo.WindowPoStPartitionSectors,
		ConsensusFaultElapsed:      inInfo.ConsensusFaultElapsed,
		PendingOwnerAddress:        inInfo.PendingOwnerAddress,
		Beneficiary:                inInfo.Beneficiary,
		BeneficiaryTerm:            miner11.BeneficiaryTerm(inInfo.BeneficiaryTerm),
		PendingBeneficiaryTerm:     (*miner11.PendingBeneficiaryChange)(inInfo.PendingBeneficiaryTerm),
	}

	outInfoCid, err := store.Put(ctx, &outInfo)
	if err != nil {
		return nil, xerrors.Errorf("failed to write new miner info: %w", err)
	}

	outState := miner11.State{
		Info:                       outInfoCid,
		PreCommitDeposits:          inState.PreCommitDeposits,
		LockedFunds:                inState.LockedFunds,
		VestingFunds:               inState.VestingFunds,
		FeeDebt:                    inState.FeeDebt,
		InitialPledge:              inState.InitialPledge,
		PreCommittedSectors:        inState.PreCommittedSectors,
		PreCommittedSectorsCleanUp: inState.PreCommittedSectorsCleanUp,
		AllocatedSectors:           inState.AllocatedSectors,
		Sectors:                    inState.Sectors,
		ProvingPeriodStart:         inState.ProvingPeriodStart,
		CurrentDeadline:            inState.CurrentDeadline,
		Deadlines:                  inState.Deadlines,
		EarlyTerminations:          inState.EarlyTerminations,
		DeadlineCronActive:         inState.DeadlineCronActive,
	}

	newHead, err := store.Put(ctx, &outState)
	if err != nil {
		return nil, xerrors.Errorf("failed to put new state: %w", err)
	}

	return &migration.ActorMigrationResult{
		NewCodeCID: m.MigratedCodeCID(),
		NewHead:    newHead,
	}, nil
}

func (m minerMigrator) Deferred() bool {
	return false
}
