package v8

import (
	"context"
	"errors"
	"strings"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
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
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/pundiai/fx-core/v8/app/keepers"
	"github.com/pundiai/fx-core/v8/app/upgrades/store"
	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	bsctypes "github.com/pundiai/fx-core/v8/x/bsc/types"
	crosschainkeeper "github.com/pundiai/fx-core/v8/x/crosschain/keeper"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20keeper "github.com/pundiai/fx-core/v8/x/erc20/keeper"
	erc20v8 "github.com/pundiai/fx-core/v8/x/erc20/migrations/v8"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	fxevmkeeper "github.com/pundiai/fx-core/v8/x/evm/keeper"
	"github.com/pundiai/fx-core/v8/x/gov/keeper"
	fxgovv8 "github.com/pundiai/fx-core/v8/x/gov/migrations/v8"
	layer2types "github.com/pundiai/fx-core/v8/x/layer2/types"
	fxstakingv8 "github.com/pundiai/fx-core/v8/x/staking/migrations/v8"
)

func CreateUpgradeHandler(cdc codec.Codec, mm *module.Manager, configurator module.Configurator, app *keepers.AppKeepers) upgradetypes.UpgradeHandler {
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

		if err = migrateEvmParams(cacheCtx, app.EvmKeeper); err != nil {
			return fromVM, err
		}

		store.RemoveStoreKeys(cacheCtx, app.GetKey(stakingtypes.StoreKey), fxstakingv8.GetRemovedStoreKeys())

		if err = migrationGovCustomParam(cacheCtx, app.GovKeeper, app.GetKey(govtypes.StoreKey)); err != nil {
			return fromVM, err
		}

		if err = migrateBridgeBalance(cacheCtx, app.BankKeeper, app.AccountKeeper); err != nil {
			return fromVM, err
		}

		if err = migrateERC20TokenToCrosschain(cacheCtx, app.BankKeeper, app.Erc20Keeper); err != nil {
			return fromVM, err
		}

		updateMetadata(cacheCtx, app.BankKeeper)

		store.RemoveStoreKeys(cacheCtx, app.GetKey(erc20types.StoreKey), erc20v8.GetRemovedStoreKeys())

		if err = mintPurseBridgeToken(cacheCtx, app.Erc20Keeper, app.BankKeeper); err != nil {
			return fromVM, err
		}

		acc := app.AccountKeeper.GetModuleAddress(evmtypes.ModuleName)
		moduleAddress := common.BytesToAddress(acc.Bytes())

		if err = deployBridgeFeeContract(
			cacheCtx,
			app.EvmKeeper,
			app.Erc20Keeper,
			app.CrosschainKeepers.EthKeeper,
			moduleAddress,
		); err != nil {
			return fromVM, err
		}

		if err = deployAccessControlContract(cacheCtx, app.EvmKeeper, moduleAddress); err != nil {
			return fromVM, err
		}

		fixBaseOracleStatus(cacheCtx, app.CrosschainKeepers.Layer2Keeper)

		commit()
		cacheCtx.Logger().Info("upgrade complete", "module", "upgrade")
		return toVM, nil
	}
}

func migrateEvmParams(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper) error {
	params := evmKeeper.GetParams(ctx)
	params.HeaderHashNum = evmtypes.DefaultHeaderHashNum
	return evmKeeper.SetParams(ctx, params)
}

func migrationGovCustomParam(ctx sdk.Context, keeper *keeper.Keeper, storeKey *storetypes.KVStoreKey) error {
	// 1. delete fxParams key
	store.RemoveStoreKeys(ctx, storeKey, fxgovv8.GetRemovedStoreKeys())

	// 2. init custom params
	return keeper.InitCustomParams(ctx)
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
		srcDenoms = append(srcDenoms, md.DenomUnits[0].Aliases...)
		if len(srcDenoms) == 0 {
			continue
		}

		for _, srcDenom := range srcDenoms {
			if err := migrateAccountBalance(ctx, bankKeeper, accountKeeper, srcDenom, dstBase); err != nil {
				return err
			}
		}
	}
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
	return err
}

func migrateERC20TokenToCrosschain(ctx sdk.Context, bankKeeper bankkeeper.Keeper, erc20Keeper erc20keeper.Keeper) error {
	balances := bankKeeper.GetAllBalances(ctx, types.NewModuleAddress(erc20types.ModuleName))
	migrateCoins := sdk.NewCoins()
	for _, bal := range balances {
		has, err := erc20Keeper.HasToken(ctx, bal.Denom)
		if err != nil {
			return err
		}
		if !has {
			continue
		}
		migrateCoins = migrateCoins.Add(bal)
	}
	ctx.Logger().Info("migrate erc20 bridge/ibc token to crosschain", "coins", migrateCoins.String())
	return bankKeeper.SendCoinsFromModuleToModule(ctx, erc20types.ModuleName, crosschaintypes.ModuleName, migrateCoins)
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

func mintPurseBridgeToken(ctx sdk.Context, erc20Keeper erc20keeper.Keeper, bankKeeper bankkeeper.Keeper) error {
	pxEscrowPurse, err := getPundixEscrowPurseAmount(ctx)
	if err != nil {
		return err
	}

	ibcToken, err := erc20Keeper.GetIBCToken(ctx, "channel-0", "purse")
	if err != nil {
		return err
	}
	bscPurseToken, err := erc20Keeper.GetBridgeToken(ctx, bsctypes.ModuleName, "purse")
	if err != nil {
		return err
	}
	ibcTokenSupply := bankKeeper.GetSupply(ctx, ibcToken.GetIbcDenom())
	bscPurseAmount := sdk.NewCoin(bscPurseToken.BridgeDenom(), pxEscrowPurse.Sub(ibcTokenSupply.Amount))
	return bankKeeper.MintCoins(ctx, bsctypes.ModuleName, sdk.NewCoins(bscPurseAmount))
}

func deployBridgeFeeContract(
	cacheCtx sdk.Context,
	evmKeeper *fxevmkeeper.Keeper,
	erc20Keeper erc20keeper.Keeper,
	crosschainKeeper crosschainkeeper.Keeper,
	evmModuleAddress common.Address,
) error {
	quoteKeeper := contract.NewBridgeFeeQuoteKeeper(evmKeeper, contract.BridgeFeeAddress)
	oracleKeeper := contract.NewBridgeFeeOracleKeeper(evmKeeper, contract.BridgeFeeOracleAddress)

	chains := crosschaintypes.GetSupportChains()
	bridgeDenoms := make([]contract.BridgeDenoms, len(chains))
	for index, chain := range chains {
		denoms := make([]string, 0)
		bridgeTokens, err := erc20Keeper.GetBridgeTokens(cacheCtx, chain)
		if err != nil {
			return err
		}
		for _, token := range bridgeTokens {
			denoms = append(denoms, token.GetDenom())
		}
		bridgeDenoms[index] = contract.BridgeDenoms{
			ChainName: chain,
			Denoms:    denoms,
		}
	}

	oracles := crosschainKeeper.GetAllOracles(cacheCtx, true)
	if oracles.Len() <= 0 {
		return errors.New("no oracle found")
	}
	return contract.DeployBridgeFeeContract(
		cacheCtx,
		evmKeeper,
		quoteKeeper,
		oracleKeeper,
		bridgeDenoms,
		evmModuleAddress,
		getContractOwner(cacheCtx),
		common.HexToAddress(oracles[0].ExternalAddress),
	)
}

func deployAccessControlContract(
	cacheCtx sdk.Context,
	evmKeeper *fxevmkeeper.Keeper,
	evmModuleAddress common.Address,
) error {
	accessControl := contract.NewAccessControlKeeper(evmKeeper, contract.AccessControlAddress)
	return contract.DeployAccessControlContract(
		cacheCtx,
		evmKeeper,
		accessControl,
		evmModuleAddress,
		getContractOwner(cacheCtx),
	)
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
