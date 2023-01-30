package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/testutil/helpers"
)

func TestUpgradeStateOnGenesis(t *testing.T) {
	myApp := helpers.Setup(false, false)

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
