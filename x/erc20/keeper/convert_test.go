package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func (suite *KeeperTestSuite) TestKeeper_ConvertCoin() {
	testCases := []struct {
		name          string
		symbol        string
		isNativeErc20 bool
		err           error
	}{
		{
			name:   "success - origin token",
			symbol: fxtypes.DefaultSymbol,
			err:    nil,
		},
		{
			name:   "success - native token",
			symbol: helpers.NewRandSymbol(),
			err:    nil,
		},
		{
			name:          "success - erc20 token",
			symbol:        helpers.NewRandSymbol(),
			isNativeErc20: true,
			err:           nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			amount := helpers.NewRandAmount()
			symbolOrAddr := tc.symbol
			if tc.isNativeErc20 {
				tokenAddr := suite.erc20TokenSuite.DeployERC20Token(suite.Ctx, suite.signer.Address(), tc.symbol)
				erc20ModuleAddr := common.BytesToAddress(authtypes.NewModuleAddress(erc20types.ModuleName).Bytes())
				suite.erc20TokenSuite.WithContract(tokenAddr).
					Mint(suite.Ctx, suite.signer.Address(), erc20ModuleAddr, amount.BigInt())
				symbolOrAddr = tokenAddr.String()
			}
			bridgeToken := suite.AddBridgeToken(ethtypes.ModuleName, symbolOrAddr, !tc.isNativeErc20)

			coin := sdk.NewCoin(bridgeToken.Denom, amount)
			sender := helpers.GenAccAddress()
			suite.MintToken(sender, coin)

			receiver := helpers.GenHexAddress()
			erc20Addr, err := suite.GetKeeper().ConvertCoin(suite.Ctx, suite.App.EvmKeeper, sender, receiver, coin)
			if tc.err != nil {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, tc.err)
				return
			}
			suite.Require().NoError(err)
			balance := suite.erc20TokenSuite.WithContract(common.HexToAddress(erc20Addr)).
				BalanceOf(suite.Ctx, receiver)
			suite.Require().Equal(coin.Amount.String(), balance.String())
		})
	}
}
