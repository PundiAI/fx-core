package v3

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerr "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

type StakingKeeper interface {
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator
	Delegate(ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdk.Int, tokenSrc stakingtypes.BondStatus,
		validator stakingtypes.Validator, subtractAccount bool) (newShares sdk.Dec, err error)
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
	GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, found bool)
	RemoveDelegation(ctx sdk.Context, delegation stakingtypes.Delegation)
}

type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

func MigrateDepositToStaking(ctx sdk.Context, moduleName string, stakingKeeper StakingKeeper, bankKeeper BankKeeper,
	oracles types.Oracles, proposalOracle types.ProposalOracle, delegateValAddr sdk.ValAddress) error {
	if moduleName != bsctypes.ModuleName && moduleName != polygontypes.ModuleName && moduleName != trontypes.ModuleName {
		return fmt.Errorf("not support module name: %s", moduleName)
	}
	validator, found := stakingKeeper.GetValidator(ctx, delegateValAddr)
	if !found {
		return stakingtypes.ErrNoValidatorFound
	}
	logger := ctx.Logger().With("module", "x/"+moduleName)

	for i, oracle := range oracles {
		if i == 0 {
			continue
		}

		if validator.OperatorAddress != oracle.DelegateValidator {
			return sdkerr.Wrap(types.ErrInvalid, "delegate validator")
		}

		delegateAddr := oracle.GetDelegateAddress(moduleName)

		delegation, found := stakingKeeper.GetDelegation(ctx, delegateAddr, delegateValAddr)
		if !found {
			logger.Info("no found delegating on migrate", "delegate", delegateAddr, "validator", delegateValAddr.String())
			continue
		}
		stakingKeeper.RemoveDelegation(ctx, delegation)

		var sendName string
		switch {
		case validator.IsBonded():
			sendName = stakingtypes.BondedPoolName
		case validator.IsUnbonding(), validator.IsUnbonded():
			sendName = stakingtypes.NotBondedPoolName
		default:
			panic("invalid validator status")
		}

		if err := bankKeeper.SendCoinsFromModuleToAccount(ctx,
			sendName, delegateAddr, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, oracle.DelegateAmount))); err != nil {
			return err
		}
		if !oracle.Online {
			continue
		}

		newShares, err := stakingKeeper.Delegate(ctx, delegateAddr, oracle.DelegateAmount, stakingtypes.Unbonded, validator, true)
		if err != nil {
			return err
		}
		//notice: Each delegate should be followed by an update of the validator `Tokens` and `DelegatorShares`
		validator, _ = validator.AddTokensFromDel(oracle.DelegateAmount)

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
