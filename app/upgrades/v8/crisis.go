package v8

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"

	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func migrateCrisisModule(ctx sdk.Context, crisisKeeper *crisiskeeper.Keeper) error {
	constantFee, err := crisisKeeper.ConstantFee.Get(ctx)
	if err != nil {
		return err
	}
	constantFee.Denom = fxtypes.DefaultDenom
	constantFee.Amount = sdkmath.NewInt(133).MulRaw(1e18)
	return crisisKeeper.ConstantFee.Set(ctx, constantFee)
}
