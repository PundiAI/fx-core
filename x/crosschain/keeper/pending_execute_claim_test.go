package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (s *KeeperMockSuite) TestKeeper_SavePendingExecuteClaim() {
	tests := []struct {
		name  string
		claim types.ExternalClaim
	}{
		{
			name: "msg bridge call claim",
			claim: &types.MsgBridgeCallClaim{
				ChainName:      s.moduleName,
				BridgerAddress: helpers.GenAccAddress().String(),
				EventNonce:     1,
				BlockHeight:    100,
				Sender:         helpers.GenExternalAddr(s.moduleName),
				Refund:         helpers.GenExternalAddr(s.moduleName),
				TokenContracts: []string{helpers.GenExternalAddr(s.moduleName)},
				Amounts:        []sdkmath.Int{sdkmath.NewInt(1)},
				To:             helpers.GenExternalAddr(s.moduleName),
				Data:           "",
				Value:          sdkmath.NewInt(0),
				Memo:           "",
				TxOrigin:       "",
			},
		},
		{
			name: "msg send to fx claim",
			claim: &types.MsgSendToFxClaim{
				EventNonce:     1,
				BlockHeight:    100,
				TokenContract:  helpers.GenExternalAddr(s.moduleName),
				Amount:         sdkmath.NewInt(1),
				Sender:         helpers.GenExternalAddr(s.moduleName),
				Receiver:       helpers.GenExternalAddr(s.moduleName),
				TargetIbc:      "",
				BridgerAddress: helpers.GenAccAddress().String(),
				ChainName:      s.moduleName,
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.crosschainKeeper.SavePendingExecuteClaim(s.ctx, tt.claim)

			claim, found := s.crosschainKeeper.GetPendingExecuteClaim(s.ctx, tt.claim.GetEventNonce())
			s.Require().Equal(true, found)
			s.Require().Equal(claim, tt.claim)
		})
	}
}
