package keeper_test

import (
	"go.uber.org/mock/gomock"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (s *KeeperTestSuite) TestKeeper_BridgeCallResultHandler() {
	tests := []struct {
		name     string
		initData func(msg *types.MsgBridgeCallResultClaim)
	}{
		{
			name: "success",
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			msg := &types.MsgBridgeCallResultClaim{
				ChainName:      s.moduleName,
				BridgerAddress: helpers.GenAccAddress().String(),
				EventNonce:     1,
				BlockHeight:    1,
				Sender:         helpers.GenExternalAddr(s.moduleName),
				Receiver:       helpers.GenExternalAddr(s.moduleName),
				To:             helpers.GenExternalAddr(s.moduleName),
				Nonce:          1,
				TxOrigin:       helpers.GenExternalAddr(s.moduleName),
				Success:        true,
				Cause:          "",
			}
			if tt.initData != nil {
				tt.initData(msg)
			}
			s.NoError(msg.ValidateBasic())

			s.accountKeeper.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			s.accountKeeper.EXPECT().NewAccountWithAddress(gomock.Any(), gomock.Any()).Times(1)

			s.crosschainKeeper.SetOutgoingBridgeCall(s.ctx, &types.OutgoingBridgeCall{
				Sender:      helpers.GenExternalAddr(s.moduleName),
				Receiver:    "",
				Tokens:      nil,
				To:          "",
				Data:        "",
				Memo:        "",
				Nonce:       msg.Nonce,
				Timeout:     0,
				BlockHeight: 0,
			})
			s.crosschainKeeper.BridgeCallResultHandler(s.ctx, msg)
			outgoingBridgeCall, found := s.crosschainKeeper.GetOutgoingBridgeCallByNonce(s.ctx, msg.Nonce)
			s.False(found)
			s.Nil(outgoingBridgeCall)
		})
	}
}
