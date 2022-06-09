package v045

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	fxtypes "github.com/functionx/fx-core/types"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
)

type StakingKeeper interface {
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator
	Delegate(ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdk.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (newShares sdk.Dec, err error)
}

func MigrateDepositToStaking(ctx sdk.Context, oracles crosschaintypes.Oracles, stakingKeeper StakingKeeper) error {

	validators := stakingKeeper.GetBondedValidatorsByPower(ctx)

	for i, oracle := range oracles {

		oracleAddr, err := sdk.AccAddressFromBech32(oracle.OracleAddress)
		if err != nil {
			return err
		}

		if oracle.DelegateAmount.Denom != fxtypes.DefaultDenom {
			return fmt.Errorf("invalid stake amount: %s", oracle.DelegateAmount.String())
		}
		validator := validators[len(validators)%(i+1)]

		newShares, err := stakingKeeper.Delegate(ctx,
			oracleAddr, oracle.DelegateAmount.Amount, stakingtypes.Unbonded, validator, true)
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
