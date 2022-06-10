package v045

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	v042 "github.com/functionx/fx-core/x/crosschain/legacy/v042"
	"github.com/functionx/fx-core/x/crosschain/types"
)

func MigrateParams(ctx sdk.Context, paramStore *paramtypes.Subspace) error {
	params := types.Params{}
	paramSetPairs := v042.GetParamSetPairs(&params)
	for _, pair := range paramSetPairs {
		paramStore.Get(ctx, pair.Key, pair.Value)
	}
	if err := params.ValidateBasic(); err != nil {
		return err
	}
	paramStore.SetParamSet(ctx, &params)
	return nil
}
