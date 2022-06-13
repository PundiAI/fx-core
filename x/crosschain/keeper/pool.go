package keeper

import (
	"encoding/binary"
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/crosschain/types"
)

// AddToOutgoingPool
// - checks a counterpart denominator exists for the given voucher type
// - burns the voucher for transfer amount and fees
// - persists an OutgoingTx
// - adds the TX to the `available` TX pool via a second index
func (k Keeper) AddToOutgoingPool(ctx sdk.Context, sender sdk.AccAddress, receiver string, amount sdk.Coin, fee sdk.Coin) (uint64, error) {
	totalAmount := amount.Add(fee)
	totalInVouchers := sdk.Coins{totalAmount}

	// If the coin is a gravity voucher, burn the coins. If not, check if there is a deployed ERC20 contract representing it.
	// If there is, lock the coins.

	bridgeToken := k.GetDenomByBridgeToken(ctx, amount.Denom)
	if bridgeToken == nil {
		return 0, sdkerrors.Wrap(types.ErrInvalid, "bridge token is not exist")
	}

	// If it is an ethereum-originated asset we burn it
	// send coins to module in prep for burn
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, k.moduleName, totalInVouchers); err != nil {
		return 0, err
	}

	// burn vouchers to send them back to ETH
	if err := k.bankKeeper.BurnCoins(ctx, k.moduleName, totalInVouchers); err != nil {
		return 0, err
	}

	// get next tx id from keeper
	nextTxID := k.autoIncrementID(ctx, types.KeyLastTxPoolID)

	tokenContract := bridgeToken.Token
	bridgeTokenFee := types.NewERC20Token(fee.Amount, tokenContract)

	// construct outgoing tx, as part of this process we represent
	// the token as an ERC20 token since it is preparing to go to ETH
	// rather than the denom that is the input to this function.
	outgoing := &types.OutgoingTransferTx{
		Id:          nextTxID,
		Sender:      sender.String(),
		DestAddress: receiver,
		Token:       types.NewERC20Token(amount.Amount, tokenContract),
		Fee:         bridgeTokenFee,
	}

	if err := k.addUnbatchedTX(ctx, outgoing); err != nil {
		return 0, nil
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSendToExternal,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyOutgoingTxID, fmt.Sprint(nextTxID)),
	))

	return nextTxID, nil
}

// RemoveFromOutgoingPoolAndRefund
// - checks that the provided tx actually exists
// - deletes the unbatched tx from the pool
// - issues the tokens back to the sender
func (k Keeper) RemoveFromOutgoingPoolAndRefund(ctx sdk.Context, txId uint64, sender sdk.AccAddress) error {
	if ctx.IsZero() || txId < 1 || sender.Empty() {
		return sdkerrors.Wrap(types.ErrInvalid, "arguments")
	}
	// check that we actually have a tx with that id and what it's details are
	tx, err := k.GetUnbatchedTxById(ctx, txId)
	if err != nil {
		return err
	}

	// Check that this user actually sent the transaction, this prevents someone from refunding someone
	// else transaction to themselves.
	txSender, err := sdk.AccAddressFromBech32(tx.Sender)
	if err != nil {
		panic("Invalid address in store!")
	}
	if !txSender.Equals(sender) {
		return sdkerrors.Wrapf(types.ErrInvalid, "Sender %s did not send Id %d", sender, txId)
	}

	// An inconsistent entry should never enter the store, but this is the ideal place to exploit
	// it such a bug if it did ever occur, so we should double check to be really sure
	if tx.Fee.Contract != tx.Token.Contract {
		return sdkerrors.Wrapf(types.ErrInvalid, "Inconsistent tokens to cancel!: %s %s", tx.Fee.Contract, tx.Token.Contract)
	}

	// delete this tx from the pool
	err = k.removeUnbatchedTX(ctx, tx.Fee, txId)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalid, "txId %d not in unbatched index! Must be in a batch!", txId)
	}
	// Make sure the tx was removed
	oldTx, oldTxErr := k.GetUnbatchedTxByFeeAndId(ctx, tx.Fee, tx.Id)
	if oldTx != nil || oldTxErr == nil {
		return sdkerrors.Wrapf(types.ErrInvalid, "tx with id %d was not fully removed from the pool, a duplicate must exist", txId)
	}

	// query denom, if not exist, return error
	bridgeToken := k.GetBridgeTokenDenom(ctx, tx.Token.Contract)
	if bridgeToken == nil {
		return sdkerrors.Wrapf(types.ErrInvalid, "Invalid token, contract %s", tx.Token.Contract)
	}
	// reissue the amount and the fee
	totalToRefund := sdk.NewCoin(bridgeToken.Denom, tx.Token.Amount)
	totalToRefund.Amount = totalToRefund.Amount.Add(tx.Fee.Amount)
	totalToRefundCoins := sdk.NewCoins(totalToRefund)

	// If it is an ethereum-originated asset we have to mint it (see Handle in attestation_handler.go)
	// mint coins in module for prep to send
	if err = k.bankKeeper.MintCoins(ctx, k.moduleName, totalToRefundCoins); err != nil {
		return sdkerrors.Wrapf(err, "mint vouchers coins: %s", totalToRefundCoins)
	}
	if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, sender, totalToRefundCoins); err != nil {
		return sdkerrors.Wrap(err, "transfer vouchers")
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSendToExternalCanceled,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyOutgoingTxID, fmt.Sprint(txId)),
	))
	return nil
}

// addUnbatchedTx creates a new transaction in the pool
// WARNING: Do not make this function public
func (k Keeper) addUnbatchedTX(ctx sdk.Context, val *types.OutgoingTransferTx) error {
	store := ctx.KVStore(k.storeKey)
	idxKey := types.GetOutgoingTxPoolKey(val.Fee, val.Id)
	if store.Has(idxKey) {
		return sdkerrors.Wrap(types.ErrDuplicate, "transaction already in pool")
	}

	bz, err := k.cdc.Marshal(val)
	if err != nil {
		return err
	}

	store.Set(idxKey, bz)
	return err
}

// removeUnbatchedTXIndex removes the tx from the pool
// WARNING: Do not make this function public
func (k Keeper) removeUnbatchedTX(ctx sdk.Context, fee types.ERC20Token, txID uint64) error {
	store := ctx.KVStore(k.storeKey)
	idxKey := types.GetOutgoingTxPoolKey(fee, txID)
	if !store.Has(idxKey) {
		return sdkerrors.Wrap(types.ErrUnknown, "pool transaction")
	}
	store.Delete(idxKey)
	return nil
}

// GetUnbatchedTxByFeeAndId grabs a tx from the pool given its fee and txID
func (k Keeper) GetUnbatchedTxByFeeAndId(ctx sdk.Context, fee types.ERC20Token, txID uint64) (*types.OutgoingTransferTx, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOutgoingTxPoolKey(fee, txID))
	if bz == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "pool transaction")
	}
	var r types.OutgoingTransferTx
	err := k.cdc.Unmarshal(bz, &r)
	return &r, err
}

// GetUnbatchedTxById grabs a tx from the pool given only the txID
// note that due to the way unbatched txs are indexed, the GetUnbatchedTxByFeeAndId method is much faster
func (k Keeper) GetUnbatchedTxById(ctx sdk.Context, txID uint64) (*types.OutgoingTransferTx, error) {
	var r *types.OutgoingTransferTx = nil
	k.IterateUnbatchedTransactions(ctx, types.OutgoingTxPoolKey, func(_ []byte, tx *types.OutgoingTransferTx) bool {
		if tx.Id == txID {
			r = tx
			return true
		}
		return false // iterating DESC, exit early
	})

	if r == nil {
		// We have no return tx, it was either batched or never existed
		return nil, sdkerrors.Wrap(types.ErrUnknown, "pool transaction")
	}
	return r, nil
}

// GetUnbatchedTransactionsByContract grabs all unbatched transactions from the tx pool for the given contract
// unbatched transactions are sorted by fee amount in DESC order
func (k Keeper) GetUnbatchedTransactionsByContract(ctx sdk.Context, contractAddress string) []*types.OutgoingTransferTx {
	return k.collectUnbatchedTransactions(ctx, types.GetOutgoingTxPoolContractPrefix(contractAddress))
}

// GetUnbatchedTransactions grabs all transactions from the tx pool, useful for queries or genesis save/load
func (k Keeper) GetUnbatchedTransactions(ctx sdk.Context) []*types.OutgoingTransferTx {
	return k.collectUnbatchedTransactions(ctx, types.OutgoingTxPoolKey)
}

// Aggregates all unbatched transactions in the store with a given prefix
func (k Keeper) collectUnbatchedTransactions(ctx sdk.Context, prefixKey []byte) (out []*types.OutgoingTransferTx) {
	k.IterateUnbatchedTransactions(ctx, prefixKey, func(_ []byte, tx *types.OutgoingTransferTx) bool {
		out = append(out, tx)
		return false
	})
	return
}

// IterateUnbatchedTransactionsByContract iterates through unbatched transactions from the tx pool for the given contract
// unbatched transactions are sorted by fee amount in DESC order
func (k Keeper) IterateUnbatchedTransactionsByContract(ctx sdk.Context, contractAddress string, cb func(key []byte, tx *types.OutgoingTransferTx) bool) {
	k.IterateUnbatchedTransactions(ctx, types.GetOutgoingTxPoolContractPrefix(contractAddress), cb)
}

// IterateUnbatchedTransactions iterates through all unbatched transactions whose keys begin with prefixKey in DESC order
func (k Keeper) IterateUnbatchedTransactions(ctx sdk.Context, prefixKey []byte, cb func(key []byte, tx *types.OutgoingTransferTx) bool) {
	prefixStore := ctx.KVStore(k.storeKey)
	iter := prefixStore.ReverseIterator(prefixRange(prefixKey))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var transact types.OutgoingTransferTx
		k.cdc.MustUnmarshal(iter.Value(), &transact)
		// cb returns true to stop early
		if cb(iter.Key(), &transact) {
			break
		}
	}
}

// GetBatchFeesByTokenType gets the fee the next batch of a given token type would
// have if created right now. This info is both presented to relayers for the purpose of determining
// when to request batches and also used by the batch creation process to decide not to create
// a new batch (fees must be increasing)
func (k Keeper) GetBatchFeesByTokenType(ctx sdk.Context, tokenContractAddr string, maxElements uint, baseFee sdk.Int) *types.BatchFees {
	batchFee := &types.BatchFees{TokenContract: tokenContractAddr, TotalFees: sdk.NewInt(0), TotalAmount: sdk.NewInt(0)}

	k.IterateUnbatchedTransactions(ctx, types.GetOutgoingTxPoolContractPrefix(tokenContractAddr), func(_ []byte, tx *types.OutgoingTransferTx) bool {
		if tx.Fee.Contract != tokenContractAddr {
			panic(fmt.Errorf("unexpected fee contract %s when getting batch fees for contract %s", tx.Fee.Contract, tokenContractAddr))
		}
		if tx.Fee.Amount.LT(baseFee) {
			// sort by fee and use ReverseIterator, so the fee behind is less than base fee
			return true
		}
		batchFee.TotalFees = batchFee.TotalFees.Add(tx.Fee.Amount)
		batchFee.TotalAmount = batchFee.TotalAmount.Add(tx.Token.Amount)
		batchFee.TotalTxs += 1
		return batchFee.TotalTxs == uint64(maxElements)
	})
	return batchFee
}

// GetAllBatchFees creates a fee entry for every batch type currently in the store
// this can be used by relayers to determine what batch types are desirable to request
func (k Keeper) GetAllBatchFees(ctx sdk.Context, maxElements uint, minBatchFees []types.MinBatchFee) (batchFees []*types.BatchFees) {
	batchFeesMap := k.createBatchFees(ctx, maxElements, minBatchFees)
	// create array of batchFees
	for _, batchFee := range batchFeesMap {
		batchFees = append(batchFees, batchFee)
	}

	// quick sort by token to make this function safe for use
	// in consensus computations
	sort.Slice(batchFees, func(i, j int) bool {
		return batchFees[i].TokenContract < batchFees[j].TokenContract
	})

	return batchFees
}

// createBatchFees iterates over the unbatched transaction pool and creates batch token fee map
// Implicitly creates batches with the highest potential fee because the transaction keys enforce an order which goes
// fee contract address -> fee amount -> transaction nonce
func (k Keeper) createBatchFees(ctx sdk.Context, maxElements uint, minBatchFees []types.MinBatchFee) map[string]*types.BatchFees {
	batchFeesMap := make(map[string]*types.BatchFees)
	txCountMap := make(map[string]int)
	baseFees := types.MinBatchFeeToBaseFees(minBatchFees)

	k.IterateUnbatchedTransactions(ctx, types.OutgoingTxPoolKey, func(_ []byte, tx *types.OutgoingTransferTx) bool {
		fee := tx.Fee

		baseFee, ok := baseFees[fee.Contract]
		if ok && fee.Amount.LT(baseFee) {
			return false //sort by token address and fee, behind have other token
		}

		if txCountMap[fee.Contract] < int(maxElements) {
			addFeeToMap(tx.Token, fee, batchFeesMap, txCountMap)
		}
		return false
	})

	return batchFeesMap
}

func (k Keeper) autoIncrementID(ctx sdk.Context, idKey []byte) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(idKey)
	var id uint64 = 1
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	bz = sdk.Uint64ToBigEndian(id + 1)
	store.Set(idKey, bz)
	return id
}

// Helper method for creating batch fees
func addFeeToMap(amt, fee types.ERC20Token, batchFeesMap map[string]*types.BatchFees, txCountMap map[string]int) {
	txCountMap[fee.Contract] = txCountMap[fee.Contract] + 1

	// add fee amount
	if _, ok := batchFeesMap[fee.Contract]; ok {
		batchFees := batchFeesMap[fee.Contract]
		batchFees.TotalFees = batchFees.TotalFees.Add(fee.Amount)
		batchFees.TotalTxs = batchFees.TotalTxs + 1
		batchFees.TotalAmount = batchFees.TotalAmount.Add(amt.Amount)
		batchFeesMap[fee.Contract] = batchFees
	} else {
		batchFeesMap[fee.Contract] = &types.BatchFees{
			TokenContract: fee.Contract,
			TotalFees:     fee.Amount,
			TotalTxs:      1,
			TotalAmount:   amt.Amount,
		}
	}
}
