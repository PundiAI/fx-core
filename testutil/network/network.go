package network

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	tmcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/node"
	tmclient "github.com/cometbft/cometbft/rpc/client"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/evmos/ethermint/server/config"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	fxcfg "github.com/pundiai/fx-core/v8/server/config"
)

// package-wide network lock to only allow one test network at a time
var lock = new(sync.Mutex)

// AppConstructor defines a function which accepts a network configuration and
// creates an ABCI Application to provide to Tendermint.
type AppConstructor = func(appConfig *fxcfg.Config, ctx *server.Context) servertypes.Application

// Config defines the necessary configuration used to bootstrap and start an
// in-process local testing network.
type Config struct {
	Codec                codec.Codec
	InterfaceRegistry    codectypes.InterfaceRegistry
	TxConfig             client.TxConfig
	AccountRetriever     client.AccountRetriever
	AppConstructor       AppConstructor             // the ABCI application constructor
	GenesisState         map[string]json.RawMessage // custom gensis state to provide
	TimeoutCommit        time.Duration              // the consensus commitment timeout
	StakingTokens        sdkmath.Int                // the amount of tokens each validator has available to stake
	BondedTokens         sdkmath.Int                // the amount of tokens each validator stakes
	NumValidators        int                        // the total number of validators to create and bond
	Mnemonics            []string                   // custom user-provided validator operator mnemonics
	ChainID              string                     // the network chain-id
	BondDenom            string                     // the staking bond denomination
	MinGasPrices         string                     // the minimum gas prices each validator will accept
	PruningStrategy      string                     // the pruning strategy each validator will have
	RPCAddress           string                     // RPC listen address (including port)
	JSONRPCAddress       string                     // JSON-RPC listen address (including port)
	APIAddress           string                     // REST API listen address (including port)
	GRPCAddress          string                     // GRPC server listen address (including port)
	EnableJSONRPC        bool                       // enable JSON-RPC service
	EnableAPI            bool                       // enable REST API service
	EnableTMLogging      bool                       // enable Tendermint logging to STDOUT
	CleanupDir           bool                       // remove base temporary directory during cleanup
	SigningAlgo          string                     // signing algorithm for keys
	KeyringOptions       []keyring.Option           // keyring configuration options
	BypassMinFeeMsgTypes []string                   // bypass minimum fee check for the given message types
}

// Network defines a local in-process testing network using SimApp. It can be
// configured to start any number of validators, each with its own RPC and API
// clients. Typically, this test network would be used in client and integration
// testing where user input is expected.
//
// Note, due to Tendermint constraints in regards to RPC functionality, there
// may only be one test network running at a time. Thus, any caller must be
// sure to Cleanup after testing is finished in order to allow other tests
// to create networks. In addition, only the first validator will have a valid
// RPC and API server/client.
type Network struct {
	logger     log.Logger
	BaseDir    string
	Config     Config
	Validators []*Validator
	ctx        context.Context
	cancel     context.CancelFunc
}

// Validator defines an in-process Tendermint validator node. Through this object,
// a client can make RPC and API calls and interact with any client command
// or handler.
type Validator struct {
	AppConfig     *fxcfg.Config
	ClientCtx     client.Context
	Ctx           *server.Context
	NodeID        string
	PubKey        cryptotypes.PubKey
	Address       sdk.AccAddress
	ValAddress    sdk.ValAddress
	RPCClient     tmclient.Client
	JSONRPCClient *ethclient.Client

	tmNode  *node.Node
	api     *api.Server
	grpc    *grpc.Server
	jsonrpc *http.Server
}

// New creates a new Network for integration tests or in-process testnets run via the CLI
func New(t *testing.T, cfg Config) *Network {
	t.Helper()
	// only one caller/test can create and use a network at a time
	t.Log("acquiring test network lock")
	lock.Lock()

	baseDir := t.TempDir()

	t.Logf("preparing test network with chain-id \"%s\"\n", cfg.ChainID)
	startTime := time.Now()

	logger := log.NewNopLogger()
	if cfg.EnableTMLogging {
		filterFunc, _ := log.ParseLogLevel("info")
		logger = log.NewLogger(os.Stdout, log.FilterOption(filterFunc))
	}
	validators, err := GenerateGenesisAndValidators(baseDir, &cfg, logger)
	require.NoError(t, err)

	network := &Network{
		ctx:        context.Background(),
		logger:     logger,
		BaseDir:    baseDir,
		Config:     cfg,
		Validators: validators,
	}

	network.ctx, network.cancel = context.WithCancel(network.ctx)
	errGroup, ctx := errgroup.WithContext(network.ctx)

	t.Log("starting test network...")
	for _, val := range network.Validators {
		if err = StartInProcess(ctx, errGroup, network.Config.AppConstructor, val); err != nil {
			require.NoError(t, err)
		}
	}

	// Ensure we cleanup incase any test was abruptly halted (e.g. SIGINT) as any
	// defer in a test would not be called.
	network.TrapSignal(errGroup)

	t.Logf("started test network %fs", time.Since(startTime).Seconds())
	return network
}

// GenerateGenesisAndValidators
//
//nolint:gocyclo // for testing
func GenerateGenesisAndValidators(baseDir string, cfg *Config, logger log.Logger) ([]*Validator, error) {
	srvconfig.SetConfigTemplate(fxcfg.DefaultConfigTemplate())

	if cfg.NumValidators < 1 || cfg.NumValidators > 100 {
		return nil, fmt.Errorf("the number of validators must be between [1,100]")
	}

	if cfg.NumValidators > 1 {
		if cfg.TimeoutCommit > 0 && cfg.TimeoutCommit < 500*time.Millisecond {
			return nil, fmt.Errorf("timeout commit is too small")
		}
	}
	var (
		nodeIDs    = make([]string, cfg.NumValidators)
		valPubKeys = make([]cryptotypes.PubKey, cfg.NumValidators)
		validators = make([]*Validator, cfg.NumValidators)

		genAccounts = make([]authtypes.GenesisAccount, cfg.NumValidators)
		genBalances = make([]banktypes.Balance, cfg.NumValidators)
		genFiles    = make([]string, cfg.NumValidators)
	)

	// generate private keys, node IDs, and initial transactions
	for i := 0; i < cfg.NumValidators; i++ {
		srvCtx := server.NewDefaultContext()
		srvCtx.Logger = logger
		srvCtx.Config.DBBackend = string(dbm.MemDBBackend)
		if cfg.NumValidators == 1 {
			srvCtx.Config.Consensus = tmcfg.TestConsensusConfig()
		}
		if cfg.TimeoutCommit > 0 {
			srvCtx.Config.Consensus.SkipTimeoutCommit = false
			srvCtx.Config.Consensus.TimeoutCommit = cfg.TimeoutCommit
		}
		srvCtx.Config.RPC.PprofListenAddress = ""
		srvCtx.Config.RPC.ListenAddress = ""
		srvCtx.Config.Instrumentation.Prometheus = false

		appCfg := fxcfg.DefaultConfig()
		appCfg.Pruning = cfg.PruningStrategy
		appCfg.MinGasPrices = cfg.MinGasPrices
		appCfg.Telemetry.Enabled = false
		appCfg.Telemetry.GlobalLabels = [][]string{{"chain_id", cfg.ChainID}}
		appCfg.API.Enable = false
		appCfg.API.Swagger = false
		appCfg.GRPC.Enable = false
		appCfg.GRPCWeb.Enable = false
		appCfg.JSONRPC.Enable = false
		appCfg.BypassMinFee.MsgMaxGasUsage = 500_000
		appCfg.BypassMinFee.MsgTypes = cfg.BypassMinFeeMsgTypes

		if i == 0 {
			if cfg.EnableAPI {
				if cfg.APIAddress != "" {
					appCfg.API.Address = cfg.APIAddress
				} else {
					apiAddr, _, err := FreeTCPAddr()
					if err != nil {
						return nil, err
					}
					appCfg.API.Address = apiAddr
				}
				appCfg.API.Enable = true
			}

			if cfg.RPCAddress != "" {
				srvCtx.Config.RPC.ListenAddress = cfg.RPCAddress
			} else {
				rpcAddr, _, err := FreeTCPAddr()
				if err != nil {
					return nil, err
				}
				srvCtx.Config.RPC.ListenAddress = rpcAddr
			}

			if cfg.GRPCAddress != "" {
				appCfg.GRPC.Address = cfg.GRPCAddress
			} else {
				_, grpcPort, err := FreeTCPAddr()
				if err != nil {
					return nil, err
				}
				appCfg.GRPC.Address = fmt.Sprintf("0.0.0.0:%s", grpcPort)
			}
			appCfg.GRPC.Enable = true

			if cfg.EnableJSONRPC {
				if cfg.JSONRPCAddress != "" {
					appCfg.JSONRPC.Address = cfg.JSONRPCAddress
				} else {
					_, jsonRPCPort, err := FreeTCPAddr()
					if err != nil {
						return nil, err
					}
					appCfg.JSONRPC.Address = fmt.Sprintf("0.0.0.0:%s", jsonRPCPort)
				}
				appCfg.JSONRPC.Enable = true
			}

			appCfg.JSONRPC.API = config.GetAPINamespaces()
			appCfg.JSONRPC.WsAddress = ""
		}

		nodeName := fmt.Sprintf("node%d", i)
		nodeDir := filepath.Join(baseDir, nodeName, strings.ToLower(cfg.ChainID))
		clientDir := filepath.Join(baseDir, nodeName, strings.ToLower(cfg.ChainID))

		if err := os.MkdirAll(filepath.Join(nodeDir, "config"), 0o750); err != nil {
			return nil, err
		}

		if err := os.MkdirAll(clientDir, 0o750); err != nil {
			return nil, err
		}

		srvCtx.Config.SetRoot(nodeDir)
		srvCtx.Config.Moniker = nodeName

		_, p2pPort, err := FreeTCPAddr()
		if err != nil {
			return nil, err
		}
		srvCtx.Config.P2P.ListenAddress = fmt.Sprintf("127.0.0.1:%s", p2pPort)
		srvCtx.Config.P2P.AddrBookStrict = false
		srvCtx.Config.P2P.AllowDuplicateIP = true

		nodeID, pubKey, err := genutil.InitializeNodeValidatorFiles(srvCtx.Config)
		if err != nil {
			return nil, err
		}
		nodeIDs[i] = nodeID
		valPubKeys[i] = pubKey

		kb, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendMemory, clientDir, nil, cfg.Codec, cfg.KeyringOptions...)
		if err != nil {
			return nil, err
		}

		keyringAlgos, _ := kb.SupportedAlgorithms()
		algo, err := keyring.NewSigningAlgoFromString(cfg.SigningAlgo, keyringAlgos)
		if err != nil {
			return nil, err
		}

		var mnemonic string
		if i < len(cfg.Mnemonics) {
			mnemonic = cfg.Mnemonics[i]
		}

		valAddr, secret, err := testutil.GenerateSaveCoinKey(kb, nodeName, mnemonic, true, algo)
		if err != nil {
			return nil, err
		}
		if i >= len(cfg.Mnemonics) {
			cfg.Mnemonics = append(cfg.Mnemonics, secret)
		}

		balances := sdk.NewCoins(sdk.NewCoin(cfg.BondDenom, cfg.StakingTokens))

		genFiles[i] = srvCtx.Config.GenesisFile()
		genBalances[i] = banktypes.Balance{Address: valAddr.String(), Coins: balances.Sort()}
		genAccounts[i] = authtypes.NewBaseAccount(valAddr, nil, 0, 0)

		createValMsg, err := stakingtypes.NewMsgCreateValidator(
			sdk.ValAddress(valAddr).String(),
			valPubKeys[i],
			sdk.NewCoin(cfg.BondDenom, cfg.BondedTokens),
			stakingtypes.NewDescription(nodeName, "", "", "", ""),
			stakingtypes.NewCommissionRates(sdkmath.LegacyNewDecWithPrec(5, 1), sdkmath.LegacyOneDec(), sdkmath.LegacyOneDec()), // 5%
			sdkmath.OneInt(),
		)
		if err != nil {
			return nil, err
		}

		memo := fmt.Sprintf("%s@%s", nodeIDs[i], srvCtx.Config.P2P.ListenAddress)
		fee := sdk.NewCoins(sdk.NewCoin(cfg.BondDenom, sdkmath.NewInt(0)))
		txBuilder := cfg.TxConfig.NewTxBuilder()
		if err = txBuilder.SetMsgs(createValMsg); err != nil {
			return nil, err
		}
		txBuilder.SetFeeAmount(fee)    // Arbitrary fee
		txBuilder.SetGasLimit(1000000) // Need at least 100386
		txBuilder.SetMemo(memo)

		txFactory := tx.Factory{}
		txFactory = txFactory.
			WithChainID(cfg.ChainID).
			WithMemo(memo).
			WithKeybase(kb).
			WithTxConfig(cfg.TxConfig)

		if err = tx.Sign(context.Background(), txFactory, nodeName, txBuilder, true); err != nil {
			return nil, err
		}

		txBz, err := cfg.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
		if err != nil {
			return nil, err
		}

		gentxsDir := filepath.Join(baseDir, "gentxs")
		if err = writeFile(fmt.Sprintf("%v.json", nodeName), gentxsDir, txBz); err != nil {
			return nil, err
		}

		srvconfig.WriteConfigFile(filepath.Join(nodeDir, "config/app.toml"), appCfg)
		srvCtx.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
		srvCtx.Viper.SetConfigFile(filepath.Join(nodeDir, "config/app.toml"))
		if err = srvCtx.Viper.ReadInConfig(); err != nil {
			return nil, err
		}

		clientCtx := client.Context{}.
			WithKeyringDir(clientDir).
			WithKeyring(kb).
			WithHomeDir(srvCtx.Config.RootDir).
			WithChainID(cfg.ChainID).
			WithInterfaceRegistry(cfg.InterfaceRegistry).
			WithCodec(cfg.Codec).
			WithTxConfig(cfg.TxConfig).
			WithAccountRetriever(cfg.AccountRetriever).
			WithKeyringOptions(cfg.KeyringOptions...)

		validators[i] = &Validator{
			AppConfig:  appCfg,
			ClientCtx:  clientCtx,
			Ctx:        srvCtx,
			NodeID:     nodeID,
			PubKey:     pubKey,
			Address:    valAddr,
			ValAddress: sdk.ValAddress(valAddr),
		}
	}

	if err := initGenFiles(*cfg, genAccounts, genBalances, genFiles); err != nil {
		return nil, err
	}

	if err := collectGenFiles(*cfg, validators, baseDir); err != nil {
		return nil, err
	}
	return validators, nil
}

func (n *Network) GetContext() context.Context {
	return n.ctx
}

// LatestHeight returns the latest height of the network or an error if the
// query fails or no validators exist.
func (n *Network) LatestHeight() (int64, error) {
	if len(n.Validators) == 0 {
		return 0, errors.New("no validators available")
	}

	status, err := n.Validators[0].RPCClient.Status(n.ctx)
	if err != nil {
		return 0, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

// WaitForHeight performs a blocking check where it waits for a block to be
// committed after a given block. If that height is not reached within a timeout,
// an error is returned. Regardless, the latest height queried is returned.
func (n *Network) WaitForHeight(h int64) (int64, error) {
	return n.WaitForHeightWithTimeout(h, 10*time.Second)
}

func (n *Network) WaitNumberBlock(number int64) (int64, error) {
	lastBlock, err := n.LatestHeight()
	if err != nil {
		return 0, err
	}
	return n.WaitForHeightWithTimeout(lastBlock+number, time.Duration(3*number)*n.Config.TimeoutCommit)
}

// WaitForHeightWithTimeout is the same as WaitForHeight except the caller can
// provide a custom timeout.
func (n *Network) WaitForHeightWithTimeout(h int64, t time.Duration) (int64, error) {
	ticker := time.NewTicker(n.Config.TimeoutCommit)
	timeout := time.After(t)

	if len(n.Validators) == 0 {
		return 0, errors.New("no validators available")
	}

	var latestHeight int64
	val := n.Validators[0]

	for {
		select {
		case <-timeout:
			ticker.Stop()
			return latestHeight, errors.New("timeout exceeded waiting for block")
		case <-ticker.C:
			status, err := val.RPCClient.Status(n.ctx)
			if err == nil && status != nil {
				latestHeight = status.SyncInfo.LatestBlockHeight
				if latestHeight >= h {
					return latestHeight, nil
				}
			}
		}
	}
}

// WaitForNextBlock waits for the next block to be committed, returning an error
// upon failure.
func (n *Network) WaitForNextBlock() error {
	lastBlock, err := n.LatestHeight()
	if err != nil {
		return err
	}

	_, err = n.WaitForHeight(lastBlock + 1)
	if err != nil {
		return err
	}

	return err
}

// Cleanup removes the root testing (temporary) directory and stops both the
// Tendermint and API services. It allows other callers to create and start
// test networks. This method must be called when a test is finished, typically
// in a defer.
func (n *Network) Cleanup() {
	defer func() {
		lock.Unlock()
		n.logger.Info("released test network lock")
	}()
	n.logger.Info("cleaning up test network...")
	startTime := time.Now()
	shutdownCtx, cancelFn := context.WithTimeout(n.ctx, 5*time.Second)
	defer cancelFn()
	n.cancel()
	for _, v := range n.Validators {
		if v.tmNode != nil && v.tmNode.IsRunning() {
			_ = v.tmNode.Stop()
		}

		if v.api != nil {
			_ = v.api.Close()
		}

		if v.grpc != nil {
			v.grpc.Stop()
		}

		if v.jsonrpc != nil {
			if err := v.jsonrpc.Shutdown(shutdownCtx); err != nil {
				n.logger.Error("HTTP server shutdown produced a warning", "error", err.Error())
			}
		}
	}

	if n.Config.CleanupDir {
		if err := os.RemoveAll(n.BaseDir); err != nil {
			n.logger.Error("remove base dir", "error", err.Error())
		}
	}

	n.logger.Info("finished cleaning up test network", "time", time.Since(startTime).Seconds())
}

func (n *Network) TrapSignal(group *errgroup.Group) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	group.Go(func() error {
		<-sigs
		n.Cleanup()
		return nil
	})
}
