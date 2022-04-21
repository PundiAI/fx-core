package types

import banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

func GetFxBankMetaData(denom string) banktypes.Metadata {
	return banktypes.Metadata{
		Description: "Function X",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    denom,
				Exponent: 18,
				Aliases:  nil,
			},
		},
		Base:    denom,
		Display: denom,
	}
}

func GetMetadata() []banktypes.Metadata {
	return []banktypes.Metadata{
		{
			Description: "Wrapped Function X", // name
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    "WFX",
					Exponent: 18,             //decimals
					Aliases:  []string{"FX"}, //
				},
			},
			Base:    "WFX",
			Display: "WFX", // symbol
		},
		//{
		//	Description: "Pundi X Token", // name
		//	DenomUnits: []*banktypes.DenomUnit{
		//		{
		//			Denom:    "PUNDIX",
		//			Exponent: 18,                //decimals
		//			Aliases:  []string{"eth0x"}, //
		//		},
		//	},
		//	Base:    "PUNDIX",
		//	Display: "PUNDIX", // symbol
		//},
	}
}
