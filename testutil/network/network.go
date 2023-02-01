package network

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

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
	ethermint "github.com/evmos/ethermint/types"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	db "github.com/tendermint/tm-db"
	"google.golang.org/grpc"

	fxcfg "github.com/functionx/fx-core/v3/server/config"
	fxtypes "github.com/functionx/fx-core/v3/types"
)

// package-wide network lock to only allow one test network at a time
var lock = new(sync.Mutex)

// AppConstructor defines a function which accepts a network configuration and
// creates an ABCI Application to provide to Tendermint.
type AppConstructor = func(val Validator) servertypes.Application

// Config defines the necessary configuration used to bootstrap and start an
// in-process local testing network.
type Config struct {
	Codec             codec.Codec
	LegacyAmino       *codec.LegacyAmino
	InterfaceRegistry codectypes.InterfaceRegistry
	TxConfig          client.TxConfig
	AccountRetriever  client.AccountRetriever
	AppConstructor    AppConstructor             // the ABCI application constructor
	GenesisState      map[string]json.RawMessage // custom gensis state to provide
	TimeoutCommit     time.Duration              // the consensus commitment timeout
	AccountTokens     sdk.Int                    // the amount of unique validator tokens (e.g. 1000node0)
	StakingTokens     sdk.Int                    // the amount of tokens each validator has available to stake
	BondedTokens      sdk.Int                    // the amount of tokens each validator stakes
	NumValidators     int                        // the total number of validators to create and bond
	Mnemonics         []string                   // custom user-provided validator operator mnemonics
	ChainID           string                     // the network chain-id
	BondDenom         string                     // the staking bond denomination
	MinGasPrices      string                     // the minimum gas prices each validator will accept
	PruningStrategy   string                     // the pruning strategy each validator will have
	RPCAddress        string                     // RPC listen address (including port)
	JSONRPCAddress    string                     // JSON-RPC listen address (including port)
	APIAddress        string                     // REST API listen address (including port)
	GRPCAddress       string                     // GRPC server listen address (including port)
	EnableTMLogging   bool                       // enable Tendermint logging to STDOUT
	CleanupDir        bool                       // remove base temporary directory during cleanup
	SigningAlgo       string                     // signing algorithm for keys
	KeyringOptions    []keyring.Option           // keyring configuration options
	PrintMnemonic     bool                       // print the mnemonic of first validator as log output for testing
}

type (
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
	Network struct {
		Logger     Logger
		BaseDir    string
		Validators []*Validator

		Config Config
	}

	// Validator defines an in-process Tendermint validator node. Through this object,
	// a client can make RPC and API calls and interact with any client command
	// or handler.
	Validator struct {
		AppConfig     *fxcfg.Config
		ClientCtx     client.Context
		Ctx           *server.Context
		Dir           string
		NodeID        string
		PubKey        cryptotypes.PubKey
		Moniker       string
		APIAddress    string
		RPCAddress    string
		P2PAddress    string
		Address       sdk.AccAddress
		ValAddress    sdk.ValAddress
		RPCClient     tmclient.Client
		JSONRPCClient *ethclient.Client

		tmNode      *node.Node
		api         *api.Server
		grpc        *grpc.Server
		grpcWeb     *http.Server
		jsonrpc     *http.Server
		jsonrpcDone chan struct{}
	}
)

// Logger is a network logger interface that exposes testnet-level Log() methods for an in-process testing network
// This is not to be confused with logging that may happen at an individual node or validator level
type Logger interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

var (
	_ Logger = (*testing.T)(nil)
	_ Logger = (*CLILogger)(nil)
)

type CLILogger struct {
	cmd *cobra.Command
}

func (s CLILogger) Log(args ...interface{}) {
	s.cmd.Println(args...)
}

func (s CLILogger) Logf(format string, args ...interface{}) {
	s.cmd.Printf(format, args...)
}

func NewCLILogger(cmd *cobra.Command) CLILogger {
	return CLILogger{cmd}
}

// New creates a new Network for integration tests or in-process testnets run via the CLI
//
//gocyclo:ignore
func New(l Logger, baseDir string, cfg Config) (*Network, error) {
	// only one caller/test can create and use a network at a time
	l.Log("acquiring test network lock")
	lock.Lock()

	if !ethermint.IsValidChainID(fxtypes.ChainIdWithEIP155()) {
		return nil, fmt.Errorf("invalid chain-id: %s", cfg.ChainID)
	}
	l.Logf("preparing test network with chain-id \"%s\"\n", cfg.ChainID)

	var (
		network = &Network{
			Logger:     l,
			BaseDir:    baseDir,
			Validators: make([]*Validator, cfg.NumValidators),
			Config:     cfg,
		}

		nodeIDs    = make([]string, network.Config.NumValidators)
		valPubKeys = make([]cryptotypes.PubKey, network.Config.NumValidators)

		genAccounts = make([]authtypes.GenesisAccount, network.Config.NumValidators)
		genBalances = make([]banktypes.Balance, network.Config.NumValidators)
		genFiles    = make([]string, network.Config.NumValidators)
		startTime   = time.Now()
	)
	if network.Config.NumValidators > 1 {
		if network.Config.TimeoutCommit > 0 && network.Config.TimeoutCommit < 500*time.Millisecond {
			return nil, fmt.Errorf("timeout commit is too small")
		}
	}

	// generate private keys, node IDs, and initial transactions
	for i := 0; i < network.Config.NumValidators; i++ {
		ctx := server.NewDefaultContext()
		ctx.Logger = log.NewNopLogger()
		tmCfg := ctx.Config
		tmCfg.DBBackend = string(db.MemDBBackend)
		if network.Config.NumValidators == 1 {
			tmCfg.Consensus = tmcfg.TestConsensusConfig()
		}
		if network.Config.TimeoutCommit > 0 {
			tmCfg.Consensus.SkipTimeoutCommit = false
			tmCfg.Consensus.TimeoutCommit = network.Config.TimeoutCommit
		}
		tmCfg.RPC.ListenAddress = ""
		tmCfg.Instrumentation.Prometheus = false

		appCfg := fxcfg.DefaultConfig()
		appCfg.Pruning = network.Config.PruningStrategy
		appCfg.MinGasPrices = network.Config.MinGasPrices
		appCfg.Telemetry.Enabled = false
		appCfg.Telemetry.GlobalLabels = [][]string{{"chain_id", network.Config.ChainID}}
		appCfg.Rosetta.Enable = false
		appCfg.API.Enable = false
		appCfg.API.Swagger = false
		appCfg.GRPC.Enable = false
		appCfg.GRPCWeb.Enable = false
		appCfg.JSONRPC.Enable = false

		if i == 0 {
			if network.Config.APIAddress != "" {
				appCfg.API.Address = network.Config.APIAddress
			} else {
				apiAddr, _, err := server.FreeTCPAddr()
				if err != nil {
					return nil, err
				}
				appCfg.API.Address = apiAddr
			}
			appCfg.API.Enable = true

			if network.Config.RPCAddress != "" {
				tmCfg.RPC.ListenAddress = network.Config.RPCAddress
			} else {
				rpcAddr, _, err := server.FreeTCPAddr()
				if err != nil {
					return nil, err
				}
				tmCfg.RPC.ListenAddress = rpcAddr
			}

			if network.Config.GRPCAddress != "" {
				appCfg.GRPC.Address = network.Config.GRPCAddress
			} else {
				_, grpcPort, err := server.FreeTCPAddr()
				if err != nil {
					return nil, err
				}
				appCfg.GRPC.Address = fmt.Sprintf("0.0.0.0:%s", grpcPort)
			}
			appCfg.GRPC.Enable = true

			// _, grpcWebPort, err := server.FreeTCPAddr()
			// if err != nil {
			//	return nil, err
			// }
			// appCfg.GRPCWeb.Address = fmt.Sprintf("0.0.0.0:%s", grpcWebPort)
			// appCfg.GRPCWeb.Enable = true

			if network.Config.JSONRPCAddress != "" {
				appCfg.JSONRPC.Address = network.Config.JSONRPCAddress
			} else {
				_, jsonRPCPort, err := server.FreeTCPAddr()
				if err != nil {
					return nil, err
				}
				appCfg.JSONRPC.Address = fmt.Sprintf("0.0.0.0:%s", jsonRPCPort)
			}
			appCfg.JSONRPC.Enable = true
			appCfg.JSONRPC.API = config.GetAPINamespaces()

			if network.Config.EnableTMLogging {
				ctx.Logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
				var err error
				ctx.Logger, err = tmflags.ParseLogLevel("info", ctx.Logger, tmcfg.DefaultLogLevel)
				if err != nil {
					return nil, err
				}
			}
		}

		nodeDirName := fmt.Sprintf("node%d", i)
		nodeDir := filepath.Join(network.BaseDir, nodeDirName, strings.ToLower(network.Config.ChainID))
		clientDir := filepath.Join(network.BaseDir, nodeDirName, strings.ToLower(network.Config.ChainID))

		if err := os.MkdirAll(filepath.Join(nodeDir, "config"), 0o750); err != nil {
			return nil, err
		}

		if err := os.MkdirAll(clientDir, 0o750); err != nil {
			return nil, err
		}

		tmCfg.SetRoot(nodeDir)
		tmCfg.Moniker = nodeDirName

		// proxyAddr, _, err := server.FreeTCPAddr()
		// if err != nil {
		//	return nil, err
		// }
		// tmCfg.ProxyApp = proxyAddr

		_, p2pPort, err := server.FreeTCPAddr()
		if err != nil {
			return nil, err
		}
		tmCfg.P2P.ListenAddress = fmt.Sprintf("127.0.0.1:%s", p2pPort)
		tmCfg.P2P.AddrBookStrict = false
		tmCfg.P2P.AllowDuplicateIP = true

		nodeID, pubKey, err := genutil.InitializeNodeValidatorFiles(tmCfg)
		if err != nil {
			return nil, err
		}
		nodeIDs[i] = nodeID
		valPubKeys[i] = pubKey

		kb, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendMemory, clientDir, nil, network.Config.KeyringOptions...)
		if err != nil {
			return nil, err
		}

		keyringAlgos, _ := kb.SupportedAlgorithms()
		algo, err := keyring.NewSigningAlgoFromString(network.Config.SigningAlgo, keyringAlgos)
		if err != nil {
			return nil, err
		}

		var mnemonic string
		if i < len(network.Config.Mnemonics) {
			mnemonic = network.Config.Mnemonics[i]
		}

		valAddr, secret, err := testutil.GenerateSaveCoinKey(kb, nodeDirName, mnemonic, true, algo)
		if err != nil {
			return nil, err
		}
		if i >= len(network.Config.Mnemonics) {
			network.Config.Mnemonics = append(network.Config.Mnemonics, secret)
		}

		// if PrintMnemonic is set to true, we print the first validator node's secret to the network's logger
		// for debugging and manual testing
		if network.Config.PrintMnemonic && i == 0 {
			printMnemonic(l, secret)
		}

		/*info := map[string]string{"secret": secret}
		infoBz, err := json.Marshal(info)
		if err != nil {
			return nil, err
		}

		// save private key seed words
		err = writeFile(fmt.Sprintf("%v.json", "key_seed"), clientDir, infoBz)
		if err != nil {
			return nil, err
		}*/

		balances := sdk.NewCoins(sdk.NewCoin(network.Config.BondDenom, network.Config.StakingTokens))

		genFiles[i] = tmCfg.GenesisFile()
		genBalances[i] = banktypes.Balance{Address: valAddr.String(), Coins: balances.Sort()}
		genAccounts[i] = authtypes.NewBaseAccount(valAddr, nil, 0, 0)

		createValMsg, err := stakingtypes.NewMsgCreateValidator(
			sdk.ValAddress(valAddr),
			valPubKeys[i],
			sdk.NewCoin(network.Config.BondDenom, network.Config.BondedTokens),
			stakingtypes.NewDescription(nodeDirName, "", "", "", ""),
			stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.OneDec(), sdk.OneDec()), // 5%
			sdk.OneInt(),
		)
		if err != nil {
			return nil, err
		}

		memo := fmt.Sprintf("%s@%s", nodeIDs[i], tmCfg.P2P.ListenAddress)
		fee := sdk.NewCoins(sdk.NewCoin(network.Config.BondDenom, sdk.NewInt(0)))
		txBuilder := network.Config.TxConfig.NewTxBuilder()
		err = txBuilder.SetMsgs(createValMsg)
		if err != nil {
			return nil, err
		}
		txBuilder.SetFeeAmount(fee)    // Arbitrary fee
		txBuilder.SetGasLimit(1000000) // Need at least 100386
		txBuilder.SetMemo(memo)

		txFactory := tx.Factory{}
		txFactory = txFactory.
			WithChainID(network.Config.ChainID).
			WithMemo(memo).
			WithKeybase(kb).
			WithTxConfig(network.Config.TxConfig)

		if err = tx.Sign(txFactory, nodeDirName, txBuilder, true); err != nil {
			return nil, err
		}

		txBz, err := network.Config.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
		if err != nil {
			return nil, err
		}

		gentxsDir := filepath.Join(network.BaseDir, "gentxs")
		if err = writeFile(fmt.Sprintf("%v.json", nodeDirName), gentxsDir, txBz); err != nil {
			return nil, err
		}

		srvconfig.SetConfigTemplate(fxcfg.DefaultConfigTemplate())
		srvconfig.WriteConfigFile(filepath.Join(nodeDir, "config/app.toml"), appCfg)

		ctx.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
		ctx.Viper.SetConfigFile(filepath.Join(nodeDir, "config/app.toml"))
		if err = ctx.Viper.ReadInConfig(); err != nil {
			return nil, err
		}

		clientCtx := client.Context{}.
			WithKeyringDir(clientDir).
			WithKeyring(kb).
			WithHomeDir(tmCfg.RootDir).
			WithChainID(network.Config.ChainID).
			WithInterfaceRegistry(network.Config.InterfaceRegistry).
			WithCodec(network.Config.Codec).
			WithLegacyAmino(network.Config.LegacyAmino).
			WithTxConfig(network.Config.TxConfig).
			WithAccountRetriever(network.Config.AccountRetriever).
			WithKeyringOptions(network.Config.KeyringOptions...)

		network.Validators[i] = &Validator{
			AppConfig:  appCfg,
			ClientCtx:  clientCtx,
			Ctx:        ctx,
			Dir:        filepath.Join(network.BaseDir, nodeDirName),
			NodeID:     nodeID,
			PubKey:     pubKey,
			Moniker:    nodeDirName,
			RPCAddress: tmCfg.RPC.ListenAddress,
			P2PAddress: tmCfg.P2P.ListenAddress,
			APIAddress: appCfg.API.Address,
			Address:    valAddr,
			ValAddress: sdk.ValAddress(valAddr),
		}
	}

	if err := initGenFiles(network.Config, genAccounts, genBalances, genFiles); err != nil {
		return nil, err
	}

	if err := collectGenFiles(network.Config, network.Validators, network.BaseDir); err != nil {
		return nil, err
	}

	l.Log("starting test network...")
	for _, v := range network.Validators {
		if err := startInProcess(network.Config, v); err != nil {
			return nil, err
		}
	}

	// Ensure we cleanup incase any test was abruptly halted (e.g. SIGINT) as any
	// defer in a test would not be called.
	server.TrapSignal(network.Cleanup)

	l.Logf("started test network %fs", time.Since(startTime).Seconds())
	return network, nil
}

// LatestHeight returns the latest height of the network or an error if the
// query fails or no validators exist.
func (n *Network) LatestHeight() (int64, error) {
	if len(n.Validators) == 0 {
		return 0, errors.New("no validators available")
	}

	status, err := n.Validators[0].RPCClient.Status(context.Background())
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
			status, err := val.RPCClient.Status(context.Background())
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
		n.Logger.Log("released test network lock")
	}()

	n.Logger.Log("cleaning up test network...")
	startTime := time.Now()
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
		if v.grpcWeb != nil {
			_ = v.grpcWeb.Close()
		}

		if v.jsonrpc != nil {
			shutdownCtx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
			if err := v.jsonrpc.Shutdown(shutdownCtx); err != nil {
				n.Logger.Log("HTTP server shutdown produced a warning", "error", err.Error())
			} else {
				n.Logger.Log("HTTP server shut down, waiting...")
				select {
				case <-time.Tick(1 * time.Second):
				case <-v.jsonrpcDone:
				}
			}
			cancelFn()
		}
	}

	if n.Config.CleanupDir {
		if err := os.RemoveAll(n.BaseDir); err != nil {
			n.Logger.Log("remove base dir", "error", err.Error())
		}
	}

	n.Logger.Logf("finished cleaning up test network %fs", time.Since(startTime).Seconds())
}
