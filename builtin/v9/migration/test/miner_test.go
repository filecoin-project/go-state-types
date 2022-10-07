package migration

import (
	"context"
	"testing"

	commp "github.com/filecoin-project/go-commp-utils/nonffi"
	market8 "github.com/filecoin-project/go-state-types/builtin/v8/market"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	miner8 "github.com/filecoin-project/go-state-types/builtin/v8/miner"
	"github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
	"github.com/filecoin-project/go-state-types/builtin/v9/migration"
	miner9 "github.com/filecoin-project/go-state-types/builtin/v9/miner"
	system9 "github.com/filecoin-project/go-state-types/builtin/v9/system"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/stretchr/testify/require"
)

func TestMinerMigration(t *testing.T) {
	ctx := context.Background()
	bs := cbor.NewCborStore(NewSyncBlockStoreInMemory())
	adtStore := adt.WrapStore(ctx, bs)

	startRoot := makeInputTree(ctx, t, adtStore)

	oldStateTree, err := migration.LoadTree(adtStore, startRoot)
	require.NoError(t, err)

	oldSystemActor, found, err := oldStateTree.GetActor(builtin.SystemActorAddr)
	require.NoError(t, err)
	require.True(t, found, "system actor not found")

	var oldSystemState system9.State
	err = adtStore.Get(ctx, oldSystemActor.Head, &oldSystemState)
	require.NoError(t, err)

	oldManifestDataCid := oldSystemState.BuiltinActors
	oldManifest := manifest.Manifest{
		Version: 1,
		Data:    oldManifestDataCid,
	}
	require.NoError(t, oldManifest.Load(ctx, adtStore), "failed to load old manifest")

	baseAddr := uint64(10000)
	baseAddrId, err := address.NewIDAddress(baseAddr)
	require.NoError(t, err)
	baseWorkerAddrId, err := address.NewIDAddress(baseAddr + 100)
	require.NoError(t, err)

	// create 3 deal proposals

	oldMarketAct, ok, err := oldStateTree.GetActor(builtin.StorageMarketActorAddr)
	require.NoError(t, err)
	require.True(t, ok)

	var oldMarketSt market8.State
	require.NoError(t, adtStore.Get(ctx, oldMarketAct.Head, &oldMarketSt))

	proposals, err := market8.AsDealProposalArray(adtStore, oldMarketSt.Proposals)
	baseDeal := market8.DealProposal{
		PieceCID:             cid.Undef,
		PieceSize:            512,
		VerifiedDeal:         false,
		Client:               baseAddrId,
		Provider:             baseAddrId,
		Label:                market8.EmptyDealLabel,
		StartEpoch:           0,
		EndEpoch:             0,
		StoragePricePerEpoch: big.Zero(),
		ProviderCollateral:   big.Zero(),
		ClientCollateral:     big.Zero(),
	}

	deal0 := baseDeal
	deal0.PieceCID = MakeCID("0", &market8.PieceCIDPrefix)
	require.NoError(t, err)

	deal1 := baseDeal
	deal1.PieceCID = MakeCID("1", &market8.PieceCIDPrefix)
	require.NoError(t, err)

	deal2 := baseDeal
	deal2.PieceCID = MakeCID("2", &market8.PieceCIDPrefix)
	require.NoError(t, err)

	require.NoError(t, proposals.Set(abi.DealID(0), &deal0))
	require.NoError(t, proposals.Set(abi.DealID(1), &deal1))
	require.NoError(t, proposals.Set(abi.DealID(2), &deal2))

	proposalsCID, err := proposals.Root()
	require.NoError(t, err)
	oldMarketSt.Proposals = proposalsCID
	oldMarketStCID, err := adtStore.Put(ctx, &oldMarketSt)
	require.NoError(t, err)
	oldMarketAct.Head = oldMarketStCID
	require.NoError(t, oldStateTree.SetActor(builtin.StorageMarketActorAddr, oldMarketAct))

	// base stuff to create miners

	oldMinerCID, ok := oldManifest.Get("storageminer")
	require.True(t, ok)

	baseMinerSt := makeBaseMinerState(ctx, t, adtStore, baseAddrId, baseWorkerAddrId)

	basePrecommit := miner8.SectorPreCommitOnChainInfo{
		Info: miner8.SectorPreCommitInfo{
			SealProof:              abi.RegisteredSealProof_StackedDrg32GiBV1_1,
			SectorNumber:           0,
			SealedCID:              MakeCID("100", &miner8.SealedCIDPrefix),
			SealRandEpoch:          0,
			DealIDs:                nil,
			Expiration:             0,
			ReplaceCapacity:        false,
			ReplaceSectorDeadline:  0,
			ReplaceSectorPartition: 0,
			ReplaceSectorNumber:    0,
		},
		PreCommitDeposit:   big.Zero(),
		PreCommitEpoch:     0,
		DealWeight:         big.Zero(),
		VerifiedDealWeight: big.Zero(),
	}

	// make 3 miners
	// miner1 has no precommits at all
	// miner2 has 4 precommits, but with no deals
	// miner3 has 3 precommits, with deals [0], [1,2], and []

	// miner1 has no precommits at all

	miner1StCid, err := adtStore.Put(ctx, &baseMinerSt)
	require.NoError(t, err)

	miner1 := migration.Actor{
		Code:       oldMinerCID,
		Head:       miner1StCid,
		CallSeqNum: 0,
		Balance:    big.Zero(),
	}

	addr1, err := address.NewIDAddress(baseAddr + 1)
	require.NoError(t, err)
	require.NoError(t, oldStateTree.SetActor(addr1, &miner1))

	// miner2 has precommits, but they have no deals

	precommits2, err := adt.AsMap(adtStore, baseMinerSt.PreCommittedSectors, builtin.DefaultHamtBitwidth)
	require.NoError(t, err)
	require.NoError(t, precommits2.Put(miner8.SectorKey(0), &basePrecommit))
	require.NoError(t, precommits2.Put(miner8.SectorKey(1), &basePrecommit))
	require.NoError(t, precommits2.Put(miner8.SectorKey(2), &basePrecommit))
	require.NoError(t, precommits2.Put(miner8.SectorKey(3), &basePrecommit))

	precommits2CID, err := precommits2.Root()
	require.NoError(t, err)
	miner2St := baseMinerSt
	miner2St.PreCommittedSectors = precommits2CID

	miner2StCid, err := adtStore.Put(ctx, &miner2St)
	require.NoError(t, err)

	miner2 := migration.Actor{
		Code:       oldMinerCID,
		Head:       miner2StCid,
		CallSeqNum: 0,
		Balance:    big.Zero(),
	}

	addr2, err := address.NewIDAddress(baseAddr + 2)
	require.NoError(t, err)
	require.NoError(t, oldStateTree.SetActor(addr2, &miner2))

	// miner 3 has precommits, some of which have deals

	precommits3, err := adt.AsMap(adtStore, baseMinerSt.PreCommittedSectors, builtin.DefaultHamtBitwidth)
	require.NoError(t, err)

	precommit3dotZero := basePrecommit
	precommit3dotZero.Info.DealIDs = []abi.DealID{0}
	precommit3dotZero.Info.SectorNumber = 0

	precommit3dotOne := basePrecommit
	precommit3dotOne.Info.DealIDs = []abi.DealID{1, 2}
	precommit3dotOne.Info.SectorNumber = 1

	precommit3dotTwo := basePrecommit
	precommit3dotTwo.Info.SectorNumber = 2

	require.NoError(t, precommits3.Put(miner8.SectorKey(0), &precommit3dotZero))
	require.NoError(t, precommits3.Put(miner8.SectorKey(1), &precommit3dotOne))
	require.NoError(t, precommits3.Put(miner8.SectorKey(2), &precommit3dotTwo))

	precommits3CID, err := precommits3.Root()
	require.NoError(t, err)
	miner3St := baseMinerSt
	miner3St.PreCommittedSectors = precommits3CID

	miner3StCid, err := adtStore.Put(ctx, &miner3St)
	require.NoError(t, err)

	miner3 := migration.Actor{
		Code:       oldMinerCID,
		Head:       miner3StCid,
		CallSeqNum: 0,
		Balance:    big.Zero(),
	}

	addr3, err := address.NewIDAddress(baseAddr + 3)
	require.NoError(t, err)
	require.NoError(t, oldStateTree.SetActor(addr3, &miner3))

	startRoot, err = oldStateTree.Flush()
	require.NoError(t, err)

	newManifestCid, _ := makeTestManifest(t, adtStore, "fil/9/")
	log := TestLogger{TB: t}

	cache := migration.NewMemMigrationCache()
	_, err = migration.MigrateStateTree(ctx, adtStore, newManifestCid, startRoot, 200, migration.Config{MaxWorkers: 1}, log, cache)
	require.NoError(t, err)

	cacheRoot, err := migration.MigrateStateTree(ctx, adtStore, newManifestCid, startRoot, 200, migration.Config{MaxWorkers: 1}, log, cache)
	require.NoError(t, err)

	noCacheRoot, err := migration.MigrateStateTree(ctx, adtStore, newManifestCid, startRoot, 200, migration.Config{MaxWorkers: 1}, log, migration.NewMemMigrationCache())
	require.NoError(t, err)
	require.True(t, cacheRoot.Equals(noCacheRoot))

	// check that the actor states were correctly updated

	newStateTree, err := migration.LoadTree(adtStore, cacheRoot)
	require.NoError(t, err)

	// miner 1 is just empty precommits

	newMiner1Actor, ok, err := newStateTree.GetActor(addr1)
	require.NoError(t, err)
	require.True(t, ok)

	var newMiner1St miner9.State
	require.NoError(t, adtStore.Get(ctx, newMiner1Actor.Head, &newMiner1St))
	require.Equal(t, newMiner1St.PreCommittedSectors, baseMinerSt.PreCommittedSectors)

	// miner 2's precommits all have nil unsealedCID

	newMiner2Actor, ok, err := newStateTree.GetActor(addr2)
	require.NoError(t, err)
	require.True(t, ok)

	var newMiner2St miner9.State
	require.NoError(t, adtStore.Get(ctx, newMiner2Actor.Head, &newMiner2St))
	require.NotEqual(t, newMiner2St.PreCommittedSectors, baseMinerSt.PreCommittedSectors)

	newPrecommits2, err := adt.AsMap(adtStore, newMiner2St.PreCommittedSectors, builtin.DefaultHamtBitwidth)
	require.NoError(t, err)
	var pc miner9.SectorPreCommitOnChainInfo
	require.NoError(t,
		newPrecommits2.ForEach(&pc, func(key string) error {
			require.Nil(t, pc.Info.UnsealedCid)
			return nil
		}))

	// miner 3's precommits must be checked one at a time
	// Sector 0 has deal [0]
	// Sector 1 has deals [1,2]
	// Sector 2 is empty

	newMiner3Actor, ok, err := newStateTree.GetActor(addr3)
	require.NoError(t, err)
	require.True(t, ok)

	var newMiner3St miner9.State
	require.NoError(t, adtStore.Get(ctx, newMiner3Actor.Head, &newMiner3St))
	require.NotEqual(t, newMiner3St.PreCommittedSectors, baseMinerSt.PreCommittedSectors)

	newPrecommits3, err := adt.AsMap(adtStore, newMiner3St.PreCommittedSectors, builtin.DefaultHamtBitwidth)
	require.NoError(t, err)

	ok, err = newPrecommits3.Get(miner9.SectorKey(0), &pc)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, pc.Info.UnsealedCid)
	expectedCid, err := commp.GenerateUnsealedCID(abi.RegisteredSealProof_StackedDrg32GiBV1_1, []abi.PieceInfo{{PieceCID: deal0.PieceCID, Size: deal0.PieceSize}})
	require.NoError(t, err)

	require.Equal(t, expectedCid, *pc.Info.UnsealedCid)

	ok, err = newPrecommits3.Get(miner9.SectorKey(1), &pc)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, pc.Info.UnsealedCid)
	expectedCid, err = commp.GenerateUnsealedCID(abi.RegisteredSealProof_StackedDrg32GiBV1_1, []abi.PieceInfo{
		{PieceCID: deal1.PieceCID, Size: deal1.PieceSize},
		{PieceCID: deal2.PieceCID, Size: deal2.PieceSize}})
	require.NoError(t, err)

	require.Equal(t, expectedCid, *pc.Info.UnsealedCid)

	ok, err = newPrecommits3.Get(miner9.SectorKey(2), &pc)
	require.NoError(t, err)
	require.True(t, ok)
	require.Nil(t, pc.Info.UnsealedCid)

	// we didn't introduce a 4th sector somehow

	ok, err = newPrecommits3.Get(miner9.SectorKey(3), &pc)
	require.NoError(t, err)
	require.False(t, ok)
}

func TestFip0029MinerMigration(t *testing.T) {
	ctx := context.Background()
	bs := cbor.NewCborStore(NewSyncBlockStoreInMemory())
	adtStore := adt.WrapStore(ctx, bs)

	startRoot := makeInputTree(ctx, t, adtStore)

	oldStateTree, err := migration.LoadTree(adtStore, startRoot)
	require.NoError(t, err)

	oldSystemActor, found, err := oldStateTree.GetActor(builtin.SystemActorAddr)
	require.NoError(t, err)
	require.True(t, found, "system actor not found")

	var oldSystemState system9.State
	err = adtStore.Get(ctx, oldSystemActor.Head, &oldSystemState)
	require.NoError(t, err)

	oldManifestDataCid := oldSystemState.BuiltinActors
	oldManifest := manifest.Manifest{
		Version: 1,
		Data:    oldManifestDataCid,
	}
	require.NoError(t, oldManifest.Load(ctx, adtStore), "failed to load old manifest")

	addr, err := address.NewIDAddress(10000)
	require.NoError(t, err)
	workerAddr, err := address.NewIDAddress(20000)
	require.NoError(t, err)

	oldMinerCID, ok := oldManifest.Get("storageminer")
	require.True(t, ok)

	minerSt := makeBaseMinerState(ctx, t, adtStore, addr, workerAddr)

	minerStCid, err := adtStore.Put(ctx, &minerSt)
	require.NoError(t, err)

	var minerInfo miner8.MinerInfo
	require.NoError(t, adtStore.Get(ctx, minerSt.Info, &minerInfo))

	miner := migration.Actor{
		Code:       oldMinerCID,
		Head:       minerStCid,
		CallSeqNum: 0,
		Balance:    big.Zero(),
	}

	require.NoError(t, oldStateTree.SetActor(addr, &miner))

	startRoot, err = oldStateTree.Flush()
	require.NoError(t, err)

	newManifestCid, _ := makeTestManifest(t, adtStore, "fil/9/")
	log := TestLogger{TB: t}

	cache := migration.NewMemMigrationCache()
	_, err = migration.MigrateStateTree(ctx, adtStore, newManifestCid, startRoot, 200, migration.Config{MaxWorkers: 1}, log, cache)
	require.NoError(t, err)

	cacheRoot, err := migration.MigrateStateTree(ctx, adtStore, newManifestCid, startRoot, 200, migration.Config{MaxWorkers: 1}, log, cache)
	require.NoError(t, err)

	noCacheRoot, err := migration.MigrateStateTree(ctx, adtStore, newManifestCid, startRoot, 200, migration.Config{MaxWorkers: 1}, log, migration.NewMemMigrationCache())
	require.NoError(t, err)
	require.True(t, cacheRoot.Equals(noCacheRoot))

	// check that the actor states were correctly updated

	newStateTree, err := migration.LoadTree(adtStore, cacheRoot)
	require.NoError(t, err)

	newMinerActor, ok, err := newStateTree.GetActor(addr)
	require.NoError(t, err)
	require.True(t, ok)

	var newMinerSt miner9.State
	require.NoError(t, adtStore.Get(ctx, newMinerActor.Head, &newMinerSt))

	var newMinerInfo miner9.MinerInfo
	require.NoError(t, adtStore.Get(ctx, newMinerSt.Info, &newMinerInfo))
	require.Equal(t, newMinerInfo.Owner, minerInfo.Owner)
	require.Equal(t, newMinerInfo.Worker, minerInfo.Worker)
	require.Equal(t, newMinerInfo.Beneficiary, minerInfo.Owner)
	require.Nil(t, newMinerInfo.PendingBeneficiaryTerm)
}

func makeBaseMinerState(ctx context.Context, t *testing.T, adtStore adt.Store, baseAddrId address.Address, baseWorkerAddrId address.Address) miner8.State {
	emptyPrecommitMapCid, err := adt.StoreEmptyMap(adtStore, builtin.DefaultHamtBitwidth)
	require.NoError(t, err, "failed to construct empty map")

	emptyMinerInfo := miner8.MinerInfo{
		Owner:                      baseAddrId,
		Worker:                     baseWorkerAddrId,
		ControlAddresses:           nil,
		PendingWorkerKey:           nil,
		PeerId:                     nil,
		Multiaddrs:                 nil,
		WindowPoStProofType:        0,
		SectorSize:                 0,
		WindowPoStPartitionSectors: 0,
		ConsensusFaultElapsed:      0,
		PendingOwnerAddress:        nil,
	}

	emptyMinerInfoCID, err := adtStore.Put(ctx, &emptyMinerInfo)
	require.NoError(t, err)

	emptyVestingFundsCID, err := adtStore.Put(ctx, &miner8.VestingFunds{Funds: nil})
	require.NoError(t, err)

	emptyPrecommitsCleanUpArrayCid, err := adt.StoreEmptyArray(adtStore, miner8.PrecommitCleanUpAmtBitwidth)
	require.NoError(t, err)
	emptySectorsArrayCid, err := adt.StoreEmptyArray(adtStore, miner8.SectorsAmtBitwidth)
	require.NoError(t, err)

	emptyBitfield := bitfield.NewFromSet(nil)
	emptyBitfieldCid, err := adtStore.Put(ctx, emptyBitfield)
	require.NoError(t, err)

	emptyPartitionsArrayCid, err := adt.StoreEmptyArray(adtStore, miner8.DeadlinePartitionsAmtBitwidth)
	require.NoError(t, err)

	emptyDeadlineExpirationArrayCid, err := adt.StoreEmptyArray(adtStore, miner8.DeadlineExpirationAmtBitwidth)
	require.NoError(t, err)

	emptyPoStSubmissionsArrayCid, err := adt.StoreEmptyArray(adtStore, miner8.DeadlineOptimisticPoStSubmissionsAmtBitwidth)
	require.NoError(t, err)

	emptyDeadlineCID, err := adtStore.Put(ctx, &miner8.Deadline{
		Partitions:        emptyPartitionsArrayCid,
		ExpirationsEpochs: emptyDeadlineExpirationArrayCid,
		PartitionsPoSted:  bitfield.New(),
		EarlyTerminations: bitfield.New(),
		LiveSectors:       0,
		TotalSectors:      0,
		FaultyPower: miner8.PowerPair{
			Raw: big.Zero(),
			QA:  big.Zero(),
		},
		OptimisticPoStSubmissions:         emptyPoStSubmissionsArrayCid,
		SectorsSnapshot:                   emptySectorsArrayCid,
		PartitionsSnapshot:                emptyPartitionsArrayCid,
		OptimisticPoStSubmissionsSnapshot: emptyPoStSubmissionsArrayCid,
	})
	require.NoError(t, err)

	emptyDeadlines := new(miner8.Deadlines)
	for i := range emptyDeadlines.Due {
		emptyDeadlines.Due[i] = emptyDeadlineCID
	}

	emptyDeadlinesCID, err := adtStore.Put(ctx, emptyDeadlines)
	require.NoError(t, err)

	return miner8.State{
		Info:                       emptyMinerInfoCID,
		PreCommitDeposits:          big.Zero(),
		LockedFunds:                big.Zero(),
		VestingFunds:               emptyVestingFundsCID,
		FeeDebt:                    big.Zero(),
		InitialPledge:              big.Zero(),
		PreCommittedSectors:        emptyPrecommitMapCid,
		PreCommittedSectorsCleanUp: emptyPrecommitsCleanUpArrayCid,
		AllocatedSectors:           emptyBitfieldCid,
		Sectors:                    emptySectorsArrayCid,
		ProvingPeriodStart:         0,
		CurrentDeadline:            0,
		Deadlines:                  emptyDeadlinesCID,
		EarlyTerminations:          emptyBitfield,
		DeadlineCronActive:         false,
	}
}
