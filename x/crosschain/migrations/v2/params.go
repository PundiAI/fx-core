package v2

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	crosschainv1 "github.com/functionx/fx-core/v3/x/crosschain/migrations/v1"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

func MigrateParams(ctx sdk.Context, moduleName string, legacyAmino *codec.LegacyAmino, paramsKey sdk.StoreKey) error {
	params := types.Params{}
	paramSetPairs := crosschainv1.GetParamSetPairs(&params)
	paramsStore := prefix.NewStore(ctx.KVStore(paramsKey), append([]byte(moduleName), '/'))
	for _, pair := range paramSetPairs {
		bz := paramsStore.Get(pair.Key)
		if err := legacyAmino.UnmarshalJSON(bz, pair.Value); err != nil {
			return fmt.Errorf("%s: %s", err.Error(), pair.Key)
		}
	}
	paramsStore.Delete(crosschainv1.ParamStoreOracles)
	paramsStore.Delete(crosschainv1.ParamOracleDepositThreshold)

	params.DelegateMultiple = types.DefaultOracleDelegateThreshold
	params.AverageBlockTime = 7_000
	ctx.Logger().Debug("migrate params", "module", moduleName, "averageBlockTime", params.AverageBlockTime)
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

func CheckInitialize(ctx sdk.Context, moduleName string, paramsKey sdk.StoreKey) bool {
	paramsStore := prefix.NewStore(ctx.KVStore(paramsKey), append([]byte(moduleName), '/'))
	return paramsStore.Has(types.ParamsStoreKeyGravityID)
}
