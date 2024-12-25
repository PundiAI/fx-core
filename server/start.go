package server

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"

	cmtdbm "github.com/cometbft/cometbft-db"
	abciserver "github.com/cometbft/cometbft/abci/server"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/crypto/tmhash"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	pvm "github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	"github.com/cometbft/cometbft/rpc/client/local"
	"github.com/cometbft/cometbft/store"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servergrpc "github.com/cosmos/cosmos-sdk/server/grpc"
	servercmtlog "github.com/cosmos/cosmos-sdk/server/log"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	ethmetricsexp "github.com/ethereum/go-ethereum/metrics/exp"
	"github.com/evmos/ethermint/indexer"
	ethdebug "github.com/evmos/ethermint/rpc/namespaces/ethereum/debug"
	ethermintserver "github.com/evmos/ethermint/server"
	ethermintconfig "github.com/evmos/ethermint/server/config"
	srvflags "github.com/evmos/ethermint/server/flags"
	ethermint "github.com/evmos/ethermint/types"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	fxtypes "github.com/pundiai/fx-core/v8/types"
)

// StartCmd runs the service passed in, either stand-alone or in-process with
// CometBFT.
//

func StartCmd(appCreator types.AppCreator, defaultNodeHome string) *cobra.Command {
	opts := ethermintserver.NewDefaultStartOptions(appCreator, defaultNodeHome)
	startCmd := ethermintserver.StartCmd(opts)
	startCmd.PreRunE = func(cmd *cobra.Command, _ []string) error {
		serverCtx := server.GetServerContextFromCmd(cmd)

		if _, err := server.GetPruningOptionsFromFlags(serverCtx.Viper); err != nil {
			return err
		}

		genDocFile := serverCtx.Config.GenesisFile()
		genesisBytes, err := os.ReadFile(genDocFile)
		if err != nil {
			return fmt.Errorf("couldn't read GenesisDoc file: %w", err)
		}
		expectGenesisHash := serverCtx.Viper.GetString("genesis_hash")
		actualGenesisHash := hex.EncodeToString(tmhash.Sum(genesisBytes))
		if len(expectGenesisHash) != 0 && actualGenesisHash != expectGenesisHash {
			return fmt.Errorf("--genesis_hash=%s does not match %s hash: %s", expectGenesisHash, genDocFile, actualGenesisHash)
		}

		appGenesis, err := genutiltypes.AppGenesisFromFile(genDocFile)
		if err != nil {
			return err
		}
		if err = checkMainnetAndBlock(appGenesis, actualGenesisHash, serverCtx.Config); err != nil {
			return err
		}

		clientCtx := client.GetClientContextFromCmd(cmd)
		if len(clientCtx.ChainID) == 0 {
			clientCtx.ChainID = appGenesis.ChainID
		}
		if len(clientCtx.HomeDir) == 0 {
			clientCtx.HomeDir = serverCtx.Config.RootDir
		}
		if err = client.SetCmdClientContext(cmd, clientCtx); err != nil {
			return err
		}
		serverCtx.Viper.Set(flags.FlagChainID, clientCtx.ChainID)
		return server.SetCmdServerContext(cmd, serverCtx)
	}
	startCmd.RunE = func(cmd *cobra.Command, args []string) error {
		serverCtx := server.GetServerContextFromCmd(cmd)
		clientCtx, err := client.GetClientQueryContext(cmd)
		if err != nil {
			return err
		}

		withTM, _ := cmd.Flags().GetBool(srvflags.WithCometBFT)
		if !withTM {
			serverCtx.Logger.Info("starting ABCI without CometBFT")
			return wrapCPUProfile(serverCtx, func() error {
				return startStandAlone(serverCtx, opts)
			})
		}

		serverCtx.Logger.Info("Unlocking keyring")

		// fire unlock precess for keyring
		krBackend := clientCtx.Keyring.Backend()
		if krBackend == keyring.BackendFile {
			_, err = clientCtx.Keyring.List()
			if err != nil {
				return err
			}
		}

		// set keyring backend type to the server context
		serverCtx.Viper.Set(flags.FlagKeyringBackend, krBackend)

		serverCtx.Logger.Info("starting ABCI with CometBFT")

		// amino is needed here for backwards compatibility of REST routes
		err = wrapCPUProfile(serverCtx, func() error {
			return startInProcess(serverCtx, clientCtx, opts)
		})
		if err != nil {
			return err
		}

		return nil
	}
	crisis.AddModuleInitFlags(startCmd)
	return startCmd
}

func checkMainnetAndBlock(genesisDoc *genutiltypes.AppGenesis, genesisHash string, config *cmtcfg.Config) error {
	if genesisDoc.InitialHeight > 1 || genesisDoc.ChainID != fxtypes.MainnetChainId || config.StateSync.Enable {
		return nil
	}
	genesisTime, err := time.Parse("2006-01-02T15:04:05Z", "2021-07-05T04:00:00Z")
	if err != nil {
		return err
	}
	blockStoreDB, err := cmtdbm.NewDB("blockstore", cmtdbm.BackendType(config.DBBackend), config.DBDir())
	if err != nil {
		return err
	}
	defer func() {
		_ = blockStoreDB.Close()
	}()
	blockStore := store.NewBlockStore(blockStoreDB)
	if genesisDoc.GenesisTime.Equal(genesisTime) {
		if genesisHash != fxtypes.MainnetGenesisHash {
			return nil
		}
		if blockStore.Height() < fxtypes.MainnetBlockHeightV2 {
			return errors.New("invalid version: The current block height is less than the fxv2 upgrade height(5_713_000), " +
				"sync block from scratch please use use fxcored v1.x.x")
		}
		if blockStore.Height() < fxtypes.MainnetBlockHeightV3 {
			return errors.New("invalid version: The current block height is less than the fxv3 upgrade height(8_756_000)," +
				" please use the v2.x.x version to synchronize the block or download the latest snapshot")
		}
		if blockStore.Height() < fxtypes.MainnetBlockHeightV4 {
			return errors.New("invalid version: The current block height is less than the v4.2.0 upgrade height(10_477_500)," +
				" please use the v3.x.x version to synchronize the block or download the latest snapshot")
		}
		if blockStore.Height() < fxtypes.MainnetBlockHeightV5 {
			return errors.New("invalid version: The current block height is less than the v5.0.0 upgrade height(11_601_700)," +
				" please use the v4.x.x version to synchronize the block or download the latest snapshot")
		}
		if blockStore.Height() < fxtypes.MainnetBlockHeightV6 {
			return errors.New("invalid version: The current block height is less than the v6.0.0 upgrade height(13_598_000)," +
				" please use the v5.x.x version to synchronize the block or download the latest snapshot")
		}
		if blockStore.Height() < fxtypes.MainnetBlockHeightV7 {
			return errors.New("invalid version: The current block height is less than the v7.5.0 upgrade height(16_838_000)," +
				" please use the v6.x.x version to synchronize the block or download the latest snapshot")
		}
		// TODO: The line of code below must be removed before the release.
		return errors.New("invalid version: The current version is not released, please use the corresponding version")
	}
	return nil
}

func startStandAlone(svrCtx *server.Context, opts ethermintserver.StartOptions) error {
	addr := svrCtx.Viper.GetString(srvflags.Address)
	transport := svrCtx.Viper.GetString(srvflags.Transport)
	home := svrCtx.Viper.GetString(flags.FlagHome)

	db, err := opts.DBOpener(svrCtx.Viper, home, server.GetAppDBBackend(svrCtx.Viper))
	if err != nil {
		return err
	}

	traceWriterFile := svrCtx.Viper.GetString(srvflags.TraceStore)
	traceWriter, err := openTraceWriter(traceWriterFile)
	if err != nil {
		return err
	}

	app := opts.AppCreator(svrCtx.Logger, db, traceWriter, svrCtx.Viper)
	defer func() {
		if err = app.Close(); err != nil {
			svrCtx.Logger.Error("close application failed", "error", err.Error())
		}
	}()

	config, err := ethermintconfig.GetConfig(svrCtx.Viper)
	if err != nil {
		svrCtx.Logger.Error("failed to get server config", "error", err.Error())
		return err
	}

	if err = config.ValidateBasic(); err != nil {
		svrCtx.Logger.Error("invalid server config", "error", err.Error())
		return err
	}

	_, err = startTelemetry(config)
	if err != nil {
		return err
	}

	cmtApp := server.NewCometABCIWrapper(app)
	svr, err := abciserver.NewServer(addr, transport, cmtApp)
	if err != nil {
		return fmt.Errorf("error creating listener: %w", err)
	}

	svr.SetLogger(servercmtlog.CometLoggerWrapper{Logger: svrCtx.Logger.With("server", "abci")})
	g, ctx := getCtx(svrCtx, false)

	g.Go(func() error {
		if err = svr.Start(); err != nil {
			svrCtx.Logger.Error("failed to start out-of-process ABCI server", "err", err)
			return err
		}

		// Wait for the calling process to be canceled or close the provided context,
		// so we can gracefully stop the ABCI server.
		<-ctx.Done()
		svrCtx.Logger.Info("stopping the ABCI server...")
		return svr.Stop()
	})

	return g.Wait()
}

//nolint:gocyclo // need to refactor
func startInProcess(svrCtx *server.Context, clientCtx client.Context, opts ethermintserver.StartOptions) (err error) {
	cfg := svrCtx.Config
	home := cfg.RootDir
	logger := svrCtx.Logger

	g, ctx := getCtx(svrCtx, true)

	db, err := opts.DBOpener(svrCtx.Viper, home, server.GetAppDBBackend(svrCtx.Viper))
	if err != nil {
		logger.Error("failed to open DB", "error", err.Error())
		return err
	}

	traceWriterFile := svrCtx.Viper.GetString(srvflags.TraceStore)
	traceWriter, err := openTraceWriter(traceWriterFile)
	if err != nil {
		logger.Error("failed to open trace writer", "error", err.Error())
		return err
	}

	config, err := ethermintconfig.GetConfig(svrCtx.Viper)
	if err != nil {
		logger.Error("failed to get server config", "error", err.Error())
		return err
	}

	if err = config.ValidateBasic(); err != nil {
		logger.Error("invalid server config", "error", err.Error())
		return err
	}

	app := opts.AppCreator(svrCtx.Logger, db, traceWriter, svrCtx.Viper)
	defer func() {
		if err = app.Close(); err != nil {
			logger.Error("close application failed", "error", err.Error())
		}
	}()

	var tmNode *node.Node
	gRPCOnly := svrCtx.Viper.GetBool(srvflags.GRPCOnly)

	if gRPCOnly {
		logger.Info("starting node in query only mode; CometBFT is disabled")
		config.GRPC.Enable = true
		config.JSONRPC.EnableIndexer = false
	} else {
		logger.Info("starting node with ABCI CometBFT in-process")

		var nodeKey *p2p.NodeKey
		nodeKey, err = p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
		if err != nil {
			logger.Error("failed load or gen node key", "error", err.Error())
			return err
		}

		cmtApp := server.NewCometABCIWrapper(app)
		tmNode, err = node.NewNodeWithContext(
			ctx,
			cfg,
			pvm.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
			nodeKey,
			proxy.NewLocalClientCreator(cmtApp),
			ethermintserver.GenDocProvider(cfg),
			cmtcfg.DefaultDBProvider,
			node.DefaultMetricsProvider(cfg.Instrumentation),
			servercmtlog.CometLoggerWrapper{Logger: svrCtx.Logger.With("server", "node")},
		)
		if err != nil {
			logger.Error("failed init node", "error", err.Error())
			return err
		}

		if err = tmNode.Start(); err != nil {
			logger.Error("failed start tendermint server", "error", err.Error())
			return err
		}

		defer func() {
			if tmNode.IsRunning() {
				_ = tmNode.Stop()
			}
		}()
	}

	// Add the tx service to the gRPC router. We only need to register this
	// service if API or gRPC or JSONRPC is enabled, and avoid doing so in the general
	// case, because it spawns a new local tendermint RPC client.
	if (config.API.Enable || config.GRPC.Enable || config.JSONRPC.Enable || config.JSONRPC.EnableIndexer) && tmNode != nil {
		clientCtx = clientCtx.WithClient(local.New(tmNode))

		app.RegisterTxService(clientCtx)
		app.RegisterTendermintService(clientCtx)
		app.RegisterNodeService(clientCtx, config.Config)
	}

	metrics, err := startTelemetry(config)
	if err != nil {
		return err
	}

	// Enable metrics if JSONRPC is enabled and --metrics is passed
	// Flag not added in config to avoid user enabling in config without passing in CLI
	if config.JSONRPC.Enable && svrCtx.Viper.GetBool(srvflags.JSONRPCEnableMetrics) {
		ethmetricsexp.Setup(config.JSONRPC.MetricsAddress)
	}

	var idxer ethermint.EVMTxIndexer
	if config.JSONRPC.EnableIndexer {
		idxDB, err := ethermintserver.OpenIndexerDB(home, server.GetAppDBBackend(svrCtx.Viper))
		if err != nil {
			logger.Error("failed to open evm indexer DB", "error", err.Error())
			return err
		}

		idxLogger := logger.With("indexer", "evm")
		idxer = indexer.NewKVIndexer(idxDB, idxLogger, clientCtx)
		indexerService := ethermintserver.NewEVMIndexerService(idxer, clientCtx.Client.(rpcclient.Client), config.JSONRPC.AllowIndexerGap)
		indexerService.SetLogger(servercmtlog.CometLoggerWrapper{Logger: idxLogger})

		g.Go(func() error {
			return indexerService.Start()
		})
		g.Go(func() error {
			<-ctx.Done()
			return indexerService.Stop()
		})
	}

	if config.API.Enable || config.JSONRPC.Enable {
		chainID := svrCtx.Viper.GetString(flags.FlagChainID)
		clientCtx = clientCtx.
			WithHomeDir(home).
			WithChainID(chainID)
	}

	grpcSrv, clientCtx, err := startGrpcServer(ctx, svrCtx, clientCtx, g, config.GRPC, app)
	if err != nil {
		return err
	}
	if grpcSrv != nil {
		defer grpcSrv.GracefulStop()
	}

	apiSrv := startAPIServer(ctx, svrCtx, clientCtx, g, config.Config, app, grpcSrv, metrics)
	if apiSrv != nil {
		defer func() {
			logger.Info("Closing API server", "err", apiSrv.Close())
		}()
	}

	_, err = startJSONRPCServer(ctx, svrCtx, clientCtx, g, config, idxer, app)
	if err != nil {
		return err
	}

	return g.Wait()
}

func wrapCPUProfile(ctx *server.Context, callback func() error) error {
	if cpuProfile := ctx.Viper.GetString(srvflags.CPUProfile); cpuProfile != "" {
		fp, err := ethdebug.ExpandHome(cpuProfile)
		if err != nil {
			ctx.Logger.Debug("failed to get filepath for the CPU profile file", "error", err.Error())
			return err
		}
		f, err := os.Create(fp)
		if err != nil {
			return err
		}

		ctx.Logger.Info("starting CPU profiler", "profile", cpuProfile)
		if err = pprof.StartCPUProfile(f); err != nil {
			return err
		}

		defer func() {
			ctx.Logger.Info("stopping CPU profiler", "profile", cpuProfile)
			pprof.StopCPUProfile()
			if err = f.Close(); err != nil {
				ctx.Logger.Info("failed to close cpu-profile file", "profile", cpuProfile, "err", err.Error())
			}
		}()
	}

	return callback()
}

func getCtx(svrCtx *server.Context, block bool) (*errgroup.Group, context.Context) {
	ctx, cancelFn := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	// listen for quit signals so the calling parent process can gracefully exit
	server.ListenForQuitSignals(g, block, cancelFn, svrCtx.Logger)
	return g, ctx
}

func openTraceWriter(traceWriterFile string) (w io.Writer, err error) {
	if traceWriterFile == "" {
		return
	}

	filePath := filepath.Clean(traceWriterFile)
	return os.OpenFile(
		filePath,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0o600,
	)
}

func startTelemetry(cfg ethermintconfig.Config) (*telemetry.Metrics, error) {
	if !cfg.Telemetry.Enabled {
		return nil, nil
	}
	return telemetry.New(cfg.Telemetry)
}

func startGrpcServer(
	ctx context.Context,
	svrCtx *server.Context,
	clientCtx client.Context,
	g *errgroup.Group,
	config serverconfig.GRPCConfig,
	app types.Application,
) (*grpc.Server, client.Context, error) {
	if !config.Enable {
		// return grpcServer as nil if gRPC is disabled
		return nil, clientCtx, nil
	}
	_, _, err := net.SplitHostPort(config.Address)
	if err != nil {
		return nil, clientCtx, fmt.Errorf("invalid grpc address %s, err: %s", config.Address, err.Error())
	}

	maxSendMsgSize := config.MaxSendMsgSize
	if maxSendMsgSize == 0 {
		maxSendMsgSize = serverconfig.DefaultGRPCMaxSendMsgSize
	}

	maxRecvMsgSize := config.MaxRecvMsgSize
	if maxRecvMsgSize == 0 {
		maxRecvMsgSize = serverconfig.DefaultGRPCMaxRecvMsgSize
	}

	// if gRPC is enabled, configure gRPC client for gRPC gateway and json-rpc
	grpcClient, err := grpc.NewClient(
		config.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.ForceCodec(codec.NewProtoCodec(clientCtx.InterfaceRegistry).GRPCCodec()),
			grpc.MaxCallRecvMsgSize(maxRecvMsgSize),
			grpc.MaxCallSendMsgSize(maxSendMsgSize),
		),
	)
	if err != nil {
		return nil, clientCtx, err
	}
	// Set `GRPCClient` to `clientCtx` to enjoy concurrent grpc query.
	// only use it if gRPC server is enabled.
	clientCtx = clientCtx.WithGRPCClient(grpcClient)
	svrCtx.Logger.Debug("gRPC client assigned to client context", "address", config.Address)

	grpcSrv, err := servergrpc.NewGRPCServer(clientCtx, app, config)
	if err != nil {
		return nil, clientCtx, err
	}

	// Start the gRPC server in a goroutine. Note, the provided ctx will ensure
	// that the server is gracefully shut down.
	g.Go(func() error {
		return servergrpc.StartGRPCServer(ctx, svrCtx.Logger.With("module", "grpc-server"), config, grpcSrv)
	})
	return grpcSrv, clientCtx, nil
}

func startAPIServer(
	ctx context.Context,
	svrCtx *server.Context,
	clientCtx client.Context,
	g *errgroup.Group,
	svrCfg serverconfig.Config,
	app types.Application,
	grpcSrv *grpc.Server,
	metrics *telemetry.Metrics,
) *api.Server {
	if !svrCfg.API.Enable {
		return nil
	}

	apiSrv := api.New(clientCtx, svrCtx.Logger.With("server", "api"), grpcSrv)
	app.RegisterAPIRoutes(apiSrv, svrCfg.API)

	if svrCfg.Telemetry.Enabled {
		apiSrv.SetTelemetry(metrics)
	}

	g.Go(func() error {
		return apiSrv.Start(ctx, svrCfg)
	})
	return apiSrv
}

func startJSONRPCServer(
	ctx context.Context,
	svrCtx *server.Context,
	clientCtx client.Context,
	g *errgroup.Group,
	config ethermintconfig.Config,
	idxer ethermint.EVMTxIndexer,
	app types.Application,
) (httpSrv *http.Server, err error) {
	if !config.JSONRPC.Enable {
		return
	}

	txApp, ok := app.(ethermintserver.AppWithPendingTxStream)
	if !ok {
		return httpSrv, fmt.Errorf("json-rpc server requires AppWithPendingTxStream")
	}

	clientCtx = clientCtx.WithChainID(fxtypes.ChainIdWithEIP155(clientCtx.ChainID))
	return StartJSONRPC(ctx, svrCtx, clientCtx, g, &config, idxer, txApp)
}
