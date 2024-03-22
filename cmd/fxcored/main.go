package main

import (
	"errors"
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/functionx/fx-core/v7/cmd"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

func main() {
	if err := svrcmd.Execute(cmd.NewRootCmd(), fxtypes.EnvPrefix, fxtypes.GetDefaultNodeHome()); err != nil {

		var e server.ErrorCode
		switch {
		case errors.As(err, &e):
			os.Exit(e.Code)
		default:
			os.Exit(1)
		}
	}
}
