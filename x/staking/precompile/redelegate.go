package precompile

import (
	"errors"
	"fmt"
	"math/big"
	"time"

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

func (c *Contract) Redelegation(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("redelegate method not readonly")
	}
	var args RedelegateArgs
	err := types.ParseMethodArgs(RedelegateMethod, &args, contract.Input[4:])
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)

	var result []byte
	err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error { // withdraw src reward
		sender := sdk.AccAddress(contract.Caller().Bytes())
		evmDenom := c.evmKeeper.GetParams(ctx).EvmDenom
		shares := sdkmath.LegacyNewDecFromBigInt(args.Shares)

		valSrcAddr := args.GetValidatorSrc()
		// check src validator
		validatorSrc, found := c.stakingKeeper.GetValidator(ctx, valSrcAddr)
		if !found {
			return fmt.Errorf("validator src not found: %s", valSrcAddr.String())
		}

		// check delegation
		delegation, found := c.stakingKeeper.GetDelegation(ctx, sender, valSrcAddr)
		if !found {
			return fmt.Errorf("delegation not found")
		}
		if delegation.Shares.LT(shares) {
			return fmt.Errorf("insufficient shares to redelegate")
		}

		// check dst validator
		valDstAddr := args.GetValidatorDst()
		if _, found = c.stakingKeeper.GetValidator(ctx, valDstAddr); !found {
			return fmt.Errorf("validator dst not found: %s", valDstAddr.String())
		}

		withdrawAddr := c.distrKeeper.GetDelegatorWithdrawAddr(ctx, sender)
		beforeDelBalance := c.bankKeeper.GetBalance(ctx, withdrawAddr, evmDenom)

		// redelegate
		completionTime, err := c.stakingKeeper.BeginRedelegation(ctx, sender, valSrcAddr, valDstAddr, shares)
		if err != nil {
			return err
		}

		redelAmount := validatorSrc.TokensFromShares(shares).TruncateInt()
		afterDelBalance := c.bankKeeper.GetBalance(ctx, withdrawAddr, evmDenom)
		rewardCoin := afterDelBalance.Sub(beforeDelBalance)

		// add redelegate log
		if err = c.AddLog(evm, RedelegateEvent, []common.Hash{contract.Caller().Hash()},
			args.ValidatorSrc, args.ValidatorDst, args.Shares, redelAmount.BigInt(), big.NewInt(completionTime.Unix())); err != nil {
			return err
		}
		// add redelegate event
		RedelegateEmitEvents(ctx, sender, valSrcAddr, valDstAddr, redelAmount, completionTime)

		result, err = UndelegateMethod.Outputs.Pack(redelAmount.BigInt(), rewardCoin.Amount.BigInt(), big.NewInt(completionTime.Unix()))
		return err
	})
	return result, err
}

func (c *Contract) RedelegationV2(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("redelegate method not readonly")
	}
	var args RedelegateV2Args
	err := types.ParseMethodArgs(RedelegateV2Method, &args, contract.Input[4:])
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)

	var result []byte
	err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error { // withdraw src reward
		sender := sdk.AccAddress(contract.Caller().Bytes())
		evmDenom := c.evmKeeper.GetParams(ctx).EvmDenom
		reDelCoin := sdk.NewCoin(evmDenom, sdkmath.NewIntFromBigInt(args.Amount))

		fxStakingKeeper, ok := c.stakingKeeper.(fxstakingkeeper.Keeper)
		if !ok {
			return errortypes.ErrNotSupported
		}
		impl := stakingkeeper.NewMsgServerImpl(fxStakingKeeper.Keeper)
		redelResp, err := impl.BeginRedelegate(sdk.WrapSDKContext(ctx), &stakingtypes.MsgBeginRedelegate{
			DelegatorAddress:    sender.String(),
			ValidatorSrcAddress: args.ValidatorSrc,
			ValidatorDstAddress: args.ValidatorDst,
			Amount:              reDelCoin,
		})
		if err != nil {
			return err
		}

		// add redelegate log
		if err = c.AddLog(evm, RedelegateV2Event, []common.Hash{contract.Caller().Hash()},
			args.ValidatorSrc, args.ValidatorDst, args.Amount, big.NewInt(redelResp.CompletionTime.Unix())); err != nil {
			return err
		}

		result, err = UndelegateMethod.Outputs.Pack(true)
		return err
	})
	return result, err
}

func RedelegateEmitEvents(ctx sdk.Context, delegator sdk.AccAddress, validatorSrc, validatorDst sdk.ValAddress, amount sdkmath.Int, completionTime time.Time) {
	if amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, evmtypes.ModuleName, "redelegate")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", evmtypes.TypeMsgEthereumTx},
				float32(amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", fxtypes.DefaultDenom)},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeRedelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeySrcValidator, validatorSrc.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyDstValidator, validatorDst.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, evmtypes.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, delegator.String()),
		),
	})
}
