package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v3/app/helpers"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestLastPendingBatchRequestByAddr() {

	testCases := []struct {
		Name              string
		OracleAddress     sdk.AccAddress
		BridgerAddress    sdk.AccAddress
		StartHeight       int64
		ExpectStartHeight uint64
	}{
		{
			Name:              "oracle start height with 1, expect oracle set block 3",
			OracleAddress:     suite.oracles[0],
			BridgerAddress:    suite.bridgers[0],
			StartHeight:       1,
			ExpectStartHeight: 3,
		},
		{
			Name:              "oracle start height with 2, expect oracle set block 2",
			OracleAddress:     suite.oracles[1],
			BridgerAddress:    suite.bridgers[1],
			StartHeight:       2,
			ExpectStartHeight: 3,
		},
		{
			Name:              "oracle start height with 3, expect oracle set block 1",
			OracleAddress:     suite.oracles[2],
			BridgerAddress:    suite.bridgers[2],
			StartHeight:       3,
			ExpectStartHeight: 3,
		},
	}
	for i := uint64(1); i <= 3; i++ {
		suite.ctx = suite.ctx.WithBlockHeight(int64(i))
		err := suite.Keeper().StoreBatch(suite.ctx, &types.OutgoingTxBatch{
			Block:      i,
			BatchNonce: i,
			Transactions: types.OutgoingTransferTxs{{
				Id:          i,
				Sender:      sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				DestAddress: helpers.GenerateAddress().Hex(),
			}},
		})
		require.NoError(suite.T(), err)
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

		response, err := suite.Keeper().LastPendingBatchRequestByAddr(wrapSDKContext,
			&types.QueryLastPendingBatchRequestByAddrRequest{
				BridgerAddress: testCase.BridgerAddress.String(),
			})
		require.NoError(suite.T(), err, testCase.Name)
		require.NotNil(suite.T(), response, testCase.Name)
		require.NotNil(suite.T(), response.Batch, testCase.Name)
		require.EqualValues(suite.T(), testCase.ExpectStartHeight, response.Batch.Block, testCase.Name)
	}
}
