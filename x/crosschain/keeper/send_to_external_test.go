package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	"go.uber.org/mock/gomock"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
)

func (s *KeeperTestSuite) TestSendToExternal() {
	bridgeTokenAddress := helpers.GenHexAddress().String()
	bridgeToken := s.AddBridgeToken(bridgeTokenAddress)

	senderAddr := helpers.GenAccAddress()
	sendMsg := types.MsgSendToExternal{
		Sender:    senderAddr.String(),
		Amount:    sdk.NewCoin("usdt", sdk.NewInt(int64(tmrand.Uint32()))),
		BridgeFee: sdk.NewCoin("usdt", sdk.NewInt(int64(tmrand.Uint32()))),
		ChainName: s.moduleName,
	}
	s.erc20Keeper.EXPECT().ConvertDenomToTarget(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sdk.NewCoin(bridgeToken.Denom, sendMsg.Amount.Amount), erc20types.ErrInsufficientLiquidity).
		Times(1)
	external, err := s.msgServer.SendToExternal(s.ctx, &sendMsg)
	require.NoError(s.T(), err)
	require.EqualValues(s.T(), 1, external.OutgoingTxId)

	pendingSendToExternal, found := s.crosschainKeeper.GetPendingPoolTxById(s.ctx, external.OutgoingTxId)
	require.True(s.T(), found)
	require.EqualValues(s.T(), sendMsg.Amount, pendingSendToExternal.Token)
	require.EqualValues(s.T(), sendMsg.BridgeFee, pendingSendToExternal.Fee)

	// add liquidity
	erc20ModuleAddr := helpers.GenAccAddress()
	s.accountKeeper.EXPECT().GetModuleAddress(erc20types.ModuleName).Return(erc20ModuleAddr).Times(1)
	sendToken := sdk.NewCoin(bridgeToken.Denom, sendMsg.Amount.Amount.Add(sendMsg.BridgeFee.Amount))
	s.bankKeeper.EXPECT().HasBalance(gomock.Any(), erc20ModuleAddr, sendToken).Return(true).Times(1)
	s.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), erc20types.ModuleName, senderAddr, sdk.NewCoins(sendToken)).Return(nil).Times(1)
	s.erc20Keeper.EXPECT().IsOriginOrConvertedDenom(gomock.Any(), bridgeToken.Denom).Return(false).Times(1)

	s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), senderAddr, s.moduleName, sdk.NewCoins(sendToken)).Return(nil).Times(1)
	s.bankKeeper.EXPECT().BurnCoins(gomock.Any(), s.moduleName, sdk.NewCoins(sendToken)).Return(nil).Times(1)

	s.crosschainKeeper.HandlePendingOutgoingTx(s.ctx, helpers.GenHexAddress().Bytes(), 1, bridgeToken)

	// check pending send to external tx is removed
	_, found = s.crosschainKeeper.GetPendingPoolTxById(s.ctx, external.OutgoingTxId)
	s.Require().False(found)
}
