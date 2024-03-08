package keeper

import (
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) AddRefundRecord(ctx sdk.Context, receiver string, eventNonce uint64, tokens []types.ERC20Token) error {
	oracleSet := k.GetLatestOracleSet(ctx)
	if oracleSet == nil {
		return errorsmod.Wrap(types.ErrInvalid, "no oracle set")
	}
	snapshotOracle, found := k.GetSnapshotOracle(ctx, oracleSet.Nonce)
	if !found {
		snapshotOracle = &types.SnapshotOracle{
			OracleSetNonce: oracleSet.Nonce,
			Members:        oracleSet.Members,
			EventNonces:    []uint64{},
		}
	}
	snapshotOracle.EventNonces = append(snapshotOracle.EventNonces, eventNonce)
	k.SetSnapshotOracle(ctx, snapshotOracle)

	k.SetRefundRecord(ctx, &types.RefundRecord{
		EventNonce:     eventNonce,
		Receiver:       receiver,
		Timeout:        uint64(ctx.BlockTime().Add(time.Hour * 24 * 7).Second()), // TODO need to be configurable
		OracleSetNonce: oracleSet.Nonce,
		Tokens:         tokens,
	})
	return nil
}

func (k Keeper) SetRefundRecord(ctx sdk.Context, refundRecord *types.RefundRecord) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetBridgeCallRefundEventNonceKey(refundRecord.EventNonce), k.cdc.MustMarshal(refundRecord))
	store.Set(types.GetBridgeCallRefundKey(refundRecord.Receiver, refundRecord.EventNonce), sdk.Uint64ToBigEndian(refundRecord.OracleSetNonce))
}

func (k Keeper) GetRefundRecord(ctx sdk.Context, eventNonce uint64) (*types.RefundRecord, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBridgeCallRefundEventNonceKey(eventNonce))
	if bz == nil {
		return nil, false
	}
	refundRecord := new(types.RefundRecord)
	k.cdc.MustUnmarshal(bz, refundRecord)
	return refundRecord, true
}

func (k Keeper) SetSnapshotOracle(ctx sdk.Context, snapshotOracleKey *types.SnapshotOracle) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetSnapshotOracleKey(snapshotOracleKey.OracleSetNonce), k.cdc.MustMarshal(snapshotOracleKey))
}

func (k Keeper) GetSnapshotOracle(ctx sdk.Context, oracleSetNonce uint64) (*types.SnapshotOracle, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSnapshotOracleKey(oracleSetNonce))
	if bz == nil {
		return nil, false
	}
	snapshotOracle := new(types.SnapshotOracle)
	k.cdc.MustUnmarshal(bz, snapshotOracle)
	return snapshotOracle, true
}

func (k Keeper) GetRefundConfirm(ctx sdk.Context, nonce uint64, addr sdk.AccAddress) (*types.MsgConfirmRefund, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetRefundConfirmKey(nonce, addr))
	if bz == nil {
		return nil, false
	}
	var msg types.MsgConfirmRefund
	k.cdc.MustUnmarshal(bz, &msg)
	return &msg, true
}

func (k Keeper) SetRefundConfirm(ctx sdk.Context, addr sdk.AccAddress, msg *types.MsgConfirmRefund) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetRefundConfirmKey(msg.Nonce, addr), k.cdc.MustMarshal(msg))
}
