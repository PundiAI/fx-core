package crosschain_test

import (
	"math/big"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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
	bridgeFeeSuite helpers.BridgeFeeSuite
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
	suite.bridgeFeeSuite = helpers.NewBridgeFeeSuite(suite.Require(), suite.App.EvmKeeper)

	chainNames := fxtypes.GetSupportChains()
	suite.chainName = chainNames[tmrand.Intn(len(chainNames))]

	suite.App.CrosschainKeepers.GetKeeper(suite.chainName).
		SetLastObservedBlockHeight(suite.Ctx, 100, 10)
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

func (suite *CrosschainPrecompileTestSuite) GetERC20TokenKeeper() contract.ERC20TokenKeeper {
	return contract.NewERC20TokenKeeper(suite.App.EvmKeeper)
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

func (suite *CrosschainPrecompileTestSuite) AddBridgeToken(symbolOrAddr string, isNativeOrOrigin bool) {
	keeper := suite.App.Erc20Keeper
	var erc20Token erc20types.ERC20Token
	var err error
	if isNativeOrOrigin {
		erc20Token, err = keeper.RegisterNativeCoin(suite.Ctx, symbolOrAddr, symbolOrAddr, 18)
	} else {
		erc20Token, err = keeper.RegisterNativeERC20(suite.Ctx, common.HexToAddress(symbolOrAddr))
	}
	suite.Require().NoError(err)
	if symbolOrAddr == fxtypes.DefaultSymbol {
		isNativeOrOrigin = false
	}
	err = keeper.AddBridgeToken(suite.Ctx, erc20Token.Denom, suite.chainName, helpers.GenExternalAddr(suite.chainName), isNativeOrOrigin)
	suite.Require().NoError(err)
}

func (suite *CrosschainPrecompileTestSuite) DepositBridgeToken(erc20token erc20types.ERC20Token, amount sdkmath.Int) {
	minter := common.BytesToAddress(authtypes.NewModuleAddress(erc20types.ModuleName).Bytes())
	_, err := suite.GetERC20TokenKeeper().Mint(suite.Ctx, erc20token.GetERC20Contract(), minter, suite.signer.Address(), big.NewInt(100))
	suite.Require().NoError(err)
	suite.MintTokenToModule(erc20types.ModuleName, sdk.NewCoin(erc20token.Denom, amount))
}

func (suite *CrosschainPrecompileTestSuite) Quote(denom string) {
	quoteQuoteInput := contract.IBridgeFeeQuoteQuoteInput{
		Cap:       0,
		GasLimit:  21000,
		Expiry:    uint64(time.Now().Add(time.Hour).Unix()),
		ChainName: contract.MustStrToByte32(suite.chainName),
		TokenName: contract.MustStrToByte32(denom),
		Amount:    big.NewInt(1),
	}

	// add token if not exist
	tokens, err := suite.bridgeFeeSuite.GetTokens(suite.Ctx, quoteQuoteInput.ChainName)
	suite.Require().NoError(err)
	found := false
	for _, token := range tokens {
		if token == quoteQuoteInput.TokenName {
			found = true
		}
	}
	if !found {
		_, err = suite.bridgeFeeSuite.AddToken(suite.Ctx, quoteQuoteInput.ChainName, []common.Hash{quoteQuoteInput.TokenName})
		suite.Require().NoError(err)
	}

	suite.bridgeFeeSuite.Quote(suite.Ctx, quoteQuoteInput)
}
