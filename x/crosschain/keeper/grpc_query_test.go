package keeper_test

import (
	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (s *KeeperMockSuite) TestQueryServer_BridgeCalls() {
	data1 := types.OutgoingBridgeCall{
		Nonce:  tmrand.Uint64(),
		Sender: helpers.GenAccAddress().String(),
	}
	data2 := types.OutgoingBridgeCall{
		Nonce:  tmrand.Uint64(),
		Sender: helpers.GenAccAddress().String(),
	}

	s.crosschainKeeper.SetOutgoingBridgeCall(s.ctx, &data1)
	s.crosschainKeeper.SetOutgoingBridgeCall(s.ctx, &data2)
	actual, err := s.queryClient.BridgeCalls(s.ctx, &types.QueryBridgeCallsRequest{
		ChainName: s.chainName,
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      1,
			CountTotal: false,
		},
	})
	s.NoError(err)
	s.Equal(len(actual.BridgeCalls), 1)

	actual, err = s.queryClient.BridgeCalls(s.ctx, &types.QueryBridgeCallsRequest{
		ChainName: s.chainName,
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      2,
			CountTotal: false,
		},
	})
	s.NoError(err)
	s.Equal(len(actual.BridgeCalls), 2)
}

func (s *KeeperMockSuite) TestQueryServer_BridgeCallsByFeeReceiver() {
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
	s.crosschainKeeper.SetOutgoingBridgeCall(s.ctx, &data1)
	s.crosschainKeeper.SetOutgoingBridgeCall(s.ctx, &data2)
	s.crosschainKeeper.SetOutgoingBridgeCallQuoteInfo(s.ctx, data1.Nonce, quote)
	s.crosschainKeeper.SetOutgoingBridgeCallQuoteInfo(s.ctx, data2.Nonce, quote)

	actual, err := s.queryClient.BridgeCallsByFeeReceiver(s.ctx, &types.QueryBridgeCallsByFeeReceiverRequest{
		ChainName:   s.chainName,
		FeeReceiver: quote.Oracle,
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      1,
			CountTotal: false,
		},
	})
	s.NoError(err)
	s.Equal(len(actual.BridgeCalls), 1)

	actual, err = s.queryClient.BridgeCallsByFeeReceiver(s.ctx, &types.QueryBridgeCallsByFeeReceiverRequest{
		ChainName:   s.chainName,
		FeeReceiver: quote.Oracle,
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      2,
			CountTotal: false,
		},
	})
	s.NoError(err)
	s.Equal(len(actual.BridgeCalls), 2)
}
