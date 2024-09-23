package keeper_test

import (
	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
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
		ChainName: s.moduleName,
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      1,
			CountTotal: false,
		},
	})
	s.NoError(err)
	s.Equal(len(actual.BridgeCalls), 1)

	actual, err = s.queryClient.BridgeCalls(s.ctx, &types.QueryBridgeCallsRequest{
		ChainName: s.moduleName,
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      2,
			CountTotal: false,
		},
	})
	s.NoError(err)
	s.Equal(len(actual.BridgeCalls), 2)
}

func (s *KeeperMockSuite) TestQueryServer_PendingBridgeCalls() {
	data1 := types.PendingOutgoingBridgeCall{
		OutgoinBridgeCall: &types.OutgoingBridgeCall{
			Nonce:  tmrand.Uint64(),
			Sender: helpers.GenExternalAddr(s.moduleName),
		},
	}
	data2 := types.PendingOutgoingBridgeCall{
		OutgoinBridgeCall: &types.OutgoingBridgeCall{
			Nonce:  tmrand.Uint64(),
			Sender: helpers.GenExternalAddr(s.moduleName),
		},
	}

	s.crosschainKeeper.SetPendingOutgoingBridgeCall(s.ctx, &data1)
	s.crosschainKeeper.SetPendingOutgoingBridgeCall(s.ctx, &data2)
	actual, err := s.queryClient.PendingBridgeCalls(s.ctx, &types.QueryPendingBridgeCallsRequest{
		ChainName: s.moduleName,
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      1,
			CountTotal: false,
		},
	})
	s.NoError(err)
	s.Equal(len(actual.PendingBridgeCalls), 1)

	actual, err = s.queryClient.PendingBridgeCalls(s.ctx, &types.QueryPendingBridgeCallsRequest{
		ChainName: s.moduleName,
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      2,
			CountTotal: false,
		},
	})
	s.NoError(err)
	s.Equal(len(actual.PendingBridgeCalls), 2)

	// test sender is not empty
	actual, err = s.queryClient.PendingBridgeCalls(s.ctx, &types.QueryPendingBridgeCallsRequest{
		ChainName:     s.moduleName,
		SenderAddress: data1.OutgoinBridgeCall.Sender,
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      2,
			CountTotal: false,
		},
	})
	s.NoError(err)
	s.Equal(len(actual.PendingBridgeCalls), 1)

	// test sender addr is invalid
	_, err = s.queryClient.PendingBridgeCalls(s.ctx, &types.QueryPendingBridgeCallsRequest{
		ChainName:     s.moduleName,
		SenderAddress: "invalid",
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      2,
			CountTotal: false,
		},
	})
	s.Error(err)
	s.EqualError(err, "rpc error: code = InvalidArgument desc = sender address")
}
