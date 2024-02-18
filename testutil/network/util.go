package network

import (
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server/api"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servergrpc "github.com/cosmos/cosmos-sdk/server/grpc"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/ethereum/go-ethereum/ethclient"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	pvm "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/rpc/client/local"
	"github.com/tendermint/tendermint/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/server"
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

		if a, ok := myApp.(servertypes.ApplicationQueryService); ok {
			a.RegisterNodeService(val.ClientCtx)
		}
	}

	if val.AppConfig.GRPC.Enable {
		_, port, err := net.SplitHostPort(val.AppConfig.GRPC.Address)
		if err != nil {
			return err
		}

		maxSendMsgSize := val.AppConfig.GRPC.MaxSendMsgSize
		if maxSendMsgSize == 0 {
			maxSendMsgSize = serverconfig.DefaultGRPCMaxSendMsgSize
		}

		maxRecvMsgSize := val.AppConfig.GRPC.MaxRecvMsgSize
		if maxRecvMsgSize == 0 {
			maxRecvMsgSize = serverconfig.DefaultGRPCMaxRecvMsgSize
		}

		grpcAddress := fmt.Sprintf("127.0.0.1:%s", port)

		// If grpc is enabled, configure grpc client for grpc gateway.
		grpcClient, err := grpc.Dial(
			grpcAddress,
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
		val.grpc, err = servergrpc.StartGRPCServer(val.ClientCtx, myApp, val.AppConfig.GRPC)
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

		tmEndpoint := "/websocket"
		tmRPCAddr := val.RPCAddress

		clientCtx := val.ClientCtx.WithChainID(fxtypes.ChainIdWithEIP155())
		val.jsonrpc, err = server.StartJSONRPC(val.Ctx, clientCtx, tmRPCAddr, tmEndpoint, val.AppConfig.ToEthermintConfig(), nil)
		if err != nil {
			return err
		}
		ln, err := server.Listen(val.jsonrpc.Addr, val.AppConfig.JSONRPC.MaxOpenConnections)
		if err != nil {
			return err
		}
		go func() {
			_ = val.jsonrpc.Serve(ln)
		}()
		val.JSONRPCClient, err = ethclient.Dial(fmt.Sprintf("http://%s", val.AppConfig.JSONRPC.Address))
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

		_, err = genutil.GenAppStateFromConfig(cfg.Codec, cfg.TxConfig, tmCfg, initCfg, *genDoc, banktypes.GenesisBalancesIterator{})
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

	customConsensusParams := app.CustomGenesisConsensusParams()
	customConsensusParams.Block.TimeIotaMs = cfg.TimeoutCommit.Milliseconds()
	genDoc := types.GenesisDoc{
		GenesisTime:     time.Now(),
		ChainID:         cfg.ChainID,
		InitialHeight:   1,
		ConsensusParams: customConsensusParams,
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

// printMnemonic prints a provided mnemonic seed phrase for debugging and manual testing
func printMnemonic(secret string) {
	lines := []string{
		"THIS MNEMONIC IS FOR TESTING PURPOSES ONLY",
		"DO NOT USE IN PRODUCTION",
		"",
		strings.Join(strings.Fields(secret)[0:8], " "),
		strings.Join(strings.Fields(secret)[8:16], " "),
		strings.Join(strings.Fields(secret)[16:24], " "),
	}

	lineLengths := make([]int, len(lines))
	for i, line := range lines {
		lineLengths[i] = len(line)
	}

	maxLineLength := 0
	for _, lineLen := range lineLengths {
		if lineLen > maxLineLength {
			maxLineLength = lineLen
		}
	}

	fmt.Print("\n")
	fmt.Print(strings.Repeat("+", maxLineLength+8))
	for _, line := range lines {
		fmt.Printf("++  %s  ++\n", centerText(line, maxLineLength))
	}
	fmt.Print(strings.Repeat("+", maxLineLength+8))
	fmt.Print("\n")
}

// centerText centers text across a fixed width, filling either side with whitespace buffers
func centerText(text string, width int) string {
	textLen := len(text)
	leftBuffer := strings.Repeat(" ", (width-textLen)/2)
	rightBuffer := strings.Repeat(" ", (width-textLen)/2+(width-textLen)%2)

	return fmt.Sprintf("%s%s%s", leftBuffer, text, rightBuffer)
}
