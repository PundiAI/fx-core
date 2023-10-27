package types

import (
	crosschaintypes "github.com/functionx/fx-core/v6/x/crosschain/types"
)

func DefaultGenesisState() *crosschaintypes.GenesisState {
	params := crosschaintypes.DefaultParams()
	params.GravityId = "fx-tron-bridge"
	params.AverageExternalBlockTime = 3_000
	return &crosschaintypes.GenesisState{
		Params: params,
	}
}
