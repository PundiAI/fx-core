package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

const (
	PundixWrapDenom = "pundix"
	PundixChannel   = "channel-0"
)

const (
	MainnetPundixUnWrapDenom = "eth0x0FD10b9899882a6f2fcb5c371E17e70FdEe00C38"
)

const (
	TestnetPundixUnWrapDenom = "eth0xd9EEd31F5731DfC3Ca18f09B487e200F50a6343B"
)

var (
	MainnetOnRecvWrap = map[string]string{
		MainnetPundixUnWrapDenom: PundixWrapDenom,
		LegacyFXDenom:            DefaultDenom,
	}

	MainnetSendPacketWrap = map[string]string{
		SendPacketDenomWrapKey(PundixChannel, PundixWrapDenom): MainnetPundixUnWrapDenom,
	}

	MainnetAckPacketWrap = map[string]string{
		SendPacketDenomWrapKey(PundixChannel, MainnetPundixUnWrapDenom): PundixWrapDenom,
	}
)

var (
	TestnetOnRecvWrap = map[string]string{
		TestnetPundixUnWrapDenom: PundixWrapDenom,
		LegacyFXDenom:            DefaultDenom,
	}
	TestnetSendPacketWrap = map[string]string{
		SendPacketDenomWrapKey(PundixChannel, PundixWrapDenom): TestnetPundixUnWrapDenom,
	}
	TestnetAckPacketWrap = map[string]string{
		SendPacketDenomWrapKey(PundixChannel, TestnetPundixUnWrapDenom): PundixWrapDenom,
	}
)

func OnRecvDenomNeedWrap(chainId, sourcePort, sourceChannel, denom string) (needWrap bool, wrapDenom, packetDenom string) {
	if !transfertypes.ReceiverChainIsSource(sourcePort, sourceChannel, denom) {
		return false, "", ""
	}

	voucherPrefix := transfertypes.GetDenomPrefix(sourcePort, sourceChannel)
	unprefixedDenom := denom[len(voucherPrefix):]

	if chainId == MainnetChainId {
		wrapDenom, needWrap = MainnetOnRecvWrap[unprefixedDenom]
	} else {
		wrapDenom, needWrap = TestnetOnRecvWrap[unprefixedDenom]
	}
	if needWrap {
		packetDenom = fmt.Sprintf("%s%s", voucherPrefix, wrapDenom)
	}
	return needWrap, wrapDenom, packetDenom
}

func SendPacketDenomNeedWrap(chainId, sourceChannel, denom string) (bool, string) {
	var needWrap bool
	var wrapDenom string
	denomWrapKey := SendPacketDenomWrapKey(sourceChannel, denom)
	if chainId == MainnetChainId {
		wrapDenom, needWrap = MainnetSendPacketWrap[denomWrapKey]
	} else {
		wrapDenom, needWrap = TestnetSendPacketWrap[denomWrapKey]
	}
	return needWrap, wrapDenom
}

func AckPacketDenomNeedWrap(chainId, sourceChannel, denom string) (bool, string) {
	var needWrap bool
	var wrapDenom string
	denomWrapKey := SendPacketDenomWrapKey(sourceChannel, denom)
	if chainId == MainnetChainId {
		wrapDenom, needWrap = MainnetAckPacketWrap[denomWrapKey]
	} else {
		wrapDenom, needWrap = TestnetAckPacketWrap[denomWrapKey]
	}
	return needWrap, wrapDenom
}

func SendPacketDenomWrapKey(sourceChannel, denonm string) string {
	return fmt.Sprintf("%s:%s", sourceChannel, denonm)
}

func OnRecvAmountCovert(wrapDenom, amountStr string) (string, error) {
	if wrapDenom != DefaultDenom {
		return amountStr, nil
	}

	amount, ok := sdkmath.NewIntFromString(amountStr)
	if !ok {
		return amountStr, transfertypes.ErrInvalidAmount.Wrapf("unable to parse transfer amount (%s) into sdkmath.Int", amountStr)
	}
	return SwapAmount(amount).String(), nil
}
