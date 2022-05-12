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

var (
	/*
		//// example of origin denom
		{
			Description: "Function X",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    "FX",
					Exponent: 18,
				},
			},
			Base:    "FX",
			Display: "FX",
		}
		//// example of other denom
		{
			Description: "Pundi X Token",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    "eth0x338E7A8687AdA7274Dc87C95D94f920d8F4185AE",
					Exponent: 18,
					Aliases:  []string{"PUNDIX"}, //symbol
				},
			},
			Base:    "eth0x338E7A8687AdA7274Dc87C95D94f920d8F4185AE",
			Display: "eth0x338E7A8687AdA7274Dc87C95D94f920d8F4185AE",
		}
	*/

	wfxMetadata = banktypes.Metadata{Description: "Function X", DenomUnits: []*banktypes.DenomUnit{{Denom: "FX", Exponent: 18}}, Base: "FX", Display: "FX"}

	devnetPUNDIXMetadata = banktypes.Metadata{Description: "Pundi X Token", DenomUnits: []*banktypes.DenomUnit{{Denom: "eth0x338E7A8687AdA7274Dc87C95D94f920d8F4185AE", Exponent: 18, Aliases: []string{"PUNDIX"}}}, Base: "eth0x338E7A8687AdA7274Dc87C95D94f920d8F4185AE", Display: "eth0x338E7A8687AdA7274Dc87C95D94f920d8F4185AE"}
	devnetPURSEMetadata  = banktypes.Metadata{Description: "PURSE TOKEN", DenomUnits: []*banktypes.DenomUnit{{Denom: "ibc/B1861D0C2E4BAFA42A61739291975B7663F278FFAF579F83C9C4AD3890D09CA0", Exponent: 18, Aliases: []string{"PURSE"}}}, Base: "ibc/B1861D0C2E4BAFA42A61739291975B7663F278FFAF579F83C9C4AD3890D09CA0", Display: "ibc/B1861D0C2E4BAFA42A61739291975B7663F278FFAF579F83C9C4AD3890D09CA0"}
	devnetUSDTMetadata   = banktypes.Metadata{Description: "USD COIN", DenomUnits: []*banktypes.DenomUnit{{Denom: "eth0x1BE1f78d417B1C4A199bb8ad4c946Ca248f7A83e", Exponent: 18, Aliases: []string{"USDT"}}}, Base: "eth0x1BE1f78d417B1C4A199bb8ad4c946Ca248f7A83e", Display: "eth0x1BE1f78d417B1C4A199bb8ad4c946Ca248f7A83e"}

	testnetPUNDIXMetadata = banktypes.Metadata{Description: "Pundi X Token", DenomUnits: []*banktypes.DenomUnit{{Denom: "PUNDIX", Exponent: 18, Aliases: []string{"eth0xd9EEd31F5731DfC3Ca18f09B487e200F50a6343B"}}}, Base: "PUNDIX", Display: "PUNDIX"}
	testnetPURSEMetadata  = banktypes.Metadata{Description: "PURSE TOKEN", DenomUnits: []*banktypes.DenomUnit{{Denom: "PURSE", Exponent: 18, Aliases: []string{"ibc/4757BC3AA2C696F7083C825BD3951AE3D1631F2A272EA7AFB9B3E1CCCA8560D4"}}}, Base: "PURSE", Display: "PURSE"}
	testnetUSDJMetadata   = banktypes.Metadata{Description: "JUST Stablecoin", DenomUnits: []*banktypes.DenomUnit{{Denom: "USDJ", Exponent: 18, Aliases: []string{"tronTLBaRhANQoJFTqre9Nf1mjuwNWjCJeYqUL"}}}, Base: "USDJ", Display: "USDJ"}
	testnetUSDFMetadata   = banktypes.Metadata{Description: "FX USD", DenomUnits: []*banktypes.DenomUnit{{Denom: "USDF", Exponent: 18, Aliases: []string{"tronTK1pM7NtkLohgRgKA6LeocW2znwJ8JtLrQ"}}}, Base: "USDF", Display: "USDF"}
	testnetLINKMetadata   = banktypes.Metadata{Description: "ChainLink Token", DenomUnits: []*banktypes.DenomUnit{{Denom: "LINK", Exponent: 18, Aliases: []string{"polygon0x326C977E6efc84E512bB9C30f76E30c160eD06FB"}}}, Base: "LINK", Display: "LINK"}

	mainnetPUNDIXMetadata      = banktypes.Metadata{Description: "Pundi X Token", DenomUnits: []*banktypes.DenomUnit{{Denom: "eth0x0FD10b9899882a6f2fcb5c371E17e70FdEe00C38", Exponent: 18, Aliases: []string{"PUNDIX"}}}, Base: "eth0x0FD10b9899882a6f2fcb5c371E17e70FdEe00C38", Display: "eth0x0FD10b9899882a6f2fcb5c371E17e70FdEe00C38"}
	mainnetPURSEMetadata       = banktypes.Metadata{Description: "PURSE TOKEN", DenomUnits: []*banktypes.DenomUnit{{Denom: "ibc/F08B62C2C1BE9E52942617489CAB1E94537FE3849F8EEC910B142468C340EB0D", Exponent: 18, Aliases: []string{"PURSE"}}}, Base: "ibc/F08B62C2C1BE9E52942617489CAB1E94537FE3849F8EEC910B142468C340EB0D", Display: "ibc/F08B62C2C1BE9E52942617489CAB1E94537FE3849F8EEC910B142468C340EB0D"}
	mainnetTronUSDTMetadata    = banktypes.Metadata{Description: "Tether USD", DenomUnits: []*banktypes.DenomUnit{{Denom: "tronTR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", Exponent: 18, Aliases: []string{"USDT"}}}, Base: "tronTR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", Display: "tronTR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"}
	mainnetPolygonUSDTMetadata = banktypes.Metadata{Description: "(PoS) Tether USD", DenomUnits: []*banktypes.DenomUnit{{Denom: "polygon0xc2132D05D31c914a87C6611C10748AEb04B58e8F", Exponent: 18, Aliases: []string{"USDT"}}}, Base: "polygon0xc2132D05D31c914a87C6611C10748AEb04B58e8F", Display: "polygon0xc2132D05D31c914a87C6611C10748AEb04B58e8F"}
)

func GetMetadata() []banktypes.Metadata {
	if NetworkDevnet() == Network() {
		return []banktypes.Metadata{wfxMetadata, devnetPUNDIXMetadata, devnetPURSEMetadata, devnetUSDTMetadata}
	} else if NetworkTestnet() == Network() {
		return []banktypes.Metadata{wfxMetadata, testnetPUNDIXMetadata, testnetPURSEMetadata, testnetUSDJMetadata, testnetUSDFMetadata, testnetLINKMetadata}
	} else {
		return []banktypes.Metadata{wfxMetadata, mainnetPUNDIXMetadata, mainnetPURSEMetadata, mainnetTronUSDTMetadata, mainnetPolygonUSDTMetadata}
	}
}
