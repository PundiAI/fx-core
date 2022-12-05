package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	v020 "github.com/functionx/fx-core/v2/x/crosschain/legacy/v2"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper      Keeper
	sk          v020.StakingKeeper
	bk          v020.BankKeeper
	paramsKey   sdk.StoreKey
	legacyAmino *codec.LegacyAmino
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, sk v020.StakingKeeper, bk v020.BankKeeper, legacyAmino *codec.LegacyAmino, paramsKey sdk.StoreKey) Migrator {
	return Migrator{
		keeper:      keeper,
		sk:          sk,
		bk:          bk,
		paramsKey:   paramsKey,
		legacyAmino: legacyAmino,
	}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	if err := v020.MigrateParams(ctx, m.keeper.moduleName, m.legacyAmino, m.paramsKey); err != nil {
		return err
	}
	oracles, delegateValAddr, err := v020.MigrateOracle(ctx, m.keeper.cdc, m.keeper.storeKey, m.sk)
	if err != nil {
		return err
	}
	if err := v020.MigrateDepositToStaking(ctx, m.keeper.moduleName, m.sk, m.bk, oracles, delegateValAddr); err != nil {
		return err
	}
	v020.MigrateStore(ctx, m.keeper.storeKey)
	return nil
}
