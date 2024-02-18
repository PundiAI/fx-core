package keeper

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) RelayTransferHandler(ctx sdk.Context, eventNonce uint64, targetHex string, receiver sdk.AccAddress, coin sdk.Coin) error {
	// ignore hex decode error
	targetByte, _ := hex.DecodeString(targetHex)
	fxTarget := fxtypes.ParseFxTarget(string(targetByte))

	if fxTarget.IsIBC() {
		// transfer to ibc
		cacheCtx, commit := ctx.CacheContext()
		targetIBCCoin, err1 := k.erc20Keeper.ConvertDenomToTarget(cacheCtx, receiver, coin, fxTarget)
		var err2 error
		if err1 == nil {
			if err2 = k.transferIBCHandler(cacheCtx, eventNonce, receiver, targetIBCCoin, fxTarget); err2 == nil {
				commit()
				return nil
			}
		}
		k.Logger(ctx).Info("failed to transfer ibc", "err1", err1, "err2", err2)
	}

	if fxTarget.GetTarget() == fxtypes.ERC20Target {
		// transfer to evm
		cacheCtx, commit := ctx.CacheContext()
		if err := k.transferErc20Handler(cacheCtx, eventNonce, receiver, coin); err != nil {
			return err
		}
		commit()
	}
	return nil
}

func (k Keeper) transferErc20Handler(ctx sdk.Context, eventNonce uint64, receiver sdk.AccAddress, coin sdk.Coin) error {
	receiverEthAddr := common.BytesToAddress(receiver.Bytes())
	if err := k.erc20Keeper.TransferAfter(ctx, receiver, receiverEthAddr.String(), coin, sdk.NewCoin(coin.Denom, sdkmath.ZeroInt()), false); err != nil {
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

func (k Keeper) transferIBCHandler(ctx sdk.Context, eventNonce uint64, receive sdk.AccAddress, coin sdk.Coin, target fxtypes.FxTarget) error {
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

	// Note: Height is fixed for 5 seconds
	ibcTransferTimeoutHeight := k.GetIbcTransferTimeoutHeight(ctx) * 5
	ibcTimeoutTime := ctx.BlockTime().Add(time.Second * time.Duration(ibcTransferTimeoutHeight))

	response, err := k.ibcTransferKeeper.Transfer(sdk.WrapSDKContext(ctx),
		transfertypes.NewMsgTransfer(
			target.SourcePort,
			target.SourceChannel,
			coin,
			receive.String(),
			ibcReceiveAddress,
			ibcclienttypes.ZeroHeight(),
			uint64(ibcTimeoutTime.UnixNano()),
			"",
		),
	)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeIbcTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(eventNonce)),
		sdk.NewAttribute(types.AttributeKeyIbcSendSequence, fmt.Sprint(response.GetSequence())),
		sdk.NewAttribute(types.AttributeKeyIbcSourcePort, target.SourcePort),
		sdk.NewAttribute(types.AttributeKeyIbcSourceChannel, target.SourceChannel),
	))
	return err
}
