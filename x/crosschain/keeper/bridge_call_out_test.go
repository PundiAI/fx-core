package keeper_test

import (
	"encoding/hex"
	"errors"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestKeeper_BridgeCallResultHandler() {
	tests := []struct {
		name     string
		initData func(msg *types.MsgBridgeCallResultClaim, outCall *types.OutgoingBridgeCall)
		err      error
	}{
		{
			name: "success",
			initData: func(msg *types.MsgBridgeCallResultClaim, _ *types.OutgoingBridgeCall) {
				msg.Success = true
			},
		},
		{
			name: "fail",
			initData: func(msg *types.MsgBridgeCallResultClaim, outCall *types.OutgoingBridgeCall) {
				msg.Success = false
				msg.Cause = hex.EncodeToString([]byte("revert"))
				outCall.To = fxtypes.ExternalAddrToStr(suite.chainName, common.Address{}.Bytes())
			},
		},
		{
			name: "fail with OnRevert",
			initData: func(msg *types.MsgBridgeCallResultClaim, outCall *types.OutgoingBridgeCall) {
				tokenAddr := suite.erc20TokenSuite.DeployERC20Token(suite.Ctx, helpers.NewRandSymbol())
				msg.Success = false
				msg.Cause = hex.EncodeToString([]byte("revert"))
				outCall.To = fxtypes.ExternalAddrToStr(suite.chainName, tokenAddr.Bytes())
			},
			err: errors.New("execution reverted: evm transaction execution failed"),
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
			}
			outCall := &types.OutgoingBridgeCall{
				Sender:      helpers.GenExternalAddr(suite.chainName),
				Refund:      helpers.GenExternalAddr(suite.chainName),
				Tokens:      nil,
				To:          "",
				Data:        "",
				Memo:        "",
				Nonce:       msg.Nonce,
				Timeout:     0,
				BlockHeight: 0,
			}
			tt.initData(msg, outCall)
			suite.NoError(msg.ValidateBasic())
			suite.Keeper().SetOutgoingBridgeCall(suite.Ctx, outCall)

			err := suite.Keeper().BridgeCallResultExecuted(suite.Ctx, suite.App.EvmKeeper, msg)
			if tt.err != nil {
				suite.Require().Equal(tt.err.Error(), err.Error())
			} else {
				suite.Require().NoError(err)
				outgoingBridgeCall, found := suite.Keeper().GetOutgoingBridgeCallByNonce(suite.Ctx, msg.Nonce)
				suite.False(found)
				suite.Nil(outgoingBridgeCall)
			}
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
