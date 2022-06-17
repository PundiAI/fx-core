package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var ChainID = "fxcore"

const (
	Name          = "fxcore"
	AddressPrefix = "fx"

	DefaultDenom = "FX"
	// BaseDenomUnit defines the base denomination unit for Photons.
	// 1 FX = 1x10^{BaseDenomUnit} fx
	BaseDenomUnit = 18

	// DefaultGasPrice is default gas price for evm transactions 500Gwei
	DefaultGasPrice = 500000000000
)

func init() {
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/-]{1,127}`
	})
}
