package keeper_test

import (
	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestKeeper_BridgeToken() {
	tokenContract := helpers.GenHexAddress().Hex()
	denom, err := suite.Keeper().SetIbcDenomTrace(suite.Ctx, tokenContract, "")
	suite.NoError(err)
	suite.Equal(types.NewBridgeDenom(suite.chainName, tokenContract), denom)

	suite.Keeper().AddBridgeToken(suite.Ctx, tokenContract, denom)

	bridgeToken := &types.BridgeToken{Token: tokenContract, Denom: denom}
	suite.EqualValues(bridgeToken, suite.Keeper().GetBridgeTokenDenom(suite.Ctx, tokenContract))

	suite.EqualValues(bridgeToken, suite.Keeper().GetDenomBridgeToken(suite.Ctx, denom))

	suite.Keeper().IterateBridgeTokenToDenom(suite.Ctx, func(bt *types.BridgeToken) bool {
		suite.Equal(bt.Token, tokenContract)
		suite.Equal(bt.Denom, denom)
		return false
	})
}
