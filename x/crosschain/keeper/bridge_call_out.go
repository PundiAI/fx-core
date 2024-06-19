package keeper

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	gogotypes "github.com/gogo/protobuf/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
)

func (k Keeper) BridgeCallCoinsToERC20Token(ctx sdk.Context, sender sdk.AccAddress, coins sdk.Coins) ([]types.ERC20Token, sdk.Coins, error) {
	tokens := make([]types.ERC20Token, 0, len(coins))
	notLiquidCoins := sdk.NewCoins()
	for _, coin := range coins {
		targetCoin, err := k.erc20Keeper.ConvertDenomToTarget(ctx, sender, coin, fxtypes.ParseFxTarget(k.moduleName))
		if err != nil && !erc20types.IsInsufficientLiquidityErr(err) {
			return nil, nil, err
		}
		bridgeToken := k.GetDenomBridgeToken(ctx, targetCoin.Denom)
		if bridgeToken == nil {
			return nil, nil, errorsmod.Wrap(types.ErrInvalid, "bridge token not found")
		}
		tokens = append(tokens, types.NewERC20Token(targetCoin.Amount, bridgeToken.Token))
		if erc20types.IsInsufficientLiquidityErr(err) {
			notLiquidCoins = notLiquidCoins.Add(targetCoin)
			continue
		}
		if err = k.TransferBridgeCoinToExternal(ctx, sender, targetCoin); err != nil {
			return nil, nil, err
		}
	}
	return tokens, notLiquidCoins, nil
}

func (k Keeper) AddOutgoingBridgeCall(ctx sdk.Context, sender, refundAddr common.Address, tokens []types.ERC20Token, to common.Address, data, memo []byte, eventNonce uint64) (uint64, error) {
	outCall, err := k.BuildOutgoingBridgeCall(ctx, sender, refundAddr, tokens, to, data, memo, eventNonce)
	if err != nil {
		return 0, err
	}
	return k.AddOutgoingBridgeCallWithoutBuild(ctx, outCall), nil
}

func (k Keeper) AddOutgoingBridgeCallWithoutBuild(ctx sdk.Context, outCall *types.OutgoingBridgeCall) uint64 {
	k.SetOutgoingBridgeCall(ctx, outCall)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCall,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, outCall.Sender),
		sdk.NewAttribute(types.AttributeKeyBridgeCallNonce, fmt.Sprint(outCall.Nonce)),
	))

	return outCall.Nonce
}

func (k Keeper) BuildOutgoingBridgeCall(ctx sdk.Context, sender common.Address, refundAddr common.Address, tokens []types.ERC20Token, to common.Address, data []byte, memo []byte, eventNonce uint64) (*types.OutgoingBridgeCall, error) {
	bridgeCallTimeout := k.CalExternalTimeoutHeight(ctx, GetBridgeCallTimeout)
	if bridgeCallTimeout <= 0 {
		return nil, errorsmod.Wrap(types.ErrInvalid, "bridge call timeout height")
	}

	nextID := k.autoIncrementID(ctx, types.KeyLastBridgeCallID)

	outCall := &types.OutgoingBridgeCall{
		Nonce:       nextID,
		Timeout:     bridgeCallTimeout,
		BlockHeight: uint64(ctx.BlockHeight()),
		Sender:      types.ExternalAddrToStr(k.moduleName, sender.Bytes()),
		Refund:      types.ExternalAddrToStr(k.moduleName, refundAddr.Bytes()),
		Tokens:      tokens,
		To:          types.ExternalAddrToStr(k.moduleName, to.Bytes()),
		Data:        hex.EncodeToString(data),
		Memo:        hex.EncodeToString(memo),
		EventNonce:  eventNonce,
	}
	return outCall, nil
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

func (k Keeper) SetOutgoingBridgeCall(ctx sdk.Context, outCall *types.OutgoingBridgeCall) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOutgoingBridgeCallNonceKey(outCall.Nonce), k.cdc.MustMarshal(outCall))
	// value is just a placeholder
	store.Set(
		types.GetOutgoingBridgeCallAddressAndNonceKey(outCall.Sender, outCall.Nonce),
		k.cdc.MustMarshal(&gogotypes.BoolValue{Value: true}),
	)
}

func (k Keeper) HasOutgoingBridgeCall(ctx sdk.Context, nonce uint64) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetOutgoingBridgeCallNonceKey(nonce))
}

func (k Keeper) HasOutgoingBridgeCallAddressAndNonce(ctx sdk.Context, sender string, nonce uint64) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetOutgoingBridgeCallAddressAndNonceKey(sender, nonce))
}

func (k Keeper) GetOutgoingBridgeCallByNonce(ctx sdk.Context, nonce uint64) (*types.OutgoingBridgeCall, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOutgoingBridgeCallNonceKey(nonce))
	if bz == nil {
		return nil, false
	}
	var outCall types.OutgoingBridgeCall
	k.cdc.MustUnmarshal(bz, &outCall)
	return &outCall, true
}

func (k Keeper) DeleteOutgoingBridgeCall(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	outCall, found := k.GetOutgoingBridgeCallByNonce(ctx, nonce)
	if !found {
		return
	}
	store.Delete(types.GetOutgoingBridgeCallNonceKey(nonce))
	store.Delete(types.GetOutgoingBridgeCallAddressAndNonceKey(outCall.Sender, outCall.Nonce))
}

func (k Keeper) IterateOutgoingBridgeCalls(ctx sdk.Context, cb func(outCall *types.OutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OutgoingBridgeCallNonceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var outCall types.OutgoingBridgeCall
		k.cdc.MustUnmarshal(iterator.Value(), &outCall)
		if cb(&outCall) {
			break
		}
	}
}

func (k Keeper) IterateOutgoingBridgeCallsByAddress(ctx sdk.Context, senderAddr string, cb func(outCall *types.OutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetOutgoingBridgeCallAddressKey(senderAddr))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		nonce := types.ParseOutgoingBridgeCallNonce(iterator.Key(), senderAddr)
		outCall, found := k.GetOutgoingBridgeCallByNonce(ctx, nonce)
		if !found {
			continue
		}
		if cb(outCall) {
			break
		}
	}
}

func (k Keeper) IterateOutgoingBridgeCallByNonce(ctx sdk.Context, startNonce uint64, cb func(outCall *types.OutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	startKey := append(types.OutgoingBridgeCallNonceKey, sdk.Uint64ToBigEndian(startNonce)...)
	endKey := append(types.OutgoingBridgeCallNonceKey, sdk.Uint64ToBigEndian(math.MaxUint64)...)
	iter := store.Iterator(startKey, endKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		outCall := new(types.OutgoingBridgeCall)
		k.cdc.MustUnmarshal(iter.Value(), outCall)
		if cb(outCall) {
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
