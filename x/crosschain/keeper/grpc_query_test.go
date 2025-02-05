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
