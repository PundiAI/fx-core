package crosschain_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/contract"
	testscontract "github.com/pundiai/fx-core/v8/tests/contract"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
)

type CrosschainPrecompileTestSuite struct {
	helpers.BaseSuite

	signer         *helpers.Signer
	crosschainAddr common.Address
	chainName      string

	helpers.CrosschainPrecompileSuite
}

func TestCrosschainPrecompileTestSuite(t *testing.T) {
	suite.Run(t, &CrosschainPrecompileTestSuite{
		crosschainAddr: common.HexToAddress(contract.CrosschainAddress),
	})
}

func TestCrosschainPrecompileTestSuite_Contract(t *testing.T) {
	suite.Run(t, new(CrosschainPrecompileTestSuite))
}

func (suite *CrosschainPrecompileTestSuite) SetupTest() {
	suite.BaseSuite.SetupTest()
	suite.Ctx = suite.Ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(fxtypes.DefaultDenom, sdkmath.OneInt())))
	suite.Ctx = suite.Ctx.WithBlockGasMeter(storetypes.NewGasMeter(1e18))

	suite.signer = suite.AddTestSigner(10_000)

	if !suite.IsCallPrecompile() {
		crosschainContract, err := suite.App.EvmKeeper.DeployContract(suite.Ctx, suite.signer.Address(), contract.MustABIJson(testscontract.CrosschainTestMetaData.ABI), contract.MustDecodeHex(testscontract.CrosschainTestMetaData.Bin))
		suite.Require().NoError(err)
		suite.crosschainAddr = crosschainContract
		account := suite.App.EvmKeeper.GetAccount(suite.Ctx, crosschainContract)
		suite.Require().True(account.IsContract())
	}

	suite.CrosschainPrecompileSuite = helpers.NewCrosschainPrecompileSuite(suite.Require(), suite.signer, suite.App.EvmKeeper, suite.crosschainAddr)

	chainNames := fxtypes.GetSupportChains()
	suite.chainName = chainNames[tmrand.Intn(len(chainNames))]
}

func (suite *CrosschainPrecompileTestSuite) SetupSubTest() {
	suite.SetupTest()
}

func (suite *CrosschainPrecompileTestSuite) GetSender() common.Address {
	if suite.IsCallPrecompile() {
		return suite.signer.Address()
	}
	return suite.crosschainAddr
}

func (suite *CrosschainPrecompileTestSuite) IsCallPrecompile() bool {
	return suite.crosschainAddr.String() == contract.CrosschainAddress
}

func (suite *CrosschainPrecompileTestSuite) SetOracle(online bool) crosschaintypes.Oracle {
	oracle := crosschaintypes.Oracle{
		OracleAddress:     helpers.GenAccAddress().String(),
		BridgerAddress:    helpers.GenAccAddress().String(),
		ExternalAddress:   helpers.GenExternalAddr(suite.chainName),
		DelegateAmount:    sdkmath.NewInt(1e18).MulRaw(1000),
		StartHeight:       1,
		Online:            online,
		DelegateValidator: sdk.ValAddress(helpers.GenAccAddress()).String(),
		SlashTimes:        0,
	}
	keeper := suite.App.CrosschainKeepers.GetKeeper(suite.chainName)
	keeper.SetOracle(suite.Ctx, oracle)
	keeper.SetOracleAddrByExternalAddr(suite.Ctx, oracle.ExternalAddress, oracle.GetOracle())
	keeper.SetOracleAddrByBridgerAddr(suite.Ctx, oracle.GetBridger(), oracle.GetOracle())
	return oracle
}

func (suite *CrosschainPrecompileTestSuite) AddBridgeToken(symbolOrAddr string, isNative bool) common.Address {
	keeper := suite.App.Erc20Keeper
	var erc20Token erc20types.ERC20Token
	var err error
	if isNative {
		erc20Token, err = keeper.RegisterNativeCoin(suite.Ctx, symbolOrAddr, symbolOrAddr, 18)
	} else {
		erc20Token, err = keeper.RegisterNativeERC20(suite.Ctx, common.HexToAddress(symbolOrAddr))
	}
	suite.Require().NoError(err)
	err = keeper.AddBridgeToken(suite.Ctx, erc20Token.Denom, suite.chainName, erc20Token.Erc20Address, isNative)
	suite.Require().NoError(err)
	return common.HexToAddress(erc20Token.Erc20Address)
}
