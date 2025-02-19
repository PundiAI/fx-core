package v8

import (
	"context"

	"cosmossdk.io/x/feegrant"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/pundiai/fx-core/v8/app/keepers"
	"github.com/pundiai/fx-core/v8/app/upgrades/store"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschainkeeper "github.com/pundiai/fx-core/v8/x/crosschain/keeper"
	erc20v8 "github.com/pundiai/fx-core/v8/x/erc20/migrations/v8"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	layer2types "github.com/pundiai/fx-core/v8/x/layer2/types"
	fxstakingv8 "github.com/pundiai/fx-core/v8/x/staking/migrations/v8"
)

func CreateUpgradeHandler(codec codec.Codec, mm *module.Manager, configurator module.Configurator, app *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := sdk.UnwrapSDKContext(ctx).CacheContext()

		var err error
		var toVM module.VersionMap
		if cacheCtx.ChainID() == fxtypes.TestnetChainId {
			return fromVM, nil
		} else {
			toVM, err = upgradeMainnet(cacheCtx, codec, mm, configurator, app, fromVM, plan)
			if err != nil {
				return fromVM, err
			}
		}
		commit()
		cacheCtx.Logger().Info("upgrade complete", "module", "upgrade")
		return toVM, nil
	}
}

func upgradeMainnet(
	ctx sdk.Context,
	codec codec.Codec,
	mm *module.Manager,
	configurator module.Configurator,
	app *keepers.AppKeepers,
	fromVM module.VersionMap,
	plan upgradetypes.Plan,
) (module.VersionMap, error) {
	if err := migrateCrosschainModuleAccount(ctx, app.AccountKeeper); err != nil {
		return fromVM, err
	}
	if err := migrateCrosschainParams(ctx, app.CrosschainKeepers); err != nil {
		return fromVM, err
	}

	ctx.Logger().Info("start to run migrations...", "module", "upgrade", "plan", plan.Name)
	toVM, err := mm.RunMigrations(ctx, configurator, fromVM)
	if err != nil {
		return fromVM, err
	}

	store.RemoveStoreKeys(ctx, app.GetKey(stakingtypes.StoreKey), fxstakingv8.GetRemovedStoreKeys())
	store.RemoveStoreKeys(ctx, app.GetKey(erc20types.StoreKey), erc20v8.GetRemovedStoreKeys())
	fixBaseOracleStatus(ctx, app.CrosschainKeepers.Layer2Keeper)

	if err = migrateGovCustomParam(ctx, app.GovKeeper, app.GetKey(govtypes.StoreKey)); err != nil {
		return fromVM, err
	}
	if err = migrateGovDefaultParams(ctx, app.GovKeeper); err != nil {
		return fromVM, err
	}
	if err = removeGovPendingProposal(ctx, app.GovKeeper); err != nil {
		return fromVM, err
	}
	if err = migrateBridgeToken(ctx, app.EvmKeeper, app.Erc20Keeper, app.BankKeeper, app.AccountKeeper); err != nil {
		return fromVM, err
	}
	if err = updateMetadata(ctx, app.BankKeeper); err != nil {
		return fromVM, err
	}
	if err = updatePundiAI(ctx, app); err != nil {
		return fromVM, err
	}
	if err = updateContract(ctx, app); err != nil {
		return fromVM, err
	}
	if err = migrateModulesData(ctx, codec, app); err != nil {
		return fromVM, err
	}

	initBridgeAccount(ctx, app.AccountKeeper)
	return toVM, nil
}

func migrateModulesData(ctx sdk.Context, codec codec.Codec, app *keepers.AppKeepers) error {
	migrateWFXToWPUNDIAI(ctx, app.EvmKeeper)
	migrateTransferTokenInEscrow(ctx, app.IBCTransferKeeper)
	migrateOracleDelegateAmount(ctx, app.CrosschainKeepers)

	if err := migrateFeemarketGasPrice(ctx, app.FeeMarketKeeper); err != nil {
		return err
	}
	if err := migrateMetadataDisplay(ctx, app.BankKeeper); err != nil {
		return err
	}
	if err := migrateMetadataFXToPundiAI(ctx, app.BankKeeper); err != nil {
		return err
	}
	if err := migrateStakingModule(ctx, app.StakingKeeper.Keeper); err != nil {
		return err
	}
	if err := migrateEvmParams(ctx, app.EvmKeeper); err != nil {
		return err
	}
	if err := migrateMintParams(ctx, app.MintKeeper); err != nil {
		return err
	}
	if err := MigrateFeegrant(ctx, codec, runtime.NewKVStoreService(app.GetKey(feegrant.StoreKey)), app.AccountKeeper); err != nil {
		return err
	}
	if err := migrateDistribution(ctx, app.StakingKeeper, app.DistrKeeper); err != nil {
		return err
	}
	if err := migrateBankModule(ctx, app.BankKeeper); err != nil {
		return err
	}
	return migrateCrisisModule(ctx, app.CrisisKeeper)
}

func fixBaseOracleStatus(ctx sdk.Context, crosschainKeeper crosschainkeeper.Keeper) {
	if crosschainKeeper.ModuleName() != layer2types.ModuleName {
		return
	}
	oracles := crosschainKeeper.GetAllOracles(ctx, false)
	for _, oracle := range oracles {
		oracle.Online = true
		oracle.SlashTimes = 0
		oracle.StartHeight = ctx.BlockHeight()
		crosschainKeeper.SetOracle(ctx, oracle)
	}
	crosschainKeeper.SetLastTotalPower(ctx)
}
