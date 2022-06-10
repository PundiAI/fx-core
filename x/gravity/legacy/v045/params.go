package v045

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v042 "github.com/functionx/fx-core/x/gravity/legacy/v042"
	"github.com/functionx/fx-core/x/gravity/types"

	fxtypes "github.com/functionx/fx-core/types"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
)

func MigrateParams(ctx sdk.Context, cdc codec.BinaryCodec, legacyAmino *codec.LegacyAmino, paramsKey sdk.StoreKey, ethKeeper EthKeeper, oracles []string) error {

	paramsStore := prefix.NewStore(ctx.KVStore(paramsKey), append([]byte(types.ModuleName), '/'))
	gravityParams := &v042.Params{}
	for _, pair := range gravityParams.ParamSetPairs() {
		bz := paramsStore.Get(pair.Key)
		if err := legacyAmino.UnmarshalJSON(bz, pair.Value); err != nil {
			panic(err)
		}
		paramsStore.Delete(pair.Key)
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
		Oracles:                           oracles,
		DelegateThreshold:                 sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1000).Mul(sdk.DefaultPowerReduction)),
		DelegateMultiple:                  crosschaintypes.DefaultOracleDelegateThreshold,
	}

	if err := params.ValidateBasic(); err != nil {
		return err
	}

	ethKeeper.SetParams(ctx, params)
	// remove gravity params
	return nil
}
