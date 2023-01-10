package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcgotesting "github.com/cosmos/ibc-go/v3/testing"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v3/app"
	avalanchetypes "github.com/functionx/fx-core/v3/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	fxtransfer "github.com/functionx/fx-core/v3/x/ibc/applications/transfer"
	fxtransfertypes "github.com/functionx/fx-core/v3/x/ibc/applications/transfer/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func (suite *KeeperTestSuite) TestOnRecvPacket() {
	baseDenom := "stake"
	senderAddr := sdk.AccAddress(rand.Bytes(20))
	receiveAddr := sdk.AccAddress(rand.Bytes(20))
	transferAmount := sdk.NewInt(rand.Int63n(100000000000))
	ibcDenomTrace := transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: baseDenom,
	}

	testCases := []struct {
		name          string
		malleate      func(fxIbcTransferMsg *channeltypes.Packet)
		expPass       bool
		errorStr      string
		checkBalance  bool
		checkCoinAddr sdk.AccAddress
		expCoins      sdk.Coins
	}{
		{
			"pass - normal - ibc transfer packet",
			func(packet *channeltypes.Packet) {
			},
			true,
			"",
			true,
			receiveAddr,
			sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), transferAmount)),
		},
		{
			"pass - normal - fx ibc transfer packet",
			func(packet *channeltypes.Packet) {
				packetData := fxtransfertypes.FungibleTokenPacketData{}
				fxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Router = ""
				packetData.Fee = sdk.ZeroInt().String()
				packet.Data = packetData.GetBytes()
			},
			true,
			"",
			true,
			receiveAddr,
			sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), transferAmount)),
		},
		{
			"pass - normal - router is bsc, sender is 0xAddress",
			func(packet *channeltypes.Packet) {
				bscKeeper := suite.GetApp(suite.chainA.App).BscKeeper
				bscKeeper.AddBridgeToken(suite.chainA.GetContext(), common.BytesToAddress(rand.Bytes(20)).String(), ibcDenomTrace.IBCDenom())
				packetData := fxtransfertypes.FungibleTokenPacketData{}
				fxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Sender = common.BytesToAddress(senderAddr.Bytes()).String()
				packetData.Router = bsctypes.ModuleName
				packetData.Fee = sdk.ZeroInt().String()
				packetData.Receiver = common.BytesToAddress(receiveAddr).String()
				packet.Data = packetData.GetBytes()
			},
			true,
			"",
			true,
			senderAddr,
			sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), sdk.ZeroInt())),
		},
		{
			name: "error - normal - transferAfter return error, receive address is error",
			malleate: func(packet *channeltypes.Packet) {
				packetData := fxtransfertypes.FungibleTokenPacketData{}
				fxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Receiver = rand.Str(20)
				routes := []string{ethtypes.ModuleName, bsctypes.ModuleName, trontypes.ModuleName, polygontypes.ModuleName, avalanchetypes.ModuleName, erc20types.ModuleName}
				packetData.Router = routes[rand.Int63n(int64(len(routes)))]
				packet.Data = packetData.GetBytes()
			},
			expPass:       false,
			errorStr:      "ABCI code: 2: error handling packet on destination chain: see events for details",
			checkBalance:  true,
			checkCoinAddr: senderAddr,
			expCoins:      sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), sdk.ZeroInt())),
		},
		{
			"error - normal - router not exists",
			func(packet *channeltypes.Packet) {
				packetData := fxtransfertypes.FungibleTokenPacketData{}
				fxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Router = rand.Str(8)
				packetData.Fee = sdk.ZeroInt().String()
				packet.Data = packetData.GetBytes()
			},
			false,
			// 103: router not found error
			"ABCI code: 103: error handling packet on destination chain: see events for details",
			true,
			senderAddr,
			sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), sdk.ZeroInt())),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			chain := suite.GetApp(suite.chainA.App)
			transferIBCModule := transfer.NewIBCModule(chain.IBCTransferKeeper)
			fxIBCMiddleware := fxtransfer.NewIBCMiddleware(chain.FxTransferKeeper, transferIBCModule)
			packetData := transfertypes.NewFungibleTokenPacketData(baseDenom, transferAmount.String(), senderAddr.String(), receiveAddr.String())
			// only use timeout height
			packet := channeltypes.NewPacket(packetData.GetBytes(), 1, ibcgotesting.TransferPort, "channel-0", ibcgotesting.TransferPort, "channel-0", clienttypes.Height{
				RevisionNumber: 100,
				RevisionHeight: 100000,
			}, 0)
			tc.malleate(&packet)

			cacheCtx, writeFn := suite.chainA.GetContext().CacheContext()
			ackI := fxIBCMiddleware.OnRecvPacket(cacheCtx, packet, nil)
			if ackI == nil || ackI.Success() {
				// write application state changes for asynchronous and successful acknowledgements
				writeFn()
			}
			suite.Require().NotNil(ackI)

			ack, ok := ackI.(channeltypes.Acknowledgement)
			suite.chainA.GetContext().EventManager().EmitEvents(cacheCtx.EventManager().Events())

			if tc.expPass {
				suite.Require().Truef(ack.Success(), "error:%s,packetData:%s", ack.GetError(), string(packet.GetData()))
			} else {
				suite.Require().False(ack.Success())
				suite.Require().True(ok)
				suite.Require().Equalf(tc.errorStr, ack.GetError(), "packetData:%s", string(packet.GetData()))
			}

			if tc.checkBalance {
				bankKeeper := suite.GetApp(suite.chainA.App).BankKeeper
				actualCoins := bankKeeper.GetAllBalances(suite.chainA.GetContext(), tc.checkCoinAddr)
				suite.Require().True(tc.expCoins.IsEqual(actualCoins), "exp:%s,actual:%s", tc.expCoins, actualCoins)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOnAcknowledgementPacket() {
	baseDenom := "stake"
	senderAddr := sdk.AccAddress(rand.Bytes(20))
	receiveAddr := sdk.AccAddress(rand.Bytes(20))
	transferAmount := sdk.NewInt(rand.Int63n(100000000000))
	testCases := []struct {
		name         string
		malleate     func(fxIbcTransferMsg *channeltypes.Packet, ack *channeltypes.Acknowledgement)
		expPass      bool
		errorStr     string
		checkBalance bool
		expCoins     sdk.Coins
	}{
		{
			"pass - success ack - ibc transfer packet",
			func(packet *channeltypes.Packet, ack *channeltypes.Acknowledgement) {
				escrowAddress := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
				mintCoin(suite.T(), suite.chainA.GetContext(), suite.GetApp(suite.chainA.App), escrowAddress, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)))
			},
			true,
			"",
			true,
			sdk.NewCoins(),
		},
		{
			"pass - error ack - ibc transfer packet",
			func(packet *channeltypes.Packet, ack *channeltypes.Acknowledgement) {
				*ack = channeltypes.NewErrorAcknowledgement("test")

				escrowAddress := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
				mintCoin(suite.T(), suite.chainA.GetContext(), suite.GetApp(suite.chainA.App), escrowAddress, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)))
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			chain := suite.GetApp(suite.chainA.App)
			transferIBCModule := transfer.NewIBCModule(chain.IBCTransferKeeper)
			fxIBCMiddleware := fxtransfer.NewIBCMiddleware(chain.FxTransferKeeper, transferIBCModule)
			packetData := transfertypes.NewFungibleTokenPacketData(baseDenom, transferAmount.String(), senderAddr.String(), receiveAddr.String())
			// only use timeout height
			packet := channeltypes.NewPacket(packetData.GetBytes(), 1, ibcgotesting.TransferPort, "channel-0", ibcgotesting.TransferPort, "channel-0", clienttypes.Height{
				RevisionNumber: 100,
				RevisionHeight: 100000,
			}, 0)

			ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
			tc.malleate(&packet, &ack)

			err := fxIBCMiddleware.OnAcknowledgementPacket(suite.chainA.GetContext(), packet, ack.Acknowledgement(), nil)
			if tc.expPass {
				suite.Require().NoError(err, "packetData:%s", string(packet.GetData()))
			} else {
				suite.Require().NotNil(err)
				suite.Require().Equalf(tc.errorStr, err.Error(), "packetData:%s", string(packet.GetData()))
			}

			if tc.checkBalance {
				bankKeeper := suite.GetApp(suite.chainA.App).BankKeeper
				senderAddrCoins := bankKeeper.GetAllBalances(suite.chainA.GetContext(), senderAddr)
				suite.Require().True(tc.expCoins.IsEqual(senderAddrCoins), "exp:%s,actual:%s", tc.expCoins, senderAddrCoins)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOnTimeoutPacket() {
	baseDenom := "stake"
	senderAddr := sdk.AccAddress(rand.Bytes(20))
	receiveAddr := sdk.AccAddress(rand.Bytes(20))
	transferAmount := sdk.NewInt(rand.Int63n(100000000000))
	ibcDenomTrace := transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: rand.Str(6),
	}
	testCases := []struct {
		name         string
		malleate     func(fxIbcTransferMsg *channeltypes.Packet)
		expPass      bool
		errorStr     string
		checkBalance bool
		expCoins     sdk.Coins
	}{
		{
			"pass - normal - ibc transfer packet",
			func(packet *channeltypes.Packet) {
				escrowAddress := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
				mintCoin(suite.T(), suite.chainA.GetContext(), suite.GetApp(suite.chainA.App), escrowAddress, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)))
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)),
		},
		{
			"pass - normal - fx ibc transfer packet",
			func(packet *channeltypes.Packet) {
				packetData := fxtransfertypes.FungibleTokenPacketData{}
				fxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Router = ""
				packetData.Fee = sdk.ZeroInt().String()
				packet.Data = packetData.GetBytes()

				escrowAddress := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
				mintCoin(suite.T(), suite.chainA.GetContext(), suite.GetApp(suite.chainA.App), escrowAddress, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)))
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)),
		},
		{
			"pass - normal - ibc mint token - router is empty",
			func(packet *channeltypes.Packet) {
				packetData := fxtransfertypes.FungibleTokenPacketData{}
				fxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Denom = ibcDenomTrace.GetFullDenomPath()
				packetData.Router = ""
				packetData.Fee = sdk.ZeroInt().String()
				packet.Data = packetData.GetBytes()
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), transferAmount)),
		},
		{
			"pass - router not empty | amount + zero fee",
			func(packet *channeltypes.Packet) {
				packetData := fxtransfertypes.FungibleTokenPacketData{}
				fxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Router = rand.Str(4)
				packetData.Fee = sdk.ZeroInt().String()
				packet.Data = packetData.GetBytes()

				amount := sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount))
				escrowAddress := transfertypes.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
				mintCoin(suite.T(), suite.chainA.GetContext(), suite.GetApp(suite.chainA.App), escrowAddress, amount)
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)),
		},
		{
			"pass - router not empty | amount + fee",
			func(packet *channeltypes.Packet) {
				packetData := fxtransfertypes.FungibleTokenPacketData{}
				fxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Router = rand.Str(4)
				fee := sdk.NewInt(50)
				packetData.Fee = fee.String()
				packet.Data = packetData.GetBytes()

				amount := sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.Add(fee)))
				escrowAddress := transfertypes.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
				mintCoin(suite.T(), suite.chainA.GetContext(), suite.GetApp(suite.chainA.App), escrowAddress, amount)
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.Add(sdk.NewInt(50)))),
		},
		{
			"error - escrow address insufficient 10coin",
			func(packet *channeltypes.Packet) {
				packetData := fxtransfertypes.FungibleTokenPacketData{}
				fxtransfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Router = ""
				packetData.Fee = sdk.ZeroInt().String()
				packet.Data = packetData.GetBytes()

				escrowAddress := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
				mintCoin(suite.T(), suite.chainA.GetContext(), suite.GetApp(suite.chainA.App), escrowAddress, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.Sub(sdk.NewInt(10)))))
			},
			false,
			fmt.Sprintf("unable to unescrow tokens, this may be caused by a malicious counterparty module or a bug: please open an issue on counterparty module: %d%s is smaller than %d%s: insufficient funds", transferAmount.Sub(sdk.NewInt(10)).Uint64(), baseDenom, transferAmount.Uint64(), baseDenom),
			true,
			sdk.NewCoins(),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			chain := suite.GetApp(suite.chainA.App)
			transferIBCModule := transfer.NewIBCModule(chain.IBCTransferKeeper)
			fxIBCMiddleware := fxtransfer.NewIBCMiddleware(chain.FxTransferKeeper, transferIBCModule)
			packetData := transfertypes.NewFungibleTokenPacketData(baseDenom, transferAmount.String(), senderAddr.String(), receiveAddr.String())
			// only use timeout height
			packet := channeltypes.NewPacket(packetData.GetBytes(), 1, ibcgotesting.TransferPort, "channel-0", ibcgotesting.TransferPort, "channel-0", clienttypes.Height{
				RevisionNumber: 100,
				RevisionHeight: 100000,
			}, 0)
			tc.malleate(&packet)

			err := fxIBCMiddleware.OnTimeoutPacket(suite.chainA.GetContext(), packet, nil)
			if tc.expPass {
				suite.Require().NoError(err, "packetData:%s", string(packet.GetData()))
			} else {
				suite.Require().NotNil(err)
				suite.Require().Equalf(tc.errorStr, err.Error(), "packetData:%s", string(packet.GetData()))
			}

			if tc.checkBalance {
				bankKeeper := suite.GetApp(suite.chainA.App).BankKeeper
				senderAddrCoins := bankKeeper.GetAllBalances(suite.chainA.GetContext(), senderAddr)
				suite.Require().True(tc.expCoins.IsEqual(senderAddrCoins), "exp:%s,actual:%s", tc.expCoins, senderAddrCoins)
			}
		})
	}
}

func mintCoin(t *testing.T, ctx sdk.Context, chain *app.App, address sdk.AccAddress, coins sdk.Coins) {
	require.NoError(t, chain.BankKeeper.MintCoins(ctx, transfertypes.ModuleName, coins))
	require.NoError(t, chain.BankKeeper.SendCoinsFromModuleToAccount(ctx, transfertypes.ModuleName, address, coins))
}
