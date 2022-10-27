package keeper_test

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v2/app/helpers"
	fxtypes "github.com/functionx/fx-core/v2/types"
	"github.com/functionx/fx-core/v2/x/erc20/types"
)

var (
	fxTokenPair = types.TokenPair{
		Erc20Address:  "0x80b5a32E4F032B2a058b4F29EC95EEfEEB87aDcd",
		Denom:         "FX",
		Enabled:       true,
		ContractOwner: 1,
	}
	pundixTokenPair = types.TokenPair{
		Erc20Address:  "0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75",
		Denom:         "eth0x0FD10b9899882a6f2fcb5c371E17e70FdEe00C38",
		Enabled:       true,
		ContractOwner: 1,
	}
	purseTokenPair = types.TokenPair{
		Erc20Address:  "0x5FD55A1B9FC24967C4dB09C513C3BA0DFa7FF687",
		Denom:         "ibc/F08B62C2C1BE9E52942617489CAB1E94537FE3849F8EEC910B142468C340EB0D",
		Enabled:       true,
		ContractOwner: 1,
	}
)

func (suite *KeeperTestSuite) TestGetAllTokenPairs() {
	var expRes []types.TokenPair

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"5 pair registered", func() {
				expRes = []types.TokenPair{fxTokenPair, pundixTokenPair, purseTokenPair}
			},
		},
		{
			"6 pair registered",
			func() {
				pair := types.NewTokenPair(helpers.GenerateAddress(), "coin", true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)

				expRes = []types.TokenPair{pair, fxTokenPair, pundixTokenPair, purseTokenPair}
			},
		},
		{
			"7 pairs registered",
			func() {
				pair := types.NewTokenPair(helpers.GenerateAddress(), "coin", true, types.OWNER_MODULE)
				pair2 := types.NewTokenPair(helpers.GenerateAddress(), "coin2", true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair2)

				expRes = []types.TokenPair{pair, pair2, fxTokenPair, pundixTokenPair, purseTokenPair}
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()
			res := suite.app.Erc20Keeper.GetAllTokenPairs(suite.ctx)

			suite.Require().ElementsMatch(expRes, res, tc.name)
		})
	}
}

func (suite *KeeperTestSuite) TestGetTokenPairID() {
	pair := types.NewTokenPair(helpers.GenerateAddress(), fxtypes.DefaultDenom, true, types.OWNER_MODULE)
	suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)

	testCases := []struct {
		name  string
		token string
		expId []byte
	}{
		{"nil token", "", nil},
		{"valid hex token", helpers.GenerateAddress().Hex(), []byte{}},
		{"valid hex token", helpers.GenerateAddress().String(), []byte{}},
	}
	for _, tc := range testCases {
		id := suite.app.Erc20Keeper.GetTokenPairID(suite.ctx, tc.token)
		if id != nil {
			suite.Require().Equal(tc.expId, id, tc.name)
		} else {
			suite.Require().Nil(id)
		}
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
	id := pair.GetID()
	suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)
	suite.app.Erc20Keeper.SetERC20Map(suite.ctx, pair.GetERC20Contract(), id)
	suite.app.Erc20Keeper.SetDenomMap(suite.ctx, pair.Denom, id)

	testCases := []struct {
		name     string
		id       []byte
		malleate func()
		ok       bool
	}{
		{"nil id", nil, func() {}, false},
		{"pair not found", []byte{}, func() {}, false},
		{"valid id", id, func() {}, true},
		{
			"detete tokenpair",
			id,
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
	addr := helpers.GenerateAddress()
	pair := types.NewTokenPair(addr, "coin", true, types.OWNER_MODULE)
	suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)
	suite.app.Erc20Keeper.SetERC20Map(suite.ctx, addr, pair.GetID())
	suite.app.Erc20Keeper.SetDenomMap(suite.ctx, pair.Denom, pair.GetID())

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
	addr := helpers.GenerateAddress()
	pair := types.NewTokenPair(addr, "coin", true, types.OWNER_MODULE)
	suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)
	suite.app.Erc20Keeper.SetERC20Map(suite.ctx, addr, pair.GetID())
	suite.app.Erc20Keeper.SetDenomMap(suite.ctx, pair.Denom, pair.GetID())

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
