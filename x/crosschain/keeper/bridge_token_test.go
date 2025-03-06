package keeper_test

import (
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (suite *KeeperTestSuite) TestKeeper_AddBridgeTokenExecuted() {
	testCases := []struct {
		name    string
		initMsg func(msg *types.MsgBridgeTokenClaim)
		error   error
	}{
		{
			name: "success - add origin token",
			initMsg: func(msg *types.MsgBridgeTokenClaim) {
				msg.Symbol = fxtypes.DefaultSymbol
			},
		},
		{
			name: "success - add origin tokens for multiple chains",
			initMsg: func(msg *types.MsgBridgeTokenClaim) {
				msg.Symbol = fxtypes.DefaultSymbol
				for _, k := range suite.App.CrosschainKeepers.ToSlice() {
					if k.ModuleName() != suite.chainName {
						msg2 := *msg
						msg2.TokenContract = helpers.GenExternalAddr(k.ModuleName())
						suite.Require().NoError(k.AddBridgeTokenExecuted(suite.Ctx, &msg2))
						break
					}
				}
				tokens, err := suite.App.Erc20Keeper.GetBaseBridgeTokens(suite.Ctx, fxtypes.DefaultDenom)
				suite.NoError(err)
				suite.Len(tokens, 1)
			},
		},
		{
			name: "success - add origin token repeatedly",
			initMsg: func(msg *types.MsgBridgeTokenClaim) {
				msg.Symbol = fxtypes.DefaultSymbol
				suite.Require().NoError(suite.Keeper().AddBridgeTokenExecuted(suite.Ctx, msg))
			},
		},
		{
			name: "success - native coin",
			initMsg: func(msg *types.MsgBridgeTokenClaim) {
				msg.Symbol = helpers.NewRandSymbol()
			},
		},
		{
			name: "success - erc20 token",
			initMsg: func(msg *types.MsgBridgeTokenClaim) {
				metadata := fxtypes.NewMetadata(msg.Symbol, msg.Symbol, 18)
				_, err := suite.App.Erc20Keeper.AddERC20Token(suite.Ctx, metadata, helpers.GenHexAddress(), erc20types.OWNER_MODULE)
				suite.NoError(err)
			},
		},
		{
			name: "success - external erc20 token",
			initMsg: func(msg *types.MsgBridgeTokenClaim) {
				metadata := fxtypes.NewMetadata(msg.Symbol, msg.Symbol, 18)
				_, err := suite.App.Erc20Keeper.AddERC20Token(suite.Ctx, metadata, helpers.GenHexAddress(), erc20types.OWNER_EXTERNAL)
				suite.NoError(err)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := types.MsgBridgeTokenClaim{
				Name:          "Test Token",
				Symbol:        helpers.NewRandSymbol(),
				Decimals:      18,
				TokenContract: helpers.GenExternalAddr(suite.chainName),
			}
			tc.initMsg(&msg)
			err := suite.Keeper().AddBridgeTokenExecuted(suite.Ctx, &msg)
			if tc.error != nil {
				suite.Require().Error(err)
				suite.ErrorIs(err, tc.error)
				return
			}
			suite.Require().NoError(err)
			baseDenom, err := suite.App.Erc20Keeper.GetBaseDenom(suite.Ctx, erc20types.NewBridgeDenom(suite.chainName, msg.TokenContract))
			suite.Require().NoError(err)
			suite.Equal(msg.GetBaseDenom(), baseDenom)
			erc20Token, err := suite.App.Erc20Keeper.GetERC20Token(suite.Ctx, baseDenom)
			suite.Require().NoError(err)
			bridgeToken, err := suite.App.Erc20Keeper.GetBridgeToken(suite.Ctx, suite.chainName, baseDenom)
			suite.Require().NoError(err)
			suite.Equal(erc20Token.Denom, bridgeToken.Denom)
			suite.Equal(erc20Token.IsNativeERC20(), bridgeToken.IsNative)
		})
	}
}
