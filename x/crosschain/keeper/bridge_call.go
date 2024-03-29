package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
		sdk.NewAttribute(types.AttributeKeyOutgoingBridgeCallId, fmt.Sprint(outCall.Nonce)),
	))

	return outCall, nil
}

func (k Keeper) SetOutgoingBridgeCall(ctx sdk.Context, out *types.OutgoingBridgeCall) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOutgoingBridgeCallKey(out.Nonce), k.cdc.MustMarshal(out))
}
