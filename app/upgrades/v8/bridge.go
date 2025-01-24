package v8

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	bsctypes "github.com/pundiai/fx-core/v8/x/bsc/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20keeper "github.com/pundiai/fx-core/v8/x/erc20/keeper"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
)

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
	ctx.Logger().Info("migrate erc20 bridge/ibc token to crosschain", "module", "upgrade", "coins", migrateCoins.String())
	return bankKeeper.SendCoinsFromModuleToModule(ctx, erc20types.ModuleName, crosschaintypes.ModuleName, migrateCoins)
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

func migrateBridgeBalance(ctx sdk.Context, bankKeeper bankkeeper.Keeper, accountKeeper authkeeper.AccountKeeper) error {
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
