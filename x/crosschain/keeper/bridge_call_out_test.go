package keeper_test

import (
	tmrand "github.com/cometbft/cometbft/libs/rand"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestKeeper_BridgeCallResultHandler() {
	tests := []struct {
		name     string
		initData func(msg *types.MsgBridgeCallResultClaim)
	}{
		{
			name: "success",
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			msg := &types.MsgBridgeCallResultClaim{
				ChainName:      suite.chainName,
				BridgerAddress: helpers.GenAccAddress().String(),
				EventNonce:     1,
				BlockHeight:    1,
				Nonce:          1,
				TxOrigin:       helpers.GenExternalAddr(suite.chainName),
				Success:        true,
				Cause:          "",
			}
			if tt.initData != nil {
				tt.initData(msg)
			}
			suite.NoError(msg.ValidateBasic())

			suite.Keeper().SetOutgoingBridgeCall(suite.Ctx, &types.OutgoingBridgeCall{
				Sender:      helpers.GenExternalAddr(suite.chainName),
				Refund:      "",
				Tokens:      nil,
				To:          "",
				Data:        "",
				Memo:        "",
				Nonce:       msg.Nonce,
				Timeout:     0,
				BlockHeight: 0,
			})
			err := suite.Keeper().BridgeCallResultExecuted(suite.Ctx, suite.App.EvmKeeper, msg)
			suite.Require().NoError(err)
			outgoingBridgeCall, found := suite.Keeper().GetOutgoingBridgeCallByNonce(suite.Ctx, msg.Nonce)
			suite.False(found)
			suite.Nil(outgoingBridgeCall)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_DeleteOutgoingBridgeCall() {
	outCall := &types.OutgoingBridgeCall{
		Sender: helpers.GenHexAddress().String(),
		Nonce:  tmrand.Uint64(),
	}
	outCallNonce := suite.Keeper().AddOutgoingBridgeCallWithoutBuild(suite.Ctx, outCall)
	suite.Require().EqualValues(outCall.Nonce, outCallNonce)

	suite.Require().True(suite.Keeper().HasOutgoingBridgeCall(suite.Ctx, outCall.Nonce))
	suite.Require().True(suite.Keeper().HasOutgoingBridgeCallAddressAndNonce(suite.Ctx, outCall.Sender, outCall.Nonce))

	suite.Keeper().DeleteOutgoingBridgeCall(suite.Ctx, outCall.Nonce)

	suite.Require().False(suite.Keeper().HasOutgoingBridgeCall(suite.Ctx, outCall.Nonce))
	suite.Require().False(suite.Keeper().HasOutgoingBridgeCallAddressAndNonce(suite.Ctx, outCall.Sender, outCall.Nonce))
}
