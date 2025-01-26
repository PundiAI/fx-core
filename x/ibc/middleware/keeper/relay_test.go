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
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	"github.com/pundiai/fx-core/v8/x/ibc/middleware/keeper"
	"github.com/pundiai/fx-core/v8/x/ibc/middleware/types"
)

func (suite *KeeperTestSuite) TestOnAcknowledgementPacket() {
	coin := sdk.NewCoin("stake", sdkmath.NewInt(tmrand.Int63n(100000000000)))

	testCases := []struct {
		name         string
		malleate     func(packet *channeltypes.Packet, ack *channeltypes.Acknowledgement, packetData transfertypes.FungibleTokenPacketData)
		expPass      bool
		errorStr     string
		checkBalance bool
		expCoins     sdk.Coins
	}{
		{
			name:         "pass - success ack - ibc transfer packet",
			expPass:      true,
			checkBalance: true,
			expCoins:     sdk.Coins{},
		},
		{
			name: "pass - error ack - ibc transfer packet",
			malleate: func(packet *channeltypes.Packet, ack *channeltypes.Acknowledgement, packetData transfertypes.FungibleTokenPacketData) {
				*ack = ackWithErr()

				suite.mintCoinEscrowAddr(packet.SourceChannel, coin)
				suite.Require().NoError(suite.App.Erc20Keeper.SetCache(suite.Ctx, crosschaintypes.NewIBCTransferKey(packet.SourceChannel, packet.Sequence), coin.Amount))
			},
			expPass:      true,
			checkBalance: true,
			expCoins:     sdk.NewCoins(coin),
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
			expCoins:     sdk.NewCoins(fxtypes.SwapCoin(sdk.NewCoin(fxtypes.LegacyFXDenom, coin.Amount))),
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
				suite.Require().Equalf(tc.errorStr, err.Error(), "packet: %s", string(packet.GetData()))
			}
			if tc.checkBalance {
				senderAddrCoins := suite.App.BankKeeper.GetAllBalances(suite.Ctx, senderAddr)
				suite.Require().Equal(tc.expCoins.String(), senderAddrCoins.String())
			}
		})
	}
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

// func (suite *KeeperTestSuite) TestOnTimeoutPacket() {
//	coin := sdk.NewCoin("stake", sdkmath.NewInt(tmrand.Int63n(100000000000)))
//	testCases := []struct {
//		name         string
//		malleate     func(packet *channeltypes.Packet)
//		expPass      bool
//		errorStr     string
//		checkBalance bool
//		expCoins     sdk.Coins
//	}{
//		{
//			name: "pass - normal - ibc transfer packet",
//			malleate: func(packet *channeltypes.Packet) {
//			},
//			expPass:      true,
//			checkBalance: true,
//			expCoins:     sdk.Coins{},
//		},
//		{
//			name: "pass - normal - ibc mint token - router is empty",
//			malleate: func(packet *channeltypes.Packet) {
//			},
//			expPass:      true,
//			checkBalance: true,
//			expCoins:     sdk.Coins{},
//		},
//		{
//			name: "error - escrow address insufficient 10coin",
//			malleate: func(packet *channeltypes.Packet) {
//			},
//			errorStr:     fmt.Sprintf("unable to unescrow tokens, this may be caused by a malicious counterparty module or a bug: please open an issue on counterparty module: spendable balance %d%s is smaller than %d%s: insufficient funds", coin.Amount.Sub(sdkmath.NewInt(10)).Uint64(), coin.Denom, coin.Amount.Uint64(), coin.Denom),
//			checkBalance: true,
//			expCoins:     sdk.Coins{},
//		},
//	}
//
//	for _, tc := range testCases {
//		suite.Run(tc.name, func() {
//
//		})
//	}
//}

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
