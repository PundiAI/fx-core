package types

import (
	crosschaintypes "github.com/functionx/fx-core/v4/x/crosschain/types"
)

func DefaultGenesisState() *crosschaintypes.GenesisState {
	params := crosschaintypes.DefaultParams()
	params.GravityId = "fx-arbitrum-bridge"
	params.AverageExternalBlockTime = 2_000
	return &crosschaintypes.GenesisState{
		Params: params,
	}
}
