package migrate

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/x/migrate/keeper"
	"github.com/functionx/fx-core/x/migrate/types"
)

// InitGenesis import module genesis
func InitGenesis(_ sdk.Context, _ keeper.Keeper, _ types.GenesisState) {
}

// ExportGenesis export module status
func ExportGenesis(_ sdk.Context, _ keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{}
}
