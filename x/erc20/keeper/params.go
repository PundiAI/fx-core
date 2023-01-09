package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v3/x/erc20/types"
)

// GetParams returns the total set of erc20 parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// GetEnableErc20 returns the EnableErc20 parameter.
func (k Keeper) GetEnableErc20(ctx sdk.Context) bool {
	var enableErc20 bool
	k.paramSpace.GetIfExists(ctx, types.ParamStoreKeyEnableErc20, &enableErc20)
	return enableErc20
}

// GetEnableEVMHook returns the EnableEVMHook parameter.
func (k Keeper) GetEnableEVMHook(ctx sdk.Context) bool {
	var enableEVMHook bool
	k.paramSpace.GetIfExists(ctx, types.ParamStoreKeyEnableEVMHook, &enableEVMHook)
	return enableEVMHook
}

// GetIbcTimeout returns the IbcTimeout parameter.
func (k Keeper) GetIbcTimeout(ctx sdk.Context) time.Duration {
	var ibcTimeout time.Duration
	k.paramSpace.GetIfExists(ctx, types.ParamStoreKeyIBCTimeout, &ibcTimeout)
	return ibcTimeout
}

// SetParams sets the erc20 parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
