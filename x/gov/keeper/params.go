package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/gov/types"
)

// GetFXParams gets the gov module's parameters.
func (keeper Keeper) GetFXParams(ctx sdk.Context, msgType string) (params types.Params) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.ParamsByMsgTypeKey(msgType))
	if bz != nil {
		keeper.cdc.MustUnmarshal(bz, &params)
		return params
	}
	v1Params := keeper.Keeper.GetParams(ctx)
	params = *types.NewParam(
		msgType,
		v1Params.GetMinDeposit(),
		sdk.NewCoin(fxtypes.DefaultDenom, types.DefaultMinInitialDeposit),
		v1Params.VotingPeriod,
		v1Params.Quorum,
		v1Params.MaxDepositPeriod,
		v1Params.Threshold,
		v1Params.VetoThreshold,
		v1Params.MinInitialDepositRatio,
		v1Params.BurnVoteQuorum,
		v1Params.BurnProposalDepositPrevote,
		v1Params.BurnVoteVeto,
	)
	return params
}

// SetFXParams sets the gov module's parameters.
func (keeper Keeper) SetFXParams(ctx sdk.Context, params *types.Params) error {
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
		if err := keeper.SetFXParams(ctx, param); err != nil {
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
	return keeper.GetFXParams(ctx, msgType).MinInitialDeposit
}

func (keeper Keeper) GetMinDeposit(ctx sdk.Context, msgType string) []sdk.Coin {
	return keeper.GetFXParams(ctx, msgType).MinDeposit
}

func (keeper Keeper) GetMaxDepositPeriod(ctx sdk.Context, msgType string) *time.Duration {
	return keeper.GetFXParams(ctx, msgType).MaxDepositPeriod
}

func (keeper Keeper) GetVotingPeriod(ctx sdk.Context, msgType string) *time.Duration {
	return keeper.GetFXParams(ctx, msgType).VotingPeriod
}

func (keeper Keeper) GetQuorum(ctx sdk.Context, msgType string) string {
	return keeper.GetFXParams(ctx, msgType).Quorum
}

func (keeper Keeper) GetThreshold(ctx sdk.Context, msgType string) string {
	return keeper.GetFXParams(ctx, msgType).Threshold
}

func (keeper Keeper) GetVetoThreshold(ctx sdk.Context, msgType string) string {
	return keeper.GetFXParams(ctx, msgType).VetoThreshold
}
