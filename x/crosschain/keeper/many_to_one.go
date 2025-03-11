package keeper

import (
	"context"
	"strings"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func (k Keeper) SwapBridgeToken(ctx context.Context, holder sdk.AccAddress, bridgeToken erc20types.BridgeToken, amount sdkmath.Int) (erc20types.BridgeToken, sdkmath.Int, error) {
	if bridgeToken.Denom != fxtypes.FXDenom {
		return bridgeToken, amount, nil
	}
	defBridgeToken, err := k.erc20Keeper.GetBridgeToken(ctx, k.moduleName, fxtypes.DefaultDenom)
	if err != nil {
		return erc20types.BridgeToken{}, sdkmath.Int{}, err
	}
	// transfer bridgeDenom from holder to module
	bridgeDenom := bridgeToken.BridgeDenom()
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, holder, k.moduleName, sdk.NewCoins(sdk.NewCoin(bridgeDenom, amount)))
	if err != nil {
		return erc20types.BridgeToken{}, sdkmath.Int{}, err
	}
	swapAmount := fxtypes.SwapAmount(amount)
	if !swapAmount.IsPositive() {
		return defBridgeToken, swapAmount, nil
	}
	// transfer defaultDenom from eth module to holder
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, ethtypes.ModuleName, holder, sdk.NewCoins(sdk.NewCoin(defBridgeToken.BridgeDenom(), swapAmount)))
	return defBridgeToken, swapAmount, err
}

func (k Keeper) DepositBridgeTokenToBaseCoin(ctx context.Context, holder sdk.AccAddress, amount sdkmath.Int, tokenAddr string) (sdk.Coin, error) {
	bridgeToken, err := k.DepositBridgeToken(ctx, holder, amount, tokenAddr)
	if err != nil {
		return sdk.Coin{}, err
	}
	bridgeToken, amount, err = k.SwapBridgeToken(ctx, holder, bridgeToken, amount)
	if err != nil {
		return sdk.Coin{}, err
	}
	baseCoin, err := k.BridgeTokenToBaseCoin(ctx, holder, amount, bridgeToken)
	return baseCoin, err
}

func (k Keeper) BridgeTokenToBaseCoin(ctx context.Context, holder sdk.AccAddress, amount sdkmath.Int, bridgeToken erc20types.BridgeToken) (baseCoin sdk.Coin, err error) {
	if !amount.IsPositive() {
		return sdk.Coin{}, err
	}
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
	if bridgeToken.IsOrigin() {
		return bridgeToken, nil
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
	bridgeToken, err = k.GetBridgeToken(ctx, types.NewBridgeDenom(k.moduleName, tokenAddr))
	if err != nil {
		return bridgeToken, err
	}

	bridgeCoins := sdk.NewCoins(sdk.NewCoin(bridgeToken.BridgeDenom(), amount))
	if bridgeToken.IsOrigin() {
		return bridgeToken, k.bankKeeper.SendCoinsFromModuleToAccount(ctx, ethtypes.ModuleName, holder, bridgeCoins)
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
	if bridgeToken.IsOrigin() {
		return k.bankKeeper.SendCoinsFromAccountToModule(ctx, holder, ethtypes.ModuleName, bridgeCoins)
	}
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, holder, k.moduleName, bridgeCoins); err != nil {
		return err
	}
	if bridgeToken.IsNative {
		return nil
	}
	return k.bankKeeper.BurnCoins(ctx, k.moduleName, bridgeCoins)
}

func (k Keeper) IBCCoinToBaseCoin(ctx context.Context, holder sdk.AccAddress, ibcCoin sdk.Coin) (foundBase bool, baseDenom string, err error) {
	isNative := !strings.HasPrefix(ibcCoin.Denom, ibctransfertypes.DenomPrefix+"/")
	if isNative {
		return true, ibcCoin.Denom, nil
	}
	baseDenom, err = k.erc20Keeper.GetBaseDenom(ctx, ibcCoin.Denom)
	if err != nil {
		// NOTE: if not found in IBCToken
		if errors.IsOf(err, collections.ErrNotFound) {
			return false, "", nil
		}
		return false, "", err
	}
	if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, holder, types.ModuleName, sdk.NewCoins(ibcCoin)); err != nil {
		return false, "", err
	}
	baseCoins := sdk.NewCoins(sdk.NewCoin(baseDenom, ibcCoin.Amount))
	if err = k.bankKeeper.MintCoins(ctx, types.ModuleName, baseCoins); err != nil {
		return false, "", err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, holder, baseCoins)
	return true, baseDenom, err
}

func (k Keeper) BaseCoinToIBCCoin(ctx context.Context, holder sdk.AccAddress, baseCoin sdk.Coin, channel string) (sdk.Coin, error) {
	ibcToken, err := k.erc20Keeper.GetIBCToken(ctx, baseCoin.Denom, channel)
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

func (k Keeper) IBCCoinToEvm(ctx sdk.Context, holderAddr string, ibcCoin sdk.Coin) error {
	if ibcCoin.GetDenom() == fxtypes.DefaultDenom {
		return nil
	}
	// parse receive address, compatible with evm addresses
	holder, isEvmAddr, err := fxtypes.ParseAddress(holderAddr)
	if err != nil {
		return err
	}

	if !isEvmAddr && !IsEthSecp256k1(k.ak.GetAccount(ctx, holder)) {
		return sdkerrors.ErrInvalidAddress.Wrap("only support hex address")
	}
	found, baseDenom, err := k.IBCCoinToBaseCoin(ctx, holder, ibcCoin)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	_, err = k.erc20Keeper.BaseCoinToEvm(ctx, k.evmKeeper, common.BytesToAddress(holder.Bytes()), sdk.NewCoin(baseDenom, ibcCoin.Amount))
	return err
}

func (k Keeper) IBCCoinRefund(ctx sdk.Context, holder sdk.AccAddress, ibcCoin sdk.Coin, ibcChannel string, ibcSequence uint64) error {
	found, baseDenom, err := k.IBCCoinToBaseCoin(ctx, holder, ibcCoin)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	ibcTransferKey := types.NewIBCTransferKey(ibcChannel, ibcSequence)
	found, err = k.erc20Keeper.HasCache(ctx, ibcTransferKey)
	if err != nil {
		return err
	}
	if found {
		return k.erc20Keeper.DeleteCache(ctx, ibcTransferKey)
	}
	_, err = k.erc20Keeper.BaseCoinToEvm(ctx, k.evmKeeper, common.BytesToAddress(holder.Bytes()), sdk.NewCoin(baseDenom, ibcCoin.Amount))
	return err
}

func (k Keeper) AfterIBCAckSuccess(ctx sdk.Context, ibcChannel string, ibcSequence uint64) error {
	ibcTransferKey := types.NewIBCTransferKey(ibcChannel, ibcSequence)
	return k.erc20Keeper.DeleteCache(ctx, ibcTransferKey)
}

func (k Keeper) GetBridgeToken(ctx context.Context, bridgeDenom string) (erc20types.BridgeToken, error) {
	baseDenom, err := k.erc20Keeper.GetBaseDenom(ctx, bridgeDenom)
	if err != nil {
		return erc20types.BridgeToken{}, err
	}
	return k.erc20Keeper.GetBridgeToken(ctx, k.moduleName, baseDenom)
}

func IsEthSecp256k1(account sdk.AccountI) bool {
	if account == nil {
		return false
	}
	if account.GetPubKey() == nil && account.GetSequence() > 0 {
		return true
	}

	if account.GetPubKey() != nil && account.GetPubKey().Type() == new(ethsecp256k1.PubKey).Type() {
		return true
	}
	return false
}
