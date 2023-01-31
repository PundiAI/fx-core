package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdkcfg "github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/snapshots"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/evmos/ethermint/crypto/hd"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/client/cli"
	fxserver "github.com/functionx/fx-core/v3/server"
	fxcfg "github.com/functionx/fx-core/v3/server/config"
	"github.com/functionx/fx-core/v3/server/grpc/base/gasprice"
	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschaincli "github.com/functionx/fx-core/v3/x/crosschain/client/cli"
)

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() *cobra.Command {
	fxtypes.SetConfig(false)

	encodingConfig := app.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(fxtypes.GetDefaultNodeHome()).
		WithViper("FX").
		WithKeyringOptions(hd.EthSecp256k1Option())

	rootCmd := &cobra.Command{
		Use:   fxtypes.Name + "d",
		Short: "FunctionX Core BlockChain App",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			// read flag
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// read client.toml
			initClientCtx, err = sdkcfg.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			// set clientCtx
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := fxcfg.AppConfig(fmt.Sprintf("4000000000000%s", fxtypes.DefaultDenom))
			if err := server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig); err != nil {
				return err
			}
			return nil
		},
	}

	initRootCmd(rootCmd, encodingConfig, fxtypes.GetDefaultNodeHome())
	overwriteFlagDefaults(rootCmd, map[string]string{
		flags.FlagChainID:        fxtypes.ChainId(),
		flags.FlagKeyringBackend: keyring.BackendOS,
		flags.FlagGas:            "100000",
		// flags.FlagGas:            "auto",
		// flags.FlagGasAdjustment:  "1.5",
	})
	return rootCmd
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig app.EncodingConfig, defaultNodeHome string) {
	rootCmd.AddCommand(
		cli.InitCmd(defaultNodeHome, app.NewDefAppGenesisByDenom(fxtypes.DefaultDenom, encodingConfig.Codec), app.CustomConsensusParams()),
		cli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, defaultNodeHome),
		cli.GenTxCmd(app.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, defaultNodeHome),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		cli.AddGenesisAccountCmd(defaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		testnetCmd(),
		cli.Debug(),
		configCmd(),
	)

	appCreator := appCreator{encodingConfig}

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		keyCommands(defaultNodeHome),
		cli.StatusCommand(),
		queryCommand(),
		txCommand(),
		version.NewVersionCommand(),
		server.NewRollbackCmd(appCreator.newApp, defaultNodeHome),
		fxserver.DataCmd(),
		fxserver.ExportSateCmd(appCreator.appExport, defaultNodeHome),
		fxserver.StartCmd(appCreator.newApp, defaultNodeHome),
		tendermintCommand(),
		app.GetUpgrades().GetLatest().PreUpgradeCmd,
	)

	// add rosetta
	rootCmd.AddCommand(server.RosettaCommand(encodingConfig.InterfaceRegistry, encodingConfig.Codec))
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
		cli.ValidatorCommand(),
		cli.BlockCommand(),
		cli.QueryTxsByEventsCmd(),
		cli.QueryTxCmd(),
		cli.QueryStoreCmd(),
		cli.QueryValidatorByConsAddr(),
		cli.QueryBlockResultsCmd(),
		gasprice.QueryCmd(),
		crosschaincli.GetQueryCmd(),
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
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		crosschaincli.GetTxCmd(),
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
		baseapp.SetIAVLCacheSize(cast.ToInt(appOpts.Get(server.FlagIAVLCacheSize))),
		baseapp.SetIAVLDisableFastNode(cast.ToBool(appOpts.Get(server.FlagDisableIAVLFastNode))),
	)
}

// appExport creates a new simapp (optionally at a given height)
func (a appCreator) appExport(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
) (servertypes.ExportedApp, error) {
	var anApp *app.App
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	if height != -1 {
		anApp = app.New(logger, db, traceStore, false, map[int64]bool{},
			homePath, cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)), a.encCfg, appOpts,
		)

		if err := anApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		anApp = app.New(logger, db, traceStore, true, map[int64]bool{},
			homePath, cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)), a.encCfg, appOpts,
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
			if key == flags.FlagGasPrices {
				f.Usage = "Gas prices in decimal format to determine the transaction fee"
			}
			if key == flags.FlagGas {
				f.Usage = "gas limit to set per-transaction; set to 'auto' to calculate sufficient gas automatically"
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
