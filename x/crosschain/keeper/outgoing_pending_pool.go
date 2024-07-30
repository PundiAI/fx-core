package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
)

func (k Keeper) AddToOutgoingPendingPool(ctx sdk.Context, sender sdk.AccAddress, receiver string, amount sdk.Coin, fee sdk.Coin) (uint64, error) {
	bridgeToken := k.GetDenomBridgeToken(ctx, amount.Denom)
	if bridgeToken == nil {
		return 0, errorsmod.Wrap(types.ErrInvalid, "bridge token is not exist")
	}

	// add pending pool switch
	if !k.GetParams(ctx).EnableSendToExternalPending {
		return 0, types.ErrInvalid.Wrapf("not enough liquidity")
	}

	nextTxID := k.autoIncrementID(ctx, types.KeyLastTxPoolID)

	pendingOutgoingTx := types.NewPendingOutgoingTx(nextTxID, sender, receiver, bridgeToken.Token, amount, fee, sdk.NewCoins())
	k.SetPendingTx(ctx, &pendingOutgoingTx)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSendToExternal,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyPendingOutgoingTxID, fmt.Sprint(nextTxID)),
	))
	return nextTxID, nil
}

func (k Keeper) RemovePendingOutgoingTx(context sdk.Context, tokenContract string, txId uint64) {
	store := context.KVStore(k.storeKey)
	store.Delete(types.GetOutgoingPendingTxPoolKey(tokenContract, txId))
}

func (k Keeper) HandleAddPendingPoolReward(ctx sdk.Context, id uint64, reward sdk.Coin) (success bool) {
	pendingPoolTx, found := k.GetPendingPoolTxById(ctx, id)
	if !found {
		return false
	}

	pendingPoolTx.Rewards = sdk.NewCoins(pendingPoolTx.GetRewards()...).Add(reward)
	k.SetPendingTx(ctx, pendingPoolTx)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeAddPendingRewards,
		sdk.NewAttribute(types.AttributeKeyPendingID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(types.AttributeKeyPendingRewards, reward.String()),
		sdk.NewAttribute(types.AttributeKeyPendingType, types.PendingTypeOutgoingTransferTx),
	))
	return true
}

func (k Keeper) HandlePendingOutgoingTx(ctx sdk.Context, liquidityProvider sdk.AccAddress, eventNonce uint64, bridgeToken *types.BridgeToken) {
	cacheContext, commit := ctx.CacheContext()

	erc20ModuleAddress := k.ak.GetModuleAddress(erc20types.ModuleName)
	var err error
	var txId uint64
	var provideLiquidityTxIds []uint64
	var rewards sdk.Coins
	liquidationSize := 0
	// iterator pending outgoing tx by bridgeToken contract address
	k.IteratorPendingOutgoingTxByBridgeTokenContractAddr(cacheContext, bridgeToken.Token, func(pendingOutgoingTx types.PendingOutgoingTransferTx) bool {
		// only allow to provide liquidity for MaxLiquidationSize times, avoid to exceed the limit
		liquidationSize++
		if liquidationSize >= types.MaxLiquidationSize {
			return true
		}
		// 1. check erc20 module has enough balance
		transferCoin := sdk.NewCoin(bridgeToken.Denom, pendingOutgoingTx.Token.Amount.Add(pendingOutgoingTx.Fee.Amount))
		if !k.bankKeeper.HasBalance(ctx, erc20ModuleAddress, transferCoin) {
			return false
		}

		// 2. transfer coin from erc20 module to sender
		sender := sdk.MustAccAddressFromBech32(pendingOutgoingTx.Sender)
		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, erc20types.ModuleName, sender, sdk.NewCoins(transferCoin)); err != nil {
			k.Logger(ctx).Info("failed to transfer coin from erc20 module to sender", "error", err)
			return true
		}

		// 3. remove pending outgoing tx
		k.RemovePendingOutgoingTx(cacheContext, bridgeToken.Token, pendingOutgoingTx.Id)

		// 4. add to outgoing tx
		if txId, err = k.AddToOutgoingPool(cacheContext, sender, pendingOutgoingTx.DestAddress, pendingOutgoingTx.Token, pendingOutgoingTx.Fee); err != nil {
			k.Logger(ctx).Info("failed to add to outgoing pool", "error", err)
			return true
		}
		provideLiquidityTxIds = append(provideLiquidityTxIds, txId)
		for _, reward := range pendingOutgoingTx.Rewards {
			rewards = rewards.Add(reward)
		}
		return false
	})

	if len(provideLiquidityTxIds) > 0 && err == nil {
		// 5. transfer rewards
		if !rewards.Empty() {
			if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, liquidityProvider, rewards); err != nil {
				k.Logger(ctx).Info("failed to transfer rewards", "error", err)
				return
			}

			for _, reward := range rewards {
				if _, err = k.erc20Keeper.ConvertDenomToTarget(ctx, liquidityProvider, reward, fxtypes.ParseFxTarget(fxtypes.ERC20Target)); err != nil {
					k.Logger(ctx).Info("failed to convert reward to target coin", "error", err)
					return
				}
			}
		}

		// 6. emit event & commit
		var eventIds string
		for _, id := range provideLiquidityTxIds {
			eventIds += fmt.Sprintf("%d,", id)
		}

		if len(eventIds) > 0 {
			eventIds = eventIds[:len(eventIds)-1]
			cacheContext.EventManager().EmitEvent(
				sdk.NewEvent(types.EventTypeProvideLiquidity,
					sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprintf("%d", eventNonce)),
					sdk.NewAttribute(types.AttributeKeyProvideLiquidityTxIds, eventIds),
				))
		}

		commit()
	}
}

func (k Keeper) handleRemoveFromOutgoingPendingPoolAndRefund(ctx sdk.Context, txId uint64, sender sdk.AccAddress) (sdk.Coin, error) {
	// 1. find pending outgoing tx by txId, and check sender
	tx, found := k.GetPendingPoolTxById(ctx, txId)
	if !found {
		return sdk.Coin{}, errorsmod.Wrap(types.ErrUnknown, "pool transaction")
	}

	txSender := sdk.MustAccAddressFromBech32(tx.Sender)
	if !txSender.Equals(sender) {
		return sdk.Coin{}, errorsmod.Wrapf(types.ErrInvalid, "Sender %s did not send Id %d", sender, txId)
	}

	// 2. delete pending outgoing tx
	k.RemovePendingOutgoingTx(ctx, tx.TokenContract, txId)

	// 3. refund rewards
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, txSender, tx.Rewards); err != nil {
		return sdk.Coin{}, err
	}

	// 4. refund token to sender
	return k.handleCancelPendingPoolRefund(ctx, txId, sender, tx.TokenContract, tx.Token.Amount.Add(tx.Fee.Amount))
}

func (k Keeper) SetPendingTx(ctx sdk.Context, outgoing *types.PendingOutgoingTransferTx) {
	store := ctx.KVStore(k.storeKey)
	idxKey := types.GetOutgoingPendingTxPoolKey(outgoing.TokenContract, outgoing.Id)
	store.Set(idxKey, k.cdc.MustMarshal(outgoing))
}

func (k Keeper) IteratorPendingOutgoingTxByBridgeTokenContractAddr(ctx sdk.Context, tokenContract string, cb func(pendingOutgoingTx types.PendingOutgoingTransferTx) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetOutgoingPendingTxPoolContractPrefix(tokenContract))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var pendingOutgoingTx types.PendingOutgoingTransferTx
		k.cdc.MustUnmarshal(iter.Value(), &pendingOutgoingTx)
		if cb(pendingOutgoingTx) {
			break
		}
	}
}

func (k Keeper) IteratorPendingOutgoingTx(ctx sdk.Context, cb func(pendingOutgoingTx types.PendingOutgoingTransferTx) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.PendingOutgoingTxPoolKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var pendingOutgoingTx types.PendingOutgoingTransferTx
		k.cdc.MustUnmarshal(iter.Value(), &pendingOutgoingTx)
		if cb(pendingOutgoingTx) {
			break
		}
	}
}

func (k Keeper) GetPendingPoolTxById(ctx sdk.Context, txId uint64) (*types.PendingOutgoingTransferTx, bool) {
	var tx types.PendingOutgoingTransferTx
	k.IteratorPendingOutgoingTx(ctx, func(pendingOutgoingTx types.PendingOutgoingTransferTx) bool {
		if pendingOutgoingTx.Id == txId {
			tx = pendingOutgoingTx
			return true
		}
		return false
	})
	return &tx, tx.Id == txId
}
