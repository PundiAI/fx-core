package keeper

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-metrics"

	fxtelemetry "github.com/functionx/fx-core/v8/telemetry"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
)

func (k Keeper) BridgeCallHandler(ctx sdk.Context, msg *types.MsgBridgeCallClaim) error {
	k.CreateBridgeAccount(ctx, msg.TxOrigin)

	tokens := msg.GetTokensAddr()
	erc20Token, err := types.NewERC20Tokens(k.moduleName, tokens, msg.GetAmounts())
	if err != nil {
		return err
	}
	refundAddr := msg.GetRefundAddr()
	cacheCtx, commit := ctx.CacheContext()
	if err = k.BridgeCallTransferAndCallEvm(cacheCtx, msg.GetSenderAddr(), refundAddr, erc20Token, msg.GetToAddr(), msg.MustData(), msg.MustMemo(), msg.Value); err != nil {
		ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeBridgeCallEvent, sdk.NewAttribute(types.AttributeKeyErrCause, err.Error())))
	} else {
		commit()
	}

	if !ctx.IsCheckTx() {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "bridge_call_in"},
			float32(1),
			[]metrics.Label{
				telemetry.NewLabel("module", k.moduleName),
				telemetry.NewLabel("success", strconv.FormatBool(err == nil)),
			},
		)
	}

	if err != nil && len(tokens) > 0 {
		// new outgoing bridge call to refund
		outCallNonce, err := k.AddOutgoingBridgeCall(ctx, refundAddr, refundAddr, erc20Token, common.Address{}, nil, nil, msg.EventNonce)
		if err != nil {
			return err
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeBridgeCallRefundOut,
			sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprintf("%d", msg.EventNonce)),
			sdk.NewAttribute(types.AttributeKeyBridgeCallNonce, fmt.Sprintf("%d", outCallNonce)),
		))
	}

	if err == nil {
		for i := 0; i < len(erc20Token); i++ {
			bridgeDenom, found := k.GetBridgeDenomByContract(ctx, erc20Token[i].Contract)
			if !found {
				continue
			}
			// no need for a double check here, as the bridge token should exist
			k.HandlePendingOutgoingTx(ctx, refundAddr.Bytes(), msg.EventNonce, bridgeDenom, erc20Token[i].Contract)
			k.HandlePendingOutgoingBridgeCall(ctx, refundAddr.Bytes(), msg.EventNonce, bridgeDenom)
		}
	}
	return nil
}

func (k Keeper) BridgeCallTransferAndCallEvm(ctx sdk.Context, sender, refundAddr common.Address, tokens []types.ERC20Token, to common.Address, data, memo []byte, value sdkmath.Int) error {
	if senderAccount := k.ak.GetAccount(ctx, sender.Bytes()); senderAccount != nil {
		if _, ok := senderAccount.(sdk.ModuleAccountI); ok {
			return errorsmod.Wrap(types.ErrInvalid, "sender is module account")
		}
	}
	isMemoSendCallTo := types.IsMemoSendCallTo(memo)
	receiverTokenAddr := to
	if isMemoSendCallTo {
		receiverTokenAddr = sender
	}
	coins, err := k.bridgeCallTransferCoins(ctx, receiverTokenAddr.Bytes(), tokens)
	if err != nil {
		return err
	}

	if !ctx.IsCheckTx() {
		fxtelemetry.SetGaugeLabelsWithCoins(
			[]string{types.ModuleName, "bridge_call_in_amount"},
			coins,
			telemetry.NewLabel("module", k.moduleName),
		)
	}

	if err = k.bridgeCallTransferTokens(ctx, receiverTokenAddr.Bytes(), receiverTokenAddr.Bytes(), coins); err != nil {
		return err
	}
	return k.BridgeCallEvm(ctx, sender, refundAddr, coins, to, data, memo, value, isMemoSendCallTo)
}

func (k Keeper) bridgeCallTransferCoins(ctx sdk.Context, sender sdk.AccAddress, tokens []types.ERC20Token) (sdk.Coins, error) {
	mintCoins := sdk.NewCoins()
	unlockCoins := sdk.NewCoins()
	for i := 0; i < len(tokens); i++ {
		bridgeDenom, found := k.GetBridgeDenomByContract(ctx, tokens[i].Contract)
		if !found {
			return nil, errorsmod.Wrap(types.ErrInvalid, "bridge token is not exist")
		}
		if !tokens[i].Amount.IsPositive() {
			continue
		}
		coin := sdk.NewCoin(bridgeDenom, tokens[i].Amount)
		isOriginOrConverted := k.erc20Keeper.IsOriginOrConvertedDenom(ctx, bridgeDenom)
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
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, sender, unlockCoins); err != nil {
			return nil, errorsmod.Wrap(err, "transfer vouchers")
		}
	}

	targetCoins := sdk.NewCoins()
	for _, coin := range unlockCoins {
		targetCoin, err := k.erc20Keeper.ConvertDenomToTarget(ctx, sender, coin, fxtypes.ParseFxTarget(fxtypes.ERC20Target))
		if err != nil {
			return nil, errorsmod.Wrap(err, "convert to target coin")
		}
		targetCoins = targetCoins.Add(targetCoin)
	}
	return targetCoins, nil
}

func (k Keeper) bridgeCallTransferTokens(ctx sdk.Context, sender sdk.AccAddress, receiver []byte, coins sdk.Coins) error {
	for _, coin := range coins {
		if coin.Denom == fxtypes.DefaultDenom {
			if bytes.Equal(sender, receiver) {
				continue
			}
			if err := k.bankKeeper.SendCoins(ctx, sender, receiver, sdk.NewCoins(coin)); err != nil {
				return err
			}
			continue
		}
		if _, err := k.erc20Keeper.ConvertCoin(ctx, &erc20types.MsgConvertCoin{
			Coin:     coin,
			Receiver: common.BytesToAddress(receiver).String(),
			Sender:   sender.String(),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) CoinsToBridgeCallTokens(ctx sdk.Context, coins sdk.Coins) ([]common.Address, []*big.Int) {
	tokens := make([]common.Address, 0, len(coins))
	amounts := make([]*big.Int, 0, len(coins))
	for _, coin := range coins {
		amounts = append(amounts, coin.Amount.BigInt())
		if coin.Denom == fxtypes.DefaultDenom {
			tokens = append(tokens, common.Address{})
			continue
		}
		// bridgeCallTransferTokens().ConvertCoin hava already checked.
		pair, _ := k.erc20Keeper.GetTokenPair(ctx, coin.Denom)
		tokens = append(tokens, common.HexToAddress(pair.Erc20Address))
	}
	return tokens, amounts
}
