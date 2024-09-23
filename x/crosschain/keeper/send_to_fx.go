package keeper

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-metrics"

	"github.com/functionx/fx-core/v8/contract"
	fxtelemetry "github.com/functionx/fx-core/v8/telemetry"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (k Keeper) SendToFxExecuted(ctx sdk.Context, claim *types.MsgSendToFxClaim) error {
	bridgeToken := k.GetBridgeTokenDenom(ctx, claim.TokenContract)
	if bridgeToken == nil {
		return errorsmod.Wrap(types.ErrInvalid, "bridge token is not exist")
	}

	coin := sdk.NewCoin(bridgeToken.Denom, claim.Amount)
	if !ctx.IsCheckTx() {
		defer func() {
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, "send_to_fx"},
				float32(1),
				[]metrics.Label{
					telemetry.NewLabel("module", k.moduleName),
				},
			)
			fxtelemetry.SetGaugeLabelsWithDenom(
				[]string{types.ModuleName, "send_to_fx_amount"},
				coin.Denom, coin.Amount.BigInt(),
				telemetry.NewLabel("module", k.moduleName),
			)
		}()
	}
	receiveAddr, err := sdk.AccAddressFromBech32(claim.Receiver)
	if err != nil {
		return errorsmod.Wrap(types.ErrInvalid, "receiver address")
	}
	isOriginOrConverted := k.erc20Keeper.IsOriginOrConvertedDenom(ctx, bridgeToken.Denom)
	if !isOriginOrConverted {
		// If it is not fxcore originated, mint the coins (aka vouchers)
		if err = k.bankKeeper.MintCoins(ctx, k.moduleName, sdk.NewCoins(coin)); err != nil {
			return errorsmod.Wrapf(err, "mint vouchers coins")
		}
	}
	if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, receiveAddr, sdk.NewCoins(coin)); err != nil {
		return errorsmod.Wrap(err, "transfer vouchers")
	}

	// convert to base denom
	cacheCtx, commit := ctx.CacheContext()
	targetCoin, err := k.erc20Keeper.ConvertDenomToTarget(cacheCtx, receiveAddr, coin, fxtypes.ParseFxTarget(fxtypes.ERC20Target))
	if err != nil {
		k.Logger(ctx).Info("failed to convert base denom", "error", err)
		return nil
	}
	commit()

	// relay transfer
	if err = k.RelayTransferHandler(ctx, claim.EventNonce, claim.TargetIbc, receiveAddr, targetCoin); err != nil {
		k.Logger(ctx).Info("failed to relay transfer", "error", err)
		return nil
	}

	k.HandlePendingOutgoingTx(ctx, receiveAddr, claim.GetEventNonce(), bridgeToken)
	k.HandlePendingOutgoingBridgeCall(ctx, receiveAddr, claim.GetEventNonce(), bridgeToken)
	return nil
}

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
		receiverHex := common.BytesToAddress(receiver.Bytes())
		if err := k.erc20Keeper.TransferAfter(cacheCtx, receiver, receiverHex.String(), coin, sdk.NewCoin(coin.Denom, sdkmath.ZeroInt()), false, false); err != nil {
			k.Logger(cacheCtx).Error("transfer convert denom failed", "error", err.Error())
			return err
		}
		cacheCtx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeEvmTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
			sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(eventNonce)),
		))
		commit()
	}
	return nil
}

func (k Keeper) transferIBCHandler(ctx sdk.Context, eventNonce uint64, receive sdk.AccAddress, coin sdk.Coin, target fxtypes.FxTarget) error {
	var ibcReceiveAddress string
	if strings.ToLower(target.Prefix) == contract.EthereumAddressPrefix {
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

	response, err := k.ibcTransferKeeper.Transfer(ctx,
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
		sdk.NewAttribute(types.AttributeKeyIbcSendSequence, fmt.Sprint(response.Sequence)),
		sdk.NewAttribute(types.AttributeKeyIbcSourcePort, target.SourcePort),
		sdk.NewAttribute(types.AttributeKeyIbcSourceChannel, target.SourceChannel),
	))
	return err
}
