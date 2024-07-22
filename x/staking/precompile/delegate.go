package precompile

import (
	"errors"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/evm/types"
	fxstakingkeeper "github.com/functionx/fx-core/v7/x/staking/keeper"
)

func (c *Contract) Delegate(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("delegate method not readonly")
	}
	var args DelegateArgs
	err := types.ParseMethodArgs(DelegateMethod, &args, contract.Input[4:])
	if err != nil {
		return nil, err
	}
	amount := contract.Value()
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("invalid delegate amount: %s", amount.String())
	}
	valAddr := args.GetValidator()

	stateDB := evm.StateDB.(types.ExtStateDB)

	var result []byte
	err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		val, found := c.stakingKeeper.GetValidator(ctx, valAddr)
		if !found {
			return fmt.Errorf("validator not found: %s", valAddr.String())
		}

		sender := sdk.AccAddress(contract.Caller().Bytes())
		evmDenom := c.evmKeeper.GetParams(ctx).EvmDenom
		delCoin := sdk.NewCoin(evmDenom, sdkmath.NewIntFromBigInt(amount))
		if err = c.bankKeeper.SendCoinsFromAccountToModule(ctx, contract.Address().Bytes(), evmtypes.ModuleName, sdk.NewCoins(delCoin)); err != nil {
			return err
		}
		if err = c.bankKeeper.SendCoinsFromModuleToAccount(ctx, evmtypes.ModuleName, sender, sdk.NewCoins(delCoin)); err != nil {
			return err
		}

		withdrawAddr := c.distrKeeper.GetDelegatorWithdrawAddr(ctx, sender)
		beforeDelBalance := c.bankKeeper.GetBalance(ctx, withdrawAddr, evmDenom)
		if withdrawAddr.Equals(sender) {
			beforeDelBalance = beforeDelBalance.Sub(delCoin)
		}

		// delegate amount
		shares, err := c.stakingKeeper.Delegate(ctx, sender, sdkmath.NewIntFromBigInt(amount), stakingtypes.Unbonded, val, true)
		if err != nil {
			return err
		}

		afterDelBalance := c.bankKeeper.GetBalance(ctx, withdrawAddr, evmDenom)
		rewardCoin := afterDelBalance.Sub(beforeDelBalance)

		// add delegate log
		if err = c.AddLog(evm, DelegateEvent, []common.Hash{contract.Caller().Hash()},
			args.Validator, amount, shares.TruncateInt().BigInt()); err != nil {
			return err
		}
		// add delegate event
		DelegateEmitEvents(ctx, sender, valAddr, amount, shares)

		result, err = DelegateMethod.Outputs.Pack(shares.TruncateInt().BigInt(), rewardCoin.Amount.BigInt())
		return err
	})
	return result, err
}

func (c *Contract) DelegateV2(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("delegate method not readonly")
	}
	var args DelegateV2Args
	err := types.ParseMethodArgs(DelegateV2Method, &args, contract.Input[4:])
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)

	var result []byte
	err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		sender := sdk.AccAddress(contract.Caller().Bytes())
		evmDenom := c.evmKeeper.GetParams(ctx).EvmDenom
		delCoin := sdk.NewCoin(evmDenom, sdkmath.NewIntFromBigInt(args.Amount))

		fxStakingKeeper, ok := c.stakingKeeper.(*fxstakingkeeper.Keeper)
		if !ok {
			return errortypes.ErrNotSupported
		}
		impl := stakingkeeper.NewMsgServerImpl(fxStakingKeeper.Keeper)
		if _, err = impl.Delegate(sdk.WrapSDKContext(ctx), &stakingtypes.MsgDelegate{
			DelegatorAddress: sender.String(),
			ValidatorAddress: args.Validator,
			Amount:           delCoin,
		}); err != nil {
			return err
		}

		// add delegate log
		if err = c.AddLog(evm, DelegateV2Event, []common.Hash{contract.Caller().Hash()}, args.Validator, args.Amount); err != nil {
			return err
		}
		result, err = DelegateV2Method.Outputs.Pack(true)
		return err
	})
	return result, err
}

func DelegateEmitEvents(ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress, amount *big.Int, newShares sdk.Dec) {
	if amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, evmtypes.ModuleName, "delegate")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", evmtypes.TypeMsgEthereumTx},
				float32(amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", fxtypes.DefaultDenom)},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeDelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, validator.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, evmtypes.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, delegator.String()),
		),
	})
}
