package keeper_test

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"go.uber.org/mock/gomock"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
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

func (s *KeeperTestSuite) TestKeeper_BridgeCallCoinsToERC20Token() {
	type Data struct {
		sender sdk.AccAddress
		coin   sdk.Coin
	}
	tests := []struct {
		name    string
		data    Data
		mock    func(data Data) (want types.ERC20Token)
		wantErr bool
	}{
		{
			name: "success - FX",
			data: Data{
				sender: helpers.GenAccAddress(),
				coin:   sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1)),
			},
			mock: func(data Data) (want types.ERC20Token) {
				s.erc20Keeper.EXPECT().ConvertDenomToTarget(gomock.Any(), data.sender, data.coin, fxtypes.ParseFxTarget(s.moduleName)).Return(data.coin, nil).Times(1)
				s.erc20Keeper.EXPECT().IsOriginOrConvertedDenom(gomock.Any(), data.coin.Denom).Return(true).Times(1)
				s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), data.sender, s.moduleName, sdk.NewCoins(data.coin)).Return(nil).Times(1)

				return types.ERC20Token{
					Contract: s.wfxTokenAddr,
					Amount:   data.coin.Amount,
				}
			},
			wantErr: false,
		},
		{
			name: "success - bridge denom",
			data: Data{
				sender: helpers.GenAccAddress(),
				coin:   sdk.NewCoin(types.NewBridgeDenom(s.moduleName, helpers.GenExternalAddr(s.moduleName)), sdk.NewInt(1)),
			},
			mock: func(data Data) (want types.ERC20Token) {
				contract := data.coin.Denom[len(s.moduleName):]
				s.AddBridgeToken(contract)

				s.erc20Keeper.EXPECT().ConvertDenomToTarget(gomock.Any(), data.sender, data.coin, fxtypes.ParseFxTarget(s.moduleName)).Return(data.coin, nil).Times(1)
				s.erc20Keeper.EXPECT().IsOriginOrConvertedDenom(gomock.Any(), data.coin.Denom).Return(false).Times(1)
				s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), data.sender, s.moduleName, sdk.NewCoins(data.coin)).Return(nil).Times(1)
				s.bankKeeper.EXPECT().BurnCoins(gomock.Any(), s.moduleName, sdk.NewCoins(data.coin)).Return(nil).Times(1)

				return types.ERC20Token{
					Contract: contract,
					Amount:   data.coin.Amount,
				}
			},
			wantErr: false,
		},
		{
			name: "success - base denom",
			data: Data{
				sender: helpers.GenAccAddress(),
				coin:   sdk.NewCoin("usdt", sdk.NewInt(1)),
			},
			mock: func(data Data) (want types.ERC20Token) {
				contract := helpers.GenHexAddress().String()
				bridgeToken := s.AddBridgeToken(contract)

				targetCoin := sdk.NewCoin(bridgeToken.Denom, data.coin.Amount)
				s.erc20Keeper.EXPECT().ConvertDenomToTarget(gomock.Any(), data.sender, data.coin, fxtypes.ParseFxTarget(s.moduleName)).Return(targetCoin, nil).Times(1)
				s.erc20Keeper.EXPECT().IsOriginOrConvertedDenom(gomock.Any(), targetCoin.Denom).Return(false).Times(1)
				s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), data.sender, s.moduleName, sdk.NewCoins(targetCoin)).Return(nil).Times(1)
				s.bankKeeper.EXPECT().BurnCoins(gomock.Any(), s.moduleName, sdk.NewCoins(targetCoin)).Return(nil).Times(1)

				return types.ERC20Token{
					Contract: contract,
					Amount:   data.coin.Amount,
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			want := tt.mock(tt.data)
			got, err := s.crosschainKeeper.BridgeCallCoinsToERC20Token(s.ctx, tt.data.sender, sdk.NewCoins(tt.data.coin))
			if (err != nil) != tt.wantErr {
				s.T().Errorf("BridgeCallCoinsToERC20Token() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, []types.ERC20Token{want}) {
				s.T().Errorf("BridgeCallCoinsToERC20Token() got = %v, want %v", got, want)
			}
		})
	}
}
