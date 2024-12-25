package types

import (
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func DefaultGenesisState() *crosschaintypes.GenesisState {
	params := crosschaintypes.DefaultParams()
	params.GravityId = "fx-arbitrum-bridge"
	params.AverageExternalBlockTime = 250
	return &crosschaintypes.GenesisState{
		Params: params,
	}
}
