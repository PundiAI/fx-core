package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	crosschainkeeper "github.com/functionx/fx-core/x/crosschain/keeper"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/functionx/fx-core/x/gravity/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	legacyMsgServer  legacyMsgServer
	legacyParamSpace paramtypes.Subspace

	crosschainkeeper.EthereumMsgServer
	stakingKeeper types.StakingKeeper
}

// NewKeeper returns a new instance of the gravity keeper
func NewKeeper(paramSpace paramtypes.Subspace, msgServer legacyMsgServer) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		legacyMsgServer:  msgServer,
		legacyParamSpace: paramSpace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetParams returns the parameters from the store
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.legacyParamSpace.GetParamSet(ctx, &params)
	return
}
