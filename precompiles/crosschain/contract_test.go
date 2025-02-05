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
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/precompiles/crosschain"
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
	bridgeFeeSuite  helpers.BridgeFeeSuite
	erc20TokenSuite helpers.ERC20TokenSuite
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
	suite.erc20TokenSuite = helpers.NewERC20Suite(suite.Require(), suite.signer, suite.App.EvmKeeper)

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

func (suite *CrosschainPrecompileTestSuite) GetERC20Token(baseDenom string) *erc20types.ERC20Token {
	erc20token, err := suite.App.Erc20Keeper.GetERC20Token(suite.Ctx, baseDenom)
	suite.Require().NoError(err)
	return &erc20token
}

func (suite *CrosschainPrecompileTestSuite) GetBridgeToken(baseDenom string) erc20types.BridgeToken {
	bridgeToken, err := suite.App.Erc20Keeper.GetBridgeToken(suite.Ctx, suite.chainName, baseDenom)
	suite.Require().NoError(err)
	return bridgeToken
}

func (suite *CrosschainPrecompileTestSuite) AddBridgeToken(symbolOrAddr string, isNativeCoin bool, isIBC ...bool) erc20types.BridgeToken {
	keeper := suite.App.Erc20Keeper
	var baseDenom string
	isNative := false
	if symbolOrAddr == fxtypes.LegacyFXDenom {
		baseDenom = fxtypes.FXDenom
	} else if isNativeCoin || symbolOrAddr == fxtypes.DefaultSymbol {
		erc20Token, err := keeper.RegisterNativeCoin(suite.Ctx, symbolOrAddr, symbolOrAddr, 18)
		suite.Require().NoError(err)
		baseDenom = erc20Token.Denom
	} else {
		isNative = true
		erc20Token, err := keeper.RegisterNativeERC20(suite.Ctx, common.HexToAddress(symbolOrAddr))
		suite.Require().NoError(err)
		baseDenom = erc20Token.Denom
	}
	if len(isIBC) > 0 && isIBC[0] {
		isNative = true
	}
	err := keeper.AddBridgeToken(suite.Ctx, baseDenom, suite.chainName, helpers.GenExternalAddr(suite.chainName), isNative)
	suite.Require().NoError(err)
	return suite.GetBridgeToken(baseDenom)
}

func (suite *CrosschainPrecompileTestSuite) AddNativeERC20ToEVM(baseDenom string, amount sdkmath.Int) {
	minter := common.BytesToAddress(authtypes.NewModuleAddress(erc20types.ModuleName).Bytes())
	erc20token := suite.GetERC20Token(baseDenom)
	if erc20token.IsNativeERC20() {
		minter = suite.signer.Address()
	}
	suite.erc20TokenSuite.Mint(suite.Ctx, minter, suite.signer.Address(), amount.BigInt())
	balance := suite.erc20TokenSuite.BalanceOf(suite.Ctx, suite.signer.Address())
	suite.Equal(balance.String(), amount.BigInt().String())
}

func (suite *CrosschainPrecompileTestSuite) AddNativeCoinToEVM(baseDenom string, amount sdkmath.Int, isIBC ...bool) {
	suite.AddNativeERC20ToEVM(baseDenom, amount)

	suite.MintTokenToModule(erc20types.ModuleName, sdk.NewCoin(baseDenom, amount))

	if len(isIBC) == 0 || !isIBC[0] {
		bridgeToken := suite.GetBridgeToken(baseDenom)
		suite.MintTokenToModule(crosschaintypes.ModuleName, sdk.NewCoin(bridgeToken.BridgeDenom(), amount))
	}
}

func (suite *CrosschainPrecompileTestSuite) NewBridgeCallArgs(erc20Contract common.Address, amount *big.Int) contract.BridgeCallArgs {
	args := contract.BridgeCallArgs{
		DstChain: suite.chainName,
		Refund:   suite.signer.Address(),
		Tokens:   []common.Address{erc20Contract},
		Amounts:  []*big.Int{amount},
		To:       suite.signer.Address(),
		Data:     []byte{},
		QuoteId:  big.NewInt(1),
		GasLimit: big.NewInt(0),
		Memo:     []byte{},
	}
	if contract.IsZeroEthAddress(erc20Contract) {
		args.Tokens = []common.Address{}
	}
	if amount == nil {
		args.Amounts = []*big.Int{}
	}
	return args
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

func (suite *CrosschainPrecompileTestSuite) executeClaim(claim crosschaintypes.ExternalClaim) *evmtypes.MsgEthereumTxResponse {
	keeper := suite.App.CrosschainKeepers.GetKeeper(suite.chainName)
	err := keeper.SavePendingExecuteClaim(suite.Ctx, claim)
	suite.Require().NoError(err)

	txResponse := suite.ExecuteClaim(suite.Ctx, suite.signer.Address(),
		contract.ExecuteClaimArgs{Chain: suite.chainName, EventNonce: big.NewInt(int64(claim.GetEventNonce()))},
	)
	suite.NotNil(txResponse)
	suite.Empty(txResponse.VmError)

	event, err := crosschain.NewExecuteClaimABI().
		UnpackEvent(txResponse.Logs[len(txResponse.Logs)-1].ToEthereum())
	suite.Require().NoError(err)
	suite.Equal(suite.GetSender(), event.Sender)
	suite.Equal(claim.GetEventNonce(), event.EventNonce.Uint64())
	suite.Equal(claim.GetChainName(), event.Chain)
	suite.Empty(event.ErrReason)
	return txResponse
}
