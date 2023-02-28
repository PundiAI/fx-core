package v2

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

type StakingKeeper interface {
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator
	Delegate(ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdk.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (newShares sdk.Dec, err error)
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
}

type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

func MigrateDepositToStaking(ctx sdk.Context, moduleName string, stakingKeeper StakingKeeper, bankKeeper BankKeeper, oracles types.Oracles, delegateValAddr sdk.ValAddress) error {
	validator, found := stakingKeeper.GetValidator(ctx, delegateValAddr)
	if !found {
		return stakingtypes.ErrNoValidatorFound
	}

	for i, oracle := range oracles {
		if validator.OperatorAddress != oracle.DelegateValidator {
			return errorsmod.Wrap(types.ErrInvalid, "delegate validator")
		}

		delegateAddr := oracle.GetDelegateAddress(moduleName)

		if err := bankKeeper.SendCoinsFromModuleToAccount(ctx,
			moduleName, delegateAddr, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, oracle.DelegateAmount))); err != nil {
			return err
		}

		newShares, err := stakingKeeper.Delegate(ctx,
			delegateAddr, oracle.DelegateAmount, stakingtypes.Unbonded, validator, true)
		if err != nil {
			return err
		}
		// notice: Each delegate should be followed by an update of the validator `Tokens` and `DelegatorShares`
		// validator, _ = validator.AddTokensFromDel(oracle.DelegateAmount)

		if i != 0 {
			continue
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			stakingtypes.EventTypeDelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, oracle.DelegateValidator),
			sdk.NewAttribute(sdk.AttributeKeyAmount, oracle.DelegateAmount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
		))
	}
	return nil
}
