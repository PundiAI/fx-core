package app_test

import (
	"testing"

	"github.com/functionx/fx-core/v8/testutil/helpers"
)

func TestNewDefaultGenesisByDenom(t *testing.T) {
	myApp := helpers.NewApp()
	genAppState := myApp.DefaultGenesis()

	helpers.AssertJsonFile(t, "./genesis.json", genAppState)
}
