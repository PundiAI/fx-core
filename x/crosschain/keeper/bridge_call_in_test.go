package keeper_test

import (
	"errors"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/mock/gomock"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
)

func (s *KeeperMockSuite) TestBridgeCallHandler() {
	callEvmMock := func(msg *types.MsgBridgeCallClaim, sender common.Address, getTokenPairTimes int) {
		s.crosschainKeeper.SetLastObservedBlockHeight(s.ctx, 1000, msg.BlockHeight-1)

		s.erc20Keeper.EXPECT().GetTokenPair(gomock.Any(), gomock.Any()).Return(erc20types.TokenPair{}, true).
			Times(getTokenPairTimes)
		contract := types.ExternalAddrToHexAddr(msg.ChainName, msg.To)
		s.evmKeeper.EXPECT().IsContract(gomock.Any(), contract).Return(true).Times(1)
		s.evmKeeper.EXPECT().CallEVM(gomock.Any(),
			sender,
			&contract,
			big.NewInt(0),
			uint64(types.MaxGasLimit),
			gomock.Any(),
			true,
		).Return(nil, errors.New("call evm error")).Times(1)
	}

	tests := []struct {
		name       string
		initData   func(msg *types.MsgBridgeCallClaim)
		customMock func(msg *types.MsgBridgeCallClaim)
		error      string
		refund     bool
	}{
		{
			name: "ok - pass",
			customMock: func(msg *types.MsgBridgeCallClaim) {
				s.evmKeeper.EXPECT().IsContract(gomock.Any(), gomock.Any()).Return(false).Times(1)
			},
		},
		{
			name: "ok - pass - no token",
			initData: func(msg *types.MsgBridgeCallClaim) {
				msg.TokenContracts = []string{}
				msg.Amounts = []sdkmath.Int{}
			},
			customMock: func(msg *types.MsgBridgeCallClaim) {
				s.evmKeeper.EXPECT().IsContract(gomock.Any(), gomock.Any()).Return(false).Times(1)
			},
		},
		{
			name: "ok - call evm error refund",
			initData: func(msg *types.MsgBridgeCallClaim) {
				msg.To = helpers.GenExternalAddr(s.moduleName)
			},
			customMock: func(msg *types.MsgBridgeCallClaim) {
				// data = "transfer(address,uint256)" "0x0000000000000000000000000000000000000000" 1
				msg.Data = "a9059cbb00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"
				callEvmMock(msg, s.crosschainKeeper.GetCallbackFrom(), len(msg.TokenContracts))
			},
			refund: true,
		},
		{
			name: "ok - memo is send to call evm",
			initData: func(msg *types.MsgBridgeCallClaim) {
				msg.Memo = "0000000000000000000000000000000000000000000000000000000000010000"
				// data = "transfer(address,uint256)" "0x0000000000000000000000000000000000000000" 1
				msg.Data = "a9059cbb00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"
			},
			customMock: func(msg *types.MsgBridgeCallClaim) {
				callEvmMock(msg, types.ExternalAddrToHexAddr(msg.ChainName, msg.Sender), 0)
			},
		},
	}
	for _, t := range tests {
		s.Run(t.name, func() {
			msg := &types.MsgBridgeCallClaim{
				ChainName: s.moduleName,
				Sender:    helpers.GenExternalAddr(s.moduleName),
				Refund:    helpers.GenExternalAddr(s.moduleName),
				To:        helpers.GenExternalAddr(s.moduleName),
				Data:      "",
				Memo:      "",
				Value:     sdkmath.NewInt(0),
				TokenContracts: []string{
					helpers.GenExternalAddr(s.moduleName),
					helpers.GenExternalAddr(s.moduleName),
				},
				Amounts: []sdkmath.Int{
					sdkmath.NewInt(1e18),
					sdkmath.NewInt(1e18).Mul(sdkmath.NewInt(2)),
				},
				EventNonce:     10,
				BlockHeight:    100,
				BridgerAddress: helpers.GenAccAddress().String(),
				TxOrigin:       helpers.GenExternalAddr(s.moduleName),
			}
			if t.initData != nil {
				t.initData(msg)
			}
			s.accountKeeper.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			s.accountKeeper.EXPECT().NewAccountWithAddress(gomock.Any(), gomock.Any()).Times(1)

			s.MockBridgeCallToken(msg.GetERC20Tokens())

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
					if event.Type == types.EventTypeBridgeCallRefundOut {
						refundEvent = true
					}
				}
				s.True(refundEvent)
			}
		})
	}
}

func (s *KeeperMockSuite) Test_CoinsToBridgeCallTokens() {
	input := sdk.Coins{
		sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1e18)),
		sdk.NewCoin("aaa", sdkmath.NewInt(2e18)),
	}
	s.erc20Keeper.EXPECT().GetTokenPair(gomock.Any(), "aaa").Return(erc20types.TokenPair{
		Erc20Address: "0x0000000000000000000000000000000000000001",
	}, true).Times(1)
	tokens, amounts := s.crosschainKeeper.CoinsToBridgeCallTokens(s.ctx, input)
	s.Require().EqualValues(len(tokens), len(amounts))
	s.Len(tokens, input.Len())

	expectTokens := []common.Address{
		common.HexToAddress("0x0000000000000000000000000000000000000000"),
		common.HexToAddress("0x0000000000000000000000000000000000000001"),
	}
	expectAmount := []*big.Int{big.NewInt(1e18), big.NewInt(2e18)}
	s.Require().EqualValues(expectTokens, tokens)
	s.Require().EqualValues(expectAmount, amounts)
}

func (s *KeeperMockSuite) MockBridgeCallToken(erc20Tokens []types.ERC20Token) {
	if len(erc20Tokens) == 0 {
		return
	}
	s.accountKeeper.EXPECT().GetModuleAddress(erc20types.ModuleName).Return(authtypes.NewEmptyModuleAccount(erc20types.ModuleName).GetAddress()).AnyTimes()

	s.erc20Keeper.EXPECT().IsOriginOrConvertedDenom(gomock.Any(), gomock.Any()).Return(false).Times(len(erc20Tokens))
	s.bankKeeper.EXPECT().MintCoins(gomock.Any(), s.moduleName, gomock.Any()).Return(nil).Times(1)
	s.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), s.moduleName, gomock.Any(), gomock.Any()).Return(nil).Times(1)

	for _, erc20Token := range erc20Tokens {
		baseDenom := helpers.NewRandDenom()
		bridgeToken := s.AddBridgeToken(erc20Token.Contract)
		bridgeCoin := sdk.NewCoin(bridgeToken.Denom, erc20Token.Amount)
		targetCoin := sdk.NewCoin(baseDenom, erc20Token.Amount)

		s.erc20Keeper.EXPECT().ConvertDenomToTarget(gomock.Any(), gomock.Any(), bridgeCoin, gomock.Any()).Return(targetCoin, nil).Times(1)
		s.erc20Keeper.EXPECT().ConvertCoin(gomock.Any(), gomock.Any()).Return(&erc20types.MsgConvertCoinResponse{}, nil).Times(1)
	}
}
