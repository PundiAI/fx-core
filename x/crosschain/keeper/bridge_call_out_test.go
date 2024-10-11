package keeper_test

import (
	tmrand "github.com/cometbft/cometbft/libs/rand"
	"go.uber.org/mock/gomock"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (s *KeeperMockSuite) TestKeeper_BridgeCallResultHandler() {
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
				Refund:      "",
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

func (s *KeeperMockSuite) TestKeeper_DeleteOutgoingBridgeCall() {
	outCall := &types.OutgoingBridgeCall{
		Sender: helpers.GenHexAddress().String(),
		Nonce:  tmrand.Uint64(),
	}
	outCallNonce := s.crosschainKeeper.AddOutgoingBridgeCallWithoutBuild(s.ctx, outCall)
	s.Require().EqualValues(outCall.Nonce, outCallNonce)

	s.Require().True(s.crosschainKeeper.HasOutgoingBridgeCall(s.ctx, outCall.Nonce))
	s.Require().True(s.crosschainKeeper.HasOutgoingBridgeCallAddressAndNonce(s.ctx, outCall.Sender, outCall.Nonce))

	s.crosschainKeeper.DeleteOutgoingBridgeCall(s.ctx, outCall.Nonce)

	s.Require().False(s.crosschainKeeper.HasOutgoingBridgeCall(s.ctx, outCall.Nonce))
	s.Require().False(s.crosschainKeeper.HasOutgoingBridgeCallAddressAndNonce(s.ctx, outCall.Sender, outCall.Nonce))
}
