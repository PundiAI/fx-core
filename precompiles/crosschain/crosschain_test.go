package crosschain_test

import (
	"fmt"
	"math/big"
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

func (suite *CrosschainPrecompileTestSuite) TestContract_Crosschain() {
	testCases := []struct {
		name                       string
		malleate                   func() *erc20types.ERC20Token
		transferAmount             *big.Int
		erc20ModuleAmount          sdkmath.Int // default base denom amount
		crosschainModuleAmount     sdkmath.Int // default bridge denom amount
		crosschainModuleBaseAmount sdkmath.Int
		chainNameAmount            sdkmath.Int // default bridge denom amount
	}{
		{
			name: "native coin",
			malleate: func() *erc20types.ERC20Token {
				bridgeToken := suite.AddBridgeToken(helpers.NewRandSymbol(), true)

				suite.Quote(bridgeToken.Denom)

				suite.AddNativeCoinToEVM(bridgeToken.Denom, sdkmath.NewInt(100))

				return suite.GetERC20Token(bridgeToken.Denom)
			},
			transferAmount:             big.NewInt(2),
			erc20ModuleAmount:          sdkmath.NewInt(98),
			crosschainModuleAmount:     sdkmath.NewInt(98),
			crosschainModuleBaseAmount: sdkmath.NewInt(0),
			chainNameAmount:            sdkmath.NewInt(0),
		},
		{
			name: "native erc20",
			malleate: func() *erc20types.ERC20Token {
				erc20TokenAddr := suite.erc20TokenSuite.DeployERC20Token(suite.Ctx, suite.signer.Address(), helpers.NewRandSymbol())
				bridgeToken := suite.AddBridgeToken(erc20TokenAddr.String(), false)

				suite.Quote(bridgeToken.Denom)

				suite.AddNativeERC20ToEVM(bridgeToken.Denom, sdkmath.NewInt(100))

				return suite.GetERC20Token(bridgeToken.Denom)
			},
			transferAmount:             big.NewInt(2),
			erc20ModuleAmount:          sdkmath.NewInt(0),
			crosschainModuleAmount:     sdkmath.NewInt(0),
			crosschainModuleBaseAmount: sdkmath.NewInt(2),
			chainNameAmount:            sdkmath.NewInt(2),
		},
		{
			name: "IBC Token",
			malleate: func() *erc20types.ERC20Token {
				bridgeToken := suite.AddBridgeToken(helpers.NewRandSymbol(), true, true)

				suite.Quote(bridgeToken.Denom)

				suite.AddNativeCoinToEVM(bridgeToken.Denom, sdkmath.NewInt(100), true)

				return suite.GetERC20Token(bridgeToken.Denom)
			},
			transferAmount:             big.NewInt(2),
			erc20ModuleAmount:          sdkmath.NewInt(98),
			crosschainModuleAmount:     sdkmath.NewInt(0),
			crosschainModuleBaseAmount: sdkmath.NewInt(2),
			chainNameAmount:            sdkmath.NewInt(2),
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			erc20Token := tc.malleate()

			erc20TokenSuite := suite.erc20TokenSuite.WithContract(erc20Token.GetERC20Contract())
			erc20TokenSuite.Approve(suite.Ctx, suite.signer.Address(), suite.crosschainAddr, big.NewInt(2))

			txResponse := suite.Crosschain(suite.Ctx, nil, suite.signer.Address(),
				contract.CrosschainArgs{
					Token:   erc20Token.GetERC20Contract(),
					Receipt: helpers.GenExternalAddr(suite.chainName),
					Amount:  big.NewInt(1),
					Fee:     big.NewInt(1),
					Target:  contract.MustStrToByte32(suite.chainName),
					Memo:    "",
				},
			)
			suite.NotNil(txResponse)
			suite.GreaterOrEqual(len(txResponse.Logs), 2)

			balance := erc20TokenSuite.BalanceOf(suite.Ctx, suite.signer.Address())
			suite.Equal(big.NewInt(98), balance)

			bridgeToken := suite.GetBridgeToken(erc20Token.Denom)
			bridgeCoin := sdk.NewCoin(bridgeToken.BridgeDenom(), tc.crosschainModuleAmount)
			suite.AssertBalance(authtypes.NewModuleAddress(crosschaintypes.ModuleName), bridgeCoin)

			baseCoin := sdk.NewCoin(bridgeToken.Denom, tc.crosschainModuleBaseAmount)
			suite.AssertBalance(authtypes.NewModuleAddress(crosschaintypes.ModuleName), baseCoin)

			bridgeCoin = sdk.NewCoin(bridgeToken.BridgeDenom(), tc.chainNameAmount)
			suite.AssertBalance(authtypes.NewModuleAddress(suite.chainName), bridgeCoin)

			baseCoin = sdk.NewCoin(erc20Token.Denom, sdkmath.NewInt(0))
			suite.AssertBalance(authtypes.NewModuleAddress(suite.chainName), baseCoin)

			bridgeCoin = sdk.NewCoin(bridgeToken.BridgeDenom(), sdkmath.NewInt(0))
			suite.AssertBalance(authtypes.NewModuleAddress(erc20types.ModuleName), bridgeCoin)

			baseCoin = sdk.NewCoin(erc20Token.Denom, tc.erc20ModuleAmount)
			suite.AssertBalance(authtypes.NewModuleAddress(erc20types.ModuleName), baseCoin)
		})
	}
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
