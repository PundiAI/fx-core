package keeper

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/v2/x/crosschain/types"
)

const OutgoingTxBatchSize = 100

// BuildOutgoingTxBatch starts the following process chain:
// - find bridged denominator for given voucher type
// - determine if a an unExecuted batch is already waiting for this token type, if so confirm the new batch would
//   have a higher total fees. If not exit without creating a batch
// - select available transactions from the outgoing transaction pool sorted by fee desc
// - persist an outgoing batch object with an incrementing ID = nonce
// - emit an event
func (k Keeper) BuildOutgoingTxBatch(ctx sdk.Context, tokenContract, feeReceive string, maxElements uint, minimumFee, baseFee sdk.Int) (*types.OutgoingTxBatch, error) {
	if maxElements == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "max elements value")
	}

	// if there is a more profitable batch for this token type do not create a new batch
	if lastBatch := k.GetLastOutgoingBatchByTokenType(ctx, tokenContract); lastBatch != nil {
		currentFees := k.GetBatchFeesByTokenType(ctx, tokenContract, maxElements, baseFee)
		if lastBatch.GetFees().GT(currentFees.TotalFees) {
			return nil, sdkerrors.Wrap(types.ErrInvalid, "new batch would not be more profitable")
		}
	}
	selectedTx, err := k.pickUnBatchedTx(ctx, tokenContract, maxElements, baseFee)
	if err != nil {
		return nil, err
	}
	if len(selectedTx) == 0 {
		return nil, sdkerrors.Wrap(types.ErrEmpty, "no batch tx")
	}
	if types.OutgoingTransferTxs(selectedTx).TotalFee().LT(minimumFee) {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "total fee less than minimum fee")
	}
	batchTimeout := k.GetBatchTimeoutHeight(ctx)
	if batchTimeout <= 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "batch timeout height")
	}
	nextID := k.autoIncrementID(ctx, types.KeyLastOutgoingBatchID)
	batch := &types.OutgoingTxBatch{
		BatchNonce:    nextID,
		BatchTimeout:  batchTimeout,
		Transactions:  selectedTx,
		TokenContract: tokenContract,
		FeeReceive:    feeReceive,
		Block:         uint64(ctx.BlockHeight()), // set the current block height when storing the batch
	}
	if err = k.StoreBatch(ctx, batch); err != nil {
		return nil, err
	}

	checkpoint, err := batch.GetCheckpoint(k.GetGravityID(ctx))
	if err != nil {
		panic(err)
	}
	k.SetPastExternalSignatureCheckpoint(ctx, checkpoint)

	eventBatchNonceTxIds := strings.Builder{}
	eventBatchNonceTxIds.WriteString(fmt.Sprintf("%d", selectedTx[0].Id))
	for _, tx := range selectedTx[1:] {
		_, _ = eventBatchNonceTxIds.WriteString(fmt.Sprintf(",%d", tx.Id))
	}
	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatch,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchNonce, fmt.Sprint(nextID)),
		sdk.NewAttribute(types.AttributeKeyOutgoingTxIds, eventBatchNonceTxIds.String()),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchTimeout, fmt.Sprint(batch.BatchTimeout)),
	)
	ctx.EventManager().EmitEvent(batchEvent)
	return batch, nil
}

// GetBatchTimeoutHeight This gets the batch timeout height in External blocks.
func (k Keeper) GetBatchTimeoutHeight(ctx sdk.Context) uint64 {
	currentFxHeight := ctx.BlockHeight()
	params := k.GetParams(ctx)
	// we store the last observed Cosmos and Ethereum heights, we do not concern ourselves if these values
	// are zero because no batch can be produced if the last Ethereum block height is not first populated by a deposit event.
	heights := k.GetLastObservedBlockHeight(ctx)
	if heights.ExternalBlockHeight == 0 {
		return 0
	}
	// we project how long it has been in milliseconds since the last Ethereum block height was observed
	projectedMillis := (uint64(currentFxHeight) - heights.BlockHeight) * params.AverageBlockTime
	// we convert that projection into the current Ethereum height using the average Ethereum block time in millis
	projectedCurrentEthereumHeight := (projectedMillis / params.AverageExternalBlockTime) + heights.ExternalBlockHeight
	// we convert our target time for block timeouts (lets say 12 hours) into a number of blocks to
	// place on top of our projection of the current Ethereum block height.
	blocksToAdd := params.ExternalBatchTimeout / params.AverageExternalBlockTime
	return projectedCurrentEthereumHeight + blocksToAdd
}

// OutgoingTxBatchExecuted is run when the Cosmos chain detects that a batch has been executed on Ethereum
// It frees all the transactions in the batch, then cancels all earlier batches
func (k Keeper) OutgoingTxBatchExecuted(ctx sdk.Context, tokenContract string, nonce uint64) {
	batch := k.GetOutgoingTxBatch(ctx, tokenContract, nonce)
	if batch == nil {
		panic(fmt.Sprintf("unknown batch nonce for outgoing tx batch %s %d", tokenContract, nonce))
	}

	// Iterate through remaining batches
	k.IterateOutgoingTxBatches(ctx, func(key []byte, iter_batch *types.OutgoingTxBatch) bool {
		// If the iterated batches nonce is lower than the one that was just executed, cancel it
		if iter_batch.BatchNonce < batch.BatchNonce && iter_batch.TokenContract == tokenContract {
			_ = k.CancelOutgoingTxBatch(ctx, tokenContract, iter_batch.BatchNonce)
		}
		return false
	})

	// Delete batch since it is finished
	k.DeleteBatch(ctx, batch)
}

// StoreBatch stores a transaction batch
func (k Keeper) StoreBatch(ctx sdk.Context, batch *types.OutgoingTxBatch) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce)
	store.Set(key, k.cdc.MustMarshal(batch))

	blockKey := types.GetOutgoingTxBatchBlockKey(batch.Block)

	if store.Has(blockKey) {
		return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("block:[%v] has batch request", batch.Block))
	}
	store.Set(blockKey, k.cdc.MustMarshal(batch))
	return nil
}

// DeleteBatch deletes an outgoing transaction batch
func (k Keeper) DeleteBatch(ctx sdk.Context, batch *types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce))
	store.Delete(types.GetOutgoingTxBatchBlockKey(batch.Block))
}

// pickUnBatchedTx find Tx in pool and remove from "available" second index
func (k Keeper) pickUnBatchedTx(ctx sdk.Context, tokenContract string, maxElements uint, baseFee sdk.Int) ([]*types.OutgoingTransferTx, error) {
	var selectedTx []*types.OutgoingTransferTx
	var err error
	k.IterateUnbatchedTransactionsByContract(ctx, tokenContract, func(_ []byte, tx *types.OutgoingTransferTx) bool {
		if tx.Fee.Amount.LT(baseFee) {
			return true
		}
		selectedTx = append(selectedTx, tx)
		err = k.removeUnbatchedTx(ctx, tx.Fee, tx.Id)
		oldTx, oldTxErr := k.GetUnbatchedTxByFeeAndId(ctx, tx.Fee, tx.Id)
		if oldTx != nil || oldTxErr == nil {
			panic("picked a duplicate transaction from the pool, duplicates should never exist!")
		}
		return err != nil || uint(len(selectedTx)) == maxElements
	})
	return selectedTx, err
}

// GetOutgoingTxBatch loads a batch object. Returns nil when not exists.
func (k Keeper) GetOutgoingTxBatch(ctx sdk.Context, tokenContract string, nonce uint64) *types.OutgoingTxBatch {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxBatchKey(tokenContract, nonce)
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil
	}
	var b types.OutgoingTxBatch
	k.cdc.MustUnmarshal(bz, &b)
	for _, tx := range b.Transactions {
		tx.Token.Contract = tokenContract
		tx.Fee.Contract = tokenContract
	}
	return &b
}

// CancelOutgoingTxBatch releases all TX in the batch and deletes the batch
func (k Keeper) CancelOutgoingTxBatch(ctx sdk.Context, tokenContract string, nonce uint64) error {
	batch := k.GetOutgoingTxBatch(ctx, tokenContract, nonce)
	if batch == nil {
		return types.ErrUnknown
	}
	for _, tx := range batch.Transactions {
		if err := k.AddUnbatchedTx(ctx, tx); err != nil {
			panic(sdkerrors.Wrapf(err, "unable to add batched transaction back into pool %v", tx))
		}
	}

	// Delete batch since it is finished
	k.DeleteBatch(ctx, batch)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeOutgoingBatchCanceled,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchNonce, fmt.Sprint(nonce)),
	))
	return nil
}

// IterateOutgoingTxBatches iterates through all outgoing batches in DESC order.
func (k Keeper) IterateOutgoingTxBatches(ctx sdk.Context, cb func(key []byte, batch *types.OutgoingTxBatch) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OutgoingTxBatchKey)
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var batch types.OutgoingTxBatch
		k.cdc.MustUnmarshal(iter.Value(), &batch)
		// cb returns true to stop early
		if cb(iter.Key(), &batch) {
			break
		}
	}
}

// GetOutgoingTxBatches returns the outgoing tx batches
func (k Keeper) GetOutgoingTxBatches(ctx sdk.Context) (out []*types.OutgoingTxBatch) {
	k.IterateOutgoingTxBatches(ctx, func(_ []byte, batch *types.OutgoingTxBatch) bool {
		out = append(out, batch)
		return false
	})
	return
}

// GetLastOutgoingBatchByTokenType gets the latest outgoing tx batch by token type
func (k Keeper) GetLastOutgoingBatchByTokenType(ctx sdk.Context, token string) *types.OutgoingTxBatch {
	batches := k.GetOutgoingTxBatches(ctx)
	var lastBatch *types.OutgoingTxBatch = nil
	lastNonce := uint64(0)
	for _, batch := range batches {
		if batch.TokenContract == token && batch.BatchNonce > lastNonce {
			lastBatch = batch
			lastNonce = batch.BatchNonce
		}
	}
	return lastBatch
}

// SetLastSlashedBatchBlock sets the latest slashed Batch block height
func (k Keeper) SetLastSlashedBatchBlock(ctx sdk.Context, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastSlashedBatchBlock, sdk.Uint64ToBigEndian(blockHeight))
}

// GetLastSlashedBatchBlock returns the latest slashed Batch block
func (k Keeper) GetLastSlashedBatchBlock(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastSlashedBatchBlock)

	if len(bytes) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(bytes)
}

// GetUnSlashedBatches returns all the unSlashed batches in state
func (k Keeper) GetUnSlashedBatches(ctx sdk.Context, maxHeight uint64) (outgoingTxBatches types.OutgoingTxBatches) {
	lastSlashedBatchBlock := k.GetLastSlashedBatchBlock(ctx)
	k.IterateBatchBySlashedBatchBlock(ctx, lastSlashedBatchBlock, maxHeight, func(_ []byte, batch *types.OutgoingTxBatch) bool {
		if batch.Block > lastSlashedBatchBlock {
			outgoingTxBatches = append(outgoingTxBatches, batch)
		}
		return false
	})
	sort.Sort(outgoingTxBatches)
	return
}

// IterateBatchBySlashedBatchBlock iterates through all Batch by last slashed Batch block in ASC order
func (k Keeper) IterateBatchBySlashedBatchBlock(ctx sdk.Context, lastSlashedBatchBlock uint64, maxHeight uint64, cb func([]byte, *types.OutgoingTxBatch) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OutgoingTxBatchBlockKey)
	iter := prefixStore.Iterator(sdk.Uint64ToBigEndian(lastSlashedBatchBlock), sdk.Uint64ToBigEndian(maxHeight))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var Batch types.OutgoingTxBatch
		k.cdc.MustUnmarshal(iter.Value(), &Batch)
		// cb returns true to stop early
		if cb(iter.Key(), &Batch) {
			break
		}
	}
}

/////////////////////////////
//      BATCH CONFIRMS     //
/////////////////////////////

// GetBatchConfirm returns a batch confirmation given its nonce, the token contract, and a oracle address
func (k Keeper) GetBatchConfirm(ctx sdk.Context, nonce uint64, tokenContract string, oracleAddr sdk.AccAddress) *types.MsgConfirmBatch {
	store := ctx.KVStore(k.storeKey)
	entity := store.Get(types.GetBatchConfirmKey(tokenContract, nonce, oracleAddr))
	if entity == nil {
		return nil
	}
	confirm := types.MsgConfirmBatch{}
	k.cdc.MustUnmarshal(entity, &confirm)
	return &confirm
}

// SetBatchConfirm sets a batch confirmation by a oracle
func (k Keeper) SetBatchConfirm(ctx sdk.Context, oracleAddr sdk.AccAddress, batch *types.MsgConfirmBatch) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetBatchConfirmKey(batch.TokenContract, batch.Nonce, oracleAddr)
	store.Set(key, k.cdc.MustMarshal(batch))
}

// IterateBatchConfirmByNonceAndTokenContract iterates through all batch confirmations
// MARK finish-batches: this is where the key is iterated in the old (presumed working) code
func (k Keeper) IterateBatchConfirmByNonceAndTokenContract(ctx sdk.Context, nonce uint64, tokenContract string, cb func([]byte, types.MsgConfirmBatch) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.BatchConfirmKey)
	prefixKey := append([]byte(tokenContract), sdk.Uint64ToBigEndian(nonce)...)
	iter := prefixStore.Iterator(prefixRange(prefixKey))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		confirm := types.MsgConfirmBatch{}
		k.cdc.MustUnmarshal(iter.Value(), &confirm)
		// cb returns true to stop early
		if cb(iter.Key(), confirm) {
			break
		}
	}
}

// GetBatchConfirmByNonceAndTokenContract returns the batch confirms
func (k Keeper) GetBatchConfirmByNonceAndTokenContract(ctx sdk.Context, nonce uint64, tokenContract string) (out []types.MsgConfirmBatch) {
	k.IterateBatchConfirmByNonceAndTokenContract(ctx, nonce, tokenContract, func(_ []byte, msg types.MsgConfirmBatch) bool {
		out = append(out, msg)
		return false
	})
	return
}
