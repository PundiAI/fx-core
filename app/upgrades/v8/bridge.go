package v8

import (
	"fmt"
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	bsctypes "github.com/pundiai/fx-core/v8/x/bsc/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20keeper "github.com/pundiai/fx-core/v8/x/erc20/keeper"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
	fxevmkeeper "github.com/pundiai/fx-core/v8/x/evm/keeper"
)

func migrateBridgeToken(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper, erc20Keeper erc20keeper.Keeper, bankKeeper bankkeeper.Keeper, accountKeeper authkeeper.AccountKeeper) error {
	if err := migrateAccountTokenBalance(ctx, bankKeeper, accountKeeper); err != nil {
		return err
	}
	if err := migrateERC20ModulePundix(ctx, evmKeeper, erc20Keeper, bankKeeper); err != nil {
		return err
	}
	if err := migrateERC20ModulePurse(ctx, evmKeeper, erc20Keeper, bankKeeper); err != nil {
		return err
	}
	if err := migrateEthModulePurse(ctx, erc20Keeper, bankKeeper); err != nil {
		return err
	}
	if err := migrateBscModulePurse(ctx, erc20Keeper, bankKeeper); err != nil {
		return err
	}
	return migrateERC20ModuleTokens(ctx, evmKeeper, erc20Keeper, bankKeeper)
}

func migrateAccountTokenBalance(ctx sdk.Context, bankKeeper bankkeeper.Keeper, accountKeeper authkeeper.AccountKeeper) error {
	mds := bankKeeper.GetAllDenomMetaData(ctx)
	for _, md := range mds {
		if md.Base == fxtypes.LegacyFXDenom || (len(md.DenomUnits) == 0 || len(md.DenomUnits[0].Aliases) == 0) && md.Symbol != pundixSymbol {
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

		ctx.Logger().Info("migrate coin", "module", "upgrade", "address", address.String(), "src-denom", srcBase, "dst-denom", dstBase, "amount", coin.Amount.String())
		if err = bankKeeper.SendCoinsFromAccountToModule(ctx, address, crosschaintypes.ModuleName, sdk.NewCoins(coin)); err != nil {
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

func migrateERC20ModulePundix(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper, erc20Keeper erc20keeper.Keeper, bankKeeper bankkeeper.Keeper) error {
	erc20Contract := contract.NewERC20TokenKeeper(evmKeeper)
	erc20Token, err := erc20Keeper.GetERC20Token(ctx, pundixBaseDenom)
	if err != nil {
		return err
	}
	erc20Supply, err := erc20Contract.TotalSupply(ctx, erc20Token.GetERC20Contract())
	if err != nil {
		return err
	}
	bridgeToken, err := erc20Keeper.GetBridgeToken(ctx, ethtypes.ModuleName, pundixBaseDenom)
	if err != nil {
		return err
	}
	erc20ModuleBalance := bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(erc20types.ModuleName), bridgeToken.BridgeDenom())
	if !erc20ModuleBalance.Amount.Equal(sdkmath.NewIntFromBigInt(erc20Supply)) {
		return fmt.Errorf("erc20 module pundix supply not equal erc20 balance, moduel %s, erc20 %s",
			erc20ModuleBalance.Amount.String(), erc20Supply.String())
	}

	baseCoin := sdk.NewCoins(sdk.NewCoin(pundixBaseDenom, erc20ModuleBalance.Amount))
	ctx.Logger().Info("migrate erc20 module", "bridge-coin", erc20ModuleBalance.String(), "base-coin", baseCoin.String())
	if err = bankKeeper.SendCoinsFromModuleToModule(ctx, erc20types.ModuleName, crosschaintypes.ModuleName, sdk.NewCoins(erc20ModuleBalance)); err != nil {
		return err
	}
	return bankKeeper.MintCoins(ctx, erc20types.ModuleName, baseCoin)
}

func migrateERC20ModulePurse(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper, erc20Keeper erc20keeper.Keeper, bankKeeper bankkeeper.Keeper) error {
	erc20Contract := contract.NewERC20TokenKeeper(evmKeeper)
	purseToken, err := erc20Keeper.GetERC20Token(ctx, purseBaseDenom)
	if err != nil {
		return err
	}
	erc20Supply, err := erc20Contract.TotalSupply(ctx, purseToken.GetERC20Contract())
	if err != nil {
		return err
	}
	ibcToken, err := erc20Keeper.GetIBCToken(ctx, fxtypes.PundixChannel, purseBaseDenom)
	if err != nil {
		return err
	}
	erc20ModuleBalance := bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(erc20types.ModuleName), ibcToken.GetIbcDenom())
	if erc20ModuleBalance.Amount.LT(sdkmath.NewIntFromBigInt(erc20Supply)) {
		return fmt.Errorf("erc20 module ibc purse supply small than contract balance, moduel %s, erc20 %s",
			erc20ModuleBalance.Amount.String(), erc20Supply.String())
	}

	fixModuleBalance := sdk.NewCoins(sdk.NewCoin(erc20ModuleBalance.GetDenom(), sdkmath.NewIntFromBigInt(erc20Supply)))
	erc20BaseSupply := sdk.NewCoin(purseBaseDenom, sdkmath.NewIntFromBigInt(erc20Supply))
	ctx.Logger().Info("migrate erc20 module", "bridge-coin", fixModuleBalance.String(), "base-coin", erc20BaseSupply.String())
	if err = bankKeeper.SendCoinsFromModuleToModule(ctx, erc20types.ModuleName, crosschaintypes.ModuleName, fixModuleBalance); err != nil {
		return err
	}
	return bankKeeper.MintCoins(ctx, erc20types.ModuleName, sdk.NewCoins(erc20BaseSupply))
}

func migrateEthModulePurse(ctx sdk.Context, erc20Keeper erc20keeper.Keeper, bankKeeper bankkeeper.Keeper) error {
	ethToken, err := erc20Keeper.GetBridgeToken(ctx, ethtypes.ModuleName, purseBaseDenom)
	if err != nil {
		return err
	}
	ethTokenSupply := bankKeeper.GetSupply(ctx, ethToken.BridgeDenom())

	ibcToken, err := erc20Keeper.GetIBCToken(ctx, fxtypes.PundixChannel, purseBaseDenom)
	if err != nil {
		return err
	}
	erc20ModuleBalance := bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(erc20types.ModuleName), ibcToken.GetIbcDenom())
	if erc20ModuleBalance.Amount.LT(ethTokenSupply.Amount) {
		return fmt.Errorf("erc20 module ibc purse supply small than eth purse supply, erc20 %s, eth %s",
			erc20ModuleBalance.String(), ethTokenSupply.String())
	}

	fixModuleBalance := sdk.NewCoins(sdk.NewCoin(ibcToken.GetIbcDenom(), ethTokenSupply.Amount))
	addPurseBaseSupply := sdk.NewCoin(purseBaseDenom, ethTokenSupply.Amount)
	ctx.Logger().Info("migrate eth module", "bridge-coin", fixModuleBalance.String(), "base-coin", addPurseBaseSupply.String())
	if err = bankKeeper.SendCoinsFromModuleToModule(ctx, erc20types.ModuleName, crosschaintypes.ModuleName, fixModuleBalance); err != nil {
		return err
	}
	return bankKeeper.MintCoins(ctx, crosschaintypes.ModuleName, sdk.NewCoins(addPurseBaseSupply))
}

func migrateBscModulePurse(ctx sdk.Context, erc20Keeper erc20keeper.Keeper, bankKeeper bankkeeper.Keeper) error {
	bscToken, err := erc20Keeper.GetBridgeToken(ctx, bsctypes.ModuleName, purseBaseDenom)
	if err != nil {
		return err
	}
	bscTokenSupply := bankKeeper.GetSupply(ctx, bscToken.BridgeDenom())
	if !bscTokenSupply.IsZero() {
		return fmt.Errorf("bsc purse supply not empty %s", bscTokenSupply.String())
	}

	ibcToken, err := erc20Keeper.GetIBCToken(ctx, fxtypes.PundixChannel, purseBaseDenom)
	if err != nil {
		return err
	}
	ibcDenomSupply := bankKeeper.GetSupply(ctx, ibcToken.GetIbcDenom())
	baseDenomSupply := bankKeeper.GetSupply(ctx, purseBaseDenom)
	if !ibcDenomSupply.Amount.Equal(baseDenomSupply.Amount) {
		return fmt.Errorf("ibc purse supply not eqaul base supply %s != %s", ibcDenomSupply.String(), baseDenomSupply.String())
	}

	pxEscrowPurse, err := getPundixEscrowPurseAmount()
	if err != nil {
		return err
	}
	bscTokeSupply := pxEscrowPurse.Sub(ibcDenomSupply.Amount)
	fixBscCoin := sdk.NewCoins(sdk.NewCoin(bscToken.BridgeDenom(), bscTokeSupply))
	fixIbcCoin := sdk.NewCoins(sdk.NewCoin(ibcToken.GetIbcDenom(), bscTokeSupply))
	fixBaseCoin := sdk.NewCoins(sdk.NewCoin(purseBaseDenom, bscTokeSupply))
	ctx.Logger().Info("migrate bsc module", "bridge-coin", fixBscCoin.String(), "base-coin", fixBaseCoin.String())
	if err = bankKeeper.MintCoins(ctx, bsctypes.ModuleName, fixBscCoin); err != nil {
		return err
	}
	if err = bankKeeper.MintCoins(ctx, crosschaintypes.ModuleName, fixIbcCoin); err != nil {
		return err
	}
	return bankKeeper.MintCoins(ctx, crosschaintypes.ModuleName, fixBaseCoin)
}

func migrateERC20ModuleTokens(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper, erc20Keeper erc20keeper.Keeper, bankKeeper bankkeeper.Keeper) error {
	erc20Contract := contract.NewERC20TokenKeeper(evmKeeper)
	migrateCoins := sdk.NewCoins()
	balances := bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(erc20types.ModuleName))
	for _, bal := range balances {
		has, err := erc20Keeper.HasToken(ctx, bal.Denom)
		if err != nil {
			return err
		}
		if has {
			migrateCoins = migrateCoins.Add(bal)
			continue
		}
		// check erc20 supply
		erc20Token, err := erc20Keeper.GetERC20Token(ctx, bal.Denom)
		if err != nil {
			return err
		}
		var supply *big.Int
		if erc20Token.IsNativeCoin() {
			supply, err = erc20Contract.TotalSupply(ctx, erc20Token.GetERC20Contract())
		} else {
			supply, err = erc20Contract.BalanceOf(ctx, erc20Token.GetERC20Contract(), common.BytesToAddress(authtypes.NewModuleAddress(erc20types.ModuleName)))
		}
		if err != nil {
			return err
		}
		if !bal.Amount.Equal(sdkmath.NewIntFromBigInt(supply)) {
			return fmt.Errorf("%s erc20 supply not equal module balance %s %s", bal.Denom, bal.Amount.String(), supply.String())
		}
	}
	ctx.Logger().Info("migrate erc20 bridge/ibc token to crosschain", "module", "upgrade", "coins", migrateCoins.String())
	return bankKeeper.SendCoinsFromModuleToModule(ctx, erc20types.ModuleName, crosschaintypes.ModuleName, migrateCoins)
}
