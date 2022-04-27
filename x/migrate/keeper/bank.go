package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	migratetypes "github.com/functionx/fx-core/x/migrate/types"
)

type BankMigrate struct {
	bankKeeper migratetypes.BankKeeper
}

func NewBankMigrate(bk migratetypes.BankKeeper) MigrateI {
	return &BankMigrate{bankKeeper: bk}
}

func (m *BankMigrate) Validate(_ sdk.Context, _ Keeper, _, _ sdk.AccAddress) error {
	return nil
}

func (m *BankMigrate) Execute(ctx sdk.Context, _ Keeper, from, to sdk.AccAddress) error {
	balances := m.bankKeeper.GetAllBalances(ctx, from)
	if balances.IsZero() {
		return nil
	}
	return m.bankKeeper.SendCoins(ctx, from, to, balances)
}
