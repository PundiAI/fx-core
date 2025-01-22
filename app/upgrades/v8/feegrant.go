package v8

import (
	"cosmossdk.io/core/store"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	types1 "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func MigrateFeegrant(ctx sdk.Context, cdc codec.BinaryCodec, storeService store.KVStoreService, authKeeper authkeeper.AccountKeeper) error {
	var err error
	keeper := feegrantkeeper.NewKeeper(cdc, storeService, authKeeper)
	iterErr := keeper.IterateAllFeeAllowances(ctx, func(grant feegrant.Grant) bool {
		grant.Allowance, err = swapAllowance(grant.Allowance)
		if err != nil {
			return true
		}
		err = updateAllowance(ctx, cdc, storeService, authKeeper, grant)
		return err != nil
	})
	if iterErr != nil {
		return iterErr
	}
	if err != nil {
		return err
	}
	return err
}

func swapAllowance(any *types1.Any) (*types1.Any, error) {
	switch allowance := any.GetCachedValue().(type) {
	case *feegrant.BasicAllowance:
		allowance.SpendLimit = fxtypes.SwapCoins(allowance.SpendLimit)
		return types1.NewAnyWithValue(allowance)
	case *feegrant.PeriodicAllowance:
		allowance.Basic.SpendLimit = fxtypes.SwapCoins(allowance.Basic.SpendLimit)
		allowance.PeriodCanSpend = fxtypes.SwapCoins(allowance.PeriodCanSpend)
		allowance.PeriodSpendLimit = fxtypes.SwapCoins(allowance.PeriodSpendLimit)
		return types1.NewAnyWithValue(allowance)
	case *feegrant.AllowedMsgAllowance:
		newAny, err := swapAllowance(allowance.Allowance)
		if err != nil {
			return nil, err
		}
		allowance.Allowance = newAny
		return types1.NewAnyWithValue(allowance)
	default:
		return nil, sdkerrors.ErrPackAny.Wrap("failed to swap allowance")
	}
}

func updateAllowance(ctx sdk.Context, cdc codec.BinaryCodec, storeService store.KVStoreService, authKeeper authkeeper.AccountKeeper, grant feegrant.Grant) error {
	kvStore := storeService.OpenKVStore(ctx)

	granter, err := authKeeper.AddressCodec().StringToBytes(grant.Granter)
	if err != nil {
		return err
	}
	grantee, err := authKeeper.AddressCodec().StringToBytes(grant.Grantee)
	if err != nil {
		return err
	}
	key := feegrant.FeeAllowanceKey(granter, grantee)

	bz, err := cdc.Marshal(&grant)
	if err != nil {
		return err
	}
	return kvStore.Set(key, bz)
}
