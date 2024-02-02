package staking

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/evm/types"
)

func (c *Contract) Redelegation(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("redelegate method not readonly")
	}
	// parse args
	var args RedelegateArgs
	err := types.ParseMethodArgs(RedelegateMethod, &args, contract.Input[4:])
	if err != nil {
		return nil, err
	}
	sender := sdk.AccAddress(contract.Caller().Bytes())
	evmDenom := c.evmKeeper.GetParams(ctx).EvmDenom

	// withdraw rewards if delegation exist, add reward to evm state balance
	reward := big.NewInt(0)

	valSrcAddr := args.GetValidatorSrc()
	shares := sdk.NewDecFromBigInt(args.Shares)

	// check src validator
	validatorSrc, found := c.stakingKeeper.GetValidator(ctx, valSrcAddr)
	if !found {
		return nil, fmt.Errorf("validator src not found: %s", valSrcAddr.String())
	}

	// check delegation
	delegation, found := c.stakingKeeper.GetDelegation(ctx, sender, valSrcAddr)
	if !found {
		return nil, fmt.Errorf("delegation not found")
	}
	if delegation.Shares.LT(shares) {
		return nil, fmt.Errorf("insufficient shares to redelegate")
	}
	// withdraw src reward
	rewardSrc, err := c.withdraw(ctx, evm, contract.Caller(), valSrcAddr, evmDenom)
	if err != nil {
		return nil, err
	}
	reward = big.NewInt(0).Add(reward, rewardSrc)

	// check dst validator
	valDstAddr := args.GetValidatorDst()
	if _, found = c.stakingKeeper.GetValidator(ctx, valDstAddr); !found {
		return nil, fmt.Errorf("validator dst not found: %s", valDstAddr.String())
	}
	// withdraw dst reward
	if _, found = c.stakingKeeper.GetDelegation(ctx, sender, valDstAddr); found {
		rewardDst, err := c.withdraw(ctx, evm, contract.Caller(), valDstAddr, evmDenom)
		if err != nil {
			return nil, err
		}
		reward = big.NewInt(0).Add(reward, rewardDst)
	}

	// redelegate
	completionTime, err := c.stakingKeeper.BeginRedelegation(ctx, sender, valSrcAddr, valDstAddr, shares)
	if err != nil {
		return nil, err
	}
	redelAmount := validatorSrc.TokensFromShares(shares).TruncateInt()
	// add redelegate log
	if err := c.AddLog(evm, RedelegateEvent, []common.Hash{contract.Caller().Hash()},
		args.ValidatorSrc, args.ValidatorDst, args.Shares, redelAmount.BigInt(), big.NewInt(completionTime.Unix())); err != nil {
		return nil, err
	}
	// add redelegate event
	RedelegateEmitEvents(ctx, sender, valSrcAddr, valDstAddr, redelAmount, completionTime)

	return UndelegateMethod.Outputs.Pack(redelAmount.BigInt(), reward, big.NewInt(completionTime.Unix()))
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
