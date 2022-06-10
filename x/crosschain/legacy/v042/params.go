package v042

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/functionx/fx-core/x/crosschain/types"
)

var ParamOracleDepositThreshold = []byte("OracleDepositThreshold")

func GetParamSetPairs(params *types.Params) paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(types.ParamsStoreKeyGravityID, &params.GravityId, nil),
		paramtypes.NewParamSetPair(types.ParamsStoreKeyAverageBlockTime, &params.AverageBlockTime, nil),
		paramtypes.NewParamSetPair(types.ParamsStoreKeyExternalBatchTimeout, &params.ExternalBatchTimeout, nil),
		paramtypes.NewParamSetPair(types.ParamsStoreKeyAverageExternalBlockTime, &params.AverageExternalBlockTime, nil),
		paramtypes.NewParamSetPair(types.ParamsStoreKeySignedWindow, &params.SignedWindow, nil),
		paramtypes.NewParamSetPair(types.ParamsStoreSlashFraction, &params.SlashFraction, nil),
		paramtypes.NewParamSetPair(types.ParamStoreOracleSetUpdatePowerChangePercent, &params.OracleSetUpdatePowerChangePercent, nil),
		paramtypes.NewParamSetPair(types.ParamStoreIbcTransferTimeoutHeight, &params.IbcTransferTimeoutHeight, nil),
		paramtypes.NewParamSetPair(types.ParamStoreOracles, &params.Oracles, nil),
		paramtypes.NewParamSetPair(types.ParamOracleDelegateThreshold, &params.DelegateThreshold, nil),
		paramtypes.NewParamSetPair(ParamOracleDepositThreshold, &params.DelegateMultiple, nil),
	}
}
