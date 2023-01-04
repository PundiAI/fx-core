package keeper

import (
	"encoding/hex"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

var targetEvmPrefix = hex.EncodeToString([]byte("module/evm"))

func (k Keeper) HandlerRelayTransfer(ctx sdk.Context, eventNonce uint64, targetIbc string, receiver sdk.AccAddress, coin sdk.Coin) {
	cacheCtx, commit := ctx.CacheContext()
	targetCoin, err := k.erc20Keeper.ConvertDenomToOne(cacheCtx, receiver, coin)
	if err != nil {
		k.Logger(ctx).Error("convert denom symbol", "address", receiver, "coin", coin, "error", err.Error())
		return
	}
	commit()

	if targetIbc == targetEvmPrefix {
		k.handlerEvmTransfer(ctx, eventNonce, receiver, targetCoin)
		return
	}
	target, ok := fxtypes.ParseHexTargetIBC(targetIbc)
	if !ok {
		return
	}
	k.handleIbcTransfer(ctx, eventNonce, receiver, targetCoin, target)
}

func (k Keeper) handlerEvmTransfer(ctx sdk.Context, eventNonce uint64, receiver sdk.AccAddress, coin sdk.Coin) {
	receiverEthAddr := common.BytesToAddress(receiver.Bytes())
	if err := k.erc20Keeper.TransferAfter(ctx, receiver.String(), receiverEthAddr.String(), coin, sdk.Coin{}); err != nil {
		k.Logger(ctx).Error("transfer convert denom failed", "error", err.Error())
		return
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEvmTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(eventNonce)),
	))
}

func (k Keeper) handleIbcTransfer(ctx sdk.Context, eventNonce uint64, receive sdk.AccAddress, coin sdk.Coin, target fxtypes.TargetIBC) {
	logger := k.Logger(ctx)

	ibcReceiveAddress, err := covertIbcPacketReceiveAddress(target.Prefix, receive)
	if err != nil {
		logger.Error("convert ibc transfer receive address error!!!", "error", err)
		return
	}

	_, clientState, err := k.ibcChannelKeeper.GetChannelClientState(ctx, target.SourcePort, target.SourceChannel)
	if err != nil {
		logger.Error("get channel client state error!!!", "error", err)
		return
	}

	ibcTransferTimeoutHeight := k.GetIbcTransferTimeoutHeight(ctx)
	clientStateHeight := clientState.GetLatestHeight()
	destTimeoutHeight := clientStateHeight.GetRevisionHeight() + ibcTransferTimeoutHeight
	ibcTimeoutHeight := ibcclienttypes.Height{
		RevisionNumber: clientStateHeight.GetRevisionNumber(),
		RevisionHeight: destTimeoutHeight,
	}

	logger.Info("crosschain start ibc transfer", "sender", receive, "receive", ibcReceiveAddress,
		"coin", coin, "destCurrentHeight", clientStateHeight.GetRevisionHeight(), "destTimeoutHeight", destTimeoutHeight)

	transferMsg := transfertypes.NewMsgTransfer(
		target.SourcePort,
		target.SourceChannel,
		coin,
		receive.String(),
		ibcReceiveAddress,
		ibcTimeoutHeight,
		0,
	)
	transferResponse, err := k.ibcTransferKeeper.Transfer(sdk.WrapSDKContext(ctx), transferMsg)
	if err != nil {
		logger.Error("ibc transfer failed", "error", err)
		return
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeIbcTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(eventNonce)),
		sdk.NewAttribute(types.AttributeKeyIbcSendSequence, fmt.Sprint(transferResponse.GetSequence())),
		sdk.NewAttribute(types.AttributeKeyIbcSourcePort, target.SourcePort),
		sdk.NewAttribute(types.AttributeKeyIbcSourceChannel, target.SourceChannel),
	))
}

func covertIbcPacketReceiveAddress(targetIbcPrefix string, receiver sdk.AccAddress) (string, error) {
	if strings.ToLower(targetIbcPrefix) == fxtypes.EthereumAddressPrefix {
		return common.BytesToAddress(receiver.Bytes()).String(), nil
	}
	return bech32.ConvertAndEncode(targetIbcPrefix, receiver)
}
