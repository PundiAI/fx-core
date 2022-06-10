package v045

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/functionx/fx-core/x/crosschain/types"
)

type StakingKeeper interface {
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator
	Delegate(ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdk.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (newShares sdk.Dec, err error)
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
}

func MigrateDepositToStaking(ctx sdk.Context, oracles types.Oracles, stakingKeeper StakingKeeper) error {
	validators := stakingKeeper.GetBondedValidatorsByPower(ctx)
	if len(validators) <= 0 {
		panic("no found bonded validator")
	}
	validator := validators[0]

	for _, oracle := range oracles {

		oracleAddr, err := sdk.AccAddressFromBech32(oracle.OracleAddress)
		if err != nil {
			return err
		}

		oracle.DelegateValidator = validator.OperatorAddress
		oracle.IsValidator = false

		newShares, err := stakingKeeper.Delegate(ctx,
			oracleAddr, oracle.DelegateAmount, stakingtypes.Unbonded, validator, true)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				stakingtypes.EventTypeDelegate,
				sdk.NewAttribute(stakingtypes.AttributeKeyValidator, validator.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, oracle.DelegateAmount.String()),
				sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
			),
		})
	}
	return nil
}
