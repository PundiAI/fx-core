package v4

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/functionx/fx-core/v4/app/upgrades"
	fxtypes "github.com/functionx/fx-core/v4/types"
	arbitrumtypes "github.com/functionx/fx-core/v4/x/arbitrum/types"
	gravitytypes "github.com/functionx/fx-core/v4/x/gravity/types"
	optimismtypes "github.com/functionx/fx-core/v4/x/optimism/types"
)

var (
	mainnetRemoveOracles = []string{
		"fx1d6xj9yekmpmyafwx9tvgzl2gmzcpe04n54853f",
		"fx1jjpwenetj0u70peyyjy9lewwzjyg2y7suq3hct",
	}

	testnetUpdateDenomAlias = []DenomAlias{
		{Denom: "weth", Alias: "arbitrum0x57b1E4C85B0f141aDE38b5573907BA8eF9aC2298"},
		{Denom: "usdt", Alias: "arbitrum0xEa99760Ecc3460154670B86E202233974883b153"},
		{Denom: "weth", Alias: "optimism0xd0fABb17BD2999A4A9fDF0F05c2386e7dF6519bb"},
		{Denom: "usdt", Alias: "optimism0xeb62B336778ac9E9CF1Aacfd268E0Eb013019DC5"},
	}
	mainnetUpdateDenomAlias = []DenomAlias{
		{Denom: "weth", Alias: "arbitrum0x82aF49447D8a07e3bd95BD0d56f35241523fBab1"},
		{Denom: "usdt", Alias: "arbitrum0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9"},
		{Denom: "weth", Alias: "optimism0x4200000000000000000000000000000000000006"},
		{Denom: "usdt", Alias: "optimism0x94b008aA00579c1307B0EF2c499aD98a8ce58e58"},
	}
)

// Deprecated: Please use v4.1.x
var Upgrade = upgrades.Upgrade{
	UpgradeName:          "fxv4",
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: func() *storetypes.StoreUpgrades {
		return &storetypes.StoreUpgrades{
			Added: []string{
				arbitrumtypes.ModuleName,
				optimismtypes.ModuleName,
			},
			Deleted: []string{
				gravitytypes.ModuleName,
			},
		}
	},
}

type DenomAlias struct {
	Denom string
	Alias string
}

func GetUpdateDenomAlias(chainId string) []DenomAlias {
	if fxtypes.TestnetChainId == chainId {
		return testnetUpdateDenomAlias
	} else if chainId == fxtypes.MainnetChainId {
		return mainnetUpdateDenomAlias
	} else {
		panic("invalid chainId:" + chainId)
	}
}

func GetBscRemoveOracles(chainId string) []string {
	if fxtypes.TestnetChainId == chainId {
		return []string{}
	} else if chainId == fxtypes.MainnetChainId {
		return mainnetRemoveOracles
	} else {
		panic("invalid chainId:" + chainId)
	}
}
