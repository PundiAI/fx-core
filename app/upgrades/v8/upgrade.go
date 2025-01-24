package v8

import (
	"context"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/x/feegrant"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
	fxevmkeeper "github.com/pundiai/fx-core/v8/x/evm/keeper"
	layer2types "github.com/pundiai/fx-core/v8/x/layer2/types"
	fxstakingv8 "github.com/pundiai/fx-core/v8/x/staking/migrations/v8"
)

func CreateUpgradeHandler(codec codec.Codec, mm *module.Manager, configurator module.Configurator, app *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := sdk.UnwrapSDKContext(ctx).CacheContext()

		var err error
		var toVM module.VersionMap
		if cacheCtx.ChainID() == fxtypes.TestnetChainId {
			if err = upgradeTestnet(cacheCtx, codec, app); err != nil {
				return fromVM, err
			}
			toVM = fromVM
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

func upgradeTestnet(ctx sdk.Context, codec codec.Codec, app *keepers.AppKeepers) error {
	fixBaseOracleStatus(ctx, app.CrosschainKeepers.Layer2Keeper)
	updateWPUNDIAILogicCode(ctx, app.EvmKeeper)
	updateERC20LogicCode(ctx, app.EvmKeeper)

	if err := fixPundixCoin(ctx, app.EvmKeeper, app.Erc20Keeper, app.BankKeeper); err != nil {
		return err
	}
	if err := fixPurseCoin(ctx, app.EvmKeeper, app.Erc20Keeper, app.BankKeeper); err != nil {
		return err
	}
	if err := fixTestnetTokenAmount(ctx, app.BankKeeper, app.EvmKeeper, app.Erc20Keeper); err != nil {
		return err
	}
	if err := migrateGovDefaultParams(ctx, app.GovKeeper); err != nil {
		return err
	}
	if err := redeployTestnetContract(ctx, app.AccountKeeper, app.EvmKeeper, app.Erc20Keeper, app.EthKeeper); err != nil {
		return err
	}
	if err := migrateCrosschainParams(ctx, app.CrosschainKeepers); err != nil {
		return err
	}

	return migrateModulesData(ctx, codec, app)
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

	if err = migrateGovCustomParam(ctx, app.GovKeeper, app.GetKey(govtypes.StoreKey)); err != nil {
		return fromVM, err
	}
	if err = migrateGovDefaultParams(ctx, app.GovKeeper); err != nil {
		return fromVM, err
	}
	if err = migrateBridgeBalance(ctx, app.BankKeeper, app.AccountKeeper); err != nil {
		return fromVM, err
	}
	if err = migrateERC20TokenToCrosschain(ctx, app.BankKeeper, app.Erc20Keeper); err != nil {
		return fromVM, err
	}
	if err = fixPundixCoin(ctx, app.EvmKeeper, app.Erc20Keeper, app.BankKeeper); err != nil {
		return fromVM, err
	}
	if err = fixPurseCoin(ctx, app.EvmKeeper, app.Erc20Keeper, app.BankKeeper); err != nil {
		return fromVM, err
	}
	if err = updateMetadata(ctx, app.BankKeeper); err != nil {
		return fromVM, err
	}

	store.RemoveStoreKeys(ctx, app.GetKey(erc20types.StoreKey), erc20v8.GetRemovedStoreKeys())

	if err = mintPurseBridgeToken(ctx, app.Erc20Keeper, app.BankKeeper); err != nil {
		return fromVM, err
	}

	acc := app.AccountKeeper.GetModuleAddress(evmtypes.ModuleName)
	moduleAddress := common.BytesToAddress(acc.Bytes())

	if err = deployBridgeFeeContract(
		ctx,
		app.EvmKeeper,
		app.Erc20Keeper,
		app.CrosschainKeepers.EthKeeper,
		moduleAddress,
	); err != nil {
		return fromVM, err
	}

	if err = deployAccessControlContract(ctx, app.EvmKeeper, moduleAddress); err != nil {
		return fromVM, err
	}

	fixBaseOracleStatus(ctx, app.CrosschainKeepers.Layer2Keeper)
	updateWPUNDIAILogicCode(ctx, app.EvmKeeper)
	updateERC20LogicCode(ctx, app.EvmKeeper)

	if err = migrateModulesData(ctx, codec, app); err != nil {
		return fromVM, err
	}
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
	if err := migrateErc20FXToPundiAI(ctx, app.Erc20Keeper); err != nil {
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

func fixPundixCoin(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper, erc20Keeper erc20keeper.Keeper, bankKeeper bankkeeper.Keeper) error {
	erc20Contract := contract.NewERC20TokenKeeper(evmKeeper)
	erc20Token, err := erc20Keeper.GetERC20Token(ctx, pundixBaseDenom)
	if err != nil {
		return err
	}
	erc20PundixSupply, err := erc20Contract.TotalSupply(ctx, erc20Token.GetERC20Contract())
	if err != nil {
		return err
	}
	erc20ModulePurseBalance := bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(erc20types.ModuleName), pundixBaseDenom)
	if !erc20ModulePurseBalance.IsZero() {
		erc20PundixSupply = big.NewInt(0).Sub(erc20PundixSupply, erc20ModulePurseBalance.Amount.BigInt())
	}

	fixCoins := sdk.NewCoins(sdk.NewCoin(pundixBaseDenom, sdkmath.NewIntFromBigInt(erc20PundixSupply)))
	return bankKeeper.MintCoins(ctx, erc20types.ModuleName, fixCoins)
}

func fixPurseCoin(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper, erc20Keeper erc20keeper.Keeper, bankKeeper bankkeeper.Keeper) error {
	erc20Contract := contract.NewERC20TokenKeeper(evmKeeper)
	purseToken, err := erc20Keeper.GetERC20Token(ctx, purseBaseDenom)
	if err != nil {
		return err
	}
	contractPurseSupply, err := erc20Contract.TotalSupply(ctx, purseToken.GetERC20Contract())
	if err != nil {
		return err
	}
	erc20ModulePurseBalance := bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(erc20types.ModuleName), purseBaseDenom)
	if !erc20ModulePurseBalance.IsZero() {
		contractPurseSupply = big.NewInt(0).Sub(contractPurseSupply, erc20ModulePurseBalance.Amount.BigInt())
	}

	ethPurseToken, err := erc20Keeper.GetBridgeToken(ctx, ethtypes.ModuleName, purseBaseDenom)
	if err != nil {
		return err
	}
	ethPurseTokenSupply := bankKeeper.GetSupply(ctx, ethPurseToken.BridgeDenom())

	bscPurseToken, err := erc20Keeper.GetBridgeToken(ctx, bsctypes.ModuleName, purseBaseDenom)
	if err != nil {
		return err
	}
	bscPurseTokenSupply := bankKeeper.GetSupply(ctx, bscPurseToken.BridgeDenom())

	ibcPurseToken, err := erc20Keeper.GetIBCToken(ctx, "channel-0", purseBaseDenom)
	if err != nil {
		return err
	}

	erc20PurseSupply := sdk.NewCoin(purseBaseDenom, sdkmath.NewIntFromBigInt(contractPurseSupply))
	if err = bankKeeper.MintCoins(ctx, erc20types.ModuleName, sdk.NewCoins(erc20PurseSupply)); err != nil {
		return err
	}

	crosschainPurseSupply := sdk.NewCoin(purseBaseDenom, bscPurseTokenSupply.Amount.Add(ethPurseTokenSupply.Amount))
	if err = bankKeeper.MintCoins(ctx, crosschaintypes.ModuleName, sdk.NewCoins(crosschainPurseSupply)); err != nil {
		return err
	}

	basePurseSupply := bankKeeper.GetSupply(ctx, purseBaseDenom)
	ibcPurseSupply := bankKeeper.GetSupply(ctx, ibcPurseToken.GetIbcDenom())

	needIBCPurseSupply := sdk.NewCoin(ibcPurseToken.GetIbcDenom(), basePurseSupply.Amount.Sub(ibcPurseSupply.Amount))
	return bankKeeper.MintCoins(ctx, crosschaintypes.ModuleName, sdk.NewCoins(needIBCPurseSupply))
}

func fixTestnetTokenAmount(ctx sdk.Context, bankKeeper bankkeeper.Keeper, evmKeeper *fxevmkeeper.Keeper, erc20Keeper erc20keeper.Keeper) error {
	// fx1ntaua8eyzefqwva6evmsx9wn9d4jcs7klnvvzn 0x9afbcE9F2416520733BAcb370315D32B6B2c43d6
	fixAddress := authtypes.NewModuleAddress("testnet")
	fixTokens := getTestnetTokenAmount(ctx)
	for denom, amount := range fixTokens {
		coins := sdk.NewCoins(sdk.NewCoin(denom, amount.Abs()))
		if amount.IsNegative() {
			if err := bankKeeper.MintCoins(ctx, erc20types.ModuleName, coins); err != nil {
				return err
			}
			if err := bankKeeper.SendCoinsFromModuleToAccount(ctx, erc20types.ModuleName, fixAddress, coins); err != nil {
				return err
			}
			continue
		}

		erc20Token, err := erc20Keeper.GetERC20Token(ctx, denom)
		if err != nil {
			return err
		}
		if !erc20Token.IsNativeCoin() {
			return fmt.Errorf("token %s is not native coin", denom)
		}
		tokenKeeper := contract.NewERC20TokenKeeper(evmKeeper)
		if _, err = tokenKeeper.Burn(ctx, erc20Token.GetERC20Contract(), common.BytesToAddress(fixAddress.Bytes()), amount.BigInt()); err != nil {
			return err
		}
		if err = bankKeeper.BurnCoins(ctx, erc20types.ModuleName, coins); err != nil {
			return err
		}
	}
	return nil
}
