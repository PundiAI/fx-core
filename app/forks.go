package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	feemarkettypes "github.com/functionx/fx-core/x/feemarket/types"

	"github.com/functionx/fx-core/app/forks"
	fxtypes "github.com/functionx/fx-core/types"
	erc20types "github.com/functionx/fx-core/x/erc20/types"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
)

func BeginBlockForks(ctx sdk.Context, fxcore *App) {
	switch ctx.BlockHeight() {
	case fxtypes.EvmV1SupportBlock():
		// update FX meta data
		forks.UpdateFXMetadata(ctx, fxcore.BankKeeper, fxcore.keys[banktypes.StoreKey])
		// clear evm v0 kv stores
		forks.ClearEvmV0KVStores(ctx, fxcore.keys)
		// init evm module
		if err := forks.InitSupportEvm(ctx, fxcore.AccountKeeper,
			fxcore.FeeMarketKeeper, feemarkettypes.DefaultParams(),
			fxcore.EvmKeeper, evmtypes.DefaultParams(),
			fxcore.Erc20Keeper, erc20types.DefaultParams(),
		); err != nil {
			panic(err)
		}
	}
}
