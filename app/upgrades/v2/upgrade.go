package v2

import (
	"fmt"
	"strings"

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

	fxtypes "github.com/functionx/fx-core/types"
	erc20keeper "github.com/functionx/fx-core/x/erc20/keeper"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v2
func CreateUpgradeHandler(
	mm *module.Manager, configurator module.Configurator,
	bankStoreKey *sdk.KVStoreKey, bankKeeper bankKeeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	ibcKeeper *ibckeeper.Keeper, erc20Keeper erc20keeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := ctx.CacheContext()

		// update FX metadata
		if err := UpdateFXMetadata(cacheCtx, bankKeeper, bankStoreKey); err != nil {
			return nil, err
		}

		// migrate base account to eth account
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
		if name == erc20types.ModuleName ||
			name == migratetypes.ModuleName {
			continue
		}
		orders = append(orders, name)
	}
	orders = append(orders, []string{
		erc20types.ModuleName,
		migratetypes.ModuleName,
	}...)
	return orders
}
