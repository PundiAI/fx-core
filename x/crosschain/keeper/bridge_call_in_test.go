package keeper_test

import (
	"math/big"
	"strings"

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

func (suite *KeeperTestSuite) TestBridgeCallHandler() {
	testCases := []struct {
		Name              string
		Msg               types.MsgBridgeCallClaim
		TokenIsNativeCoin []bool
		Success           bool
		CallContract      bool
	}{
		{
			Name: "success - token",
			Msg: types.MsgBridgeCallClaim{
				ChainName:      suite.chainName,
				BridgerAddress: helpers.GenAccAddress().String(),
				EventNonce:     1,
				BlockHeight:    1,
				Sender:         helpers.GenExternalAddr(suite.chainName),
				Refund:         helpers.GenExternalAddr(suite.chainName),
				TokenContracts: []string{
					helpers.GenExternalAddr(suite.chainName),
					helpers.GenExternalAddr(suite.chainName),
				},
				Amounts: []sdkmath.Int{
					helpers.NewRandAmount(),
					helpers.NewRandAmount(),
				},
				To:       helpers.GenExternalAddr(suite.chainName),
				Data:     "",
				Value:    sdkmath.ZeroInt(),
				Memo:     "",
				TxOrigin: helpers.GenExternalAddr(suite.chainName),
			},
			TokenIsNativeCoin: []bool{true, true},
			Success:           true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.Name, func() {
			_, _, erc20Addrs := suite.BridgeCallClaimInitialize(tc.Msg, tc.TokenIsNativeCoin)

			err := suite.Keeper().BridgeCallHandler(suite.Ctx, &tc.Msg)
			if tc.Success {
				suite.Require().NoError(err)
				if !tc.CallContract {
					for i, addr := range erc20Addrs {
						suite.CheckBalanceOf(addr, tc.Msg.GetToAddr(), tc.Msg.Amounts[i].BigInt())
					}
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) BridgeCallClaimInitialize(msg types.MsgBridgeCallClaim, tokenIsNativeCoin []bool) (baseDenoms, bridgeDenoms []string, erc20Addrs []common.Address) {
	suite.Require().Equal(len(tokenIsNativeCoin), len(msg.TokenContracts))

	baseDenoms = make([]string, 0, len(msg.TokenContracts))
	bridgeDenoms = make([]string, 0, len(msg.TokenContracts))
	erc20Addrs = make([]common.Address, 0, len(msg.TokenContracts))
	for i, c := range msg.TokenContracts {
		baseDenom := helpers.NewRandDenom()
		bridgeDenom := types.NewBridgeDenom(suite.chainName, c)
		suite.SetToken(strings.ToUpper(baseDenom), bridgeDenom)
		suite.AddBridgeToken(c, strings.ToLower(baseDenom))
		erc20Addr := suite.AddTokenPair(baseDenom, tokenIsNativeCoin[i])

		baseDenoms = append(baseDenoms, baseDenom)
		bridgeDenoms = append(bridgeDenoms, bridgeDenom)
		erc20Addrs = append(erc20Addrs, erc20Addr)

		if !tokenIsNativeCoin[i] {
			suite.MintTokenToModule(suite.chainName, sdk.NewCoin(bridgeDenom, msg.Amounts[i]))
		}
	}
	return baseDenoms, bridgeDenoms, erc20Addrs
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
