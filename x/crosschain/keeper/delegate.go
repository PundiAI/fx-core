package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetOracleDelegateToken(ctx sdk.Context, delegateAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdkmath.Int, error) {
	delegation, err := k.stakingKeeper.GetDelegation(ctx, delegateAddr, valAddr)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	delegateToken := validator.TokensFromSharesTruncated(delegation.GetShares()).TruncateInt()
	sharesTruncated, err := validator.SharesFromTokensTruncated(delegateToken)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	delShares := delegation.GetShares()
	if sharesTruncated.GT(delShares) {
		delegateToken = validator.TokensFromSharesTruncated(sharesTruncated).TruncateInt()
	}
	return delegateToken, nil
}
