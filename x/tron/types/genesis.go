package types

import (
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
)

func DefaultGenesisState() *crosschaintypes.GenesisState {
	params := crosschaintypes.DefaultParams()
	params.GravityId = "fx-bridge-tron"
	return &crosschaintypes.GenesisState{
		Params: params,
	}
}
