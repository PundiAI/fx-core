package keeper

import (
	"context"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	gravitytypes "github.com/functionx/fx-core/v3/x/gravity/types"
)

var _ types.MsgServer = &Keeper{}

// ConvertCoin converts Cosmos-native Coins into ERC20 tokens for both
// Cosmos-native and ERC20 TokenPair Owners
func (k Keeper) ConvertCoin(goCtx context.Context, msg *types.MsgConvertCoin) (*types.MsgConvertCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Error checked during msg validation
	receiver := common.HexToAddress(msg.Receiver)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)

	pair, err := k.MintingEnabled(ctx, sender, receiver.Bytes(), msg.Coin.Denom)
	if err != nil {
		return nil, err
	}

	// Remove token pair if contract is suicided
	erc20 := common.HexToAddress(pair.Erc20Address)
	acc := k.evmKeeper.GetAccountWithoutBalance(ctx, erc20)
	if acc == nil || !acc.IsContract() {
		k.RemoveTokenPair(ctx, pair)
		k.Logger(ctx).Debug("deleting selfdestructed token pair from state", "contract", pair.Erc20Address)
		// NOTE: return nil error to persist the changes from the deletion
		return &types.MsgConvertCoinResponse{}, nil
	}

	newCtx := ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	// Check ownership and execute conversion
	switch {
	case pair.IsNativeCoin():
		return k.ConvertCoinNativeCoin(newCtx, pair, msg, receiver, sender) // case 1.1
	case pair.IsNativeERC20():
		return k.ConvertCoinNativeERC20(newCtx, pair, msg, receiver, sender) // case 2.2
	default:
		return nil, types.ErrUndefinedOwner
	}
}

// ConvertERC20 converts ERC20 tokens into Cosmos-native Coins for both
// Cosmos-native and ERC20 TokenPair Owners
func (k Keeper) ConvertERC20(goCtx context.Context, msg *types.MsgConvertERC20) (*types.MsgConvertERC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Error checked during msg validation
	receiver := sdk.MustAccAddressFromBech32(msg.Receiver)
	/* sender have two cases
	* 1. cosmos - m/44'/118'/0'/0/0
	* private key ---> cosmos address ---> hex sender
	* hex sender convert from cosmos address
	* private key can authenticate cosmos address in fxcore (verify with cosmos_spec256k1)
	* private key can **NOT** authenticate hex sender address in EVM (verify with eth_spec256k1)
	*
	* 2. eth - m/44'/60'/0'/0/1
	* private key ---> eth public key ---> hex sender ---> cosmos address
	*						 | ---> cosmos public key
	* cosmos address equal to hex sender, generate by eth public key
	* private key can authenticate cosmos address in fxcore (verify with eth_spec256k1)
	* private key can authenticate hex sender address in EVM (verify with eth_spec256k1)
	 */
	sender := common.HexToAddress(msg.Sender)

	pair, err := k.MintingEnabled(ctx, sender.Bytes(), receiver, msg.ContractAddress)
	if err != nil {
		return nil, err
	}

	// Remove token pair if contract is suicided
	erc20 := common.HexToAddress(pair.Erc20Address)
	acc := k.evmKeeper.GetAccountWithoutBalance(ctx, erc20)
	if acc == nil || !acc.IsContract() {
		k.RemoveTokenPair(ctx, pair)
		k.Logger(ctx).Debug("deleting selfdestructed token pair from state", "contract", pair.Erc20Address)
		// NOTE: return nil error to persist the changes from the deletion
		return &types.MsgConvertERC20Response{}, nil
	}

	newCtx := ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	// Check ownership
	switch {
	case pair.IsNativeCoin():
		return k.ConvertERC20NativeCoin(newCtx, pair, msg, receiver, sender) // case 1.2
	case pair.IsNativeERC20():
		return k.ConvertERC20NativeToken(newCtx, pair, msg, receiver, sender) // case 2.1
	default:
		return nil, types.ErrUndefinedOwner
	}
}

// ConvertDenom converts coin into other coin, use for multiple chains in the same currency
func (k Keeper) ConvertDenom(goCtx context.Context, msg *types.MsgConvertDenom) (*types.MsgConvertDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Error checked during msg validation
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	receiver := sdk.MustAccAddressFromBech32(msg.Receiver)

	var coin sdk.Coin
	var err error
	if len(msg.Target) > 0 {
		// convert one to many
		coin, err = k.ConvertDenomToMany(ctx, sender, msg.Coin, msg.Target)
	} else {
		coin, err = k.ConvertDenomToOne(ctx, sender, msg.Coin)
	}
	if err != nil {
		return nil, err
	}

	if coin.Denom == msg.Coin.Denom {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "denom %s not support", msg.Coin.Denom)
	}

	if !sender.Equals(receiver) {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(coin)); err != nil {
			return nil, err
		}
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(coin)); err != nil {
			return nil, err
		}
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tx", "msg", "convert", "denom", "total"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("denom", msg.Coin.Denom),
				telemetry.NewLabel("target_denom", coin.Denom),
			},
		)
	}()

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeConvertDenom,
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
				sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
			),
		},
	)

	return &types.MsgConvertDenomResponse{}, nil
}

// ConvertCoinNativeCoin handles the Coin conversion flow for a native coin
// token pair:
//   - Escrow Coins on module account (Coins are not burned)
//   - Mint Tokens and send to receiver
//   - Check if token balance increased by amount
func (k Keeper) ConvertCoinNativeCoin(ctx sdk.Context, pair types.TokenPair, msg *types.MsgConvertCoin, receiver common.Address, sender sdk.AccAddress) (*types.MsgConvertCoinResponse, error) {
	// NOTE: ignore validation from NewCoin constructor
	coins := sdk.Coins{msg.Coin}
	erc20 := fxtypes.GetERC20().ABI
	contract := pair.GetERC20Contract()

	// Escrow Coins on module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to escrow coins")
	}

	// Mint Tokens and send to receiver
	_, err := k.CallEVM(ctx, erc20, types.ModuleAddress, contract, true, "mint", receiver, msg.Coin.Amount.BigInt())
	if err != nil {
		return nil, err
	}

	if pair.Denom == fxtypes.DefaultDenom {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, contract.Bytes(), coins); err != nil {
			return nil, sdkerrors.Wrap(err, "failed to transfer escrow coins to origin denom")
		}
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tx", "msg", "convert", "coin", "total"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("denom", pair.Denom),
				telemetry.NewLabel("erc20", pair.Erc20Address),
			},
		)
	}()

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeConvertCoin,
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
				sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
				sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Coin.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyDenom, msg.Coin.Denom),
				sdk.NewAttribute(types.AttributeKeyTokenAddress, pair.Erc20Address),
			),
		},
	)

	return &types.MsgConvertCoinResponse{}, nil
}

// ConvertERC20NativeCoin handles the erc20 conversion flow for a native coin token pair:
//   - Burn escrowed tokens
//   - Unescrow coins that have been previously escrowed with ConvertCoin
//   - Check if coin balance increased by amount
//   - Check if token balance decreased by amount
func (k Keeper) ConvertERC20NativeCoin(ctx sdk.Context, pair types.TokenPair, msg *types.MsgConvertERC20, receiver sdk.AccAddress, sender common.Address) (*types.MsgConvertERC20Response, error) {
	// NOTE: coin fields already validated
	coins := sdk.Coins{sdk.Coin{Denom: pair.Denom, Amount: msg.Amount}}

	erc20 := fxtypes.GetERC20().ABI
	contract := pair.GetERC20Contract()

	// Burn escrowed tokens
	_, err := k.CallEVM(ctx, erc20, types.ModuleAddress, contract, true, "burn", sender, msg.Amount.BigInt())
	if err != nil {
		return nil, err
	}

	// Transfer origin denom to module
	if pair.Denom == fxtypes.DefaultDenom {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, contract.Bytes(), types.ModuleName, coins); err != nil {
			return nil, sdkerrors.Wrap(err, "failed to transfer origin denom to module")
		}
	}

	// Unescrow Coins and send to receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, coins); err != nil {
		return nil, err
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tx", "msg", "convert", "erc20", "total"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("denom", pair.Denom),
				telemetry.NewLabel("erc20", pair.Erc20Address),
			},
		)
	}()

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeConvertERC20,
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
				sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
				sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyDenom, pair.Denom),
				sdk.NewAttribute(types.AttributeKeyTokenAddress, msg.ContractAddress),
			),
		},
	)

	return &types.MsgConvertERC20Response{}, nil
}

// ConvertERC20NativeToken handles the erc20 conversion flow for a native erc20 token pair:
//   - Escrow tokens on module account (Don't burn as module is not contract owner)
//   - Mint coins on module
//   - Send minted coins to the receiver
//   - Check if coin balance increased by amount
//   - Check if token balance decreased by amount
//   - Check for unexpected `appove` event in logs
func (k Keeper) ConvertERC20NativeToken(ctx sdk.Context, pair types.TokenPair, msg *types.MsgConvertERC20, receiver sdk.AccAddress, sender common.Address) (*types.MsgConvertERC20Response, error) {
	// NOTE: coin fields already validated
	coins := sdk.Coins{sdk.Coin{Denom: pair.Denom, Amount: msg.Amount}}
	erc20 := fxtypes.GetERC20().ABI
	contract := pair.GetERC20Contract()

	// Escrow tokens on module account
	transferData, err := erc20.Pack("transfer", types.ModuleAddress, msg.Amount.BigInt())
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrABIPack, "failed to pack transfer: %s", err.Error())
	}

	res, err := k.evmKeeper.CallEVMWithData(ctx, sender, &contract, transferData, true)
	if err != nil {
		return nil, err
	}

	// Check unpackedRet execution
	var unpackedRet types.ERC20BoolResponse
	if err := erc20.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrABIUnpack, "failed to unpack transfer: %s", err.Error())
	}

	if !unpackedRet.Value {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "failed to execute transfer")
	}

	// Mint coins
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return nil, err
	}

	// Send minted coins to the receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, coins); err != nil {
		return nil, err
	}

	// Check for unexpected `appove` event in logs
	if err = k.monitorApprovalEvent(res); err != nil {
		return nil, err
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tx", "msg", "convert", "erc20", "total"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("coin", pair.Denom),
				telemetry.NewLabel("erc20", pair.Erc20Address),
			},
		)
	}()

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeConvertERC20,
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
				sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
				sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyDenom, pair.Denom),
				sdk.NewAttribute(types.AttributeKeyTokenAddress, msg.ContractAddress),
			),
		},
	)

	return &types.MsgConvertERC20Response{}, nil
}

// ConvertCoinNativeERC20 handles the Coin conversion flow for a native ERC20
// token pair:
//   - Escrow Coins on module account
//   - Unescrow Tokens that have been previously escrowed with ConvertERC20 and send to receiver
//   - Burn escrowed Coins
//   - Check if token balance increased by amount
//   - Check for unexpected `appove` event in logs
func (k Keeper) ConvertCoinNativeERC20(ctx sdk.Context, pair types.TokenPair, msg *types.MsgConvertCoin, receiver common.Address, sender sdk.AccAddress) (*types.MsgConvertCoinResponse, error) {
	// NOTE: ignore validation from NewCoin constructor
	coins := sdk.Coins{msg.Coin}

	erc20 := fxtypes.GetERC20().ABI
	contract := pair.GetERC20Contract()

	// Escrow Coins on module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to escrow coins")
	}

	// Unescrow Tokens and send to receiver
	res, err := k.CallEVM(ctx, erc20, types.ModuleAddress, contract, true, "transfer", receiver, msg.Coin.Amount.BigInt())
	if err != nil {
		return nil, err
	}

	// Check unpackedRet execution
	var unpackedRet types.ERC20BoolResponse
	if err := erc20.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrABIUnpack, "failed to unpack transfer: %s", err.Error())
	}

	if !unpackedRet.Value {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "failed to execute unescrow tokens from user")
	}

	// Burn escrowed Coins
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to burn coins")
	}

	// Check for unexpected `appove` event in logs
	if err = k.monitorApprovalEvent(res); err != nil {
		return nil, err
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tx", "msg", "convert", "coin", "total"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("denom", pair.Denom),
				telemetry.NewLabel("erc20", pair.Erc20Address),
			},
		)
	}()

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeConvertCoin,
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
				sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
				sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Coin.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyDenom, msg.Coin.Denom),
				sdk.NewAttribute(types.AttributeKeyTokenAddress, pair.Erc20Address),
			),
		},
	)

	return &types.MsgConvertCoinResponse{}, nil
}

// ConvertDenomToMany handles the Denom conversion flow for one to many
// token pair:
//   - Escrow Coins on module account
//   - Unescrow Tokens that have been previously escrowed with ConvertDenomToMany and send to receiver
//   - Burn escrowed Coins
//   - Check if token balance increased by amount
func (k Keeper) ConvertDenomToMany(ctx sdk.Context, from sdk.AccAddress, coin sdk.Coin, target string) (sdk.Coin, error) {
	// metadata has alias
	md, found := k.HasDenomAlias(ctx, coin.Denom)
	if !found {
		return coin, nil
	}

	// denom registered
	if !k.IsDenomRegistered(ctx, coin.Denom) {
		return sdk.Coin{}, sdkerrors.Wrapf(types.ErrInvalidDenom, "denom %s not registered", coin.Denom)
	}

	// convert target to denom prefix
	denomPrefix := target
	if target == gravitytypes.ModuleName {
		denomPrefix = ethtypes.ModuleName
	}

	aliases := md.DenomUnits[0].Aliases
	targetDenom := ""
	for _, alias := range aliases {
		if strings.HasPrefix(alias, denomPrefix) {
			targetDenom = alias
			break
		}
	}
	if len(targetDenom) == 0 {
		return sdk.Coin{}, sdkerrors.Wrapf(types.ErrInvalidTarget, "target %s denom not exist", target)
	}

	var err error
	targetCoin := sdk.NewCoin(targetDenom, coin.Amount)
	// send symbol denom to module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return sdk.Coin{}, err
	}
	// send alias denom to from addr
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, from, sdk.NewCoins(targetCoin)) //ibc0xxx
	if err != nil {
		return sdk.Coin{}, err
	}
	// burn symbol coin
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return sdk.Coin{}, err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeConvertDenomToMany,
				sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, coin.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyDenom, coin.Denom),
				sdk.NewAttribute(types.AttributeKeyTargetDenom, targetCoin.Denom),
			),
		},
	)

	return targetCoin, nil
}

// ConvertDenomToOne handles the Denom conversion flow for many to one
// token pair:
//   - Escrow Coins on module account (Coins are not burned)
//   - Mint Tokens and send to from address
//   - Check if token balance increased by amount
func (k Keeper) ConvertDenomToOne(ctx sdk.Context, from sdk.AccAddress, coin sdk.Coin) (sdk.Coin, error) {
	// denom not register
	if k.IsDenomRegistered(ctx, coin.Denom) {
		return sdk.Coin{}, sdkerrors.Wrapf(types.ErrInvalidDenom, "denom %s already registered", coin.Denom)
	}
	// alias register
	aliasDenomBytes := k.GetAliasDenom(ctx, coin.Denom)
	if len(aliasDenomBytes) == 0 {
		return sdk.Coin{}, sdkerrors.Wrapf(types.ErrInvalidDenom, "alias %s not registered", coin.Denom)
	}
	// has denom alias
	if _, found := k.HasDenomAlias(ctx, string(aliasDenomBytes)); !found {
		return sdk.Coin{}, sdkerrors.Wrapf(types.ErrInvalidMetadata, "not support with %s", string(aliasDenomBytes))
	}

	var err error
	targetCoin := sdk.NewCoin(string(aliasDenomBytes), coin.Amount)
	// send alias denom to module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return sdk.Coin{}, err
	}
	//mint symbol denom
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(targetCoin))
	if err != nil {
		return sdk.Coin{}, err
	}
	//send symbol denom to from addr
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, from, sdk.NewCoins(targetCoin))
	if err != nil {
		return sdk.Coin{}, err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeConvertDenomToOne,
				sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, coin.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyDenom, coin.Denom),
				sdk.NewAttribute(types.AttributeKeyTargetDenom, targetCoin.Denom),
			),
		},
	)

	return targetCoin, nil
}
