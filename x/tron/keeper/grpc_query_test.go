package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
	trontypes "github.com/pundiai/fx-core/v8/x/tron/types"
)

func (suite *KeeperTestSuite) TestKeeper_GetOracleByExternalAddr() {
	var (
		request  *types.QueryOracleByExternalAddrRequest
		response *types.QueryOracleResponse
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "external address is error",
			malleate: func() {
				request = &types.QueryOracleByExternalAddrRequest{
					ChainName:       trontypes.ModuleName,
					ExternalAddress: helpers.GenHexAddress().Hex(),
				}
			},
			expPass: false,
		},
		{
			name: "nonexistent external address",
			malleate: func() {
				request = &types.QueryOracleByExternalAddrRequest{
					ChainName:       trontypes.ModuleName,
					ExternalAddress: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
				}
			},
			expPass: false,
		},
		{
			name: "normal external address and oracle",
			malleate: func() {
				oracle, bridger, externalKey := suite.NewOracleByBridger()
				request = &types.QueryOracleByExternalAddrRequest{
					ChainName:       trontypes.ModuleName,
					ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
				}
				response = &types.QueryOracleResponse{Oracle: &types.Oracle{
					OracleAddress:   oracle.String(),
					BridgerAddress:  bridger.String(),
					ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
					DelegateAmount:  sdkmath.ZeroInt(),
				}}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		suite.Run(testCase.name, func() {
			suite.SetupTest()
			testCase.malleate()
			res, err := suite.queryServer.GetOracleByExternalAddr(suite.Ctx, request)
			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(response.Oracle, res.Oracle)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
