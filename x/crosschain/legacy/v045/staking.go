package v045

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerr "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/x/crosschain/types"
)

type StakingKeeper interface {
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator
	Delegate(ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdk.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (newShares sdk.Dec, err error)
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
}

func MigrateDepositToStaking(ctx sdk.Context, moduleName string, stakingKeeper StakingKeeper, oracles types.Oracles, delegateValidator stakingtypes.Validator) error {
	for _, oracle := range oracles {
		if delegateValidator.OperatorAddress != oracle.DelegateValidator {
			return sdkerr.Wrap(types.ErrInvalid, "delegate validator")
		}

		delegateAddress := oracle.GetDelegateAddress(moduleName)
		newShares, err := stakingKeeper.Delegate(ctx,
			delegateAddress, oracle.DelegateAmount, stakingtypes.Unbonded, delegateValidator, true)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				stakingtypes.EventTypeDelegate,
				sdk.NewAttribute(stakingtypes.AttributeKeyValidator, oracle.DelegateValidator),
				sdk.NewAttribute(sdk.AttributeKeyAmount, oracle.DelegateAmount.String()),
				sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
			),
		})
	}
	return nil
}
