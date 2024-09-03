package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

func TestAppDB(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	myApp, chainId := buildApp(t, fxtypes.MainnetChainId)
	require.NoError(t, myApp.LoadLatestVersion())
	ctx := newContext(t, myApp, chainId, false)

	_ = ctx

	// do something ...
}
