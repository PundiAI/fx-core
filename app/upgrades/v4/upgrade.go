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

	"github.com/functionx/fx-core/v4/app/keepers"
	fxcfg "github.com/functionx/fx-core/v4/server/config"
	fxtypes "github.com/functionx/fx-core/v4/types"
	crosschainkeeper "github.com/functionx/fx-core/v4/x/crosschain/keeper"
	erc20keeper "github.com/functionx/fx-core/v4/x/erc20/keeper"
	evmkeeper "github.com/functionx/fx-core/v4/x/evm/keeper"
	"github.com/functionx/fx-core/v4/x/gov/keeper"
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

		// 3. update Logoic code
		updateLogicCode(cacheCtx, app.EvmKeeper)

		// 4. remove bsc oracles
		removeBscOracle(cacheCtx, app.BscKeeper)

		ctx.Logger().Info("start to run v4 migrations...", "module", "upgrade")
		toVM, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			return fromVM, err
		}

		// update arbitrum and optimism denom alias, after bank module migration, because bank module migrates to fixing the bank denom bug
		// discovered in https://github.com/cosmos/cosmos-sdk/pull/13821
		UpdateDenomAliases(cacheCtx, app.Erc20Keeper)

		commit()
		ctx.Logger().Info("Upgrade complete")
		return toVM, nil
	}
}

func removeBscOracle(ctx sdk.Context, bscKeeper crosschainkeeper.Keeper) {
	bscRemoveOracles := GetBscRemoveOracles(ctx.ChainID())
	if len(bscRemoveOracles) <= 0 {
		return
	}

	proposalOracle, found := bscKeeper.GetProposalOracle(ctx)
	oracles := proposalOracle.Oracles
	if !found || len(oracles) <= 0 {
		return
	}

	removeOracleMap := make(map[string]bool, len(bscRemoveOracles))
	for _, oracle := range bscRemoveOracles {
		removeOracleMap[oracle] = true
	}

	newOracle := []string{}
	for _, oracle := range oracles {
		if _, ok := removeOracleMap[oracle]; ok {
			continue
		}
		newOracle = append(newOracle, oracle)
	}

	if len(newOracle) == len(oracles) {
		return
	}
	err := bscKeeper.UpdateChainOracles(ctx, newOracle)
	if err != nil && ctx.ChainID() == fxtypes.TestnetChainId {
		panic(err)
	}
}

func updateLogicCode(ctx sdk.Context, evmKeeper *evmkeeper.Keeper) {
	updateFIP20LogicCode(ctx, evmKeeper)
	updateWFXLogicCode(ctx, evmKeeper)
}

func updateFIP20LogicCode(ctx sdk.Context, k *evmkeeper.Keeper) {
	fip20 := fxtypes.GetFIP20()
	if err := k.UpdateContractCode(ctx, fip20.Address, fip20.Code); err != nil {
		panic(fmt.Sprintf("update fip logic code error: %s", err.Error()))
	}
	ctx.Logger().Info("update FIP20 contract", "module", "upgrade", "codeHash", fip20.CodeHash())
}

func updateWFXLogicCode(ctx sdk.Context, k *evmkeeper.Keeper) {
	wfx := fxtypes.GetWFX()
	if err := k.UpdateContractCode(ctx, wfx.Address, wfx.Code); err != nil {
		panic(fmt.Sprintf("update wfx logic code error: %s", err.Error()))
	}
	ctx.Logger().Info("update WFX contract", "module", "upgrade", "codeHash", wfx.CodeHash())
}

func InitGovFXParams(ctx sdk.Context, keeper keeper.Keeper) {
	// init fx gov params
	if err := keeper.InitFxGovParams(ctx); err != nil {
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
