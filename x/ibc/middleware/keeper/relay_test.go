package keeper_test

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestOnRecvPacket() {
	testCases := []struct {
		name          string
		malleate      func(packet *channeltypes.Packet)
		expPass       bool
		errorStr      string
		checkBalance  bool
		checkCoinAddr common.Address
		expCoins      sdk.Coins
		afterFn       func(packetData transfertypes.FungibleTokenPacketData)
	}{
		{
			name: "pass - normal - ibc transfer packet",
		},
		{
			name: "pass - normal - receive address is 0xAddress, coin is DefaultCoin",
		},
		{
			name: "pass - normal - receive address is 0xAddress",
		},
		{
			name: "error - normal - receive address is 0xAddress but coin not registered",
		},
		{
			name: "pass - any memo",
		},
		{
			name: "pass - ibc call evm",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
		})
	}
}

func (suite *KeeperTestSuite) TestOnAcknowledgementPacket() {
	coin := sdk.NewCoin("stake", sdkmath.NewInt(tmrand.Int63n(100000000000)))

	testCases := []struct {
		name         string
		malleate     func(packet *channeltypes.Packet, ack *channeltypes.Acknowledgement)
		expPass      bool
		errorStr     string
		checkBalance bool
		expCoins     sdk.Coins
	}{
		{
			"pass - success ack - ibc transfer packet",
			func(packet *channeltypes.Packet, ack *channeltypes.Acknowledgement) {
				escrowAddress := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
				suite.mintCoin(escrowAddress, coin)
			},
			true,
			"",
			true,
			sdk.Coins{},
		},
		{
			"pass - error ack - ibc transfer packet",
			func(packet *channeltypes.Packet, ack *channeltypes.Acknowledgement) {
				*ack = channeltypes.NewErrorAcknowledgement(fmt.Errorf("test"))

				escrowAddress := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
				suite.mintCoin(escrowAddress, coin)
				suite.Require().NoError(suite.App.Erc20Keeper.SetCache(suite.Ctx, crosschaintypes.NewIBCTransferKey(packet.SourceChannel, 1), coin.Amount))
			},
			true,
			"",
			true,
			sdk.NewCoins(coin),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			senderAddr, packet := suite.mockPacket(coin, "")

			ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
			tc.malleate(&packet, &ack)

			err := suite.ibcMiddleware.OnAcknowledgementPacket(suite.Ctx, packet, ack.Acknowledgement(), nil)
			if tc.expPass {
				suite.Require().NoError(err, "packet: %s", string(packet.GetData()))
			} else {
				suite.Require().NotNil(err)
				suite.Require().Equalf(tc.errorStr, err.Error(), "packet: %s", string(packet.GetData()))
			}
			if tc.checkBalance {
				senderAddrCoins := suite.bankKeeper.GetAllBalances(suite.Ctx, senderAddr)
				suite.Require().Equal(tc.expCoins.String(), senderAddrCoins.String())
			}
		})
	}
}

func (suite *KeeperTestSuite) mockPacket(coin sdk.Coin, memo string) (sdk.AccAddress, channeltypes.Packet) {
	senderAddr := helpers.GenAccAddress()
	receiveAddr := helpers.GenAccAddress()

	suite.ibcTransferKeeper.SetTotalEscrowForDenom(suite.Ctx, coin)
	packetData := transfertypes.NewFungibleTokenPacketData(coin.Denom, coin.Amount.String(), senderAddr.String(), receiveAddr.String(), memo)
	packet := channeltypes.NewPacket(packetData.GetBytes(),
		1, transfertypes.PortID, "channel-0", transfertypes.PortID, "channel-0",
		clienttypes.Height{RevisionNumber: 100, RevisionHeight: 100000}, 0,
	)
	return senderAddr, packet
}

func (suite *KeeperTestSuite) TestOnTimeoutPacket() {
	coin := sdk.NewCoin("stake", sdkmath.NewInt(tmrand.Int63n(100000000000)))
	testCases := []struct {
		name         string
		malleate     func(packet *channeltypes.Packet)
		expPass      bool
		errorStr     string
		checkBalance bool
		expCoins     sdk.Coins
	}{
		{
			"pass - normal - ibc transfer packet",
			func(packet *channeltypes.Packet) {
			},
			true,
			"",
			true,
			sdk.Coins{},
		},
		{
			"pass - normal - ibc mint token - router is empty",
			func(packet *channeltypes.Packet) {
			},
			true,
			"",
			true,
			sdk.Coins{},
		},
		{
			"error - escrow address insufficient 10coin",
			func(packet *channeltypes.Packet) {
			},
			false,
			fmt.Sprintf("unable to unescrow tokens, this may be caused by a malicious counterparty module or a bug: please open an issue on counterparty module: spendable balance %d%s is smaller than %d%s: insufficient funds", coin.Amount.Sub(sdkmath.NewInt(10)).Uint64(), coin.Denom, coin.Amount.Uint64(), coin.Denom),
			true,
			sdk.Coins{},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
		})
	}
}

func (suite *KeeperTestSuite) mintCoin(address sdk.AccAddress, coins ...sdk.Coin) {
	suite.Require().NoError(suite.bankKeeper.MintCoins(suite.Ctx, transfertypes.ModuleName, coins))
	suite.Require().NoError(suite.bankKeeper.SendCoinsFromModuleToAccount(suite.Ctx, transfertypes.ModuleName, address, coins))
}
