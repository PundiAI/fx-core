package v045

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v042 "github.com/functionx/fx-core/x/crosschain/legacy/v042"
	"github.com/functionx/fx-core/x/crosschain/types"
)

func MigrateParams(ctx sdk.Context, moduleName string, legacyAmino *codec.LegacyAmino, paramsKey sdk.StoreKey) error {
	params := types.Params{}
	paramSetPairs := v042.GetParamSetPairs(&params)
	paramsStore := prefix.NewStore(ctx.KVStore(paramsKey), append([]byte(moduleName), '/'))
	for _, pair := range paramSetPairs {
		bz := paramsStore.Get(pair.Key)
		if err := legacyAmino.UnmarshalJSON(bz, pair.Value); err != nil {
			return err
		}
	}
	paramsStore.Delete(v042.ParamStoreOracles)
	paramsStore.Delete(v042.ParamOracleDepositThreshold)

	params.DelegateMultiple = types.DefaultOracleDelegateThreshold
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
