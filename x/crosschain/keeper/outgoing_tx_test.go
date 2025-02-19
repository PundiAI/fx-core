package keeper_test

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestKeeper_IterateOutgoingTxBatches() {
	batchNumber := 3
	suite.addOutgoingTxBatchs(batchNumber, helpers.GenExternalAddr(suite.chainName))
	batchList := make([]*types.OutgoingTxBatch, 0)
	suite.Keeper().IterateOutgoingTxBatches(suite.Ctx, func(batch *types.OutgoingTxBatch) bool {
		batchList = append(batchList, batch)
		return false
	})
	for i := 0; i < batchNumber; i++ {
		suite.Equal(batchList[i].BatchNonce, uint64(i+1))
	}
}

func (suite *KeeperTestSuite) TestResendTimeoutOutgoingTxBatch() {
	tokenContract := helpers.GenExternalAddr(suite.chainName)
	suite.addOutgoingTxBatchs(5, tokenContract)
	storeKey := suite.App.GetKVStoreKey()[suite.Keeper().ModuleName()]
	// set last batch id to 6
	suite.Ctx.KVStore(storeKey).Set(types.KeyLastOutgoingBatchID, sdk.Uint64ToBigEndian(6))

	bz := suite.Ctx.KVStore(storeKey).Get(types.KeyLastOutgoingBatchID)
	suite.EqualValues(6, sdk.BigEndianToUint64(bz))

	suite.Commit(3)
	// if batchNonce 5 is executed, 1,2,3,4 will be resend
	err := suite.Keeper().OutgoingTxBatchExecuted(suite.Ctx, tokenContract, 5)
	suite.NoError(err)

	actualBatchList := make([]uint64, 0, 4)
	suite.Keeper().IterateOutgoingTxBatches(suite.Ctx, func(batch *types.OutgoingTxBatch) bool {
		suite.EqualValues(uint64(suite.Ctx.BlockHeight()), batch.Block)
		actualBatchList = append(actualBatchList, batch.BatchNonce)
		return false
	})
	// expected batch 6,7,8,9
	suite.EqualValues([]uint64{6, 7, 8, 9}, actualBatchList)
}

func (suite *KeeperTestSuite) addOutgoingTxBatchs(batchNumber int, tokenContract string) {
	for i := 1; i <= batchNumber; i++ {
		err := suite.Keeper().StoreBatch(suite.Ctx, &types.OutgoingTxBatch{
			// save batch with same block height

			Block:         1,
			BatchNonce:    uint64(i),
			TokenContract: tokenContract,
			Transactions: []*types.OutgoingTransferTx{
				{
					Id:          1,
					Sender:      helpers.GenAccAddress().String(),
					DestAddress: helpers.GenExternalAddr(suite.chainName),
					Token: types.ERC20Token{
						Contract: tokenContract,
						Amount:   sdkmath.NewInt(tmrand.Int63()),
					},
					Fee: types.ERC20Token{
						Contract: tokenContract,
						Amount:   sdkmath.NewInt(tmrand.Int63()),
					},
				},
			},
		})
		suite.NoError(err)
	}
}

func (suite *KeeperTestSuite) TestKeeper_OutgoingTxBatchExecuted() {
	batchNumber := 1
	tokenContract := helpers.GenExternalAddr(suite.chainName)
	suite.addOutgoingTxBatchs(batchNumber, tokenContract)

	err := suite.Keeper().OutgoingTxBatchExecuted(suite.Ctx, tokenContract, 1)
	suite.Require().NoError(err)
	batch := suite.Keeper().GetOutgoingTxBatch(suite.Ctx, tokenContract, 1)
	suite.Empty(batch)

	err = suite.Keeper().OutgoingTxBatchExecuted(suite.Ctx, tokenContract, 1)
	suite.Error(err)
	suite.EqualError(err, fmt.Sprintf("unknown batch nonce for outgoing tx batch %s %d", tokenContract, 1))
}
