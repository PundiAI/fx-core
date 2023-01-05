package keeper

import (
	"encoding/hex"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

func (k Keeper) RelayTransferHandler(ctx sdk.Context, eventNonce uint64, targetHex string, receiver sdk.AccAddress, coin sdk.Coin) error {
	// ignore hex decode error
	targetByte, _ := hex.DecodeString(targetHex)
	target := string(targetByte)
	targetCoin, isToERC20, err := k.erc20Keeper.ConvertDenomToTarget(ctx, receiver, coin, target)
	if err != nil {
		return err
	}
	if strings.HasPrefix(target, ibchost.ModuleName) {
		targetIBC, ok := fxtypes.ParseTargetIBC(target)
		if !ok {
			return nil
		}
		return k.transferIBCHandler(ctx, eventNonce, receiver, targetCoin, targetIBC)
	}
	if isToERC20 {
		return k.transferErc20Handler(ctx, eventNonce, receiver, targetCoin)
	}
	return nil
}

func (k Keeper) transferErc20Handler(ctx sdk.Context, eventNonce uint64, receiver sdk.AccAddress, coin sdk.Coin) error {
	receiverEthAddr := common.BytesToAddress(receiver.Bytes())
	if err := k.erc20Keeper.TransferAfter(ctx, receiver.String(), receiverEthAddr.String(), coin, sdk.NewCoin(coin.Denom, sdk.ZeroInt())); err != nil {
		k.Logger(ctx).Error("transfer convert denom failed", "error", err.Error())
		return err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEvmTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(eventNonce)),
	))
	return nil
}

func (k Keeper) transferIBCHandler(ctx sdk.Context, eventNonce uint64, receive sdk.AccAddress, coin sdk.Coin, target fxtypes.TargetIBC) error {
	var ibcReceiveAddress string
	if strings.ToLower(target.Prefix) == fxtypes.EthereumAddressPrefix {
		ibcReceiveAddress = common.BytesToAddress(receive.Bytes()).String()
	} else {
		var err error
		ibcReceiveAddress, err = bech32.ConvertAndEncode(target.Prefix, receive)
		if err != nil {
			return err
		}
	}

	_, clientState, err := k.ibcChannelKeeper.GetChannelClientState(ctx, target.SourcePort, target.SourceChannel)
	if err != nil {
		return err
	}

	ibcTransferTimeoutHeight := k.GetIbcTransferTimeoutHeight(ctx)
	clientStateHeight := clientState.GetLatestHeight()
	destTimeoutHeight := clientStateHeight.GetRevisionHeight() + ibcTransferTimeoutHeight
	ibcTimeoutHeight := ibcclienttypes.Height{
		RevisionNumber: clientStateHeight.GetRevisionNumber(),
		RevisionHeight: destTimeoutHeight,
	}

	k.Logger(ctx).Info("crosschain start ibc transfer", "sender", receive, "receive", ibcReceiveAddress,
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
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeIbcTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(eventNonce)),
		sdk.NewAttribute(types.AttributeKeyIbcSendSequence, fmt.Sprint(transferResponse.GetSequence())),
		sdk.NewAttribute(types.AttributeKeyIbcSourcePort, target.SourcePort),
		sdk.NewAttribute(types.AttributeKeyIbcSourceChannel, target.SourceChannel),
	))
	return err
}
