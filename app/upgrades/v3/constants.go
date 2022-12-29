package v3

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/functionx/fx-core/v3/app/upgrades"
	fxtypes "github.com/functionx/fx-core/v3/types"
	avalanchetypes "github.com/functionx/fx-core/v3/x/avalanche/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
)

const (
	UpgradeName = "fxv3"

	DAIDenom  = "dai"
	EURSDenom = "eurs"
	LINKDenom = "link"
	C98Denom  = "c98"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: func() *store.StoreUpgrades {
		return &store.StoreUpgrades{
			Added: []string{
				avalanchetypes.ModuleName,
				ethtypes.ModuleName,
			},
			Deleted: []string{
				"other",
			},
		}
	},
}

func GetMetadata(chainId string) []banktypes.Metadata {
	var (
		//fxv3 token TODO update testnet denom
		testnetWAVAXMetadata = fxtypes.GetCrossChainMetadataManyToOne("Wrapped AVAX", "WAVAX", 18, "avalanche0x0000000000000000000000000000000000000001")
		testnetSAVAXMetadata = fxtypes.GetCrossChainMetadataManyToOne("Staked AVAX", "sAVAX", 18, "avalanche0x0000000000000000000000000000000000000002")
		testnetQIMetadata    = fxtypes.GetCrossChainMetadataManyToOne("BENQI", "QI", 18, "avalanche0x0000000000000000000000000000000000000003")
		testnetBAVAMetadata  = fxtypes.GetCrossChainMetadataManyToOne("BavaToken", "BAVA", 18, "avalanche0x0000000000000000000000000000000000000004")
		testnetWBTCMetadata  = fxtypes.GetCrossChainMetadataManyToOne("Wrapped BTC", "WBTC", 8, "eth0x0000000000000000000000000000000000000005")

		mainnetWAVAXMetadata = fxtypes.GetCrossChainMetadataManyToOne("Wrapped AVAX", "WAVAX", 18, "avalanche0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7")
		mainnetSAVAXMetadata = fxtypes.GetCrossChainMetadataManyToOne("Staked AVAX", "sAVAX", 18, "avalanche0x2b2C81e08f1Af8835a78Bb2A90AE924ACE0eA4bE")
		mainnetQIMetadata    = fxtypes.GetCrossChainMetadataManyToOne("BENQI", "QI", 18, "avalanche0x8729438EB15e2C8B576fCc6AeCdA6A148776C0F5")
		mainnetBAVAMetadata  = fxtypes.GetCrossChainMetadataManyToOne("BavaToken", "BAVA", 18, "avalanche0xe19A1684873faB5Fb694CfD06607100A632fF21c")
		mainnetWBTCMetadata  = fxtypes.GetCrossChainMetadataManyToOne("Wrapped BTC", "WBTC", 8, "eth0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599")
	)
	if fxtypes.TestnetChainId == chainId {
		return []banktypes.Metadata{testnetWAVAXMetadata, testnetSAVAXMetadata, testnetQIMetadata, testnetBAVAMetadata, testnetWBTCMetadata}
	} else {
		return []banktypes.Metadata{mainnetWAVAXMetadata, mainnetSAVAXMetadata, mainnetQIMetadata, mainnetBAVAMetadata, mainnetWBTCMetadata}
	}
}
