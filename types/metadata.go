package types

import (
	"fmt"
	"strings"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func GetFXMetaData() banktypes.Metadata {
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

func GetCrossChainMetadataManyToOne(name, symbol string, decimals uint32, aliases ...string) banktypes.Metadata {
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

func GetCrossChainMetadataOneToOne(name, denom, symbol string, decimals uint32) banktypes.Metadata {
	return banktypes.Metadata{
		Description: "The cross chain token of the Function X",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    denom,
				Exponent: 0,
			},
			{
				Denom:    symbol,
				Exponent: decimals,
			},
		},
		Base:    denom,
		Display: denom,
		Name:    name,
		Symbol:  symbol,
	}
}

func ValidateMetadata(md banktypes.Metadata) error {
	decimals := uint8(0)
	for _, du := range md.DenomUnits {
		if du.Denom == md.Symbol {
			decimals = uint8(du.Exponent)
			break
		}
	}
	if md.Base == DefaultDenom {
		decimals = DenomUnit
	}
	if len(md.Name) == 0 {
		return fmt.Errorf("invalid name %s", md.Name)
	}
	if len(md.Symbol) == 0 {
		return fmt.Errorf("invalid symbol %s", md.Symbol)
	}
	if decimals == 0 {
		return fmt.Errorf("invalid decimals %d", decimals)
	}
	return nil
}
