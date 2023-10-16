package app_test

import (
	"encoding/json"
	"io"
	"math/big"
	"os"
	"strconv"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	ethermint "github.com/evmos/ethermint/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v5/app"
	fxtypes "github.com/functionx/fx-core/v5/types"
)

const (
	// todo
	l2StakingContractAddr = "0x"
	// todo
	l2StakingContractByteCode = "0x"
)

var (
	ethPublicKeyType = new(ethsecp256k1.PubKey).Type()

	coinOne = sdk.NewDecFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	_       = coinOne
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

	myApp, ctx, l2GenesisData := buildAppContext(t)

	addrCanMigrateMap := make(map[string]bool)
	valMap := make(map[string]stakingtypes.Validator)
	burnBondedAmount := sdk.ZeroInt()
	burnNotBondedAmount := sdk.ZeroInt()
	userDelegateAmountMap := make(map[string]sdkmath.Int)

	myApp.StakingKeeper.IterateAllDelegations(ctx, func(delegation stakingtypes.Delegation) (stop bool) {
		delegateAddr := delegation.GetDelegatorAddr().String()
		can, found := addrCanMigrateMap[delegateAddr]
		if !found {
			can = canMigrate(t, ctx, myApp, delegation.GetDelegatorAddr())
			addrCanMigrateMap[delegateAddr] = can
		}
		if !can {
			return false
		}

		valAddrStr := delegation.GetValidatorAddr().String()
		val, foundVal := valMap[valAddrStr]
		if !foundVal {
			val, foundVal = myApp.StakingKeeper.GetValidator(ctx, delegation.GetValidatorAddr())
			require.Truef(t, foundVal, "validator not found: %s", valAddrStr)
			valMap[valAddrStr] = val
		}

		unBondAmount, err := myApp.StakingKeeper.Unbond(ctx, delegation.GetDelegatorAddr(), delegation.GetValidatorAddr(), delegation.GetShares())
		require.NoError(t, err)

		mapDelegateAmount, found := userDelegateAmountMap[delegateAddr]
		if !found {
			userDelegateAmountMap[delegateAddr] = unBondAmount
		} else {
			userDelegateAmountMap[delegateAddr] = mapDelegateAmount.Add(unBondAmount)
		}

		if val.IsBonded() {
			burnBondedAmount = burnBondedAmount.Add(unBondAmount)
		} else {
			burnNotBondedAmount = burnNotBondedAmount.Add(unBondAmount)
		}
		return false
	})

	t.Logf("totalDelegateAddr: %d", len(userDelegateAmountMap))
	if !burnBondedAmount.IsZero() {
		err := myApp.BankKeeper.BurnCoins(ctx, stakingtypes.BondedPoolName, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, burnBondedAmount)))
		require.NoError(t, err)
	}

	if !burnNotBondedAmount.IsZero() {
		err := myApp.BankKeeper.BurnCoins(ctx, stakingtypes.NotBondedPoolName, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, burnNotBondedAmount)))
		require.NoError(t, err)
	}

	storage := make(map[string]string)
	for addr, amount := range userDelegateAmountMap {
		accAddr := sdk.MustAccAddressFromBech32(addr)
		ethAddress := common.BytesToAddress(accAddr)
		sequence, err := myApp.AccountKeeper.GetSequence(ctx, accAddr)
		require.NoErrorf(t, err, "get sequence error: %s", addr)
		l2GenesisData.Genesis = append(l2GenesisData.Genesis, genesisAccountFromJSON{
			Balance: "0",
			Nonce:   strconv.Itoa(int(sequence)),
			Address: ethAddress.String(),
		})
		storage[ethAddress.String()] = "0x" + common.Bytes2Hex(amount.BigInt().Bytes())
	}

	l2GenesisData.Genesis = append(l2GenesisData.Genesis, genesisAccountFromJSON{
		Balance:      "0",
		Nonce:        "0",
		Address:      l2StakingContractAddr,
		Bytecode:     l2StakingContractByteCode,
		Storage:      storage,
		ContractName: "L2 Staking Contract",
	})
}

func buildAppContext(t *testing.T) (*app.App, sdk.Context, *genesisFromJSON) {
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
	return myApp, ctx, l2GenesisData
}

func canMigrate(t *testing.T, ctx sdk.Context, myApp *app.App, addr sdk.AccAddress) bool {
	accountI := myApp.AccountKeeper.GetAccount(ctx, addr)
	require.NotNilf(t, accountI, "account not found: %s", addr.String())
	_, ok := accountI.(*ethermint.EthAccount)
	if ok {
		return false
	}

	pubKey := accountI.GetPubKey()
	if pubKey == nil {
		return accountI.GetSequence() > 0
	}

	return pubKey.Type() == ethPublicKeyType
}
