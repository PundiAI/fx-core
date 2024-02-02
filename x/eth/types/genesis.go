package types

import (
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
)

func DefaultGenesisState() *crosschaintypes.GenesisState {
	params := crosschaintypes.DefaultParams()
	params.GravityId = "fx-bridge-eth"
	params.AverageExternalBlockTime = 15_000
	return &crosschaintypes.GenesisState{
		Params: params,
	}
}
