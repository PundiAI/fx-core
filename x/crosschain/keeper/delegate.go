package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (k Keeper) GetOracleDelegateToken(ctx sdk.Context, delegateAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdkmath.Int, error) {
	delegation, err := k.stakingKeeper.GetDelegation(ctx, delegateAddr, valAddr)
	if err != nil {
		return sdkmath.ZeroInt(), errorsmod.Wrap(err, "no delegation for (address, validator) tuple")
	}
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return sdkmath.ZeroInt(), errorsmod.Wrap(err, "no validator for address")
	}

	delegateToken := validator.TokensFromSharesTruncated(delegation.GetShares()).TruncateInt()
	sharesTruncated, err := validator.SharesFromTokensTruncated(delegateToken)
	if err != nil {
		return sdkmath.ZeroInt(), errorsmod.Wrapf(types.ErrInvalid, "shares from tokens:%v", delegateToken)
	}
	delShares := delegation.GetShares()
	if sharesTruncated.GT(delShares) {
		delegateToken = validator.TokensFromSharesTruncated(sharesTruncated).TruncateInt()
	}
	return delegateToken, nil
}
