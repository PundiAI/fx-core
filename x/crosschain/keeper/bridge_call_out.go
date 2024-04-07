package keeper

import (
	"fmt"
	"math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"

	"github.com/functionx/fx-core/v7/contract"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) AddOutgoingBridgeCall(ctx sdk.Context, msg *types.MsgBridgeCall) (*types.OutgoingBridgeCall, error) {
	params := k.GetParams(ctx)
	batchTimeout := k.CalExternalTimeoutHeight(ctx, params, params.ExternalBatchTimeout)
	if batchTimeout <= 0 {
		return nil, errorsmod.Wrap(types.ErrInvalid, "bridge call timeout height")
	}

	nextID := k.autoIncrementID(ctx, types.KeyLastBridgeCallID)

	senderAddr := sdk.MustAccAddressFromBech32(msg.Sender)
	outCall := &types.OutgoingBridgeCall{
		Nonce:    nextID,
		Timeout:  batchTimeout,
		Sender:   fxtypes.AddressToStr(senderAddr.Bytes(), k.moduleName),
		Receiver: msg.Receiver,
		To:       msg.To,
		Asset:    msg.Asset,
		Message:  msg.Message,
		Value:    msg.Value,
		GasLimit: msg.GasLimit,
	}
	k.SetOutgoingBridgeCall(ctx, outCall)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCall,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyBridgeCallNonce, fmt.Sprint(outCall.Nonce)),
	))

	return outCall, nil
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

func (k Keeper) HandleOutgoingBridgeCallRefund(ctx sdk.Context, data *types.OutgoingBridgeCall) {
	receiveAddr := types.ExternalAddressToAccAddress(k.moduleName, data.GetSender())
	if err := k.bridgeCallAssetRefundHandler(ctx, receiveAddr, data.Asset); err != nil {
		panic(err)
	}
}

func (k Keeper) bridgeCallAssetRefundHandler(ctx sdk.Context, receive sdk.AccAddress, asset string) error {
	assetType, assetData, err := types.UnpackAssetType(asset)
	if err != nil {
		return errorsmod.Wrap(types.ErrInvalid, "asset")
	}

	switch assetType {
	case contract.AssetERC20:
		tokenAddresses, amounts, err := contract.UnpackERC20Asset(assetData)
		if err != nil {
			return errorsmod.Wrap(types.ErrInvalid, "erc20 token")
		}
		tokens, err := types.NewERC20Tokens(k.moduleName, tokenAddresses, amounts)
		if err != nil {
			return errorsmod.Wrap(types.ErrInvalid, err.Error())
		}
		coins, err := k.bridgeCallTransferToSender(ctx, receive, tokens)
		if err != nil {
			return err
		}
		return k.bridgeCallTransferToReceiver(ctx, receive, receive, coins)
	default:
		return errorsmod.Wrap(types.ErrInvalid, "asset type")
	}
}

func (k Keeper) IterateBridgeCallByNonce(ctx sdk.Context, startNonce uint64, cb func(bridgeCall *types.OutgoingBridgeCall) bool) {
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
