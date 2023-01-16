package builtin

import (
	"context"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"
)

///// Code shared by multiple built-in actors. /////

// Default log2 of branching factor for HAMTs.
// This value has been empirically chosen, but the optimal value for maps with different mutation profiles may differ.
const DefaultHamtBitwidth = 5

const DefaultTokenActorBitwidth = 3

type BigFrac struct {
	Numerator   big.Int
	Denominator big.Int
}

func MakeEmptyState() (cid.Cid, error) {
	store := cbor.NewMemCborStore()
	emptyObject, err := store.Put(context.TODO(), []struct{}{})
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to make empty object: %w", err)
	}

	return emptyObject, nil
}
