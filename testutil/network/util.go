package network

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/server/api"
	servergrpc "github.com/cosmos/cosmos-sdk/server/grpc"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	mintypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/evmos/ethermint/server"
	"github.com/evmos/ethermint/server/config"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	pvm "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/rpc/client/local"
	"github.com/tendermint/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
	fxtypes "github.com/functionx/fx-core/v3/types"
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
		logger.With("moniker", val.Ctx.Config.Moniker),
	)
	if err != nil {
		return err
	}

	if err = tmNode.Start(); err != nil {
		return err
	}

	val.tmNode = tmNode

	if val.RPCAddress != "" {
		val.RPCClient = local.New(tmNode)
	}

	// We'll need a RPC client if the validator exposes a gRPC or REST endpoint.
	if val.APIAddress != "" || val.AppConfig.GRPC.Enable {
		val.ClientCtx = val.ClientCtx.
			WithClient(val.RPCClient)

		// Add the tx service in the gRPC router.
		myApp.RegisterTxService(val.ClientCtx)

		// Add the tendermint queries service in the gRPC router.
		myApp.RegisterTendermintService(val.ClientCtx)
	}

	errCh := make(chan error, 8)
	wg := sync.WaitGroup{}

	if val.AppConfig.API.Enable && val.APIAddress != "" {
		val.api = api.New(val.ClientCtx, logger.With("module", "api-server"))
		myApp.RegisterAPIRoutes(val.api, val.AppConfig.API)

		go func() {
			if err = val.api.Start(val.AppConfig.Config); err != nil {
				errCh <- err
			}
		}()
	}

	if val.AppConfig.GRPC.Enable {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val.grpc, err = servergrpc.StartGRPCServer(val.ClientCtx, myApp, val.AppConfig.GRPC.Address)
			if err != nil {
				errCh <- err
			}
		}()

		if val.AppConfig.GRPCWeb.Enable {
			wg.Add(1)
			go func() {
				defer wg.Done()
				val.grpcWeb, err = servergrpc.StartGRPCWeb(val.grpc, val.AppConfig.Config)
				if err != nil {
					errCh <- err
				}
			}()
		}
	}

	if val.AppConfig.JSONRPC.Enable && val.AppConfig.JSONRPC.Address != "" {
		if val.Ctx == nil || val.Ctx.Viper == nil {
			return fmt.Errorf("validator %s context is nil", val.Ctx.Config.Moniker)
		}

		tmEndpoint := "/websocket"
		tmRPCAddr := val.RPCAddress

		wg.Add(1)
		go func() {
			defer wg.Done()
			clientCtx := val.ClientCtx.WithChainID(fxtypes.ChainIdWithEIP155())
			val.jsonrpc, val.jsonrpcDone, err = server.StartJSONRPC(val.Ctx, clientCtx, tmRPCAddr, tmEndpoint, &config.Config{
				Config:  val.AppConfig.Config,
				EVM:     val.AppConfig.EVM,
				JSONRPC: val.AppConfig.JSONRPC,
				TLS:     val.AppConfig.TLS,
			}, nil)
			if err != nil {
				errCh <- err
				return
			}
			val.JSONRPCClient, err = ethclient.Dial(fmt.Sprintf("http://%s", val.AppConfig.JSONRPC.Address))
			if err != nil {
				errCh <- fmt.Errorf("failed to dial JSON-RPC at %s: %w", val.AppConfig.JSONRPC.Address, err)
			}
		}()
	}

	wg.Wait()
	select {
	case err = <-errCh:
		return err
	default:
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

	var stakingGenState stakingtypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[stakingtypes.ModuleName], &stakingGenState)

	stakingGenState.Params.BondDenom = cfg.BondDenom
	cfg.GenesisState[stakingtypes.ModuleName] = cfg.Codec.MustMarshalJSON(&stakingGenState)

	var govGenState govtypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[govtypes.ModuleName], &govGenState)

	govGenState.DepositParams.MinDeposit[0].Denom = cfg.BondDenom
	cfg.GenesisState[govtypes.ModuleName] = cfg.Codec.MustMarshalJSON(&govGenState)

	var mintGenState mintypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[mintypes.ModuleName], &mintGenState)

	mintGenState.Params.MintDenom = cfg.BondDenom
	cfg.GenesisState[mintypes.ModuleName] = cfg.Codec.MustMarshalJSON(&mintGenState)

	var crisisGenState crisistypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[crisistypes.ModuleName], &crisisGenState)

	crisisGenState.ConstantFee.Denom = cfg.BondDenom
	cfg.GenesisState[crisistypes.ModuleName] = cfg.Codec.MustMarshalJSON(&crisisGenState)

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
	file := filepath.Join(dir, name)

	err := tmos.EnsureDir(dir, 0o755)
	if err != nil {
		return err
	}

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
