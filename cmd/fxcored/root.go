package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdkCfg "github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/snapshots"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingcli "github.com/cosmos/cosmos-sdk/x/auth/vesting/client/cli"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tendermint/tendermint/config"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	store2 "github.com/tendermint/tendermint/store"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/app"
	appCli "github.com/functionx/fx-core/app/cli"
	fxtypes "github.com/functionx/fx-core/types"
)

const EnvPrefix = "FX"

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() *cobra.Command {

	encodingConfig := app.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(app.DefaultNodeHome).
		WithViper(EnvPrefix)

	rootCmd := &cobra.Command{
		Use:   fxtypes.Name + "d",
		Short: "FunctionX Core Chain App",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = sdkCfg.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			if err := server.InterceptConfigsPreRunHandler(cmd, "", nil); err != nil {
				return err
			}
			// add log filter
			return appCli.AddCmdLogWrapFilterLogType(cmd)
		},
	}

	rootCmd.PersistentFlags().StringSlice(appCli.FlagLogFilter, []string{}, `The logging filter can discard custom log type (ABCIQuery)`)
	initRootCmd(rootCmd, encodingConfig)
	overwriteFlagDefaults(rootCmd, map[string]string{
		flags.FlagChainID: fxtypes.ChainId(),
	})
	for _, command := range rootCmd.Commands() {
		if command.Use == "" {
			rootCmd.RemoveCommand(command)
		}
	}
	return rootCmd
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig app.EncodingConfig) {
	authclient.Codec = encodingConfig.Marshaler

	sdkCfgCmd := sdkCfg.Cmd()
	sdkCfgCmd.AddCommand(appCli.AppTomlCmd(), appCli.ConfigTomlCmd())

	rootCmd.AddCommand(
		InitCmd(),
		appCli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		genutilcli.MigrateGenesisCmd(),
		appCli.GenTxCmd(app.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		appCli.AddGenesisAccountCmd(app.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		TestnetCmd(),
		appCli.Debug(),
		Network(),
		sdkCfgCmd,
	)

	appCreator := appCreator{encodingConfig}
	server.AddCommands(rootCmd, app.DefaultNodeHome, appCreator.newApp, appCreator.appExport, addStartFlags)

	// add keybase, auxiliary RPC, query, and tx child commands
	rpcStatusCmd := rpc.StatusCommand()
	rpcStatusCmd.SetOut(os.Stdout)
	rootCmd.AddCommand(
		appCli.KeyCommands(app.DefaultNodeHome),
		rpcStatusCmd,
		queryCommand(),
		txCommand(),
	)

	for _, command := range rootCmd.Commands() {
		// tendermint add update validator key command
		if command.Use == "tendermint" {
			command.AddCommand(appCli.UpdateValidatorKeyCmd())
			command.AddCommand(appCli.UpdateNodeKeyCmd())
		}
	}
}

func addStartFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)

	preRun := startCmd.PreRunE
	startCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if err := preRun(cmd, args); err != nil {
			return err
		}
		serverCtx := server.GetServerContextFromCmd(cmd)
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
}

func checkMainnetAndBlock(genesisDoc *tmtypes.GenesisDoc, config *config.Config) error {
	if genesisDoc.InitialHeight != 0 || genesisDoc.ChainID != fxtypes.MainnetChainId {
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
	blockStore := store2.NewBlockStore(blockStoreDB)
	if genesisDoc.GenesisTime.Equal(genesisTime) && blockStore.Height() >= 5713000 {
		return errors.New("invalid version: Please use a later version of fxcored for sync blocks")
	}
	return nil
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
		appCli.QueryStoreCmd(),
		appCli.QueryBlockResultsCmd(),
	)

	app.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		vestingcli.GetTxCmd(),
	)

	app.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

type appCreator struct {
	encCfg app.EncodingConfig
}

// newApp is an AppCreator
func (a appCreator) newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	var cache sdk.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	snapshotDir := filepath.Join(cast.ToString(appOpts.Get(flags.FlagHome)), "data", "snapshots")
	snapshotDB, err := sdk.NewLevelDB("metadata", snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	gasPrice := cast.ToString(appOpts.Get(server.FlagMinGasPrices))
	if strings.Contains(gasPrice, ".") {
		panic("Invalid gas price, cannot contain decimals")
	}
	return app.New(
		logger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		a.encCfg,
		// this line is used by starport scaffolding # stargate/root/appArgument
		appOpts,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(gasPrice),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshotStore(snapshotStore),
		baseapp.SetSnapshotInterval(cast.ToUint64(appOpts.Get(server.FlagStateSyncSnapshotInterval))),
		baseapp.SetSnapshotKeepRecent(cast.ToUint32(appOpts.Get(server.FlagStateSyncSnapshotKeepRecent))),
	)
}

// appExport creates a new simapp (optionally at a given height)
func (a appCreator) appExport(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailAllowedAddrs []string,
	appOpts servertypes.AppOptions) (servertypes.ExportedApp, error) {

	var anApp *app.App

	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	if height != -1 {
		anApp = app.New(
			logger,
			db,
			traceStore,
			false,
			map[int64]bool{},
			homePath,
			uint(1),
			a.encCfg,
			// this line is used by starport scaffolding # stargate/root/exportArgument
			appOpts,
		)

		if err := anApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		anApp = app.New(
			logger,
			db,
			traceStore,
			true,
			map[int64]bool{},
			homePath,
			uint(1),
			a.encCfg,
			// this line is used by starport scaffolding # stargate/root/noHeightExportArgument
			appOpts,
		)
	}

	return anApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)
}

func overwriteFlagDefaults(c *cobra.Command, defaults map[string]string) {
	set := func(s *pflag.FlagSet, key, val string) {
		if f := s.Lookup(key); f != nil {
			f.DefValue = val
			if err := f.Value.Set(val); err != nil {
				panic(err)
			}
		}
	}
	for key, val := range defaults {
		set(c.Flags(), key, val)
		set(c.PersistentFlags(), key, val)
	}
	for _, c := range c.Commands() {
		overwriteFlagDefaults(c, defaults)
	}
}
