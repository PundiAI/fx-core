package staking

import (
	"errors"
	"math/big"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v7/x/evm/types"
)

func (c *Contract) Withdraw(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("withdraw method not readonly")
	}
	// parse args
	var args WithdrawArgs
	if err := types.ParseMethodArgs(WithdrawMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}

	evmDenom := c.evmKeeper.GetParams(ctx).EvmDenom
	rewardAmount, err := c.withdraw(ctx, evm, contract.Caller(), args.GetValidator(), evmDenom)
	if err != nil {
		return nil, err
	}
	return WithdrawMethod.Outputs.Pack(rewardAmount)
}

func (c *Contract) withdraw(ctx sdk.Context, evm *vm.EVM, sender common.Address, valAddr sdk.ValAddress, withDrawDenom string) (*big.Int, error) {
	delAddr := sdk.AccAddress(sender.Bytes())
	withdrawAddr := c.distrKeeper.GetDelegatorWithdrawAddr(ctx, delAddr)
	// cache withdraw address state, before withdraw rewards
	evm.StateDB.GetBalance(common.BytesToAddress(withdrawAddr.Bytes()))

	rewards, err := c.distrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
	if err != nil {
		return nil, err
	}

	// add withdraw event
	WithdrawEmitEvents(ctx, delAddr, rewards)

	rewardAmount := rewards.AmountOf(withDrawDenom).BigInt()
	if rewardAmount.Cmp(big.NewInt(0)) == 0 {
		return rewardAmount, nil
	}

	if withdrawAddr.Equals(delAddr) {
		evm.StateDB.AddBalance(sender, rewardAmount)
	} else {
		evm.StateDB.AddBalance(common.BytesToAddress(withdrawAddr.Bytes()), rewardAmount)
	}

	// add withdraw log
	if err := c.AddLog(evm, WithdrawEvent, []common.Hash{sender.Hash()}, valAddr.String(), rewardAmount); err != nil {
		return nil, err
	}

	return rewardAmount, nil
}

func WithdrawEmitEvents(ctx sdk.Context, delegator sdk.AccAddress, amount sdk.Coins) {
	defer func() {
		for _, a := range amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "withdraw_reward"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, evmtypes.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, delegator.String()),
		),
	)
}
