package app_test

import (
	"testing"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
)

func TestNewDefaultGenesisByDenom(t *testing.T) {
	myApp := helpers.NewApp()
	genAppState := myApp.DefaultGenesis()

	helpers.AssertJsonFile(t, "./genesis.json", genAppState)
}
