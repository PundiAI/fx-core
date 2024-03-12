package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

type Migrator struct {
	keeper Keeper
}

func NewMigrator(k Keeper) Migrator {
	return Migrator{
		keeper: k,
	}
}

func (m Migrator) Migrate(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	params.BridgeCallRefundTimeout = types.DefaultBridgeCallRefundTimeout
	return m.keeper.SetParams(ctx, &params)
}
