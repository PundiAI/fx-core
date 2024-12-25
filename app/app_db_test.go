package app_test

import (
	"testing"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func TestAppDB(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	chainId := fxtypes.MainnetChainId
	myApp := buildApp(t)
	ctx := newContext(t, myApp, chainId, false)

	_ = ctx

	// do something ...
}
