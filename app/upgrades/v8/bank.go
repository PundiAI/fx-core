package v8

import (
	"errors"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func migrateBankModule(ctx sdk.Context, bankKeeper bankkeeper.Keeper) error {
	sendEnabledEntry, found := bankKeeper.GetSendEnabledEntry(ctx, fxtypes.LegacyFXDenom)
	if found {
		bankKeeper.DeleteSendEnabled(ctx, fxtypes.LegacyFXDenom)
		bankKeeper.SetSendEnabled(ctx, fxtypes.DefaultDenom, sendEnabledEntry.Enabled)
	}

	var err error
	fxSupply := bankKeeper.GetSupply(ctx, fxtypes.LegacyFXDenom)
	apundiaiSupply := sdkmath.ZeroInt()

	bk, ok := bankKeeper.(bankkeeper.BaseKeeper)
	if !ok {
		return errors.New("bank keeper not implement bank.BaseKeeper")
	}
	bk.IterateAllBalances(ctx, func(address sdk.AccAddress, coin sdk.Coin) (stop bool) {
		if coin.Denom != fxtypes.LegacyFXDenom {
			return false
		}
		if err = bk.Balances.Remove(ctx, collections.Join(address, coin.Denom)); err != nil {
			return true
		}
		coin.Denom = fxtypes.DefaultDenom
		coin.Amount = fxtypes.SwapAmount(coin.Amount)
		if !coin.IsPositive() {
			return false
		}
		apundiaiSupply = apundiaiSupply.Add(coin.Amount)
		if err = bk.Balances.Set(ctx, collections.Join(address, coin.Denom), coin.Amount); err != nil {
			return true
		}
		return false
	})
	if err != nil {
		return err
	}

	ctx.Logger().Info("migrate fx to apundiai", "module", "upgrade", "FX supply", fxSupply.Amount.String(), "apundiai supply", apundiaiSupply.String())
	if err = bk.Supply.Remove(ctx, fxtypes.LegacyFXDenom); err != nil {
		return err
	}
	return bk.Supply.Set(ctx, fxtypes.DefaultDenom, apundiaiSupply)
}
