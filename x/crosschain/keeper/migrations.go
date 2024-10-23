package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/x/crosschain/migrations/v8"
)

type Migrator struct {
	keeper Keeper
}

func NewMigrator(k Keeper) Migrator {
	return Migrator{
		keeper: k,
	}
}

func (m Migrator) Migrate7to8(ctx sdk.Context) error {
	return v8.Migrate(ctx, m.keeper.storeKey, m.keeper.cdc)
}

func (m Migrator) Migrate7to8WithArbExternalBlockTime(ctx sdk.Context) error {
	if err := m.Migrate7to8(ctx); err != nil {
		return err
	}
	params := m.keeper.GetParams(ctx)
	params.AverageExternalBlockTime = 250
	return m.keeper.SetParams(ctx, &params)
}
