package keeper_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	tenderminttypes "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/ibc/middleware/types"
)

func (suite *KeeperTestSuite) TestOnRecvPacket() {
	baseDenom := "stake"
	senderAddr := sdk.AccAddress(tmrand.Bytes(20))
	receiveAddr := sdk.AccAddress(tmrand.Bytes(20))
	transferAmount := sdkmath.NewInt(tmrand.Int63n(100000000000))
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
		afterFn       func(packetData transfertypes.FungibleTokenPacketData)
	}{
		{
			name: "pass - normal - ibc transfer packet",
			malleate: func(packet *channeltypes.Packet) {
			},
			expPass:       true,
			errorStr:      "",
			checkBalance:  true,
			checkCoinAddr: receiveAddr,
			expCoins:      sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), transferAmount)),
		},
		{
			name: "pass - normal - receive address is 0xAddress, coin is DefaultCoin",
			malleate: func(packet *channeltypes.Packet) {
				protID := "transfer"
				channelID := "channel-0"
				coin := sdk.NewCoin(fxtypes.DefaultDenom, transferAmount)
				coins := sdk.NewCoins(coin)
				err := suite.bankKeeper.MintCoins(suite.Ctx, transfertypes.ModuleName, coins)
				suite.Require().NoError(err)
				portChannelAddr := transfertypes.GetEscrowAddress(protID, channelID)
				err = suite.bankKeeper.SendCoinsFromModuleToAccount(suite.Ctx, transfertypes.ModuleName, portChannelAddr, coins)
				suite.Require().NoError(err)

				packetData := transfertypes.FungibleTokenPacketData{}
				transfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Receiver = common.BytesToAddress(receiveAddr.Bytes()).String()
				packetData.Denom = transfertypes.DenomTrace{
					BaseDenom: fxtypes.DefaultDenom,
					Path:      fmt.Sprintf("%s/%s", protID, channelID),
				}.GetFullDenomPath()
				packet.Data = packetData.GetBytes()

				suite.ibcTransferKeeper.SetTotalEscrowForDenom(suite.Ctx, coin)
			},
			expPass:       true,
			checkBalance:  true,
			checkCoinAddr: receiveAddr,
			expCoins:      sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, transferAmount)),
		},
		{
			name: "pass - normal - receive address is 0xAddress",
			malleate: func(packet *channeltypes.Packet) {
				packetData := transfertypes.FungibleTokenPacketData{}
				transfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Receiver = common.BytesToAddress(receiveAddr.Bytes()).String()
				packet.Data = packetData.GetBytes()
				meta := banktypes.Metadata{
					// token -> contracts
					Base: ibcDenomTrace.IBCDenom(),
					// evm token: name
					Name: ibcDenomTrace.GetFullDenomPath(),
					// evm token: symbol
					Symbol: strings.ToUpper(baseDenom),
					// evm token decimal - denomunits.denom == symbol
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    strings.ToLower(baseDenom),
							Exponent: 0,
						},
						{
							Denom:    strings.ToUpper(baseDenom),
							Exponent: uint32(tmrand.Int31n(19)),
						},
					},
				}
				_, err := suite.erc20Keeper.RegisterNativeCoin(suite.Ctx, meta)
				suite.Require().NoError(err)
			},
			expPass:       true,
			checkBalance:  true,
			checkCoinAddr: receiveAddr,
			expCoins:      sdk.NewCoins(),
			afterFn: func(packetData transfertypes.FungibleTokenPacketData) {
				expectBalance, ok := sdkmath.NewIntFromString(packetData.Amount)
				suite.Require().True(ok)
				erc20TokenAddr, found := suite.erc20Keeper.GetTokenPair(suite.Ctx, ibcDenomTrace.IBCDenom())
				suite.Require().True(found)
				toAddress := common.HexToAddress(packetData.Receiver)
				var balanceRes struct{ Value *big.Int }
				err := suite.evmKeeper.QueryContract(suite.Ctx, common.Address{}, common.HexToAddress(erc20TokenAddr.Erc20Address), contract.GetFIP20().ABI, "balanceOf", &balanceRes, toAddress)
				suite.Require().NoError(err)
				suite.Require().EqualValues(expectBalance.String(), sdkmath.NewIntFromBigInt(balanceRes.Value).String())
			},
		},
		{
			name: "error - normal - receive address is 0xAddress but coin not registered",
			malleate: func(packet *channeltypes.Packet) {
				packetData := transfertypes.FungibleTokenPacketData{}
				transfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Receiver = common.BytesToAddress(receiveAddr.Bytes()).String()
				packet.Data = packetData.GetBytes()
			},
			// 4: token pair not found
			errorStr:      "ABCI code: 4: error handling packet: see events for details",
			checkBalance:  true,
			checkCoinAddr: senderAddr,
			expCoins:      sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), sdkmath.ZeroInt())),
		},
		{
			name: "pass - any memo",
			malleate: func(packet *channeltypes.Packet) {
				packetData := transfertypes.FungibleTokenPacketData{}
				transfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Memo = "0000"
				packet.Data = packetData.GetBytes()
			},
			expPass:       true,
			errorStr:      "",
			checkBalance:  true,
			checkCoinAddr: senderAddr,
			expCoins:      sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), sdkmath.ZeroInt())),
		},
		{
			name: "pass - ibc call evm",
			malleate: func(packet *channeltypes.Packet) {
				packetData := transfertypes.FungibleTokenPacketData{}
				transfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)

				hexIbcSender := types.IntermediateSender(transfertypes.PortID, "channel-0", senderAddr.String())
				ibcCallBaseAcc := authtypes.NewBaseAccountWithAddress(hexIbcSender.Bytes())
				suite.NoError(ibcCallBaseAcc.SetSequence(0))
				acc := suite.accountKeeper.NewAccount(suite.Ctx, ibcCallBaseAcc)
				suite.accountKeeper.SetAccount(suite.Ctx, acc)
				evmPacket := types.IbcCallEvmPacket{
					To:    common.BigToAddress(big.NewInt(0)).String(),
					Value: sdkmath.ZeroInt(),
					Data:  "",
				}
				bz, err := suite.App.AppCodec().MarshalInterfaceJSON(&evmPacket)
				suite.Require().NoError(err)
				packetData.Memo = string(bz)
				packet.Data = packetData.GetBytes()
			},
			expPass:       true,
			errorStr:      "",
			checkBalance:  true,
			checkCoinAddr: senderAddr,
			expCoins:      sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), sdkmath.ZeroInt())),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.ibcTransferKeeper.SetTotalEscrowForDenom(suite.Ctx, sdk.NewCoin(baseDenom, transferAmount))
			packetData := transfertypes.NewFungibleTokenPacketData(baseDenom, transferAmount.String(), senderAddr.String(), receiveAddr.String(), "")
			// only use timeout height
			packet := channeltypes.NewPacket(packetData.GetBytes(), 1, transfertypes.PortID, "channel-0", transfertypes.PortID, "channel-0", clienttypes.Height{
				RevisionNumber: 100,
				RevisionHeight: 100000,
			}, 0)
			tc.malleate(&packet)

			cacheCtx, writeFn := suite.Ctx.CacheContext()
			cacheCtx = cacheCtx.WithConsensusParams(tenderminttypes.ConsensusParams{Block: &tenderminttypes.BlockParams{MaxGas: 5000000}})
			ackI := suite.ibcMiddleware.OnRecvPacket(cacheCtx, packet, nil)
			if ackI == nil || ackI.Success() {
				// write application state changes for asynchronous and successful acknowledgements
				writeFn()
			}
			suite.Require().NotNil(ackI)

			ack, ok := ackI.(channeltypes.Acknowledgement)
			suite.Ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())

			if tc.expPass {
				suite.Require().Truef(ack.Success(), "ackError:%s,causeError:%s,packetData:%s", ack.GetError(), getOnRecvPacketErrorByEvent(cacheCtx), string(packet.GetData()))
			} else {
				suite.Require().False(ack.Success())
				suite.Require().True(ok)
				suite.Require().Equalf(tc.errorStr, ack.GetError(), "packetData:%s", string(packet.GetData()))
			}

			if tc.checkBalance {
				actualCoins := suite.bankKeeper.GetAllBalances(suite.Ctx, tc.checkCoinAddr)
				suite.Require().True(tc.expCoins.Equal(actualCoins), "exp:%s,actual:%s", tc.expCoins, actualCoins)
			}
		})
	}
}

func getOnRecvPacketErrorByEvent(ctx sdk.Context) string {
	events := ctx.EventManager().Events()
	for _, event := range events {
		if event.Type == transfertypes.EventTypePacket {
			for _, attr := range event.Attributes {
				if attr.Key == types.AttributeKeyRecvError {
					return attr.Value
				}
			}
		}
	}
	return ""
}

func (suite *KeeperTestSuite) TestOnAcknowledgementPacket() {
	baseDenom := "stake"
	senderAddr := sdk.AccAddress(tmrand.Bytes(20))
	receiveAddr := sdk.AccAddress(tmrand.Bytes(20))
	transferAmount := sdkmath.NewInt(tmrand.Int63n(100000000000))
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
				mintCoin(suite.T(), suite.Ctx, suite.bankKeeper, escrowAddress, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)))
			},
			true,
			"",
			true,
			sdk.NewCoins(),
		},
		{
			"pass - error ack - ibc transfer packet",
			func(packet *channeltypes.Packet, ack *channeltypes.Acknowledgement) {
				*ack = channeltypes.NewErrorAcknowledgement(fmt.Errorf("test"))

				escrowAddress := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
				mintCoin(suite.T(), suite.Ctx, suite.bankKeeper, escrowAddress, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)))
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.ibcTransferKeeper.SetTotalEscrowForDenom(suite.Ctx, sdk.NewCoin(baseDenom, transferAmount))
			packetData := transfertypes.NewFungibleTokenPacketData(baseDenom, transferAmount.String(), senderAddr.String(), receiveAddr.String(), "")
			// only use timeout height
			packet := channeltypes.NewPacket(packetData.GetBytes(), 1, transfertypes.PortID, "channel-0", transfertypes.PortID, "channel-0", clienttypes.Height{
				RevisionNumber: 100,
				RevisionHeight: 100000,
			}, 0)

			ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
			tc.malleate(&packet, &ack)

			err := suite.ibcMiddleware.OnAcknowledgementPacket(suite.Ctx, packet, ack.Acknowledgement(), nil)
			if tc.expPass {
				suite.Require().NoError(err, "packetData:%s", string(packet.GetData()))
			} else {
				suite.Require().NotNil(err)
				suite.Require().Equalf(tc.errorStr, err.Error(), "packetData:%s", string(packet.GetData()))
			}

			if tc.checkBalance {
				bankKeeper := suite.bankKeeper
				senderAddrCoins := bankKeeper.GetAllBalances(suite.Ctx, senderAddr)
				suite.Require().True(tc.expCoins.Equal(senderAddrCoins), "exp:%s,actual:%s", tc.expCoins, senderAddrCoins)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOnTimeoutPacket() {
	baseDenom := "stake"
	senderAddr := sdk.AccAddress(tmrand.Bytes(20))
	receiveAddr := sdk.AccAddress(tmrand.Bytes(20))
	transferAmount := sdkmath.NewInt(tmrand.Int63n(100000000000))
	ibcDenomTrace := transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: tmrand.Str(6),
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
				mintCoin(suite.T(), suite.Ctx, suite.bankKeeper, escrowAddress, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)))
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount)),
		},
		{
			"pass - normal - ibc mint token - router is empty",
			func(packet *channeltypes.Packet) {
				packetData := transfertypes.FungibleTokenPacketData{}
				transfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packetData.Denom = ibcDenomTrace.GetFullDenomPath()
				packet.Data = packetData.GetBytes()
			},
			true,
			"",
			true,
			sdk.NewCoins(sdk.NewCoin(ibcDenomTrace.IBCDenom(), transferAmount)),
		},
		{
			"error - escrow address insufficient 10coin",
			func(packet *channeltypes.Packet) {
				packetData := transfertypes.FungibleTokenPacketData{}
				transfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), &packetData)
				packet.Data = packetData.GetBytes()

				escrowAddress := transfertypes.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
				mintCoin(suite.T(), suite.Ctx, suite.bankKeeper, escrowAddress, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.Sub(sdkmath.NewInt(10)))))
			},
			false,
			fmt.Sprintf("unable to unescrow tokens, this may be caused by a malicious counterparty module or a bug: please open an issue on counterparty module: spendable balance %d%s is smaller than %d%s: insufficient funds", transferAmount.Sub(sdkmath.NewInt(10)).Uint64(), baseDenom, transferAmount.Uint64(), baseDenom),
			true,
			sdk.NewCoins(),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.ibcTransferKeeper.SetTotalEscrowForDenom(suite.Ctx, sdk.NewCoin(baseDenom, transferAmount))
			packetData := transfertypes.NewFungibleTokenPacketData(baseDenom, transferAmount.String(), senderAddr.String(), receiveAddr.String(), "")
			// only use timeout height
			packet := channeltypes.NewPacket(packetData.GetBytes(), 1, transfertypes.PortID, "channel-0", transfertypes.PortID, "channel-0", clienttypes.Height{
				RevisionNumber: 100,
				RevisionHeight: 100000,
			}, 0)
			tc.malleate(&packet)

			err := suite.ibcMiddleware.OnTimeoutPacket(suite.Ctx, packet, nil)
			if tc.expPass {
				suite.Require().NoError(err, "packetData:%s", string(packet.GetData()))
			} else {
				suite.Require().NotNil(err)
				suite.Require().Equalf(tc.errorStr, err.Error(), "packetData:%s", string(packet.GetData()))
			}

			if tc.checkBalance {
				bankKeeper := suite.bankKeeper
				senderAddrCoins := bankKeeper.GetAllBalances(suite.Ctx, senderAddr)
				suite.Require().True(tc.expCoins.Equal(senderAddrCoins), "exp:%s,actual:%s", tc.expCoins, senderAddrCoins)
			}
		})
	}
}

func mintCoin(t *testing.T, ctx sdk.Context, bankKeeper bankkeeper.Keeper, address sdk.AccAddress, coins sdk.Coins) {
	require.NoError(t, bankKeeper.MintCoins(ctx, transfertypes.ModuleName, coins))
	require.NoError(t, bankKeeper.SendCoinsFromModuleToAccount(ctx, transfertypes.ModuleName, address, coins))
}
