package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

	fxtypes "github.com/functionx/fx-core/v8/types"
	ibctesting "github.com/functionx/fx-core/v8/x/ibc/testing"
)

func (suite *KeeperTestSuite) TestMsgTransfer() {
	var msg *types.MsgTransfer

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {},
			true,
		},
		{
			"send transfers disabled",
			func() {
				suite.GetApp(suite.chainA.App).IBCTransferKeeper.SetParams(suite.chainA.GetContext(),
					types.Params{
						SendEnabled: false,
					},
				)
			},
			false,
		},
		{
			"invalid sender",
			func() {
				msg.Sender = "address"
			},
			false,
		},
		{
			"sender is a blocked address",
			func() {
				msg.Sender = suite.GetApp(suite.chainA.App).AccountKeeper.GetModuleAddress(types.ModuleName).String()
			},
			false,
		},
		{
			"channel does not exist",
			func() {
				msg.SourceChannel = "channel-100"
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			path := ibctesting.NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			coin := sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(100))
			msg = types.NewMsgTransfer(
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
				coin, suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String(),
				suite.chainB.GetTimeoutHeight(), 0, // only use timeout height
				"memo",
			)

			tc.malleate()

			res, err := suite.GetApp(suite.chainA.App).IBCTransferKeeper.Transfer(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().NotEqual(res.Sequence, uint64(0))
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}
