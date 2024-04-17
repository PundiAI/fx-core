package tests_test

import (
	"encoding/hex"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/keeper"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestKeeper_BridgeCallRefund() {
	suite.bondedOracle()
	suite.Commit()

	bridgeToken := helpers.GenerateAddress()
	bridgeTokenStr := fxtypes.AddressToStr(bridgeToken.Bytes(), suite.chainName)
	suite.addBridgeToken(bridgeTokenStr, fxtypes.GetCrossChainMetadataManyToOne("test token", "TTT", 18))

	suite.registerCoin(keeper.NewBridgeDenom(suite.chainName, bridgeTokenStr))

	fxAddr1 := helpers.GenerateAddress()
	randomBlock := tmrand.Int63n(1000000000)
	randomAmount := tmrand.Int63n(1000000000)
	_, err := suite.MsgServer().SendToFxClaim(suite.ctx, &types.MsgSendToFxClaim{
		EventNonce:     suite.Keeper().GetLastEventNonceByOracle(suite.ctx, suite.oracleAddrs[0]) + 1,
		BlockHeight:    uint64(randomBlock),
		TokenContract:  bridgeTokenStr,
		Amount:         sdk.NewInt(randomAmount),
		Sender:         fxtypes.AddressToStr(helpers.GenerateAddress().Bytes(), suite.chainName),
		Receiver:       sdk.AccAddress(fxAddr1.Bytes()).String(),
		TargetIbc:      "",
		BridgerAddress: suite.bridgerAddrs[0].String(),
	})
	suite.NoError(err)

	pair, b := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, "ttt")
	suite.True(b)
	suite.Equal(sdkmath.NewInt(randomAmount), suite.app.BankKeeper.GetBalance(suite.ctx, fxAddr1.Bytes(), pair.Denom).Amount)

	suite.NoError(err)
	_, err = suite.MsgServer().BridgeCall(suite.ctx, &types.MsgBridgeCall{
		Sender:   sdk.AccAddress(fxAddr1.Bytes()).String(),
		Receiver: helpers.GenerateAddressByModule(suite.chainName),
		Coins:    sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), sdkmath.NewInt(randomAmount))),
		Value:    sdk.ZeroInt(),
	})
	suite.NoError(err)

	suite.Equal(sdkmath.NewInt(0), suite.app.BankKeeper.GetBalance(suite.ctx, fxAddr1.Bytes(), pair.Denom).Amount)

	var outgoingBridgeCall *types.OutgoingBridgeCall
	suite.Keeper().IterateOutgoingBridgeCallsByAddress(suite.ctx, fxtypes.AddressToStr(fxAddr1.Bytes(), suite.chainName), func(value *types.OutgoingBridgeCall) bool {
		outgoingBridgeCall = value
		return true
	})
	suite.NotNil(outgoingBridgeCall)

	// Triggering the SendtoFx claim once is just to trigger timeout
	_, err = suite.MsgServer().SendToFxClaim(suite.ctx, &types.MsgSendToFxClaim{
		EventNonce:     suite.Keeper().GetLastEventNonceByOracle(suite.ctx, suite.oracleAddrs[0]) + 1,
		BlockHeight:    outgoingBridgeCall.Timeout,
		TokenContract:  bridgeTokenStr,
		Amount:         sdk.NewInt(randomAmount),
		Sender:         fxtypes.AddressToStr(helpers.GenerateAddress().Bytes(), suite.chainName),
		Receiver:       sdk.AccAddress(fxAddr1.Bytes()).String(),
		TargetIbc:      hex.EncodeToString([]byte(fxtypes.ERC20Target)),
		BridgerAddress: suite.bridgerAddrs[0].String(),
	})
	suite.NoError(err)
	// expect balance = sendToFx value + outgointBridgeCallRefund value
	suite.checkBalanceOf(pair.GetERC20Contract(), fxAddr1, big.NewInt(randomAmount+randomAmount))
}

func (suite *KeeperTestSuite) checkBalanceOf(contractAddr, address common.Address, expectBalance *big.Int) {
	var balanceRes struct {
		Value *big.Int
	}
	err := suite.app.EvmKeeper.QueryContract(suite.ctx, contractAddr, contractAddr, contract.GetFIP20().ABI, "balanceOf", &balanceRes, address)
	suite.Require().NoError(err)
	suite.EqualValuesf(expectBalance.Cmp(balanceRes.Value), 0, "expect balance %s, got %s", expectBalance.String(), balanceRes.Value.String())
}
