package types

import (
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func DefaultGenesisState() *crosschaintypes.GenesisState {
	params := crosschaintypes.DefaultParams()
	params.GravityId = "fx-bridge-eth"
	params.AverageExternalBlockTime = 15_000
	return &crosschaintypes.GenesisState{
		Params: params,
	}
}
