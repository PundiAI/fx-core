package types

import (
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
)

func DefaultGenesisState() *crosschaintypes.GenesisState {
	params := crosschaintypes.DefaultParams()
	params.GravityId = "fx-tron-bridge"
	return &crosschaintypes.GenesisState{
		Params: params,
	}
}
