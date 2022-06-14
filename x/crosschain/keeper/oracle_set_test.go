package keeper_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/functionx/fx-core/x/crosschain/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestLastPendingOracleSetRequestByAddr() {
	testCases := []struct {
		OracleAddress  sdk.AccAddress
		BridgerAddress sdk.AccAddress
		StartHeight    int64

		ExpectOracleSetSize int
	}{
		{
			OracleAddress:       suite.oracles[0],
			BridgerAddress:      suite.bridgers[0],
			StartHeight:         1,
			ExpectOracleSetSize: 3,
		},
		{
			OracleAddress:       suite.oracles[1],
			BridgerAddress:      suite.bridgers[1],
			StartHeight:         2,
			ExpectOracleSetSize: 2,
		},
		{
			OracleAddress:       suite.oracles[2],
			BridgerAddress:      suite.bridgers[2],
			StartHeight:         3,
			ExpectOracleSetSize: 1,
		},
	}

	for i := 1; i <= 3; i++ {
		suite.Keeper().StoreOracleSet(suite.ctx, &types.OracleSet{
			Nonce: uint64(i),
			Members: types.BridgeValidators{{
				Power:           uint64(i),
				ExternalAddress: fmt.Sprintf("0x%d", i),
			}},
			Height: uint64(i),
		})
	}

	wrapSDKContext := sdk.WrapSDKContext(suite.ctx)
	for _, testCase := range testCases {
		oracle := types.Oracle{
			OracleAddress:  testCase.OracleAddress.String(),
			BridgerAddress: testCase.BridgerAddress.String(),
			StartHeight:    testCase.StartHeight,
		}
		// save oracle
		suite.Keeper().SetOracle(suite.ctx, oracle)
		suite.Keeper().SetOracleByBridger(suite.ctx, oracle.GetOracle(), testCase.BridgerAddress)

		pendingOracleSetRequestByAddr, err := suite.Keeper().LastPendingOracleSetRequestByAddr(wrapSDKContext, &types.QueryLastPendingOracleSetRequestByAddrRequest{
			BridgerAddress: testCase.BridgerAddress.String(),
		})
		require.NoError(suite.T(), err)
		require.EqualValues(suite.T(), testCase.ExpectOracleSetSize, len(pendingOracleSetRequestByAddr.OracleSets))
	}
}

func (suite *KeeperTestSuite) TestGetUnSlashedOracleSets() {
	for i := 1; i <= 3; i++ {
		suite.Keeper().StoreOracleSet(suite.ctx, &types.OracleSet{
			Nonce: uint64(i),
			Members: types.BridgeValidators{{
				Power:           uint64(i),
				ExternalAddress: fmt.Sprintf("0x%d", i),
			}},
			Height: uint64(1000 + i),
		})
	}
	slashOracleSetHeight := 1003
	sets := suite.Keeper().GetUnSlashedOracleSets(suite.ctx, uint64(slashOracleSetHeight))
	require.NotNil(suite.T(), sets)
	require.EqualValues(suite.T(), 2, sets.Len())

	suite.Keeper().SetLastSlashedOracleSetNonce(suite.ctx, 1)
	slashOracleSetHeight = 1003
	sets = suite.Keeper().GetUnSlashedOracleSets(suite.ctx, uint64(slashOracleSetHeight))
	require.NotNil(suite.T(), sets)
	require.EqualValues(suite.T(), 1, sets.Len())

	slashOracleSetHeight = 1004
	sets = suite.Keeper().GetUnSlashedOracleSets(suite.ctx, uint64(slashOracleSetHeight))
	require.NotNil(suite.T(), sets)
	require.EqualValues(suite.T(), 2, sets.Len())

}
