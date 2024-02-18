package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	ethermint "github.com/evmos/ethermint/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	migratetypes "github.com/functionx/fx-core/v7/x/migrate/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx              sdk.Context
	app              *app.App
	secp256k1PrivKey cryptotypes.PrivKey
	queryClient      migratetypes.QueryClient
	govAddr          string
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// Test helpers
func (suite *KeeperTestSuite) SetupTest() {
	// init app
	initBalances := sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(20000))
	validator, genesisAccounts, balances := helpers.GenerateGenesisValidator(3, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, initBalances)))
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), validator, genesisAccounts, balances...)

	suite.ctx = suite.app.NewContext(false, tmproto.Header{Height: 1, ChainID: "fxcore", ProposerAddress: validator.Validators[0].Address, Time: time.Now().UTC()})
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(fxtypes.DefaultDenom, sdkmath.OneInt())))
	suite.ctx = suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(1e18))

	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	migratetypes.RegisterQueryServer(queryHelperEvm, suite.app.MigrateKeeper)
	suite.queryClient = migratetypes.NewQueryClient(queryHelperEvm)

	// account key
	suite.secp256k1PrivKey = secp256k1.GenPrivKey()
	acc := &ethermint.EthAccount{
		BaseAccount: authtypes.NewBaseAccount(suite.secp256k1PrivKey.PubKey().Address().Bytes(), nil, 0, 0),
		CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
	}
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

	amount := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18)))
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.NewCoins(amount))
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoins(amount))
	suite.Require().NoError(err)

	// update staking unbonding time
	stakingParams := suite.app.StakingKeeper.GetParams(suite.ctx)
	stakingParams.UnbondingTime = 5 * time.Minute
	suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams)

	suite.govAddr = authtypes.NewModuleAddress(govtypes.ModuleName).String()
}

func (suite *KeeperTestSuite) Commit() {
	_ = suite.app.Commit()
	header := suite.ctx.BlockHeader()
	header.Height += 1
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: header,
	})

	// update ctx
	suite.ctx = suite.app.NewContext(false, header)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	migratetypes.RegisterQueryServer(queryHelper, suite.app.MigrateKeeper)
	suite.queryClient = migratetypes.NewQueryClient(queryHelper)
}

func (suite *KeeperTestSuite) GenerateAcc(num int) []cryptotypes.PrivKey {
	keys := make([]cryptotypes.PrivKey, 0, num)
	for i := 0; i < num; i++ {
		privateKey := secp256k1.GenPrivKey()
		amount := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100000).Mul(sdkmath.NewInt(1e18)))
		err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.NewCoins(amount))
		suite.Require().NoError(err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, privateKey.PubKey().Address().Bytes(), sdk.NewCoins(amount))
		suite.Require().NoError(err)

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
		amount := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100000).Mul(sdkmath.NewInt(1e18)))
		err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.NewCoins(amount))
		suite.Require().NoError(err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, privateKey.PubKey().Address().Bytes(), sdk.NewCoins(amount))
		suite.Require().NoError(err)

		keys = append(keys, privateKey)
	}
	return keys
}

func (suite *KeeperTestSuite) mintToken(module string, address sdk.AccAddress, amount sdk.Coin) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, module, sdk.NewCoins(amount))
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, module, address.Bytes(), sdk.NewCoins(amount))
	suite.Require().NoError(err)
}
