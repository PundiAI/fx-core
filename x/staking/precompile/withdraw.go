package precompile

import (
	"errors"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v7/x/evm/types"
)

func (c *Contract) Withdraw(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("withdraw method not readonly")
	}

	var args WithdrawArgs
	if err := types.ParseMethodArgs(WithdrawMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)
	var result []byte
	err := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		evmDenom := c.evmKeeper.GetParams(ctx).EvmDenom
		sender := sdk.AccAddress(contract.Caller().Bytes())

		impl := distrkeeper.NewMsgServerImpl(c.distrKeeper.(distrkeeper.Keeper))
		rewardResp, err := impl.WithdrawDelegatorReward(sdk.WrapSDKContext(ctx), &distrtypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: sender.String(),
			ValidatorAddress: args.GetValidator().String(),
		})
		if err != nil {
			return err
		}
		// add withdraw event
		WithdrawEmitEvents(ctx, sender, rewardResp.Amount)

		// add withdraw log
		if err = c.AddLog(evm, WithdrawEvent, []common.Hash{contract.Caller().Hash()}, args.GetValidator().String(), rewardResp.Amount.AmountOf(evmDenom).BigInt()); err != nil {
			return err
		}

		result, err = WithdrawMethod.Outputs.Pack(rewardResp.Amount.AmountOf(evmDenom).BigInt())
		return err
	})

	return result, err
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
