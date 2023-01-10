package keeper

import (
	"context"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

var _ types.MsgServer = &Keeper{}

// ConvertCoin converts Cosmos-native Coins into ERC20 tokens for both
// Cosmos-native and ERC20 TokenPair Owners
func (k Keeper) ConvertCoin(goCtx context.Context, msg *types.MsgConvertCoin) (*types.MsgConvertCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	receiver := common.HexToAddress(msg.Receiver)

	pair, err := k.MintingEnabled(ctx, receiver.Bytes(), msg.Coin.Denom)
	if err != nil {
		return nil, err
	}

	// Remove token pair if contract is suicided
	if acc := k.evmKeeper.GetAccountWithoutBalance(ctx, pair.GetERC20Contract()); acc == nil || !acc.IsContract() {
		k.RemoveTokenPair(ctx, pair)
		k.Logger(ctx).Debug("deleting selfdestructed token pair from state", "contract", pair.Erc20Address)
		// NOTE: return nil error to persist the changes from the deletion
		return &types.MsgConvertCoinResponse{}, nil
	}

	// Check ownership and execute conversion
	newCtx := ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	switch {
	case pair.IsNativeCoin():
		err = k.ConvertCoinNativeCoin(newCtx, pair, sender, receiver, msg.Coin)
	case pair.IsNativeERC20():
		err = k.ConvertCoinNativeERC20(newCtx, pair, sender, receiver, msg.Coin)
	default:
		return nil, types.ErrUndefinedOwner
	}
	if err != nil {
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

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeConvertCoin,
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
		sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Coin.Amount.String()),
		sdk.NewAttribute(types.AttributeKeyDenom, msg.Coin.Denom),
		sdk.NewAttribute(types.AttributeKeyTokenAddress, pair.Erc20Address),
	))
	return &types.MsgConvertCoinResponse{}, nil
}

// ConvertERC20 converts ERC20 tokens into Cosmos-native Coins for both
// Cosmos-native and ERC20 TokenPair Owners
func (k Keeper) ConvertERC20(goCtx context.Context, msg *types.MsgConvertERC20) (*types.MsgConvertERC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender := common.HexToAddress(msg.Sender)
	receiver := sdk.MustAccAddressFromBech32(msg.Receiver)

	pair, err := k.MintingEnabled(ctx, receiver, msg.ContractAddress)
	if err != nil {
		return nil, err
	}

	// Remove token pair if contract is suicided
	if acc := k.evmKeeper.GetAccountWithoutBalance(ctx, pair.GetERC20Contract()); acc == nil || !acc.IsContract() {
		k.RemoveTokenPair(ctx, pair)
		k.Logger(ctx).Debug("deleting selfdestructed token pair from state", "contract", pair.Erc20Address)
		// NOTE: return nil error to persist the changes from the deletion
		return &types.MsgConvertERC20Response{}, nil
	}

	// Check ownership and execute conversion
	newCtx := ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	switch {
	case pair.IsNativeCoin():
		err = k.ConvertERC20NativeCoin(newCtx, pair, sender, receiver, msg.Amount)
	case pair.IsNativeERC20():
		err = k.ConvertERC20NativeToken(newCtx, pair, sender, receiver, msg.Amount)
	default:
		return nil, types.ErrUndefinedOwner
	}
	if err != nil {
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

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeConvertERC20,
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
		sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		sdk.NewAttribute(types.AttributeKeyDenom, pair.Denom),
		sdk.NewAttribute(types.AttributeKeyTokenAddress, msg.ContractAddress),
	))
	return &types.MsgConvertERC20Response{}, nil
}

// ConvertDenom converts coin into other coin, use for multiple chains in the same currency
func (k Keeper) ConvertDenom(goCtx context.Context, msg *types.MsgConvertDenom) (*types.MsgConvertDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Error checked during msg validation
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	receiver := sdk.MustAccAddressFromBech32(msg.Receiver)

	fxTarget := fxtypes.ParseFxTarget(msg.Target)
	targetCoin, err := k.ConvertDenomToTarget(ctx, sender, msg.Coin, fxTarget)
	if err != nil {
		return nil, err
	}
	if targetCoin.Denom == msg.Coin.Denom {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "denom %s not support", msg.Coin.Denom)
	}

	if !sender.Equals(receiver) {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(targetCoin)); err != nil {
			return nil, err
		}
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(targetCoin)); err != nil {
			return nil, err
		}
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tx", "msg", "convert", "denom", "total"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("denom", msg.Coin.Denom),
				telemetry.NewLabel("target", msg.Target),
			},
		)
	}()

	// ctx.EventManager().EmitEvent(sdk.NewEvent(
	//	types.EventTypeConvertDenom,
	//	sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	//	sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
	// ))

	return &types.MsgConvertDenomResponse{}, nil
}

// ConvertCoinNativeCoin handles the Coin conversion flow for a native coin
// token pair:
//   - Escrow Coins on module account (Coins are not burned)
//   - Mint Tokens and send to receiver
//   - Check if token balance increased by amount
func (k Keeper) ConvertCoinNativeCoin(ctx sdk.Context, pair types.TokenPair, sender sdk.AccAddress, receiver common.Address, coin sdk.Coin) error {
	// NOTE: ignore validation from NewCoin constructor
	coins := sdk.Coins{coin}

	// Escrow Coins on module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins); err != nil {
		return sdkerrors.Wrap(err, "failed to escrow coins")
	}

	erc20 := fxtypes.GetERC20().ABI
	contract := pair.GetERC20Contract()

	// Mint Tokens and send to receiver
	_, err := k.CallEVM(ctx, erc20, k.moduleAddress, contract, true, "mint", receiver, coin.Amount.BigInt())
	if err != nil {
		return err
	}

	if pair.Denom == fxtypes.DefaultDenom {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, contract.Bytes(), coins); err != nil {
			return sdkerrors.Wrap(err, "failed to transfer escrow coins to origin denom")
		}
	}
	return nil
}

// ConvertERC20NativeCoin handles the erc20 conversion flow for a native coin token pair:
//   - Burn escrowed tokens
//   - Unescrow coins that have been previously escrowed with ConvertCoin
//   - Check if coin balance increased by amount
//   - Check if token balance decreased by amount
func (k Keeper) ConvertERC20NativeCoin(ctx sdk.Context, pair types.TokenPair, sender common.Address, receiver sdk.AccAddress, amount sdk.Int) error {
	erc20 := fxtypes.GetERC20().ABI
	contract := pair.GetERC20Contract()

	// Burn escrowed tokens
	_, err := k.CallEVM(ctx, erc20, k.moduleAddress, contract, true, "burn", sender, amount.BigInt())
	if err != nil {
		return err
	}

	// NOTE: coin fields already validated
	coins := sdk.Coins{sdk.Coin{Denom: pair.Denom, Amount: amount}}

	// Transfer origin denom to module
	if pair.Denom == fxtypes.DefaultDenom {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, contract.Bytes(), types.ModuleName, coins); err != nil {
			return sdkerrors.Wrap(err, "failed to transfer origin denom to module")
		}
	}

	// Unescrow Coins and send to receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, coins); err != nil {
		return err
	}
	return nil
}

// ConvertERC20NativeToken handles the erc20 conversion flow for a native erc20 token pair:
//   - Escrow tokens on module account (Don't burn as module is not contract owner)
//   - Mint coins on module
//   - Send minted coins to the receiver
//   - Check if coin balance increased by amount
//   - Check if token balance decreased by amount
//   - Check for unexpected `approve` event in logs
func (k Keeper) ConvertERC20NativeToken(ctx sdk.Context, pair types.TokenPair, sender common.Address, receiver sdk.AccAddress, amount sdk.Int) error {
	erc20 := fxtypes.GetERC20().ABI

	// Escrow tokens on module account
	transferData, err := erc20.Pack("transfer", k.moduleAddress, amount.BigInt())
	if err != nil {
		return sdkerrors.Wrapf(types.ErrABIPack, "failed to pack transfer: %s", err.Error())
	}

	contract := pair.GetERC20Contract()
	res, err := k.evmKeeper.CallEVMWithData(ctx, sender, &contract, transferData, true)
	if err != nil {
		return err
	}

	// Check unpackedRet execution
	var unpackedRet struct{ Value bool }
	if err := erc20.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
		return sdkerrors.Wrapf(types.ErrABIUnpack, "failed to unpack transfer: %s", err.Error())
	}

	if !unpackedRet.Value {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "failed to execute transfer")
	}

	// Mint coins
	// NOTE: coin fields already validated
	coins := sdk.Coins{sdk.Coin{Denom: pair.Denom, Amount: amount}}
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	// Send minted coins to the receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, coins); err != nil {
		return err
	}

	// Check for unexpected `approve` event in logs
	if err = k.monitorApprovalEvent(res); err != nil {
		return err
	}
	return nil
}

// ConvertCoinNativeERC20 handles the Coin conversion flow for a native ERC20
// token pair:
//   - Escrow Coins on module account
//   - Unescrow Tokens that have been previously escrowed with ConvertERC20 and send to receiver
//   - Burn escrowed Coins
//   - Check if token balance increased by amount
//   - Check for unexpected `approve` event in logs
func (k Keeper) ConvertCoinNativeERC20(ctx sdk.Context, pair types.TokenPair, sender sdk.AccAddress, receiver common.Address, coin sdk.Coin) error {
	// NOTE: ignore validation from NewCoin constructor
	coins := sdk.Coins{coin}

	// Escrow Coins on module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins); err != nil {
		return sdkerrors.Wrap(err, "failed to escrow coins")
	}

	erc20 := fxtypes.GetERC20().ABI
	contract := pair.GetERC20Contract()

	// Unescrow Tokens and send to receiver
	res, err := k.CallEVM(ctx, erc20, k.moduleAddress, contract, true, "transfer", receiver, coin.Amount.BigInt())
	if err != nil {
		return err
	}

	// Check unpackedRet execution
	var unpackedRet struct{ Value bool }
	if err := erc20.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
		return sdkerrors.Wrapf(types.ErrABIUnpack, "failed to unpack transfer: %s", err.Error())
	}

	if !unpackedRet.Value {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "failed to execute unescrow tokens from user")
	}

	// Burn escrowed Coins
	if err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins); err != nil {
		return sdkerrors.Wrap(err, "failed to burn coins")
	}

	// Check for unexpected `approve` event in logs
	if err = k.monitorApprovalEvent(res); err != nil {
		return err
	}
	return nil
}

func (k Keeper) ConvertDenomToTarget(ctx sdk.Context, from sdk.AccAddress, coin sdk.Coin, fxTarget fxtypes.FxTarget) (sdk.Coin, error) {
	if coin.Denom == fxtypes.DefaultDenom {
		return coin, nil
	}
	var metadata banktypes.Metadata
	if k.IsDenomRegistered(ctx, coin.Denom) {
		// is base denom
		var found bool
		metadata, found = k.HasDenomAlias(ctx, coin.Denom)
		if !found { // no convert required
			return coin, nil
		}
	} else {
		// is alias denom
		denom, found := k.GetAliasDenom(ctx, coin.Denom)
		if !found { // no convert required
			return coin, nil
		}
		metadata, found = k.HasDenomAlias(ctx, denom)
		if !found { // no convert required
			return coin, nil
		}
	}

	targetDenom := ToTargetDenom(coin.Denom, metadata.Base, metadata.DenomUnits[0].Aliases, fxTarget)
	if coin.Denom == targetDenom {
		return coin, nil
	}

	targetCoin := sdk.NewCoin(targetDenom, coin.Amount)
	// send denom to module
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return sdk.Coin{}, err
	}

	if coin.Denom == metadata.Base {
		// burn coin
		if err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
			return sdk.Coin{}, err
		}
	} else if targetDenom == metadata.Base {
		// mint denom
		if err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(targetCoin)); err != nil {
			return sdk.Coin{}, err
		}
	}

	// send alias denom to from addr
	if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, from, sdk.NewCoins(targetCoin)); err != nil {
		return sdk.Coin{}, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeConvertDenom,
		sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
		sdk.NewAttribute(sdk.AttributeKeyAmount, coin.Amount.String()),
		sdk.NewAttribute(types.AttributeKeyDenom, coin.Denom),
		sdk.NewAttribute(types.AttributeKeyTargetDenom, targetCoin.Denom),
	))
	return targetCoin, nil
}

func ToTargetDenom(denom, base string, aliases []string, fxTarget fxtypes.FxTarget) string {
	// erc20
	if len(fxTarget.GetTarget()) <= 0 || fxTarget.GetTarget() == types.ModuleName {
		return base
	}
	if len(aliases) <= 0 {
		return denom
	}

	// ibc
	target := fxTarget.GetTarget()
	if fxTarget.IsIBC() {
		target = ibchost.ModuleName
	}

	for _, alias := range aliases {
		if strings.HasPrefix(alias, target) {
			return alias
		}
	}
	return denom
}
