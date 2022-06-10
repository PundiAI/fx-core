package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	fxtypes "github.com/functionx/fx-core/types"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
)

func DefaultGenesisState() *crosschaintypes.GenesisState {
	return &crosschaintypes.GenesisState{
		Params: &crosschaintypes.Params{
			GravityId:                         "fx-bridge-tron",
			AverageBlockTime:                  5 * 1e3,
			ExternalBatchTimeout:              24 * 3600 * 1e3,
			AverageExternalBlockTime:          5 * 1e3,
			SignedWindow:                      20 * 1e3,
			SlashFraction:                     sdk.NewDecWithPrec(1, 3),
			OracleSetUpdatePowerChangePercent: sdk.NewDecWithPrec(1, 1),
			IbcTransferTimeoutHeight:          20 * 1e3,
			Oracles:                           []string{},
			DelegateThreshold:                 sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10*1e3)),
			DelegateMultiple:                  10,
		},
	}
}
