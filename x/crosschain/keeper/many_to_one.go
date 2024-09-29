package keeper

import (
	"context"
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v8/types"
)

func (k Keeper) BridgeTokenToBaseCoin(ctx context.Context, tokenAddr string, amount *big.Int, holder sdk.AccAddress) (sdk.Coin, error) {
	bridgeDenom, found := k.GetBridgeDenomByContract(sdk.UnwrapSDKContext(ctx), tokenAddr)
	if !found {
		return sdk.Coin{}, errortypes.ErrInvalidCoins.Wrapf("bridge denom not found %s", tokenAddr)
	}
	bridgeToken := sdk.NewCoin(bridgeDenom, sdkmath.NewIntFromBigInt(amount))
	if err := k.DepositBridgeToken(ctx, bridgeToken, holder); err != nil {
		return sdk.Coin{}, err
	}
	baseDenom, err := k.ManyToOne(ctx, bridgeToken.Denom)
	if err != nil {
		return sdk.Coin{}, err
	}
	if err = k.ConversionCoin(ctx, holder, bridgeToken, baseDenom, baseDenom); err != nil {
		return sdk.Coin{}, err
	}
	return sdk.NewCoin(baseDenom, bridgeToken.Amount), nil
}

func (k Keeper) BaseCoinToBridgeToken(ctx context.Context, module string, coin sdk.Coin, holder sdk.AccAddress) (common.Address, *big.Int, error) {
	bridgeDenom, err := k.ManyToOne(ctx, coin.Denom, module)
	if err != nil {
		return common.Address{}, nil, err
	}
	if err = k.ConversionCoin(ctx, holder, coin, coin.Denom, bridgeDenom); err != nil {
		return common.Address{}, nil, err
	}
	if err = k.WithdrawBridgeToken(ctx, sdk.NewCoin(bridgeDenom, coin.Amount), holder); err != nil {
		return common.Address{}, nil, err
	}
	tokenAddr, found := k.GetContractByBridgeDenom(sdk.UnwrapSDKContext(ctx), bridgeDenom)
	if !found {
		return common.Address{}, nil, err
	}
	return common.HexToAddress(tokenAddr), coin.Amount.BigInt(), nil
}

// DepositBridgeToken get bridge token from crosschain module
func (k Keeper) DepositBridgeToken(ctx context.Context, bridgeToken sdk.Coin, holder sdk.AccAddress) error {
	if bridgeToken.Denom == fxtypes.DefaultDenom {
		return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, holder, sdk.NewCoins(bridgeToken))
	}
	baseDenom, err := k.GetBaseDenom(ctx, bridgeToken.Denom)
	if err != nil {
		return err
	}
	tokenPair, found := k.erc20Keeper.GetTokenPair(sdk.UnwrapSDKContext(ctx), baseDenom)
	if !found {
		return errortypes.ErrInvalidCoins.Wrapf("token pair not found: %s", baseDenom)
	}

	if tokenPair.IsNativeCoin() && tokenPair.GetDenom() != fxtypes.DefaultDenom {
		if err := k.bankKeeper.MintCoins(ctx, k.moduleName, sdk.NewCoins(bridgeToken)); err != nil {
			return err
		}
	}
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, holder, sdk.NewCoins(bridgeToken))
}

// WithdrawBridgeToken put bridge token to crosschain module
func (k Keeper) WithdrawBridgeToken(ctx context.Context, bridgeToken sdk.Coin, holder sdk.AccAddress) error {
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, holder, k.moduleName, sdk.NewCoins(bridgeToken)); err != nil {
		return err
	}
	if bridgeToken.Denom == fxtypes.DefaultDenom {
		return nil
	}
	baseDenom, err := k.GetBaseDenom(ctx, bridgeToken.Denom)
	if err != nil {
		return err
	}
	tokenPair, found := k.erc20Keeper.GetTokenPair(sdk.UnwrapSDKContext(ctx), baseDenom)
	if !found {
		return errortypes.ErrInvalidCoins.Wrapf("token pair not found: %s", baseDenom)
	}
	if tokenPair.IsNativeERC20() {
		return nil
	}
	return k.bankKeeper.BurnCoins(ctx, k.moduleName, sdk.NewCoins(bridgeToken))
}

// ManyToOne get target denom by denom and target argument
func (k Keeper) ManyToOne(ctx context.Context, denom string, targets ...string) (string, error) {
	// NOTE: if empty target, convert to base denom
	target := ""
	if len(targets) > 0 && len(targets[0]) > 0 {
		target = targets[0]
	}

	// 1. check base or bridge
	baseDenom := denom
	found, err := k.HasToken(ctx, denom)
	if err != nil {
		return "", err
	}
	if !found {
		if baseDenom, err = k.GetBaseDenom(ctx, denom); err != nil {
			return "", err
		}
	}

	// 2. not need convert
	if baseDenom == fxtypes.DefaultDenom || len(target) == 0 {
		return baseDenom, nil
	}

	// 3. get target denom
	targetDenom := baseDenom
	if len(target) != 0 {
		if targetDenom, err = k.BaseDenomToBridgeDenom(ctx, baseDenom, target); err != nil {
			return "", err
		}
	}
	return targetDenom, nil
}

func (k Keeper) BaseDenomToBridgeDenom(ctx context.Context, baseDenom, target string) (string, error) {
	bridgeDenom, err := k.GetBridgeDenom(ctx, baseDenom)
	if err != nil {
		return "", err
	}
	ibcPrefix := ibctransfertypes.DenomPrefix + "/"
	for _, bd := range bridgeDenom {
		if strings.HasPrefix(bd, ibcPrefix) && strings.HasPrefix(target, ibcPrefix) {
			// TODO ibc token
			continue
		}
		if strings.HasPrefix(bd, target) {
			return bd, nil
		}
	}
	return "", errortypes.ErrInvalidCoins.Wrapf("not found bridge denom: %s, %s", baseDenom, target)
}

// ConversionCoin Convert coin between base and bridge
func (k Keeper) ConversionCoin(ctx context.Context, holder sdk.AccAddress, coin sdk.Coin, baseDenom, targetDenom string) error {
	if coin.Denom == fxtypes.DefaultDenom {
		return nil
	}
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, holder, k.moduleName, sdk.NewCoins(coin)); err != nil {
		return err
	}
	tokenPair, found := k.erc20Keeper.GetTokenPair(sdk.UnwrapSDKContext(ctx), baseDenom)
	if !found {
		return errortypes.ErrInvalidCoins.Wrapf("token pair not found %s", baseDenom)
	}

	targetCoin := sdk.NewCoin(targetDenom, coin.Amount)
	if tokenPair.IsNativeERC20() {
		if err := k.bankKeeper.BurnCoins(ctx, k.moduleName, sdk.NewCoins(coin)); err != nil {
			return err
		}
		if err := k.bankKeeper.MintCoins(ctx, k.moduleName, sdk.NewCoins(targetCoin)); err != nil {
			return err
		}
		return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, holder, sdk.NewCoins(targetCoin))
	}
	if coin.Denom == baseDenom {
		if err := k.bankKeeper.BurnCoins(ctx, k.moduleName, sdk.NewCoins(coin)); err != nil {
			return err
		}
		return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, holder, sdk.NewCoins(targetCoin))
	}
	if err := k.bankKeeper.MintCoins(ctx, k.moduleName, sdk.NewCoins(targetCoin)); err != nil {
		return err
	}
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, holder, sdk.NewCoins(targetCoin))
}
