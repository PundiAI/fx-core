package v4

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	avalanchetypes "github.com/functionx/fx-core/v3/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	gravitytypes "github.com/functionx/fx-core/v3/x/gravity/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"

	"github.com/functionx/fx-core/v3/app/keepers"
	fxcfg "github.com/functionx/fx-core/v3/server/config"
)

func createUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	app *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := ctx.CacheContext()

		// 1. Set param key table for params module migration
		SetParamKeyTableMigration(app.ParamsKeeper)

		ctx.Logger().Info("start to run v4 migrations...", "module", "upgrade")
		toVM, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			panic(fmt.Sprintf("run migrations: %s", err.Error()))
		}
		commit()
		ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
		return toVM, nil
	}
}

func SetParamKeyTableMigration(keeper paramskeeper.Keeper) {
	for _, subspace := range keeper.GetSubspaces() {
		subspace := subspace
		var keyTable paramstypes.KeyTable
		switch subspace.Name() {
		case avalanchetypes.ModuleName:
			keyTable = crosschaintypes.ParamKeyTable() //nolint:staticcheck
		case bsctypes.ModuleName:
			keyTable = crosschaintypes.ParamKeyTable() //nolint:staticcheck
		case ethtypes.ModuleName:
			keyTable = crosschaintypes.ParamKeyTable() //nolint:staticcheck
		case polygontypes.ModuleName:
			keyTable = crosschaintypes.ParamKeyTable() //nolint:staticcheck
		case trontypes.ModuleName:
			keyTable = crosschaintypes.ParamKeyTable() //nolint:staticcheck
		case erc20types.ModuleName:
			keyTable = erc20types.ParamKeyTable()
		}
		if !subspace.HasKeyTable() && !strings.EqualFold(subspace.Name(), gravitytypes.ModuleName) {
			subspace.WithKeyTable(keyTable)
		}
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
