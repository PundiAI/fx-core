package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/functionx/fx-core/v4/cmd"
	fxtypes "github.com/functionx/fx-core/v4/types"
)

func main() {
	if err := svrcmd.Execute(cmd.NewRootCmd(), fxtypes.EnvPrefix, fxtypes.GetDefaultNodeHome()); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)
		default:
			os.Exit(1)
		}
	}
}
