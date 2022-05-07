package transfer_test

import (
	"fmt"
	"testing"
	"time"

	fxtypes "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/ibc/applications/transfer"
	ibctesting "github.com/functionx/fx-core/x/ibc/testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"

	"github.com/functionx/fx-core/x/ibc/applications/transfer/types"
)

var (
	zeroFee  = sdk.NewCoin(fxtypes.DefaultDenom, sdk.ZeroInt())
	noRouter = ""
	noFeeStr = "0"
)

type TransferTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	chainC *ibctesting.TestChain
}

func (suite *TransferTestSuite) SetupTest() {
	fxtypes.ChangeNetworkForTest(fxtypes.NetworkDevnet())
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 3)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainC = suite.coordinator.GetChain(ibctesting.GetChainID(2))
}

// constructs a send from chainA to chainB on the established channel/connection
// and sends the same coin back from chainB to chainA.
func (suite *TransferTestSuite) TestHandleMsgTransfer() {
	// setup between chainA and chainB
	clientA, clientB, connA, connB := suite.coordinator.SetupClientConnections(suite.chainA, suite.chainB, exported.Tendermint)
	channelA, channelB := suite.coordinator.CreateTransferChannels(suite.chainA, suite.chainB, connA, connB, channeltypes.UNORDERED)
	//	originalBalance := suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), sdk.DefaultBondDenom)
	timeoutHeight := clienttypes.NewHeight(0, 110)

	coinToSendToB := sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(100))

	// send from chainA to chainB
	msg := types.NewMsgTransfer(channelA.PortID, channelA.ID, coinToSendToB, suite.chainA.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0, noRouter, zeroFee)

	err := suite.coordinator.SendMsg(suite.chainA, suite.chainB, clientB, msg)
	suite.Require().NoError(err) // message committed

	// relay send
	fungibleTokenPacket := types.NewFungibleTokenPacketData(coinToSendToB.Denom, coinToSendToB.Amount.String(), suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String(), noRouter, noFeeStr)
	packet := channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelA.PortID, channelA.ID, channelB.PortID, channelB.ID, timeoutHeight, 0)
	ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	err = suite.coordinator.RelayPacket(suite.chainA, suite.chainB, clientA, clientB, packet, ack.GetBytes())
	suite.Require().NoError(err) // relay committed

	// check that voucher exists on chain B
	voucherDenomTrace := types.ParseDenomTrace(types.GetPrefixedDenom(packet.GetDestPort(), packet.GetDestChannel(), fxtypes.DefaultDenom))
	balance := suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())

	coinSentFromAToB := types.GetTransferCoin(channelB.PortID, channelB.ID, fxtypes.DefaultDenom, 100)
	suite.Require().Equal(coinSentFromAToB, balance)

	// setup between chainB to chainC
	clientOnBForC, clientOnCForB, connOnBForC, connOnCForB := suite.coordinator.SetupClientConnections(suite.chainB, suite.chainC, exported.Tendermint)
	channelOnBForC, channelOnCForB := suite.coordinator.CreateTransferChannels(suite.chainB, suite.chainC, connOnBForC, connOnCForB, channeltypes.UNORDERED)

	// send from chainB to chainC
	msg = types.NewMsgTransfer(channelOnBForC.PortID, channelOnBForC.ID, coinSentFromAToB, suite.chainB.SenderAccount.GetAddress(), suite.chainC.SenderAccount.GetAddress().String(), timeoutHeight, 0, noRouter, sdk.NewCoin(coinSentFromAToB.Denom, sdk.ZeroInt()))

	err = suite.coordinator.SendMsg(suite.chainB, suite.chainC, clientOnCForB, msg)
	suite.Require().NoError(err) // message committed

	// relay send
	// NOTE: fungible token is prefixed with the full trace in order to verify the packet commitment
	fullDenomPath := types.GetPrefixedDenom(channelOnCForB.PortID, channelOnCForB.ID, voucherDenomTrace.GetFullDenomPath())
	fungibleTokenPacket = types.NewFungibleTokenPacketData(voucherDenomTrace.GetFullDenomPath(), coinSentFromAToB.Amount.String(), suite.chainB.SenderAccount.GetAddress().String(), suite.chainC.SenderAccount.GetAddress().String(), noRouter, noFeeStr)
	packet = channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelOnBForC.PortID, channelOnBForC.ID, channelOnCForB.PortID, channelOnCForB.ID, timeoutHeight, 0)
	err = suite.coordinator.RelayPacket(suite.chainB, suite.chainC, clientOnBForC, clientOnCForB, packet, ack.GetBytes())
	suite.Require().NoError(err) // relay committed

	coinSentFromBToC := sdk.NewInt64Coin(types.ParseDenomTrace(fullDenomPath).IBCDenom(), 100)
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), coinSentFromBToC.Denom)

	// check that the balance is updated on chainC
	suite.Require().Equal(coinSentFromBToC, balance)

	// check that balance on chain B is empty
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
	suite.Require().Zero(balance.Amount.Int64())

	// send from chainC back to chainB
	msg = types.NewMsgTransfer(channelOnCForB.PortID, channelOnCForB.ID, coinSentFromBToC, suite.chainC.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0, noRouter, sdk.NewCoin(coinSentFromBToC.Denom, sdk.ZeroInt()))

	err = suite.coordinator.SendMsg(suite.chainC, suite.chainB, clientOnBForC, msg)
	suite.Require().NoError(err) // message committed

	// relay send
	// NOTE: fungible token is prefixed with the full trace in order to verify the packet commitment
	fungibleTokenPacket = types.NewFungibleTokenPacketData(fullDenomPath, coinSentFromBToC.Amount.String(), suite.chainC.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String(), noRouter, noFeeStr)
	packet = channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelOnCForB.PortID, channelOnCForB.ID, channelOnBForC.PortID, channelOnBForC.ID, timeoutHeight, 0)
	err = suite.coordinator.RelayPacket(suite.chainC, suite.chainB, clientOnCForB, clientOnBForC, packet, ack.GetBytes())
	suite.Require().NoError(err) // relay committed

	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromAToB.Denom)

	// check that the balance on chainA returned back to the original state
	suite.Require().Equal(coinSentFromAToB, balance)

	// check that module account escrow address is empty
	escrowAddress := types.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
	balance = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), escrowAddress, fxtypes.DefaultDenom)
	suite.Require().Equal(sdk.NewCoin(fxtypes.DefaultDenom, sdk.ZeroInt()), balance)

	// check that balance on chain B is empty
	balance = suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	suite.Require().Zero(balance.Amount.Int64())
}

func (suite *TransferTestSuite) TestTransferForwardPacket() {
	var err error
	// setup between chainA and chainB
	clientA, clientB, connA, connB := suite.coordinator.SetupClientConnections(suite.chainA, suite.chainB, exported.Tendermint)
	channelA, channelB := suite.coordinator.CreateTransferChannels(suite.chainA, suite.chainB, connA, connB, channeltypes.UNORDERED)

	if suite.chainB.CurrentHeader.GetHeight() < fxtypes.EvmV1SupportBlock() {
		suite.coordinator.CommitNBlocks(suite.chainB, uint64(fxtypes.EvmV1SupportBlock()-suite.chainB.CurrentHeader.GetHeight())+1)
	}
	// setup between chainB to chainC
	clientOnBForC, clientOnCForB, connOnBForC, connOnCForB := suite.coordinator.SetupClientConnections(suite.chainB, suite.chainC, exported.Tendermint)
	channelOnBForC, channelOnCForB := suite.coordinator.CreateTransferChannels(suite.chainB, suite.chainC, connOnBForC, connOnCForB, channeltypes.UNORDERED)

	// test forward-packet-to-chain
	// send from chainA -> chainB --(forward)--> chainC
	timeoutHeight := clienttypes.NewHeight(0, 1010)
	coinToSendToC := sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(100))

	forwardPacketReceiverAddress := fmt.Sprintf("%s|%s/%s:%s", suite.chainB.SenderAccount.GetAddress().String(), channelOnBForC.PortID, channelOnBForC.ID, suite.chainC.SenderAccount.GetAddress().String())
	chainAToChainBForwardChainCMsg := types.NewMsgTransfer(channelA.PortID, channelA.ID, coinToSendToC, suite.chainA.SenderAccount.GetAddress(), forwardPacketReceiverAddress, timeoutHeight, 0, noRouter, zeroFee)

	err = suite.coordinator.SendMsg(suite.chainA, suite.chainB, clientB, chainAToChainBForwardChainCMsg)
	suite.Require().NoError(err) // message committed

	// relay chainA to chainB
	fungibleTokenPacket := types.NewFungibleTokenPacketData(coinToSendToC.Denom, coinToSendToC.Amount.String(), suite.chainA.SenderAccount.GetAddress().String(), forwardPacketReceiverAddress, noRouter, noFeeStr)
	packet := channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelA.PortID, channelA.ID, channelB.PortID, channelB.ID, timeoutHeight, 0)
	ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	err = suite.coordinator.RelayPacket(suite.chainA, suite.chainB, clientA, clientB, packet, ack.GetBytes())
	suite.Require().NoError(err)

	// update chainC clientOnCForB state
	err = suite.coordinator.UpdateClient(suite.chainC, suite.chainB, clientOnCForB, exported.Tendermint)
	suite.Require().NoError(err)

	// relay chainB to chainC
	voucherDenomTrace := types.ParseDenomTrace(types.GetPrefixedDenom(packet.GetDestPort(), packet.GetDestChannel(), fxtypes.DefaultDenom))

	fungibleTokenPacket = types.NewFungibleTokenPacketData(voucherDenomTrace.GetFullDenomPath(), coinToSendToC.Amount.String(), suite.chainB.SenderAccount.GetAddress().String(), suite.chainC.SenderAccount.GetAddress().String(), noRouter, noFeeStr)
	chainCTimeoutHeight := clienttypes.NewHeight(0, 0)
	chainCTimeoutTimestamp := uint64(suite.chainB.LastRecvPacketHeader.Time.Add(transfer.ForwardPacketTimeHour * time.Hour).UnixNano())
	packet = channeltypes.NewPacket(fungibleTokenPacket.GetBytes(), 1, channelOnBForC.PortID, channelOnBForC.ID, channelOnCForB.PortID, channelOnCForB.ID, chainCTimeoutHeight, chainCTimeoutTimestamp)
	ack = channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	err = suite.coordinator.RelayPacket(suite.chainB, suite.chainC, clientOnBForC, clientOnCForB, packet, ack.GetBytes())
	suite.Require().NoError(err) // relay committed

	// check chainC receive address balance
	voucherDenomTrace = types.ParseDenomTrace(types.GetPrefixedDenom(packet.GetDestPort(), packet.GetDestChannel(), voucherDenomTrace.GetFullDenomPath()))
	balance := suite.chainC.App.BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	chainCReceiveAddrBalance := sdk.NewCoin(voucherDenomTrace.IBCDenom(), coinToSendToC.Amount)
	suite.Require().True(chainCReceiveAddrBalance.Equal(balance))
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
