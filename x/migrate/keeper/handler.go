package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// MigrateHandler specifies the type of function that is called when a migration is applied
type MigrateHandler func(ctx sdk.Context, k Keeper, from sdk.AccAddress, to common.Address) error

type MigrateI interface {
	Validate(ctx sdk.Context, k Keeper, from sdk.AccAddress, to common.Address) error
	Execute(ctx sdk.Context, k Keeper, from sdk.AccAddress, to common.Address) error
	// todo replace the parameter Keeper with codec
}
