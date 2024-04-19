package keeper

import (
	"fmt"
	"math/big"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
)

func (k Keeper) BridgeCallHandler(
	ctx sdk.Context,
	sender common.Address,
	to *common.Address,
	receiver sdk.AccAddress,
	tokens []common.Address,
	amounts []*big.Int,
	message []byte,
	value sdkmath.Int,
	gasLimit, eventNonce uint64,
) error {
	erc20Token, err := types.NewERC20Tokens(k.moduleName, tokens, amounts)
	if err != nil {
		return err
	}
	var errCause string
	cacheCtx, commit := ctx.CacheContext()
	if err = k.bridgeCallFxCore(cacheCtx, sender, receiver, to, erc20Token, message, value, gasLimit, eventNonce); err != nil {
		errCause = err.Error()
	} else {
		commit()
	}
	if len(errCause) > 0 && len(tokens) > 0 {
		// new outgoing bridge call to refund
		outCall, err := k.AddOutgoingBridgeCall(ctx, receiver, sender.String(), common.Address{}.String(), erc20Token, "", sdkmath.ZeroInt(), 0)
		if err != nil {
			return err
		}
		// refund event
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeBridgeCallRefund,
			sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprintf("%d", eventNonce)),
			sdk.NewAttribute(types.AttributeKeyBridgeCallNonce, fmt.Sprintf("%d", outCall.Nonce)),
		))
	}

	if len(errCause) == 0 {
		for i := 0; i < len(erc20Token); i++ {
			bridgeToken := k.GetBridgeTokenDenom(ctx, erc20Token[i].Contract)
			// no need for a double check here, as the bridge token should exist
			k.HandlePendingOutgoingTx(ctx, receiver, eventNonce, bridgeToken)
		}
	}
	return nil
}

func (k Keeper) bridgeCallFxCore(
	ctx sdk.Context,
	sender common.Address,
	receiver sdk.AccAddress,
	to *common.Address,
	tokens []types.ERC20Token,
	message []byte,
	value sdkmath.Int,
	gasLimit uint64,
	eventNonce uint64,
) error {
	coins, err := k.bridgeCallTransferToSender(ctx, sender.Bytes(), tokens)
	if err != nil {
		return err
	}
	if err = k.bridgeCallTransferToReceiver(ctx, sender.Bytes(), receiver, coins); err != nil {
		return err
	}
	if len(message) > 0 || to != nil {
		evmErr, evmResult := "", false
		defer func() {
			attrs := []sdk.Attribute{
				sdk.NewAttribute(types.AttributeKeyEventNonce, strconv.FormatUint(eventNonce, 10)),
				sdk.NewAttribute(types.AttributeKeyBridgeCallResult, strconv.FormatBool(evmResult)),
			}
			if len(evmErr) > 0 {
				attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyBridgeCallError, evmErr))
			}
			ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(types.EventTypeBridgeCallEvent, attrs...)})
		}()
		txResp, err := k.evmKeeper.CallEVM(ctx, sender, to, value.BigInt(), gasLimit, message, true)
		if err != nil {
			evmErr = err.Error()
			return err
		}
		evmResult = !txResp.Failed()
		evmErr = txResp.VmError
		if txResp.Failed() {
			return errorsmod.Wrap(types.ErrInvalid, evmErr)
		}
	}
	return nil
}

func (k Keeper) bridgeCallTransferToSender(ctx sdk.Context, receiver sdk.AccAddress, tokens []types.ERC20Token) (sdk.Coins, error) {
	mintCoins := sdk.NewCoins()
	unlockCoins := sdk.NewCoins()
	for i := 0; i < len(tokens); i++ {
		bridgeToken := k.GetBridgeTokenDenom(ctx, tokens[i].Contract)
		if bridgeToken == nil {
			return nil, errorsmod.Wrap(types.ErrInvalid, "bridge token is not exist")
		}
		if !tokens[i].Amount.IsPositive() {
			continue
		}
		coin := sdk.NewCoin(bridgeToken.Denom, tokens[i].Amount)
		isOriginOrConverted := k.erc20Keeper.IsOriginOrConvertedDenom(ctx, bridgeToken.Denom)
		if !isOriginOrConverted {
			mintCoins = mintCoins.Add(coin)
		}
		unlockCoins = unlockCoins.Add(coin)
	}
	if mintCoins.IsAllPositive() {
		if err := k.bankKeeper.MintCoins(ctx, k.moduleName, mintCoins); err != nil {
			return nil, errorsmod.Wrapf(err, "mint vouchers coins")
		}
	}
	if unlockCoins.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, receiver, unlockCoins); err != nil {
			return nil, errorsmod.Wrap(err, "transfer vouchers")
		}
	}

	targetCoins := sdk.NewCoins()
	for _, coin := range unlockCoins {
		targetCoin, err := k.erc20Keeper.ConvertDenomToTarget(ctx, receiver, coin, fxtypes.ParseFxTarget(fxtypes.ERC20Target))
		if err != nil {
			return nil, errorsmod.Wrap(err, "convert to target coin")
		}
		targetCoins = targetCoins.Add(targetCoin)
	}
	return targetCoins, nil
}

func (k Keeper) bridgeCallTransferToReceiver(ctx sdk.Context, sender sdk.AccAddress, receiver []byte, coins sdk.Coins) error {
	for _, coin := range coins {
		if coin.Denom == fxtypes.DefaultDenom {
			if err := k.bankKeeper.SendCoins(ctx, sender, receiver, sdk.NewCoins(coin)); err != nil {
				return err
			}
			continue
		}
		if _, err := k.erc20Keeper.ConvertCoin(sdk.WrapSDKContext(ctx), &erc20types.MsgConvertCoin{
			Coin:     coin,
			Receiver: common.BytesToAddress(receiver).String(),
			Sender:   sender.String(),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) bridgeCallCoinsHandler(ctx sdk.Context, sender sdk.AccAddress, coins sdk.Coins) ([]types.ERC20Token, error) {
	tokens := make([]types.ERC20Token, 0, len(coins))
	for _, coin := range coins {
		targetCoin, err := k.erc20Keeper.ConvertDenomToTarget(ctx, sender, coin, fxtypes.ParseFxTarget(k.moduleName))
		if err != nil {
			return nil, err
		}
		bridgeToken := k.GetDenomBridgeToken(ctx, targetCoin.Denom)
		if bridgeToken == nil {
			return nil, errorsmod.Wrap(types.ErrInvalid, "bridge token not found")
		}

		isOriginOrConverted := k.erc20Keeper.IsOriginOrConvertedDenom(ctx, targetCoin.Denom)
		if isOriginOrConverted {
			// lock coins in module
			if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, k.moduleName, sdk.NewCoins(targetCoin)); err != nil {
				return nil, err
			}
		} else {
			// If it is an external blockchain asset we burn it send coins to module in prep for burn
			if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, k.moduleName, sdk.NewCoins(targetCoin)); err != nil {
				return nil, err
			}
			// burn vouchers to send them back to external blockchain
			if err = k.bankKeeper.BurnCoins(ctx, k.moduleName, sdk.NewCoins(targetCoin)); err != nil {
				return nil, err
			}
		}
		tokens = append(tokens, types.NewERC20Token(targetCoin.Amount, bridgeToken.Token))
	}
	return tokens, nil
}
