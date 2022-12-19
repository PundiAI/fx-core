package keeper

import (
	"encoding/hex"
	"fmt"

	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

var targetEvmPrefix = hex.EncodeToString([]byte("module/evm"))

func (k Keeper) HandlerRelayTransfer(ctx sdk.Context, claim *types.MsgSendToFxClaim, receiver sdk.AccAddress, coin sdk.Coin) {
	// convert denom
	coin = k.handlerConvertDenom(ctx, claim, receiver, coin)

	// router evm condition
	// 1. target == module/evm
	evmTarget := claim.TargetIbc == targetEvmPrefix
	if evmTarget {
		k.handlerEvmTransfer(ctx, claim, receiver, coin)
		return
	}
	k.handleIbcTransfer(ctx, claim, receiver, coin)
}

func (k Keeper) handleIbcTransfer(ctx sdk.Context, claim *types.MsgSendToFxClaim, receiveAddr sdk.AccAddress, coin sdk.Coin) {
	logger := k.Logger(ctx)
	targetIBC, ok := fxtypes.ParseHexTargetIBC(claim.TargetIbc)
	if !ok {
		logger.Error("convert target ibc data error!!!", "targetIbc", claim.GetTargetIbc())
		return
	}
	ibcReceiveAddress, err := types.CovertIbcPacketReceiveAddressByPrefix(targetIBC.Prefix, receiveAddr)
	if err != nil {
		logger.Error("convert ibc transfer receive address error!!!", "fxReceive", claim.Receiver,
			"ibcPrefix", targetIBC.Prefix, "sourcePort", targetIBC.SourcePort, "sourceChannel", targetIBC.SourceChannel, "error", err)
		return
	}

	_, clientState, err := k.ibcChannelKeeper.GetChannelClientState(ctx, targetIBC.SourcePort, targetIBC.SourceChannel)
	if err != nil {
		logger.Error("get channel client state error!!!", "sourcePort", targetIBC.SourcePort, "sourceChannel", targetIBC.SourceChannel)
		return
	}

	ibcTransferTimeoutHeight := k.GetIbcTransferTimeoutHeight(ctx)
	clientStateHeight := clientState.GetLatestHeight()
	destTimeoutHeight := clientStateHeight.GetRevisionHeight() + ibcTransferTimeoutHeight
	ibcTimeoutHeight := ibcclienttypes.Height{
		RevisionNumber: clientStateHeight.GetRevisionNumber(),
		RevisionHeight: destTimeoutHeight,
	}

	logger.Info("crosschain start ibc transfer", "sender", receiveAddr, "receive", ibcReceiveAddress, "coin", coin, "destCurrentHeight", clientStateHeight.GetRevisionHeight(), "destTimeoutHeight", destTimeoutHeight)

	transferMsg := transfertypes.NewMsgTransfer(targetIBC.SourcePort, targetIBC.SourceChannel, coin, receiveAddr.String(), ibcReceiveAddress, ibcTimeoutHeight, 0)
	transferResponse, err := k.ibcTransferKeeper.Transfer(sdk.WrapSDKContext(ctx), transferMsg)
	if err != nil {
		logger.Error("crosschain ibc transfer fail", "sender", receiveAddr, "receive", ibcReceiveAddress, "coin", coin, "err", err)
		return
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeIbcTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(claim.EventNonce)),
		sdk.NewAttribute(types.AttributeKeyIbcSendSequence, fmt.Sprint(transferResponse.GetSequence())),
		sdk.NewAttribute(types.AttributeKeyIbcSourcePort, targetIBC.SourcePort),
		sdk.NewAttribute(types.AttributeKeyIbcSourceChannel, targetIBC.SourceChannel),
	))
}

func (k Keeper) handlerEvmTransfer(ctx sdk.Context, claim *types.MsgSendToFxClaim, receiver sdk.AccAddress, coin sdk.Coin) {
	logger := k.Logger(ctx)
	receiverEthType := common.BytesToAddress(receiver.Bytes())
	logger.Info("convert denom to fip20", "sender", claim.Sender, "receiver", claim.Receiver,
		"receiver-eth-type", receiverEthType.String(), "amount", coin.String(), "target", claim.TargetIbc)
	err := k.erc20Keeper.RelayConvertCoin(ctx, receiver, receiverEthType, coin)
	if err != nil {
		logger.Error("evm transfer, convert denom to fip20 failed", "error", err.Error())
		return
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEvmTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(claim.EventNonce)),
	))
}

func (k Keeper) handlerConvertDenom(ctx sdk.Context, claim *types.MsgSendToFxClaim, receiver sdk.AccAddress, coin sdk.Coin) sdk.Coin {
	logger := k.Logger(ctx)
	logger.Info("convert denom symbol", "address", receiver.String(), "coin", coin.String())

	cacheCtx, commit := ctx.CacheContext()
	targetCoin, err := k.erc20Keeper.RelayConvertDenomToOne(cacheCtx, receiver, coin)
	if err != nil {
		logger.Error("convert denom symbol", "address", receiver.String(), "coin", coin.String(), "error", err.Error())
		//if convert err, return default coin
		return coin
	}
	commit()
	logger.Info("convert denom symbol", "address", receiver.String(), "coin", coin.String(), "target", targetCoin.String(), "module", k.moduleName)
	return targetCoin
}
