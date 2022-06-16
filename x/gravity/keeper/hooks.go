package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Hooks Wrapper struct
type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

// Hooks Create new gravity hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

func (h Hooks) AfterValidatorBeginUnbonding(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) {

	// When Validator starts Unbonding, Persist the block height in the store
	// Later in endblocker, check if there is atleast one validator who started unbonding and create a valset request.
	// The reason for creating valset requests in endblock is to create only one valset request per block if multiple validators starts unbonding at same block.

	if _, found := h.k.GetEthAddressByValidator(ctx, valAddr); !found {
		return
	}
	h.k.SetLastUnBondingBlockHeight(ctx, uint64(ctx.BlockHeight()))

}

func (h Hooks) BeforeDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {
}
func (h Hooks) AfterValidatorCreated(_ sdk.Context, _ sdk.ValAddress)                   {}
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress)                 {}
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) {}

func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {}
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.k.DelEthAddressForValidator(ctx, valAddr)
}
func (h Hooks) BeforeValidatorSlashed(_ sdk.Context, _ sdk.ValAddress, _ sdk.Dec) {}
func (h Hooks) BeforeDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {
}
func (h Hooks) AfterDelegationModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {
}
