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
	devnetUSDTMetadata   = banktypes.Metadata{Description: "USD COIN", DenomUnits: []*banktypes.DenomUnit{{Denom: "eth0x1BE1f78d417B1C4A199bb8ad4c946Ca248f7A83e", Exponent: 6, Aliases: []string{"USDT"}}}, Base: "eth0x1BE1f78d417B1C4A199bb8ad4c946Ca248f7A83e", Display: "eth0x1BE1f78d417B1C4A199bb8ad4c946Ca248f7A83e"}

	testnetPUNDIXMetadata   = banktypes.Metadata{Description: "Pundi X Token", DenomUnits: []*banktypes.DenomUnit{{Denom: "eth0xd9EEd31F5731DfC3Ca18f09B487e200F50a6343B", Exponent: 18, Aliases: []string{"PUNDIX"}}}, Base: "eth0xd9EEd31F5731DfC3Ca18f09B487e200F50a6343B", Display: "eth0xd9EEd31F5731DfC3Ca18f09B487e200F50a6343B"}
	testnetEthUSDTMetadata  = banktypes.Metadata{Description: "Tether USD", DenomUnits: []*banktypes.DenomUnit{{Denom: "eth0xD69133f9A0206b3340d9622F2eBc4571022b3b5f", Exponent: 6, Aliases: []string{"USDT"}}}, Base: "eth0xD69133f9A0206b3340d9622F2eBc4571022b3b5f", Display: "eth0xD69133f9A0206b3340d9622F2eBc4571022b3b5f"}
	testnetUSDCMetadata     = banktypes.Metadata{Description: "USD Coin", DenomUnits: []*banktypes.DenomUnit{{Denom: "eth0xeC822cd1238d946Cf0f73be57359c5cAa5512a9D", Exponent: 6, Aliases: []string{"USDC"}}}, Base: "eth0xeC822cd1238d946Cf0f73be57359c5cAa5512a9D", Display: "eth0xeC822cd1238d946Cf0f73be57359c5cAa5512a9D"}
	testnetDAIMetadata      = banktypes.Metadata{Description: "DAI StableCoin", DenomUnits: []*banktypes.DenomUnit{{Denom: "eth0x2870405E4ABF9FcCDc93d9cC83c09788296d8354", Exponent: 18, Aliases: []string{"DAI"}}}, Base: "eth0x2870405E4ABF9FcCDc93d9cC83c09788296d8354", Display: "eth0x2870405E4ABF9FcCDc93d9cC83c09788296d8354"}
	testnetPURSEMetadata    = banktypes.Metadata{Description: "PURSE TOKEN", DenomUnits: []*banktypes.DenomUnit{{Denom: "ibc/4757BC3AA2C696F7083C825BD3951AE3D1631F2A272EA7AFB9B3E1CCCA8560D4", Exponent: 18, Aliases: []string{"PURSE"}}}, Base: "ibc/4757BC3AA2C696F7083C825BD3951AE3D1631F2A272EA7AFB9B3E1CCCA8560D4", Display: "ibc/4757BC3AA2C696F7083C825BD3951AE3D1631F2A272EA7AFB9B3E1CCCA8560D4"}
	testnetUSDJMetadata     = banktypes.Metadata{Description: "JUST Stablecoin v1.0", DenomUnits: []*banktypes.DenomUnit{{Denom: "tronTLBaRhANQoJFTqre9Nf1mjuwNWjCJeYqUL", Exponent: 18, Aliases: []string{"USDJ"}}}, Base: "tronTLBaRhANQoJFTqre9Nf1mjuwNWjCJeYqUL", Display: "tronTLBaRhANQoJFTqre9Nf1mjuwNWjCJeYqUL"}
	testnetUSDFMetadata     = banktypes.Metadata{Description: "FX USD", DenomUnits: []*banktypes.DenomUnit{{Denom: "tronTK1pM7NtkLohgRgKA6LeocW2znwJ8JtLrQ", Exponent: 6, Aliases: []string{"USDF"}}}, Base: "tronTK1pM7NtkLohgRgKA6LeocW2znwJ8JtLrQ", Display: "tronTK1pM7NtkLohgRgKA6LeocW2znwJ8JtLrQ"}
	testnetTronUSDTMetadata = banktypes.Metadata{Description: "Tether USD", DenomUnits: []*banktypes.DenomUnit{{Denom: "tronTXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj", Exponent: 6, Aliases: []string{"USDT"}}}, Base: "tronTXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj", Display: "tronTXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj"}
	testnetLINKMetadata     = banktypes.Metadata{Description: "ChainLink Token", DenomUnits: []*banktypes.DenomUnit{{Denom: "polygon0x326C977E6efc84E512bB9C30f76E30c160eD06FB", Exponent: 18, Aliases: []string{"LINK"}}}, Base: "polygon0x326C977E6efc84E512bB9C30f76E30c160eD06FB", Display: "polygon0x326C977E6efc84E512bB9C30f76E30c160eD06FB"}

	mainnetPUNDIXMetadata      = banktypes.Metadata{Description: "Pundi X Token", DenomUnits: []*banktypes.DenomUnit{{Denom: "eth0x0FD10b9899882a6f2fcb5c371E17e70FdEe00C38", Exponent: 18, Aliases: []string{"PUNDIX"}}}, Base: "eth0x0FD10b9899882a6f2fcb5c371E17e70FdEe00C38", Display: "eth0x0FD10b9899882a6f2fcb5c371E17e70FdEe00C38"}
	mainnetPURSEMetadata       = banktypes.Metadata{Description: "PURSE TOKEN", DenomUnits: []*banktypes.DenomUnit{{Denom: "ibc/F08B62C2C1BE9E52942617489CAB1E94537FE3849F8EEC910B142468C340EB0D", Exponent: 18, Aliases: []string{"PURSE"}}}, Base: "ibc/F08B62C2C1BE9E52942617489CAB1E94537FE3849F8EEC910B142468C340EB0D", Display: "ibc/F08B62C2C1BE9E52942617489CAB1E94537FE3849F8EEC910B142468C340EB0D"}
	mainnetTronUSDTMetadata    = banktypes.Metadata{Description: "Tether USD", DenomUnits: []*banktypes.DenomUnit{{Denom: "tronTR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", Exponent: 6, Aliases: []string{"USDT"}}}, Base: "tronTR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", Display: "tronTR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"}
	mainnetPolygonUSDTMetadata = banktypes.Metadata{Description: "(PoS) Tether USD", DenomUnits: []*banktypes.DenomUnit{{Denom: "polygon0xc2132D05D31c914a87C6611C10748AEb04B58e8F", Exponent: 6, Aliases: []string{"USDT"}}}, Base: "polygon0xc2132D05D31c914a87C6611C10748AEb04B58e8F", Display: "polygon0xc2132D05D31c914a87C6611C10748AEb04B58e8F"}
)

func GetMetadata() []banktypes.Metadata {
	if NetworkDevnet() == Network() {
		return []banktypes.Metadata{wfxMetadata, devnetPUNDIXMetadata, devnetPURSEMetadata, devnetUSDTMetadata}
	} else if NetworkTestnet() == Network() {
		return []banktypes.Metadata{wfxMetadata, testnetPUNDIXMetadata, testnetEthUSDTMetadata, testnetUSDCMetadata, testnetDAIMetadata,
			testnetPURSEMetadata, testnetUSDJMetadata, testnetUSDFMetadata, testnetTronUSDTMetadata, testnetLINKMetadata}
	} else {
		return []banktypes.Metadata{wfxMetadata, mainnetPUNDIXMetadata, mainnetPURSEMetadata, mainnetTronUSDTMetadata, mainnetPolygonUSDTMetadata}
	}
}
