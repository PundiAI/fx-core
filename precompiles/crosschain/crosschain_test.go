package crosschain_test

import (
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/precompiles/crosschain"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func TestCrosschainABI(t *testing.T) {
	crosschainABI := crosschain.NewCrosschainABI()

	require.Len(t, crosschainABI.Method.Inputs, 6)
	require.Len(t, crosschainABI.Method.Outputs, 1)

	require.Len(t, crosschainABI.Event.Inputs, 8)
}

func (suite *CrosschainPrecompileTestSuite) TestContract_Crosschain_NativeCoin() {
	symbol := helpers.NewRandSymbol()
	suite.AddBridgeToken(symbol, true)

	baseDenom := strings.ToLower(symbol)
	suite.Quote(baseDenom)

	erc20Contract := suite.GetERC20Token(baseDenom).GetERC20Contract()
	suite.erc20TokenSuite.WithContract(erc20Contract)

	amount := sdkmath.NewInt(100)
	suite.AddNativeCoinToEVM(baseDenom, amount)

	suite.erc20TokenSuite.Approve(suite.Ctx, suite.crosschainAddr, big.NewInt(2))

	txResponse := suite.Crosschain(suite.Ctx, nil, suite.signer.Address(),
		contract.CrosschainArgs{
			Token:   erc20Contract,
			Receipt: helpers.GenExternalAddr(suite.chainName),
			Amount:  big.NewInt(1),
			Fee:     big.NewInt(1),
			Target:  contract.MustStrToByte32(suite.chainName),
			Memo:    "",
		},
	)
	suite.NotNil(txResponse)
	suite.GreaterOrEqual(len(txResponse.Logs), 2)

	balance := suite.erc20TokenSuite.BalanceOf(suite.Ctx, suite.signer.Address())
	suite.Equal(big.NewInt(98), balance)

	baseCoin := sdk.NewCoin(baseDenom, sdkmath.NewInt(98))
	suite.AssertBalance(authtypes.NewModuleAddress(erc20types.ModuleName), baseCoin)

	bridgeToken := suite.GetBridgeToken(baseDenom)
	bridgeCoin := sdk.NewCoin(bridgeToken.BridgeDenom(), sdkmath.NewInt(98))
	suite.AssertBalance(authtypes.NewModuleAddress(crosschaintypes.ModuleName), bridgeCoin)
}

func (suite *CrosschainPrecompileTestSuite) TestContract_Crosschain_NativeERC20() {
	symbol := helpers.NewRandSymbol()

	erc20TokenAddr := suite.erc20TokenSuite.DeployERC20Token(suite.Ctx, symbol)
	suite.AddBridgeToken(erc20TokenAddr.String(), false)

	baseDenom := strings.ToLower(symbol)
	suite.Quote(baseDenom)

	amount := sdkmath.NewInt(100)
	suite.AddNativeERC20ToEVM(baseDenom, amount)

	suite.erc20TokenSuite.Approve(suite.Ctx, suite.crosschainAddr, big.NewInt(2))

	txResponse := suite.Crosschain(suite.Ctx, nil, suite.signer.Address(),
		contract.CrosschainArgs{
			Token:   erc20TokenAddr,
			Receipt: helpers.GenExternalAddr(suite.chainName),
			Amount:  big.NewInt(1),
			Fee:     big.NewInt(1),
			Target:  contract.MustStrToByte32(suite.chainName),
			Memo:    "",
		},
	)
	suite.NotNil(txResponse)
	suite.GreaterOrEqual(len(txResponse.Logs), 2)

	balance := suite.erc20TokenSuite.BalanceOf(suite.Ctx, suite.signer.Address())
	suite.Equal(big.NewInt(98), balance)

	baseCoin := sdk.NewCoin(baseDenom, sdkmath.NewInt(0))
	suite.AssertBalance(authtypes.NewModuleAddress(erc20types.ModuleName), baseCoin)

	bridgeToken := suite.GetBridgeToken(baseDenom)
	bridgeCoin := sdk.NewCoin(bridgeToken.BridgeDenom(), sdkmath.NewInt(2))
	suite.AssertBalance(authtypes.NewModuleAddress(suite.chainName), bridgeCoin)
}

func (suite *CrosschainPrecompileTestSuite) TestContract_Crosschain_IBCToken() {
	symbol := helpers.NewRandSymbol()

	suite.AddBridgeToken(symbol, true, true)

	baseDenom := strings.ToLower(symbol)
	suite.Quote(baseDenom)

	erc20Contract := suite.GetERC20Token(baseDenom).GetERC20Contract()
	suite.erc20TokenSuite.WithContract(erc20Contract)

	amount := sdkmath.NewInt(100)
	suite.AddNativeCoinToEVM(baseDenom, amount, true)

	suite.erc20TokenSuite.Approve(suite.Ctx, suite.crosschainAddr, big.NewInt(2))

	txResponse := suite.Crosschain(suite.Ctx, nil, suite.signer.Address(),
		contract.CrosschainArgs{
			Token:   erc20Contract,
			Receipt: helpers.GenExternalAddr(suite.chainName),
			Amount:  big.NewInt(1),
			Fee:     big.NewInt(1),
			Target:  contract.MustStrToByte32(suite.chainName),
			Memo:    "",
		},
	)
	suite.NotNil(txResponse)
	suite.GreaterOrEqual(len(txResponse.Logs), 2)

	balance := suite.erc20TokenSuite.BalanceOf(suite.Ctx, suite.signer.Address())
	suite.Equal(big.NewInt(98), balance)

	baseCoin := sdk.NewCoin(baseDenom, sdkmath.NewInt(98))
	suite.AssertBalance(authtypes.NewModuleAddress(erc20types.ModuleName), baseCoin)

	bridgeToken := suite.GetBridgeToken(baseDenom)
	bridgeCoin := sdk.NewCoin(bridgeToken.BridgeDenom(), sdkmath.NewInt(2))
	suite.AssertBalance(authtypes.NewModuleAddress(suite.chainName), bridgeCoin)
}

func (suite *CrosschainPrecompileTestSuite) TestContract_Crosschain_OriginToken() {
	suite.AddBridgeToken(fxtypes.DefaultSymbol, false)

	suite.Quote(fxtypes.DefaultDenom)

	balance := suite.Balance(suite.signer.AccAddress())

	txResponse := suite.Crosschain(suite.Ctx, big.NewInt(2), suite.signer.Address(),
		contract.CrosschainArgs{
			Token:   common.Address{},
			Receipt: helpers.GenExternalAddr(suite.chainName),
			Amount:  big.NewInt(1),
			Fee:     big.NewInt(1),
			Target:  contract.MustStrToByte32(suite.chainName),
			Memo:    "",
		},
	)
	suite.NotNil(txResponse)
	suite.Len(txResponse.Logs, 1)

	transferCoin := helpers.NewStakingCoin(2, 0)
	suite.AssertBalance(suite.signer.AccAddress(), balance.Sub(transferCoin)...)
	suite.AssertBalance(authtypes.NewModuleAddress(ethtypes.ModuleName), transferCoin)
}
