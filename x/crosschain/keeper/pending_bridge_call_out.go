package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/cosmos/gogoproto/types"
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
		NotLiquidCoins:    notLiquidCoins,
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
		if !found {
			k.Logger(iterCtx).Error("no pending bridge call found", "nonce", bridgeCallNonce)
			return false
		}

		// 3. transfer coin from erc20 module to sender
		bridgeCall := pendingBridgeCall.OutgoinBridgeCall
		sender := types.ExternalAddrToAccAddr(k.moduleName, bridgeCall.Sender)
		if err = k.bankKeeper.SendCoinsFromModuleToAccount(iterCtx, erc20types.ModuleName, sender, notLiquidCoins); err != nil {
			k.Logger(iterCtx).Info("failed to transfer coin from erc20 module to sender", "error", err)
			return true
		}

		for _, coin := range notLiquidCoins {
			if err = k.TransferBridgeCoinToExternal(iterCtx, sender, coin); err != nil {
				k.Logger(iterCtx).Info("failed to transfer bridge coin to external", "error", err)
				return true
			}
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

func (k Keeper) HandleCancelPendingOutgoingBridgeCall(ctx sdk.Context, nonce uint64, sender sdk.AccAddress) (sdk.Coins, error) {
	pendingOutCall, found := k.GetPendingOutgoingBridgeCallByNonce(ctx, nonce)
	if !found {
		return nil, types.ErrInvalid.Wrapf("not found, nonce: %d", nonce)
	}

	outCall := pendingOutCall.OutgoinBridgeCall
	outCallSender := types.ExternalAddrToAccAddr(k.moduleName, outCall.Sender)
	if !sender.Equals(outCallSender) {
		return nil, types.ErrInvalid.Wrapf("msg.sender %s is not bridge call sender %s", sender.String(), outCallSender.String())
	}

	refundCoins := sdk.NewCoins()
	// 1. reuse refund logic
	notLiquidTargetCoins := sdk.NewCoins()
	for _, coin := range pendingOutCall.NotLiquidCoins {
		notLiquidTargetCoin, err := k.erc20Keeper.RefundLiquidity(ctx, outCallSender, coin)
		if err != nil {
			return nil, types.ErrInvalid.Wrapf("refund liquidity failed, error: %s", err)
		}
		notLiquidTargetCoins = notLiquidTargetCoins.Add(notLiquidTargetCoin)
		bridgeToken := k.GetDenomBridgeToken(ctx, coin.GetDenom())
		if bridgeToken == nil {
			return nil, types.ErrInvalid.Wrapf("bridge token not found, denom: %s", coin.GetDenom())
		}
		for i := 0; i < len(outCall.Tokens); i++ {
			if outCall.Tokens[i].Contract == bridgeToken.Token {
				outCall.Tokens = append(outCall.Tokens[:i], outCall.Tokens[i+1:]...)
				break
			}
		}
	}

	refundCoins = refundCoins.Add(notLiquidTargetCoins...)

	if !notLiquidTargetCoins.IsZero() && !k.HasBridgeCallFromMsg(ctx, nonce) {
		if err := k.bridgeCallTransferTokens(ctx, outCallSender, outCallSender, notLiquidTargetCoins); err != nil {
			panic(err)
		}
	}

	outCall.Refund = outCall.Sender
	coins := k.HandleOutgoingBridgeCallRefund(ctx, outCall)
	refundCoins = refundCoins.Add(coins...)

	// 2. refund rewards
	if !pendingOutCall.Rewards.IsZero() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, sender, pendingOutCall.Rewards); err != nil {
			return nil, err
		}
	}
	refundCoins = refundCoins.Add(pendingOutCall.Rewards...)

	// 3. delete pending outgoing bridge call
	k.DeletePendingOutgoingBridgeCall(ctx, nonce)

	// 4. delete bridge call from msg
	k.DeleteBridgeCallFromMsg(ctx, nonce)
	return refundCoins, nil
}

func (k Keeper) AddPendingPoolRewards(ctx sdk.Context, nonce uint64, sender sdk.AccAddress, rewards sdk.Coins) error {
	// 0. validate rewards coin, only support stake coin.
	reward, err := types.RewardValidator(rewards)
	if err != nil {
		return err
	}

	// 1. transfer coins to module
	if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, k.moduleName, sdk.NewCoins(reward)); err != nil {
		return err
	}
	// 2. try to handle adding pending bridge call rewards
	if addSuccess := k.HandleAddPendingBridgeCallRewards(ctx, nonce, reward); addSuccess {
		return nil
	}
	// 3. try to handle adding pending pool rewards
	if addSuccess := k.HandleAddPendingPoolReward(ctx, nonce, reward); addSuccess {
		return nil
	}
	return errors.ErrInvalidRequest.Wrap("not found pending record")
}

func (k Keeper) HandleAddPendingBridgeCallRewards(ctx sdk.Context, nonce uint64, reward sdk.Coin) (success bool) {
	// 1. find the pending outgoing bridge call by nonce
	pendingBridgeCall, found := k.GetPendingOutgoingBridgeCallByNonce(ctx, nonce)
	if !found {
		return false
	}

	// 3. update rewards
	pendingBridgeCall.Rewards = sdk.NewCoins(pendingBridgeCall.GetRewards()...).Add(reward)
	k.SetPendingOutgoingBridgeCallWithoutNotLiquidCoins(ctx, pendingBridgeCall)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeAddPendingRewards,
		sdk.NewAttribute(types.AttributeKeyPendingID, fmt.Sprintf("%d", nonce)),
		sdk.NewAttribute(types.AttributeKeyPendingRewards, reward.String()),
		sdk.NewAttribute(types.AttributeKeyPendingType, types.PendingTypeOutgoingBridgeCall),
	))
	return true
}

func (k Keeper) transferLiquidityProviderRewards(ctx sdk.Context, liquidityProvider []byte, eventNonce uint64, rewards sdk.Coins, provideLiquidityBridgeCallNonces []uint64) error {
	// transfer rewards
	if !rewards.Empty() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, liquidityProvider, rewards); err != nil {
			k.Logger(ctx).Info("failed to transfer rewards", "error", err)
			return err
		}

		for _, reward := range rewards {
			if reward.Denom == fxtypes.DefaultDenom {
				continue
			}
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
	k.SetPendingOutgoingBridgeCallWithoutNotLiquidCoins(ctx, pendingOutCall)
	store := ctx.KVStore(k.storeKey)
	nonce := pendingOutCall.OutgoinBridgeCall.Nonce

	notLiquidBz := []byte(pendingOutCall.GetNotLiquidCoins().String())
	for _, coin := range pendingOutCall.GetNotLiquidCoins() {
		store.Set(types.GetNotLiquidCoinWithIdKey(coin.Denom, nonce), notLiquidBz)
	}
}

func (k Keeper) SetPendingOutgoingBridgeCallWithoutNotLiquidCoins(ctx sdk.Context, pendingOutCall *types.PendingOutgoingBridgeCall) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetPendingOutgoingBridgeCallNonceKey(pendingOutCall.OutgoinBridgeCall.Nonce), k.cdc.MustMarshal(pendingOutCall))
	store.Set(
		types.GetPendingOutgoingBridgeCallAddressAndNonceKey(pendingOutCall.OutgoinBridgeCall.Sender, pendingOutCall.OutgoinBridgeCall.Nonce),
		k.cdc.MustMarshal(&gogotypes.BoolValue{Value: true}),
	)
}

func (k Keeper) DeletePendingOutgoingBridgeCall(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	pendingOutCall, found := k.GetPendingOutgoingBridgeCallByNonce(ctx, nonce)
	if !found {
		return
	}
	outCallKey := types.GetPendingOutgoingBridgeCallNonceKey(nonce)
	store.Delete(outCallKey)

	store.Delete(types.GetPendingOutgoingBridgeCallAddressAndNonceKey(pendingOutCall.OutgoinBridgeCall.Sender, nonce))

	for _, coin := range pendingOutCall.GetNotLiquidCoins() {
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
		coins, err := sdk.ParseCoinsNormalized(string(iter.Value()))
		if err != nil {
			break
		}
		if cb(bridgeCallNonce, coins) {
			break
		}
	}
}

func (k Keeper) IteratePendingOutgoingBridgeCallsByAddress(ctx sdk.Context, senderAddr string, cb func(outCall *types.PendingOutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetPendingOutgoingBridgeCallAddressKey(senderAddr))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		nonce := types.ParsePendingOutgoingBridgeCallNonce(iterator.Key(), senderAddr)
		outCall, found := k.GetPendingOutgoingBridgeCallByNonce(ctx, nonce)
		if !found {
			continue
		}
		if cb(outCall) {
			break
		}
	}
}
