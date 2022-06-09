package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/functionx/fx-core/x/gravity/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	ethereumMsgServer ethereumMsgServer
	paramSpace        paramtypes.Subspace
}

// NewKeeper returns a new instance of the gravity keeper
func NewKeeper(paramSpace paramtypes.Subspace, ethMsgServer ethereumMsgServer) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		ethereumMsgServer: ethMsgServer,
		paramSpace:        paramSpace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetParams returns the parameters from the store
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}
