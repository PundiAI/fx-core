package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/ibc/middleware/types"
)

func parseIBCCoinDenom(packet channeltypes.Packet, packetDenom string) string {
	// This is the prefix that would have been prefixed to the denomination
	// on sender chain IF and only if the token originally came from the
	// receiving chain.
	//
	// NOTE: We use SourcePort and SourceChannel here, because the counterparty
	// chain would have prefixed with DestPort and DestChannel when originally
	// receiving this coin as seen in the "sender chain is the source" condition.

	var receiveDenom string
	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), packetDenom) {
		// sender chain is not the source, unescrow tokens

		// remove prefix added by sender chain
		voucherPrefix := transfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())
		unprefixedDenom := packetDenom[len(voucherPrefix):]

		// coin denomination used in sending from the escrow address
		denom := unprefixedDenom

		// The denomination used to send the coins is either the native denom or the hash of the path
		// if the denomination is not native.
		denomTrace := transfertypes.ParseDenomTrace(unprefixedDenom)
		if !denomTrace.IsNativeDenom() {
			denom = denomTrace.IBCDenom()
		}
		receiveDenom = denom
	} else {
		// sender chain is the source, mint vouchers

		// since SendPacket did not prefix the denomination, we must prefix denomination here
		sourcePrefix := transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel())
		// NOTE: sourcePrefix contains the trailing "/"
		prefixedDenom := sourcePrefix + packetDenom

		// construct the denomination trace from the full raw denomination
		denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)

		receiveDenom = denomTrace.IBCDenom()
	}
	return receiveDenom
}

func parseReceiveAndAmountByPacketWithRouter(data types.FungibleTokenPacketData) (sdk.AccAddress, sdkmath.Int, sdkmath.Int, error) {
	// parse the transfer amount
	transferAmount, ok := sdkmath.NewIntFromString(data.Amount)
	if !ok {
		return nil, sdkmath.Int{}, sdkmath.Int{}, transfertypes.ErrInvalidAmount.Wrapf("unable to parse transfer amount (%s) into sdkmath.Int", data.Amount)
	}

	addressBytes, _, err := fxtypes.ParseAddress(data.Sender)
	if err != nil {
		return nil, sdkmath.Int{}, sdkmath.Int{}, err
	}
	feeAmount, ok := sdkmath.NewIntFromString(data.Fee)
	if !ok || feeAmount.IsNegative() {
		return nil, sdkmath.Int{}, sdkmath.Int{}, transfertypes.ErrInvalidAmount.Wrapf("fee amount is invalid:%s", data.Fee)
	}
	return addressBytes, transferAmount, feeAmount, nil
}

func parseAmountAndFeeByPacket(data types.FungibleTokenPacketData) (sdkmath.Int, sdkmath.Int, error) {
	// parse the transfer amount
	transferAmount, ok := sdkmath.NewIntFromString(data.Amount)
	if !ok {
		return sdkmath.Int{}, sdkmath.Int{}, transfertypes.ErrInvalidAmount.Wrapf("unable to parse transfer amount (%s) into sdkmath.Int", data.Amount)
	}

	feeAmount := sdkmath.ZeroInt()
	if data.Router != "" {
		fee, ok := sdkmath.NewIntFromString(data.Fee)
		if !ok || fee.IsNegative() {
			return sdkmath.Int{}, sdkmath.Int{}, transfertypes.ErrInvalidAmount.Wrapf("fee amount is invalid:%s", data.Fee)
		}
		feeAmount = fee
	}
	return transferAmount, feeAmount, nil
}
