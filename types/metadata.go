package types

import (
	"strings"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func GetFXMetaData(denom string) banktypes.Metadata {
	return banktypes.Metadata{
		Description: "The native staking token of the Function X",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    denom,
				Exponent: 0,
			},
		},
		Base:    denom,
		Display: denom,
		Name:    "Function X",
		Symbol:  denom,
	}
}

func GetCrossChainMetadata(name, symbol string, decimals uint32, aliases ...string) banktypes.Metadata {
	return banktypes.Metadata{
		Description: "The cross chain token of the Function X",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    strings.ToLower(symbol),
				Exponent: 0,
				Aliases:  aliases,
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
