package cmd

import (
	"errors"
	fxserver "github.com/functionx/fx-core/server"
	fxtypes "github.com/functionx/fx-core/types"

	"github.com/functionx/fx-core/crypto/hd"
	"github.com/functionx/fx-core/server/config"
	"io"
	"os"
	"path/filepath"
	"strings"

	sdkCfg "github.com/cosmos/cosmos-sdk/client/config"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/snapshots"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingcli "github.com/cosmos/cosmos-sdk/x/auth/vesting/client/cli"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/app"
	appCmd "github.com/functionx/fx-core/app/cmd"
	"github.com/functionx/fx-core/app/fxcore"
	// this line is u by starport scaffolding # stargate/root/import
)

const envPrefix = "FX"

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() *cobra.Command {

	encodingConfig := fxcore.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(fxcore.DefaultNodeHome).
		WithViper(envPrefix).
		WithKeyringOptions(hd.EthSecp256k1Option())

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

			customAppTemplate, customAppConfig := config.AppConfig(fxtypes.MintDenom)
			if err := server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig); err != nil {
				return err
			}
			// add log filter
			return app.AddCmdLogWrapFilterLogType(cmd)
		},
	}

	rootCmd.PersistentFlags().StringSlice(app.FlagLogFilter, []string{}, `The logging filter can discard custom log type (ABCIQuery)`)
	initRootCmd(rootCmd, encodingConfig)
	overwriteFlagDefaults(rootCmd, map[string]string{
		flags.FlagChainID:        fxtypes.ChainID,
		flags.FlagKeyringBackend: keyring.BackendTest,
		flags.FlagGasPrices:      "4000000000000" + fxtypes.MintDenom,
	})
	for _, command := range rootCmd.Commands() {
		if command.Use == "" {
			rootCmd.RemoveCommand(command)
		}
	}
	return rootCmd
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig fxcore.EncodingConfig) {
	sdkCfgCmd := sdkCfg.Cmd()
	sdkCfgCmd.AddCommand(appCmd.AppTomlCmd(), appCmd.ConfigTomlCmd())

	rootCmd.AddCommand(
		InitCmd(),
		appCmd.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, fxcore.DefaultNodeHome),
		genutilcli.MigrateGenesisCmd(),
		appCmd.GenTxCmd(fxcore.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, fxcore.DefaultNodeHome),
		genutilcli.ValidateGenesisCmd(fxcore.ModuleBasics),
		appCmd.AddGenesisAccountCmd(fxcore.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		TestnetCmd(),
		appCmd.Debug(),
		// this line is used by starport scaffolding # stargate/root/commands
		appCmd.Network(),
		sdkCfgCmd,
	)

	appCreator := appCreator{encodingConfig}
	fxserver.AddCommands(rootCmd, fxcore.DefaultNodeHome, appCreator.newApp, appCreator.appExport, func(startCmd *cobra.Command) {})

	// add keybase, auxiliary RPC, query, and tx child commands
	rpcStatusCmd := rpc.StatusCommand()
	rpcStatusCmd.SetOut(os.Stdout)
	rootCmd.AddCommand(
		appCmd.KeyCommands(fxcore.DefaultNodeHome),
		rpcStatusCmd,
		queryCommand(),
		txCommand(),
	)

	for _, command := range rootCmd.Commands() {
		// tendermint add update validator key command
		if command.Use == "tendermint" {
			command.AddCommand(appCmd.UpdateValidatorKeyCmd())
			command.AddCommand(appCmd.UpdateNodeKeyCmd())
		}
	}
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
		appCmd.QueryStoreCmd(),
		appCmd.QueryValidatorByConsAddr(),
		appCmd.QueryBlockResultsCmd(),
	)

	fxcore.ModuleBasics.AddQueryCommands(cmd)
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

	fxcore.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

type appCreator struct {
	encCfg fxcore.EncodingConfig
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
	return fxcore.New(
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

	var anApp *fxcore.App

	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	if height != -1 {
		anApp = fxcore.New(
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
		anApp = fxcore.New(
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
