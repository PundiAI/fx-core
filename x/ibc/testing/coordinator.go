package ibctesting

import (
	"testing"
	"time"

	ibctesting "github.com/cosmos/ibc-go/v6/testing"
)

var globalStartTime = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)

// NewCoordinator initializes Coordinator with N Fxcore TestChain's (fx-core apps) and M Cosmos chains (Simulation Apps)
func NewCoordinator(t *testing.T, nFxChains, mCosmosChains int) *ibctesting.Coordinator {
	chains := make(map[string]*ibctesting.TestChain)
	coord := &ibctesting.Coordinator{
		T:           t,
		CurrentTime: globalStartTime,
	}

	ibctesting.DefaultTestingAppInit = DefaultTestingAppInit

	for i := 1; i <= nFxChains; i++ {
		chainID := ibctesting.GetChainID(i)
		chains[chainID] = NewTestChain(t, coord, chainID)
	}

	// setup Cosmos chains
	ibctesting.DefaultTestingAppInit = ibctesting.SetupTestingApp

	for j := 1 + nFxChains; j <= nFxChains+mCosmosChains; j++ {
		chainID := ibctesting.GetChainID(j)
		chains[chainID] = ibctesting.NewTestChain(t, coord, chainID)
	}

	coord.Chains = chains

	return coord
}
