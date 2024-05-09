package keeper

import (
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v7/x/ibc/applications/transfer/types"
)

func (k Keeper) HandlerIbcCall(ctx sdk.Context, sourcePort, sourceChannel string, data types.FungibleTokenPacketData) error {
	var mp types.MemoPacket
	if err := k.cdc.UnmarshalInterfaceJSON([]byte(data.Memo), &mp); err != nil {
		return nil
	}

	if err := mp.ValidateBasic(); err != nil {
		return err
	}

	switch packet := mp.(type) {
	case *types.IbcCallEvmPacket:
		hexSender := types.IntermediateSender(sourcePort, sourceChannel, data.Sender)
		return k.HandlerIbcCallEvm(ctx, hexSender, packet)
	default:
		return errorsmod.Wrapf(types.ErrMemoNotSupport, "invalid call type %s", mp.GetType())
	}
}

func (k Keeper) HandlerIbcCallEvm(ctx sdk.Context, sender common.Address, evmPacket *types.IbcCallEvmPacket) error {
	limit := ctx.ConsensusParams().GetBlock().GetMaxGas()
	evmErr, evmResult := "", false
	defer func() {
		attrs := []sdk.Attribute{
			sdk.NewAttribute(types.AttributeKeyIBCCallType, types.IbcCallType_name[int32(evmPacket.GetType())]),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyIBCCallResult, strconv.FormatBool(evmResult)),
		}
		if len(evmErr) > 0 {
			attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyIBCCallError, evmErr))
		}
		ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(types.EventTypeIBCCall, attrs...)})
	}()
	txResp, err := k.evmKeeper.CallEVM(ctx, sender,
		evmPacket.GetToAddress(), evmPacket.Value.BigInt(), uint64(limit), evmPacket.MustGetMessage(), true)
	if err != nil {
		evmErr = err.Error()
		return err
	}
	evmResult = !txResp.Failed()
	evmErr = txResp.VmError
	if txResp.Failed() {
		return errorsmod.Wrap(evmtypes.ErrVMExecution, txResp.VmError)
	}
	return nil
}
