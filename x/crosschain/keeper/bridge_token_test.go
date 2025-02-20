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
				for _, k := range suite.App.CrosschainKeepers.ToSlice() {
					if k.ModuleName() != suite.chainName {
						msg := types.MsgBridgeTokenClaim{Name: "Test Token", Symbol: fxtypes.DefaultSymbol, Decimals: 18, TokenContract: helpers.GenExternalAddr(k.ModuleName())}
						suite.Require().NoError(k.AddBridgeTokenExecuted(suite.Ctx, &msg))
						break
					}
					suite.Require().Fail("only one crosschain keeper")
				}
				msg.Symbol = fxtypes.DefaultSymbol
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
			has, err := suite.App.Erc20Keeper.HasERC20Token(suite.Ctx, baseDenom)
			suite.Require().NoError(err)
			suite.True(has)
		})
	}
}
