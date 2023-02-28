package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var _ stakingtypes.StakingHooks = Keeper{}

type LPTokenHook struct {
	keeper Keeper
}

func (k Keeper) Hooks() LPTokenHook {
	return LPTokenHook{k}
}

// AfterValidatorCreated - call hook if registered
func (h LPTokenHook) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	_, err := h.keeper.DeployLPToken(ctx, valAddr)
	if err != nil {
		// todo - cosmos-sdk v0.46.x will return error
		panic(errortypes.ErrInvalidRequest.Wrapf("failed to deploy lp token contract: %s", err.Error()))
	}
}

// BeforeValidatorModified - call hook if registered
func (h LPTokenHook) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress) {}

// AfterValidatorRemoved - call hook if registered
func (h LPTokenHook) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) {
	if err := h.keeper.SelfDestructLPToken(ctx, valAddr); err != nil {
		// todo - cosmos-sdk v0.46.x will return error
		panic(errortypes.ErrInvalidRequest.Wrapf("failed to selfdestruct: %s", err.Error()))
	}
}

// AfterValidatorBonded - call hook if registered
func (h LPTokenHook) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) {
}

// AfterValidatorBeginUnbonding - call hook if registered
func (h LPTokenHook) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) {
}

// BeforeDelegationCreated - call hook if registered
func (h LPTokenHook) BeforeDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {
}

// BeforeDelegationSharesModified - call hook if registered
func (h LPTokenHook) BeforeDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {
}

// BeforeDelegationRemoved - call hook if registered
func (h LPTokenHook) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {
}

// AfterDelegationModified - call hook if registered
func (h LPTokenHook) AfterDelegationModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {
}

// BeforeValidatorSlashed - call hook if registered
func (h LPTokenHook) BeforeValidatorSlashed(_ sdk.Context, _ sdk.ValAddress, _ sdk.Dec) {}
