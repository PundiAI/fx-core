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
	app *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := ctx.CacheContext()

		updateWFXLogicCode(cacheCtx, app.EvmKeeper)

		updateMetadataAliasNull(cacheCtx, app.BankKeeper)

		ctx.Logger().Info("start to run v3 migrations...", "module", "upgrade")
		toVM, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			panic(fmt.Sprintf("run migrations: %s", err.Error()))
		}

		initAvalanche(cacheCtx, app.AvalancheKeeper)

		updateBSCOracles(cacheCtx, app.BscKeeper)

		registerCoin(cacheCtx, app.Erc20Keeper)

		commit()
		ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
		return toVM, nil
	}
}

func initAvalanche(ctx sdk.Context, avalancheKeeper crosschainkeeper.Keeper) {
	if ctx.ChainID() == fxtypes.TestnetChainId {
		params := avalancheKeeper.GetParams(ctx)
		params.GravityId = fmt.Sprintf("%s-test", params.GravityId)
		if err := avalancheKeeper.SetParams(ctx, &params); err != nil {
			panic(err)
		}
	}
	var newOracles []string
	chainId := ctx.ChainID()
	if chainId == fxtypes.MainnetChainId {
		newOracles = []string{
			"fx1y45zkm5fedeglz5vc9kf48x3dlfdc8mm7vgyss",
			"fx13yc08ua0vzlccefzk8pwae3lsr40rten7s067f",
			"fx18q29l6rlfp6c4h447233sth5vcwgcwx2fnhgrp",
			"fx1qd08zkx9mz6qzwm7w2m9n6uzhzzegcwpzfyyaa",
			"fx1p8utlnwa9uyul6pg5hkhw47hvpsgedg66nke8f",
			"fx1cdet8s7dfmuxeu535uvcs9hx2ldl064eq2lzyt",
			"fx1htxadktncppq685v93vz8exdzhsvu060n9d5n6",
			"fx1j4feg7g6qr27v6ytm0avj5wh0mrydsl8pcvjqj",
			"fx1am8zdzjk0z47tjtw83k47s48e3g8v2l3yh0vqe",
			"fx1905t0cyarp890gcp752fq3kpd3t3fnwf397f57",
		}
	} else if chainId == fxtypes.TestnetChainId {
		newOracles = []string{
			"fx1q4avdlyhxhzq3l2ngux2tpmz7jwl5mnkycnxve",
			"fx13s5dyfagdyv2vcf25gw5rl849w5e93kztf5t5f",
			"fx1wmakpdj7u3cf9anqq0u552qnm2uef50fgj7wnz",
			"fx1ehd44azw0cu8u2kljhkfkccfc4xkjpl7nlredz",
			"fx1fcytwf6netk6nvftan5wnu7jsv06x0gxuq4avf",
		}
	} else {
		panic("invalid chainId:" + chainId)
	}
	if len(newOracles) <= 0 {
		return
	}
	oracles := make([]string, 0, len(newOracles))
	for _, oracle := range newOracles {
		if _, err := sdk.AccAddressFromBech32(oracle); err != nil {
			ctx.Logger().Error("parse avalanche oracle address error", "module", "upgrade", "error", err.Error())
			continue
		}
		oracles = append(oracles, oracle)
	}
	if len(oracles) <= 0 {
		return
	}
	ctx.Logger().Info("init module avalanche oracles", "module", "upgrade", "number", len(oracles))
	avalancheKeeper.SetProposalOracle(ctx, &crosschaintypes.ProposalOracle{
		Oracles: oracles,
	})
}

func updateBSCOracles(ctx sdk.Context, bscKeeper crosschainkeeper.Keeper) {
	newOracles := getBSCOracleAddrs(ctx.ChainID())
	if len(newOracles) <= 0 {
		return
	}
	// append old oracle
	onlineOracles := bscKeeper.GetAllOracles(ctx, true)
	for _, onlineOracle := range onlineOracles {
		var isExist bool
		for _, newOracle := range newOracles {
			if onlineOracle.OracleAddress == newOracle {
				isExist = true
				break
			}
		}
		if !isExist {
			newOracles = append(newOracles, onlineOracle.OracleAddress)
		}
	}
	oracles := make([]string, 0, len(newOracles))
	for _, oracle := range newOracles {
		if _, err := sdk.AccAddressFromBech32(oracle); err != nil {
			ctx.Logger().Error("parse bsc oracle address error", "module", "upgrade", "error", err.Error())
			continue
		}
		oracles = append(oracles, oracle)
	}
	if len(oracles) <= 0 {
		return
	}
	ctx.Logger().Info("update module bsc oracles to", "module", "upgrade", "number", len(oracles))
	bscKeeper.SetProposalOracle(ctx, &crosschaintypes.ProposalOracle{
		Oracles: oracles,
	})
}

func registerCoin(ctx sdk.Context, k erc20keeper.Keeper) {
	for _, metadata := range GetMetadata(ctx.ChainID()) {
		cacheCtx, commit := ctx.CacheContext()
		pair, err := k.RegisterCoin(cacheCtx, metadata)
		if err != nil {
			// run time error, non-fatal, print info
			ctx.Logger().Error("failed to register coin", "module", "upgrade", "denom", metadata.Base, "error", err.Error())
			continue
		}
		commit()
		ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
		ctx.Logger().Info("add metadata successfully", "module", "upgrade", "metadata", metadata.String())
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			erc20types.EventTypeRegisterCoin,
			sdk.NewAttribute(erc20types.AttributeKeyDenom, pair.Denom),
			sdk.NewAttribute(erc20types.AttributeKeyTokenAddress, pair.Erc20Address),
		))
	}
}

func updateWFXLogicCode(ctx sdk.Context, k *evmkeeper.Keeper) {
	wfx := fxtypes.GetWFX()
	if err := k.UpdateContractCode(ctx, wfx.Address, wfx.Code); err != nil {
		panic(fmt.Sprintf("update wfx logic code error: %s", err.Error()))
	}
	ctx.Logger().Info("update WFX contract", "module", "upgrade", "codeHash", wfx.CodeHash())
}

func updateMetadataAliasNull(ctx sdk.Context, bk bankkeeper.Keeper) {
	bk.IterateAllDenomMetaData(ctx, func(md banktypes.Metadata) bool {
		if len(md.DenomUnits) != 2 || len(md.DenomUnits[1].Aliases) != 1 || md.DenomUnits[1].Aliases[0] != "null" {
			return false
		}
		ctx.Logger().Info("fix metadata alias", "module", "upgrade", "denom", md.Base)
		md.DenomUnits[1].Aliases = []string{}
		bk.SetDenomMetaData(ctx, md)
		return false
	})
}

func getBSCOracleAddrs(chainId string) []string {
	var oracles []string
	if chainId == fxtypes.MainnetChainId {
		oracles = []string{
			"fx1qa9h4hm6parfw4250mydjc9guf42u8c80c30wd",
			"fx1klx55m5pwphmfz0lwtszztm8wm5jewe6cv9prr",
			"fx13tj0ck5huf5dw0mue62vcsg2rken86lmtkynst",
			"fx10mhm00rg9hrdukl32jyypjt2rvredmfxcr2e43",
			"fx129h996ynl62z9hqh44fjc4pz46q0w4xcesvvlg",
			"fx1udd9lctranrs48ekf8t4cck9tdgvpew4r2cmag",
			"fx1589zrr7tld7usp25hn8p0clmv37nejzyhcyrj3",
			"fx170ev804dcn0vuvw9d25mwuyx6xy4cp9n2w2ykt",
			"fx1fkq2tmwcpuxdd0wv8t60y884gtpvjp43c937x2",
			"fx1sr9x23zt6sag0hdalgv8wu3m283kwuzpghthc3",
		}
	} else if chainId == fxtypes.TestnetChainId {
		oracles = []string{
			"fx1v55r4dl0l35ra4hgjsp9hq4skzkpc6z7hvnrcv",
			"fx1l2nqwmhw8xw2y68yzucs4nvs2mxdd7ff5jn53e",
			"fx1vavhtkdycrxrsa5gfr53gn90xktvma4ystd8na",
			"fx1cajtzkv5mu2jhl5q7c6qwqxh2d0zlylyuppf2f",
			"fx1zfvcvl4hk7rl0zgesuqx7n37zr0q6c9hpjk4jc",
		}
	} else {
		panic("invalid chainId:" + chainId)
	}
	return oracles
}

func GetMetadata(chainId string) []banktypes.Metadata {
	if fxtypes.TestnetChainId == chainId {
		return []banktypes.Metadata{
			fxtypes.GetCrossChainMetadata("Wrapped AVAX", "WAVAX", 18, "avalanche0xd0fABb17BD2999A4A9fDF0F05c2386e7dF6519bb"),
			fxtypes.GetCrossChainMetadata("Staked AVAX", "sAVAX", 18, "avalanche0x57b1E4C85B0f141aDE38b5573907BA8eF9aC2298"),
			fxtypes.GetCrossChainMetadata("BENQI", "QI", 18, "avalanche0xeb62B336778ac9E9CF1Aacfd268E0Eb013019DC5"),
			fxtypes.GetCrossChainMetadata("BavaToken", "BAVA", 18, "avalanche0x52583B59A9458667b82358A7ac07b6d26f83A2A4"),
			fxtypes.GetCrossChainMetadata("Wrapped BTC", "WBTC", 8, "eth0x6895a336ccC9086aD026a83B93073960622b35B4"),
		}
	} else if chainId == fxtypes.MainnetChainId {
		return []banktypes.Metadata{
			fxtypes.GetCrossChainMetadata("Wrapped AVAX", "WAVAX", 18, "avalanche0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7"),
			fxtypes.GetCrossChainMetadata("Staked AVAX", "sAVAX", 18, "avalanche0x2b2C81e08f1Af8835a78Bb2A90AE924ACE0eA4bE"),
			fxtypes.GetCrossChainMetadata("BENQI", "QI", 18, "avalanche0x8729438EB15e2C8B576fCc6AeCdA6A148776C0F5"),
			fxtypes.GetCrossChainMetadata("BavaToken", "BAVA", 18, "avalanche0xe19A1684873faB5Fb694CfD06607100A632fF21c"),
			fxtypes.GetCrossChainMetadata("Wrapped BTC", "WBTC", 8, "eth0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599"),
		}
	} else {
		panic("invalid chainId:" + chainId)
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
