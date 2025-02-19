package keeper_test

import (
	"fmt"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (suite *KeeperTestSuite) TestKeeper_AddBridgeTokenExecuted() {
	testCases := []struct {
		name     string
		malleate func(string) string
		error    func(string) string
	}{
		{
			name: "success - add origin token",
			malleate: func(_ string) string {
				return fxtypes.DefaultSymbol
			},
		},
		{
			name: "success - add origin tokens for multiple chains",
			malleate: func(_ string) string {
				for _, k := range suite.App.CrosschainKeepers.ToSlice() {
					if k.ModuleName() != suite.chainName {
						msg := types.MsgBridgeTokenClaim{Name: "Test Token", Symbol: fxtypes.DefaultSymbol, Decimals: 18, TokenContract: helpers.GenExternalAddr(k.ModuleName())}
						suite.Require().NoError(k.AddBridgeTokenExecuted(suite.Ctx, &msg))
						break
					}
					suite.Require().Fail("only one crosschain keeper")
				}
				return fxtypes.DefaultSymbol
			},
		},
		{
			name: "failed - add origin token multiple more than once",
			malleate: func(_ string) string {
				msg := types.MsgBridgeTokenClaim{Name: "Test Token", Symbol: fxtypes.DefaultSymbol, Decimals: 18, TokenContract: helpers.GenExternalAddr(suite.chainName)}
				suite.Require().NoError(suite.Keeper().AddBridgeTokenExecuted(suite.Ctx, &msg))
				return fxtypes.DefaultSymbol
			},
			error: func(tc string) string {
				return fmt.Sprintf("bridge token %s already exists: invalid request", tc)
			},
		},
		{
			name: "success - native coin",
			malleate: func(_ string) string {
				return helpers.NewRandSymbol()
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tokenContract := helpers.GenExternalAddr(suite.chainName)
			symbol := tc.malleate(tokenContract)
			msg := types.MsgBridgeTokenClaim{Name: "Test Token", Symbol: symbol, Decimals: 18, TokenContract: tokenContract}
			err := suite.Keeper().AddBridgeTokenExecuted(suite.Ctx, &msg)
			if tc.error != nil {
				suite.Require().Error(err)
				suite.EqualError(err, tc.error(tokenContract))
				return
			}
			suite.Require().NoError(err)
			baseDenom, err := suite.App.Erc20Keeper.GetBaseDenom(suite.Ctx, erc20types.NewBridgeDenom(suite.chainName, tokenContract))
			suite.Require().NoError(err)
			has, err := suite.App.Erc20Keeper.HasERC20Token(suite.Ctx, baseDenom)
			suite.Require().NoError(err)
			suite.True(has)
		})
	}
}
