package v4

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/erc20/types"
)

func MigratorParam(ctx sdk.Context, legacySubspace types.Subspace, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	if !legacySubspace.HasKeyTable() {
		legacySubspace.WithKeyTable(types.ParamKeyTable())
	}
	var currParams types.Params
	legacySubspace.GetParamSet(ctx, &currParams)
	if err := currParams.Validate(); err != nil {
		return err
	}
	bz := cdc.MustMarshal(&currParams)
	ctx.KVStore(storeKey).Set(types.ParamsKey, bz)
	return nil
}
