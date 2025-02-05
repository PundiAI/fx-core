package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestKeeper_SavePendingExecuteClaim() {
	tests := []struct {
		name  string
		claim types.ExternalClaim
	}{
		{
			name: "msg bridge call claim",
			claim: &types.MsgBridgeCallClaim{
				ChainName:      suite.chainName,
				BridgerAddress: helpers.GenAccAddress().String(),
				EventNonce:     1,
				BlockHeight:    100,
				Sender:         helpers.GenExternalAddr(suite.chainName),
				Refund:         helpers.GenExternalAddr(suite.chainName),
				TokenContracts: []string{helpers.GenExternalAddr(suite.chainName)},
				Amounts:        []sdkmath.Int{sdkmath.NewInt(1)},
				To:             helpers.GenExternalAddr(suite.chainName),
				Data:           "",
				QuoteId:        sdkmath.NewInt(0),
				Memo:           "",
				TxOrigin:       "",
			},
		},
		{
			name: "msg send to fx claim",
			claim: &types.MsgSendToFxClaim{
				EventNonce:     1,
				BlockHeight:    100,
				TokenContract:  helpers.GenExternalAddr(suite.chainName),
				Amount:         sdkmath.NewInt(1),
				Sender:         helpers.GenExternalAddr(suite.chainName),
				Receiver:       helpers.GenExternalAddr(suite.chainName),
				TargetIbc:      "",
				BridgerAddress: helpers.GenAccAddress().String(),
				ChainName:      suite.chainName,
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := suite.App.GetKeeper(suite.chainName).SavePendingExecuteClaim(suite.Ctx, tt.claim)
			suite.Require().NoError(err)

			claim, err := suite.App.GetKeeper(suite.chainName).GetPendingExecuteClaim(suite.Ctx, tt.claim.GetEventNonce())
			suite.Require().NoError(err)
			suite.Require().Equal(claim, tt.claim)
		})
	}
}
