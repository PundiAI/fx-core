package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v7/testutil/helpers"
)

func TestUpgradeStateOnGenesis(t *testing.T) {
	valSet, genAccs, balances := helpers.GenerateGenesisValidator(1, nil)
	myApp := helpers.SetupWithGenesisValSet(t, valSet, genAccs, balances...)
	// make sure the upgrade keeper has version map in state
	ctx := myApp.NewContext(false, tmproto.Header{Height: myApp.LastBlockHeight()})
	vm := myApp.UpgradeKeeper.GetModuleVersionMap(ctx)
	modules := myApp.GetModules()
	require.Equal(t, len(vm), len(modules))
	for k, module := range modules {
		require.Equal(t, k, module.Name())
		require.Equal(t, vm[module.Name()], module.ConsensusVersion())
	}
}
