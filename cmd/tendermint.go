package cmd

import (
	"errors"
	"time"

	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/evmos/ethermint/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/store"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/functionx/fx-core/v3/client/cli"
	fxtypes "github.com/functionx/fx-core/v3/types"
)

func tendermintCommand() *cobra.Command {
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
		cli.UnsafeRestPrivValidatorCmd(),
		cli.UnsafeResetNodeKeyCmd(),
		cli.ReplayCmd(),
		cli.ReplayConsoleCmd(),
		tmcmd.ResetAllCmd,
		tmcmd.ResetStateCmd,
		tmcmd.GenValidatorCmd,
		tmcmd.GenNodeKeyCmd,
		//tmcmd.ResetPrivValidatorCmd,
	)
	return tendermintCmd
}

func startCommand(appCreator types.AppCreator, defaultNodeHome string) *cobra.Command {
	startCmd := server.StartCmd(appCreator, defaultNodeHome)
	crisis.AddModuleInitFlags(startCmd)

	startCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		serverCtx := sdkserver.GetServerContextFromCmd(cmd)

		if zeroLog, ok := serverCtx.Logger.(sdkserver.ZeroLogWrapper); ok {
			filterLogTypes, _ := cmd.Flags().GetStringSlice(cli.FlagLogFilter)
			if len(filterLogTypes) > 0 {
				serverCtx.Logger = cli.NewFxZeroLogWrapper(zeroLog, filterLogTypes)
			}
		}

		// Bind flags to the Context's Viper so the app construction can set
		// options accordingly.
		if err := serverCtx.Viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		if _, err := sdkserver.GetPruningOptionsFromFlags(serverCtx.Viper); err != nil {
			return err
		}

		genesisDoc, err := tmtypes.GenesisDocFromFile(serverCtx.Config.GenesisFile())
		if err != nil {
			return err
		}
		if err = checkMainnetAndBlock(genesisDoc, serverCtx.Config); err != nil {
			return err
		}
		fxtypes.SetChainId(genesisDoc.ChainID)
		return nil
	}
	startCmd.Flags().StringSlice(cli.FlagLogFilter, nil, `The logging filter can discard custom log type (ABCIQuery)`)
	return startCmd
}

func checkMainnetAndBlock(genesisDoc *tmtypes.GenesisDoc, config *config.Config) error {
	if genesisDoc.InitialHeight > 1 || genesisDoc.ChainID != fxtypes.MainnetChainId || config.StateSync.Enable {
		return nil
	}
	genesisTime, err := time.Parse("2006-01-02T15:04:05Z", "2021-07-05T04:00:00Z")
	if err != nil {
		return err
	}
	blockStoreDB, err := node.DefaultDBProvider(&node.DBContext{ID: "blockstore", Config: config})
	if err != nil {
		return err
	}
	defer blockStoreDB.Close()
	blockStore := store.NewBlockStore(blockStoreDB)
	if genesisDoc.GenesisTime.Equal(genesisTime) && blockStore.Height() <= 0 {
		return errors.New("invalid version: Sync block from scratch please use use fxcored v1.1.x")
	}
	return nil
}
