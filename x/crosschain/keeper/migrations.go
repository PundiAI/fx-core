package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

type Migrator struct {
	keeper Keeper
}

func NewMigrator(k Keeper) Migrator {
	return Migrator{
		keeper: k,
	}
}

func (m Migrator) Migrate(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)

	params.BridgeCallTimeout = types.DefBridgeCallTimeout
	params.BridgeCallMaxGasLimit = types.MaxGasLimit

	enablePending := false
	if ctx.ChainID() == fxtypes.TestnetChainId {
		enablePending = true
	}
	params.EnableSendToExternalPending = enablePending
	params.EnableBridgeCallPending = enablePending

	if err := m.keeper.SetParams(ctx, &params); err != nil {
		return err
	}
	return nil
}
