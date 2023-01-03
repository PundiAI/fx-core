package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/functionx/fx-core/v3/app/helpers"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (suite *KeeperTestSuite) TestMintingEnabled() {
	expPair := types.NewTokenPair(helpers.GenerateAddress(), "coin", true, types.OWNER_MODULE)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"intrarelaying is disabled globally",
			func() {
				params := types.DefaultParams()
				params.EnableErc20 = false
				suite.app.Erc20Keeper.SetParams(suite.ctx, params)
			},
			false,
		},
		{
			"token pair not found",
			func() {},
			false,
		},
		{
			"intrarelaying is disabled for the given pair",
			func() {
				expPair.Enabled = false
				suite.app.Erc20Keeper.AddTokenPair(suite.ctx, expPair)
			},
			false,
		},
		{
			"token transfers are disabled",
			func() {
				expPair.Enabled = true
				suite.app.Erc20Keeper.AddTokenPair(suite.ctx, expPair)

				params := banktypes.DefaultParams()
				params.SendEnabled = []*banktypes.SendEnabled{
					{Denom: expPair.Denom, Enabled: false},
				}
				suite.app.BankKeeper.SetParams(suite.ctx, params)
			},
			false,
		},
		{
			"ok",
			func() {
				suite.app.Erc20Keeper.AddTokenPair(suite.ctx, expPair)
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			receiver := sdk.AccAddress(helpers.GenerateAddress().Bytes())
			pair, err := suite.app.Erc20Keeper.MintingEnabled(suite.ctx, receiver, expPair.Erc20Address)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expPair, pair)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
