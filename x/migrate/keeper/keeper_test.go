package keeper_test

import (
	"testing"
	"time"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	ethermint "github.com/tharsis/ethermint/types"

	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/app/helpers"
	upgradesv2 "github.com/functionx/fx-core/app/upgrades/v2"
	fxtypes "github.com/functionx/fx-core/types"
	bsctypes "github.com/functionx/fx-core/x/bsc/types"
	migratetypes "github.com/functionx/fx-core/x/migrate/types"
)

var (
	DevnetPurseDenom = "ibc/B1861D0C2E4BAFA42A61739291975B7663F278FFAF579F83C9C4AD3890D09CA0"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx                 sdk.Context
	app                 *app.App
	clientCtx           client.Context
	secp256k1PrivKey    cryptotypes.PrivKey
	ethsecp256k1PrivKey cryptotypes.PrivKey
	queryClient         migratetypes.QueryClient
	checkTx             bool
	purseBalance        sdk.Int
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// Test helpers
func (suite *KeeperTestSuite) DoSetupTest(t require.TestingT) {
	var err error
	// account key
	suite.secp256k1PrivKey = secp256k1.GenPrivKey()
	suite.ethsecp256k1PrivKey, err = ethsecp256k1.GenerateKey()
	require.NoError(t, err)

	// init app
	initBalances := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := helpers.GenerateGenesisValidator(3, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, initBalances)))
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), validator, genesisAccounts, balances...)

	suite.ctx = suite.app.BaseApp.NewContext(suite.checkTx, tmproto.Header{Height: 1, ChainID: "fxcore", ProposerAddress: validator.Validators[0].Address, Time: time.Now().UTC()})
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(fxtypes.DefaultDenom, sdk.OneInt())))
	suite.ctx = suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(1e18))

	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	migratetypes.RegisterQueryServer(queryHelperEvm, suite.app.MigrateKeeper)
	suite.queryClient = migratetypes.NewQueryClient(queryHelperEvm)

	acc := &ethermint.EthAccount{
		BaseAccount: authtypes.NewBaseAccount(suite.secp256k1PrivKey.PubKey().Address().Bytes(), nil, 0, 0),
		CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
	}
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

	amount := sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1000).Mul(sdk.NewInt(1e18)))
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.NewCoins(amount))
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoins(amount))
	suite.Require().NoError(err)

	if !suite.purseBalance.IsNil() {
		amount := sdk.NewCoin(DevnetPurseDenom, sdk.NewInt(1000).Mul(sdk.NewInt(1e18)))
		err = suite.app.BankKeeper.MintCoins(suite.ctx, bsctypes.ModuleName, sdk.NewCoins(amount))
		suite.Require().NoError(err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, bsctypes.ModuleName, suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoins(amount))
		suite.Require().NoError(err)
	}

	//update staking unbonding time
	stakingParams := suite.app.StakingKeeper.GetParams(suite.ctx)
	stakingParams.UnbondingTime = 5 * time.Minute
	suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams)

	encodingConfig := app.MakeEncodingConfig()
	suite.clientCtx = client.Context{}.WithTxConfig(encodingConfig.TxConfig)

	//update metadata
	err = upgradesv2.UpdateFXMetadata(suite.ctx, suite.app.BankKeeper, suite.app.GetKey(banktypes.StoreKey))
	require.NoError(t, err)
	//migrate account to eth
	upgradesv2.MigrateAccountToEth(suite.ctx, suite.app.AccountKeeper)
	// init logic contract
	for _, contract := range fxtypes.GetInitContracts() {
		require.True(t, len(contract.Code) > 0)
		require.True(t, contract.Address != common.HexToAddress(fxtypes.EmptyEvmAddress))
		err := suite.app.Erc20Keeper.CreateContractWithCode(suite.ctx, contract.Address, contract.Code)
		require.NoError(t, err)
	}

	// register coin
	for _, metadata := range fxtypes.GetMetadata() {
		_, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, metadata)
		require.NoError(t, err)
	}
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.DoSetupTest(suite.T())
}

func (suite *KeeperTestSuite) Commit() {
	_ = suite.app.Commit()
	header := suite.ctx.BlockHeader()
	header.Height += 1
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: header,
	})

	// update ctx
	suite.ctx = suite.app.BaseApp.NewContext(false, header)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	migratetypes.RegisterQueryServer(queryHelper, suite.app.MigrateKeeper)
	suite.queryClient = migratetypes.NewQueryClient(queryHelper)
}

func (suite *KeeperTestSuite) GenerateAcc(num int) []cryptotypes.PrivKey {
	keys := make([]cryptotypes.PrivKey, 0, num)
	for i := 0; i < num; i++ {
		privateKey := secp256k1.GenPrivKey()
		amount := sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(100000).Mul(sdk.NewInt(1e18)))
		err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.NewCoins(amount))
		suite.Require().NoError(err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, privateKey.PubKey().Address().Bytes(), sdk.NewCoins(amount))
		suite.Require().NoError(err)

		if !suite.purseBalance.IsNil() {
			amount := sdk.NewCoin(DevnetPurseDenom, sdk.NewInt(100000).Mul(sdk.NewInt(1e18)))
			err = suite.app.BankKeeper.MintCoins(suite.ctx, bsctypes.ModuleName, sdk.NewCoins(amount))
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, bsctypes.ModuleName, privateKey.PubKey().Address().Bytes(), sdk.NewCoins(amount))
			suite.Require().NoError(err)
		}

		acc := &ethermint.EthAccount{
			BaseAccount: authtypes.NewBaseAccount(privateKey.PubKey().Address().Bytes(), privateKey.PubKey(), 0, 0),
			CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
		}
		suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

		keys = append(keys, privateKey)
	}
	return keys
}

func (suite *KeeperTestSuite) GenerateEthAcc(num int) []cryptotypes.PrivKey {
	keys := make([]cryptotypes.PrivKey, 0, num)
	for i := 0; i < num; i++ {
		privateKey, err := ethsecp256k1.GenerateKey()
		suite.Require().NoError(err)
		amount := sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(100000).Mul(sdk.NewInt(1e18)))
		err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.NewCoins(amount))
		suite.Require().NoError(err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, privateKey.PubKey().Address().Bytes(), sdk.NewCoins(amount))
		suite.Require().NoError(err)

		keys = append(keys, privateKey)
	}
	return keys
}
