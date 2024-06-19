package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
)

func (k Keeper) AddPendingOutgoingBridgeCall(ctx sdk.Context, sender, refundAddr common.Address, tokens []types.ERC20Token, to common.Address, data, memo []byte, eventNonce uint64, notLiquidCoins sdk.Coins) (uint64, error) {
	// try to calculate the bridge call timeout height, Avoid failure to calculate timeout when liquidity is sufficient
	outCall, err := k.BuildOutgoingBridgeCall(ctx, sender, refundAddr, tokens, to, data, memo, eventNonce)
	if err != nil {
		return 0, err
	}
	outCall.Timeout = 0

	pendingOutCall := &types.PendingOutgoingBridgeCall{
		OutgoinBridgeCall: outCall,
		NotLiquidCoins:    types.NotLiquidCoins{NotLiquidCoins: notLiquidCoins},
		Rewards:           sdk.NewCoins(),
	}
	k.SetPendingOutgoingBridgeCall(ctx, pendingOutCall)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypePendingBridgeCall,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, outCall.Sender),
		sdk.NewAttribute(types.AttributeKeyBridgeCallNonce, fmt.Sprint(outCall.Nonce)),
	))

	return outCall.Nonce, nil
}

func (k Keeper) HandlePendingOutgoingBridgeCall(ctx sdk.Context, liquidityProvider []byte, eventNonce uint64, bridgeToken *types.BridgeToken) {
	cacheContext, commit := ctx.CacheContext()

	erc20ModuleAddress := k.ak.GetModuleAddress(erc20types.ModuleName)
	var err error
	var provideLiquidityBridgeCallNonces []uint64
	var rewards sdk.Coins
	liquidationSize := 0
	// iterator pending outgoing tx by bridgeToken contract address
	k.IteratorBridgeCallNotLiquidsByDenom(cacheContext, bridgeToken.Denom, func(bridgeCallNonce uint64, notLiquidCoins sdk.Coins) bool {
		iterCtx, iterCommit := cacheContext.CacheContext()
		// only allow to provide liquidity for MaxLiquidationSize times, avoid to exceed the limit
		liquidationSize++
		if liquidationSize >= types.MaxLiquidationSize {
			return true
		}
		// 1. check erc20 module has enough balance
		for _, coin := range notLiquidCoins {
			if !k.bankKeeper.HasBalance(iterCtx, erc20ModuleAddress, coin) {
				return false
			}
		}

		// 2. check bridgeCall notLiquidCoins has balances
		pendingBridgeCall, found := k.GetPendingOutgoingBridgeCallByNonce(iterCtx, bridgeCallNonce)
		if found {
			k.Logger(iterCtx).Error("no pending bridge call found", "nonce", bridgeCallNonce)
			return false
		}

		// 3. transfer coin from erc20 module to sender
		bridgeCall := pendingBridgeCall.OutgoinBridgeCall
		sender := sdk.MustAccAddressFromBech32(bridgeCall.Sender)
		if err = k.bankKeeper.SendCoinsFromModuleToAccount(iterCtx, erc20types.ModuleName, sender, notLiquidCoins); err != nil {
			k.Logger(iterCtx).Info("failed to transfer coin from erc20 module to sender", "error", err)
			return true
		}

		// 4. remove pending outgoing tx
		k.DeletePendingOutgoingBridgeCall(iterCtx, bridgeCall.Nonce)

		// 5. add to outgoing bridge call
		bridgeCallTimeout := k.CalExternalTimeoutHeight(iterCtx, GetBridgeCallTimeout)
		if bridgeCallTimeout <= 0 {
			k.Logger(iterCtx).Error("failed calc bridge call external timeout height", "err", err, "nonce", bridgeCall.Nonce)
			return true
		}
		bridgeCall.Timeout = bridgeCallTimeout
		k.AddOutgoingBridgeCallWithoutBuild(iterCtx, bridgeCall)

		// 6. rewards
		provideLiquidityBridgeCallNonces = append(provideLiquidityBridgeCallNonces, bridgeCall.Nonce)
		for _, reward := range pendingBridgeCall.Rewards {
			rewards = rewards.Add(reward)
		}
		iterCommit()
		return false
	})

	if len(provideLiquidityBridgeCallNonces) > 0 && err == nil {
		if err = k.transferLiquidityProviderRewards(cacheContext, liquidityProvider, eventNonce, rewards, provideLiquidityBridgeCallNonces); err != nil {
			return
		}
	}

	if err == nil {
		commit()
	}
}

func (k Keeper) transferLiquidityProviderRewards(ctx sdk.Context, liquidityProvider []byte, eventNonce uint64, rewards sdk.Coins, provideLiquidityBridgeCallNonces []uint64) error {
	// transfer rewards
	if !rewards.Empty() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, liquidityProvider, rewards); err != nil {
			k.Logger(ctx).Info("failed to transfer rewards", "error", err)
			return err
		}

		for _, reward := range rewards {
			if _, err := k.erc20Keeper.ConvertDenomToTarget(ctx, liquidityProvider, reward, fxtypes.ParseFxTarget(fxtypes.ERC20Target)); err != nil {
				k.Logger(ctx).Info("failed to convert reward to target coin", "error", err)
				return err
			}
		}
	}

	// emit event & commit
	var eventIds string
	for _, id := range provideLiquidityBridgeCallNonces {
		eventIds += fmt.Sprintf("%d,", id)
	}

	if len(eventIds) > 0 {
		eventIds = eventIds[:len(eventIds)-1]
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventTypeProvideLiquidity,
				sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprintf("%d", eventNonce)),
				sdk.NewAttribute(types.AttributeKeyProvideLiquidityTxIds, eventIds),
			))
	}
	return nil
}

func (k Keeper) SetPendingOutgoingBridgeCall(ctx sdk.Context, pendingOutCall *types.PendingOutgoingBridgeCall) {
	store := ctx.KVStore(k.storeKey)
	nonce := pendingOutCall.OutgoinBridgeCall.Nonce
	store.Set(types.GetPendingOutgoingBridgeCallNonceKey(nonce), k.cdc.MustMarshal(pendingOutCall))

	notLiquidCoinsBytes := k.cdc.MustMarshal(&pendingOutCall.NotLiquidCoins)
	for _, coin := range pendingOutCall.NotLiquidCoins.GetNotLiquidCoins() {
		store.Set(types.GetNotLiquidCoinWithIdKey(coin.Denom, nonce), notLiquidCoinsBytes)
	}
}

func (k Keeper) DeletePendingOutgoingBridgeCall(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	pendingOutCall, found := k.GetPendingOutgoingBridgeCallByNonce(ctx, nonce)
	if !found {
		return
	}
	outCallKey := types.GetPendingOutgoingBridgeCallNonceKey(nonce)
	store.Delete(outCallKey)

	for _, coin := range pendingOutCall.NotLiquidCoins.GetNotLiquidCoins() {
		store.Delete(types.GetNotLiquidCoinWithIdKey(coin.Denom, nonce))
	}
}

func (k Keeper) GetPendingOutgoingBridgeCallByNonce(ctx sdk.Context, nonce uint64) (*types.PendingOutgoingBridgeCall, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPendingOutgoingBridgeCallNonceKey(nonce))
	if bz == nil {
		return nil, false
	}
	var outCall types.PendingOutgoingBridgeCall
	k.cdc.MustUnmarshal(bz, &outCall)
	return &outCall, true
}

func (k Keeper) IteratorBridgeCallNotLiquidsByDenom(ctx sdk.Context, denom string,
	cb func(bridgeCallNonce uint64, notLiquidCoins sdk.Coins) bool,
) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetNotLiquidCoinKey(denom))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		bridgeCallNonce := types.ParseBridgeCallNotLiquidNonce(iter.Key(), denom)
		var notLiquidCoins types.NotLiquidCoins
		k.cdc.MustUnmarshal(iter.Value(), &notLiquidCoins)
		if cb(bridgeCallNonce, notLiquidCoins.NotLiquidCoins) {
			break
		}
	}
}
