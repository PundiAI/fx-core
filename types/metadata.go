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

func GetCrossChainMetadata(name, symbol string, decimals uint32, denom string) banktypes.Metadata {
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

func GetCrossChainMetadataManyToOne(name, symbol string, decimals uint32, denom string) banktypes.Metadata {
	return banktypes.Metadata{
		Description: "The cross chain token of the Function X",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    strings.ToLower(symbol),
				Exponent: 0,
				Aliases: []string{
					denom,
				},
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

/*
		//// example of origin denom
		{
			"description": "The native staking token of the Function X",
			"denom_units": [
				{
					"denom": "FX",
					"exponent": 0,
					"aliases": []
				}
			],
			"base": "FX",
			"display": "FX",
			"name": "Function X",
			"symbol": "FX"
		}
		//// example of other denom
		{
			"description":"The cross chain token of the Function X",
			"denom_units":[
				{
					"denom":"eth0x338E7A8687AdA7274Dc87C95D94f920d8F4185AE",
					"exponent":0,
					"aliases":[

					]
				},
				{
					"denom":"PUNDIX",
					"exponent":18,
					"aliases":[

					]
				}
			],
			"base":"eth0x338E7A8687AdA7274Dc87C95D94f920d8F4185AE",
			"display":"eth0x338E7A8687AdA7274Dc87C95D94f920d8F4185AE",
			"name":"Pundi X Token",
			"symbol":"PUNDIX"
		}
		//// example many to one denom
	    {
	        "description":"The cross chain token of the Function X",
	        "denom_units":[
	            {
	                "denom":"usdt",
	                "exponent":0,
	                "aliases":[
	                    "tronTR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
	                    "polygon0xc2132D05D31c914a87C6611C10748AEb04B58e8F"
	                ]
	            },
	            {
	                "denom":"USDT",
	                "exponent":6,
	                "aliases":[]
	            }
	        ],
	        "base":"usdt",
	        "display":"usdt",
	        "name":"Tether USD",
	        "symbol":"USDT"
	    }
*/
