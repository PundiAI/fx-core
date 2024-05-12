package keeper

import (
	"fmt"
	"math"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) bridgeCallCoinsToERC20Token(ctx sdk.Context, sender sdk.AccAddress, coins sdk.Coins) ([]types.ERC20Token, error) {
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

func (k Keeper) AddOutgoingBridgeCall(
	ctx sdk.Context,
	sender sdk.AccAddress,
	receiver string,
	tokens []types.ERC20Token,
	to string,
	data string,
	memo string,
) (*types.OutgoingBridgeCall, error) {
	params := k.GetParams(ctx)
	bridgeCallTimeout := k.CalExternalTimeoutHeight(ctx, params, params.BridgeCallTimeout)
	if bridgeCallTimeout <= 0 {
		return nil, errorsmod.Wrap(types.ErrInvalid, "bridge call timeout height")
	}

	nextID := k.autoIncrementID(ctx, types.KeyLastBridgeCallID)

	bridgeCall := &types.OutgoingBridgeCall{
		Nonce:       nextID,
		Timeout:     bridgeCallTimeout,
		BlockHeight: uint64(ctx.BlockHeight()),
		Sender:      types.ExternalAddrToStr(k.moduleName, sender),
		Receiver:    receiver,
		Tokens:      tokens,
		To:          to,
		Data:        data,
		Memo:        memo,
	}
	k.SetOutgoingBridgeCall(ctx, bridgeCall)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCall,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, bridgeCall.Sender),
		sdk.NewAttribute(types.AttributeKeyBridgeCallNonce, fmt.Sprint(bridgeCall.Nonce)),
	))

	return bridgeCall, nil
}

func (k Keeper) BridgeCallResultHandler(ctx sdk.Context, claim *types.MsgBridgeCallResultClaim) {
	k.CreateBridgeAccount(ctx, claim.TxOrigin)

	outgoingBridgeCall, found := k.GetOutgoingBridgeCallByNonce(ctx, claim.Nonce)
	if !found {
		panic(fmt.Errorf("bridge call not found for nonce %d", claim.Nonce))
	}
	if !claim.Success {
		k.HandleOutgoingBridgeCallRefund(ctx, outgoingBridgeCall)
	}
	k.DeleteOutgoingBridgeCallRecord(ctx, claim.Nonce)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCallResult,
		sdk.NewAttribute(types.AttributeKeyEventNonce, strconv.FormatInt(int64(claim.Nonce), 10)),
		sdk.NewAttribute(types.AttributeKeyStateSuccess, strconv.FormatBool(claim.Success)),
		sdk.NewAttribute(types.AttributeKeyErrCause, claim.Cause),
	))
}

func (k Keeper) SetOutgoingBridgeCall(ctx sdk.Context, out *types.OutgoingBridgeCall) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOutgoingBridgeCallNonceKey(out.Nonce), k.cdc.MustMarshal(out))
	// value is just a placeholder
	store.Set(types.GetOutgoingBridgeCallAddressAndNonceKey(out.Sender, out.Nonce), k.cdc.MustMarshal(&gogotypes.BoolValue{Value: true}))
}

func (k Keeper) GetOutgoingBridgeCallByNonce(ctx sdk.Context, nonce uint64) (*types.OutgoingBridgeCall, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOutgoingBridgeCallNonceKey(nonce))
	if bz == nil {
		return nil, false
	}
	var out types.OutgoingBridgeCall
	k.cdc.MustUnmarshal(bz, &out)
	return &out, true
}

func (k Keeper) DeleteOutgoingBridgeCall(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOutgoingBridgeCallNonceKey(nonce))
}

func (k Keeper) IterateOutgoingBridgeCalls(ctx sdk.Context, cb func(*types.OutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OutgoingBridgeCallNonceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var value types.OutgoingBridgeCall
		k.cdc.MustUnmarshal(iterator.Value(), &value)
		if cb(&value) {
			break
		}
	}
}

func (k Keeper) IterateOutgoingBridgeCallsByAddress(ctx sdk.Context, addr string, cb func(record *types.OutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetOutgoingBridgeCallAddressKey(addr))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		nonce := types.ParseOutgoingBridgeCallNonce(iterator.Key(), addr)
		record, found := k.GetOutgoingBridgeCallByNonce(ctx, nonce)
		if !found {
			continue
		}
		if cb(record) {
			break
		}
	}
}

func (k Keeper) IterateOutgoingBridgeCallByNonce(ctx sdk.Context, startNonce uint64, cb func(bridgeCall *types.OutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	startKey := append(types.OutgoingBridgeCallNonceKey, sdk.Uint64ToBigEndian(startNonce)...)
	endKey := append(types.OutgoingBridgeCallNonceKey, sdk.Uint64ToBigEndian(math.MaxUint64)...)
	iter := store.Iterator(startKey, endKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		value := new(types.OutgoingBridgeCall)
		k.cdc.MustUnmarshal(iter.Value(), value)
		if cb(value) {
			break
		}
	}
}

func (k Keeper) SetBridgeCallFromMsg(ctx sdk.Context, txID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetBridgeCallFromMsgKey(txID), []byte{})
}

func (k Keeper) DeleteBridgeCallFromMsg(ctx sdk.Context, txID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetBridgeCallFromMsgKey(txID))
}

func (k Keeper) HasBridgeCallFromMsg(ctx sdk.Context, txID uint64) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetBridgeCallFromMsgKey(txID))
}
