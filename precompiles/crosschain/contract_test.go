package crosschain_test

import (
	"math/big"
	"testing"

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

	// todo: need refactor
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
	return suite.BaseSuite.SetOracle(suite.chainName, online)
}

func (suite *CrosschainPrecompileTestSuite) GetERC20Token(baseDenom string) *erc20types.ERC20Token {
	return suite.BaseSuite.GetERC20Token(baseDenom)
}

func (suite *CrosschainPrecompileTestSuite) GetBridgeToken(baseDenom string) erc20types.BridgeToken {
	return suite.BaseSuite.GetBridgeToken(suite.chainName, baseDenom)
}

func (suite *CrosschainPrecompileTestSuite) AddBridgeToken(symbolOrAddr string, isNativeCoin bool, isIBC ...bool) erc20types.BridgeToken {
	return suite.BaseSuite.AddBridgeToken(suite.chainName, symbolOrAddr, isNativeCoin, isIBC...)
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
	suite.bridgeFeeSuite.MockQuote(suite.Ctx, suite.chainName, denom)
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
