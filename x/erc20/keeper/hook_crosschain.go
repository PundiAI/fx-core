package keeper

import (
	"strings"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (h Hooks) HookTransferCrossChainEvent(ctx sdk.Context, relayTransferCrossChains []types.RelayTransferCrossChain) (err error) {
	for _, relay := range relayTransferCrossChains {
		h.k.Logger(ctx).Info("transfer cross chain", "token", relay.TokenContract.String(), "denom", relay.Denom)

		balances := h.k.bankKeeper.GetAllBalances(ctx, relay.From.Bytes())
		if !balances.IsAllGTE(relay.TotalAmount(relay.Denom)) {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "%s is smaller than %s", balances.String(), relay.TotalAmount(relay.Denom).String())
		}

		targetStr := fxtypes.Byte32ToString(relay.Target)
		fxTarget := fxtypes.ParseFxTarget(targetStr)

		targetCoin, err := h.k.ConvertDenomToTarget(ctx, relay.GetFrom(), relay.GetAmount(relay.Denom).Add(relay.GetFee(relay.Denom)), fxTarget)
		if err != nil {
			return err
		}

		amount := relay.GetAmount(targetCoin.Denom)
		fee := relay.GetFee(targetCoin.Denom)

		if fxTarget.IsIBC() {
			err = h.transferIBCHandler(ctx, relay.GetFrom(), relay.Recipient, amount, fee, fxTarget)
		} else {
			err = h.transferCrossChainHandler(ctx, relay.GetFrom(), relay.Recipient, amount, fee, fxTarget)
		}
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeRelayTransferCrossChain,
			sdk.NewAttribute(types.AttributeKeyFrom, relay.From.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, relay.Recipient),
			sdk.NewAttribute(sdk.AttributeKeyAmount, relay.Amount.String()),
			sdk.NewAttribute(sdk.AttributeKeyFee, relay.Fee.String()),
			sdk.NewAttribute(types.AttributeKeyTarget, fxtypes.Byte32ToString(relay.Target)),
			sdk.NewAttribute(types.AttributeKeyTokenAddress, relay.TokenContract.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, relay.Denom),
		))

		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "relay_transfer_cross_chain"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("erc20", relay.TokenContract.String()),
				telemetry.NewLabel("denom", relay.Denom),
				telemetry.NewLabel("target", fxtypes.Byte32ToString(relay.Target)),
				telemetry.NewLabel("amount", relay.GetAmount(relay.Denom).String()),
			},
		)
	}
	return nil
}

func (h Hooks) transferCrossChainHandler(ctx sdk.Context, from sdk.AccAddress, to string, amount, fee sdk.Coin, fxTarget fxtypes.FxTarget) error {
	h.k.Logger(ctx).Info("transfer cross-chain handler", "from", from, "to", to, "amount", amount.String(), "fee", fee.String(), "target", fxTarget.GetTarget())
	if h.k.router == nil {
		return sdkerrors.Wrapf(types.ErrInternalRouter, "transfer chain router not set")
	}
	route, has := h.k.router.GetRoute(fxTarget.GetTarget())
	if !has {
		return sdkerrors.Wrapf(types.ErrInvalidTarget, "target %s not support", fxTarget.GetTarget())
	}
	return route.TransferAfter(ctx, from.String(), to, amount, fee)
}

func (h Hooks) transferIBCHandler(ctx sdk.Context, from sdk.AccAddress, to string, amount, fee sdk.Coin, fxTarget fxtypes.FxTarget) error {
	h.k.Logger(ctx).Info("transfer ibc handler", "from", from, "to", to, "amount", amount.String(), "fee", fee.String(), "target", fxTarget.GetTarget())
	if !fee.IsZero() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "ibc transfer fee must be zero: %s", fee.Amount.String())
	}
	if strings.ToLower(fxTarget.Prefix) == fxtypes.EthereumAddressPrefix {
		if err := fxtypes.ValidateEthereumAddress(to); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid to address %s, error %s", to, err.Error())
		}
	}
	if _, err := sdk.GetFromBech32(to, fxTarget.Prefix); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid to address %s, error %s", to, err.Error())
	}

	ibcTimeoutTimestamp := uint64(ctx.BlockTime().UnixNano()) + uint64(h.k.GetIbcTimeout(ctx))
	transferResponse, err := h.k.ibcTransferKeeper.Transfer(sdk.WrapSDKContext(ctx),
		transfertypes.NewMsgTransfer(
			fxTarget.SourcePort,
			fxTarget.SourceChannel,
			amount,
			from.String(),
			to,
			ibcclienttypes.ZeroHeight(),
			ibcTimeoutTimestamp,
		),
	)
	if err != nil {
		return err
	}
	h.k.SetIBCTransferRelation(ctx, fxTarget.SourceChannel, transferResponse.GetSequence())
	return nil
}
