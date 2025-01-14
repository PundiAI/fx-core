package keeper

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-metrics"

	"github.com/pundiai/fx-core/v8/contract"
	fxtelemetry "github.com/pundiai/fx-core/v8/telemetry"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (k Keeper) AddOutgoingBridgeCall(ctx sdk.Context, sender, refundAddr common.Address, baseCoins sdk.Coins, to common.Address, data, memo []byte, gasLimit, eventNonce uint64) (uint64, error) {
	tokens := make([]types.ERC20Token, 0, len(baseCoins))
	for _, baseCoin := range baseCoins {
		bridgeToken, err := k.BaseCoinToBridgeToken(ctx, sender.Bytes(), baseCoin)
		if err != nil {
			return 0, err
		}
		if err = k.WithdrawBridgeToken(ctx, sender.Bytes(), baseCoin.Amount, bridgeToken); err != nil {
			return 0, err
		}
		tokens = append(tokens, types.NewERC20Token(baseCoin.Amount, bridgeToken.Contract))
	}
	outCall, err := k.BuildOutgoingBridgeCall(ctx, sender, refundAddr, tokens, to, data, memo, gasLimit, eventNonce)
	if err != nil {
		return 0, err
	}
	return k.AddOutgoingBridgeCallWithoutBuild(ctx, outCall), nil
}

func (k Keeper) AddOutgoingBridgeCallWithoutBuild(ctx sdk.Context, outCall *types.OutgoingBridgeCall) uint64 {
	k.SetOutgoingBridgeCall(ctx, outCall)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCall,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, outCall.Sender),
		sdk.NewAttribute(types.AttributeKeyBridgeCallNonce, fmt.Sprint(outCall.Nonce)),
	))

	if !ctx.IsCheckTx() {
		for _, t := range outCall.Tokens {
			fxtelemetry.SetGaugeLabelsWithDenom(
				[]string{types.ModuleName, "bridge_call_out_amount"},
				t.Contract, t.Amount.BigInt(),
				telemetry.NewLabel("module", k.moduleName),
			)
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, "bridge_call_out"},
				float32(1),
				[]metrics.Label{
					telemetry.NewLabel("module", k.moduleName),
					telemetry.NewLabel("contract", t.Contract),
				},
			)
		}
	}

	return outCall.Nonce
}

func (k Keeper) BuildOutgoingBridgeCall(ctx sdk.Context, sender, refundAddr common.Address, tokens []types.ERC20Token, to common.Address, data, memo []byte, gasLimit, eventNonce uint64) (*types.OutgoingBridgeCall, error) {
	bridgeCallTimeout := k.CalExternalTimeoutHeight(ctx, GetBridgeCallTimeout)
	if bridgeCallTimeout <= 0 {
		return nil, types.ErrInvalid.Wrapf("bridge call timeout height")
	}

	nextID := k.autoIncrementID(ctx, types.KeyLastBridgeCallID)

	outCall := &types.OutgoingBridgeCall{
		Nonce:       nextID,
		Timeout:     bridgeCallTimeout,
		BlockHeight: uint64(ctx.BlockHeight()),
		Sender:      fxtypes.ExternalAddrToStr(k.moduleName, sender.Bytes()),
		Refund:      fxtypes.ExternalAddrToStr(k.moduleName, refundAddr.Bytes()),
		Tokens:      tokens,
		To:          fxtypes.ExternalAddrToStr(k.moduleName, to.Bytes()),
		Data:        hex.EncodeToString(data),
		Memo:        hex.EncodeToString(memo),
		GasLimit:    gasLimit,
		EventNonce:  eventNonce,
	}
	return outCall, nil
}

func (k Keeper) BridgeCallResultExecuted(ctx sdk.Context, caller contract.Caller, claim *types.MsgBridgeCallResultClaim) error {
	k.CreateBridgeAccount(ctx, claim.TxOrigin)

	outgoingBridgeCall, found := k.GetOutgoingBridgeCallByNonce(ctx, claim.Nonce)
	if !found {
		return fmt.Errorf("bridge call not found for nonce %d", claim.Nonce)
	}
	if !claim.Success {
		if err := k.RefundOutgoingBridgeCall(ctx, caller, outgoingBridgeCall); err != nil {
			return err
		}

		if err := k.BridgeCallOnRevert(ctx, caller, claim.Nonce, outgoingBridgeCall.To, claim.Cause); err != nil {
			return err
		}
	}
	if err := k.DeleteOutgoingBridgeCallRecord(ctx, claim.Nonce); err != nil {
		return err
	}

	if err := k.TransferBridgeFeeToRelayer(ctx, caller, claim.Nonce); err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCallResult,
		sdk.NewAttribute(types.AttributeKeyEventNonce, strconv.FormatInt(int64(claim.Nonce), 10)),
		sdk.NewAttribute(types.AttributeKeyStateSuccess, strconv.FormatBool(claim.Success)),
		sdk.NewAttribute(types.AttributeKeyErrCause, claim.Cause),
	))
	return nil
}

func (k Keeper) BridgeCallOnRevert(ctx sdk.Context, caller contract.Caller, nonce uint64, toAddr, cause string) error {
	args, err := contract.PackOnRevert(nonce, []byte(cause))
	if err != nil {
		return err
	}
	callEvmSender := k.GetCallbackFrom()
	gasLimit := k.GetBridgeCallMaxGasLimit(ctx)
	to := common.HexToAddress(toAddr)

	txResp, err := caller.ExecuteEVM(ctx, callEvmSender, &to, nil, gasLimit, args)
	if err != nil {
		return err
	}
	if txResp.Failed() {
		return types.ErrInvalid.Wrap(txResp.VmError)
	}
	return nil
}

func (k Keeper) RefundOutgoingBridgeCall(ctx sdk.Context, caller contract.Caller, data *types.OutgoingBridgeCall) error {
	refund := fxtypes.ExternalAddrToAccAddr(k.moduleName, data.GetRefund())
	baseCoins := sdk.NewCoins()
	for _, token := range data.Tokens {
		bridgeToken, err := k.DepositBridgeToken(ctx, refund.Bytes(), token.Amount, token.Contract)
		if err != nil {
			return err
		}
		baseCoin, err := k.BridgeTokenToBaseCoin(ctx, refund.Bytes(), token.Amount, bridgeToken)
		if err != nil {
			return err
		}
		baseCoins = baseCoins.Add(baseCoin)
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCallRefund,
		sdk.NewAttribute(types.AttributeKeyRefund, refund.String()),
	))

	originAmount, err := k.erc20Keeper.GetCache(ctx, types.NewOriginTokenKey(k.moduleName, data.Nonce))
	if err != nil && !errors.IsOf(err, collections.ErrNotFound) {
		return err
	}
	if errors.IsOf(err, collections.ErrNotFound) {
		originAmount = sdkmath.ZeroInt()
	}

	originCoin := sdk.NewCoin(fxtypes.DefaultDenom, originAmount)
	if !baseCoins.IsAllGTE(sdk.NewCoins(originCoin)) {
		return types.ErrInvalid.Wrapf("bridge call coin less than origin amount")
	}
	baseCoins = baseCoins.Sub(originCoin)

	for _, coin := range baseCoins {
		_, err = k.erc20Keeper.BaseCoinToEvm(ctx, caller, common.BytesToAddress(refund.Bytes()), coin)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) BridgeCallBaseCoin(ctx sdk.Context, caller contract.Caller, from, refund, to common.Address, coins sdk.Coins, data, memo []byte, quoteId, gasLimit *big.Int, fxTarget *types.FxTarget, originTokenAmount sdkmath.Int) (uint64, error) {
	var cacheKey string
	var nonce uint64
	if fxTarget.IsIBC() {
		if !coins.IsValid() || len(coins) != 1 {
			return 0, sdkerrors.ErrInvalidCoins.Wrapf("ibc transfer with coins: %s", coins.String())
		}
		baseCoin := coins[0]
		toAddr, err := fxTarget.ReceiveAddrToStr(to.Bytes())
		if err != nil {
			return 0, err
		}
		ibcCoin, err := k.BaseCoinToIBCCoin(ctx, from.Bytes(), baseCoin, fxTarget.IBCChannel)
		if err != nil {
			return 0, err
		}
		nonce, err = k.IBCTransfer(ctx, from.Bytes(), toAddr, ibcCoin, fxTarget.IBCChannel, string(memo))
		if err != nil {
			return 0, err
		}
		cacheKey = types.NewIBCTransferKey(fxTarget.IBCChannel, nonce)
	} else {
		var err error
		if nonce, err = k.AddOutgoingBridgeCall(ctx, from, refund, coins, to, data, memo, gasLimit.Uint64(), 0); err != nil {
			return 0, err
		}
		cacheKey = types.NewOriginTokenKey(k.moduleName, nonce)

		if err = k.HandlerBridgeCallOutFee(ctx, caller, from, nonce, quoteId, gasLimit.Uint64()); err != nil {
			return 0, err
		}
	}

	if originTokenAmount.IsPositive() {
		if err := k.erc20Keeper.SetCache(ctx, cacheKey, originTokenAmount); err != nil {
			return 0, err
		}
	}
	return nonce, nil
}

// Deprecated: precompile crosschain api has been deprecated
func (k Keeper) CrosschainBaseCoin(ctx sdk.Context, caller contract.Caller, from sdk.AccAddress, receipt string, amount, fee sdk.Coin, fxTarget *types.FxTarget, memo string, originToken bool) error {
	if fxTarget.IsIBC() {
		if !fee.IsZero() {
			return sdkerrors.ErrInvalidRequest.Wrap("ibc transfer fee must be zero")
		}
		sequence, err := k.IBCTransfer(ctx, from.Bytes(), receipt, amount, fxTarget.IBCChannel, memo)
		if err != nil {
			return err
		}
		if originToken {
			return k.erc20Keeper.SetCache(ctx, types.NewIBCTransferKey(fxTarget.IBCChannel, sequence), amount.Amount)
		}
	} else {
		if _, err := k.BuildOutgoingTxBatch(ctx, caller, from, receipt, amount, fee); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) IBCTransfer(ctx sdk.Context, from sdk.AccAddress, to string, amount sdk.Coin, channel, memo string) (uint64, error) {
	timeout := 12 * time.Hour
	transferResponse, err := k.ibcTransferKeeper.Transfer(ctx,
		transfertypes.NewMsgTransfer(
			transfertypes.ModuleName,
			channel,
			amount,
			from.String(),
			to,
			ibcclienttypes.ZeroHeight(),
			uint64(ctx.BlockTime().Add(timeout).UnixNano()),
			memo,
		),
	)
	if err != nil {
		return 0, fmt.Errorf("ibc transfer error: %s", err.Error())
	}
	return transferResponse.Sequence, nil
}

func (k Keeper) ResendBridgeCall(ctx sdk.Context, bridgeCall types.OutgoingBridgeCall, quoteInfo types.QuoteInfo) error {
	bridgeCallTimeout := k.CalExternalTimeoutHeight(ctx, GetBridgeCallTimeout)
	if bridgeCallTimeout <= 0 {
		return types.ErrInvalid.Wrapf("bridge call timeout height")
	}

	oldBridgeCallNonce := bridgeCall.Nonce
	k.DeleteOutgoingBridgeCallQuoteInfo(ctx, oldBridgeCallNonce)

	newBridgeCallNonce := k.autoIncrementID(ctx, types.KeyLastBridgeCallID)
	bridgeCall.Nonce = newBridgeCallNonce
	bridgeCall.Timeout = bridgeCallTimeout
	k.SetOutgoingBridgeCall(ctx, &bridgeCall)
	k.SetOutgoingBridgeCallQuoteInfo(ctx, newBridgeCallNonce, quoteInfo)

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeBridgeCallResend,
		sdk.NewAttribute(types.AttributeKeyBridgeCallResendOldNonce, fmt.Sprintf("%d", oldBridgeCallNonce)),
		sdk.NewAttribute(types.AttributeKeyBridgeCallResendNewNonce, fmt.Sprintf("%d", newBridgeCallNonce))),
	)
	return k.erc20Keeper.ReSetCache(ctx, types.NewOriginTokenKey(k.moduleName, oldBridgeCallNonce), types.NewOriginTokenKey(k.moduleName, newBridgeCallNonce))
}
