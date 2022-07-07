package migrate

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/types"

	"github.com/functionx/fx-core/x/migrate/keeper"
	"github.com/functionx/fx-core/x/migrate/types"
)

// InitGenesis import module genesis
func InitGenesis(ctx sdk.Context, k keeper.Keeper, state types.GenesisState) {
	for _, record := range state.MigrateRecords {
		fromAddr, err := sdk.AccAddressFromBech32(record.From)
		if err != nil {
			panic(err)
		}
		if err = fxtypes.ValidateEthereumAddress(record.To); err != nil {
			panic(err)
		}
		k.SetMigrateRecord(ctx, fromAddr, common.HexToAddress(record.To))
	}
}

// ExportGenesis export module status
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesisState := &types.GenesisState{}
	k.IterateMigrateRecords(ctx, func(record types.MigrateRecord) bool {
		genesisState.MigrateRecords = append(genesisState.MigrateRecords, record)
		return false
	})
	return genesisState
}
