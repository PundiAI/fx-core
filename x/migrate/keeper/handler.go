package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MigrateHandler specifies the type of function that is called when a migration is applied
type MigrateHandler func(ctx sdk.Context, k Keeper, from, to sdk.AccAddress) error

type MigrateI interface {
	Validate(ctx sdk.Context, k Keeper, from, to sdk.AccAddress) error
	Execute(ctx sdk.Context, k Keeper, from, to sdk.AccAddress) error
}
