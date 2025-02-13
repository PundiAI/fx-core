package keeper_test

import (
	"bytes"
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/jsonpb"
	"github.com/cosmos/gogoproto/proto"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	"github.com/pundiai/fx-core/v8/x/ibc/middleware/keeper"
	"github.com/pundiai/fx-core/v8/x/ibc/middleware/types"
)

func (suite *KeeperTestSuite) TestOnRecvPacket() {
	testCases := []struct {
		name         string
		coin         sdk.Coin
		isOurCoin    bool
		malleate     func(packet channeltypes.Packet, packetData transfertypes.FungibleTokenPacketData) channeltypes.Packet
		err          error
		checkBalance bool
		expCoin      sdk.Coin
	}{
		{
			name:         "pass - send FX to ours",
			coin:         sdk.NewCoin(fxtypes.LegacyFXDenom, sdkmath.NewInt(1000)),
			isOurCoin:    true,
			checkBalance: true,
			expCoin:      sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10)),
		},
		{
			name:      "pass - mainnet pundix chain send pundix to ours",
			coin:      sdk.NewCoin(fxtypes.MainnetPundixUnWrapDenom, sdkmath.NewInt(1000)),
			isOurCoin: true,
			malleate: func(packet channeltypes.Packet, _ transfertypes.FungibleTokenPacketData) channeltypes.Packet {
				suite.Ctx = suite.Ctx.WithChainID(fxtypes.MainnetChainId)
				return packet
			},
			checkBalance: true,
			expCoin:      sdk.NewCoin(fxtypes.PundixWrapDenom, sdkmath.NewInt(1000)),
		},
		{
			name:      "pass - testnet pundix chain send pundix to ours",
			coin:      sdk.NewCoin(fxtypes.TestnetPundixUnWrapDenom, sdkmath.NewInt(1000)),
			isOurCoin: true,
			malleate: func(packet channeltypes.Packet, _ transfertypes.FungibleTokenPacketData) channeltypes.Packet {
				suite.Ctx = suite.Ctx.WithChainID(fxtypes.TestnetChainId)
				return packet
			},
			checkBalance: true,
			expCoin:      sdk.NewCoin(fxtypes.PundixWrapDenom, sdkmath.NewInt(1000)),
		},
		{
			name:         "pass - send defaultDenom to ours",
			coin:         sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000)),
			isOurCoin:    true,
			checkBalance: true,
			expCoin:      sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000)),
		},
		{
			name:         "pass - send dest chain coin to ours",
			coin:         sdk.NewCoin("stake", sdkmath.NewInt(1000)),
			isOurCoin:    false,
			checkBalance: false,
		},
		{
			name:      "failed - dest address is error",
			coin:      sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000)),
			isOurCoin: false,
			malleate: func(packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) channeltypes.Packet {
				packetData := packetDataToFxPacketData(data)
				packetData.Router = "tron"
				packetData.Receiver = "errAddr"
				packet.Data = mustProtoMarshalJSON(&packetData)
				return packet
			},
			err: fmt.Errorf("wrong length"),
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.ibcMiddleware.Keeper = suite.App.IBCMiddlewareKeeper.SetCrosschainKeeper(mockCrosschainKeeper{})
			packet, packetData := suite.mockOnRecvPacket(tc.coin, tc.isOurCoin)
			if tc.malleate != nil {
				packet = tc.malleate(packet, packetData)
			}
			if tc.isOurCoin {
				suite.mintCoinEscrowAddr(packet.DestinationChannel, tc.expCoin)
			}
			acknowledgement := suite.ibcMiddleware.OnRecvPacket(suite.Ctx, packet, sdk.AccAddress{})
			if tc.err != nil {
				suite.Require().False(acknowledgement.Success(), tc.err.Error())
				suite.Contains(tc.err.Error(), getErrByEvent(suite.Ctx))
				return
			}
			suite.Require().True(acknowledgement.Success(), suite.Ctx.EventManager().Events())
			if tc.checkBalance {
				suite.AssertBalance(common.HexToAddress(packetData.Receiver).Bytes(), tc.expCoin)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOnAcknowledgementPacket() {
	coin := sdk.NewCoin("stake", sdkmath.NewInt(tmrand.Int63n(100000000000)))

	testCases := []struct {
		name         string
		malleate     func(packet *channeltypes.Packet, ack *channeltypes.Acknowledgement, packetData transfertypes.FungibleTokenPacketData)
		expPass      bool
		checkBalance bool
		expCoin      sdk.Coin
	}{
		{
			name:         "pass - success ack - ibc transfer packet",
			expPass:      true,
			checkBalance: false,
			expCoin:      sdk.Coin{},
		},
		{
			name: "pass - error ack - ibc transfer packet",
			malleate: func(packet *channeltypes.Packet, ack *channeltypes.Acknowledgement, _ transfertypes.FungibleTokenPacketData) {
				*ack = ackWithErr()

				suite.mintCoinEscrowAddr(packet.SourceChannel, coin)
				suite.Require().NoError(suite.App.Erc20Keeper.SetCache(suite.Ctx, crosschaintypes.NewIBCTransferKey(packet.SourceChannel, packet.Sequence), coin.Amount))
			},
			expPass:      true,
			checkBalance: true,
			expCoin:      coin,
		},
		{
			name: "pass - error ack - denom is FX",
			malleate: func(packet *channeltypes.Packet, ack *channeltypes.Acknowledgement, packetData transfertypes.FungibleTokenPacketData) {
				*ack = ackWithErr()

				packetData.Denom = fxtypes.LegacyFXDenom
				packet.Data = packetData.GetBytes()
				swapCoin := fxtypes.SwapCoin(sdk.NewCoin(packetData.Denom, coin.Amount))

				suite.mintCoinEscrowAddr(packet.SourceChannel, swapCoin)
				suite.Require().NoError(suite.App.Erc20Keeper.SetCache(suite.Ctx, crosschaintypes.NewIBCTransferKey(packet.SourceChannel, packet.Sequence), swapCoin.Amount))
			},
			expPass:      true,
			checkBalance: true,
			expCoin:      fxtypes.SwapCoin(sdk.NewCoin(fxtypes.LegacyFXDenom, coin.Amount)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			senderAddr, packet, packetData := suite.mockPacket(coin, "")

			ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
			if tc.malleate != nil {
				tc.malleate(&packet, &ack, packetData)
			}

			err := suite.ibcMiddleware.OnAcknowledgementPacket(suite.Ctx, packet, ack.Acknowledgement(), nil)
			if tc.expPass {
				suite.Require().NoError(err, "packet: %s", string(packet.GetData()))
			} else {
				suite.Require().Error(err)
			}
			if tc.checkBalance {
				suite.AssertBalance(senderAddr, tc.expCoin)
			}
		})
	}
}

func (suite *KeeperTestSuite) mockOnRecvPacket(coin sdk.Coin, ibcCoin bool) (channeltypes.Packet, transfertypes.FungibleTokenPacketData) {
	senderAddr := helpers.GenAccAddress()
	receiveAddr := helpers.GenHexAddress()
	destChannel := fmt.Sprintf("channel-%d", tmrand.Int63n(10000))
	sourceChannel := fmt.Sprintf("channel-%d", tmrand.Int63n(10000))
	denom := coin.Denom
	if ibcCoin {
		denom = transfertypes.GetPrefixedDenom(transfertypes.PortID, sourceChannel, denom)
	}

	packetData := transfertypes.NewFungibleTokenPacketData(denom, coin.Amount.String(), senderAddr.String(), receiveAddr.String(), "")
	packet := channeltypes.NewPacket(packetData.GetBytes(), uint64(tmrand.Int63n(1000000)), transfertypes.PortID, sourceChannel, transfertypes.PortID, destChannel,
		clienttypes.Height{RevisionNumber: 100, RevisionHeight: 100000}, 0,
	)
	return packet, packetData
}

func (suite *KeeperTestSuite) mockPacket(coin sdk.Coin, memo string) (sdk.AccAddress, channeltypes.Packet, transfertypes.FungibleTokenPacketData) {
	senderAddr := helpers.GenAccAddress()
	receiveAddr := helpers.GenAccAddress()

	packetData := transfertypes.NewFungibleTokenPacketData(coin.Denom, coin.Amount.String(), senderAddr.String(), receiveAddr.String(), memo)
	packet := channeltypes.NewPacket(packetData.GetBytes(),
		1, transfertypes.PortID, "channel-0", transfertypes.PortID, "channel-0",
		clienttypes.Height{RevisionNumber: 100, RevisionHeight: 100000}, 0,
	)
	return senderAddr, packet, packetData
}

func (suite *KeeperTestSuite) TestOnTimeoutPacket() {
	coin := sdk.NewCoin("stake", sdkmath.NewInt(tmrand.Int63n(100000000000)))
	testCases := []struct {
		name         string
		malleate     func(packet *channeltypes.Packet, packetData transfertypes.FungibleTokenPacketData)
		expPass      bool
		checkBalance bool
		expCoins     sdk.Coins
	}{
		{
			name: "pass - normal - ibc transfer packet",
			malleate: func(packet *channeltypes.Packet, _ transfertypes.FungibleTokenPacketData) {
				suite.mintCoinEscrowAddr(packet.SourceChannel, coin)
				suite.Require().NoError(suite.App.Erc20Keeper.SetCache(suite.Ctx, crosschaintypes.NewIBCTransferKey(packet.SourceChannel, packet.Sequence), coin.Amount))
			},
			expPass:      true,
			checkBalance: true,
			expCoins:     sdk.NewCoins(coin),
		},
		{
			name: "pass - normal - denom is FX",
			malleate: func(packet *channeltypes.Packet, packetData transfertypes.FungibleTokenPacketData) {
				packetData.Denom = fxtypes.LegacyFXDenom
				packet.Data = packetData.GetBytes()
				swapCoin := fxtypes.SwapCoin(sdk.NewCoin(packetData.Denom, coin.Amount))

				suite.mintCoinEscrowAddr(packet.SourceChannel, swapCoin)
				suite.Require().NoError(suite.App.Erc20Keeper.SetCache(suite.Ctx, crosschaintypes.NewIBCTransferKey(packet.SourceChannel, packet.Sequence), swapCoin.Amount))
			},
			expPass:      true,
			checkBalance: true,
			expCoins:     sdk.NewCoins(fxtypes.SwapCoin(sdk.NewCoin(fxtypes.LegacyFXDenom, coin.Amount))),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			senderAddr, packet, packetData := suite.mockPacket(coin, "")

			if tc.malleate != nil {
				tc.malleate(&packet, packetData)
			}
			err := suite.ibcMiddleware.OnTimeoutPacket(suite.Ctx, packet, nil)
			if tc.expPass {
				suite.Require().NoError(err, "packet: %s", string(packet.GetData()))
			} else {
				suite.Require().Error(err)
			}
			if tc.checkBalance {
				senderAddrCoins := suite.App.BankKeeper.GetAllBalances(suite.Ctx, senderAddr)
				suite.Require().Equal(tc.expCoins.String(), senderAddrCoins.String())
			}
		})
	}
}

func (suite *KeeperTestSuite) mintCoinEscrowAddr(channel string, coin sdk.Coin) {
	suite.MintToken(transfertypes.GetEscrowAddress(transfertypes.PortID, channel), sdk.NewCoins(coin)...)
	suite.App.IBCTransferKeeper.SetTotalEscrowForDenom(suite.Ctx, coin)
}

func TestUnmarshalAckPacketData(t *testing.T) {
	normalData := types.FungibleTokenPacketData{
		Denom:  fxtypes.LegacyFXDenom,
		Amount: "1000",
	}
	normalExpected := transfertypes.FungibleTokenPacketData{
		Denom:  fxtypes.DefaultDenom,
		Amount: "10",
	}

	testCases := []struct {
		name     string
		malleate func(data types.FungibleTokenPacketData, exp transfertypes.FungibleTokenPacketData) (types.FungibleTokenPacketData, transfertypes.FungibleTokenPacketData)
		exp      transfertypes.FungibleTokenPacketData
		isZero   bool
		err      error
	}{
		{
			name:   "normal",
			exp:    normalExpected,
			isZero: false,
		},
		{
			name: "normal - a pundiai is zero",
			malleate: func(data types.FungibleTokenPacketData, exp transfertypes.FungibleTokenPacketData) (types.FungibleTokenPacketData, transfertypes.FungibleTokenPacketData) {
				data.Amount = "99"
				exp.Amount = "0"
				return data, exp
			},
			isZero: true,
		},
		{
			name: "normal - amount + fee",
			malleate: func(data types.FungibleTokenPacketData, exp transfertypes.FungibleTokenPacketData) (types.FungibleTokenPacketData, transfertypes.FungibleTokenPacketData) {
				data.Amount = "1000"
				data.Fee = "100"
				data.Router = tmrand.Str(10)
				exp.Amount = "11"
				return data, exp
			},
			isZero: false,
		},
		{
			name: "normal - another denom with fee",
			malleate: func(data types.FungibleTokenPacketData, exp transfertypes.FungibleTokenPacketData) (types.FungibleTokenPacketData, transfertypes.FungibleTokenPacketData) {
				denom := tmrand.Str(10)
				data.Denom = denom
				data.Amount = "1000"
				data.Fee = "1000"
				data.Router = tmrand.Str(10)

				exp.Denom = denom
				exp.Amount = "2000"
				return data, exp
			},
			isZero: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tcData := normalData
			if tc.malleate != nil {
				tcData, tc.exp = tc.malleate(tcData, normalExpected)
			}
			packetDateByte, err := sdk.SortJSON(mustProtoMarshalJSON(&tcData))
			require.NoError(t, err)
			data, isZeroAmount, err := keeper.UnmarshalAckPacketData(packetDateByte)
			if tc.err != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.exp, data)
			require.Equal(t, tc.isZero, isZeroAmount)
		})
	}
}

func mustProtoMarshalJSON(msg proto.Message) []byte {
	anyResolver := codectypes.NewInterfaceRegistry()

	// EmitDefaults is set to false to prevent marshaling of unpopulated fields (memo)
	// OrigName and the anyResovler match the fields the original SDK function would expect
	// in order to minimize changes.

	// OrigName is true since there is no particular reason to use camel case
	// The any resolver is empty, but provided anyways.
	jm := &jsonpb.Marshaler{OrigName: true, EmitDefaults: false, AnyResolver: anyResolver}

	err := codectypes.UnpackInterfaces(msg, codectypes.ProtoJSONPacker{JSONPBMarshaler: jm})
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	if err := jm.Marshal(buf, msg); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func ackWithErr() channeltypes.Acknowledgement {
	return channeltypes.NewErrorAcknowledgement(fmt.Errorf("test"))
}

func packetDataToFxPacketData(packetData transfertypes.FungibleTokenPacketData) types.FungibleTokenPacketData {
	return types.FungibleTokenPacketData{
		Denom:    packetData.Denom,
		Amount:   packetData.Amount,
		Sender:   packetData.Sender,
		Receiver: packetData.Receiver,
		Fee:      "0",
		Memo:     packetData.Memo,
	}
}

func getErrByEvent(ctx sdk.Context) string {
	events := ctx.EventManager().Events()
	for _, event := range events {
		if event.Type != types.EventTypeReceive {
			continue
		}
		for _, attr := range event.Attributes {
			if attr.Key == types.AttributeKeyError {
				return attr.Value
			}
		}
	}
	return ""
}
