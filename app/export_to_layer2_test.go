package app_test

import (
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/auth/types"
	ethermint "github.com/evmos/ethermint/types"
	"github.com/functionx/fx-core/v5/app"
	fxtypes "github.com/functionx/fx-core/v5/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

// GenesisFromJSON is the config file for network_custom
type genesisFromJSON struct {
	// L1: root hash of the genesis block
	Root string `json:"root"`
	// L1: block number of the genesis block
	GenesisBlockNum uint64 `json:"genesisBlockNumber"`
	// L2:  List of states contracts used to populate merkle tree at initial state
	Genesis []genesisAccountFromJSON `json:"genesis"`
	// L1: configuration of the network
	L1Config l1Config
}

// L1Config represents the configuration of the network used in L1
type l1Config struct {
	// Chain ID of the L1 network
	L1ChainID uint64 `json:"chainId"`
	// Address of the L1 contract
	CDKValidiumAddr common.Address `json:"cdkValidiumAddress"`
	// Address of the L1 Matic token Contract
	MaticAddr common.Address `json:"maticTokenAddress"`
	// Address of the L1 GlobalExitRootManager contract
	GlobalExitRootManagerAddr common.Address `json:"polygonZkEVMGlobalExitRootAddress"`
	// Address of the data availability committee contract
	DataCommitteeAddr common.Address `json:"cdkDataCommitteeContract"`
}

type genesisAccountFromJSON struct {
	// Address of the account
	Balance string `json:"balance,omitempty"`
	// Nonce of the account
	Nonce string `json:"nonce"`
	// Address of the contract
	Address string `json:"address"`
	// Byte code of the contract
	Bytecode string `json:"bytecode,omitempty"`
	// Initial storage of the contract
	Storage map[string]string `json:"storage,omitempty"`
	// Name of the contract in L1 (e.g. "PolygonZkEVMDeployer", "PolygonZkEVMBridge",...)
	ContractName string `json:"contractName,omitempty"`
}

func TestExportDataToLayer2Genesis(t *testing.T) {

	t.SkipNow()

	// export FX_DATA_DIR=${HOME}/.fxcore/data L2_GENESIS_FILE=${HOME}/.l2/config/genesis.json

	fxtypes.SetConfig(false)
	fxDataDir := os.Getenv("FX_DATA_DIR")
	require.NotEmptyf(t, fxDataDir, "env var FX_DATA_DIR is empty")
	db, err := dbm.NewDB("application", dbm.GoLevelDBBackend, os.ExpandEnv(fxDataDir))
	require.NoError(t, err)

	l2GenesisFile := os.Getenv("L2_GENESIS_FILE")
	require.NotEmptyf(t, l2GenesisFile, "env var L2_GENESIS_FILE is empty")

	open, err := os.Open(os.ExpandEnv(l2GenesisFile))
	require.NoError(t, err)
	l2GenesisBytes, err := io.ReadAll(open)
	require.NoError(t, err)
	l2GenesisData := &genesisFromJSON{}
	err = json.Unmarshal(l2GenesisBytes, l2GenesisData)
	require.NoError(t, err)

	makeEncodingConfig := app.MakeEncodingConfig()

	myApp := app.New(log.NewFilter(log.NewTMLogger(os.Stdout), log.AllowAll()),
		db, nil, false, map[int64]bool{}, fxtypes.GetDefaultNodeHome(), 0,
		makeEncodingConfig, app.EmptyAppOptions{})

	err = myApp.LoadLatestVersion()
	require.NoError(t, err)

	ctx := newContext(t, myApp)

	var ethAccountCount int
	myApp.AccountKeeper.IterateAccounts(ctx, func(account types.AccountI) bool {
		ethAccount, ok := account.(ethermint.EthAccountI)
		if !ok {
			return false
		}
		ethAccountCount++

		bytecode := myApp.EvmKeeper.GetCode(ctx, ethAccount.GetCodeHash())

		ethAddress := ethAccount.EthAddress()
		balance := myApp.BankKeeper.GetBalance(ctx, ethAddress.Bytes(), fxtypes.DefaultDenom)
		t.Logf("address: %s, balance: %s", ethAccount.EthAddress().String(), balance.String())
		storage := make(map[string]string)
		myApp.EvmKeeper.ForEachStorage(ctx, ethAccount.EthAddress(), func(key, value common.Hash) bool {
			storage[key.String()] = value.String()
			return true
		})
		l2GenesisData.Genesis = append(l2GenesisData.Genesis, genesisAccountFromJSON{
			Balance:      "0",
			Nonce:        strconv.Itoa(int(ethAccount.GetSequence())),
			Address:      ethAddress.String(),
			Bytecode:     "0x" + hex.EncodeToString(bytecode),
			Storage:      storage,
			ContractName: "",
		})
		return false
	})

	// update the genesis file
	l2GenesisBytes, err = json.MarshalIndent(l2GenesisData, "", " ")
	require.NoError(t, err)
	err = os.WriteFile(os.ExpandEnv(l2GenesisFile), l2GenesisBytes, 0644)
	require.NoError(t, err)
	t.Logf("ethAccountCount: %d", ethAccountCount)
}
