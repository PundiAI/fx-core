package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
)

func TestAppDB(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	chainId := fxtypes.MainnetChainId
	myApp := buildApp(t)
	require.NoError(t, myApp.LoadLatestVersion())
	ctx := newContext(t, myApp, chainId, false)

	_ = ctx

	// do something ...
}
