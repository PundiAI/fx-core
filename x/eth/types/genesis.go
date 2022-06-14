package types

import (
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
)

func DefaultGenesisState() *crosschaintypes.GenesisState {
	params := crosschaintypes.DefaultParams()
	params.GravityId = "fx-bridge-eth"
	params.AverageExternalBlockTime = 15 * 1e3
	return &crosschaintypes.GenesisState{
		Params: params,
	}
}
