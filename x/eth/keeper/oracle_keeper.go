package keeper

import (
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
	"github.com/functionx/fx-core/x/eth/types"
)

type ValOracleKeeper struct {
	cdc           codec.BinaryCodec
	storeKey      sdk.StoreKey
	stakingKeeper types.StakingKeeper
}

func NewCrossOracleKeeper() {

}

// SetOracle save Oracle data
func (k ValOracleKeeper) SetOracle(ctx sdk.Context, oracle crosschaintypes.Oracle) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&oracle)
	store.Set(crosschaintypes.GetOracleKey(oracle.GetOracle()), bz)
}

// GetOracle get Oracle data
func (k ValOracleKeeper) GetOracle(ctx sdk.Context, addr sdk.AccAddress) (oracle crosschaintypes.Oracle, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(crosschaintypes.GetOracleKey(addr))
	if value == nil {
		return oracle, false
	}
	k.cdc.MustUnmarshal(value, &oracle)
	return oracle, true
}

func (k ValOracleKeeper) DelOracle(ctx sdk.Context, oracle sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := crosschaintypes.GetOracleKey(oracle)
	if !store.Has(key) {
		return
	}
	store.Delete(key)
}

// GetAllOracles
func (k ValOracleKeeper) GetAllOracles(ctx sdk.Context) (oracles crosschaintypes.Oracles) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, crosschaintypes.OracleKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var oracle crosschaintypes.Oracle
		k.cdc.MustUnmarshal(iterator.Value(), &oracle)
		oracles = append(oracles, oracle)
	}
	sort.Sort(oracles)
	return oracles
}

func (k ValOracleKeeper) GetAllActiveOracles(ctx sdk.Context) (oracles crosschaintypes.Oracles) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, crosschaintypes.OracleKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var oracle crosschaintypes.Oracle
		k.cdc.MustUnmarshal(iterator.Value(), &oracle)
		if oracle.Jailed {
			continue
		}
		oracles = append(oracles, oracle)
	}
	sort.Sort(oracles)
	return oracles
}

func (k ValOracleKeeper) SlashOracle(ctx sdk.Context, oracleAddress string, slashFraction sdk.Dec) {
	oracleAddr, err := sdk.AccAddressFromBech32(oracleAddress)
	if err != nil {
		panic(err)
	}
	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		panic(fmt.Sprintf("not found oracle:%s", oracleAddr.String()))
	}
	if oracle.Jailed {
		return
	}
	//slashAmount := oracle.DepositAmount.Amount.ToDec().Mul(slashFraction).TruncateInt()
	//tokensToBurn := sdk.MinInt(slashAmount, oracle.DepositAmount.Amount)
	//tokensToBurn = sdk.MaxInt(tokensToBurn, sdk.ZeroInt())
	//slashCoin := sdk.NewCoin(oracle.DepositAmount.Denom, tokensToBurn)
	//if slashCoin.IsPositive() {
	//	oracle.DepositAmount = oracle.DepositAmount.Sub(slashCoin)
	//	if err = k.bankKeeper.BurnCoins(ctx, k.moduleName, sdk.NewCoins(slashCoin)); err != nil {
	//		panic(err)
	//	}
	//}
	//
	//oracle.Jailed = true
	//oracle.JailedHeight = ctx.BlockHeight()
	//k.SetOracle(ctx, oracle)
	//k.SetLastOracleSlashBlockHeight(ctx, uint64(ctx.BlockHeight()))
}
