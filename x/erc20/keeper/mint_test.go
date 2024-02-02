package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/erc20/types"
)

func (suite *KeeperTestSuite) TestMintingEnabled() {
	testCases := []struct {
		name     string
		malleate func() types.TokenPair
		expPass  bool
	}{
		{
			"intrarelaying is disabled globally",
			func() types.TokenPair {
				params := types.DefaultParams()
				params.EnableErc20 = false
				err := suite.app.Erc20Keeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				return types.NewTokenPair(helpers.GenerateAddress(), "coin", true, types.OWNER_MODULE)
			},
			false,
		},
		{
			"token pair not found",
			func() types.TokenPair {
				return types.NewTokenPair(helpers.GenerateAddress(), "coin", true, types.OWNER_MODULE)
			},
			false,
		},
		{
			"intrarelaying is disabled for the given pair",
			func() types.TokenPair {
				expPair := types.NewTokenPair(helpers.GenerateAddress(), "coin", true, types.OWNER_MODULE)
				expPair.Enabled = false
				suite.app.Erc20Keeper.AddTokenPair(suite.ctx, expPair)
				return expPair
			},
			false,
		},
		{
			"token transfers are disabled",
			func() types.TokenPair {
				expPair := types.NewTokenPair(helpers.GenerateAddress(), "coin", true, types.OWNER_MODULE)
				expPair.Enabled = true
				suite.app.Erc20Keeper.AddTokenPair(suite.ctx, expPair)

				params := banktypes.DefaultParams()
				params.SendEnabled = []*banktypes.SendEnabled{
					{Denom: expPair.Denom, Enabled: false},
				}
				suite.app.BankKeeper.SetParams(suite.ctx, params)
				return expPair
			},
			false,
		},
		{
			"ok",
			func() types.TokenPair {
				expPair := types.NewTokenPair(helpers.GenerateAddress(), "coin", true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.AddTokenPair(suite.ctx, expPair)
				return expPair
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			expPair := tc.malleate()

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
