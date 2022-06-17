package v2

import (
	"fmt"
	"strings"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/tharsis/ethermint/types"

	migratetypes "github.com/functionx/fx-core/x/migrate/types"

	erc20types "github.com/functionx/fx-core/x/erc20/types"

	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v3/modules/core/03-connection/types"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"

	paramstypesproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	fxtypes "github.com/functionx/fx-core/types"
	erc20keeper "github.com/functionx/fx-core/x/erc20/keeper"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v2
func CreateUpgradeHandler(kvStoreKeyMap map[string]*sdk.KVStoreKey, mm *module.Manager, configurator module.Configurator, bankStoreKey *sdk.KVStoreKey, bankKeeper bankKeeper.Keeper, accountKeeper authkeeper.AccountKeeper, paramsKeeper paramskeeper.Keeper, ibcKeeper *ibckeeper.Keeper, erc20Keeper erc20keeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := ctx.CacheContext()

		// 1. clean testnet data
		clearTestnetKVStores(cacheCtx, kvStoreKeyMap)

		// 2. update FX metadata
		if err := UpdateFXMetadata(cacheCtx, bankKeeper, bankStoreKey); err != nil {
			return nil, err
		}

		// 3. update block params (max_gas:3000000000)
		if err := updateBlockParams(cacheCtx, paramsKeeper); err != nil {
			return nil, err
		}

		// 4. migrate base account to eth account
		MigrateAccountToEth(cacheCtx, accountKeeper)

		// set max expected block time parameter. Replace the default with your expected value
		// https://github.com/cosmos/ibc-go/blob/release/v1.0.x/docs/ibc/proto-docs.md#params-2
		ibcKeeper.ConnectionKeeper.SetParams(cacheCtx, ibcconnectiontypes.DefaultParams())

		// cosmos-sdk 0.42.x from version must be empty
		if len(fromVM) != 0 {
			panic("invalid from version map")
		}

		for n, m := range mm.Modules {
			//NOTE: fromVM empty
			if initGenesis[n] {
				continue
			}
			if v, ok := runMigrates[n]; ok {
				fromVM[n] = v
				continue
			}
			fromVM[n] = m.ConsensusVersion()
		}

		if mm.OrderMigrations == nil {
			mm.OrderMigrations = migrationsOrder(mm.ModuleNames())
		}
		cacheCtx.Logger().Info("start to run module v2 migrations...")
		toVersion, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			return nil, fmt.Errorf("run migrations error %s", err.Error())
		}

		// register coin
		for _, metadata := range fxtypes.GetMetadata() {
			cacheCtx.Logger().Info("add metadata", "coin", metadata.String())
			pair, err := erc20Keeper.RegisterCoin(cacheCtx, metadata)
			if err != nil {
				return nil, fmt.Errorf("register %s error %s", metadata.Base, err.Error())
			}
			cacheCtx.EventManager().EmitEvent(sdk.NewEvent(
				erc20types.EventTypeRegisterCoin,
				sdk.NewAttribute(erc20types.AttributeKeyDenom, pair.Denom),
				sdk.NewAttribute(erc20types.AttributeKeyTokenAddress, pair.Erc20Address),
			))
		}

		//commit upgrade
		commit()

		return toVersion, nil
	}
}

func UpdateFXMetadata(ctx sdk.Context, bankKeeper bankKeeper.Keeper, key *types.KVStoreKey) error {
	md := fxtypes.GetFXMetaData(fxtypes.DefaultDenom)
	if err := md.Validate(); err != nil {
		return fmt.Errorf("invalid %s metadata", fxtypes.DefaultDenom)
	}
	ctx.Logger().Info("update metadata", "metadata", md.String())
	//delete fx
	deleteMetadata(ctx, key, strings.ToLower(fxtypes.DefaultDenom))
	//set FX
	bankKeeper.SetDenomMetaData(ctx, md)
	return nil
}

func MigrateAccountToEth(ctx sdk.Context, ak authkeeper.AccountKeeper) {
	ctx.Logger().Info("update v2 migrate account to eth")
	// migrate base account to eth account
	ak.IterateAccounts(ctx, func(account authtypes.AccountI) (stop bool) {
		if _, ok := account.(ethermint.EthAccountI); ok {
			return false
		}
		baseAccount, ok := account.(*authtypes.BaseAccount)
		if !ok {
			ctx.Logger().Info("migrate account", "address", account.GetAddress(), "ignore type", fmt.Sprintf("%T", account))
			return false
		}
		ethAccount := &ethermint.EthAccount{
			BaseAccount: baseAccount,
			CodeHash:    common.BytesToHash(emptyCodeHash).String(),
		}
		ak.SetAccount(ctx, ethAccount)
		ctx.Logger().Info("migrate account to eth", "address", account.GetAddress())
		return false
	})
}

func updateBlockParams(ctx sdk.Context, pk paramskeeper.Keeper) error {
	baseappSubspace, found := pk.GetSubspace(baseapp.Paramspace)
	if !found {
		return sdkerrors.Wrap(paramstypesproposal.ErrUnknownSubspace, baseapp.Paramspace)
	}
	var bp abci.BlockParams
	baseappSubspace.Get(ctx, baseapp.ParamStoreKeyBlockParams, &bp)

	bp.MaxGas = blockParamsMaxGas
	baseappSubspace.Set(ctx, baseapp.ParamStoreKeyBlockParams, bp)
	return nil
}

func deleteMetadata(ctx sdk.Context, key *types.KVStoreKey, base ...string) {
	store := ctx.KVStore(key)
	for _, b := range base {
		store.Delete(banktypes.DenomMetadataKey(b))
	}
}

func migrationsOrder(modules []string) []string {
	modules = module.DefaultMigrationsOrder(modules)
	orders := make([]string, 0, len(modules))
	for _, name := range modules {
		if name == feemarkettypes.ModuleName ||
			name == evmtypes.ModuleName ||
			name == erc20types.ModuleName ||
			name == migratetypes.ModuleName {
			continue
		}
		orders = append(orders, name)
	}
	orders = append(orders, []string{
		feemarkettypes.ModuleName, evmtypes.ModuleName,
		erc20types.ModuleName, migratetypes.ModuleName,
	}...)
	return orders
}

func clearTestnetKVStores(ctx sdk.Context, keys map[string]*types.KVStoreKey) {
	logger := ctx.Logger()
	if fxtypes.NetworkTestnet() != fxtypes.Network() {
		return
	}
	logger.Info("clear kv store", "network", fxtypes.Network())
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
		logger.Info("clear kv store done", "storesName", storeName, "consumeMs", time.Now().UnixMilli()-startTime.UnixMilli())
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
