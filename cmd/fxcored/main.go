package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/functionx/fx-core/app/fxcore"
	"github.com/functionx/fx-core/cmd/fxcored/cmd"
)

func main() {
	if err := svrcmd.Execute(cmd.NewRootCmd(), fxcore.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
