package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) GetOracleDelegateToken(ctx sdk.Context, delegateAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdkmath.Int, error) {
	delegation, found := k.stakingKeeper.GetDelegation(ctx, delegateAddr, valAddr)
	if !found {
		return sdkmath.ZeroInt(), errorsmod.Wrap(types.ErrInvalid, "no delegation for (address, validator) tuple")
	}
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return sdkmath.ZeroInt(), stakingtypes.ErrNoValidatorFound
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
