package main

import (
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	types2 "github.com/tendermint/tendermint/types"

	fxtypes "github.com/functionx/fx-core/v2/types"

	"github.com/functionx/fx-core/v2/app/cli"

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
		tmcmd.RollbackStateCmd,
		tmcmd.ReplayCmd,
		tmcmd.ReplayConsoleCmd,
		tmcmd.GenValidatorCmd,
		tmcmd.GenNodeKeyCmd,
		tmcmd.VersionCmd,
		//tmcmd.ResetPrivValidatorCmd
	)

	startCmd := server.StartCmd(appCreator, defaultNodeHome)
	crisis.AddModuleInitFlags(startCmd)
	preRun := startCmd.PreRunE
	startCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if err := preRun(cmd, args); err != nil {
			return err
		}
		serverCtx := sdkserver.GetServerContextFromCmd(cmd)
		genesisDoc, err := types2.GenesisDocFromFile(serverCtx.Config.GenesisFile())
		if err != nil {
			return err
		}
		fxtypes.SetChainId(genesisDoc.ChainID)
		return nil
	}
	rootCmd.AddCommand(
		startCmd,
		tendermintCmd,
		cli.ExportSateCmd(appExport, defaultNodeHome),
		version.NewVersionCommand(),
		sdkserver.NewRollbackCmd(defaultNodeHome),
	)
}
