package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"

	"github.com/functionx/fx-core/x/crosschain/types"
)

// SetOracle save Oracle data
func (k Keeper) SetOracle(ctx sdk.Context, oracle types.Oracle) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&oracle)
	store.Set(types.GetOracleKey(oracle.GetOracle()), bz)
}

// GetOracle get Oracle data
func (k Keeper) GetOracle(ctx sdk.Context, addr sdk.AccAddress) (oracle types.Oracle, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetOracleKey(addr))
	if value == nil {
		return oracle, false
	}
	k.cdc.MustUnmarshal(value, &oracle)
	return oracle, true
}

func (k Keeper) DelOracle(ctx sdk.Context, oracle sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOracleKey(oracle)
	if !store.Has(key) {
		return
	}
	store.Delete(key)
}

// GetAllOracles
func (k Keeper) GetAllOracles(ctx sdk.Context) (oracles types.Oracles) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OracleKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var oracle types.Oracle
		k.cdc.MustUnmarshal(iterator.Value(), &oracle)
		oracles = append(oracles, oracle)
	}
	sort.Sort(oracles)
	return oracles
}

func (k Keeper) GetAllActiveOracles(ctx sdk.Context) (oracles types.Oracles) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OracleKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var oracle types.Oracle
		k.cdc.MustUnmarshal(iterator.Value(), &oracle)
		if oracle.Jailed {
			continue
		}
		oracles = append(oracles, oracle)
	}
	sort.Sort(oracles)
	return oracles
}

func (k Keeper) SlashOracle(ctx sdk.Context, oracle types.Oracle, slashFraction sdk.Dec) {
	if oracle.Jailed {
		return
	}
	slashAmount := oracle.DelegateAmount.Amount.ToDec().Mul(slashFraction).TruncateInt()
	tokensToBurn := sdk.MinInt(slashAmount, oracle.DelegateAmount.Amount)
	tokensToBurn = sdk.MaxInt(tokensToBurn, sdk.ZeroInt())
	slashCoin := sdk.NewCoin(oracle.DelegateAmount.Denom, tokensToBurn)
	if slashCoin.IsPositive() {
		oracle.DelegateAmount = oracle.DelegateAmount.Sub(slashCoin)
		//if err = k.bankKeeper.BurnCoins(ctx, k.moduleName, sdk.NewCoins(slashCoin)); err != nil {
		//	panic(err)
		//}

		//ctx.EventManager().EmitEvents(sdk.Events{
		//	sdk.NewEvent(
		//		stakingtypes.EventTypeUnbond,
		//		sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValidatorAddress),
		//		sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		//		sdk.NewAttribute(stakingtypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		//	),
		//})
	}

	oracle.Jailed = true
	oracle.JailedHeight = ctx.BlockHeight()
	k.SetOracle(ctx, oracle)
	k.SetLastOracleSlashBlockHeight(ctx, uint64(ctx.BlockHeight()))
}
