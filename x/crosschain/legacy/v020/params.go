package v020

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	v010 "github.com/functionx/fx-core/x/crosschain/legacy/v010"
	"github.com/functionx/fx-core/x/crosschain/types"
)

func MigrateParams(ctx sdk.Context, moduleName string, legacyAmino *codec.LegacyAmino, paramsKey sdk.StoreKey) error {
	params := types.Params{}
	paramSetPairs := v010.GetParamSetPairs(&params)
	paramsStore := prefix.NewStore(ctx.KVStore(paramsKey), append([]byte(moduleName), '/'))
	for _, pair := range paramSetPairs {
		bz := paramsStore.Get(pair.Key)
		if err := legacyAmino.UnmarshalJSON(bz, pair.Value); err != nil {
			return err
		}
	}
	paramsStore.Delete(v010.ParamStoreOracles)
	paramsStore.Delete(v010.ParamOracleDepositThreshold)

	params.DelegateMultiple = types.DefaultOracleDelegateThreshold
	params.AverageBlockTime = 7_000
	logger := ctx.Logger().With("module", "x/"+moduleName)
	logger.Debug("migrate params", "averageBlockTime", params.AverageBlockTime)
	if err := params.ValidateBasic(); err != nil {
		return err
	}

	for _, pair := range params.ParamSetPairs() {
		bz, err := legacyAmino.MarshalJSON(pair.Value)
		if err != nil {
			return err
		}
		paramsStore.Set(pair.Key, bz)
	}
	return nil
}
