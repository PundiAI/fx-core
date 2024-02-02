package keeper_test

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/erc20/types"
)

func (suite *KeeperTestSuite) TestGetTokenPair() {
	pair := types.NewTokenPair(helpers.GenerateAddress(), fxtypes.DefaultDenom, true, types.OWNER_MODULE)
	suite.app.Erc20Keeper.AddTokenPair(suite.ctx, pair)

	testCases := []struct {
		name  string
		token string
		expId []byte
	}{
		{"nil token", "", nil},
		{"valid token", pair.Erc20Address, pair.GetID()},
		{"invalid token", helpers.GenerateAddress().String(), nil},
	}
	for _, tc := range testCases {
		tokenPair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, tc.token)
		if len(tc.expId) > 0 {
			suite.True(found)
			suite.Require().Equal(tc.expId, tokenPair.GetID(), tc.name)
		}
	}
}

func (suite *KeeperTestSuite) TestIsERC20Registered() {
	pair := types.NewTokenPair(helpers.GenerateAddress(), "coin", true, types.OWNER_MODULE)
	suite.app.Erc20Keeper.AddTokenPair(suite.ctx, pair)

	testCases := []struct {
		name     string
		erc20    common.Address
		malleate func()
		ok       bool
	}{
		{"nil erc20 address", common.Address{}, func() {}, false},
		{"valid erc20 address", pair.GetERC20Contract(), func() {}, true},
		{
			"deleted erc20 map", pair.GetERC20Contract(), func() {
				suite.app.Erc20Keeper.RemoveTokenPair(suite.ctx, pair)
			}, false,
		},
	}
	for _, tc := range testCases {
		tc.malleate()

		found := suite.app.Erc20Keeper.IsERC20Registered(suite.ctx, tc.erc20)
		if tc.ok {
			suite.Require().True(found, tc.name)
		} else {
			suite.Require().False(found, tc.name)
		}
	}
}

func (suite *KeeperTestSuite) TestIsDenomRegistered() {
	pair := types.NewTokenPair(helpers.GenerateAddress(), "coin", true, types.OWNER_MODULE)
	suite.app.Erc20Keeper.AddTokenPair(suite.ctx, pair)

	testCases := []struct {
		name     string
		denom    string
		malleate func()
		ok       bool
	}{
		{"empty denom", "", func() {}, false},
		{"valid denom", pair.GetDenom(), func() {}, true},
		{
			"deleted denom map", pair.GetDenom(), func() {
				suite.app.Erc20Keeper.RemoveTokenPair(suite.ctx, pair)
			}, false,
		},
	}
	for _, tc := range testCases {
		tc.malleate()

		found := suite.app.Erc20Keeper.IsDenomRegistered(suite.ctx, tc.denom)
		if tc.ok {
			suite.Require().True(found, tc.name)
		} else {
			suite.Require().False(found, tc.name)
		}
	}
}
