package legacy

import (
	"math/big"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// Fip20CrossChainEvents use for fip20 cross chain
// Deprecated
func Fip20CrossChainEvents(ctx sdk.Context, from, token common.Address, recipient, target, denom string, amount, fee *big.Int) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		EventTypeRelayTransferCrossChain,
		sdk.NewAttribute(AttributeKeyFrom, from.String()),
		sdk.NewAttribute(AttributeKeyRecipient, recipient),
		sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
		sdk.NewAttribute(sdk.AttributeKeyFee, fee.String()),
		sdk.NewAttribute(AttributeKeyTarget, target),
		sdk.NewAttribute(AttributeKeyTokenAddress, token.String()),
		sdk.NewAttribute(AttributeKeyDenom, denom),
	))

	telemetry.IncrCounterWithLabels(
		[]string{"relay_transfer_cross_chain"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("erc20", token.String()),
			telemetry.NewLabel("denom", denom),
			telemetry.NewLabel("target", target),
		},
	)
}

// CrossChainEvents
// Deprecated
func CrossChainEvents(ctx sdk.Context, from, token common.Address, recipient, target, denom, memo string, amount, fee *big.Int) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		EventTypeCrossChain,
		sdk.NewAttribute(AttributeKeyFrom, from.String()),
		sdk.NewAttribute(AttributeKeyRecipient, recipient),
		sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
		sdk.NewAttribute(sdk.AttributeKeyFee, fee.String()),
		sdk.NewAttribute(AttributeKeyTarget, target),
		sdk.NewAttribute(AttributeKeyTokenAddress, token.String()),
		sdk.NewAttribute(AttributeKeyDenom, denom),
		sdk.NewAttribute(AttributeKeyMemo, memo),
	))
}

const (
	// EventTypeRelayTransferCrossChain
	// Deprecated
	EventTypeRelayTransferCrossChain = "relay_transfer_cross_chain"
	// EventTypeCrossChain new cross chain event type
	EventTypeCrossChain = "cross_chain"

	AttributeKeyDenom        = "coin"
	AttributeKeyTokenAddress = "token_address"
	AttributeKeyFrom         = "from"
	AttributeKeyRecipient    = "recipient"
	AttributeKeyTarget       = "target"
	AttributeKeyMemo         = "memo"
)
