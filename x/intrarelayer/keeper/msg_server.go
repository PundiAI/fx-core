package keeper

import (
	"bytes"
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/functionx/fx-core/x/intrarelayer/types"
	"github.com/functionx/fx-core/x/intrarelayer/types/contracts"
)

var _ types.MsgServer = &Keeper{}

// ConvertCoin converts Cosmos-native Coins into FIP20 tokens for both
// Cosmos-native and FIP20 TokenPair Owners
func (k Keeper) ConvertCoin(goCtx context.Context, msg *types.MsgConvertCoin) (*types.MsgConvertCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	k.evmKeeper.WithContext(ctx)

	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	receiver := common.HexToAddress(msg.Receiver)

	err := k.ConvertDenomToFIP20(ctx, sender, receiver, msg.Coin)
	return &types.MsgConvertCoinResponse{}, err
}

// ConvertFIP20 converts FIP20 tokens into Cosmos-native Coins for both
// Cosmos-native and FIP20 TokenPair Owners
func (k Keeper) ConvertFIP20(goCtx context.Context, msg *types.MsgConvertFIP20) (*types.MsgConvertFIP20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	k.evmKeeper.WithContext(ctx)

	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	receiver, _ := sdk.AccAddressFromBech32(msg.Receiver)
	ethSender := common.BytesToAddress(sender)

	pubKey, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, msg.PubKey)
	if bytes.Equal(sender, pubKey.Address()) {
		decompressPubkey, _ := crypto.DecompressPubkey(pubKey.Bytes())
		ethSender = crypto.PubkeyToAddress(*decompressPubkey)
	}

	err := k.ConvertFIP20ToDenom(ctx, msg.ContractAddress, ethSender, receiver, msg.Amount)
	return &types.MsgConvertFIP20Response{}, err
}

func (k Keeper) ConvertDenomToFIP20(
	ctx sdk.Context,
	sender sdk.AccAddress,
	receiver common.Address,
	coin sdk.Coin) error {
	pair, err := k.MintingEnabled(ctx, sender, receiver.Bytes(), coin.Denom)
	if err != nil {
		return err
	}
	switch {
	case pair.IsNativeCoin():
		return k.convertDenomNativeCoin(ctx, pair, sender, receiver, coin)
	case pair.IsNativeFIP20():
		return k.convertDenomNativeFIP20(ctx, pair, sender, receiver, coin)
	default:
		return types.ErrUndefinedOwner
	}
}

func (k Keeper) ConvertFIP20ToDenom(
	ctx sdk.Context,
	contract string,
	sender common.Address,
	receiver sdk.AccAddress,
	amount sdk.Int) error {
	pair, err := k.MintingEnabled(ctx, sdk.AccAddress(sender.Bytes()), receiver, contract)
	if err != nil {
		return err
	}
	//check fip20 balance
	balanceOf, err := k.QueryFIP20BalanceOf(ctx, pair.GetFIP20Contract(), sender)
	if err != nil {
		return err
	}
	if balanceOf.Cmp(amount.BigInt()) < 0 {
		return fmt.Errorf("insufficient balance of %s at token %s", sender, pair.GetFIP20Contract().Hex())
	}
	switch {
	case pair.IsNativeCoin():
		return k.convertFIP20NativeDenom(ctx, pair, sender, receiver, amount)
	case pair.IsNativeFIP20():
		return k.convertFIP20NativeToken(ctx, pair, sender, receiver, amount)
	default:
		return types.ErrUndefinedOwner
	}
}

// convertDenomNativeCoin handles the Coin conversion flow for a native coin token pair:
//  - Escrow Coins on module account (Coins are not burned)
//  - Mint Tokens and send to receiver
func (k Keeper) convertDenomNativeCoin(ctx sdk.Context, pair types.TokenPair, sender sdk.AccAddress, receiver common.Address, coin sdk.Coin) error {
	coins := sdk.Coins{coin}
	fip20 := contracts.FIP20Contract.ABI
	contract := pair.GetFIP20Contract()

	// Escrow Coins on module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins); err != nil {
		return sdkerrors.Wrap(err, "failed to escrow coins")
	}

	// Mint Tokens and send to receiver
	_, err := k.CallEVMWithModule(ctx, fip20, contract, "mint", receiver, coin.Amount.BigInt())
	if err != nil {
		return sdkerrors.Wrap(err, "failed to call mint function with module")
	}

	// Event
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeConvertCoin,
				sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
				sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, coin.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyCosmosCoin, coin.Denom),
				sdk.NewAttribute(types.AttributeKeyFIP20Token, pair.Fip20Address),
			),
		},
	)
	return nil
}

// convertDenomNativeFIP20 handles the Coin conversion flow for a native FIP20 token pair:
//  - Escrow Coins on module account
//  - Unescrow Tokens that have been previously escrowed with ConvertFIP20 and send to receiver
//  - Burn escrowed Coins
func (k Keeper) convertDenomNativeFIP20(ctx sdk.Context, pair types.TokenPair, sender sdk.AccAddress, receiver common.Address, coin sdk.Coin) error {
	coins := sdk.Coins{coin}
	fip20 := contracts.FIP20Contract.ABI
	contract := pair.GetFIP20Contract()

	// Escrow Coins on module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins); err != nil {
		return sdkerrors.Wrap(err, "failed to escrow coins")
	}

	// Unescrow Tokens and send to receiver
	res, err := k.CallEVMWithModule(ctx, fip20, contract, "transfer", receiver, coin.Amount.BigInt())
	if err != nil {
		return sdkerrors.Wrap(err, "failed to call transfer function with module")
	}

	// Check unpackedRet execution
	var unpackedRet types.FIP20BoolResponse
	if err := fip20.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
		return sdkerrors.Wrap(err, "failed to unpack transfer return data")
	}
	if !unpackedRet.Value {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "failed to execute unescrow tokens from user")
	}

	// Burn escrowed Coins
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to burn escrowed coins")
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeConvertCoin,
				sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
				sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, coin.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyCosmosCoin, coin.Denom),
				sdk.NewAttribute(types.AttributeKeyFIP20Token, pair.Fip20Address),
			),
		},
	)
	return nil
}

// convertFIP20NativeDenom handles the fip20 conversion flow for a native coin token pair:
//  - Escrow tokens on module account
//  - Burn escrowed tokens
//  - Unescrow coins that have been previously escrowed with ConvertCoin
func (k Keeper) convertFIP20NativeDenom(ctx sdk.Context, pair types.TokenPair, sender common.Address, receiver sdk.AccAddress, amount sdk.Int) error {
	coins := sdk.Coins{sdk.Coin{Denom: pair.Denom, Amount: amount}}
	fip20 := contracts.FIP20Contract.ABI
	contract := pair.GetFIP20Contract()

	// Call evm to burn amount
	_, err := k.CallEVMWithModule(ctx, fip20, contract, "burn", sender, amount.BigInt())
	if err != nil {
		return sdkerrors.Wrap(err, "failed to call burn function with module")
	}

	// Unescrow Coins and send to receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, coins); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeConvertFIP20,
				sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
				sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
				sdk.NewAttribute(types.AttributeKeyCosmosCoin, pair.Denom),
				sdk.NewAttribute(types.AttributeKeyFIP20Token, contract.String()),
			),
		},
	)
	return nil
}

// convertFIP20NativeToken handles the fip20 conversion flow for a native fip20 token pair:
//  - Escrow tokens on module account (Don't burn as module is not contract owner)
//  - Mint coins on module
//  - Send minted coins to the receiver
func (k Keeper) convertFIP20NativeToken(ctx sdk.Context, pair types.TokenPair, sender common.Address, receiver sdk.AccAddress, amount sdk.Int) error {
	coins := sdk.Coins{sdk.Coin{Denom: pair.Denom, Amount: amount}}
	fip20 := contracts.FIP20Contract.ABI
	contract := pair.GetFIP20Contract()

	// Escrow tokens on module account
	transferData, err := fip20.Pack("transfer", types.ModuleAddress, amount.BigInt())
	if err != nil {
		return sdkerrors.Wrap(err, "failed to pack transfer")
	}
	// Call evm with eip55 address
	res, err := k.CallEVMWithPayload(ctx, sender, &contract, transferData)
	if err != nil {
		return sdkerrors.Wrap(err, fmt.Sprintf("failed to call transfer function with %s", sender.String()))
	}

	// Check unpackedRet execution
	var unpackedRet types.FIP20BoolResponse
	if err := fip20.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
		return sdkerrors.Wrap(err, "failed to unpack transfer return data")
	}

	if !unpackedRet.Value {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "failed to execute transfer")
	}

	// Mint coins
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return sdkerrors.Wrap(err, fmt.Sprintf("failed to mint coins %s", coins.String()))
	}

	// Send minted coins to the receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, coins); err != nil {
		return sdkerrors.Wrap(err, "failed to send coin from module")
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeConvertFIP20,
				sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
				sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
				sdk.NewAttribute(types.AttributeKeyCosmosCoin, pair.Denom),
				sdk.NewAttribute(types.AttributeKeyFIP20Token, contract.String()),
			),
		},
	)
	return nil
}
