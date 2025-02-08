package keeper

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/hashicorp/go-metrics"

	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (k Keeper) BridgeCallExecuted(ctx sdk.Context, caller contract.Caller, msg *types.MsgBridgeCallClaim) error {
	k.CreateBridgeAccount(ctx, msg.TxOrigin)
	if senderAccount := k.ak.GetAccount(ctx, msg.GetSenderAddr().Bytes()); senderAccount != nil {
		if _, ok := senderAccount.(sdk.ModuleAccountI); ok {
			return types.ErrInvalid.Wrap("sender is module account")
		}
	}
	receiverAddr := msg.GetReceiverAddr()

	baseCoins := sdk.NewCoins()
	for i, tokenAddr := range msg.TokenContracts {
		baseCoin, err := k.DepositBridgeTokenToBaseCoin(ctx, receiverAddr.Bytes(), msg.Amounts[i], tokenAddr)
		if err != nil {
			return err
		}
		baseCoins = baseCoins.Add(baseCoin)
	}

	tokens := make([]common.Address, 0, baseCoins.Len())
	amounts := make([]*big.Int, 0, baseCoins.Len())
	for _, coin := range baseCoins {
		tokenContract, err := k.erc20Keeper.BaseCoinToEvm(ctx, caller, receiverAddr, coin)
		if err != nil {
			return err
		}
		tokens = append(tokens, common.HexToAddress(tokenContract))
		amounts = append(amounts, coin.Amount.BigInt())
	}
	err := k.HandlerBridgeCallInFee(ctx, caller, msg.GetSenderAddr(), msg.QuoteId.BigInt(), msg.GasLimit.Uint64())
	if err != nil {
		return err
	}

	if !k.evmKeeper.IsContract(ctx, msg.GetToAddr()) {
		return nil
	}

	var callEvmSender common.Address
	var args []byte
	if msg.IsMemoSendCallTo() {
		args = msg.MustData()
		callEvmSender = msg.GetSenderAddr()
	} else {
		args, err = contract.PackOnBridgeCall(msg.GetSenderAddr(), msg.GetRefundAddr(), tokens, amounts, msg.MustData(), msg.MustMemo())
		if err != nil {
			return err
		}
		callEvmSender = k.GetCallbackFrom()
	}

	cacheCtx, commit := sdk.UnwrapSDKContext(ctx).CacheContext()
	err = k.BridgeCallEvm(cacheCtx, caller, callEvmSender, msg.GetToAddr(), args, msg.GetGasLimit())
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
	revertMsg := err.Error()

	// refund bridge-call case of error
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCallEvent,
		sdk.NewAttribute(types.AttributeKeyErrCause, err.Error()),
	))

	erc20TokenKeeper := contract.NewERC20TokenKeeper(caller)
	for i, tokenAddr := range tokens {
		_, err = erc20TokenKeeper.Transfer(ctx, tokenAddr, receiverAddr, msg.GetRefundAddr(), amounts[i])
		if err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCallFailed,
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprintf("%d", msg.EventNonce)),
		sdk.NewAttribute(types.AttributeKeyBridgeCallFailedRefundAddr, msg.GetRefundAddr().Hex()),
	))

	// onRevert bridgeCall
	_, err = k.AddOutgoingBridgeCall(ctx, k.GetCallbackFrom(), common.Address{}, sdk.Coins{},
		msg.GetSenderAddr(), []byte(revertMsg), []byte{}, 0, msg.EventNonce)
	return err
}

func (k Keeper) BridgeCallEvm(ctx sdk.Context, caller contract.Caller, sender, to common.Address, args []byte, gasLimit uint64) error {
	if gasLimit == 0 {
		gasLimit = k.GetBridgeCallMaxGasLimit(ctx)
	}
	txResp, err := caller.ExecuteEVM(ctx, sender, &to, nil, gasLimit, args)
	if err != nil {
		return err
	}
	if txResp.Failed() {
		errStr := txResp.VmError
		if txResp.VmError == vm.ErrExecutionReverted.Error() {
			if cause, unpackErr := abi.UnpackRevert(common.CopyBytes(txResp.Ret)); unpackErr == nil {
				errStr = cause
			}
		}
		return evmtypes.ErrVMExecution.Wrap(errStr)
	}
	return nil
}

func (k Keeper) CreateBridgeAccount(ctx sdk.Context, address string) {
	accAddress := fxtypes.ExternalAddrToAccAddr(k.moduleName, address)
	if k.ak.HasAccount(ctx, accAddress) {
		return
	}
	k.ak.SetAccount(ctx, k.ak.NewAccountWithAddress(ctx, accAddress))
}
