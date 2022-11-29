package types

import (
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
)

func DefaultGenesisState() *crosschaintypes.GenesisState {
	params := crosschaintypes.DefaultParams()
	params.GravityId = "fx-eth-bridge"
	params.AverageExternalBlockTime = 15_000
	return &crosschaintypes.GenesisState{
		Params: params,
	}
}
