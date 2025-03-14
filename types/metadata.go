package types

import (
	"strings"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func NewDefaultMetadata() banktypes.Metadata {
	// Ref: https://github.com/cosmos/chain-registry/blob/master/cosmoshub/assetlist.json#L14
	return banktypes.Metadata{
		Description: "The native staking token of the Pundi AIFX",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    DefaultDenom,
				Exponent: 0,
			},
			{
				Denom:    strings.ToLower(DefaultSymbol),
				Exponent: DenomUnit,
			},
		},
		Base:    DefaultDenom,
		Display: strings.ToLower(DefaultSymbol),
		Name:    "Pundi AIFX Token",
		Symbol:  DefaultSymbol,
	}
}

func NewMetadata(name, symbol string, decimals uint32) banktypes.Metadata {
	return banktypes.Metadata{
		Description: "The crosschain token of the Pundi AIFX",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    strings.ToLower(symbol),
				Exponent: 0,
			},
			{
				Denom:    symbol,
				Exponent: decimals,
			},
		},
		Base:    strings.ToLower(symbol),
		Display: symbol,
		Name:    name,
		Symbol:  symbol,
	}
}
