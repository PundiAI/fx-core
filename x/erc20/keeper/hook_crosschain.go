package keeper

import (
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (h Hooks) HookTransferCrossChain(ctx sdk.Context, relayTransferCrossChains []types.RelayTransferCrossChain, from common.Address, to *common.Address, txHash common.Hash) (err error) {
	logger := h.k.Logger(ctx)
	for _, relay := range relayTransferCrossChains {
		logger.Info("transfer cross", "tx-hash", txHash.Hex(), "from", from.Hex(), "to", to.Hex(), "token", relay.TokenContract.String(), "denom", relay.Denom)

		balances := h.k.bankKeeper.GetAllBalances(ctx, relay.From.Bytes())
		if !balances.IsAllGTE(relay.TotalAmount(relay.Denom)) {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "%s is smaller than %s", balances.String(), relay.TotalAmount(relay.Denom).String())
		}

		targetType, target := relay.GetTarget()
		amount := relay.GetAmount(relay.Denom)
		fee := relay.GetFee(relay.Denom)

		switch targetType {
		case types.FIP20TargetChain:
			err = h.transferChainHandler(ctx, relay.GetFrom(), relay.Recipient, amount, fee, target)
		case types.FIP20TargetIBC:
			err = h.transferIBCHandler(ctx, relay.GetFrom(), relay.Recipient, amount, fee, target, txHash)
		default:
			err = sdkerrors.Wrapf(types.ErrInvalidTarget, "target type %s", targetType)
		}
		if err != nil {
			logger.Error("failed to transfer cross chain", "tx-hash", txHash.Hex(), "error", err.Error())
			return err
		}

		ctx.EventManager().EmitEvents(
			sdk.Events{
				sdk.NewEvent(
					types.EventTypeRelayTransferCrossChain,
					sdk.NewAttribute(sdk.AttributeKeySender, from.String()),
					sdk.NewAttribute(types.AttributeKeyTo, to.String()),
					sdk.NewAttribute(types.AttributeKeyEvmTxHash, txHash.String()),
					sdk.NewAttribute(types.AttributeKeyFrom, relay.From.String()),
					sdk.NewAttribute(types.AttributeKeyRecipient, relay.Recipient),
					sdk.NewAttribute(sdk.AttributeKeyAmount, relay.Amount.String()),
					sdk.NewAttribute(sdk.AttributeKeyFee, relay.Fee.String()),
					sdk.NewAttribute(types.AttributeKeyTarget, fxtypes.Byte32ToString(relay.Target)),
					sdk.NewAttribute(types.AttributeKeyTokenAddress, relay.TokenContract.String()),
					sdk.NewAttribute(types.AttributeKeyDenom, relay.Denom),
				),
			},
		)

		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "relay_transfer_cross_chain"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("erc20", relay.TokenContract.String()),
				telemetry.NewLabel("denom", relay.Denom),
				telemetry.NewLabel("type", targetType.String()),
				telemetry.NewLabel("target", target),
				telemetry.NewLabel("amount", relay.GetAmount(relay.Denom).String()),
			},
		)
	}
	return nil
}

func (h Hooks) transferChainHandler(ctx sdk.Context, from sdk.AccAddress, to string, amount, fee sdk.Coin, target string) error {
	h.k.Logger(ctx).Info("transfer chain handler", "from", from, "to", to, "amount", amount.String(), "fee", fee.String(), "target", target)
	if h.k.router == nil {
		return sdkerrors.Wrapf(types.ErrInvalidTarget, "not set router")
	}
	route, has := h.k.router.GetRoute(target)
	if !has {
		return sdkerrors.Wrapf(types.ErrInvalidTarget, "target %s not support", target)
	}
	return route.TransferAfter(ctx, from.String(), to, amount, fee)
}

func (h Hooks) transferIBCHandler(ctx sdk.Context, from sdk.AccAddress, to string, amount, fee sdk.Coin, target string, txHash common.Hash) error {
	logger := h.k.Logger(ctx)
	logger.Info("transfer ibc handler", "from", from, "to", to, "amount", amount.String(), "fee", fee.String(), "target", target)

	targetIBC, ok := fxtypes.ParseTargetIBC(target)
	if !ok {
		return sdkerrors.Wrapf(types.ErrInvalidTarget, "invalid target ibc %s", target)
	}
	if err := validateIbcReceiveAddress(targetIBC.Prefix, to); err != nil {
		logger.Error("validate ibc receive address", "address", to, "prefix", targetIBC.Prefix, "err", err.Error())
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid to address %s, error %s", to, err.Error())
	}
	if !fee.IsZero() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "ibc transfer fee must be zero: %s", fee.Amount.String())
	}
	ibcTimeout := h.k.GetIbcTimeout(ctx)
	ibcTimeoutHeight := ibcclienttypes.ZeroHeight()
	ibcTimeoutTimestamp := uint64(ctx.BlockTime().UnixNano()) + uint64(ibcTimeout)
	transferMsg := transfertypes.NewMsgTransfer(targetIBC.SourcePort, targetIBC.SourceChannel, amount, from.String(), to, ibcTimeoutHeight, ibcTimeoutTimestamp)
	transferResponse, err := h.k.IbcTransferKeeper.Transfer(sdk.WrapSDKContext(ctx), transferMsg)
	if err != nil {
		return err
	}
	h.k.SetIBCTransferHash(ctx, targetIBC.SourcePort, targetIBC.SourceChannel, transferResponse.GetSequence(), txHash)
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
