package main

import (
	"os"

	"github.com/functionx/fx-core/app"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/functionx/fx-core/cmd/fxcored/cmd"
)

func main() {
	if err := svrcmd.Execute(cmd.NewRootCmd(), app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
