package types

import banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

func GetMetadata() []banktypes.Metadata {
	return []banktypes.Metadata{
		{
			Description: "Wrap Function X", // name
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    "FX",
					Exponent: 18,
					Aliases:  []string{},
				},
				{
					Denom:    "WFX",
					Exponent: 18,
					Aliases:  []string{},
				},
			},
			Base:    "FX",
			Display: "WFX", // symbl
		},
	}
}
