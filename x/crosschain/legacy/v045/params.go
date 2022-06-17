package v045

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	v042 "github.com/functionx/fx-core/x/crosschain/legacy/v042"
	"github.com/functionx/fx-core/x/crosschain/types"
)

func MigrateParams(ctx sdk.Context, paramStore *paramtypes.Subspace, moduleName string, paramsKey sdk.StoreKey) error {
	params := types.Params{}
	paramSetPairs := v042.GetParamSetPairs(&params)
	for _, pair := range paramSetPairs {
		paramStore.Get(ctx, pair.Key, pair.Value)
	}
	params.DelegateMultiple = types.DefaultOracleDelegateThreshold
	if err := params.ValidateBasic(); err != nil {
		return err
	}
	paramStore.SetParamSet(ctx, &params)

	paramsStore := prefix.NewStore(ctx.KVStore(paramsKey), append([]byte(moduleName), '/'))
	paramsStore.Delete(v042.ParamStoreOracles)
	paramsStore.Delete(v042.ParamOracleDepositThreshold)
	return nil
}
