package keeper

import (
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

	attr := sdk.NewAttribute(types.AttributeKeyIBCCallType, types.IbcCallType_name[int32(mp.GetType())])
	ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(types.EventTypeIBCCall, attr)})
	switch packet := mp.(type) {
	case *types.IbcCallEvmPacket:
		hexSender := types.IntermediateSender(sourcePort, sourceChannel, data.Sender)
		return k.HandlerIbcCallEvm(ctx, hexSender, packet)
	default:
		return types.ErrMemoNotSupport.Wrapf("invalid call type %s", mp.GetType())
	}
}

func (k Keeper) HandlerIbcCallEvm(ctx sdk.Context, sender common.Address, evmData *types.IbcCallEvmPacket) error {
	txResp, err := k.evmKeeper.CallEVM(ctx, sender,
		evmData.MustGetToAddr(), evmData.Value.BigInt(), evmData.GasLimit, evmData.MustGetMessage(), true)
	if err != nil {
		return err
	}
	if txResp.Failed() {
		return errorsmod.Wrap(evmtypes.ErrVMExecution, txResp.VmError)
	}
	return nil
}
