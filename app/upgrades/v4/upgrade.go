package v4

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v5/app/keepers"
	fxtypes "github.com/functionx/fx-core/v5/types"
	avalanchetypes "github.com/functionx/fx-core/v5/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v5/x/bsc/types"
	crosschainkeeper "github.com/functionx/fx-core/v5/x/crosschain/keeper"
	"github.com/functionx/fx-core/v5/x/crosschain/types"
	erc20keeper "github.com/functionx/fx-core/v5/x/erc20/keeper"
	ethtypes "github.com/functionx/fx-core/v5/x/eth/types"
	evmkeeper "github.com/functionx/fx-core/v5/x/evm/keeper"
	"github.com/functionx/fx-core/v5/x/gov/keeper"
	polygontypes "github.com/functionx/fx-core/v5/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v5/x/tron/types"
)

func CreateUpgradeHandler(
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

		// 3. update Logic code
		UpdateLogicCode(cacheCtx, app.EvmKeeper)

		ctx.Logger().Info("start to run v4 migrations...", "module", "upgrade")
		toVM, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			return fromVM, err
		}

		// 4. update arbitrum and optimism denom alias, after bank module migration, because bank module migrates to fixing the bank denom bug
		// discovered in https://github.com/cosmos/cosmos-sdk/pull/13821
		UpdateDenomAliases(cacheCtx, app.Erc20Keeper)

		// 5. reset cross chain module oracle delegate, bind oracle delegate starting info
		err = ResetCrossChainModuleOracleDelegate(cacheCtx, app.CrossChainKeepers, app.StakingKeeper, app.DistrKeeper)
		if err != nil {
			return fromVM, err
		}

		// 6. remove bsc oracles
		RemoveBscOracle(cacheCtx, app.BscKeeper)

		commit()
		ctx.Logger().Info("Upgrade complete")
		return toVM, nil
	}
}

func ResetCrossChainModuleOracleDelegate(ctx sdk.Context, crossChainKeepers keepers.CrossChainKeepers, stakingKeeper types.StakingKeeper, distributionKeeper types.DistributionKeeper) error {
	needHandlerModules := []string{ethtypes.ModuleName, bsctypes.ModuleName, polygontypes.ModuleName, trontypes.ModuleName, avalanchetypes.ModuleName}
	type crossChainKeeper interface {
		GetAllOracles(ctx sdk.Context, isOnline bool) (oracles types.Oracles)
	}
	moduleHandler := map[string]crossChainKeeper{
		ethtypes.ModuleName:       crossChainKeepers.EthKeeper,
		bsctypes.ModuleName:       crossChainKeepers.BscKeeper,
		trontypes.ModuleName:      crossChainKeepers.TronKeeper,
		polygontypes.ModuleName:   crossChainKeepers.PolygonKeeper,
		avalanchetypes.ModuleName: crossChainKeepers.AvalancheKeeper,
	}
	for _, handlerModule := range needHandlerModules {
		handlerKeeper, ok := moduleHandler[handlerModule]
		if !ok {
			continue
		}
		oracles := handlerKeeper.GetAllOracles(ctx, false)
		if len(oracles) <= 0 {
			continue
		}

		for _, oracle := range oracles {
			if oracle.DelegateAmount.IsZero() {
				continue
			}

			delegateAddress := oracle.GetDelegateAddress(handlerModule)
			startingInfo := distributionKeeper.GetDelegatorStartingInfo(ctx, oracle.GetValidator(), delegateAddress)
			if startingInfo.Height > 0 {
				continue
			}
			err := stakingKeeper.BeforeDelegationCreated(ctx, delegateAddress, oracle.GetValidator())
			if err != nil {
				return err
			}
			err = stakingKeeper.AfterDelegationModified(ctx, delegateAddress, oracle.GetValidator())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func RemoveBscOracle(ctx sdk.Context, bscKeeper crosschainkeeper.Keeper) {
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

func UpdateLogicCode(ctx sdk.Context, evmKeeper *evmkeeper.Keeper) {
	UpdateFIP20LogicCode(ctx, evmKeeper)
	UpdateWFXLogicCode(ctx, evmKeeper)
}

func UpdateFIP20LogicCode(ctx sdk.Context, k *evmkeeper.Keeper) {
	fip20 := fxtypes.GetFIP20()
	if err := k.UpdateContractCode(ctx, fip20.Address, fip20.Code); err != nil {
		panic(fmt.Sprintf("update fip logic code error: %s", err.Error()))
	}
	ctx.Logger().Info("update FIP20 contract", "module", "upgrade", "codeHash", fip20.CodeHash())
}

func UpdateWFXLogicCode(ctx sdk.Context, k *evmkeeper.Keeper) {
	wfx := fxtypes.GetWFX()
	if err := k.UpdateContractCode(ctx, wfx.Address, wfx.Code); err != nil {
		panic(fmt.Sprintf("update wfx logic code error: %s", err.Error()))
	}
	ctx.Logger().Info("update WFX contract", "module", "upgrade", "codeHash", wfx.CodeHash())
}

func InitGovFXParams(ctx sdk.Context, keeper keeper.Keeper) {
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
