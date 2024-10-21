package v8

import (
	"context"
	"encoding/hex"
	"strings"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"

	"github.com/functionx/fx-core/v8/app/keepers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	crosschainkeeper "github.com/functionx/fx-core/v8/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	erc20v8 "github.com/functionx/fx-core/v8/x/erc20/migrations/v8"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
	"github.com/functionx/fx-core/v8/x/gov/keeper"
	fxgovv8 "github.com/functionx/fx-core/v8/x/gov/migrations/v8"
	fxstakingv8 "github.com/functionx/fx-core/v8/x/staking/migrations/v8"
)

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator, app *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := sdk.UnwrapSDKContext(ctx).CacheContext()

		if err := migrateCrosschainModuleAccount(cacheCtx, app.AccountKeeper); err != nil {
			return fromVM, err
		}

		cacheCtx.Logger().Info("start to run migrations...", "module", "upgrade", "plan", plan.Name)
		toVM, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			return fromVM, err
		}

		removeStoreKeys(cacheCtx, app.GetKey(stakingtypes.StoreKey), fxstakingv8.GetRemovedStoreKeys())

		if err = migrationGovCustomParam(cacheCtx, app.GovKeeper, app.GetKey(govtypes.StoreKey)); err != nil {
			return fromVM, err
		}

		if err = updateArbitrumParams(cacheCtx, app.ArbitrumKeeper); err != nil {
			return fromVM, err
		}

		if err = migrateBridgeBalance(cacheCtx, app.BankKeeper, app.AccountKeeper); err != nil {
			return fromVM, err
		}

		updateMetadata(cacheCtx, app.BankKeeper)

		removeStoreKeys(cacheCtx, app.GetKey(erc20types.StoreKey), erc20v8.GetRemovedStoreKeys())

		commit()
		cacheCtx.Logger().Info("upgrade complete", "module", "upgrade")
		return toVM, nil
	}
}

func updateArbitrumParams(ctx sdk.Context, keeper crosschainkeeper.Keeper) error {
	params := keeper.GetParams(ctx)
	params.AverageExternalBlockTime = 250
	return keeper.SetParams(ctx, &params)
}

func migrationGovCustomParam(ctx sdk.Context, keeper *keeper.Keeper, storeKey *storetypes.KVStoreKey) error {
	// 1. delete fxParams key
	removeStoreKeys(ctx, storeKey, fxgovv8.GetRemovedStoreKeys())

	// 2. init custom params
	return keeper.InitCustomParams(ctx)
}

func removeStoreKeys(ctx sdk.Context, storeKey *storetypes.KVStoreKey, keys [][]byte) {
	store := ctx.KVStore(storeKey)

	deleteFn := func(key []byte) {
		iterator := storetypes.KVStorePrefixIterator(store, key)
		defer iterator.Close()
		for ; iterator.Valid(); iterator.Next() {
			store.Delete(iterator.Key())
			ctx.Logger().Info("remove store key", "kvStore", storeKey.Name(),
				"prefix", hex.EncodeToString(key), "key", string(iterator.Key()))
		}
	}

	for _, key := range keys {
		deleteFn(key)
	}
}

func migrateCrosschainModuleAccount(ctx sdk.Context, ak authkeeper.AccountKeeper) error {
	addr, perms := ak.GetModuleAddressAndPermissions(crosschaintypes.ModuleName)
	if addr == nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("crosschain module empty permissions")
	}
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("crosschain account not exist")
	}
	baseAcc, ok := acc.(*types.BaseAccount)
	if !ok {
		return sdkerrors.ErrInvalidAddress.Wrapf("crosschain account not base account")
	}
	macc := types.NewModuleAccount(baseAcc, crosschaintypes.ModuleName, perms...)
	ak.SetModuleAccount(ctx, macc)
	return nil
}

func migrateBridgeBalance(ctx sdk.Context, bankKeeper bankkeeper.Keeper, accountKeeper authkeeper.AccountKeeper) error {
	mds := bankKeeper.GetAllDenomMetaData(ctx)
	for _, md := range mds {
		if md.Base == fxtypes.DefaultDenom || (len(md.DenomUnits) == 0 || len(md.DenomUnits[0].Aliases) == 0) && md.Symbol != "PUNDIX" {
			continue
		}
		dstBase := strings.ToLower(md.Symbol)
		srcDenoms := make([]string, 0, len(md.DenomUnits[0].Aliases)+1)
		if md.Base != dstBase {
			// pundix, purse
			srcDenoms = append(srcDenoms, md.Base)
		}
		// bridge token, exclude ibc
		bridgeTokens := make([]string, 0, len(md.DenomUnits[0].Aliases))
		for _, alias := range md.DenomUnits[0].Aliases {
			if strings.HasPrefix(alias, ibctransfertypes.DenomPrefix+"/") {
				continue
			}
			bridgeTokens = append(bridgeTokens, alias)
		}
		srcDenoms = append(srcDenoms, bridgeTokens...)
		if len(srcDenoms) == 0 {
			continue
		}

		for _, srcDenom := range srcDenoms {
			if err := migrateAccountBalance(ctx, bankKeeper, accountKeeper, srcDenom, dstBase); err != nil {
				return err
			}
		}
	}
	// todo migrate bridge token and ibc token balance to crosschain module
	return nil
}

func migrateAccountBalance(ctx sdk.Context, bankKeeper bankkeeper.Keeper, accountKeeper authkeeper.AccountKeeper, srcBase, dstBase string) error {
	var err error
	bankKeeper.IterateAllBalances(ctx, func(address sdk.AccAddress, coin sdk.Coin) (stop bool) {
		if coin.Denom != srcBase {
			return false
		}

		account := accountKeeper.GetAccount(ctx, address)
		if _, ok := account.(sdk.ModuleAccountI); ok {
			return false
		}

		ctx.Logger().Info("migrate coin", "address", address.String(), "src-denom", srcBase, "dst-denom", dstBase, "amount", coin.Amount.String())
		if err = bankKeeper.SendCoinsFromAccountToModule(ctx, address, erc20types.ModuleName, sdk.NewCoins(coin)); err != nil {
			return true
		}
		coin.Denom = dstBase
		if err = bankKeeper.MintCoins(ctx, crosschaintypes.ModuleName, sdk.NewCoins(coin)); err != nil {
			return true
		}
		if err = bankKeeper.SendCoinsFromModuleToAccount(ctx, crosschaintypes.ModuleName, address, sdk.NewCoins(coin)); err != nil {
			return true
		}

		return false
	})
	return nil
}

func updateMetadata(ctx sdk.Context, bankKeeper bankkeeper.Keeper) {
	mds := bankKeeper.GetAllDenomMetaData(ctx)
	for _, md := range mds {
		if md.Base == fxtypes.DefaultDenom || (len(md.DenomUnits) == 0 || len(md.DenomUnits[0].Aliases) == 0) && md.Symbol != "PUNDIX" {
			continue
		}
		// remove alias
		md.DenomUnits[0].Aliases = []string{}

		newBase := strings.ToLower(md.Symbol)
		// update pundix/purse base denom
		if md.Base != newBase && !strings.Contains(md.Base, newBase) && !strings.HasPrefix(md.Display, ibctransfertypes.ModuleName+"/"+ibcchanneltypes.ChannelPrefix) {
			md.Base = newBase
			md.Display = newBase
			md.DenomUnits[0].Denom = newBase
		}

		bankKeeper.SetDenomMetaData(ctx, md)
	}
}
