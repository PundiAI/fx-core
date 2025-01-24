package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	fxevmkeeper "github.com/pundiai/fx-core/v8/x/evm/keeper"
)

func migrateEvmParams(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper) error {
	params := evmKeeper.GetParams(ctx)
	params.EvmDenom = fxtypes.DefaultDenom
	params.HeaderHashNum = evmtypes.DefaultHeaderHashNum
	return evmKeeper.SetParams(ctx, params)
}
