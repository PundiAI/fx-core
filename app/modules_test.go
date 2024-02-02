package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v7/testutil/helpers"
)

func TestOnOrderBeginBlockers(t *testing.T) {
	valSet, genAccs, balances := helpers.GenerateGenesisValidator(1, nil)
	myApp := helpers.SetupWithGenesisValSet(t, valSet, genAccs, balances...)
	modules := myApp.GetModules()
	orderBeginBlockersModules := myApp.GetOrderBeginBlockersModules()
	require.Equal(t, len(orderBeginBlockersModules), len(modules))
	for _, moduleName := range orderBeginBlockersModules {
		_, ok := modules[moduleName]
		require.True(t, ok)
	}
}

func TestOnOrderEndBlockers(t *testing.T) {
	valSet, genAccs, balances := helpers.GenerateGenesisValidator(1, nil)
	myApp := helpers.SetupWithGenesisValSet(t, valSet, genAccs, balances...)
	modules := myApp.GetModules()
	orderEndBlockersModules := myApp.GetOrderEndBlockersModules()
	require.Equal(t, len(orderEndBlockersModules), len(modules))
	for _, moduleName := range orderEndBlockersModules {
		_, ok := modules[moduleName]
		require.True(t, ok)
	}
}

func TestOnOrderInitGenesis(t *testing.T) {
	valSet, genAccs, balances := helpers.GenerateGenesisValidator(1, nil)
	myApp := helpers.SetupWithGenesisValSet(t, valSet, genAccs, balances...)
	modules := myApp.GetModules()
	initGenesisModules := myApp.GetOrderInitGenesisModules()
	require.Equal(t, len(initGenesisModules), len(modules))
	for _, moduleName := range initGenesisModules {
		_, ok := modules[moduleName]
		require.True(t, ok)
	}
}
