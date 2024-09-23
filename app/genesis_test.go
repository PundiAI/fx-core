package app_test

import (
	"testing"

	"github.com/functionx/fx-core/v8/app"
	"github.com/functionx/fx-core/v8/testutil/helpers"
)

func TestNewDefaultGenesisByDenom(t *testing.T) {
	myApp := helpers.NewApp()
	genAppState := app.NewDefAppGenesisByDenom(myApp.AppCodec(), myApp.ModuleBasics)

	helpers.AssertJsonFile(t, "./genesis.json", genAppState)
}
