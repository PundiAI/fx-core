package keeper_test

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (suite *KeeperTestSuite) TestGetTokenPairID() {
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
		id := suite.app.Erc20Keeper.GetTokenPairID(suite.ctx, tc.token)
		suite.Require().Equal(tc.expId, id, tc.name)
	}
}

func (suite *KeeperTestSuite) TestGetTokenPair() {
	pair := types.NewTokenPair(helpers.GenerateAddress(), fxtypes.DefaultDenom, true, types.OWNER_MODULE)
	suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)

	testCases := []struct {
		name string
		id   []byte
		ok   bool
	}{
		{"nil id", nil, false},
		{"valid id", pair.GetID(), true},
		{"pair not found", []byte{}, false},
	}
	for _, tc := range testCases {
		p, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, tc.id)
		if tc.ok {
			suite.Require().True(found, tc.name)
			suite.Require().Equal(pair, p, tc.name)
		} else {
			suite.Require().False(found, tc.name)
		}
	}
}

func (suite *KeeperTestSuite) TestDeleteTokenPair() {
	pair := types.NewTokenPair(helpers.GenerateAddress(), fxtypes.DefaultDenom, true, types.OWNER_MODULE)
	suite.app.Erc20Keeper.AddTokenPair(suite.ctx, pair)

	testCases := []struct {
		name     string
		id       []byte
		malleate func()
		ok       bool
	}{
		{"nil id", nil, func() {}, false},
		{"pair not found", []byte{}, func() {}, false},
		{"valid id", pair.GetID(), func() {}, true},
		{
			"detete tokenpair",
			pair.GetID(),
			func() {
				suite.app.Erc20Keeper.RemoveTokenPair(suite.ctx, pair)
			},
			false,
		},
	}
	for _, tc := range testCases {
		tc.malleate()

		p, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, tc.id)
		if tc.ok {
			suite.Require().True(found, tc.name)
			suite.Require().Equal(pair, p, tc.name)
		} else {
			suite.Require().False(found, tc.name)
		}
	}
}

func (suite *KeeperTestSuite) TestIsTokenPairRegistered() {
	pair := types.NewTokenPair(helpers.GenerateAddress(), fxtypes.DefaultDenom, true, types.OWNER_MODULE)
	suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)

	testCases := []struct {
		name string
		id   []byte
		ok   bool
	}{
		{"valid id", pair.GetID(), true},
		{"pair not found", []byte{}, false},
	}
	for _, tc := range testCases {
		found := suite.app.Erc20Keeper.IsTokenPairRegistered(suite.ctx, tc.id)
		if tc.ok {
			suite.Require().True(found, tc.name)
		} else {
			suite.Require().False(found, tc.name)
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
			"deleted erc20 map",
			pair.GetERC20Contract(),
			func() {
				suite.app.Erc20Keeper.RemoveTokenPair(suite.ctx, pair)
			},
			false,
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
			"deleted denom map",
			pair.GetDenom(),
			func() {
				suite.app.Erc20Keeper.RemoveTokenPair(suite.ctx, pair)
			},
			false,
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
