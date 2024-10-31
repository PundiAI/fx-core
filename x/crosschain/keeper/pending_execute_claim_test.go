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
				ChainName:      s.chainName,
				BridgerAddress: helpers.GenAccAddress().String(),
				EventNonce:     1,
				BlockHeight:    100,
				Sender:         helpers.GenExternalAddr(s.chainName),
				Refund:         helpers.GenExternalAddr(s.chainName),
				TokenContracts: []string{helpers.GenExternalAddr(s.chainName)},
				Amounts:        []sdkmath.Int{sdkmath.NewInt(1)},
				To:             helpers.GenExternalAddr(s.chainName),
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
				TokenContract:  helpers.GenExternalAddr(s.chainName),
				Amount:         sdkmath.NewInt(1),
				Sender:         helpers.GenExternalAddr(s.chainName),
				Receiver:       helpers.GenExternalAddr(s.chainName),
				TargetIbc:      "",
				BridgerAddress: helpers.GenAccAddress().String(),
				ChainName:      s.chainName,
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.crosschainKeeper.SavePendingExecuteClaim(s.ctx, tt.claim)
			s.Require().NoError(err)

			claim, err := s.crosschainKeeper.GetPendingExecuteClaim(s.ctx, tt.claim.GetEventNonce())
			s.Require().NoError(err)
			s.Require().Equal(claim, tt.claim)
		})
	}
}
