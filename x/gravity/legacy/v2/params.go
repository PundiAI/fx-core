package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	"github.com/functionx/fx-core/v3/x/gravity/types"
)

func MigrateParams(legacyAmino *codec.LegacyAmino, paramsStore sdk.KVStore, toModuleName string) error {
	oldStore := prefix.NewStore(paramsStore, append([]byte(types.ModuleName), '/'))
	gravityParams := &types.Params{} // nolint:staticcheck
	isExist := false
	for _, pair := range gravityParams.ParamSetPairs() {
		bz := oldStore.Get(pair.Key)
		if len(bz) <= 0 {
			continue
		}
		isExist = true
		if err := legacyAmino.UnmarshalJSON(bz, pair.Value); err != nil {
			panic(err)
		}
		oldStore.Delete(pair.Key)
	}
	if !isExist {
		return nil
	}
	if err := gravityParams.ValidateBasic(); err != nil {
		return err
	}
	params := crosschaintypes.Params{
		GravityId:                         gravityParams.GravityId,
		AverageBlockTime:                  7000,
		ExternalBatchTimeout:              gravityParams.TargetBatchTimeout,
		AverageExternalBlockTime:          12000,
		SignedWindow:                      30_000,
		SlashFraction:                     sdk.NewDecWithPrec(8, 1), // 80%
		OracleSetUpdatePowerChangePercent: gravityParams.ValsetUpdatePowerChangePercent,
		IbcTransferTimeoutHeight:          gravityParams.IbcTransferTimeoutHeight,
		DelegateThreshold:                 sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(100).Mul(sdk.DefaultPowerReduction)),
		DelegateMultiple:                  crosschaintypes.DefaultOracleDelegateThreshold,
	}

	if err := params.ValidateBasic(); err != nil {
		return err
	}

	newStore := prefix.NewStore(paramsStore, append([]byte(toModuleName), '/'))
	for _, pair := range params.ParamSetPairs() {
		bz, err := legacyAmino.MarshalJSON(pair.Value)
		if err != nil {
			panic(err)
		}
		newStore.Set(pair.Key, bz)
	}
	return nil
}
