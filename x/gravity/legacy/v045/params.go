package v045

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/types"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
	gravitytypes "github.com/functionx/fx-core/x/gravity/types"
)

func MigrateParams(ctx sdk.Context, gravityParams gravitytypes.Params, ethKeeper EthKeeper, oracles []string) error {
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
