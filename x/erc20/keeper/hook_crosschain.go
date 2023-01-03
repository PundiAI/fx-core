package keeper

import (
	"fmt"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (k Keeper) HookTransferCrossChain(ctx sdk.Context, tcels []*TransferCrossChainEventLog, from common.Address, to *common.Address, receipt *ethtypes.Receipt) error {
	logger := k.Logger(ctx)
	for _, tcel := range tcels {
		logger.Info("transfer cross", "tx-hash", receipt.TxHash.Hex(),
			"from", from.Hex(), "to", to.Hex(), "token", tcel.Pair.Erc20Address, "denom", tcel.Pair.Denom)

		balances := k.bankKeeper.GetAllBalances(ctx, tcel.Event.From.Bytes())
		if !balances.IsAllGTE(tcel.Event.TotalAmount(tcel.Pair.Denom)) {
			return fmt.Errorf("insufficient balance, have %s expected %s", balances.String(), tcel.Event.TotalAmount(tcel.Pair.Denom).String())
		}

		var err error
		targetType, target := tcel.Event.GetTarget()
		amount := tcel.Event.GetAmount(tcel.Pair.Denom)
		fee := tcel.Event.GetFee(tcel.Pair.Denom)

		switch targetType {
		case fxtypes.FIP20TargetChain:
			err = k.TransferChainHandler(ctx, tcel.Event.GetFrom(), tcel.Event.Recipient, amount, fee, target, receipt)
		case fxtypes.FIP20TargetIBC:
			err = k.TransferIBCHandler(ctx, tcel.Event.GetFrom(), tcel.Event.Recipient, amount, fee, target, receipt)
		default:
			err = fmt.Errorf("traget unknown %d", targetType)
		}
		if err != nil {
			logger.Error("failed to transfer cross chain", "tx-hash", receipt.TxHash.Hex(), "error", err.Error())
			return err
		}
		logger.Info("transfer cross chain success", "tx-hash", receipt.TxHash.Hex())

		ctx.EventManager().EmitEvents(
			sdk.Events{
				sdk.NewEvent(
					types.EventTypeRelayTransferCrossChain,
					sdk.NewAttribute(sdk.AttributeKeySender, from.String()),
					sdk.NewAttribute(types.AttributeKeyTo, to.String()),
					sdk.NewAttribute(types.AttributeKeyEvmTxHash, receipt.TxHash.String()),
					sdk.NewAttribute(types.AttributeKeyFrom, tcel.Event.From.String()),
					sdk.NewAttribute(types.AttributeKeyRecipient, tcel.Event.Recipient),
					sdk.NewAttribute(sdk.AttributeKeyAmount, tcel.Event.Amount.String()),
					sdk.NewAttribute(sdk.AttributeKeyFee, tcel.Event.Fee.String()),
					sdk.NewAttribute(types.AttributeKeyTarget, fxtypes.Byte32ToString(tcel.Event.Target)),
					sdk.NewAttribute(types.AttributeKeyTokenAddress, tcel.Pair.Erc20Address),
					sdk.NewAttribute(types.AttributeKeyDenom, tcel.Pair.Denom),
				),
			},
		)

		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "relay_transfer_cross_chain"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("erc20", tcel.Pair.Erc20Address),
				telemetry.NewLabel("denom", tcel.Pair.Denom),
				telemetry.NewLabel("type", targetType.String()),
				telemetry.NewLabel("target", target),
				telemetry.NewLabel("amount", tcel.Event.GetAmount(tcel.Pair.Denom).String()),
			},
		)
	}
	return nil
}

func (k Keeper) TransferChainHandler(ctx sdk.Context, from sdk.AccAddress, to string, amount, fee sdk.Coin, target string, _ *ethtypes.Receipt) error {
	k.Logger(ctx).Info("transfer chain handler", "from", from, "to", to, "amount", amount.String(), "fee", fee.String(), "target", target)
	if k.router == nil || !k.router.HasRoute(target) {
		return fmt.Errorf("target %s not support", target)
	}

	route, _ := k.router.GetRoute(target)
	return route.TransferAfter(ctx, from.String(), to, amount, fee)
}

func (k Keeper) TransferIBCHandler(ctx sdk.Context, from sdk.AccAddress, to string, amount, fee sdk.Coin, target string, receipt *ethtypes.Receipt) error {
	logger := k.Logger(ctx)
	logger.Info("transfer ibc handler", "from", from, "to", to, "amount", amount.String(), "fee", fee.String(), "target", target)

	targetIBC, ok := fxtypes.ParseTargetIBC(target)
	if !ok {
		return fmt.Errorf("invalid target ibc %s", target)
	}
	if err := validateIbcReceiveAddress(targetIBC.Prefix, to); err != nil {
		logger.Error("validate ibc receive address", "address", to, "prefix", targetIBC.Prefix, "err", err.Error())
		return fmt.Errorf("invalid to address %s", to)
	}
	if !fee.IsZero() {
		return fmt.Errorf("ibc transfer fee must be zero: %s", fee.Amount.String())
	}
	ibcTimeout := k.GetIbcTimeout(ctx)
	ibcTimeoutHeight := ibcclienttypes.ZeroHeight()
	ibcTimeoutTimestamp := uint64(ctx.BlockTime().UnixNano()) + uint64(ibcTimeout)
	transferMsg := transfertypes.NewMsgTransfer(targetIBC.SourcePort, targetIBC.SourceChannel, amount, from.String(), to, ibcTimeoutHeight, ibcTimeoutTimestamp)
	transferResponse, err := k.IbcTransferKeeper.Transfer(sdk.WrapSDKContext(ctx), transferMsg)
	if err != nil {
		return err
	}
	k.SetIBCTransferHash(ctx, targetIBC.SourcePort, targetIBC.SourceChannel, transferResponse.GetSequence(), receipt.TxHash)
	return nil
}

func validateIbcReceiveAddress(prefix, addr string) error {
	// after block support denom many-to-one, validate prefix with 0x
	if strings.ToLower(prefix) == fxtypes.EthereumAddressPrefix {
		return fxtypes.ValidateEthereumAddress(addr)
	}
	_, err := sdk.GetFromBech32(addr, prefix)
	return err
}

func (k Keeper) SetIBCTransferHash(ctx sdk.Context, port, channel string, sequence uint64, hash common.Hash) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetIBCTransferKey(port, channel, sequence), hash.Bytes())
}

func (k Keeper) DeleteIBCTransferHash(ctx sdk.Context, port, channel string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetIBCTransferKey(port, channel, sequence))
}

func (k Keeper) GetIBCTransferHash(ctx sdk.Context, port, channel string, sequence uint64) (common.Hash, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetIBCTransferKey(port, channel, sequence)
	if !store.Has(key) {
		return common.Hash{}, false
	}
	value := store.Get(key)
	return common.BytesToHash(value), true
}

func (k Keeper) HasIBCTransferHash(ctx sdk.Context, port, channel string, sequence uint64) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetIBCTransferKey(port, channel, sequence))
}

func (k Keeper) IBCTransferHashIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.KeyPrefixIBCTransfer)
}
