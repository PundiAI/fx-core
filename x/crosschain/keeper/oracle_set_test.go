package keeper_test

import (
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v3/app/helpers"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
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
		suite.Keeper().SetOracleByBridger(suite.ctx, testCase.BridgerAddress, oracle.GetOracle())

		pendingOracleSetRequestByAddr, err := suite.Keeper().LastPendingOracleSetRequestByAddr(wrapSDKContext,
			&types.QueryLastPendingOracleSetRequestByAddrRequest{
				BridgerAddress: testCase.BridgerAddress.String(),
			})
		require.NoError(suite.T(), err)
		require.EqualValues(suite.T(), testCase.ExpectOracleSetSize, len(pendingOracleSetRequestByAddr.OracleSets))
	}
}

func (suite *KeeperTestSuite) TestGetUnSlashedOracleSets() {
	height := rand.Intn(1000) + 1
	index := rand.Intn(100) + 1
	for i := 1; i <= index; i++ {
		suite.Keeper().StoreOracleSet(suite.ctx, &types.OracleSet{
			Nonce: uint64(i),
			Members: types.BridgeValidators{{
				Power:           rand.Uint64(),
				ExternalAddress: helpers.GenerateAddress().Hex(),
			}},
			Height: uint64(height + i),
		})
	}

	sets := suite.Keeper().GetUnSlashedOracleSets(suite.ctx, uint64(height+index))
	require.NotNil(suite.T(), sets)
	require.EqualValues(suite.T(), index-1, sets.Len())

	suite.Keeper().SetLastSlashedOracleSetNonce(suite.ctx, 1)
	sets = suite.Keeper().GetUnSlashedOracleSets(suite.ctx, uint64(height+index))
	require.NotNil(suite.T(), sets)
	require.EqualValues(suite.T(), index-2, sets.Len())

	sets = suite.Keeper().GetUnSlashedOracleSets(suite.ctx, uint64(height+index+1))
	require.NotNil(suite.T(), sets)
	require.EqualValues(suite.T(), index-1, sets.Len())

}
