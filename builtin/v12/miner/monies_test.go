package miner

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v12/util/smoothing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitialPledgeForPower_NetworkQAPLessThanNetworkBaseline(t *testing.T) {
	pledge := initialConsensusPledge(
		abi.NewStoragePower(5),
		abi.NewStoragePower(1000),
		smoothing.NewEstimate(big.NewInt(500), big.NewInt(0)),
		abi.NewTokenAmount(10000))

	assert.Equal(t, pledge, big.NewInt(15))
}

func TestInitialPledgeForPower_NetworkQAPGreaterThanNetworkBaseline(t *testing.T) {
	pledge := initialConsensusPledge(
		abi.NewStoragePower(5),
		abi.NewStoragePower(1000),
		smoothing.NewEstimate(big.NewInt(1500), big.NewInt(0)),
		abi.NewTokenAmount(10000))

	assert.Equal(t, pledge, big.NewInt(10))
}
