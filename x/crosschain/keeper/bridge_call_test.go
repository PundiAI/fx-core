package keeper_test

import (
	"errors"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"go.uber.org/mock/gomock"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
)

func (s *KeeperTestSuite) TestBridgeCallHandler() {
	mockBridgeCallFn := func(msg *types.MsgBridgeCallClaim) {
		if len(msg.TokenContracts) == 0 {
			return
		}
		// set oracle nonce
		s.SetOracleSet(10, 100, 90)

		// set bridge token and mock erc20
		for i, c := range msg.TokenContracts {
			denom := helpers.NewRandDenom()
			s.crosschainKeeper.AddBridgeToken(s.ctx, c, fmt.Sprintf("%s%s", s.moduleName, c))

			coin := sdk.NewCoin(fmt.Sprintf("%s%s", s.moduleName, c), msg.Amounts[i])
			targetCoin := sdk.NewCoin(denom, msg.Amounts[i])
			s.erc20Keeper.EXPECT().ConvertDenomToTarget(gomock.Any(), gomock.Any(), coin, gomock.Any()).Return(targetCoin, nil).Times(1)
			s.erc20Keeper.EXPECT().ConvertCoin(gomock.Any(), gomock.Any()).Return(&erc20types.MsgConvertCoinResponse{}, nil).Times(1)
		}
		s.erc20Keeper.EXPECT().IsOriginOrConvertedDenom(gomock.Any(), gomock.Any()).Return(false).Times(len(msg.TokenContracts))

		// mock bank
		s.bankKeeper.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		s.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		// mock auth
		s.accountKeeper.EXPECT().GetModuleAddress(erc20types.ModuleName).Return(authtypes.NewEmptyModuleAccount(erc20types.ModuleName).GetAddress()).AnyTimes()
	}
	initMsg := func() *types.MsgBridgeCallClaim {
		return &types.MsgBridgeCallClaim{
			ChainName: s.moduleName,
			Sender:    helpers.GenerateAddressByModule(s.moduleName),
			Receiver:  helpers.GenerateAddressByModule(s.moduleName),
			To:        "",
			Data:      "",
			Value:     sdkmath.NewInt(0),
			TokenContracts: []string{
				helpers.GenerateAddressByModule(s.moduleName),
				helpers.GenerateAddressByModule(s.moduleName),
			},
			Amounts: []sdkmath.Int{
				sdkmath.NewInt(1e18),
				sdkmath.NewInt(1e18).Mul(sdkmath.NewInt(2)),
			},
			EventNonce:     10,
			BlockHeight:    100,
			BridgerAddress: sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
		}
	}

	tests := []struct {
		name   string
		mock   func(claim *types.MsgBridgeCallClaim)
		msgFn  func(msg *types.MsgBridgeCallClaim) *types.MsgBridgeCallClaim
		error  string
		refund bool
	}{
		{
			name: "ok - pass",
			msgFn: func(msg *types.MsgBridgeCallClaim) *types.MsgBridgeCallClaim {
				return msg
			},
		},
		{
			name: "ok - pass - no token",
			msgFn: func(msg *types.MsgBridgeCallClaim) *types.MsgBridgeCallClaim {
				msg.TokenContracts = []string{}
				msg.Amounts = []sdkmath.Int{}
				return msg
			},
		},
		{
			name: "ok - call evm error refund",
			msgFn: func(msg *types.MsgBridgeCallClaim) *types.MsgBridgeCallClaim {
				msg.To = helpers.GenerateAddressByModule(s.moduleName)
				return msg
			},
			mock: func(msg *types.MsgBridgeCallClaim) {
				// set height
				s.crosschainKeeper.SetLastObservedBlockHeight(s.ctx, 1000, msg.BlockHeight-1)
				// mock evm
				s.evmKeeper.EXPECT().CallEVM(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("call evm error")).Times(1)
			},
			refund: true,
		},
	}
	for _, t := range tests {
		s.Run(t.name, func() {
			msg := t.msgFn(initMsg())

			// mock msg
			mockBridgeCallFn(msg)

			s.accountKeeper.EXPECT().GetAccount(gomock.Any(), msg.GetSenderAddr().Bytes()).Return(nil).Times(1)
			if t.mock != nil {
				t.mock(msg)
			}

			// call
			err := s.crosschainKeeper.BridgeCallHandler(s.ctx, msg)

			// check
			if len(t.error) > 0 {
				s.EqualError(err, t.error)
			} else {
				s.NoError(err)
			}

			if t.refund {
				refundEvent := false
				for _, event := range s.ctx.EventManager().Events().ToABCIEvents() {
					if event.Type == types.EventTypeBridgeCallRefund {
						refundEvent = true
					}
				}
				s.True(refundEvent)

			}
		})
	}
}
