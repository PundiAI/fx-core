package network

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	pvm "github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	"github.com/cometbft/cometbft/rpc/client/local"
	"github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server/api"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	servergrpc "github.com/cosmos/cosmos-sdk/server/grpc"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/evmos/ethermint/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/functionx/fx-core/v7/app"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

//gocyclo:ignore
func startInProcess(appConstructor AppConstructor, val *Validator) error {
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

	genDocProvider := node.DefaultGenesisDocProviderFunc(tmCfg)
	tmNode, err := node.NewNode(
		tmCfg,
		pvm.LoadOrGenFilePV(tmCfg.PrivValidatorKeyFile(), tmCfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(myApp),
		genDocProvider,
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(tmCfg.Instrumentation),
		logger.With("module", val.Ctx.Config.Moniker),
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
	if val.APIAddress != "" || val.AppConfig.GRPC.Enable {
		val.ClientCtx = val.ClientCtx.WithClient(val.RPCClient)

		// Add the tx service in the gRPC router.
		myApp.RegisterTxService(val.ClientCtx)

		// Add the tendermint queries service in the gRPC router.
		myApp.RegisterTendermintService(val.ClientCtx)
		myApp.RegisterNodeService(val.ClientCtx)
	}

	if val.AppConfig.GRPC.Enable {
		maxSendMsgSize := val.AppConfig.GRPC.MaxSendMsgSize
		if maxSendMsgSize == 0 {
			maxSendMsgSize = srvconfig.DefaultGRPCMaxSendMsgSize
		}

		maxRecvMsgSize := val.AppConfig.GRPC.MaxRecvMsgSize
		if maxRecvMsgSize == 0 {
			maxRecvMsgSize = srvconfig.DefaultGRPCMaxRecvMsgSize
		}

		// If grpc is enabled, configure grpc client for grpc gateway.
		grpcClient, err := grpc.Dial(
			val.AppConfig.GRPC.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.ForceCodec(codec.NewProtoCodec(val.ClientCtx.InterfaceRegistry).GRPCCodec()),
				grpc.MaxCallRecvMsgSize(maxRecvMsgSize),
				grpc.MaxCallSendMsgSize(maxSendMsgSize),
			),
		)
		if err != nil {
			return err
		}
		val.ClientCtx = val.ClientCtx.WithGRPCClient(grpcClient)
	}

	if val.AppConfig.API.Enable && val.APIAddress != "" {
		val.api = api.New(val.ClientCtx, logger.With("module", "api-server"))
		myApp.RegisterAPIRoutes(val.api, val.AppConfig.API)

		go func() {
			_ = val.api.Start(val.AppConfig.Config)
		}()
	}

	if val.AppConfig.GRPC.Enable {
		val.grpc, err = StartGRPCServer(val.ClientCtx, myApp, val.AppConfig.GRPC)
		if err != nil {
			return err
		}

		if val.AppConfig.GRPCWeb.Enable {
			val.grpcWeb, err = servergrpc.StartGRPCWeb(val.grpc, val.AppConfig.Config)
			if err != nil {
				return err
			}
		}
	}

	if val.AppConfig.JSONRPC.Enable && val.AppConfig.JSONRPC.Address != "" {
		if val.Ctx == nil || val.Ctx.Viper == nil {
			return fmt.Errorf("validator %s context is nil", val.Ctx.Config.Moniker)
		}

		val.ClientCtx = val.ClientCtx.WithChainID(fxtypes.ChainIdWithEIP155())
		val.jsonrpc, val.jsonrpcDone, err = StartJSONRPC(val.Ctx, val.ClientCtx, val.AppConfig.ToEthermintConfig(), nil, myApp.(server.AppWithPendingTxStream))
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
	for i := 0; i < cfg.NumValidators; i++ {
		tmCfg := vals[i].Ctx.Config

		gentxsDir := filepath.Join(outputDir, "gentxs")
		initCfg := genutiltypes.NewInitConfig(cfg.ChainID, gentxsDir, vals[i].NodeID, vals[i].PubKey)

		genFile := tmCfg.GenesisFile()
		genDoc, err := types.GenesisDocFromFile(genFile)
		if err != nil {
			return err
		}

		_, err = genutil.GenAppStateFromConfig(cfg.Codec, cfg.TxConfig, tmCfg, initCfg, *genDoc, banktypes.GenesisBalancesIterator{}, genutiltypes.DefaultMessageValidator)
		if err != nil {
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

	genDoc := types.GenesisDoc{
		GenesisTime:     time.Now(),
		ChainID:         cfg.ChainID,
		InitialHeight:   1,
		ConsensusParams: app.CustomGenesisConsensusParams(),
		Validators:      nil,
		AppHash:         nil,
		AppState:        appGenStateJSON,
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
