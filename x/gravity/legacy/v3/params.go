package v3

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	"github.com/functionx/fx-core/v3/x/gravity/types"
)

func MigrateParams(ctx sdk.Context, legacyAmino *codec.LegacyAmino, paramsStoreKey sdk.StoreKey, toModuleName string) error {
	paramsStore := prefix.NewStore(ctx.KVStore(paramsStoreKey), append([]byte(types.ModuleName), '/'))
	gravityParams := &types.Params{}
	isExist := false
	for _, pair := range gravityParams.ParamSetPairs() {
		bz := paramsStore.Get(pair.Key)
		if len(bz) <= 0 {
			continue
		}
		isExist = true
		if err := legacyAmino.UnmarshalJSON(bz, pair.Value); err != nil {
			panic(err)
		}
		paramsStore.Delete(pair.Key)
	}
	if !isExist {
		return nil
	}
	if err := gravityParams.ValidateBasic(); err != nil {
		return err
	}
	params := crosschaintypes.Params{
		GravityId:                         gravityParams.GravityId,
		AverageBlockTime:                  gravityParams.AverageBlockTime,
		ExternalBatchTimeout:              gravityParams.TargetBatchTimeout,
		AverageExternalBlockTime:          gravityParams.AverageEthBlockTime,
		SignedWindow:                      gravityParams.SignedValsetsWindow,
		SlashFraction:                     gravityParams.SlashFractionValset,
		OracleSetUpdatePowerChangePercent: gravityParams.ValsetUpdatePowerChangePercent,
		IbcTransferTimeoutHeight:          gravityParams.IbcTransferTimeoutHeight,
		DelegateThreshold:                 sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(100).Mul(sdk.DefaultPowerReduction)),
		DelegateMultiple:                  crosschaintypes.DefaultOracleDelegateThreshold,
	}

	if err := params.ValidateBasic(); err != nil {
		return err
	}

	store := prefix.NewStore(ctx.KVStore(paramsStoreKey), append([]byte(toModuleName), '/'))
	for _, pair := range params.ParamSetPairs() {
		bz, err := legacyAmino.MarshalJSON(pair.Value)
		if err != nil {
			panic(err)
		}
		store.Set(pair.Key, bz)
	}
	return nil
}
