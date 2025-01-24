package v8

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"

	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func migrateFeemarketGasPrice(ctx sdk.Context, feemarketKeeper feemarketkeeper.Keeper) error {
	params := feemarketKeeper.GetParams(ctx)
	params.BaseFee = sdkmath.NewInt(fxtypes.DefaultGasPrice)
	params.MinGasPrice = sdkmath.LegacyNewDec(fxtypes.DefaultGasPrice)
	return feemarketKeeper.SetParams(ctx, params)
}
