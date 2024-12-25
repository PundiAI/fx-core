package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/pundiai/fx-core/v8/cmd"
	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func main() {
	if err := svrcmd.Execute(cmd.NewRootCmd(), fxtypes.EnvPrefix, fxtypes.GetDefaultNodeHome()); err != nil {
		os.Exit(1)
	}
}
