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
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func TestCrosschainABI(t *testing.T) {
	crosschainABI := crosschain.NewCrosschainABI()

	require.Len(t, crosschainABI.Method.Inputs, 6)
	require.Len(t, crosschainABI.Method.Outputs, 1)

	require.Len(t, crosschainABI.Event.Inputs, 8)
}

func (suite *CrosschainPrecompileTestSuite) TestContract_Crosschain() {
	symbol := "USDT"
	suite.AddBridgeToken(symbol, true)

	baseDenom := strings.ToLower(symbol)
	suite.Quote(baseDenom)

	erc20token, err := suite.App.Erc20Keeper.GetERC20Token(suite.Ctx, baseDenom)
	suite.Require().NoError(err)

	amount := sdkmath.NewInt(100)
	suite.DepositBridgeToken(erc20token, amount)

	_, err = suite.GetERC20TokenKeeper().Approve(
		suite.Ctx, erc20token.GetERC20Contract(), suite.signer.Address(), suite.crosschainAddr, big.NewInt(2))
	suite.Require().NoError(err)

	txResponse := suite.Crosschain(suite.Ctx, nil, suite.signer.Address(),
		contract.CrosschainArgs{
			Token:   erc20token.GetERC20Contract(),
			Receipt: helpers.GenExternalAddr(suite.chainName),
			Amount:  big.NewInt(1),
			Fee:     big.NewInt(1),
			Target:  contract.MustStrToByte32(suite.chainName),
			Memo:    "",
		},
	)
	suite.NotNil(txResponse)
	suite.GreaterOrEqual(len(txResponse.Logs), 2)

	balance, err := suite.GetERC20TokenKeeper().
		BalanceOf(suite.Ctx, erc20token.GetERC20Contract(), suite.signer.Address())
	suite.Require().NoError(err)
	suite.Equal(big.NewInt(98), balance)

	baseCoin := sdk.NewCoin(baseDenom, sdkmath.NewInt(98))
	suite.AssertBalance(authtypes.NewModuleAddress(erc20types.ModuleName), baseCoin)

	bridgeToken, err := suite.App.Erc20Keeper.GetBridgeToken(suite.Ctx, suite.chainName, baseDenom)
	suite.Require().NoError(err)
	bridgeCoin := sdk.NewCoin(bridgeToken.BridgeDenom(), sdkmath.NewInt(2))
	suite.AssertBalance(authtypes.NewModuleAddress(suite.chainName), bridgeCoin)
}

func (suite *CrosschainPrecompileTestSuite) TestContract_Crosschain_OriginToken() {
	suite.AddBridgeToken(fxtypes.DefaultSymbol, true)
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
