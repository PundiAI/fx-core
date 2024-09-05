package app_test

import (
	"testing"

	"github.com/functionx/fx-core/v8/app"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
)

func TestNewDefaultGenesisByDenom(t *testing.T) {
	encodingConfig := app.MakeEncodingConfig()
	genAppState := app.NewDefAppGenesisByDenom(fxtypes.DefaultDenom, encodingConfig.Codec)

	helpers.AssertJsonFile(t, "./genesis.json", genAppState)
}
