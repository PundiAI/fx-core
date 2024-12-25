package keeper

import (
	"context"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (k Keeper) BridgeTokenToBaseCoin(ctx context.Context, holder sdk.AccAddress, amount sdkmath.Int, bridgeToken erc20types.BridgeToken) (baseCoin sdk.Coin, err error) {
	baseCoin = sdk.NewCoin(bridgeToken.Denom, amount)
	if bridgeToken.IsOrigin() {
		return baseCoin, nil
	}
	bridgeCoins := sdk.NewCoins(sdk.NewCoin(bridgeToken.BridgeDenom(), amount))
	if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, holder, types.ModuleName, bridgeCoins); err != nil {
		return sdk.Coin{}, err
	}
	baseCoins := sdk.NewCoins(baseCoin)
	if bridgeToken.IsNative {
		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, bridgeCoins)
	} else {
		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, baseCoins)
	}
	if err != nil {
		return sdk.Coin{}, err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, holder, baseCoins)
	return baseCoin, err
}

func (k Keeper) BaseCoinToBridgeToken(ctx context.Context, holder sdk.AccAddress, baseCoin sdk.Coin) (bridgeToken erc20types.BridgeToken, err error) {
	bridgeToken, err = k.erc20Keeper.GetBridgeToken(ctx, k.moduleName, baseCoin.Denom)
	if err != nil {
		return bridgeToken, err
	}

	baseCoins := sdk.NewCoins(baseCoin)
	if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, holder, types.ModuleName, baseCoins); err != nil {
		return bridgeToken, err
	}
	bridgeCoins := sdk.NewCoins(sdk.NewCoin(bridgeToken.BridgeDenom(), baseCoin.Amount))
	if bridgeToken.IsNative {
		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, bridgeCoins)
	} else {
		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, baseCoins)
	}
	if err != nil {
		return bridgeToken, err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, holder, bridgeCoins)
	return bridgeToken, err
}

// DepositBridgeToken get bridge token from k.moduleName
func (k Keeper) DepositBridgeToken(ctx context.Context, holder sdk.AccAddress, amount sdkmath.Int, tokenAddr string) (bridgeToken erc20types.BridgeToken, err error) {
	bridgeToken, err = k.GetBridgeToken(ctx, tokenAddr)
	if err != nil {
		return bridgeToken, err
	}

	bridgeCoins := sdk.NewCoins(sdk.NewCoin(bridgeToken.BridgeDenom(), amount))
	if bridgeToken.IsOrigin() {
		return bridgeToken, k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, holder, bridgeCoins)
	}

	if !bridgeToken.IsNative {
		if err = k.bankKeeper.MintCoins(ctx, k.moduleName, bridgeCoins); err != nil {
			return bridgeToken, err
		}
	}
	return bridgeToken, k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, holder, bridgeCoins)
}

// WithdrawBridgeToken put bridge token to k.moduleName
func (k Keeper) WithdrawBridgeToken(ctx context.Context, holder sdk.AccAddress, amount sdkmath.Int, bridgeToken erc20types.BridgeToken) error {
	bridgeCoins := sdk.NewCoins(sdk.NewCoin(bridgeToken.BridgeDenom(), amount))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, holder, k.moduleName, bridgeCoins); err != nil {
		return err
	}
	if bridgeToken.IsOrigin() || bridgeToken.IsNative {
		return nil
	}
	return k.bankKeeper.BurnCoins(ctx, k.moduleName, bridgeCoins)
}

func (k Keeper) IBCCoinToBaseCoin(ctx context.Context, holder sdk.AccAddress, ibcCoin sdk.Coin) (string, error) {
	isNative := !strings.HasPrefix(ibcCoin.Denom, ibctransfertypes.DenomPrefix+"/")
	if isNative {
		return ibcCoin.Denom, nil
	}
	baseDenom, err := k.erc20Keeper.GetBaseDenom(ctx, ibcCoin.Denom)
	if err != nil {
		// NOTE: if not found in IBCToken
		return "", nil
	}
	if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, holder, types.ModuleName, sdk.NewCoins(ibcCoin)); err != nil {
		return "", err
	}
	baseCoins := sdk.NewCoins(sdk.NewCoin(baseDenom, ibcCoin.Amount))
	if err = k.bankKeeper.MintCoins(ctx, types.ModuleName, baseCoins); err != nil {
		return "", err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, holder, baseCoins)
	return baseDenom, err
}

func (k Keeper) BaseCoinToIBCCoin(ctx context.Context, holder sdk.AccAddress, baseCoin sdk.Coin, channel string) (sdk.Coin, error) {
	ibcToken, err := k.erc20Keeper.GetIBCToken(ctx, channel, baseCoin.Denom)
	if err != nil {
		// NOTE: if not found in IBCToken
		return baseCoin, nil
	}
	baseCoins := sdk.NewCoins(baseCoin)
	if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, holder, types.ModuleName, baseCoins); err != nil {
		return sdk.Coin{}, err
	}
	if err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, baseCoins); err != nil {
		return sdk.Coin{}, err
	}
	ibcCoin := sdk.NewCoin(ibcToken.IbcDenom, baseCoin.Amount)
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, holder, sdk.NewCoins(ibcCoin))
	return ibcCoin, err
}

func (k Keeper) IBCCoinToEvm(ctx sdk.Context, holder sdk.AccAddress, ibcCoin sdk.Coin) error {
	baseDenom, err := k.IBCCoinToBaseCoin(ctx, holder, ibcCoin)
	if err != nil {
		return err
	}
	if baseDenom == "" {
		return nil
	}
	_, err = k.erc20Keeper.BaseCoinToEvm(ctx, common.BytesToAddress(holder.Bytes()), sdk.NewCoin(baseDenom, ibcCoin.Amount))
	return err
}

func (k Keeper) IBCCoinRefund(ctx sdk.Context, holder sdk.AccAddress, ibcCoin sdk.Coin, ibcChannel string, ibcSequence uint64) error {
	baseDenom, err := k.IBCCoinToBaseCoin(ctx, holder, ibcCoin)
	if err != nil {
		return err
	}
	if baseDenom == "" {
		return nil
	}
	ibcTransferKey := types.NewIBCTransferKey(ibcChannel, ibcSequence)
	found, err := k.erc20Keeper.HasCache(ctx, ibcTransferKey)
	if err != nil {
		return err
	}
	if found {
		return k.erc20Keeper.DeleteCache(ctx, ibcTransferKey)
	}
	_, err = k.erc20Keeper.BaseCoinToEvm(ctx, common.BytesToAddress(holder.Bytes()), sdk.NewCoin(baseDenom, ibcCoin.Amount))
	return err
}

func (k Keeper) AfterIBCAckSuccess(ctx sdk.Context, ibcChannel string, ibcSequence uint64) error {
	ibcTransferKey := types.NewIBCTransferKey(ibcChannel, ibcSequence)
	return k.erc20Keeper.DeleteCache(ctx, ibcTransferKey)
}

func (k Keeper) GetBridgeToken(ctx context.Context, tokenAddr string) (erc20types.BridgeToken, error) {
	baseDenom, err := k.erc20Keeper.GetBaseDenom(ctx, types.NewBridgeDenom(k.moduleName, tokenAddr))
	if err != nil {
		return erc20types.BridgeToken{}, err
	}
	return k.erc20Keeper.GetBridgeToken(ctx, k.moduleName, baseDenom)
}
