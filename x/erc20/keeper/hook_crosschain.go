package keeper

import (
	"fmt"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (h Hooks) HookTransferCrossChainEvent(ctx sdk.Context, relayTransferCrossChains []types.RelayTransferCrossChain, txHash common.Hash) (err error) {
	for _, relay := range relayTransferCrossChains {
		h.k.Logger(ctx).Info("transfer cross chain", "token", relay.TokenContract.String(), "denom", relay.Denom)

		balances := h.k.bankKeeper.GetAllBalances(ctx, relay.From.Bytes())
		if !balances.IsAllGTE(relay.TotalAmount(relay.Denom)) {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "%s is smaller than %s", balances.String(), relay.TotalAmount(relay.Denom).String())
		}

		amount := relay.GetAmount(relay.Denom)
		fee := relay.GetFee(relay.Denom)

		target := fxtypes.Byte32ToString(relay.Target)
		if strings.HasPrefix(target, ibchost.ModuleName) {
			target = strings.TrimPrefix(target, fmt.Sprintf("%s/", ibchost.ModuleName))
			err = h.transferIBCHandler(ctx, relay.GetFrom(), relay.Recipient, amount, fee, target, txHash)
		} else {
			target = strings.TrimPrefix(target, "chain/")
			err = h.transferCrossChainHandler(ctx, relay.GetFrom(), relay.Recipient, amount, fee, target)
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

func (h Hooks) transferCrossChainHandler(ctx sdk.Context, from sdk.AccAddress, to string, amount, fee sdk.Coin, target string) error {
	h.k.Logger(ctx).Info("transfer cross-chain handler", "from", from, "to", to, "amount", amount.String(), "fee", fee.String(), "target", target)
	if h.k.router == nil {
		return sdkerrors.Wrapf(types.ErrInternalRouter, "transfer chain router not set")
	}
	route, has := h.k.router.GetRoute(target)
	if !has {
		return sdkerrors.Wrapf(types.ErrInvalidTarget, "target %s not support", target)
	}
	return route.TransferAfter(ctx, from.String(), to, amount, fee)
}

func (h Hooks) transferIBCHandler(ctx sdk.Context, from sdk.AccAddress, to string, amount, fee sdk.Coin, targetStr string, txHash common.Hash) error {
	h.k.Logger(ctx).Info("transfer ibc handler", "from", from, "to", to, "amount", amount.String(), "fee", fee.String(), "target", targetStr)
	if !fee.IsZero() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "ibc transfer fee must be zero: %s", fee.Amount.String())
	}

	target := fxtypes.ParseFxTarget(targetStr)
	if !target.IsIBC() {
		return sdkerrors.Wrapf(types.ErrInvalidTarget, "invalid target ibc %s", targetStr)
	}
	if strings.ToLower(target.Prefix) == fxtypes.EthereumAddressPrefix {
		if err := fxtypes.ValidateEthereumAddress(to); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid to address %s, error %s", to, err.Error())
		}
	}
	if _, err := sdk.GetFromBech32(to, target.Prefix); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid to address %s, error %s", to, err.Error())
	}

	ibcTimeoutTimestamp := uint64(ctx.BlockTime().UnixNano()) + uint64(h.k.GetIbcTimeout(ctx))
	transferResponse, err := h.k.ibcTransferKeeper.Transfer(sdk.WrapSDKContext(ctx),
		transfertypes.NewMsgTransfer(
			target.SourcePort,
			target.SourceChannel,
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
	h.k.SetIBCTransferHash(ctx, target.SourcePort, target.SourceChannel, transferResponse.GetSequence(), txHash)
	return nil
}
