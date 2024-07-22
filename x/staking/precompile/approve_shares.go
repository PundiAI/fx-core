package precompile

import (
	"errors"

	sdkmath "cosmossdk.io/math"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v7/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v7/x/staking/types"
)

func (c *Contract) ApproveShares(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("approve method not readonly")
	}
	var args ApproveSharesArgs
	if err := types.ParseMethodArgs(ApproveSharesMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}
	stateDB := evm.StateDB.(types.ExtStateDB)
	err := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		owner := contract.Caller()
		c.stakingKeeper.SetAllowance(ctx, args.GetValidator(), owner.Bytes(), args.Spender.Bytes(), args.Shares)

		if err := c.AddLog(evm, ApproveSharesEvent, []common.Hash{owner.Hash(), args.Spender.Hash()}, args.Validator, args.Shares); err != nil {
			return err
		}

		ApproveSharesEmitEvents(ctx, args.GetValidator(), owner.Bytes(), args.Spender.Bytes(), sdkmath.NewIntFromBigInt(args.Shares))
		return nil
	})
	if err != nil {
		return nil, err
	}

	return ApproveSharesMethod.Outputs.Pack(true)
}

func ApproveSharesEmitEvents(ctx sdk.Context, validator sdk.ValAddress, owner, spender sdk.AccAddress, shares sdkmath.Int) {
	if shares.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, evmtypes.ModuleName, "approve_shares")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", evmtypes.TypeMsgEthereumTx},
				float32(shares.Int64()),
				[]metrics.Label{telemetry.NewLabel("validator", validator.String())},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			fxstakingtypes.EventTypeApproveShares,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, validator.String()),
			sdk.NewAttribute(fxstakingtypes.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(fxstakingtypes.AttributeKeySpender, spender.String()),
			sdk.NewAttribute(fxstakingtypes.AttributeKeyShares, shares.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, evmtypes.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, owner.String()),
		),
	})
}
