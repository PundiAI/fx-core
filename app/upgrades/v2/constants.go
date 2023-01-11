package v2

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"github.com/functionx/fx-core/v3/app/upgrades"
	fxtypes "github.com/functionx/fx-core/v3/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	migratetypes "github.com/functionx/fx-core/v3/x/migrate/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

// Upgrade nolint
var Upgrade = upgrades.Upgrade{
	CreateUpgradeHandler: createUpgradeHandler,
	PreUpgradeCmd:        preUpgradeCmd(),
	StoreUpgrades: func() *storetypes.StoreUpgrades {
		if fxtypes.ChainId() == fxtypes.TestnetChainId {
			return &storetypes.StoreUpgrades{
				Added: []string{
					feegrant.StoreKey,
					authzkeeper.StoreKey,
				},
			}
		}
		return &storetypes.StoreUpgrades{
			Added: []string{
				feemarkettypes.StoreKey,
				evmtypes.StoreKey,
				erc20types.StoreKey,
				migratetypes.StoreKey,
				feegrant.StoreKey,
				authzkeeper.StoreKey,
			},
		}
	},
}

var (
	initGenesis = map[string]bool{
		feemarkettypes.ModuleName: true,
		evmtypes.ModuleName:       true,
		erc20types.ModuleName:     true,
		migratetypes.ModuleName:   true,
		feegrant.ModuleName:       true,
		authz.ModuleName:          true,
	}

	runMigrates = map[string]uint64{
		authtypes.ModuleName:         1,
		banktypes.ModuleName:         1,
		distributiontypes.ModuleName: 1,
		govtypes.ModuleName:          1,
		slashingtypes.ModuleName:     1,
		stakingtypes.ModuleName:      1,
		ibchost.ModuleName:           1,
		bsctypes.ModuleName:          1,
		polygontypes.ModuleName:      1,
		trontypes.ModuleName:         1,
	}

	crossChainModule = map[string]bool{
		bsctypes.ModuleName:     true,
		polygontypes.ModuleName: true,
		trontypes.ModuleName:    true,
	}
)

func getMetadata(chainId string) []banktypes.Metadata {
	fxMetaData := fxtypes.GetFXMetaData(fxtypes.DefaultDenom)
	if fxtypes.TestnetChainId == chainId {
		return []banktypes.Metadata{
			fxMetaData,
			getCrossChainMetadata("Pundi X Token", "PUNDIX", 18, "eth0xd9EEd31F5731DfC3Ca18f09B487e200F50a6343B"),
			getCrossChainMetadata("Tether USD", "USDT", 6, "eth0xD69133f9A0206b3340d9622F2eBc4571022b3b5f"),
			getCrossChainMetadata("USD Coin", "USDC", 18, "eth0xeC822cd1238d946Cf0f73be57359c5cAa5512a9D"),
			getCrossChainMetadata("DAI StableCoin", "DAI", 18, "eth0x2870405E4ABF9FcCDc93d9cC83c09788296d8354"),
			getCrossChainMetadata("PURSE TOKEN", "PURSE", 18, "ibc/4757BC3AA2C696F7083C825BD3951AE3D1631F2A272EA7AFB9B3E1CCCA8560D4"),
			getCrossChainMetadata("JUST Stablecoin v1.0", "USDJ", 18, "tronTLBaRhANQoJFTqre9Nf1mjuwNWjCJeYqUL"),
			getCrossChainMetadata("FX USD", "USDF", 6, "tronTK1pM7NtkLohgRgKA6LeocW2znwJ8JtLrQ"),
			getCrossChainMetadata("Tether USD", "USDT", 6, "tronTXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj"),
			getCrossChainMetadata("ChainLink Token", "LINK", 18, "polygon0x326C977E6efc84E512bB9C30f76E30c160eD06FB"),
		}
	} else {
		return []banktypes.Metadata{
			fxMetaData,
			getCrossChainMetadata("Pundi X Token", "PUNDIX", 18, "eth0x0FD10b9899882a6f2fcb5c371E17e70FdEe00C38"),
			getCrossChainMetadata("PURSE TOKEN", "PURSE", 18, "ibc/F08B62C2C1BE9E52942617489CAB1E94537FE3849F8EEC910B142468C340EB0D"),
		}
	}
}

func getCrossChainMetadata(name, symbol string, decimals uint32, denom string) banktypes.Metadata {
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
