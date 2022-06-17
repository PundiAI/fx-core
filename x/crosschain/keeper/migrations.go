package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	v045 "github.com/functionx/fx-core/x/crosschain/legacy/v045"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper      Keeper
	sk          v045.StakingKeeper
	bk          v045.BankKeeper
	paramsKey   sdk.StoreKey
	legacyAmino *codec.LegacyAmino
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, sk v045.StakingKeeper, bk v045.BankKeeper, legacyAmino *codec.LegacyAmino, paramsKey sdk.StoreKey) Migrator {
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
	if err := v045.MigrateParams(ctx, m.keeper.moduleName, m.legacyAmino, m.paramsKey); err != nil {
		return err
	}
	oracles, delegateValidator, err := v045.MigrateOracle(ctx, m.keeper.cdc, m.keeper.storeKey, m.sk)
	if err != nil {
		return err
	}
	if err := v045.MigrateDepositToStaking(ctx, m.keeper.moduleName, m.sk, m.bk, oracles, delegateValidator); err != nil {
		return err
	}
	v045.MigrateStore(ctx, m.keeper.storeKey)
	return nil
}
