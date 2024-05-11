package keeper_test

import (
	"errors"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"go.uber.org/mock/gomock"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
)

func (s *KeeperTestSuite) TestBridgeCallHandler() {
	tests := []struct {
		name       string
		initMsg    func(msg *types.MsgBridgeCallClaim)
		customMock func(msg *types.MsgBridgeCallClaim)
		error      string
		refund     bool
	}{
		{
			name: "ok - pass",
		},
		{
			name: "ok - pass - no token",
			initMsg: func(msg *types.MsgBridgeCallClaim) {
				msg.TokenContracts = []string{}
				msg.Amounts = []sdkmath.Int{}
			},
		},
		{
			name: "ok - call evm error refund",
			initMsg: func(msg *types.MsgBridgeCallClaim) {
				msg.To = helpers.GenerateAddressByModule(s.moduleName)
			},
			customMock: func(msg *types.MsgBridgeCallClaim) {
				s.crosschainKeeper.SetLastObservedBlockHeight(s.ctx, 1000, msg.BlockHeight-1)

				sender := types.ExternalAddrToHexAddr(msg.ChainName, msg.Sender)
				contract := types.ExternalAddrToHexAddr(msg.ChainName, msg.To)
				s.evmKeeper.EXPECT().CallEVM(gomock.Any(),
					sender,
					&contract,
					big.NewInt(0),
					uint64(BlockGasLimit),
					[]byte{},
					true,
				).Return(nil, errors.New("call evm error")).Times(1)
			},
			refund: true,
		},
	}
	for _, t := range tests {
		s.Run(t.name, func() {
			msg := &types.MsgBridgeCallClaim{
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
			if t.initMsg != nil {
				t.initMsg(msg)
			}

			if len(msg.TokenContracts) != 0 {
				s.accountKeeper.EXPECT().GetModuleAddress(erc20types.ModuleName).Return(authtypes.NewEmptyModuleAccount(erc20types.ModuleName).GetAddress()).AnyTimes()

				s.erc20Keeper.EXPECT().IsOriginOrConvertedDenom(gomock.Any(), gomock.Any()).Return(false).Times(len(msg.TokenContracts))
				s.bankKeeper.EXPECT().MintCoins(gomock.Any(), msg.ChainName, gomock.Any()).Return(nil).Times(1)
				s.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), msg.ChainName, gomock.Any(), gomock.Any()).Return(nil).Times(1)

				for i, contract := range msg.TokenContracts {
					amount := msg.Amounts[i]
					baseDenom := helpers.NewRandDenom()
					bridgeToken := s.AddBridgeToken(contract)
					bridgeCoin := sdk.NewCoin(bridgeToken.Denom, amount)
					targetCoin := sdk.NewCoin(baseDenom, amount)

					s.erc20Keeper.EXPECT().ConvertDenomToTarget(gomock.Any(), gomock.Any(), bridgeCoin, gomock.Any()).Return(targetCoin, nil).Times(1)
					s.erc20Keeper.EXPECT().ConvertCoin(gomock.Any(), gomock.Any()).Return(&erc20types.MsgConvertCoinResponse{}, nil).Times(1)
				}
			}

			s.accountKeeper.EXPECT().GetAccount(gomock.Any(), msg.GetSenderAddr().Bytes()).Return(nil).Times(1)

			if t.customMock != nil {
				t.customMock(msg)
			}

			err := s.crosschainKeeper.BridgeCallHandler(s.ctx, msg)
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
