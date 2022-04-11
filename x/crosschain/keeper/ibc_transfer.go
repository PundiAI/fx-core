package keeper

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	fxtypes "github.com/functionx/fx-core/types"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	ibcclienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"

	"github.com/functionx/fx-core/x/crosschain/types"
	ibctransfertypes "github.com/functionx/fx-core/x/ibc/applications/transfer/types"
)

var targetEvmPrefix = hex.EncodeToString([]byte("module/evm"))

func (a Keeper) handlerRelayTransfer(ctx sdk.Context, claim *types.MsgSendToFxClaim, receiver sdk.AccAddress, coin sdk.Coin) {
	if ctx.BlockHeight() >= fxtypes.EvmSupportBlock() && claim.TargetIbc == targetEvmPrefix {
		// TODO Do not judge by TargetIbc
		a.handlerEvmTransfer(ctx, claim, receiver, coin)
		return
	}
	a.handleIbcTransfer(ctx, claim, receiver, coin)
}

func (k Keeper) handleIbcTransfer(ctx sdk.Context, claim *types.MsgSendToFxClaim, receiveAddr sdk.AccAddress, coin sdk.Coin) {
	logger := k.Logger(ctx)
	ibcPrefix, sourcePort, sourceChannel, ok := covertIbcData(claim.TargetIbc)
	if !ok {
		logger.Error("convert target ibc data error!!!", "targetIbc", claim.GetTargetIbc())
		return
	}
	ibcReceiveAddress, err := bech32.ConvertAndEncode(ibcPrefix, receiveAddr)
	if err != nil {
		logger.Error("convert ibc transfer receive address error!!!", "fxReceive:", claim.Receiver,
			"ibcPrefix:", ibcPrefix, "sourcePort:", sourcePort, "sourceChannel:", sourceChannel, "error:", err)
		return
	}
	wrapSdkContext := sdk.WrapSDKContext(ctx)

	_, clientState, err := k.ibcChannelKeeper.GetChannelClientState(ctx, sourcePort, sourceChannel)
	if err != nil {
		logger.Error("get channel client state error!!!", "sourcePort", sourcePort, "sourceChannel", sourceChannel)
		return
	}

	params := k.GetParams(ctx)
	clientStateHeight := clientState.GetLatestHeight()
	ibcTimeoutHeight := ibcclienttypes.Height{
		RevisionNumber: clientStateHeight.GetRevisionNumber(),
		RevisionHeight: clientStateHeight.GetRevisionHeight() + params.IbcTransferTimeoutHeight,
	}

	nextSequenceSend, found := k.ibcChannelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		logger.Error("ibc channel next sequence send not found!!!", "source port:", sourcePort, "source channel:", sourceChannel)
		return
	}
	logger.Info("gravity start ibc transfer", "sender:", receiveAddr, "receive:", ibcReceiveAddress, "coin:", coin, "timeout:", params.IbcTransferTimeoutHeight, "nextSequenceSend:", nextSequenceSend)

	ibcTransferMsg := ibctransfertypes.NewMsgTransfer(sourcePort, sourceChannel, coin, receiveAddr, ibcReceiveAddress, ibcTimeoutHeight, 0, "", sdk.NewCoin(coin.Denom, sdk.ZeroInt()))
	if _, err = k.ibcTransferKeeper.Transfer(wrapSdkContext, ibcTransferMsg); err != nil {
		logger.Error("gravity ibc transfer fail. ", "sender:", receiveAddr, "receive:", ibcReceiveAddress, "coin:", coin, "err:", err)
		return
	}
	k.SetIbcSequenceHeight(ctx, sourcePort, sourceChannel, nextSequenceSend, uint64(ctx.BlockHeight()))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeIbcTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(claim.EventNonce)),
		sdk.NewAttribute(types.AttributeKeyIbcSendSequence, fmt.Sprint(nextSequenceSend)),
		sdk.NewAttribute(types.AttributeKeyIbcSourcePort, sourcePort),
		sdk.NewAttribute(types.AttributeKeyIbcSourceChannel, sourceChannel),
	))
}

func (k Keeper) handlerEvmTransfer(ctx sdk.Context, claim *types.MsgSendToFxClaim, receiver sdk.AccAddress, coin sdk.Coin) {
	logger := k.Logger(ctx)
	if !k.erc20Keeper.IsDenomRegistered(ctx, coin.Denom) {
		logger.Error("evm transfer, denom not registered", "denom", coin.Denom)
		return
	}
	receiverEthType := common.BytesToAddress(receiver.Bytes())
	logger.Info("convert denom to fip20", "sender", claim.Sender, "receiver", claim.Receiver,
		"receiver-eth-type", receiverEthType.String(), "amount", coin.String(), "target", claim.TargetIbc)
	err := k.erc20Keeper.ConvertDenomToFIP20(ctx, receiver, receiverEthType, coin)
	if err != nil {
		logger.Error("evm transfer, convert denom to fip20 failed", "error", err.Error())
		return
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEvmTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(claim.EventNonce)),
	))
}

func covertIbcData(targetIbc string) (prefix, sourcePort, sourceChannel string, isOk bool) {
	// fx/transfer/channel-0
	targetIbcBytes, err := hex.DecodeString(targetIbc)
	if err != nil {
		return
	}
	ibcData := strings.Split(string(targetIbcBytes), "/")
	if len(ibcData) < 3 {
		isOk = false
		return
	}
	prefix = ibcData[0]
	sourcePort = ibcData[1]
	sourceChannel = ibcData[2]
	isOk = true
	return
}
