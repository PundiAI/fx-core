package v4_2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v5/app/keepers"
	v4 "github.com/functionx/fx-core/v5/app/upgrades/v4"
	fxtypes "github.com/functionx/fx-core/v5/types"
	crosschainkeeper "github.com/functionx/fx-core/v5/x/crosschain/keeper"
	erc20types "github.com/functionx/fx-core/v5/x/erc20/types"
)

func createUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	app *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// testnet upgrade
		if ctx.ChainID() == fxtypes.TestnetChainId {
			return fromVM, nil
		}

		// mainnet upgrade
		toVM, err := v4.CreateUpgradeHandler(mm, configurator, app)(ctx, plan, fromVM)
		if err != nil {
			return nil, err
		}

		// refund polygon USDT, after erc20 module params migrate
		for _, r := range PolygonUSDTRefunds {
			cacheCtx, commit := ctx.CacheContext()
			if err := CrossChainRefundToAccount(cacheCtx, app, app.PolygonKeeper, r.Address, r.Coins); err != nil {
				ctx.Logger().Error("refund failed", "addr", r.Address.String(), "coins", r.Coins.String(), "err", err.Error())
				continue
			}
			commit()
			ctx.Logger().Info("refund success", "addr", r.Address.String(), "coins", r.Coins.String())
		}
		return toVM, nil
	}
}

func CrossChainRefundToAccount(ctx sdk.Context, app *keepers.AppKeepers, chk crosschainkeeper.Keeper, addr common.Address, coins sdk.Coins) error {
	if err := app.BankKeeper.MintCoins(ctx, chk.ModuleName(), coins); err != nil {
		return err
	}
	if err := app.BankKeeper.SendCoinsFromModuleToAccount(ctx, chk.ModuleName(), addr.Bytes(), coins); err != nil {
		return err
	}
	for _, coin := range coins {
		if coin.Denom == fxtypes.DefaultDenom {
			ctx.Logger().Info("skip refund", "addr", addr.String(), "coin", coin.String())
			continue
		}
		targetCoin, err := app.Erc20Keeper.ConvertDenomToTarget(ctx, addr.Bytes(), coin, fxtypes.ParseFxTarget(fxtypes.ERC20Target))
		if err != nil {
			return err
		}
		if _, err = app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(ctx), &erc20types.MsgConvertCoin{
			Coin:     targetCoin,
			Receiver: addr.Hex(),
			Sender:   sdk.AccAddress(addr.Bytes()).String(),
		}); err != nil {
			return err
		}
	}
	return nil
}
