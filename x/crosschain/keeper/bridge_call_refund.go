package keeper

import (
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/contract"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) HandleOutgoingBridgeCallRefund(ctx sdk.Context, data *types.OutgoingBridgeCall) {
	sender := types.ExternalAddrToAccAddr(k.moduleName, data.GetSender())
	coins, err := k.bridgeCallTransferToSender(ctx, sender, data.Tokens)
	if err != nil {
		panic(err)
	}

	evmErrCause, evmSuccess, isCallback := "", false, false
	defer func() {
		attrs := []sdk.Attribute{
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		}
		if isCallback {
			attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyStateSuccess, strconv.FormatBool(evmSuccess)))
			if len(evmErrCause) > 0 {
				attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyErrCause, evmErrCause))
			}
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeBridgeCallRefund,
			attrs...,
		))
	}()

	if k.HasBridgeCallFromMsg(ctx, data.Nonce) {
		return
	}
	// precompile bridge call, refund to evm
	if err = k.bridgeCallTransferToReceiver(ctx, sender, sender, coins); err != nil {
		panic(err)
	}
	if data.EventNonce > 0 {
		contractAddr := common.BytesToAddress(sender.Bytes())
		account := k.evmKeeper.GetAccount(ctx, contractAddr)
		if !account.IsContract() {
			return
		}

		isCallback = true
		maxGasLimit := k.GetParams(ctx).BridgeCallMaxGasLimit
		tokens := types.ERC20Tokens(data.Tokens)
		args, err := contract.GetBridgeCallRefundCallback().Pack(
			"refundCallback",
			data.EventNonce,
			tokens.GetContracts(),
			tokens.GetAmounts(),
		)
		if err != nil {
			evmErrCause = err.Error()
			return
		}
		txResp, err := k.evmKeeper.CallEVM(
			ctx,
			k.callbackFrom,
			&contractAddr,
			big.NewInt(0),
			maxGasLimit,
			args,
			true,
		)
		if err != nil {
			evmErrCause = err.Error()
		} else {
			evmSuccess = !txResp.Failed()
			evmErrCause = txResp.VmError
		}
	}
}

func (k Keeper) DeleteOutgoingBridgeCallRecord(ctx sdk.Context, bridgeCallNonce uint64) {
	// 1. delete bridge call
	k.DeleteOutgoingBridgeCall(ctx, bridgeCallNonce)

	// 2. delete bridge call confirm
	k.DeleteBridgeCallConfirm(ctx, bridgeCallNonce)

	// 3. delete bridge call from msg
	k.DeleteBridgeCallFromMsg(ctx, bridgeCallNonce)
}
