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

	EventUpdateContract  = "update_contract"
	AttributeKeyContract = "contract"
	AttributeKeyVersion  = "version"

	DAIDenom  = "dai"
	EURSDenom = "eurs"
	LINKDenom = "link"
	C98Denom  = "c98"

	EventUpdateMetadata = "update_metadata"
	AttributeKeyDenom   = "denom"
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
		//fxv3 token TODO update testnet and mainnet token denom
		testnetWAVAXMetadata = fxtypes.GetCrossChainMetadataManyToOne("Wrapped AVAX", "WAVAX", 18, "avax0x0000000000000000000000000000000000000001")
		testnetSAVAXMetadata = fxtypes.GetCrossChainMetadataManyToOne("Staked AVAX", "sAVAX", 18, "avax0x0000000000000000000000000000000000000002")
		testnetQIMetadata    = fxtypes.GetCrossChainMetadataManyToOne("BENQI", "QI", 18, "avax0x0000000000000000000000000000000000000003")
		testnetBAVAMetadata  = fxtypes.GetCrossChainMetadataManyToOne("BavaToken", "BAVA", 18, "avax0x0000000000000000000000000000000000000004")
		testnetWBTCMetadata  = fxtypes.GetCrossChainMetadataManyToOne("Wrapped BTC", "WBTC", 8, "eth0x0000000000000000000000000000000000000005")

		mainnetWAVAXMetadata = fxtypes.GetCrossChainMetadataManyToOne("Wrapped AVAX", "WAVAX", 18, "avax0x0000000000000000000000000000000000000001")
		mainnetSAVAXMetadata = fxtypes.GetCrossChainMetadataManyToOne("Staked AVAX", "sAVAX", 18, "avax0x0000000000000000000000000000000000000002")
		mainnetQIMetadata    = fxtypes.GetCrossChainMetadataManyToOne("BENQI", "QI", 18, "avax0x0000000000000000000000000000000000000003")
		mainnetBAVAMetadata  = fxtypes.GetCrossChainMetadataManyToOne("BavaToken", "BAVA", 18, "avax0x0000000000000000000000000000000000000004")
		mainnetWBTCMetadata  = fxtypes.GetCrossChainMetadataManyToOne("Wrapped BTC", "WBTC", 8, "eth0x0000000000000000000000000000000000000005")
	)
	if fxtypes.TestnetChainId == chainId {
		return []banktypes.Metadata{testnetWAVAXMetadata, testnetSAVAXMetadata, testnetQIMetadata, testnetBAVAMetadata, testnetWBTCMetadata}
	} else {
		return []banktypes.Metadata{mainnetWAVAXMetadata, mainnetSAVAXMetadata, mainnetQIMetadata, mainnetBAVAMetadata, mainnetWBTCMetadata}
	}
}

func GetAliasNullDenom() []string {
	if fxtypes.ChainId() == fxtypes.TestnetChainId {
		return []string{}
	}
	return []string{DAIDenom, EURSDenom, LINKDenom, C98Denom}
}
