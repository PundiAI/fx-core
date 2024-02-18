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
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/snapshots"
	snapshottypes "github.com/cosmos/cosmos-sdk/snapshots/types"
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
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/client/cli"
	fxserver "github.com/functionx/fx-core/v7/server"
	fxcfg "github.com/functionx/fx-core/v7/server/config"
	fxtypes "github.com/functionx/fx-core/v7/types"
	arbitrumcli "github.com/functionx/fx-core/v7/x/arbitrum/client/cli"
	avalanchecli "github.com/functionx/fx-core/v7/x/avalanche/client/cli"
	bsccli "github.com/functionx/fx-core/v7/x/bsc/client/cli"
	crosschaincli "github.com/functionx/fx-core/v7/x/crosschain/client/cli"
	ethcli "github.com/functionx/fx-core/v7/x/eth/client/cli"
	optimismcli "github.com/functionx/fx-core/v7/x/optimism/client/cli"
	polygoncli "github.com/functionx/fx-core/v7/x/polygon/client/cli"
	troncli "github.com/functionx/fx-core/v7/x/tron/client/cli"
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
		WithViper(fxtypes.EnvPrefix).
		WithKeyringOptions(hd.EthSecp256k1Option())

	rootCmd := &cobra.Command{
		Use:   fxtypes.Name + "d",
		Short: "FunctionX Core BlockChain App",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) (err error) {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			// read flag
			initClientCtx, err = client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// read client.toml
			initClientCtx, err = sdkcfg.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			// set clientCtx
			if err = client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := fxcfg.AppConfig(fxtypes.GetDefGasPrice())
			if err = sdkserver.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, fxcfg.DefaultTendermintConfig()); err != nil {
				return err
			}
			return nil
		},
	}

	initRootCmd(rootCmd, encodingConfig, fxtypes.GetDefaultNodeHome())
	return rootCmd
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig app.EncodingConfig, defaultNodeHome string) {
	myAppCreator := appCreator{encodingConfig}
	rootCmd.AddCommand(
		cli.Debug(),
		cli.InitCmd(defaultNodeHome, app.NewDefAppGenesisByDenom(fxtypes.DefaultDenom, encodingConfig.Codec), app.CustomGenesisConsensusParams()),
		cli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, defaultNodeHome),
		cli.GenTxCmd(app.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, defaultNodeHome),
		cli.AddGenesisAccountCmd(defaultNodeHome),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		tmcli.NewCompletionCmd(rootCmd, true),
		testnetCmd(encodingConfig),
		configCmd(),
		pruningCommand(myAppCreator.newApp, defaultNodeHome),
	)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		cli.StatusCommand(),
		keyCommands(defaultNodeHome),
		queryCommand(),
		txCommand(),
		version.NewVersionCommand(),
		sdkserver.NewRollbackCmd(myAppCreator.newApp, defaultNodeHome),
		fxserver.DataCmd(),
		fxserver.ExportSateCmd(myAppCreator.appExport, defaultNodeHome),
		fxserver.StartCmd(myAppCreator.newApp, defaultNodeHome),
		fxserver.TendermintCommand(),
		fxserver.RosettaCommand(encodingConfig.InterfaceRegistry, encodingConfig.Codec),
		preUpgradeCmd(),
		doctorCmd(),
	)
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
		cli.BlockCommand(),
		cli.QueryTxsByEventsCmd(),
		cli.QueryTxCmd(),
		cli.QueryStoreCmd(),
		cli.QueryValidatorByConsAddr(),
		cli.QueryBlockResultsCmd(),
		cli.QueryGasPricesCmd(),
		crosschaincli.GetQueryCmd(
			avalanchecli.GetQueryCmd(),
			bsccli.GetQueryCmd(),
			ethcli.GetQueryCmd(),
			polygoncli.GetQueryCmd(),
			troncli.GetQueryCmd(),
			arbitrumcli.GetQueryCmd(),
			optimismcli.GetQueryCmd(),
		),
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
		crosschaincli.GetTxCmd(
			avalanchecli.GetTxCmd(),
			bsccli.GetTxCmd(),
			ethcli.GetTxCmd(),
			polygoncli.GetTxCmd(),
			troncli.GetTxCmd(),
			arbitrumcli.GetTxCmd(),
			optimismcli.GetTxCmd(),
		),
	)

	app.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func pruningCommand(appCreator servertypes.AppCreator, nodeHome string) *cobra.Command {
	pruningCmd := pruning.PruningCmd(appCreator)
	homeFlag := pruningCmd.Flag(flags.FlagHome)
	homeFlag.DefValue = nodeHome
	if err := homeFlag.Value.Set(nodeHome); err != nil {
		panic(err)
	}
	dbBackend := pruningCmd.Flag(pruning.FlagAppDBBackend)
	dbBackend.DefValue = string(dbm.GoLevelDBBackend)
	if err := dbBackend.Value.Set(string(dbm.GoLevelDBBackend)); err != nil {
		panic(err)
	}
	pruningCmd.Example = fmt.Sprintf(`$ %s prune --pruning custom --pruning-keep-recent 100`, fxtypes.Name)
	return pruningCmd
}

type appCreator struct {
	encCfg app.EncodingConfig
}

// newApp is an AppCreator
func (a appCreator) newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	var cache sdk.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(sdkserver.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(sdkserver.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := sdkserver.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	snapshotDir := filepath.Join(cast.ToString(appOpts.Get(flags.FlagHome)), "data", "snapshots")
	snapshotDB, err := dbm.NewDB("metadata", sdkserver.GetAppDBBackend(appOpts), snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	gasPrice := cast.ToString(appOpts.Get(sdkserver.FlagMinGasPrices))
	if strings.Contains(gasPrice, ".") {
		panic("Invalid gas price, cannot contain decimals")
	}

	snapshotOptions := snapshottypes.NewSnapshotOptions(
		cast.ToUint64(appOpts.Get(sdkserver.FlagStateSyncSnapshotInterval)),
		cast.ToUint32(appOpts.Get(sdkserver.FlagStateSyncSnapshotKeepRecent)),
	)
	return app.New(
		logger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(sdkserver.FlagInvCheckPeriod)),
		a.encCfg,
		appOpts,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(gasPrice),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(sdkserver.FlagMinRetainBlocks))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(sdkserver.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(sdkserver.FlagHaltTime))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(sdkserver.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(sdkserver.FlagIndexEvents))),
		baseapp.SetSnapshot(snapshotStore, snapshotOptions),
		baseapp.SetIAVLCacheSize(cast.ToInt(appOpts.Get(sdkserver.FlagIAVLCacheSize))),
		baseapp.SetIAVLDisableFastNode(cast.ToBool(appOpts.Get(sdkserver.FlagDisableIAVLFastNode))),
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
			homePath, cast.ToUint(appOpts.Get(sdkserver.FlagInvCheckPeriod)), a.encCfg, appOpts,
		)

		if err := anApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		anApp = app.New(logger, db, traceStore, true, map[int64]bool{},
			homePath, cast.ToUint(appOpts.Get(sdkserver.FlagInvCheckPeriod)), a.encCfg, appOpts,
		)
	}

	return anApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)
}
