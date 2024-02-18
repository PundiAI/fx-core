package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/gov/types"
)

// GetParams gets the gov module's parameters.
func (keeper Keeper) GetParams(ctx sdk.Context, msgType string) (params types.Params) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.ParamsByMsgTypeKey(msgType))
	if bz != nil {
		keeper.cdc.MustUnmarshal(bz, &params)
		return params
	}
	depositParams := keeper.GetDepositParams(ctx)
	votingParams := keeper.GetVotingParams(ctx)
	tallyParams := keeper.GetTallyParams(ctx)
	params = *types.NewParam(
		msgType,
		depositParams.GetMinDeposit(),
		sdk.NewCoin(fxtypes.DefaultDenom, types.DefaultMinInitialDeposit),
		votingParams.VotingPeriod,
		tallyParams.Quorum,
		depositParams.MaxDepositPeriod,
		tallyParams.Threshold,
		tallyParams.VetoThreshold,
	)
	return params
}

// SetParams sets the gov module's parameters.
func (keeper Keeper) SetParams(ctx sdk.Context, params *types.Params) error {
	store := ctx.KVStore(keeper.storeKey)
	bz, err := keeper.cdc.Marshal(params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsByMsgTypeKey(params.MsgType), bz)

	return nil
}

// SetAllParams sets batch the gov module's parameters.
func (keeper Keeper) SetAllParams(ctx sdk.Context, params []*types.Params) error {
	for _, param := range params {
		if err := keeper.SetParams(ctx, param); err != nil {
			return err
		}
	}
	return nil
}

// GetEGFParams gets the gov module's parameters.
func (keeper Keeper) GetEGFParams(ctx sdk.Context) (params types.EGFParams) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.FxEGFParamsKey)
	if bz == nil {
		return params
	}
	keeper.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetEGFParams sets the gov module's parameters.
func (keeper Keeper) SetEGFParams(ctx sdk.Context, params *types.EGFParams) error {
	store := ctx.KVStore(keeper.storeKey)
	bz, err := keeper.cdc.Marshal(params)
	if err != nil {
		return err
	}
	store.Set(types.FxEGFParamsKey, bz)
	return nil
}

func (keeper Keeper) GetMinInitialDeposit(ctx sdk.Context, msgType string) sdk.Coin {
	return keeper.GetParams(ctx, msgType).MinInitialDeposit
}

func (keeper Keeper) GetMinDeposit(ctx sdk.Context, msgType string) []sdk.Coin {
	return keeper.GetParams(ctx, msgType).MinDeposit
}

func (keeper Keeper) GetMaxDepositPeriod(ctx sdk.Context, msgType string) *time.Duration {
	return keeper.GetParams(ctx, msgType).MaxDepositPeriod
}

func (keeper Keeper) GetVotingPeriod(ctx sdk.Context, msgType string) *time.Duration {
	return keeper.GetParams(ctx, msgType).VotingPeriod
}

func (keeper Keeper) GetQuorum(ctx sdk.Context, msgType string) string {
	return keeper.GetParams(ctx, msgType).Quorum
}

func (keeper Keeper) GetThreshold(ctx sdk.Context, msgType string) string {
	return keeper.GetParams(ctx, msgType).Threshold
}

func (keeper Keeper) GetVetoThreshold(ctx sdk.Context, msgType string) string {
	return keeper.GetParams(ctx, msgType).VetoThreshold
}
