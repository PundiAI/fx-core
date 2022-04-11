package types

var ChainID = "fxcore"

const (
	Name          = "fxcore"
	AddressPrefix = "fx"

	MintDenom = "FX"
	// BaseDenomUnit defines the base denomination unit for Photons.
	// 1 FX = 1x10^{BaseDenomUnit} fx
	BaseDenomUnit = 18

	// DefaultGasPrice is default gas price for evm transactions
	DefaultGasPrice = 20
)
