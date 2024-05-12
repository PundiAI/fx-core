package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (s *KeeperTestSuite) TestKeeper_HandleOutgoingBridgeCallRefund() {
	tests := []struct {
		name     string
		initData func(outgoingBridgeCall *types.OutgoingBridgeCall)
	}{
		{
			name: "success",
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			outgoingBridgeCall := &types.OutgoingBridgeCall{
				Sender:   helpers.GenerateAddressByModule(s.moduleName),
				Receiver: helpers.GenerateAddressByModule(s.moduleName),
				Tokens: []types.ERC20Token{
					{
						Contract: helpers.GenerateAddressByModule(s.moduleName),
						Amount:   sdkmath.NewInt(1),
					},
					{
						Contract: helpers.GenerateAddressByModule(s.moduleName),
						Amount:   sdkmath.NewInt(2),
					},
				},
				To:          helpers.GenerateAddressByModule(s.moduleName),
				Data:        "",
				Memo:        "",
				Nonce:       1,
				Timeout:     0,
				BlockHeight: 1,
			}
			if tt.initData != nil {
				tt.initData(outgoingBridgeCall)
			}
			s.MockBridgeCallToken(outgoingBridgeCall.Tokens)

			s.crosschainKeeper.HandleOutgoingBridgeCallRefund(s.ctx, outgoingBridgeCall)
		})
	}
}
