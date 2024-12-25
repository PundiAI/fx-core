package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/pundiai/fx-core/v8/x/ibc/middleware/types"
)

func (k Keeper) HandlerIbcCall(ctx sdk.Context, sourcePort, sourceChannel string, data transfertypes.FungibleTokenPacketData) error {
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
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid call type %s", mp.GetType())
	}
}

func (k Keeper) HandlerIbcCallEvm(ctx sdk.Context, sender common.Address, evmPacket *types.IbcCallEvmPacket) error {
	limit := ctx.ConsensusParams().Block.GetMaxGas()
	evmErrCause, evmSuccess := "", false
	defer func() {
		attrs := []sdk.Attribute{
			sdk.NewAttribute(types.AttributeKeyIBCCallType, types.IbcCallType_name[int32(evmPacket.GetType())]),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyIBCCallSuccess, strconv.FormatBool(evmSuccess)),
		}
		if len(evmErrCause) > 0 {
			attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyIBCCallErrCause, evmErrCause))
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeIBCCall, attrs...))
	}()
	txResp, err := k.evmKeeper.ExecuteEVM(ctx, sender,
		evmPacket.GetToAddress(), evmPacket.Value.BigInt(), uint64(limit), evmPacket.MustGetData())
	if err != nil {
		evmErrCause = err.Error()
		return err
	}
	evmSuccess = !txResp.Failed()
	evmErrCause = txResp.VmError
	if txResp.Failed() {
		return evmtypes.ErrVMExecution.Wrap(txResp.VmError)
	}
	return nil
}
