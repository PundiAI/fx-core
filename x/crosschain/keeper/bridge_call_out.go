package keeper

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"time"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/cosmos/gogoproto/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-metrics"

	fxtelemetry "github.com/functionx/fx-core/v8/telemetry"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (k Keeper) AddOutgoingBridgeCall(ctx sdk.Context, sender, refundAddr common.Address, baseCoins sdk.Coins, to common.Address, data, memo []byte, eventNonce uint64) (uint64, error) {
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
	outCall, err := k.BuildOutgoingBridgeCall(ctx, sender, refundAddr, tokens, to, data, memo, eventNonce)
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

func (k Keeper) BuildOutgoingBridgeCall(ctx sdk.Context, sender common.Address, refundAddr common.Address, tokens []types.ERC20Token, to common.Address, data []byte, memo []byte, eventNonce uint64) (*types.OutgoingBridgeCall, error) {
	bridgeCallTimeout := k.CalExternalTimeoutHeight(ctx, GetBridgeCallTimeout)
	if bridgeCallTimeout <= 0 {
		return nil, types.ErrInvalid.Wrapf("bridge call timeout height")
	}

	nextID := k.autoIncrementID(ctx, types.KeyLastBridgeCallID)

	outCall := &types.OutgoingBridgeCall{
		Nonce:       nextID,
		Timeout:     bridgeCallTimeout,
		BlockHeight: uint64(ctx.BlockHeight()),
		Sender:      types.ExternalAddrToStr(k.moduleName, sender.Bytes()),
		Refund:      types.ExternalAddrToStr(k.moduleName, refundAddr.Bytes()),
		Tokens:      tokens,
		To:          types.ExternalAddrToStr(k.moduleName, to.Bytes()),
		Data:        hex.EncodeToString(data),
		Memo:        hex.EncodeToString(memo),
		EventNonce:  eventNonce,
	}
	return outCall, nil
}

func (k Keeper) BridgeCallResultHandler(ctx sdk.Context, claim *types.MsgBridgeCallResultClaim) error {
	k.CreateBridgeAccount(ctx, claim.TxOrigin)

	outgoingBridgeCall, found := k.GetOutgoingBridgeCallByNonce(ctx, claim.Nonce)
	if !found {
		return fmt.Errorf("bridge call not found for nonce %d", claim.Nonce)
	}
	if !claim.Success {
		if err := k.RefundOutgoingBridgeCall(ctx, outgoingBridgeCall); err != nil {
			return err
		}
	}
	k.DeleteOutgoingBridgeCallRecord(ctx, claim.Nonce)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCallResult,
		sdk.NewAttribute(types.AttributeKeyEventNonce, strconv.FormatInt(int64(claim.Nonce), 10)),
		sdk.NewAttribute(types.AttributeKeyStateSuccess, strconv.FormatBool(claim.Success)),
		sdk.NewAttribute(types.AttributeKeyErrCause, claim.Cause),
	))
	return nil
}

func (k Keeper) RefundOutgoingBridgeCall(ctx sdk.Context, data *types.OutgoingBridgeCall) error {
	refund := types.ExternalAddrToAccAddr(k.moduleName, data.GetRefund())
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

	for _, coin := range baseCoins {
		_, err := k.erc20Keeper.BaseCoinToEvm(ctx, common.BytesToAddress(refund.Bytes()), coin)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) DeleteOutgoingBridgeCallRecord(ctx sdk.Context, bridgeCallNonce uint64) {
	// 1. delete bridge call
	k.DeleteOutgoingBridgeCall(ctx, bridgeCallNonce)

	// 2. delete bridge call confirm
	k.DeleteBridgeCallConfirm(ctx, bridgeCallNonce)
}

func (k Keeper) SetOutgoingBridgeCall(ctx sdk.Context, outCall *types.OutgoingBridgeCall) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOutgoingBridgeCallNonceKey(outCall.Nonce), k.cdc.MustMarshal(outCall))
	// value is just a placeholder
	store.Set(
		types.GetOutgoingBridgeCallAddressAndNonceKey(outCall.Sender, outCall.Nonce),
		k.cdc.MustMarshal(&gogotypes.BoolValue{Value: true}),
	)
}

func (k Keeper) HasOutgoingBridgeCall(ctx sdk.Context, nonce uint64) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetOutgoingBridgeCallNonceKey(nonce))
}

func (k Keeper) HasOutgoingBridgeCallAddressAndNonce(ctx sdk.Context, sender string, nonce uint64) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetOutgoingBridgeCallAddressAndNonceKey(sender, nonce))
}

func (k Keeper) GetOutgoingBridgeCallByNonce(ctx sdk.Context, nonce uint64) (*types.OutgoingBridgeCall, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOutgoingBridgeCallNonceKey(nonce))
	if bz == nil {
		return nil, false
	}
	var outCall types.OutgoingBridgeCall
	k.cdc.MustUnmarshal(bz, &outCall)
	return &outCall, true
}

func (k Keeper) DeleteOutgoingBridgeCall(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	outCall, found := k.GetOutgoingBridgeCallByNonce(ctx, nonce)
	if !found {
		return
	}
	store.Delete(types.GetOutgoingBridgeCallNonceKey(nonce))
	store.Delete(types.GetOutgoingBridgeCallAddressAndNonceKey(outCall.Sender, outCall.Nonce))
}

func (k Keeper) IterateOutgoingBridgeCalls(ctx sdk.Context, cb func(outCall *types.OutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.OutgoingBridgeCallNonceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var outCall types.OutgoingBridgeCall
		k.cdc.MustUnmarshal(iterator.Value(), &outCall)
		if cb(&outCall) {
			break
		}
	}
}

func (k Keeper) IterateOutgoingBridgeCallsByAddress(ctx sdk.Context, senderAddr string, cb func(outCall *types.OutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.GetOutgoingBridgeCallAddressKey(senderAddr))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		nonce := types.ParseOutgoingBridgeCallNonce(iterator.Key(), senderAddr)
		outCall, found := k.GetOutgoingBridgeCallByNonce(ctx, nonce)
		if !found {
			continue
		}
		if cb(outCall) {
			break
		}
	}
}

func (k Keeper) IterateOutgoingBridgeCallByNonce(ctx sdk.Context, startNonce uint64, cb func(outCall *types.OutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	startKey := append(types.OutgoingBridgeCallNonceKey, sdk.Uint64ToBigEndian(startNonce)...)
	endKey := append(types.OutgoingBridgeCallNonceKey, sdk.Uint64ToBigEndian(math.MaxUint64)...)
	iter := store.Iterator(startKey, endKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		outCall := new(types.OutgoingBridgeCall)
		k.cdc.MustUnmarshal(iter.Value(), outCall)
		if cb(outCall) {
			break
		}
	}
}

func (k Keeper) BridgeCallBaseCoin(
	ctx sdk.Context,
	from, refund, to common.Address,
	coins sdk.Coins,
	data, memo []byte,
	fxTarget *types.FxTarget,
	originTokenAmount sdkmath.Int,
) (uint64, error) {
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
		sequence, err := k.IBCTransfer(ctx, from.Bytes(), toAddr, ibcCoin, fxTarget.IBCChannel, string(memo))
		if err != nil {
			return 0, err
		}
		if originTokenAmount.IsPositive() {
			ibcTransferKey := types.NewIBCTransferKey(fxTarget.IBCChannel, sequence)
			if err = k.erc20Keeper.SetCache(ctx, ibcTransferKey); err != nil {
				return 0, err
			}
		}
		return sequence, nil
	}
	// todo record origin amount
	return k.AddOutgoingBridgeCall(ctx, from, refund, coins, to, data, memo, 0)
}

func (k Keeper) CrosschainBaseCoin(
	ctx sdk.Context,
	from sdk.AccAddress,
	receipt string,
	amount, fee sdk.Coin,
	fxTarget *types.FxTarget,
	memo string,
	originToken bool,
) error {
	var cacheKey string
	if fxTarget.IsIBC() {
		sequence, err := k.IBCTransfer(ctx, from.Bytes(), receipt, amount, fxTarget.IBCChannel, memo)
		if err != nil {
			return err
		}
		cacheKey = types.NewIBCTransferKey(fxTarget.IBCChannel, sequence)
	} else {
		batchNonce, err := k.BuildOutgoingTxBatch(ctx, from, receipt, amount, fee)
		if err != nil {
			return err
		}
		cacheKey = types.NewOriginTokenKey(fxTarget.GetModuleName(), batchNonce)
	}

	if originToken {
		return k.erc20Keeper.SetCache(ctx, cacheKey)
	}
	return nil
}

func (k Keeper) IBCTransfer(
	ctx sdk.Context,
	from sdk.AccAddress,
	to string,
	amount sdk.Coin,
	channel string,
	memo string,
) (uint64, error) {
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
