package migration

import (
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/manifest"
)

func CreateEAMActor(m *manifest.Manifest, head cid.Cid) (*builtin.ActorV5, error) {
	eamCode, ok := m.Get(manifest.EamKey)
	if !ok {
		return nil, xerrors.Errorf("didn't find EAM code CID")
	}

	return &builtin.ActorV5{
		Code:       eamCode,
		Head:       head,
		CallSeqNum: 0,
		Balance:    abi.NewTokenAmount(0),
		Address:    nil,
	}, nil
}
