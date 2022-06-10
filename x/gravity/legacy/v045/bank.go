package v045

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	ethtypes "github.com/functionx/fx-core/x/eth/types"
	gravitytypes "github.com/functionx/fx-core/x/gravity/types"
)

func MigrateBank(ctx sdk.Context, accountKeeper AccountKeeper, bankKeeper BankKeeper) error {
	moduleAddr := accountKeeper.GetModuleAddress(gravitytypes.ModuleName)
	balances := bankKeeper.GetAllBalances(ctx, moduleAddr)
	if balances.IsZero() {
		return nil
	}
	return bankKeeper.SendCoinsFromModuleToModule(ctx, gravitytypes.ModuleName, ethtypes.ModuleName, balances)
}
