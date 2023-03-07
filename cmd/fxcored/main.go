package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/functionx/fx-core/v3/cmd"
	fxtypes "github.com/functionx/fx-core/v3/types"
)

func main() {
	if err := svrcmd.Execute(cmd.NewRootCmd(), fxtypes.EnvPrefix, fxtypes.GetDefaultNodeHome()); err != nil {
		os.Exit(1)
	}
}
