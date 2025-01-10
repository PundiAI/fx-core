package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-metrics"

	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (k Keeper) BaseCoinToEvm(ctx context.Context, caller contract.Caller, holder common.Address, coin sdk.Coin) (string, error) {
	erc20Address, err := k.ConvertCoin(ctx, caller, holder.Bytes(), holder, coin)
	if err != nil {
		return "", err
	}
	return erc20Address, nil
}

func (k Keeper) ConvertCoin(ctx context.Context, caller contract.Caller, sender sdk.AccAddress, receiver common.Address, coin sdk.Coin) (erc20Addr string, err error) {
	erc20Token, err := k.MintingEnabled(ctx, receiver.Bytes(), true, coin.Denom)
	if err != nil {
		return erc20Addr, err
	}

	// Check ownership and execute conversion
	switch {
	case erc20Token.IsNativeCoin():
		err = k.ConvertCoinNativeCoin(ctx, caller, erc20Token, sender, receiver, coin)
	case erc20Token.IsNativeERC20():
		err = k.ConvertCoinNativeERC20(ctx, caller, erc20Token, sender, receiver, coin)
	default:
		return erc20Addr, types.ErrUndefinedOwner
	}
	if err != nil {
		return erc20Addr, err
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tx", "msg", "convert", "coin", "total"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("denom", erc20Token.Denom),
				telemetry.NewLabel("erc20", erc20Token.Erc20Address),
			},
		)
	}()

	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeConvertCoin,
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		sdk.NewAttribute(sdk.AttributeKeyAmount, coin.Amount.String()),
		sdk.NewAttribute(types.AttributeKeyDenom, coin.Denom),
		sdk.NewAttribute(types.AttributeKeyTokenAddress, erc20Token.Erc20Address),
	))
	return erc20Token.Erc20Address, nil
}

func (k Keeper) ConvertCoinNativeCoin(ctx context.Context, caller contract.Caller, erc20Token types.ERC20Token, sender sdk.AccAddress, receiver common.Address, coin sdk.Coin) error {
	// NOTE: ignore validation from NewCoin constructor
	coins := sdk.Coins{coin}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins); err != nil {
		return err
	}

	erc20TokenKeeper := contract.NewERC20TokenKeeper(caller)
	erc20Contract := erc20Token.GetERC20Contract()
	if _, err := erc20TokenKeeper.Mint(ctx, erc20Contract, k.contractOwner, receiver, coin.Amount.BigInt()); err != nil {
		return err
	}

	if erc20Token.Denom == fxtypes.DefaultDenom {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, erc20Contract.Bytes(), coins); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) ConvertCoinNativeERC20(ctx context.Context, caller contract.Caller, erc20Token types.ERC20Token, sender sdk.AccAddress, receiver common.Address, coin sdk.Coin) error {
	// NOTE: ignore validation from NewCoin constructor
	coins := sdk.Coins{coin}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins); err != nil {
		return err
	}

	erc20TokenKeeper := contract.NewERC20TokenKeeper(caller)
	if _, err := erc20TokenKeeper.Transfer(ctx, erc20Token.GetERC20Contract(), k.contractOwner, receiver, coin.Amount.BigInt()); err != nil {
		return err
	}

	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
}
