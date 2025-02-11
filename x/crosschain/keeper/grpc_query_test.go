package keeper_test

import (
	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestQueryServer_BridgeCalls() {
	data1 := types.OutgoingBridgeCall{
		Nonce:  tmrand.Uint64(),
		Sender: helpers.GenAccAddress().String(),
	}
	data2 := types.OutgoingBridgeCall{
		Nonce:  tmrand.Uint64(),
		Sender: helpers.GenAccAddress().String(),
	}

	suite.Keeper().SetOutgoingBridgeCall(suite.Ctx, &data1)
	suite.Keeper().SetOutgoingBridgeCall(suite.Ctx, &data2)
	actual, err := suite.QueryClient().BridgeCalls(suite.Ctx, &types.QueryBridgeCallsRequest{
		ChainName: suite.chainName,
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      1,
			CountTotal: false,
		},
	})
	suite.NoError(err)
	suite.Equal(len(actual.BridgeCalls), 1)

	actual, err = suite.QueryClient().BridgeCalls(suite.Ctx, &types.QueryBridgeCallsRequest{
		ChainName: suite.chainName,
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      2,
			CountTotal: false,
		},
	})
	suite.NoError(err)
	suite.Equal(len(actual.BridgeCalls), 2)
}

func (suite *KeeperTestSuite) TestQueryServer_BridgeCallsByFeeReceiver() {
	data1 := types.OutgoingBridgeCall{
		Nonce:  tmrand.Uint64(),
		Sender: helpers.GenAccAddress().String(),
	}
	data2 := types.OutgoingBridgeCall{
		Nonce:  tmrand.Uint64(),
		Sender: helpers.GenAccAddress().String(),
	}

	quote := types.QuoteInfo{
		Oracle: helpers.GenHexAddress().Hex(),
	}
	suite.Keeper().SetOutgoingBridgeCall(suite.Ctx, &data1)
	suite.Keeper().SetOutgoingBridgeCall(suite.Ctx, &data2)
	suite.Keeper().SetOutgoingBridgeCallQuoteInfo(suite.Ctx, data1.Nonce, quote)
	suite.Keeper().SetOutgoingBridgeCallQuoteInfo(suite.Ctx, data2.Nonce, quote)

	actual, err := suite.QueryClient().BridgeCallsByFeeReceiver(suite.Ctx, &types.QueryBridgeCallsByFeeReceiverRequest{
		ChainName:   suite.chainName,
		FeeReceiver: quote.Oracle,
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      1,
			CountTotal: false,
		},
	})
	suite.NoError(err)
	suite.Equal(len(actual.BridgeCalls), 1)

	actual, err = suite.QueryClient().BridgeCallsByFeeReceiver(suite.Ctx, &types.QueryBridgeCallsByFeeReceiverRequest{
		ChainName:   suite.chainName,
		FeeReceiver: quote.Oracle,
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      2,
			CountTotal: false,
		},
	})
	suite.NoError(err)
	suite.Equal(len(actual.BridgeCalls), 2)
}

func (suite *KeeperTestSuite) TestQueryServer_OutgoingTxBatches() {
	batchNumber := 3
	token := helpers.GenExternalAddr(suite.chainName)
	testCases := []struct {
		name     string
		req      *types.QueryOutgoingTxBatchesRequest
		expCount int
		expError bool
	}{
		{
			name:     "normal test",
			req:      &types.QueryOutgoingTxBatchesRequest{},
			expCount: batchNumber,
			expError: false,
		},
		{
			name: "chain name test",
			req: &types.QueryOutgoingTxBatchesRequest{
				ChainName: suite.chainName,
			},
			expCount: batchNumber,
			expError: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			for i := 1; i <= batchNumber; i++ {
				batch := &types.OutgoingTxBatch{
					BatchNonce:    uint64(i),
					TokenContract: token,
				}
				err := suite.Keeper().StoreBatch(suite.Ctx, batch)
				suite.NoError(err)
			}

			res, err := suite.QueryClient().OutgoingTxBatches(suite.Ctx, tc.req)
			if tc.expError {
				suite.Error(err)
				return
			}

			suite.NoError(err)
			suite.Equal(tc.expCount, len(res.Batches))

			if tc.expCount > 0 {
				for i, batch := range res.Batches {
					if i > 0 {
						suite.Equal(res.Batches[i-1].BatchNonce+1, batch.BatchNonce)
					}
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryServer_LastPendingBatchRequestByAddr() {
	oracle := suite.oracleAddrs[0]
	bridger := suite.bridgerAddrs[0]
	token := helpers.GenExternalAddr(suite.chainName)
	batchNumber := 3

	testCases := []struct {
		name      string
		setup     func()
		req       *types.QueryLastPendingBatchRequestByAddrRequest
		expLength int
		expError  bool
	}{
		{
			name: "normal test",
			setup: func() {
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, bridger, oracle)
				suite.Keeper().SetOracle(suite.Ctx, types.Oracle{
					OracleAddress:  oracle.String(),
					BridgerAddress: bridger.String(),
				})
				for i := 1; i <= batchNumber; i++ {
					batch := &types.OutgoingTxBatch{
						BatchNonce:    uint64(i),
						TokenContract: token,
					}
					err := suite.Keeper().StoreBatch(suite.Ctx, batch)
					suite.NoError(err)
				}
			},
			req: &types.QueryLastPendingBatchRequestByAddrRequest{
				ChainName:      suite.chainName,
				BridgerAddress: bridger.String(),
			},
			expLength: batchNumber,
			expError:  false,
		},
		{
			name: "null test",
			setup: func() {
				suite.Keeper().SetOracleAddrByBridgerAddr(suite.Ctx, bridger, oracle)
				suite.Keeper().SetOracle(suite.Ctx, types.Oracle{
					OracleAddress:  oracle.String(),
					BridgerAddress: bridger.String(),
				})
			},
			req: &types.QueryLastPendingBatchRequestByAddrRequest{
				ChainName:      suite.chainName,
				BridgerAddress: bridger.String(),
			},
			expLength: 0,
			expError:  false,
		},
		{
			name:  "error bridger",
			setup: func() {},
			req: &types.QueryLastPendingBatchRequestByAddrRequest{
				ChainName:      suite.chainName,
				BridgerAddress: bridger.String(),
			},
			expLength: 0,
			expError:  true,
		},
		{
			name: "error chain name",
			setup: func() {
				batch := &types.OutgoingTxBatch{
					BatchNonce:    1,
					TokenContract: token,
					Transactions: []*types.OutgoingTransferTx{
						{
							Id:     1,
							Sender: helpers.GenAccAddress().String(),
						},
					},
				}
				err := suite.Keeper().StoreBatch(suite.Ctx, batch)
				suite.NoError(err)
			},
			req: &types.QueryLastPendingBatchRequestByAddrRequest{
				ChainName:      "wrong-chain",
				BridgerAddress: bridger.String(),
			},
			expLength: 0,
			expError:  true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // 重置测试环境

			if tc.setup != nil {
				tc.setup()
			}

			res, err := suite.QueryClient().LastPendingBatchRequestByAddr(suite.Ctx, tc.req)
			if tc.expError {
				suite.Error(err)
				return
			}

			suite.NoError(err)
			if tc.expLength > 0 {
				suite.Equal(tc.expLength, len(res.GetBatchs()))
				for i := 0; i < batchNumber; i++ {
					res.Batchs[i].BatchNonce = uint64(i + 1)
				}
			} else {
				suite.Nil(res.GetBatchs())
			}
		})
	}
}
