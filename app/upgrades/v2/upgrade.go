package v2

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v3/modules/core/03-connection/types"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	tmcfg "github.com/tendermint/tendermint/config"

	"github.com/functionx/fx-core/v3/app/keepers"
	fxcfg "github.com/functionx/fx-core/v3/server/config"
	fxtypes "github.com/functionx/fx-core/v3/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	crosschainv2 "github.com/functionx/fx-core/v3/x/crosschain/legacy/v2"
	erc20keeper "github.com/functionx/fx-core/v3/x/erc20/keeper"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	migratetypes "github.com/functionx/fx-core/v3/x/migrate/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

// createUpgradeHandler creates an SDK upgrade handler for v2
func createUpgradeHandler(mm *module.Manager, configurator module.Configurator, keepers *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// cache context
		cacheCtx, commit := ctx.CacheContext()

		kvStoreKeyMap := keepers.GetKVStoreKey()
		// 1. clear testnet module data
		clearTestnetModule(cacheCtx, kvStoreKeyMap)

		// 2. update FX metadata
		updateFXMetadata(cacheCtx, keepers.BankKeeper, kvStoreKeyMap[banktypes.StoreKey])

		// 3. update block params (max_gas:3000000000)
		updateBlockParams(cacheCtx, keepers.ParamsKeeper)

		// 4. migrate ibc cosmos-sdk/x/ibc -> ibc-go v3.1.0
		ibcMigrate(cacheCtx, keepers.IBCKeeper, keepers.IBCTransferKeeper)

		// 5. run migrations
		toVersion := runMigrations(cacheCtx, kvStoreKeyMap[paramstypes.StoreKey], fromVM, mm, configurator)

		// 6. clear metadata except FX
		clearTestnetDenom(cacheCtx, kvStoreKeyMap[banktypes.StoreKey])

		// 7. register coin
		registerCoin(cacheCtx, keepers.Erc20Keeper)

		// commit upgrade
		commit()
		ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())

		return toVersion, nil
	}
}

func ibcMigrate(ctx sdk.Context, ibcKeeper *ibckeeper.Keeper, transferKeeper ibctransferkeeper.Keeper) {
	// set max expected block time parameter. Replace the default with your expected value
	// https://github.com/cosmos/ibc-go/blob/release/v1.0.x/docs/ibc/proto-docs.md#params-2
	ibcKeeper.ConnectionKeeper.SetParams(ctx, ibcconnectiontypes.DefaultParams())

	// list of traces that must replace the old traces in store
	// https://github.com/cosmos/ibc-go/blob/v3.1.0/docs/migrations/support-denoms-with-slashes.md
	var newTraces []ibctransfertypes.DenomTrace
	transferKeeper.IterateDenomTraces(ctx,
		func(dt ibctransfertypes.DenomTrace) bool {
			// check if the new way of splitting FullDenom
			// into Trace and BaseDenom passes validation and
			// is the same as the current DenomTrace.
			// If it isn't then store the new DenomTrace in the list of new traces.
			newTrace := ibctransfertypes.ParseDenomTrace(dt.GetFullDenomPath())
			if err := newTrace.Validate(); err == nil && !reflect.DeepEqual(newTrace, dt) {
				newTraces = append(newTraces, newTrace)
			}

			return false
		})

	// replace the outdated traces with the new trace information
	for _, nt := range newTraces {
		transferKeeper.SetDenomTrace(ctx, nt)
	}
}

func updateFXMetadata(ctx sdk.Context, bankKeeper bankKeeper.Keeper, bankKey *sdk.KVStoreKey) {
	metaData := fxtypes.GetFXMetaData(fxtypes.DefaultDenom)
	if err := metaData.Validate(); err != nil {
		panic(fmt.Sprintf("invalid %s metadata", fxtypes.DefaultDenom))
	}
	logger := ctx.Logger()
	logger.Info("update FX metadata", "metadata", metaData.String())
	// delete fx
	fxDenom := strings.ToLower(fxtypes.DefaultDenom)
	denomMetaDataStore := prefix.NewStore(ctx.KVStore(bankKey), banktypes.DenomMetadataKey(fxDenom))
	denomMetaDataStore.Delete([]byte(fxDenom))
	// set FX
	bankKeeper.SetDenomMetaData(ctx, metaData)
}

func updateBlockParams(ctx sdk.Context, pk paramskeeper.Keeper) {
	logger := ctx.Logger()
	logger.Info("update block params", "chainId", ctx.ChainID())
	baseappSubspace, found := pk.GetSubspace(baseapp.Paramspace)
	if !found {
		panic(fmt.Sprintf("unknown subspace: %s", baseapp.Paramspace))
	}
	var bp abci.BlockParams
	baseappSubspace.Get(ctx, baseapp.ParamStoreKeyBlockParams, &bp)
	logger.Info("update block params", "before update", bp.String())
	bp.MaxGas = 30_000_000
	baseappSubspace.Set(ctx, baseapp.ParamStoreKeyBlockParams, bp)
	logger.Info("update block params", "after update", bp.String())
}

func migrationsOrder(modules []string) []string {
	modules = module.DefaultMigrationsOrder(modules)
	orders := make([]string, 0, len(modules))
	for _, name := range modules {
		if name == bsctypes.ModuleName || name == polygontypes.ModuleName || name == trontypes.ModuleName ||
			name == feemarkettypes.ModuleName || name == evmtypes.ModuleName ||
			name == erc20types.ModuleName || name == migratetypes.ModuleName {
			continue
		}
		orders = append(orders, name)
	}
	orders = append(orders, []string{
		bsctypes.ModuleName, polygontypes.ModuleName, trontypes.ModuleName,
		feemarkettypes.ModuleName, evmtypes.ModuleName,
		erc20types.ModuleName, migratetypes.ModuleName,
	}...)
	return orders
}

func runMigrations(ctx sdk.Context, paramsKey *sdk.KVStoreKey, fromVersion module.VersionMap,
	mm *module.Manager, configurator module.Configurator,
) module.VersionMap {
	if len(fromVersion) != 0 {
		panic("invalid from version map")
	}

	for n, m := range mm.Modules {
		// NOTE: fromVM empty
		if initGenesis[n] {
			continue
		}
		// if module genesis init, continue
		if needInitGenesis(ctx, n, paramsKey) {
			continue
		}
		if v, ok := runMigrates[n]; ok {
			// migrate module
			fromVersion[n] = v
			continue
		}
		fromVersion[n] = m.ConsensusVersion()
	}

	if mm.OrderMigrations == nil {
		mm.OrderMigrations = migrationsOrder(mm.ModuleNames())
	}
	ctx.Logger().Info("start to run module v2 migrations...")
	toVersion, err := mm.RunMigrations(ctx, configurator, fromVersion)
	if err != nil {
		panic(fmt.Sprintf("run migrations: %s", err.Error()))
	}
	return toVersion
}

func clearTestnetDenom(ctx sdk.Context, bankKey *types.KVStoreKey) {
	if fxtypes.TestnetChainId != ctx.ChainID() {
		return
	}
	logger := ctx.Logger()
	logger.Info("clear testnet metadata", "chainId", ctx.ChainID())
	for _, md := range getMetadata(ctx.ChainID()) {
		// remove denom except FX
		if md.Base == fxtypes.DefaultDenom {
			continue
		}
		logger.Info("clear testnet metadata", "metadata", md.String())
		denomMetaDataStore := prefix.NewStore(ctx.KVStore(bankKey), banktypes.DenomMetadataKey(md.Base))
		denomMetaDataStore.Delete([]byte(md.Base))
	}
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

func clearTestnetModule(ctx sdk.Context, keys map[string]*types.KVStoreKey) {
	logger := ctx.Logger()
	if fxtypes.TestnetChainId != ctx.ChainID() {
		return
	}
	logger.Info("clear kv store", "chainId", ctx.ChainID())
	cleanModules := []string{feemarkettypes.StoreKey, evmtypes.StoreKey, erc20types.StoreKey, migratetypes.StoreKey}
	multiStore := ctx.MultiStore()
	for _, storeName := range cleanModules {
		logger.Info("clear kv store", "storesName", storeName)
		startTime := time.Now()
		storeKey, ok := keys[storeName]
		if !ok {
			panic(fmt.Sprintf("%s store not found", storeName))
		}
		kvStore := multiStore.GetKVStore(storeKey)
		if err := deleteKVStore(kvStore); err != nil {
			panic(fmt.Sprintf("failed to delete store %s: %s", storeName, err.Error()))
		}
		logger.Info("clear kv store done", "storesName", storeName, "consumeMs", time.Now().UnixNano()-startTime.UnixNano())
	}
}

func deleteKVStore(kv types.KVStore) error {
	// Note that we cannot write while iterating, so load all keys here, delete below
	var keys [][]byte
	itr := kv.Iterator(nil, nil)
	defer itr.Close()

	for itr.Valid() {
		keys = append(keys, itr.Key())
		itr.Next()
	}

	for _, k := range keys {
		kv.Delete(k)
	}
	return nil
}

// needInitGenesis check module initialized
func needInitGenesis(ctx sdk.Context, module string, paramsKey *sdk.KVStoreKey) bool {
	// crosschain module
	if crossChainModule[module] {
		if !crosschainv2.CheckInitialize(ctx, module, paramsKey) {
			return true
		}
	}
	return false
}

// preUpgradeCmd called by cosmovisor
func preUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pre-upgrade",
		Short: "fxv2 pre-upgrade, called by cosmovisor, before migrations upgrade",
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
			return clientCtx.PrintString("fxv2 pre-upgrade success")
		},
	}
	return cmd
}
