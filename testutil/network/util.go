package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"

	"cosmossdk.io/log"
	cmtcfg "github.com/cometbft/cometbft/config"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	pvm "github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	"github.com/cometbft/cometbft/rpc/client/local"
	"github.com/cometbft/cometbft/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	cosmosserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	servergrpc "github.com/cosmos/cosmos-sdk/server/grpc"
	servercmtlog "github.com/cosmos/cosmos-sdk/server/log"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethermintserver "github.com/evmos/ethermint/server"
	"golang.org/x/sync/errgroup"

	"github.com/functionx/fx-core/v8/app"
	"github.com/functionx/fx-core/v8/server"
	fxtypes "github.com/functionx/fx-core/v8/types"
)

// StartInProcess creates and starts an in-process local test network.
//
//gocyclo:ignore
func StartInProcess(ctx context.Context, errGroup *errgroup.Group, appConstructor AppConstructor, val *Validator) error {
	logger := val.Ctx.Logger
	tmCfg := val.Ctx.Config

	if err := val.AppConfig.ValidateBasic(); err != nil {
		return err
	}

	nodeKey, err := p2p.LoadOrGenNodeKey(tmCfg.NodeKeyFile())
	if err != nil {
		return err
	}

	myApp := appConstructor(val.AppConfig, val.Ctx)

	cmtApp := cosmosserver.NewCometABCIWrapper(myApp)
	appGenesisProvider := func() (*types.GenesisDoc, error) {
		appGenesis, err := genutiltypes.AppGenesisFromFile(tmCfg.GenesisFile())
		if err != nil {
			return nil, err
		}

		return appGenesis.ToGenesisDoc()
	}
	tmNode, err := node.NewNodeWithContext(
		ctx,
		tmCfg,
		pvm.LoadOrGenFilePV(tmCfg.PrivValidatorKeyFile(), tmCfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(cmtApp),
		appGenesisProvider,
		cmtcfg.DefaultDBProvider,
		node.DefaultMetricsProvider(tmCfg.Instrumentation),
		servercmtlog.CometLoggerWrapper{Logger: logger.With("module", val.Ctx.Config.Moniker)},
	)
	if err != nil {
		return err
	}

	if err = tmNode.Start(); err != nil {
		return err
	}

	val.tmNode = tmNode
	val.RPCClient = local.New(tmNode)

	// We'll need a RPC client if the validator exposes a gRPC or REST endpoint.
	if val.AppConfig.API.Address != "" || val.AppConfig.GRPC.Enable {
		val.ClientCtx = val.ClientCtx.WithClient(val.RPCClient)

		myApp.RegisterTxService(val.ClientCtx)
		myApp.RegisterTendermintService(val.ClientCtx)
		myApp.RegisterNodeService(val.ClientCtx, val.AppConfig.Config)
	}

	if val.AppConfig.GRPC.Enable {
		val.grpc, err = servergrpc.NewGRPCServer(val.ClientCtx, myApp, val.AppConfig.GRPC)
		if err != nil {
			return err
		}
		errGroup.Go(func() error {
			return servergrpc.StartGRPCServer(context.Background(), logger.With(log.ModuleKey, "grpc-server"), val.AppConfig.GRPC, val.grpc)
		})
	}

	if val.AppConfig.API.Enable && val.AppConfig.API.Address != "" {
		val.api = api.New(val.ClientCtx, logger.With("module", "api-server"), val.grpc)
		myApp.RegisterAPIRoutes(val.api, val.AppConfig.API)

		errGroup.Go(func() error {
			return val.api.Start(context.Background(), val.AppConfig.Config)
		})
	}

	if val.AppConfig.JSONRPC.Enable && val.AppConfig.JSONRPC.Address != "" {
		if val.Ctx == nil || val.Ctx.Viper == nil {
			return fmt.Errorf("validator %s context is nil", val.Ctx.Config.Moniker)
		}

		val.ClientCtx = val.ClientCtx.WithChainID(fxtypes.ChainIdWithEIP155(val.ClientCtx.ChainID))
		val.jsonrpc, err = server.StartJSONRPC(ctx, val.Ctx, val.ClientCtx, errGroup, val.AppConfig.ToEthermintConfig(), nil, myApp.(ethermintserver.AppWithPendingTxStream))
		if err != nil {
			return err
		}

		address := fmt.Sprintf("http://%s", val.AppConfig.JSONRPC.Address)

		val.JSONRPCClient, err = ethclient.Dial(address)
		if err != nil {
			return fmt.Errorf("failed to dial JSON-RPC at %s: %w", val.AppConfig.JSONRPC.Address, err)
		}
	}
	return nil
}

func collectGenFiles(cfg Config, vals []*Validator, outputDir string) error {
	genTime := cmttime.Now()
	for i := 0; i < cfg.NumValidators; i++ {
		tmCfg := vals[i].Ctx.Config

		gentxsDir := filepath.Join(outputDir, "gentxs")
		initCfg := genutiltypes.NewInitConfig(cfg.ChainID, gentxsDir, vals[i].NodeID, vals[i].PubKey)

		genFile := tmCfg.GenesisFile()
		appGenesis, err := genutiltypes.AppGenesisFromFile(genFile)
		if err != nil {
			return err
		}
		appState, err := genutil.GenAppStateFromConfig(cfg.Codec, cfg.TxConfig,
			tmCfg, initCfg, appGenesis, banktypes.GenesisBalancesIterator{}, genutiltypes.DefaultMessageValidator,
			cfg.TxConfig.SigningContext().ValidatorAddressCodec())
		if err != nil {
			return err
		}
		appGenesis.GenesisTime = genTime
		appGenesis.AppState = appState
		if err = appGenesis.SaveAs(genFile); err != nil {
			return err
		}
	}
	return nil
}

func initGenFiles(cfg Config, genAccounts []authtypes.GenesisAccount, genBalances []banktypes.Balance, genFiles []string) error {
	// set the accounts in the genesis state
	var authGenState authtypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[authtypes.ModuleName], &authGenState)

	accounts, err := authtypes.PackAccounts(genAccounts)
	if err != nil {
		return err
	}

	authGenState.Accounts = append(authGenState.Accounts, accounts...)
	cfg.GenesisState[authtypes.ModuleName] = cfg.Codec.MustMarshalJSON(&authGenState)

	// set the balances in the genesis state
	var bankGenState banktypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[banktypes.ModuleName], &bankGenState)
	bankGenState.Balances = append(bankGenState.Balances, genBalances...)
	cfg.GenesisState[banktypes.ModuleName] = cfg.Codec.MustMarshalJSON(&bankGenState)

	appGenStateJSON, err := json.MarshalIndent(cfg.GenesisState, "", "  ")
	if err != nil {
		return err
	}

	genDoc := genutiltypes.AppGenesis{
		AppName:       fxtypes.Name,
		ChainID:       cfg.ChainID,
		InitialHeight: 1,
		Consensus: &genutiltypes.ConsensusGenesis{
			Validators: nil,
			Params:     app.CustomGenesisConsensusParams(),
		},
		AppState: appGenStateJSON,
	}

	// generate empty genesis files for each validator and save
	for i := 0; i < cfg.NumValidators; i++ {
		if err = genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}

	return nil
}

func writeFile(name string, dir string, contents []byte) error {
	if err := tmos.EnsureDir(dir, 0o755); err != nil {
		return err
	}

	file := filepath.Join(dir, name)
	return tmos.WriteFile(file, contents, 0o644)
}

func FreeTCPAddr() (addr, port string, err error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", "", err
	}

	portI := l.Addr().(*net.TCPAddr).Port
	port = fmt.Sprintf("%d", portI)
	addr = fmt.Sprintf("tcp://0.0.0.0:%s", port)
	return
}
