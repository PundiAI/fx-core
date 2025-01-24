package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/pundiai/fx-core/v8/app/keepers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func migrateCrosschainParams(ctx sdk.Context, keepers keepers.CrosschainKeepers) error {
	for _, k := range keepers.ToSlice() {
		params := k.GetParams(ctx)
		params.DelegateThreshold.Denom = fxtypes.DefaultDenom
		params.DelegateThreshold.Amount = fxtypes.SwapAmount(params.DelegateThreshold.Amount)
		if !params.DelegateThreshold.IsPositive() {
			return sdkerrors.ErrInvalidCoins.Wrapf("module %s invalid delegate threshold: %s",
				k.ModuleName(), params.DelegateThreshold.String())
		}
		if err := k.SetParams(ctx, &params); err != nil {
			return err
		}
	}
	return nil
}

func migrateCrosschainModuleAccount(ctx sdk.Context, ak authkeeper.AccountKeeper) error {
	addr, perms := ak.GetModuleAddressAndPermissions(crosschaintypes.ModuleName)
	if addr == nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("crosschain module empty permissions")
	}
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("crosschain account not exist")
	}
	baseAcc, ok := acc.(*authtypes.BaseAccount)
	if !ok {
		return sdkerrors.ErrInvalidAddress.Wrapf("crosschain account not base account")
	}
	macc := authtypes.NewModuleAccount(baseAcc, crosschaintypes.ModuleName, perms...)
	ak.SetModuleAccount(ctx, macc)
	return nil
}

func migrateOracleDelegateAmount(ctx sdk.Context, keepers keepers.CrosschainKeepers) {
	for _, k := range keepers.ToSlice() {
		k.IterateOracle(ctx, func(oracle crosschaintypes.Oracle) bool {
			oracle.DelegateAmount = fxtypes.SwapAmount(oracle.DelegateAmount)
			k.SetOracle(ctx, oracle)
			return false
		})
	}
}
