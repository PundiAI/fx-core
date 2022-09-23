package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/functionx/fx-core/app"
)

func main() {
	if err := svrcmd.Execute(NewRootCmd(), app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
