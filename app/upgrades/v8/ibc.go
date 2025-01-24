package v8

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"

	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func migrateTransferTokenInEscrow(ctx sdk.Context, transferKeeper ibctransferkeeper.Keeper) {
	escrowDenoms := GetMigrateEscrowDenoms(ctx.ChainID())
	for oldDenom, newDenom := range escrowDenoms {
		totalEscrow := transferKeeper.GetTotalEscrowForDenom(ctx, oldDenom)
		newAmount := totalEscrow.Amount
		if oldDenom == fxtypes.LegacyFXDenom {
			newAmount = fxtypes.SwapAmount(newAmount)
		}
		// first remove old denom
		transferKeeper.SetTotalEscrowForDenom(ctx, sdk.NewCoin(oldDenom, sdkmath.ZeroInt()))
		// then add new denom
		transferKeeper.SetTotalEscrowForDenom(ctx, sdk.NewCoin(newDenom, newAmount))
	}
}
