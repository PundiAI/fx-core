package server

import (
	tmcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	ethermintserver "github.com/evmos/ethermint/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TendermintCommand() *cobra.Command {
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
		UnsafeRestPrivValidatorCmd(),
		UnsafeResetNodeKeyCmd(),
		ReplayCmd(),
		ReplayConsoleCmd(),
		tmcmd.ResetAllCmd,
		tmcmd.ResetStateCmd,
		tmcmd.GenValidatorCmd,
		tmcmd.GenNodeKeyCmd,
		ethermintserver.NewIndexTxCmd(),
		// tmcmd.ResetPrivValidatorCmd,
	)
	return tendermintCmd
}

func CometCommand(appCreator types.AppCreator) *cobra.Command {
	cometCmd := &cobra.Command{
		Use:   "comet",
		Short: "CometBFT subcommands",
	}

	cometCmd.AddCommand(
		sdkserver.ShowNodeIDCmd(),
		sdkserver.ShowValidatorCmd(),
		sdkserver.ShowAddressCmd(),
		sdkserver.VersionCmd(),
		tmcmd.ResetAllCmd,
		tmcmd.ResetStateCmd,
		sdkserver.BootstrapStateCmd(appCreator),
	)
	return cometCmd
}
