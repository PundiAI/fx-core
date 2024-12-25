package v8

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	crosschainkeeper "github.com/pundiai/fx-core/v8/x/crosschain/keeper"
	"github.com/pundiai/fx-core/v8/x/erc20/keeper"
)

type Migrator struct {
	storeKey          storetypes.StoreKey
	cdc               codec.BinaryCodec
	keeper            keeper.Keeper
	bankKeeper        bankkeeper.Keeper
	crosschainKeepers []crosschainkeeper.Keeper
}

func NewMigrator(storeKey storetypes.StoreKey, cdc codec.BinaryCodec, keeper keeper.Keeper, bk bankkeeper.Keeper, cks []crosschainkeeper.Keeper) Migrator {
	return Migrator{
		storeKey:          storeKey,
		cdc:               cdc,
		keeper:            keeper,
		bankKeeper:        bk,
		crosschainKeepers: cks,
	}
}

// Migrate3to4 migrates from version 3 to 4.
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	if err := m.migrateKeys(ctx); err != nil {
		return err
	}
	return m.MigrateToken(ctx)
}
