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
	"math/big"
)

var _ types.MsgServer = &Keeper{}

// ConvertCoin converts Cosmos-native Coins into FIP20 tokens for both
// Cosmos-native and FIP20 TokenPair Owners
func (k Keeper) ConvertCoin(goCtx context.Context, msg *types.MsgConvertCoin) (*types.MsgConvertCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.HasInit(ctx) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "intrarelayer module not enable")
	}

	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	receiver := common.HexToAddress(msg.Receiver)

	err := k.ConvertDenomToFIP20(ctx, sender, receiver, msg.Coin)
	return &types.MsgConvertCoinResponse{}, err
}

// ConvertFIP20 converts FIP20 tokens into Cosmos-native Coins for both
// Cosmos-native and FIP20 TokenPair Owners
func (k Keeper) ConvertFIP20(goCtx context.Context, msg *types.MsgConvertFIP20) (*types.MsgConvertFIP20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.HasInit(ctx) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "intrarelayer module not enable")
	}

	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	receiver, _ := sdk.AccAddressFromBech32(msg.Receiver)
	ethSender := common.BytesToAddress(sender)

	pubKey, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, msg.PubKey)
	if bytes.Equal(sender, pubKey.Address()) {
		decompressPubKey, _ := crypto.DecompressPubkey(pubKey.Bytes())
		ethSender = crypto.PubkeyToAddress(*decompressPubKey)
		if _, err := k.accountKeeper.GetSequence(ctx, ethSender.Bytes()); err != nil {
			accI := k.accountKeeper.NewAccountWithAddress(ctx, ethSender.Bytes())
			k.accountKeeper.SetAccount(ctx, accI)
		}
	}

	err := k.ConvertFIP20ToDenom(ctx, msg.ContractAddress, ethSender, sender, receiver, msg.Amount)
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

func (k Keeper) ConvertFIP20ToDenom(ctx sdk.Context, contract string, ethSender common.Address, sender, receiver sdk.AccAddress, amount sdk.Int) error {
	pair, err := k.MintingEnabled(ctx, sdk.AccAddress(ethSender.Bytes()), receiver, contract)
	if err != nil {
		return err
	}
	//check fip20 balance
	balanceOf, err := k.QueryFIP20BalanceOf(ctx, pair.GetFIP20Contract(), ethSender)
	if err != nil {
		return err
	}
	if balanceOf.Cmp(amount.BigInt()) < 0 {
		return fmt.Errorf("insufficient balance of %s at token %s", ethSender, pair.GetFIP20Contract().Hex())
	}
	switch {
	case pair.IsNativeCoin():
		return k.convertFIP20NativeDenom(ctx, pair, ethSender, sender, receiver, amount)
	case pair.IsNativeFIP20():
		return k.convertFIP20NativeToken(ctx, pair, ethSender, sender, receiver, amount)
	default:
		return types.ErrUndefinedOwner
	}
}

// convertDenomNativeCoin handles the Coin conversion flow for a native coin token pair:
//  - Escrow Coins on module account (Coins are not burned)
//  - Mint Tokens and send to receiver
func (k Keeper) convertDenomNativeCoin(ctx sdk.Context, pair types.TokenPair, sender sdk.AccAddress, receiver common.Address, coin sdk.Coin) error {
	coins := sdk.Coins{coin}
	fip20ABI, found := contracts.GetABI(ctx.BlockHeight(), contracts.FIP20UpgradeType)
	if !found {
		return sdkerrors.Wrap(types.ErrInvalidContract, "fip20 contract not found")
	}
	contract := pair.GetFIP20Contract()

	//check balance
	balanceCoin, balanceToken, err := k.balanceOfConvert(ctx, pair.Denom, sender, contract, receiver)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalidBalance, err.Error())
	}

	// Escrow Coins on module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins); err != nil {
		return sdkerrors.Wrap(err, "failed to escrow coins")
	}

	// Mint Tokens and send to receiver
	_, err = k.CallEVM(ctx, fip20ABI, types.ModuleAddress, contract, "mint", receiver, coin.Amount.BigInt())
	if err != nil {
		return sdkerrors.Wrap(err, "failed to call mint function with module")
	}

	evmParams := k.evmKeeper.GetParams(ctx)
	if pair.Denom == evmParams.EvmDenom {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, contract.Bytes(), coins); err != nil {
			return sdkerrors.Wrap(err, "failed to transfer escrow coins to origin denom")
		}
	}

	//check balance
	balanceCoinAfter, balanceTokenAfter, err := k.balanceOfConvert(ctx, pair.Denom, sender, contract, receiver)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalidBalance, err.Error())
	}
	expCoin := balanceCoinAfter.Add(coin)
	if !balanceCoin.Equal(expCoin) {
		return sdkerrors.Wrapf(types.ErrInvalidBalance, "invalid coin balance - expected: %v, actual: %v", expCoin, balanceCoin)
	}
	tokens := coin.Amount.BigInt()
	expToken := big.NewInt(0).Add(balanceToken, tokens)
	if r := balanceTokenAfter.Cmp(expToken); r != 0 {
		return sdkerrors.Wrapf(types.ErrInvalidBalance, "invalid token balance - expected: %v, actual: %v", expToken, balanceTokenAfter)
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
	fip20ABI, found := contracts.GetABI(ctx.BlockHeight(), contracts.FIP20UpgradeType)
	if !found {
		return sdkerrors.Wrap(types.ErrInvalidContract, "fip20 contract not found")
	}
	contract := pair.GetFIP20Contract()
	//check balance
	balanceCoin, balanceToken, err := k.balanceOfConvert(ctx, pair.Denom, sender, contract, receiver)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalidBalance, err.Error())
	}

	// Escrow Coins on module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins); err != nil {
		return sdkerrors.Wrap(err, "failed to escrow coins")
	}

	// Unescrow Tokens and send to receiver
	res, err := k.CallEVM(ctx, fip20ABI, types.ModuleAddress, contract, "transfer", receiver, coin.Amount.BigInt())
	if err != nil {
		return sdkerrors.Wrap(err, "failed to call transfer function with module")
	}

	// Check unpackedRet execution
	var unpackedRet types.FIP20BoolResponse
	if err := fip20ABI.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
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

	//check balance
	balanceCoinAfter, balanceTokenAfter, err := k.balanceOfConvert(ctx, pair.Denom, sender, contract, receiver)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalidBalance, err.Error())
	}
	tokens := coin.Amount.BigInt()
	expCoin := balanceCoinAfter.Add(coin)
	if !balanceCoin.Equal(expCoin) {
		return sdkerrors.Wrapf(types.ErrInvalidBalance, "invalid coin balance - expected: %v, actual: %v", expCoin, balanceCoin)
	}
	expToken := big.NewInt(0).Add(balanceToken, tokens)
	if r := balanceTokenAfter.Cmp(expToken); r != 0 {
		return sdkerrors.Wrapf(types.ErrInvalidBalance, "invalid token balance - expected: %v, actual: %v", expToken, balanceTokenAfter)
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
func (k Keeper) convertFIP20NativeDenom(ctx sdk.Context, pair types.TokenPair, ethSender common.Address, sender, receiver sdk.AccAddress, amount sdk.Int) error {
	coins := sdk.Coins{sdk.Coin{Denom: pair.Denom, Amount: amount}}
	fip20ABI, found := contracts.GetABI(ctx.BlockHeight(), contracts.FIP20UpgradeType)
	if !found {
		return sdkerrors.Wrap(types.ErrInvalidContract, "fip20 contract not found")
	}
	contract := pair.GetFIP20Contract()
	//check balance
	balanceCoin, balanceToken, err := k.balanceOfConvert(ctx, pair.Denom, receiver, contract, ethSender)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalidBalance, err.Error())
	}

	// Call evm to burn amount
	_, err = k.CallEVM(ctx, fip20ABI, types.ModuleAddress, contract, "burn", ethSender, amount.BigInt())
	if err != nil {
		return sdkerrors.Wrap(err, "failed to call burn function with module")
	}

	// Transfer origin fip20 to module
	evmParams := k.evmKeeper.GetParams(ctx)
	if pair.Denom == evmParams.EvmDenom {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, contract.Bytes(), types.ModuleName, coins); err != nil {
			return sdkerrors.Wrap(err, "failed to transfer origin fip20 to module")
		}
	}

	// Unescrow Coins and send to receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, coins); err != nil {
		return sdkerrors.Wrap(err, "failed to unescrow coins")
	}

	//check balance
	balanceCoinAfter, balanceTokenAfter, err := k.balanceOfConvert(ctx, pair.Denom, receiver, contract, ethSender)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalidBalance, err.Error())
	}
	expCoin := balanceCoin.Add(sdk.NewCoin(pair.Denom, amount))
	if !balanceCoinAfter.Equal(expCoin) {
		return sdkerrors.Wrapf(types.ErrInvalidBalance, "invalid coin balance - expected: %v, actual: %v", expCoin, balanceCoinAfter)
	}
	expToken := big.NewInt(0).Add(balanceTokenAfter, amount.BigInt())
	if r := balanceToken.Cmp(expToken); r != 0 {
		return sdkerrors.Wrapf(types.ErrInvalidBalance, "invalid token balance - expected: %v, actual: %v", expToken, balanceToken)
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeConvertFIP20,
				sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
				sdk.NewAttribute(types.AttributeKeyEthSender, ethSender.String()),
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
func (k Keeper) convertFIP20NativeToken(ctx sdk.Context, pair types.TokenPair, ethSender common.Address, sender, receiver sdk.AccAddress, amount sdk.Int) error {
	coins := sdk.Coins{sdk.Coin{Denom: pair.Denom, Amount: amount}}
	fip20ABI, found := contracts.GetABI(ctx.BlockHeight(), contracts.FIP20UpgradeType)
	if !found {
		return sdkerrors.Wrap(types.ErrInvalidContract, "fip20 contract not found")
	}
	contract := pair.GetFIP20Contract()

	//check balance
	balanceCoin, balanceToken, err := k.balanceOfConvert(ctx, pair.Denom, receiver, contract, ethSender)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalidBalance, err.Error())
	}

	// Escrow tokens on module account
	transferData, err := fip20ABI.Pack("transfer", types.ModuleAddress, amount.BigInt())
	if err != nil {
		return sdkerrors.Wrap(err, "failed to pack transfer")
	}
	// Call evm with eip55 address
	res, err := k.CallEVMWithPayload(ctx, ethSender, &contract, transferData)
	if err != nil {
		return sdkerrors.Wrap(err, fmt.Sprintf("failed to call transfer function with %s", ethSender.String()))
	}

	// Check unpackedRet execution
	var unpackedRet types.FIP20BoolResponse
	if err := fip20ABI.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
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

	//check balance
	balanceCoinAfter, balanceTokenAfter, err := k.balanceOfConvert(ctx, pair.Denom, receiver, contract, ethSender)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalidBalance, err.Error())
	}
	expCoin := balanceCoin.Add(sdk.NewCoin(pair.Denom, amount))
	if !balanceCoinAfter.Equal(expCoin) {
		return sdkerrors.Wrapf(types.ErrInvalidBalance, "invalid coin balance - expected: %v, actual: %v", expCoin, balanceCoinAfter)
	}
	expToken := big.NewInt(0).Add(balanceTokenAfter, amount.BigInt())
	if r := balanceToken.Cmp(expToken); r != 0 {
		return sdkerrors.Wrapf(types.ErrInvalidBalance, "invalid token balance - expected: %v, actual: %v", expToken, balanceToken)
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeConvertFIP20,
				sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
				sdk.NewAttribute(types.AttributeKeyEthSender, ethSender.String()),
				sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
				sdk.NewAttribute(types.AttributeKeyCosmosCoin, pair.Denom),
				sdk.NewAttribute(types.AttributeKeyFIP20Token, contract.String()),
			),
		},
	)
	return nil
}

func (k Keeper) balanceOfConvert(ctx sdk.Context, denom string, acc sdk.AccAddress, fip20, addr common.Address) (sdk.Coin, *big.Int, error) {
	//check balacne
	balanceCoin := k.bankKeeper.GetBalance(ctx, acc, denom)
	balanceToken, err := k.QueryFIP20BalanceOf(ctx, fip20, addr)
	if err != nil {
		return sdk.Coin{}, nil, fmt.Errorf("failed to get balance of %s", addr.String())
	}
	return balanceCoin, balanceToken, nil
}
