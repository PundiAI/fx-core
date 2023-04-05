package v4

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"

	"github.com/functionx/fx-core/v3/app/keepers"
	fxcfg "github.com/functionx/fx-core/v3/server/config"
	fxtypes "github.com/functionx/fx-core/v3/types"
	erc20keeper "github.com/functionx/fx-core/v3/x/erc20/keeper"
	"github.com/functionx/fx-core/v3/x/gov/keeper"
	govtypes "github.com/functionx/fx-core/v3/x/gov/types"
)

func createUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	app *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := ctx.CacheContext()

		// 1. initialize the evm module account
		CreateEvmModuleAccount(cacheCtx, app.AccountKeeper)

		// 2. init go fx params
		InitGovFXParams(cacheCtx, app.GovKeeper)

		ctx.Logger().Info("start to run v4 migrations...", "module", "upgrade")
		toVM, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			panic(fmt.Sprintf("run migrations: %s", err.Error()))
		}

		// update arbitrum and optimism denom alias
		UpdateDenomAliases(cacheCtx, app.Erc20Keeper)

		commit()
		return toVM, nil
	}
}

func InitGovFXParams(ctx sdk.Context, keeper keeper.Keeper) {
	if err := keeper.SetParams(ctx, govtypes.DefaultParams()); err != nil {
		panic(err)
	}
}

func CreateEvmModuleAccount(ctx sdk.Context, k authkeeper.AccountKeeper) {
	account, _ := k.GetModuleAccountAndPermissions(ctx, evmtypes.ModuleName)
	if account == nil {
		panic("create evm module account empty")
	}
}

func UpdateDenomAliases(ctx sdk.Context, k erc20keeper.Keeper) {
	denomAlias := GetUpdateDenomAlias(ctx.ChainID())
	for _, da := range denomAlias {
		cacheCtx, commit := ctx.CacheContext()

		addFlag, err := k.UpdateDenomAliases(cacheCtx, da.Denom, da.Alias)
		if err != nil {
			ctx.Logger().Error("failed to update denom alias", "denom", da.Denom, "alias", da.Alias, "err", err.Error())
			continue
		}
		commit()
		ctx.Logger().Info("update denom alias successfully", "denom", da.Denom, "alias", da.Alias, "add-flag", strconv.FormatBool(addFlag))
	}
}

func GetUpdateDenomAlias(chainId string) []DenomAlias {
	if fxtypes.TestnetChainId == chainId {
		// todo deploy testnet contract
		return []DenomAlias{
			{Denom: "weth", Alias: "arbitrum0x0000000000000000000000000000000000000001"},
			{Denom: "usdt", Alias: "arbitrum0x0000000000000000000000000000000000000002"},
			{Denom: "weth", Alias: "optimism0x0000000000000000000000000000000000000003"},
			{Denom: "usdt", Alias: "optimism0x0000000000000000000000000000000000000004"},
		}
	} else if chainId == fxtypes.MainnetChainId {
		return []DenomAlias{
			{Denom: "weth", Alias: "arbitrum0x82aF49447D8a07e3bd95BD0d56f35241523fBab1"},
			{Denom: "usdt", Alias: "arbitrum0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9"},
			{Denom: "weth", Alias: "optimism0x4200000000000000000000000000000000000006"},
			{Denom: "usdt", Alias: "optimism0x94b008aA00579c1307B0EF2c499aD98a8ce58e58"},
		}
	} else {
		panic("invalid chainId:" + chainId)
	}
}

// preUpgradeCmd called by cosmovisor
func preUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pre-upgrade",
		Short: "fxv4 pre-upgrade, called by cosmovisor, before migrations upgrade",
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
			return clientCtx.PrintString("fxv4 pre-upgrade success")
		},
	}
	return cmd
}
