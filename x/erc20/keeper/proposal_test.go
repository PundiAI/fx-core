package keeper_test

import (
	"crypto/sha256"
	"fmt"

	"cosmossdk.io/collections"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/erc20/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func (suite *KeeperTestSuite) TestKeeper_RegisterNativeCoin() {
	testCases := []struct {
		name      string
		getSymbol func() string
		err       error
	}{
		{
			name: "success - origin token",
			getSymbol: func() string {
				return fxtypes.DefaultSymbol
			},
		},
		{
			name:      "success - native token",
			getSymbol: helpers.NewRandSymbol,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			symbol := tc.getSymbol()
			erc20Token, err := suite.GetKeeper().RegisterNativeCoin(suite.Ctx, "token name", symbol, 18)

			if tc.err != nil {
				suite.Require().Error(err)
				suite.Require().ErrorIs(tc.err, err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().True(suite.GetKeeper().HasERC20Token(suite.Ctx, erc20Token.Denom))
			suite.Require().True(suite.GetKeeper().HasToken(suite.Ctx, erc20Token.Erc20Address))
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_RegisterNativeERC20() {
	testCases := []struct {
		name   string
		symbol string
		err    error
	}{
		{
			name:   "success",
			symbol: helpers.NewRandSymbol(),
		},
		{
			name:   "success - origin token",
			symbol: fxtypes.DefaultSymbol,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tokenAddr := suite.erc20TokenSuite.DeployERC20Token(suite.Ctx, suite.signer.Address(), tc.symbol)

			erc20Token, err := suite.GetKeeper().RegisterNativeERC20(suite.Ctx, tokenAddr)
			if tc.err != nil {
				suite.Require().Error(err)
				suite.Require().ErrorIs(tc.err, err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().True(suite.GetKeeper().HasERC20Token(suite.Ctx, erc20Token.Denom))
			suite.Require().True(suite.GetKeeper().HasToken(suite.Ctx, erc20Token.Erc20Address))
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_RegisterBridgeToken() {
	testCases := []struct {
		name      string
		contract  string
		chainName string
		ibcDenom  string
		channel   string
		expError  string
	}{
		{
			name:      "success - bridge token",
			contract:  helpers.GenExternalAddr(ethtypes.ModuleName),
			chainName: ethtypes.ModuleName,
		},
		{
			name:     "success - ibc token",
			ibcDenom: fmt.Sprintf("%s/%s", ibctransfertypes.DenomPrefix, sha256.Sum256([]byte("test"))),
			channel:  "channel-1",
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			symbol := helpers.NewRandSymbol()
			md := fxtypes.NewMetadata("Test Token", symbol, 18)
			_, err := suite.GetKeeper().AddERC20Token(suite.Ctx, md, helpers.GenHexAddress(), types.OWNER_MODULE)
			suite.Require().NoError(err)

			_, err = suite.GetKeeper().RegisterBridgeToken(suite.Ctx, md.Base, tc.channel, tc.ibcDenom, tc.chainName, tc.contract, true)
			if tc.expError != "" {
				suite.Require().Error(err)
				suite.Require().ErrorContains(err, tc.expError)
				return
			}
			suite.Require().NoError(err)
			if tc.ibcDenom != "" {
				has1, _ := suite.GetKeeper().HasToken(suite.Ctx, tc.ibcDenom)
				suite.True(has1)
				has2, _ := suite.GetKeeper().IBCToken.Has(suite.Ctx, collections.Join(md.Base, tc.channel))
				suite.True(has2)
			} else {
				has1, _ := suite.GetKeeper().HasToken(suite.Ctx, types.NewBridgeDenom(tc.chainName, tc.contract))
				suite.True(has1)
				has2, _ := suite.GetKeeper().BridgeToken.Has(suite.Ctx, collections.Join(tc.chainName, md.Base))
				suite.True(has2)
			}
		})
	}
}
