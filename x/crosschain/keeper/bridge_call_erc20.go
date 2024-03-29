package keeper

import (
	"fmt"
	"math/big"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v7/contract"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
)

func (k Keeper) BridgeCallERC20Handler(
	ctx sdk.Context,
	asset []byte,
	sender common.Address,
	to *common.Address,
	receiver sdk.AccAddress,
	dstChainID string,
	message []byte,
	value sdkmath.Int,
	gasLimit, eventNonce uint64,
) error {
	tokenAddrs, amounts, err := contract.UnpackERC20Asset(asset)
	if err != nil {
		return errorsmod.Wrap(types.ErrInvalid, err.Error())
	}
	tokens, err := types.NewERC20Tokens(k.moduleName, tokenAddrs, amounts)
	if err != nil {
		return errorsmod.Wrap(types.ErrInvalid, err.Error())
	}

	var errCause string
	if dstChainID == types.FxcoreChainID {
		cacheCtx, commit := ctx.CacheContext()
		err = k.bridgeCallFxCore(cacheCtx, sender, tokens, receiver, message, to, value, gasLimit, eventNonce)
		if err != nil {
			errCause = err.Error()
		} else {
			commit()
		}
	}
	if len(errCause) > 0 && len(tokens) > 0 {
		receiverStr := fxtypes.AddressToStr(sender.Bytes(), k.moduleName)
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(types.EventTypeBridgeCallRefund,
				sdk.NewAttribute(types.AttributeKeyErrCause, errCause),
				sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(eventNonce)),
				sdk.NewAttribute(types.AttributeKeyRefundAddress, receiverStr),
			),
		})
		return k.AddRefundRecord(ctx, receiverStr, eventNonce, tokens)
	}
	return nil
}

func (k Keeper) bridgeCallFxCore(ctx sdk.Context, sender common.Address, tokens []types.ERC20Token, receiver sdk.AccAddress, message []byte, to *common.Address, value sdkmath.Int, gasLimit uint64, eventNonce uint64) error {
	coins, err := k.bridgeCallTransferToSender(ctx, sender.Bytes(), tokens)
	if err != nil {
		return err
	}
	if err = k.bridgeCallTransferToReceiver(ctx, sender.Bytes(), receiver, coins); err != nil {
		return err
	}
	if len(message) > 0 || to != nil {
		res, err := k.bridgeCallEvmHandler(ctx, sender, to, message, value, gasLimit, eventNonce)
		if err != nil {
			return err
		}
		if res.Failed() {
			return errorsmod.Wrap(types.ErrInvalid, res.VmError)
		}
	}
	return nil
}

func (k Keeper) bridgeCallTransferToSender(ctx sdk.Context, receiver sdk.AccAddress, tokens []types.ERC20Token) (sdk.Coins, error) {
	mintCoins := sdk.NewCoins()
	unlockCoins := sdk.NewCoins()
	for i := 0; i < len(tokens); i++ {
		bridgeToken := k.GetBridgeTokenDenom(ctx, tokens[i].Contract)
		if bridgeToken == nil {
			return nil, errorsmod.Wrap(types.ErrInvalid, "bridge token is not exist")
		}
		if !tokens[i].Amount.IsPositive() {
			continue
		}
		coin := sdk.NewCoin(bridgeToken.Denom, tokens[i].Amount)
		isOriginOrConverted := k.erc20Keeper.IsOriginOrConvertedDenom(ctx, bridgeToken.Denom)
		if !isOriginOrConverted {
			mintCoins = mintCoins.Add(coin)
		}
		unlockCoins = unlockCoins.Add(coin)
	}
	if mintCoins.IsAllPositive() {
		if err := k.bankKeeper.MintCoins(ctx, k.moduleName, mintCoins); err != nil {
			return nil, errorsmod.Wrapf(err, "mint vouchers coins")
		}
	}
	if unlockCoins.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, receiver, unlockCoins); err != nil {
			return nil, errorsmod.Wrap(err, "transfer vouchers")
		}
	}

	targetCoins := sdk.NewCoins()
	for _, coin := range unlockCoins {
		targetCoin, err := k.erc20Keeper.ConvertDenomToTarget(ctx, receiver, coin, fxtypes.ParseFxTarget(fxtypes.ERC20Target))
		if err != nil {
			return nil, errorsmod.Wrap(err, "convert to target coin")
		}
		targetCoins = targetCoins.Add(targetCoin)
	}
	return targetCoins, nil
}

func (k Keeper) bridgeCallTransferToReceiver(ctx sdk.Context, sender sdk.AccAddress, receiver []byte, coins sdk.Coins) error {
	for _, coin := range coins {
		if coin.Denom == fxtypes.DefaultDenom {
			if err := k.bankKeeper.SendCoins(ctx, sender, receiver, sdk.NewCoins(coin)); err != nil {
				return err
			}
			continue
		}
		if _, err := k.erc20Keeper.ConvertCoin(sdk.WrapSDKContext(ctx), &erc20types.MsgConvertCoin{
			Coin:     coin,
			Receiver: common.BytesToAddress(receiver).String(),
			Sender:   sender.String(),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) bridgeCallEvmHandler(ctx sdk.Context, sender common.Address, to *common.Address, message []byte, value sdkmath.Int, gasLimit, eventNonce uint64) (*evmtypes.MsgEthereumTxResponse, error) {
	callErr, callResult := "", false
	defer func() {
		attrs := []sdk.Attribute{
			sdk.NewAttribute(types.AttributeKeyEventNonce, strconv.FormatUint(eventNonce, 10)),
			sdk.NewAttribute(types.AttributeKeyEvmCallResult, strconv.FormatBool(callResult)),
		}
		if len(callErr) > 0 {
			attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyEvmCallError, callErr))
		}
		ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(types.EventTypeBridgeCallEvm, attrs...)})
	}()

	txResp, err := k.evmKeeper.CallEVM(ctx, sender, to, value.BigInt(), gasLimit, message, true)
	if err != nil {
		callErr = err.Error()
		return nil, err
	}

	callResult = !txResp.Failed()
	callErr = txResp.VmError
	return txResp, nil
}

func (k Keeper) bridgeCallAssetHandler(ctx sdk.Context, sender sdk.AccAddress, asset string) (string, error) {
	assetType, assetData, err := types.UnpackAssetType(asset)
	if err != nil {
		return "", errorsmod.Wrap(types.ErrInvalid, "asset")
	}

	switch assetType {
	case contract.AssetERC20:
		tokenAddrs, amounts, err := contract.UnpackERC20Asset(assetData)
		if err != nil {
			return "", errorsmod.Wrap(types.ErrInvalid, "erc20 token")
		}
		tokens, err := types.NewERC20Tokens(fxtypes.EthTarget, tokenAddrs, amounts)
		if err != nil {
			return "", errorsmod.Wrap(types.ErrInvalid, err.Error())
		}
		coins, err := k.bridgeCallERC20ToCoins(ctx, sender, tokens)
		if err != nil {
			return "", err
		}
		if coins.IsZero() {
			return "", nil
		}
		return k.BridgeTokenToAsset(ctx, sender, coins)
	default:
		return "", errorsmod.Wrap(types.ErrInvalid, "asset type")
	}
}

func (k Keeper) bridgeCallERC20ToCoins(ctx sdk.Context, sender sdk.AccAddress, tokens []types.ERC20Token) (sdk.Coins, error) {
	coins := sdk.NewCoins()
	for _, t := range tokens {
		if t.Contract == contract.EmptyEvmAddress {
			coins = coins.Add(sdk.NewCoin(fxtypes.DefaultDenom, t.Amount))
			continue
		}
		if _, err := k.erc20Keeper.ConvertERC20(sdk.WrapSDKContext(ctx), &erc20types.MsgConvertERC20{
			ContractAddress: t.Contract,
			Amount:          t.Amount,
			Receiver:        sender.String(),
			Sender:          common.BytesToAddress(sender.Bytes()).String(),
		}); err != nil {
			return nil, err
		}
		tokenPair, found := k.erc20Keeper.GetTokenPair(ctx, t.Contract)
		if !found {
			return nil, errorsmod.Wrap(erc20types.ErrTokenPairNotFound, t.Contract)
		}
		targetCoin, err := k.erc20Keeper.ConvertDenomToTarget(ctx, sender, sdk.NewCoin(tokenPair.Denom, t.Amount), fxtypes.ParseFxTarget(k.moduleName))
		if err != nil {
			return nil, err
		}
		coins = coins.Add(targetCoin)
	}
	return coins, nil
}

func (k Keeper) BridgeTokenToAsset(ctx sdk.Context, sender sdk.AccAddress, coins sdk.Coins) (string, error) {
	tokenAddrs := make([]common.Address, 0, len(coins))
	amounts := make([]*big.Int, 0, len(coins))
	coinsNeedBurn := sdk.NewCoins()
	for _, coin := range coins {
		bridgeToken := k.GetDenomBridgeToken(ctx, coin.Denom)
		if bridgeToken == nil {
			return "", errorsmod.Wrap(errortypes.ErrNotFound, "bridge token")
		}
		tokenAddr := types.ExternalAddressToAccAddress(k.moduleName, bridgeToken.Token)
		tokenAddrs = append(tokenAddrs, common.BytesToAddress(tokenAddr.Bytes()))
		amounts = append(amounts, coin.Amount.BigInt())
		if k.erc20Keeper.IsOriginOrConvertedDenom(ctx, coin.Denom) {
			coinsNeedBurn = coinsNeedBurn.Add(coin)
		}
	}
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, k.moduleName, coins); err != nil {
		return "", err
	}
	if err := k.bankKeeper.BurnCoins(ctx, k.moduleName, coinsNeedBurn); err != nil {
		return "", err
	}
	return contract.PackERC20AssetWithType(tokenAddrs, amounts)
}
