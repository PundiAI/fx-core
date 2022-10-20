package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/functionx/fx-core/v2/app"
	"github.com/functionx/fx-core/v2/cmd"
)

func main() {
	if err := svrcmd.Execute(cmd.NewRootCmd(), app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
