package v8

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"

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
	fxgovkeeper "github.com/pundiai/fx-core/v8/x/gov/keeper"
	fxgovv8 "github.com/pundiai/fx-core/v8/x/gov/migrations/v8"
	layer2types "github.com/pundiai/fx-core/v8/x/layer2/types"
	fxstakingv8 "github.com/pundiai/fx-core/v8/x/staking/migrations/v8"
)

func CreateUpgradeHandler(cdc codec.Codec, mm *module.Manager, configurator module.Configurator, app *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := sdk.UnwrapSDKContext(ctx).CacheContext()

		var err error
		var toVM module.VersionMap
		if cacheCtx.ChainID() == fxtypes.TestnetChainId {
			if err = upgradeTestnet(cacheCtx, app); err != nil {
				return fromVM, err
			}
			toVM = fromVM
		} else {
			toVM, err = upgradeMainnet(cacheCtx, mm, configurator, app, fromVM, plan)
			if err != nil {
				return fromVM, err
			}
		}
		commit()
		cacheCtx.Logger().Info("upgrade complete", "module", "upgrade")
		return toVM, nil
	}
}

func upgradeTestnet(ctx sdk.Context, app *keepers.AppKeepers) error {
	fixBaseOracleStatus(ctx, app.CrosschainKeepers.Layer2Keeper)

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
	if err := migrateFeemarketGasPrice(ctx, app.FeeMarketKeeper); err != nil {
		return err
	}
	if err := redeployTestnetContract(ctx, app.AccountKeeper, app.EvmKeeper, app.Erc20Keeper, app.EthKeeper); err != nil {
		return err
	}
	if err := migrateCrosschainParams(ctx, app.CrosschainKeepers); err != nil {
		return err
	}
	if err := migrateMetadataDisplay(ctx, app.BankKeeper); err != nil {
		return err
	}
	if err := migrateErc20FXToPundiAI(ctx, app.Erc20Keeper); err != nil {
		return err
	}

	migrationWFXToWPUNDIAI(ctx, app.EvmKeeper)

	return migrateMetadataFXToPundiAI(ctx, app.BankKeeper)
}

//nolint:gocyclo // mainnet
func upgradeMainnet(
	ctx sdk.Context,
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

	migrationWFXToWPUNDIAI(ctx, app.EvmKeeper)

	if err = migrateEvmParams(ctx, app.EvmKeeper); err != nil {
		return fromVM, err
	}

	store.RemoveStoreKeys(ctx, app.GetKey(stakingtypes.StoreKey), fxstakingv8.GetRemovedStoreKeys())

	if err = migrationGovCustomParam(ctx, app.GovKeeper, app.GetKey(govtypes.StoreKey)); err != nil {
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

	if err = migrateFeemarketGasPrice(ctx, app.FeeMarketKeeper); err != nil {
		return toVM, err
	}
	if err = migrateMetadataDisplay(ctx, app.BankKeeper); err != nil {
		return toVM, err
	}
	if err = migrateErc20FXToPundiAI(ctx, app.Erc20Keeper); err != nil {
		return toVM, err
	}
	if err = migrateMetadataFXToPundiAI(ctx, app.BankKeeper); err != nil {
		return toVM, err
	}
	return toVM, nil
}

func migrateEvmParams(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper) error {
	params := evmKeeper.GetParams(ctx)
	params.HeaderHashNum = evmtypes.DefaultHeaderHashNum
	return evmKeeper.SetParams(ctx, params)
}

func migrationGovCustomParam(ctx sdk.Context, keeper *fxgovkeeper.Keeper, storeKey *storetypes.KVStoreKey) error {
	// 1. delete fxParams key
	store.RemoveStoreKeys(ctx, storeKey, fxgovv8.GetRemovedStoreKeys())

	// 2. init custom params
	if err := keeper.InitCustomParams(ctx); err != nil {
		return err
	}

	// 3. set default params
	return migrateGovDefaultParams(ctx, keeper)
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
	baseAcc, ok := acc.(*authtypes.BaseAccount)
	if !ok {
		return sdkerrors.ErrInvalidAddress.Wrapf("crosschain account not base account")
	}
	macc := authtypes.NewModuleAccount(baseAcc, crosschaintypes.ModuleName, perms...)
	ak.SetModuleAccount(ctx, macc)
	return nil
}

func migrateBridgeBalance(ctx sdk.Context, bankKeeper bankkeeper.Keeper, accountKeeper authkeeper.AccountKeeper) error {
	mds := bankKeeper.GetAllDenomMetaData(ctx)
	for _, md := range mds {
		if md.Base == fxtypes.DefaultDenom || (len(md.DenomUnits) == 0 || len(md.DenomUnits[0].Aliases) == 0) && md.Symbol != pundixSymbol {
			continue
		}
		dstBase := md.Base
		if !strings.Contains(md.Base, strings.ToLower(md.Symbol)) {
			dstBase = strings.ToLower(md.Symbol)
		}
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
	balances := bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(erc20types.ModuleName))
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

func updateMetadata(ctx sdk.Context, bankKeeper bankkeeper.Keeper) error {
	mds := bankKeeper.GetAllDenomMetaData(ctx)

	removeMetadata := make([]string, 0, 2)
	for _, md := range mds {
		if md.Base == fxtypes.DefaultDenom || (len(md.DenomUnits) == 0 || len(md.DenomUnits[0].Aliases) == 0) && md.Symbol != pundixSymbol {
			continue
		}
		// remove alias
		md.DenomUnits[0].Aliases = []string{}

		newBase := strings.ToLower(md.Symbol)
		// update pundix/purse base denom
		if md.Base != newBase && !strings.Contains(md.Base, newBase) && !strings.HasPrefix(md.Display, ibctransfertypes.ModuleName+"/"+ibcchanneltypes.ChannelPrefix) {
			removeMetadata = append(removeMetadata, md.Base)

			md.Base = newBase
			md.Display = newBase
			md.DenomUnits[0].Denom = newBase
		}

		bankKeeper.SetDenomMetaData(ctx, md)
	}

	bk, ok := bankKeeper.(bankkeeper.BaseKeeper)
	if !ok {
		return errors.New("bank keeper not implement bank.BaseKeeper")
	}
	for _, base := range removeMetadata {
		if !bankKeeper.HasDenomMetaData(ctx, base) {
			continue
		}
		ctx.Logger().Info("remove metadata", "base", base)
		if err := bk.BaseViewKeeper.DenomMetadata.Remove(ctx, base); err != nil {
			return err
		}
	}
	return nil
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
	chains := fxtypes.GetSupportChains()
	bridgeDenoms := make([]contract.BridgeDenoms, len(chains))
	for index, chain := range chains {
		denoms := make([]common.Hash, 0)
		bridgeTokens, err := erc20Keeper.GetBridgeTokens(cacheCtx, chain)
		if err != nil {
			return err
		}
		for _, token := range bridgeTokens {
			denoms = append(denoms, contract.MustStrToByte32(token.GetDenom()))
		}
		bridgeDenoms[index] = contract.BridgeDenoms{
			ChainName: contract.MustStrToByte32(chain),
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
		bridgeDenoms,
		evmModuleAddress,
		getContractOwner(cacheCtx),
		common.HexToAddress(oracles[0].ExternalAddress),
	)
}

func deployAccessControlContract(cacheCtx sdk.Context, evmKeeper *fxevmkeeper.Keeper, evmModuleAddress common.Address) error {
	return contract.DeployAccessControlContract(
		cacheCtx,
		evmKeeper,
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

func fixPundixCoin(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper, erc20Keeper erc20keeper.Keeper, bankKeeper bankkeeper.Keeper) error {
	erc20Contract := contract.NewERC20TokenKeeper(evmKeeper)
	erc20Token, err := erc20Keeper.GetERC20Token(ctx, pundixBase)
	if err != nil {
		return err
	}
	erc20PundixSupply, err := erc20Contract.TotalSupply(ctx, erc20Token.GetERC20Contract())
	if err != nil {
		return err
	}
	erc20ModulePurseBalance := bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(erc20types.ModuleName), pundixBase)
	if !erc20ModulePurseBalance.IsZero() {
		erc20PundixSupply = big.NewInt(0).Sub(erc20PundixSupply, erc20ModulePurseBalance.Amount.BigInt())
	}

	fixCoins := sdk.NewCoins(sdk.NewCoin(pundixBase, sdkmath.NewIntFromBigInt(erc20PundixSupply)))
	return bankKeeper.MintCoins(ctx, erc20types.ModuleName, fixCoins)
}

func fixPurseCoin(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper, erc20Keeper erc20keeper.Keeper, bankKeeper bankkeeper.Keeper) error {
	erc20Contract := contract.NewERC20TokenKeeper(evmKeeper)
	purseToken, err := erc20Keeper.GetERC20Token(ctx, purseBase)
	if err != nil {
		return err
	}
	contractPurseSupply, err := erc20Contract.TotalSupply(ctx, purseToken.GetERC20Contract())
	if err != nil {
		return err
	}
	erc20ModulePurseBalance := bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(erc20types.ModuleName), purseBase)
	if !erc20ModulePurseBalance.IsZero() {
		contractPurseSupply = big.NewInt(0).Sub(contractPurseSupply, erc20ModulePurseBalance.Amount.BigInt())
	}

	ethPurseToken, err := erc20Keeper.GetBridgeToken(ctx, ethtypes.ModuleName, purseBase)
	if err != nil {
		return err
	}
	ethPurseTokenSupply := bankKeeper.GetSupply(ctx, ethPurseToken.BridgeDenom())

	bscPurseToken, err := erc20Keeper.GetBridgeToken(ctx, bsctypes.ModuleName, purseBase)
	if err != nil {
		return err
	}
	bscPurseTokenSupply := bankKeeper.GetSupply(ctx, bscPurseToken.BridgeDenom())

	ibcPurseToken, err := erc20Keeper.GetIBCToken(ctx, "channel-0", purseBase)
	if err != nil {
		return err
	}

	erc20PurseSupply := sdk.NewCoin(purseBase, sdkmath.NewIntFromBigInt(contractPurseSupply))
	if err = bankKeeper.MintCoins(ctx, erc20types.ModuleName, sdk.NewCoins(erc20PurseSupply)); err != nil {
		return err
	}

	crosschainPurseSupply := sdk.NewCoin(purseBase, bscPurseTokenSupply.Amount.Add(ethPurseTokenSupply.Amount))
	if err = bankKeeper.MintCoins(ctx, crosschaintypes.ModuleName, sdk.NewCoins(crosschainPurseSupply)); err != nil {
		return err
	}

	basePurseSupply := bankKeeper.GetSupply(ctx, purseBase)
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
		erc20ModuleAddress := common.BytesToAddress(authtypes.NewModuleAddress(erc20types.ModuleName))
		if _, err = tokenKeeper.Burn(ctx, erc20Token.GetERC20Contract(), erc20ModuleAddress,
			common.BytesToAddress(fixAddress.Bytes()), amount.BigInt()); err != nil {
			return err
		}
		if err = bankKeeper.BurnCoins(ctx, erc20types.ModuleName, coins); err != nil {
			return err
		}
	}
	return nil
}

func redeployTestnetContract(
	ctx sdk.Context,
	accountKeeper authkeeper.AccountKeeper,
	evmKeeper *fxevmkeeper.Keeper,
	erc20Keeper erc20keeper.Keeper,
	ethKeeper crosschainkeeper.Keeper,
) error {
	if err := evmKeeper.DeleteAccount(ctx, common.HexToAddress(contract.BridgeFeeAddress)); err != nil {
		return err
	}
	if err := evmKeeper.DeleteAccount(ctx, common.HexToAddress(contract.BridgeFeeOracleAddress)); err != nil {
		return err
	}

	acc := accountKeeper.GetModuleAddress(evmtypes.ModuleName)
	moduleAddress := common.BytesToAddress(acc.Bytes())
	return deployBridgeFeeContract(
		ctx,
		evmKeeper,
		erc20Keeper,
		ethKeeper,
		moduleAddress,
	)
}

func migrateGovDefaultParams(ctx sdk.Context, keeper *fxgovkeeper.Keeper) error {
	params, err := keeper.Params.Get(ctx)
	if err != nil {
		return err
	}

	minDepositAmount := sdkmath.NewInt(1e18).MulRaw(30)

	params.MinDeposit = sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, minDepositAmount))
	params.ExpeditedMinDeposit = sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, minDepositAmount.MulRaw(govv1.DefaultMinExpeditedDepositTokensRatio)))
	params.MinInitialDepositRatio = sdkmath.LegacyMustNewDecFromStr("0.33").String()
	params.MinDepositRatio = sdkmath.LegacyMustNewDecFromStr("0").String()

	return keeper.Params.Set(ctx, params)
}

func migrateFeemarketGasPrice(ctx sdk.Context, feemarketKeeper feemarketkeeper.Keeper) error {
	params := feemarketKeeper.GetParams(ctx)
	params.BaseFee = sdkmath.NewInt(fxtypes.DefaultGasPrice)
	params.MinGasPrice = sdkmath.LegacyNewDec(fxtypes.DefaultGasPrice)
	return feemarketKeeper.SetParams(ctx, params)
}

func migrateCrosschainParams(ctx sdk.Context, keepers keepers.CrosschainKeepers) error {
	for _, k := range keepers.ToSlice() {
		params := k.GetParams(ctx)
		params.DelegateThreshold.Denom = fxtypes.DefaultDenom
		if err := k.SetParams(ctx, &params); err != nil {
			return err
		}
	}
	return nil
}

func migrateMetadataDisplay(ctx sdk.Context, bankKeeper bankkeeper.Keeper) error {
	mds := bankKeeper.GetAllDenomMetaData(ctx)
	for _, md := range mds {
		if md.Display != md.Base || len(md.DenomUnits) <= 1 {
			continue
		}
		for _, dus := range md.DenomUnits {
			if dus.Denom != md.Base {
				md.Display = dus.Denom
				break
			}
		}
		if err := md.Validate(); err != nil {
			return err
		}
		bankKeeper.SetDenomMetaData(ctx, md)
	}
	return nil
}

func migrateErc20FXToPundiAI(ctx sdk.Context, keeper erc20keeper.Keeper) error {
	fxDenom := strings.ToUpper(fxtypes.FXDenom)
	erc20Token, err := keeper.GetERC20Token(ctx, fxDenom)
	if err != nil {
		return err
	}
	erc20Token.Denom = fxtypes.DefaultDenom
	if err = keeper.ERC20Token.Set(ctx, erc20Token.Denom, erc20Token); err != nil {
		return err
	}
	return keeper.ERC20Token.Remove(ctx, fxDenom)
}

func migrateMetadataFXToPundiAI(ctx sdk.Context, keeper bankkeeper.Keeper) error {
	// add pundiai metadata
	metadata := fxtypes.NewDefaultMetadata()
	keeper.SetDenomMetaData(ctx, metadata)

	// remove FX metadata
	bk, ok := keeper.(bankkeeper.BaseKeeper)
	if !ok {
		return errors.New("bank keeper not implement bank.BaseKeeper")
	}
	return bk.BaseViewKeeper.DenomMetadata.Remove(ctx, strings.ToUpper(fxtypes.FXDenom))
}
