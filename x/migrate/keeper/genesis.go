package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v6/types"
	"github.com/functionx/fx-core/v6/x/migrate/types"
)

// InitGenesis import module genesis
func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
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
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesisState := &types.GenesisState{}
	k.IterateMigrateRecords(ctx, func(record types.MigrateRecord) bool {
		genesisState.MigrateRecords = append(genesisState.MigrateRecords, record)
		return false
	})
	return genesisState
}
