package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (s *KeeperTestSuite) TestKeeper_BridgeCalls() {
	ctx := sdk.WrapSDKContext(s.ctx)
	data1 := types.OutgoingBridgeCall{
		Nonce:  tmrand.Uint64(),
		Sender: sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
	}
	data2 := types.OutgoingBridgeCall{
		Nonce:  tmrand.Uint64(),
		Sender: sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
	}

	s.crosschainKeeper.SetOutgoingBridgeCall(s.ctx, &data1)
	s.crosschainKeeper.SetOutgoingBridgeCall(s.ctx, &data2)
	actual, err := s.queryClient.BridgeCalls(ctx, &types.QueryBridgeCallsRequest{
		ChainName: s.moduleName,
		Pagination: &query.PageRequest{
			Offset:     0,
			Limit:      1,
			CountTotal: false,
		},
	})
	s.NoError(err)
	s.Equal(len(actual.BridgeCalls), 1)

	actual, err = s.queryClient.BridgeCalls(ctx, &types.QueryBridgeCallsRequest{
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
