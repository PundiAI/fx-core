package v8

import (
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/pundiai/fx-core/v8/app/upgrades/store"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	fxgovkeeper "github.com/pundiai/fx-core/v8/x/gov/keeper"
	fxgovv8 "github.com/pundiai/fx-core/v8/x/gov/migrations/v8"
)

func migrateGovCustomParam(ctx sdk.Context, keeper *fxgovkeeper.Keeper, storeKey *storetypes.KVStoreKey) error {
	// 1. delete fxParams key
	store.RemoveStoreKeys(ctx, storeKey, fxgovv8.GetRemovedStoreKeys())

	// 2. init custom params
	return keeper.InitCustomParams(ctx)
}

func migrateGovDefaultParams(ctx sdk.Context, keeper *fxgovkeeper.Keeper) error {
	params, err := keeper.Params.Get(ctx)
	if err != nil {
		return err
	}

	minDepositAmount := sdkmath.NewInt(1e18).MulRaw(30)

	params.MinDeposit = sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, minDepositAmount))
	params.ExpeditedMinDeposit = sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, minDepositAmount.MulRaw(govv1.DefaultMinExpeditedDepositTokensRatio)))
	params.MinInitialDepositRatio = sdkmath.LegacyMustNewDecFromStr("0.33").String()
	params.MinDepositRatio = sdkmath.LegacyMustNewDecFromStr("0").String()

	return keeper.Params.Set(ctx, params)
}
