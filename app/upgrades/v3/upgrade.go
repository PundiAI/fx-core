package v3

import (
	"fmt"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"

	"github.com/functionx/fx-core/v3/app/keepers"
	fxcfg "github.com/functionx/fx-core/v3/server/config"
	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschainkeeper "github.com/functionx/fx-core/v3/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	erc20keeper "github.com/functionx/fx-core/v3/x/erc20/keeper"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	evmkeeper "github.com/functionx/fx-core/v3/x/evm/keeper"
)

func createUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// cache context
		cacheCtx, commit := ctx.CacheContext()

		// update wfx logic code
		updateWFXLogicCode(cacheCtx, keepers.EvmKeeper)

		// update metadata alias null
		updateMetadataAliasNull(cacheCtx, keepers.BankKeeper)

		// run migrations
		toVM := runMigrations(cacheCtx, fromVM, mm, configurator)

		// init avalanche oracles
		initAvalancheOracles(cacheCtx, keepers.AvalancheKeeper)

		// update bsc oracles
		updateBSCOracles(cacheCtx, keepers.BscKeeper)

		// register coin
		registerCoin(cacheCtx, keepers.Erc20Keeper)

		// commit upgrade
		commit()
		ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
		return toVM, nil
	}
}

func initAvalancheOracles(ctx sdk.Context, avalancheKeeper crosschainkeeper.Keeper) {
	var oracles []string
	chainId := ctx.ChainID()
	// todo need add oracles
	if chainId == fxtypes.MainnetChainId {
		oracles = []string{}
	} else if chainId == fxtypes.TestnetChainId {
		oracles = []string{}
	} else {
		panic("invalid chainId:" + chainId)
	}
	avalancheKeeper.SetProposalOracle(ctx, &crosschaintypes.ProposalOracle{
		Oracles: oracles,
	})
}

func updateBSCOracles(ctx sdk.Context, bscKeeper crosschainkeeper.Keeper) {
	var oracles []string
	chainId := ctx.ChainID()
	// todo need add oracles
	if chainId == fxtypes.MainnetChainId {
		oracles = []string{}
	} else if chainId == fxtypes.TestnetChainId {
		oracles = []string{}
	} else {
		panic("invalid chainId:" + chainId)
	}
	if len(oracles) <= 0 {
		return
	}
	bscKeeper.SetProposalOracle(ctx, &crosschaintypes.ProposalOracle{
		Oracles: oracles,
	})
}

func runMigrations(ctx sdk.Context, fromVM module.VersionMap, mm *module.Manager, mc module.Configurator) module.VersionMap {
	ctx.Logger().Info("start to run module v3 migrations...")
	toVM, err := mm.RunMigrations(ctx, mc, fromVM)
	if err != nil {
		panic(fmt.Sprintf("run migrations: %s", err.Error()))
	}
	return toVM
}

func registerCoin(ctx sdk.Context, k erc20keeper.Keeper) {
	for _, metadata := range getMetadata(ctx.ChainID()) {
		ctx.Logger().Info("add metadata", "coin", metadata.String())
		pair, err := k.RegisterCoin(ctx, metadata)
		if err != nil {
			panic(fmt.Sprintf("register %s: %s", metadata.Base, err.Error()))
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			erc20types.EventTypeRegisterCoin,
			sdk.NewAttribute(erc20types.AttributeKeyDenom, pair.Denom),
			sdk.NewAttribute(erc20types.AttributeKeyTokenAddress, pair.Erc20Address),
		))
	}
}

func updateWFXLogicCode(ctx sdk.Context, k *evmkeeper.Keeper) {
	wfx := fxtypes.GetWFX()
	err := k.UpdateContractCode(ctx, wfx.Address, wfx.Code)
	if err != nil {
		panic(fmt.Sprintf("update wfx logic code error: %s", err.Error()))
	}
}

func updateMetadataAliasNull(ctx sdk.Context, bk bankkeeper.Keeper) {
	bk.IterateAllDenomMetaData(ctx, func(md banktypes.Metadata) bool {
		if len(md.DenomUnits) != 2 || len(md.DenomUnits[1].Aliases) != 1 || md.DenomUnits[1].Aliases[0] != "null" {
			return false
		}
		ctx.Logger().Info("fix metadata alias", "denom", md.Base)
		md.DenomUnits[1].Aliases = []string{}
		bk.SetDenomMetaData(ctx, md)
		return false
	})
}

func getMetadata(chainId string) []banktypes.Metadata {
	if fxtypes.TestnetChainId == chainId {
		return []banktypes.Metadata{
			fxtypes.GetCrossChainMetadata("Wrapped AVAX", "WAVAX", 18, "avalanche0xd0fABb17BD2999A4A9fDF0F05c2386e7dF6519bb"),
			fxtypes.GetCrossChainMetadata("Staked AVAX", "sAVAX", 18, "avalanche0x57b1E4C85B0f141aDE38b5573907BA8eF9aC2298"),
			fxtypes.GetCrossChainMetadata("BENQI", "QI", 18, "avalanche0xeb62B336778ac9E9CF1Aacfd268E0Eb013019DC5"),
			fxtypes.GetCrossChainMetadata("BavaToken", "BAVA", 18, "avalanche0x52583B59A9458667b82358A7ac07b6d26f83A2A4"),
			fxtypes.GetCrossChainMetadata("Wrapped BTC", "WBTC", 8, "eth0x6895a336ccC9086aD026a83B93073960622b35B4"),
		}
	} else {
		return []banktypes.Metadata{
			fxtypes.GetCrossChainMetadata("Wrapped AVAX", "WAVAX", 18, "avalanche0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7"),
			fxtypes.GetCrossChainMetadata("Staked AVAX", "sAVAX", 18, "avalanche0x2b2C81e08f1Af8835a78Bb2A90AE924ACE0eA4bE"),
			fxtypes.GetCrossChainMetadata("BENQI", "QI", 18, "avalanche0x8729438EB15e2C8B576fCc6AeCdA6A148776C0F5"),
			fxtypes.GetCrossChainMetadata("BavaToken", "BAVA", 18, "avalanche0xe19A1684873faB5Fb694CfD06607100A632fF21c"),
			fxtypes.GetCrossChainMetadata("Wrapped BTC", "WBTC", 8, "eth0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599"),
		}
	}
}

// preUpgradeCmd called by cosmovisor
func preUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pre-upgrade",
		Short: "fxv3 pre-upgrade, called by cosmovisor, before migrations upgrade",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			serverCtx.Logger.Info("pre-upgrade", "action", "update app.toml and config.toml")

			rootDir := serverCtx.Config.RootDir
			fileName := filepath.Join(rootDir, "config", "config.toml")
			tmcfg.WriteConfigFile(fileName, serverCtx.Config)

			config.SetConfigTemplate(fxcfg.DefaultConfigTemplate())
			appConfig := fxcfg.DefaultConfig()
			if err := serverCtx.Viper.Unmarshal(appConfig); err != nil {
				return err
			}
			fileName = filepath.Join(rootDir, "config", "app.toml")
			config.WriteConfigFile(fileName, appConfig)

			clientCtx := client.GetClientContextFromCmd(cmd)
			return clientCtx.PrintString("fxv3 pre-upgrade success")
		},
	}
	return cmd
}
