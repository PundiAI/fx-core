package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
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

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	migratetypes "github.com/functionx/fx-core/v8/x/migrate/types"
)

type KeeperTestSuite struct {
	helpers.BaseSuite

	secp256k1PrivKey cryptotypes.PrivKey
	queryClient      migratetypes.QueryClient
	govAddr          string
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.BaseSuite.MintValNumber = 3
	suite.BaseSuite.SetupTest()

	suite.Ctx = suite.Ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(fxtypes.DefaultDenom, sdkmath.OneInt())))
	suite.Ctx = suite.Ctx.WithBlockGasMeter(storetypes.NewGasMeter(1e18))

	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.Ctx, suite.App.InterfaceRegistry())
	migratetypes.RegisterQueryServer(queryHelperEvm, suite.App.MigrateKeeper)
	suite.queryClient = migratetypes.NewQueryClient(queryHelperEvm)

	// account key
	suite.secp256k1PrivKey = secp256k1.GenPrivKey()
	acc := &ethermint.EthAccount{
		BaseAccount: authtypes.NewBaseAccount(suite.secp256k1PrivKey.PubKey().Address().Bytes(), nil, suite.App.AccountKeeper.NextAccountNumber(suite.Ctx), 0),
		CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
	}
	suite.App.AccountKeeper.SetAccount(suite.Ctx, acc)

	amount := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18)))
	err := suite.App.BankKeeper.MintCoins(suite.Ctx, minttypes.ModuleName, sdk.NewCoins(amount))
	suite.Require().NoError(err)
	err = suite.App.BankKeeper.SendCoinsFromModuleToAccount(suite.Ctx, minttypes.ModuleName, suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoins(amount))
	suite.Require().NoError(err)

	// update staking unbonding time
	stakingParams, err := suite.App.StakingKeeper.GetParams(suite.Ctx)
	suite.Require().NoError(err)
	stakingParams.UnbondingTime = 5 * time.Minute
	err = suite.App.StakingKeeper.SetParams(suite.Ctx, stakingParams)
	suite.Require().NoError(err)

	suite.govAddr = authtypes.NewModuleAddress(govtypes.ModuleName).String()
}

func (suite *KeeperTestSuite) GenerateAcc(num int) []cryptotypes.PrivKey {
	keys := make([]cryptotypes.PrivKey, 0, num)
	for i := 0; i < num; i++ {
		privateKey := secp256k1.GenPrivKey()

		amount := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100000).Mul(sdkmath.NewInt(1e18)))
		suite.MintToken(privateKey.PubKey().Address().Bytes(), amount)

		acc := &ethermint.EthAccount{
			BaseAccount: authtypes.NewBaseAccount(privateKey.PubKey().Address().Bytes(), privateKey.PubKey(), suite.App.AccountKeeper.NextAccountNumber(suite.Ctx), 0),
			CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
		}
		suite.App.AccountKeeper.SetAccount(suite.Ctx, acc)

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
		suite.MintToken(privateKey.PubKey().Address().Bytes(), amount)

		keys = append(keys, privateKey)
	}
	return keys
}

func (suite *KeeperTestSuite) ValStringToVal(addrStr string) sdk.ValAddress {
	addrBytes, err := suite.App.StakingKeeper.ValidatorAddressCodec().StringToBytes(addrStr)
	suite.Require().NoError(err)
	return addrBytes
}
