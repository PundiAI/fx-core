package keeper

import (
	"fmt"
	"math/big"
	"strconv"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-metrics"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (k Keeper) BridgeCallHandler(ctx sdk.Context, msg *types.MsgBridgeCallClaim) error {
	k.CreateBridgeAccount(ctx, msg.TxOrigin)
	if senderAccount := k.ak.GetAccount(ctx, msg.GetSenderAddr().Bytes()); senderAccount != nil {
		if _, ok := senderAccount.(sdk.ModuleAccountI); ok {
			return types.ErrInvalid.Wrap("sender is module account")
		}
	}
	isMemoSendCallTo := msg.IsMemoSendCallTo()
	receiverAddr := msg.GetToAddr()
	if isMemoSendCallTo {
		receiverAddr = msg.GetSenderAddr()
	}

	baseCoins := sdk.NewCoins()
	for i, tokenAddr := range msg.TokenContracts {
		bridgeToken, err := k.DepositBridgeToken(ctx, receiverAddr.Bytes(), msg.Amounts[i], tokenAddr)
		if err != nil {
			return err
		}
		baseCoin, err := k.BridgeTokenToBaseCoin(ctx, receiverAddr.Bytes(), msg.Amounts[i], bridgeToken)
		if err != nil {
			return err
		}
		baseCoins = baseCoins.Add(baseCoin)
	}

	cacheCtx, commit := sdk.UnwrapSDKContext(ctx).CacheContext()
	err := k.BridgeCallEvm(cacheCtx, msg.GetSenderAddr(), msg.GetRefundAddr(), msg.GetToAddr(),
		receiverAddr, baseCoins, msg.MustData(), msg.MustMemo(), msg.Value, isMemoSendCallTo)
	if !ctx.IsCheckTx() {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "bridge_call_in"},
			float32(1),
			[]metrics.Label{
				telemetry.NewLabel("module", k.moduleName),
				telemetry.NewLabel("success", strconv.FormatBool(err == nil)),
			},
		)
	}
	if err == nil {
		commit()
		return nil
	}
	// refund bridge-call case of error
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeBridgeCallEvent, sdk.NewAttribute(types.AttributeKeyErrCause, err.Error())))

	refundAddr := msg.GetRefundAddr()
	outCallNonce, err := k.AddOutgoingBridgeCall(ctx, refundAddr, refundAddr, baseCoins, common.Address{}, nil, nil, msg.EventNonce)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCallRefundOut,
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprintf("%d", msg.EventNonce)),
		sdk.NewAttribute(types.AttributeKeyBridgeCallNonce, fmt.Sprintf("%d", outCallNonce)),
	))
	return nil
}

func (k Keeper) BridgeCallEvm(ctx sdk.Context, sender, refundAddr, to, receiverAddr common.Address, baseCoins sdk.Coins, data, memo []byte, value sdkmath.Int, isMemoSendCallTo bool) error {
	tokens := make([]common.Address, 0, baseCoins.Len())
	amounts := make([]*big.Int, 0, baseCoins.Len())
	for _, coin := range baseCoins {
		tokenContract, err := k.erc20Keeper.BaseCoinToEvm(ctx, receiverAddr, coin)
		if err != nil {
			return err
		}
		tokens = append(tokens, common.HexToAddress(tokenContract))
		amounts = append(amounts, coin.Amount.BigInt())
	}

	if !k.evmKeeper.IsContract(ctx, to) {
		return nil
	}
	var callEvmSender common.Address
	var args []byte

	if isMemoSendCallTo {
		args = data
		callEvmSender = sender
	} else {
		var err error
		args, err = contract.PackBridgeCallback(sender, refundAddr, tokens, amounts, data, memo)
		if err != nil {
			return err
		}
		callEvmSender = k.GetCallbackFrom()
	}

	gasLimit := k.GetBridgeCallMaxGasLimit(ctx)
	txResp, err := k.evmKeeper.ExecuteEVM(ctx, callEvmSender, &to, value.BigInt(), gasLimit, args)
	if err != nil {
		return err
	}
	if txResp.Failed() {
		return types.ErrInvalid.Wrap(txResp.VmError)
	}
	return nil
}
