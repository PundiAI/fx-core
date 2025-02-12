package keeper_test

import (
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestKeeper_IterateOutgoingTxBatches() {
	tokenContract := helpers.GenExternalAddr(suite.chainName)
	batchNumber := 3
	for i := 1; i <= batchNumber; i++ {
		err := suite.Keeper().StoreBatch(suite.Ctx, &types.OutgoingTxBatch{
			BatchNonce:    uint64(i),
			TokenContract: tokenContract,
		})
		suite.NoError(err)
	}
	batchList := make([]*types.OutgoingTxBatch, 0)
	suite.Keeper().IterateOutgoingTxBatches(suite.Ctx, func(batch *types.OutgoingTxBatch) bool {
		batchList = append(batchList, batch)
		return false
	})
	for i := 0; i < batchNumber; i++ {
		suite.Equal(batchList[i].BatchNonce, uint64(i+1))
	}
}
