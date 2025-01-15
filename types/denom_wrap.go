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
	MainnetOsmosisChannel    = "channel-19"
)

const (
	TestnetPundixUnWrapDenom = "eth0xd9EEd31F5731DfC3Ca18f09B487e200F50a6343B"
	TestnetOsmosisChannel    = "channel-119"
)

var (
	MainnetOnRecvWrap = map[string]string{
		OnRecvDenomWrapKey(PundixChannel, MainnetPundixUnWrapDenom): PundixWrapDenom,
		OnRecvDenomWrapKey(MainnetOsmosisChannel, FXDenom):          DefaultDenom,
	}

	MainnetSendPacketWrap = map[string]string{
		SendPacketDenomWrapKey(PundixChannel, PundixWrapDenom): MainnetPundixUnWrapDenom,
	}
)

var (
	TestnetOnRecvWrap = map[string]string{
		OnRecvDenomWrapKey(PundixChannel, TestnetPundixUnWrapDenom): PundixWrapDenom,
		OnRecvDenomWrapKey(TestnetOsmosisChannel, FXDenom):          DefaultDenom,
	}
	TestnetSendPacketWrap = map[string]string{
		SendPacketDenomWrapKey(PundixChannel, PundixWrapDenom): TestnetPundixUnWrapDenom,
	}
)

func OnRecvDenomNeedWrap(chainId, channel, denom string) (bool, string) {
	var needWrap bool
	var wrapDenom string
	denomWrapKey := OnRecvDenomWrapKey(channel, denom)
	if chainId == MainnetChainId {
		wrapDenom, needWrap = MainnetOnRecvWrap[denomWrapKey]
	} else {
		wrapDenom, needWrap = TestnetOnRecvWrap[denomWrapKey]
	}
	return needWrap, wrapDenom
}

func SendPacketDenomNeedWrap(chainId, channel, denom string) (bool, string) {
	var needWrap bool
	var wrapDenom string
	denomWrapKey := SendPacketDenomWrapKey(channel, denom)
	if chainId == MainnetChainId {
		wrapDenom, needWrap = MainnetSendPacketWrap[denomWrapKey]
	} else {
		wrapDenom, needWrap = TestnetSendPacketWrap[denomWrapKey]
	}
	return needWrap, wrapDenom
}

func SendPacketDenomWrapKey(sourceChannel, denonm string) string {
	return fmt.Sprintf("%s:%s", sourceChannel, denonm)
}

func OnRecvDenomWrapKey(destChannel, denonm string) string {
	return fmt.Sprintf("%s:%s", destChannel, denonm)
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
