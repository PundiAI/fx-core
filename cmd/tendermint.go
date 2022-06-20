package main

import (
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/functionx/fx-core/app/cli"

	sdkserver "github.com/cosmos/cosmos-sdk/server"

	tmcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"

	"github.com/evmos/ethermint/server"
)

func addTendermintCommands(rootCmd *cobra.Command, defaultNodeHome string, appCreator types.AppCreator, appExport types.AppExporter) {
	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint subcommands",
	}

	tmRestCmdPreRun := func(cmd *cobra.Command, args []string) {
		serverCtx := sdkserver.GetServerContextFromCmd(cmd)
		globalViper := viper.GetViper()
		for _, s := range serverCtx.Viper.AllKeys() {
			globalViper.Set(s, serverCtx.Viper.Get(s))
		}
	}

	tmcmd.ResetStateCmd.PreRun = tmRestCmdPreRun
	tmcmd.ResetAllCmd.PreRun = tmRestCmdPreRun

	tendermintCmd.AddCommand(
		sdkserver.ShowNodeIDCmd(),
		sdkserver.ShowValidatorCmd(),
		sdkserver.ShowAddressCmd(),
		sdkserver.VersionCmd(),
		tmcmd.ResetAllCmd,
		tmcmd.ResetStateCmd,
	)

	rootCmd.AddCommand(
		server.StartCmd(appCreator, defaultNodeHome),
		tendermintCmd,
		cli.ExportSateCmd(appExport, defaultNodeHome),
		version.NewVersionCommand(),
		sdkserver.NewRollbackCmd(defaultNodeHome),
	)
}
