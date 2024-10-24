package keeper

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashicorp/go-metrics"

	fxtelemetry "github.com/functionx/fx-core/v8/telemetry"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (k Keeper) BuildOutgoingTxBatch(ctx sdk.Context, sender sdk.AccAddress, receiver string, amount, fee sdk.Coin) (uint64, error) {
	bridgeToken, err := k.BaseCoinToBridgeToken(ctx, sender, amount.Add(fee))
	if err != nil {
		return 0, err
	}
	if err = k.WithdrawBridgeToken(ctx, sender, amount.Amount, bridgeToken); err != nil {
		return 0, err
	}

	feeReceive := "" // todo: query feeReceive from quote contract

	batchTimeout := k.CalExternalTimeoutHeight(ctx, GetExternalBatchTimeout)
	if batchTimeout <= 0 {
		return 0, types.ErrInvalid.Wrapf("batch timeout height")
	}
	batch := &types.OutgoingTxBatch{
		BatchNonce:   k.autoIncrementID(ctx, types.KeyLastOutgoingBatchID),
		BatchTimeout: batchTimeout,
		Transactions: []*types.OutgoingTransferTx{
			{
				Id:          k.autoIncrementID(ctx, types.KeyLastTxPoolID),
				Sender:      sender.String(),
				DestAddress: receiver,
				Token:       types.NewERC20Token(amount.Amount, bridgeToken.Contract),
				Fee:         types.NewERC20Token(fee.Amount, bridgeToken.Contract),
			},
		},
		TokenContract: bridgeToken.Contract,
		FeeReceive:    feeReceive,
		Block:         uint64(ctx.BlockHeight()), // set the current block height when storing the batch
	}
	if err = k.StoreBatch(ctx, batch); err != nil {
		return 0, err
	}

	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatch,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchNonce, fmt.Sprint(batch.BatchNonce)),
		sdk.NewAttribute(types.AttributeKeyOutgoingTxIds, fmt.Sprint(batch.Transactions[0].Id)),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchTimeout, fmt.Sprint(batch.BatchTimeout)),
	)
	ctx.EventManager().EmitEvent(batchEvent)

	if !ctx.IsCheckTx() {
		fxtelemetry.SetGaugeLabelsWithDenom(
			[]string{types.ModuleName, "outgoing_tx_amount"},
			amount.Denom, amount.Amount.BigInt(),
			telemetry.NewLabel("module", k.moduleName),
		)
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "outgoing_tx"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("module", k.moduleName),
				telemetry.NewLabel("denom", amount.Denom),
			},
		)
	}
	return batch.BatchNonce, nil
}

func (k Keeper) OutgoingTxBatchExecuted(ctx sdk.Context, tokenContract string, batchNonce uint64) (err error) {
	batch := k.GetOutgoingTxBatch(ctx, tokenContract, batchNonce)
	if batch == nil {
		return fmt.Errorf("unknown batch nonce for outgoing tx batch %s %d", tokenContract, batchNonce)
	}

	// Iterate through remaining batches
	k.IterateOutgoingTxBatches(ctx, func(iterBatch *types.OutgoingTxBatch) bool {
		// If the iterated batches nonce is lower than the one that was just executed, resend it
		if iterBatch.BatchNonce < batch.BatchNonce && iterBatch.TokenContract == tokenContract {
			if err = k.ResendTimeoutOutgoingTxBatch(ctx, iterBatch); err != nil {
				return true
			}
		}
		return false
	})
	if err != nil {
		return err
	}

	// Delete batch since it is finished
	k.DeleteBatch(ctx, batch)
	k.DeleteBatchConfirm(ctx, batch.BatchNonce, batch.TokenContract)
	return nil
}

// --- OUTGOING TX BATCH --- //

// StoreBatch stores a transaction batch
func (k Keeper) StoreBatch(ctx sdk.Context, batch *types.OutgoingTxBatch) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce)
	store.Set(key, k.cdc.MustMarshal(batch))

	blockKey := types.GetOutgoingTxBatchBlockKey(batch.Block, batch.BatchNonce)
	// Note: Only one OutgoingTxBatch can be submitted in a block
	if store.Has(blockKey) {
		return types.ErrInvalid.Wrapf("block:[%v] has batch request", batch.Block)
	}
	store.Set(blockKey, k.cdc.MustMarshal(batch))
	return nil
}

// DeleteBatch deletes an outgoing transaction batch
func (k Keeper) DeleteBatch(ctx sdk.Context, batch *types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce))
	store.Delete(types.GetOutgoingTxBatchBlockKey(batch.Block, batch.BatchNonce))
}

// GetOutgoingTxBatch loads a batch object. Returns nil when not exists.
func (k Keeper) GetOutgoingTxBatch(ctx sdk.Context, tokenContract string, batchNonce uint64) *types.OutgoingTxBatch {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxBatchKey(tokenContract, batchNonce)
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil
	}
	batch := new(types.OutgoingTxBatch)
	k.cdc.MustUnmarshal(bz, batch)
	return batch
}

func (k Keeper) ResendTimeoutOutgoingTxBatch(ctx sdk.Context, batch *types.OutgoingTxBatch) error {
	k.DeleteBatch(ctx, batch)
	k.DeleteBatchConfirm(ctx, batch.BatchNonce, batch.TokenContract)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeOutgoingBatchCanceled,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchNonce, fmt.Sprint(batch.BatchNonce)),
	))

	batchTimeout := k.CalExternalTimeoutHeight(ctx, GetExternalBatchTimeout)
	if batchTimeout <= 0 {
		return types.ErrInvalid.Wrapf("batch timeout height")
	}
	batch.BatchTimeout = batchTimeout
	batch.BatchNonce = k.autoIncrementID(ctx, types.KeyLastOutgoingBatchID)
	if err := k.StoreBatch(ctx, batch); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeOutgoingBatch,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchNonce, fmt.Sprint(batch.BatchNonce)),
		sdk.NewAttribute(types.AttributeKeyOutgoingTxIds, fmt.Sprint(batch.Transactions[0].Id)),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchTimeout, fmt.Sprint(batch.BatchTimeout)),
	))
	return nil
}

// IterateOutgoingTxBatches iterates through all outgoing batches
func (k Keeper) IterateOutgoingTxBatches(ctx sdk.Context, cb func(batch *types.OutgoingTxBatch) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStoreReversePrefixIterator(store, types.OutgoingTxBatchKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		batch := new(types.OutgoingTxBatch)
		k.cdc.MustUnmarshal(iter.Value(), batch)
		// cb returns true to stop early
		if cb(batch) {
			break
		}
	}
}

// GetOutgoingTxBatches used in testing
func (k Keeper) GetOutgoingTxBatches(ctx sdk.Context) (out []*types.OutgoingTxBatch) {
	k.IterateOutgoingTxBatches(ctx, func(batch *types.OutgoingTxBatch) bool {
		out = append(out, batch)
		return false
	})
	return
}

// GetLastOutgoingBatchByToken gets the latest outgoing tx batch by token type
func (k Keeper) GetLastOutgoingBatchByToken(ctx sdk.Context, token string) *types.OutgoingTxBatch {
	var lastBatch *types.OutgoingTxBatch
	lastNonce := uint64(0)
	k.IterateOutgoingTxBatches(ctx, func(batch *types.OutgoingTxBatch) bool {
		if batch.TokenContract == token && batch.BatchNonce > lastNonce {
			lastBatch = batch
			lastNonce = batch.BatchNonce
		}
		return false
	})
	return lastBatch
}

// IterateBatchByBlockHeight iterates through all Batch by block in the half-open interval [start,end)
func (k Keeper) IterateBatchByBlockHeight(ctx sdk.Context, start, end uint64, cb func(*types.OutgoingTxBatch) bool) {
	store := ctx.KVStore(k.storeKey)
	startKey := append(types.OutgoingTxBatchBlockKey, sdk.Uint64ToBigEndian(start)...)
	endKey := append(types.OutgoingTxBatchBlockKey, sdk.Uint64ToBigEndian(end)...)
	iter := store.Iterator(startKey, endKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		batch := new(types.OutgoingTxBatch)
		k.cdc.MustUnmarshal(iter.Value(), batch)
		// cb returns true to stop early
		if cb(batch) {
			break
		}
	}
}
