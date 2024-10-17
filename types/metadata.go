package types

import (
	"strings"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func NewFXMetaData() banktypes.Metadata {
	return banktypes.Metadata{
		Description: "The native staking token of the Function X",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    DefaultDenom,
				Exponent: 0,
			},
		},
		Base:    DefaultDenom,
		Display: DefaultDenom,
		Name:    "Function X",
		Symbol:  DefaultDenom,
	}
}

func NewMetadata(name, symbol string, decimals uint32) banktypes.Metadata {
	return banktypes.Metadata{
		Description: "The crosschain token of the Function X",
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
		Display: strings.ToLower(symbol),
		Name:    name,
		Symbol:  symbol,
	}
}
