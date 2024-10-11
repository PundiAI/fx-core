package keeper

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
)

func (k Keeper) EvmToBaseCoin(ctx context.Context, tokenAddr string, amount *big.Int, holder common.Address) (sdk.Coin, error) {
	_, err := k.erc20Keeper.ConvertERC20(ctx, &erc20types.MsgConvertERC20{
		ContractAddress: tokenAddr,
		Amount:          sdkmath.NewIntFromBigInt(amount),
		Receiver:        sdk.AccAddress(holder.Bytes()).String(),
		Sender:          holder.String(),
	})
	if err != nil {
		return sdk.Coin{}, err
	}
	tokenPair, ok := k.erc20Keeper.GetTokenPair(sdk.UnwrapSDKContext(ctx), tokenAddr)
	if !ok {
		return sdk.Coin{}, types.ErrInvalid.Wrapf("not found %s token pair", tokenAddr)
	}
	return sdk.NewCoin(tokenPair.Denom, sdkmath.NewIntFromBigInt(amount)), nil
}

func (k Keeper) BaseCoinToEvm(ctx context.Context, coin sdk.Coin, holder common.Address) (string, error) {
	_, err := k.erc20Keeper.ConvertCoin(ctx, &erc20types.MsgConvertCoin{
		Coin:     coin,
		Receiver: holder.String(),
		Sender:   sdk.AccAddress(holder.Bytes()).String(),
	})
	if err != nil {
		return "", err
	}
	tokenPair, ok := k.erc20Keeper.GetTokenPair(sdk.UnwrapSDKContext(ctx), coin.Denom)
	if !ok {
		return "", types.ErrInvalid.Wrapf("not found %s token pair", coin.Denom)
	}
	return tokenPair.Erc20Address, nil
}

func (k Keeper) BridgeTokenToBaseCoin(ctx context.Context, tokenAddr string, amount sdkmath.Int, holder sdk.AccAddress) (sdk.Coin, error) {
	bridgeDenom, found := k.GetBridgeDenomByContract(sdk.UnwrapSDKContext(ctx), tokenAddr)
	if !found {
		return sdk.Coin{}, sdkerrors.ErrInvalidCoins.Wrapf("bridge denom not found %s", tokenAddr)
	}
	bridgeToken := sdk.NewCoin(bridgeDenom, amount)
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

func (k Keeper) BaseCoinToBridgeToken(ctx context.Context, coin sdk.Coin, holder sdk.AccAddress) (string, error) {
	bridgeDenom, err := k.ManyToOne(ctx, coin.Denom, k.moduleName)
	if err != nil {
		return "", err
	}
	if err = k.ConversionCoin(ctx, holder, coin, coin.Denom, bridgeDenom); err != nil {
		return "", err
	}
	if err = k.WithdrawBridgeToken(ctx, sdk.NewCoin(bridgeDenom, coin.Amount), holder); err != nil {
		return "", err
	}
	tokenAddr, found := k.GetContractByBridgeDenom(sdk.UnwrapSDKContext(ctx), bridgeDenom)
	if !found {
		return "", sdkerrors.ErrInvalidRequest.Wrapf("bridge token not found %s", bridgeDenom)
	}
	return tokenAddr, nil
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
		return sdkerrors.ErrInvalidCoins.Wrapf("token pair not found: %s", baseDenom)
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
		return sdkerrors.ErrInvalidCoins.Wrapf("token pair not found: %s", baseDenom)
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
	bridgeDenom, err := k.GetBridgeDenoms(ctx, baseDenom)
	if err != nil {
		return "", err
	}
	ibcPrefix := ibctransfertypes.DenomPrefix + "/"
	fxTarget := fxtypes.ParseFxTarget(target)
	for _, bd := range bridgeDenom {
		if strings.HasPrefix(bd, ibcPrefix) && fxTarget.IsIBC() {
			hexHash := strings.TrimPrefix(bd, ibctransfertypes.DenomPrefix+"/")
			hash, err := ibctransfertypes.ParseHexHash(hexHash)
			if err != nil {
				return "", err
			}
			denomTrace, found := k.ibcTransferKeeper.GetDenomTrace(sdk.UnwrapSDKContext(ctx), hash)
			if !found {
				continue
			}
			if !strings.HasPrefix(denomTrace.GetPath(), fmt.Sprintf("%s/%s", fxTarget.SourcePort, fxTarget.SourceChannel)) {
				continue
			}
			return bd, nil
		}
		if strings.HasPrefix(bd, fxTarget.GetTarget()) {
			return bd, nil
		}
	}
	return "", sdkerrors.ErrInvalidCoins.Wrapf("not found bridge denom: %s, %s", baseDenom, target)
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
		return sdkerrors.ErrInvalidCoins.Wrapf("token pair not found %s", baseDenom)
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

func (k Keeper) IBCCoinToBaseCoin(ctx context.Context, coin sdk.Coin, holder sdk.AccAddress) (sdk.Coin, error) {
	if !strings.HasPrefix(coin.Denom, ibctransfertypes.DenomPrefix+"/") {
		return coin, nil
	}
	baseDenom, err := k.ManyToOne(ctx, coin.Denom)
	if err != nil {
		return sdk.Coin{}, err
	}
	baseCoin := sdk.NewCoin(baseDenom, coin.Amount)
	if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, holder, ibctransfertypes.ModuleName, sdk.NewCoins(coin)); err != nil {
		return sdk.Coin{}, err
	}
	if err = k.bankKeeper.MintCoins(ctx, ibctransfertypes.ModuleName, sdk.NewCoins(baseCoin)); err != nil {
		return sdk.Coin{}, err
	}
	if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, ibctransfertypes.ModuleName, holder, sdk.NewCoins(baseCoin)); err != nil {
		return sdk.Coin{}, err
	}
	return baseCoin, nil
}

func (k Keeper) BaseCoinToIBCCoin(ctx context.Context, coin sdk.Coin, holder sdk.AccAddress, ibcTarget string) (sdk.Coin, error) {
	if strings.HasPrefix(coin.Denom, ibctransfertypes.DenomPrefix+"/") {
		return coin, sdkerrors.ErrInvalidCoins.Wrapf("can not convert ibc denom")
	}
	ibcDenom, err := k.ManyToOne(ctx, coin.Denom, ibcTarget)
	if err != nil {
		return sdk.Coin{}, err
	}
	ibcCoin := sdk.NewCoin(ibcDenom, coin.Amount)
	if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, holder, ibctransfertypes.ModuleName, sdk.NewCoins(coin)); err != nil {
		return sdk.Coin{}, err
	}
	if err = k.bankKeeper.BurnCoins(ctx, ibctransfertypes.ModuleName, sdk.NewCoins(coin)); err != nil {
		return sdk.Coin{}, err
	}
	if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, ibctransfertypes.ModuleName, holder, sdk.NewCoins(ibcCoin)); err != nil {
		return sdk.Coin{}, err
	}
	return ibcCoin, nil
}

func (k Keeper) IBCCoinToEvm(ctx context.Context, coin sdk.Coin, holder sdk.AccAddress) error {
	baseCoin, err := k.IBCCoinToBaseCoin(ctx, coin, holder)
	if err != nil {
		return err
	}
	_, err = k.erc20Keeper.ConvertCoin(ctx, &erc20types.MsgConvertCoin{
		Coin:     baseCoin,
		Receiver: common.BytesToAddress(holder).String(),
		Sender:   holder.String(),
	})
	return err
}

func (k Keeper) IBCCoinRefund(ctx sdk.Context, coin sdk.Coin, holder sdk.AccAddress, ibcChannel string, ibcSequence uint64) error {
	baseCoin, err := k.IBCCoinToBaseCoin(ctx, coin, holder)
	if err != nil {
		return err
	}
	if !k.erc20Keeper.DeleteIBCTransferRelation(ctx, ibcChannel, ibcSequence) {
		return nil
	}
	_, err = k.BaseCoinToEvm(ctx, baseCoin, common.BytesToAddress(holder.Bytes()))
	return err
}

func (k Keeper) AfterIBCAckSuccess(ctx sdk.Context, sourceChannel string, sequence uint64) {
	k.erc20Keeper.DeleteOutgoingTransferRelation(ctx, sourceChannel, sequence)
}
