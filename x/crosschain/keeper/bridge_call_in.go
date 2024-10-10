package keeper

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-metrics"

	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
)

func (k Keeper) BridgeCallHandler(ctx sdk.Context, msg *types.MsgBridgeCallClaim) error {
	k.CreateBridgeAccount(ctx, msg.TxOrigin)
	if senderAccount := k.ak.GetAccount(ctx, msg.GetSenderAddr().Bytes()); senderAccount != nil {
		if _, ok := senderAccount.(sdk.ModuleAccountI); ok {
			return types.ErrInvalid.Wrap("sender is module account")
		}
	}
	isMemoSendCallTo := msg.IsMemoSendCallTo()
	receiverAddr := msg.GetToAddr()
	if isMemoSendCallTo {
		receiverAddr = msg.GetSenderAddr()
	}

	baseCoins := sdk.NewCoins()
	for i, address := range msg.TokenContracts {
		baseCoin, err := k.BridgeTokenToBaseCoin(ctx, address, msg.Amounts[i], receiverAddr.Bytes())
		if err != nil {
			return err
		}
		baseCoins = baseCoins.Add(baseCoin)
	}

	cacheCtx, commit := sdk.UnwrapSDKContext(ctx).CacheContext()
	var err error
	if err = k.BridgeCallEvm(cacheCtx, msg.GetSenderAddr(), msg.GetRefundAddr(), msg.GetToAddr(),
		receiverAddr, baseCoins, msg.MustData(), msg.MustMemo(), msg.Value, isMemoSendCallTo); err == nil {
		commit()
		return nil
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeBridgeCallEvent, sdk.NewAttribute(types.AttributeKeyErrCause, err.Error())))
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
	return k.BridgeCallFailedRefund(ctx, msg.GetRefundAddr(), baseCoins, msg.EventNonce)
}

func (k Keeper) BridgeCallEvm(ctx sdk.Context, sender, refundAddr, to, receiverAddr common.Address, baseCoins sdk.Coins, data, memo []byte, value sdkmath.Int, isMemoSendCallTo bool) error {
	tokens := make([]common.Address, 0, baseCoins.Len())
	amounts := make([]*big.Int, 0, baseCoins.Len())
	for _, coin := range baseCoins {
		tokenContract, err := k.BaseCoinToEvm(ctx, coin, receiverAddr)
		if err != nil {
			return err
		}
		tokens = append(tokens, common.HexToAddress(tokenContract))
		amounts = append(amounts, coin.Amount.BigInt())
	}

	if !k.evmKeeper.IsContract(ctx, to) {
		return nil
	}
	var callEvmSender common.Address
	var args []byte

	if isMemoSendCallTo {
		args = data
		callEvmSender = sender
	} else {
		var err error
		args, err = types.PackBridgeCallback(sender, refundAddr, tokens, amounts, data, memo)
		if err != nil {
			return err
		}
		callEvmSender = k.GetCallbackFrom()
	}

	gasLimit := k.GetBridgeCallMaxGasLimit(ctx)
	txResp, err := k.evmKeeper.CallEVM(ctx, callEvmSender, &to, value.BigInt(), gasLimit, args, true)
	if err != nil {
		return err
	}
	if txResp.Failed() {
		return types.ErrInvalid.Wrap(txResp.VmError)
	}
	return nil
}

func (k Keeper) BridgeCallFailedRefund(ctx sdk.Context, refundAddr common.Address, baseCoins sdk.Coins, eventNonce uint64) error {
	outCallNonce, err := k.AddOutgoingBridgeCall(ctx, refundAddr, refundAddr, baseCoins, common.Address{}, nil, nil, eventNonce)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCallRefundOut,
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprintf("%d", eventNonce)),
		sdk.NewAttribute(types.AttributeKeyBridgeCallNonce, fmt.Sprintf("%d", outCallNonce)),
	))
	return nil
}

func (k Keeper) bridgeCallTransferCoins(ctx sdk.Context, sender sdk.AccAddress, tokens []types.ERC20Token) (sdk.Coins, error) {
	mintCoins := sdk.NewCoins()
	unlockCoins := sdk.NewCoins()
	for i := 0; i < len(tokens); i++ {
		bridgeDenom, found := k.GetBridgeDenomByContract(ctx, tokens[i].Contract)
		if !found {
			return nil, types.ErrInvalid.Wrapf("bridge token is not exist")
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
			return nil, err
		}
	}
	if unlockCoins.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, sender, unlockCoins); err != nil {
			return nil, err
		}
	}

	targetCoins := sdk.NewCoins()
	for _, coin := range unlockCoins {
		targetCoin, err := k.erc20Keeper.ConvertDenomToTarget(ctx, sender, coin, fxtypes.ParseFxTarget(fxtypes.ERC20Target))
		if err != nil {
			return nil, err
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
