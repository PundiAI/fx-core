package keeper

import (
	"context"
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-metrics"

	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/erc20/types"
)

func (k Keeper) EvmToBaseCoin(ctx context.Context, holder common.Address, amount *big.Int, tokenAddr string) (sdk.Coin, error) {
	sdkInt := sdkmath.NewIntFromBigInt(amount)
	baseDenom, err := k.ConvertERC20(ctx, holder, holder.Bytes(), tokenAddr, sdkInt)
	if err != nil {
		return sdk.Coin{}, err
	}
	return sdk.NewCoin(baseDenom, sdkInt), nil
}

func (k Keeper) BaseCoinToEvm(ctx context.Context, holder common.Address, coin sdk.Coin) (string, error) {
	erc20Address, err := k.ConvertCoin(ctx, holder.Bytes(), holder, coin)
	if err != nil {
		return "", err
	}
	return erc20Address, nil
}

func (k Keeper) ConvertCoin(ctx context.Context, sender sdk.AccAddress, receiver common.Address, coin sdk.Coin) (erc20Addr string, err error) {
	erc20Token, err := k.MintingEnabled(ctx, receiver.Bytes(), true, coin.Denom)
	if err != nil {
		return erc20Addr, err
	}

	// Check ownership and execute conversion
	switch {
	case erc20Token.IsNativeCoin():
		err = k.ConvertCoinNativeCoin(ctx, erc20Token, sender, receiver, coin)
	case erc20Token.IsNativeERC20():
		err = k.ConvertCoinNativeERC20(ctx, erc20Token, sender, receiver, coin)
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

func (k Keeper) ConvertCoinNativeCoin(ctx context.Context, erc20Token types.ERC20Token, sender sdk.AccAddress, receiver common.Address, coin sdk.Coin) error {
	// NOTE: ignore validation from NewCoin constructor
	coins := sdk.Coins{coin}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins); err != nil {
		return err
	}

	erc20Contract := erc20Token.GetERC20Contract()
	if err := k.evmErc20Keeper.ERC20Mint(ctx, erc20Contract, k.contractOwner, receiver, coin.Amount.BigInt()); err != nil {
		return err
	}

	if erc20Token.Denom == fxtypes.DefaultDenom {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, erc20Contract.Bytes(), coins); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) ConvertCoinNativeERC20(ctx context.Context, erc20Token types.ERC20Token, sender sdk.AccAddress, receiver common.Address, coin sdk.Coin) error {
	// NOTE: ignore validation from NewCoin constructor
	coins := sdk.Coins{coin}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins); err != nil {
		return err
	}

	if err := k.evmErc20Keeper.ERC20Transfer(ctx, erc20Token.GetERC20Contract(), k.contractOwner, receiver, coin.Amount.BigInt()); err != nil {
		return err
	}

	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
}

func (k Keeper) ConvertERC20(ctx context.Context, sender common.Address, receiver sdk.AccAddress, contractAddr string, amount sdkmath.Int) (baseDenom string, err error) {
	erc20Token, err := k.MintingEnabled(ctx, receiver, false, contractAddr)
	if err != nil {
		return baseDenom, err
	}

	// Check ownership and execute conversion
	switch {
	case erc20Token.IsNativeCoin():
		err = k.ConvertERC20NativeCoin(ctx, erc20Token, sender, receiver, amount)
	case erc20Token.IsNativeERC20():
		err = k.ConvertERC20NativeToken(ctx, erc20Token, sender, receiver, amount)
	default:
		return baseDenom, types.ErrUndefinedOwner
	}
	if err != nil {
		return baseDenom, err
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tx", "msg", "convert", "erc20", "total"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("denom", erc20Token.Denom),
				telemetry.NewLabel("erc20", erc20Token.Erc20Address),
			},
		)
	}()

	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeConvertERC20,
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
		sdk.NewAttribute(types.AttributeKeyDenom, erc20Token.Denom),
		sdk.NewAttribute(types.AttributeKeyTokenAddress, contractAddr),
	))
	return erc20Token.Denom, nil
}

func (k Keeper) ConvertERC20NativeCoin(ctx context.Context, erc20Token types.ERC20Token, sender common.Address, receiver sdk.AccAddress, amount sdkmath.Int) error {
	erc20Contract := erc20Token.GetERC20Contract()

	if err := k.evmErc20Keeper.ERC20Burn(ctx, erc20Contract, k.contractOwner, sender, amount.BigInt()); err != nil {
		return err
	}

	// NOTE: coin fields already validated
	coins := sdk.Coins{sdk.Coin{Denom: erc20Token.Denom, Amount: amount}}

	if erc20Token.Denom == fxtypes.DefaultDenom {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, erc20Contract.Bytes(), types.ModuleName, coins); err != nil {
			return err
		}
	}

	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, coins)
}

func (k Keeper) ConvertERC20NativeToken(ctx context.Context, erc20Token types.ERC20Token, sender common.Address, receiver sdk.AccAddress, amount sdkmath.Int) error {
	erc20Contract := erc20Token.GetERC20Contract()

	if err := k.evmErc20Keeper.ERC20Transfer(ctx, erc20Contract, sender, k.contractOwner, amount.BigInt()); err != nil {
		return err
	}

	// NOTE: coin fields already validated
	coins := sdk.Coins{sdk.Coin{Denom: erc20Token.Denom, Amount: amount}}
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, coins)
}
