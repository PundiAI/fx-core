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
