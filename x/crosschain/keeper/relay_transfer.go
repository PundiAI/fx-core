package keeper

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"
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
		if err := k.transferErc20Handler(cacheCtx, eventNonce, receiver, receiver, coin); err != nil {
			return err
		}
		commit()
	}
	return nil
}

func (k Keeper) transferErc20Handler(ctx sdk.Context, eventNonce uint64, sender, receiver sdk.AccAddress, coin sdk.Coin) error {
	receiverEthAddr := common.BytesToAddress(receiver.Bytes())
	if err := k.erc20Keeper.TransferAfter(ctx, sender, receiverEthAddr.String(), coin, sdk.NewCoin(coin.Denom, sdkmath.ZeroInt()), false); err != nil {
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

func (k Keeper) bridgeCallERC20Handler(
	ctx sdk.Context,
	asset, sender, to, receiver []byte,
	dstChainID, message string,
	value sdkmath.Int,
	gasLimit, eventNonce uint64,
) error {
	tokens, amounts, err := types.UnpackERC20Asset(asset)
	if err != nil {
		return errorsmod.Wrap(types.ErrInvalid, "asset erc20")
	}
	senderAddr := common.BytesToAddress(sender)
	targetCoins, err := k.bridgeCallTargetCoinsHandler(ctx, senderAddr, tokens, amounts)
	if err != nil {
		return err
	}

	switch dstChainID {
	case types.FxcoreChainID:
		// convert coin to erc20
		for _, coin := range targetCoins {
			// not convert FX
			if coin.Denom == fxtypes.DefaultDenom {
				continue
			}
			if err = k.transferErc20Handler(ctx, eventNonce, senderAddr.Bytes(), receiver, coin); err != nil {
				return err
			}
		}
		var toAddrPtr *common.Address
		if len(to) > 0 {
			toAddr := common.BytesToAddress(to)
			toAddrPtr = &toAddr
		}
		if len(message) > 0 || toAddrPtr != nil {
			k.bridgeCallEvmHandler(ctx, senderAddr, toAddrPtr, message, value, gasLimit, eventNonce)
		}
	default:
		// not support chain, refund
	}
	// todo refund asset

	return nil
}

func (k Keeper) bridgeCallTargetCoinsHandler(ctx sdk.Context, receiver common.Address, tokens []common.Address, amounts []*big.Int) (sdk.Coins, error) {
	tokens, amounts = types.MergeDuplicationERC20(tokens, amounts)
	targetCoins := sdk.NewCoins()
	for i := 0; i < len(tokens); i++ {
		bridgeToken := k.GetBridgeTokenDenom(ctx, tokens[i].String())
		if bridgeToken == nil {
			return nil, errorsmod.Wrap(types.ErrInvalid, "bridge token is not exist")
		}
		if amounts[i].Cmp(big.NewInt(0)) <= 0 {
			continue
		}
		coin := sdk.NewCoin(bridgeToken.Denom, sdkmath.NewIntFromBigInt(amounts[i]))
		isOriginOrConverted := k.erc20Keeper.IsOriginOrConvertedDenom(ctx, bridgeToken.Denom)
		if !isOriginOrConverted {
			if err := k.bankKeeper.MintCoins(ctx, k.moduleName, sdk.NewCoins(coin)); err != nil {
				return nil, errorsmod.Wrapf(err, "mint vouchers coins")
			}
		}
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, receiver.Bytes(), sdk.NewCoins(coin)); err != nil {
			return nil, errorsmod.Wrap(err, "transfer vouchers")
		}
		targetCoin, err := k.erc20Keeper.ConvertDenomToTarget(ctx, receiver.Bytes(), coin, fxtypes.ParseFxTarget(fxtypes.ERC20Target))
		if err != nil {
			return nil, errorsmod.Wrap(err, "convert to target coin")
		}
		targetCoins = targetCoins.Add(targetCoin)
	}
	return targetCoins, nil
}

func (k Keeper) bridgeCallEvmHandler(
	ctx sdk.Context,
	sender common.Address,
	to *common.Address,
	message string, value sdkmath.Int,
	gasLimit, eventNonce uint64,
) {
	callErr, callResult := "", false
	defer func() {
		attrs := []sdk.Attribute{sdk.NewAttribute(types.AttributeKeyEvmCallResult, strconv.FormatBool(callResult))}
		if len(callErr) > 0 {
			attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyEvmCallError, callErr))
		}
		ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(types.EventTypeBridgeCallEvm, attrs...)})
	}()

	txResp, err := k.evmKeeper.CallEVM(ctx, sender, to, value.BigInt(), gasLimit, types.MustDecodeMessage(message), true)
	if err != nil {
		k.Logger(ctx).Error("bridge call evm error", "nonce", eventNonce, "error", err.Error())
		callErr = err.Error()
		return
	}

	callResult = !txResp.Failed()
	callErr = txResp.VmError
}
