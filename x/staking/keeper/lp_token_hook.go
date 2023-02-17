package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

var _ stakingtypes.StakingHooks = Keeper{}

type LPTokenHook struct {
	keeper Keeper
}

func (k Keeper) Hooks() LPTokenHook {
	return LPTokenHook{k}
}

// AfterValidatorCreated - call hook if registered
func (h LPTokenHook) AfterValidatorCreated(_ sdk.Context, _ sdk.ValAddress) {}

// BeforeValidatorModified - call hook if registered
func (h LPTokenHook) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress) {}

// AfterValidatorRemoved - call hook if registered
func (h LPTokenHook) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.keeper.deleteLPTokenContract(ctx, valAddr)

	lpTokenContract, found := h.keeper.GetLPTokenContract(ctx, valAddr)
	if !found {
		// todo - is need panic? if not found, it means that the validator has been removed
		return
	}

	lpToken := fxtypes.GetLPToken().ABI
	data, err := lpToken.Pack("selfdestruct", common.BytesToAddress(h.keeper.lpTokenModuleAddress.Bytes()))
	if err != nil {
		return
	}

	err = h.keeper.callEVM(ctx, &lpTokenContract, data)
	if err != nil {
		return
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
