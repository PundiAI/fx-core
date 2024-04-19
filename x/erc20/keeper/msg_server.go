package keeper

import (
	"context"
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/slices"

	"github.com/functionx/fx-core/v7/contract"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/erc20/types"
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
	switch {
	case pair.IsNativeCoin():
		err = k.ConvertCoinNativeCoin(ctx, pair, sender, receiver, msg.Coin)
	case pair.IsNativeERC20():
		err = k.ConvertCoinNativeERC20(ctx, pair, sender, receiver, msg.Coin)
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
	switch {
	case pair.IsNativeCoin():
		err = k.ConvertERC20NativeCoin(ctx, pair, sender, receiver, msg.Amount)
	case pair.IsNativeERC20():
		err = k.ConvertERC20NativeToken(ctx, pair, sender, receiver, msg.Amount)
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
		return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "convert to source denom: %s", msg.Coin.Denom)
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
		return errorsmod.Wrap(err, "failed to escrow coins")
	}

	erc20 := contract.GetFIP20().ABI
	contract := pair.GetERC20Contract()

	// Mint Tokens and send to receiver
	_, err := k.evmKeeper.ApplyContract(ctx, k.moduleAddress, contract, nil, erc20, "mint", receiver, coin.Amount.BigInt())
	if err != nil {
		return err
	}

	if pair.Denom == fxtypes.DefaultDenom {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, contract.Bytes(), coins); err != nil {
			return errorsmod.Wrap(err, "failed to transfer escrow coins to origin denom")
		}
	}
	return nil
}

// ConvertERC20NativeCoin handles the erc20 conversion flow for a native coin token pair:
//   - Burn escrowed tokens
//   - Unescrow coins that have been previously escrowed with ConvertCoin
//   - Check if coin balance increased by amount
//   - Check if token balance decreased by amount
func (k Keeper) ConvertERC20NativeCoin(ctx sdk.Context, pair types.TokenPair, sender common.Address, receiver sdk.AccAddress, amount sdkmath.Int) error {
	erc20 := contract.GetFIP20().ABI
	contract := pair.GetERC20Contract()

	// Burn escrowed tokens
	_, err := k.evmKeeper.ApplyContract(ctx, k.moduleAddress, contract, nil, erc20, "burn", sender, amount.BigInt())
	if err != nil {
		return err
	}

	// NOTE: coin fields already validated
	coins := sdk.Coins{sdk.Coin{Denom: pair.Denom, Amount: amount}}

	// Transfer origin denom to module
	if pair.Denom == fxtypes.DefaultDenom {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, contract.Bytes(), types.ModuleName, coins); err != nil {
			return errorsmod.Wrap(err, "failed to transfer origin denom to module")
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
func (k Keeper) ConvertERC20NativeToken(ctx sdk.Context, pair types.TokenPair, sender common.Address, receiver sdk.AccAddress, amount sdkmath.Int) error {
	erc20 := contract.GetFIP20().ABI

	// Escrow tokens on module account
	contract := pair.GetERC20Contract()
	res, err := k.evmKeeper.ApplyContract(ctx, sender, contract, nil, erc20, "transfer", k.moduleAddress, amount.BigInt())
	if err != nil {
		return err
	}

	// Check unpackedRet execution
	var unpackedRet struct{ Value bool }
	if err := erc20.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
		return errorsmod.Wrapf(types.ErrABIUnpack, "failed to unpack transfer: %s", err.Error())
	}

	if !unpackedRet.Value {
		return errorsmod.Wrap(errortypes.ErrLogic, "failed to execute transfer")
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
		return errorsmod.Wrap(err, "failed to escrow coins")
	}

	erc20 := contract.GetFIP20().ABI
	contract := pair.GetERC20Contract()

	// Unescrow Tokens and send to receiver
	res, err := k.evmKeeper.ApplyContract(ctx, k.moduleAddress, contract, nil, erc20, "transfer", receiver, coin.Amount.BigInt())
	if err != nil {
		return err
	}

	// Check unpackedRet execution
	var unpackedRet struct{ Value bool }
	if err := erc20.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
		return errorsmod.Wrapf(types.ErrABIUnpack, "failed to unpack transfer: %s", err.Error())
	}

	if !unpackedRet.Value {
		return errorsmod.Wrap(errortypes.ErrLogic, "failed to execute unescrow tokens from user")
	}

	// Burn escrowed Coins
	if err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins); err != nil {
		return errorsmod.Wrap(err, "failed to burn coins")
	}

	// Check for unexpected `approve` event in logs
	if err = k.monitorApprovalEvent(res); err != nil {
		return err
	}
	return nil
}

func (k Keeper) ConvertDenomToTarget(ctx sdk.Context, from sdk.AccAddress, coin sdk.Coin, fxTarget fxtypes.FxTarget) (sdk.Coin, error) {
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

		// NOTE: metadata must exist, and alias must exist
		metadata, found = k.HasDenomAlias(ctx, denom)
		if !found { // no convert required
			return coin, nil
		}
	}

	targetDenom := k.ToTargetDenom(ctx, coin.Denom, metadata.Base, metadata.DenomUnits[0].Aliases, fxTarget)
	if coin.Denom == targetDenom {
		return coin, nil
	}

	targetCoin := sdk.NewCoin(targetDenom, coin.Amount)
	// send denom to module
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return sdk.Coin{}, err
	}
	if err = k.convertDenomToContractOwner(ctx, targetCoin, coin, metadata); err != nil {
		return sdk.Coin{}, err
	}
	moduleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
	moduleTargetCoin := k.bankKeeper.GetBalance(ctx, moduleAddress, targetCoin.Denom)
	if moduleTargetCoin.IsLT(targetCoin) {
		return targetCoin, types.ErrInsufficientLiquidity.Wrapf("%s is smaller than %s", moduleTargetCoin, targetCoin)
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

func (k Keeper) convertDenomToContractOwner(ctx sdk.Context, targetCoin, coin sdk.Coin, metadata banktypes.Metadata) error {
	pair, found := k.GetTokenPair(ctx, metadata.Base)
	if !found {
		return errorsmod.Wrapf(types.ErrTokenPairNotFound, "convert denom: %s", metadata.Base)
	}

	// converted metadata
	if k.IsConvertedMetadata(metadata) {
		return k.convertNativeAlias(ctx, targetCoin, coin, metadata)
	}

	// native coin
	if pair.IsNativeCoin() {
		return k.convertNativeCoin(ctx, targetCoin, coin, metadata)
	}

	// native erc20
	if pair.IsNativeERC20() {
		return k.convertNativeERC20(ctx, targetCoin, coin, metadata)
	}
	return errorsmod.Wrapf(types.ErrUndefinedOwner, "convert denom:%s, pair undefined owner.", metadata.Base)
}

func (k Keeper) IsConvertedMetadata(md banktypes.Metadata) bool {
	// one-to-one metadata
	if len(md.DenomUnits) == 0 || len(md.DenomUnits[0].Aliases) == 0 {
		return false
	}
	return k.checkConvertedDenom(md.Base)
}

func (k Keeper) convertNativeAlias(ctx sdk.Context, targetCoin, coin sdk.Coin, metadata banktypes.Metadata) error {
	// targetCoin is alias, coin is base
	// lock base, mint alias
	if coin.Denom == metadata.Base &&
		slices.Contains(metadata.DenomUnits[0].Aliases, targetCoin.Denom) {
		// already lock coin

		// mint alias
		return k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(targetCoin))
	}

	// targetCoin is base, coin is alias
	// burn alias, unlock base
	if targetCoin.Denom == metadata.Base &&
		slices.Contains(metadata.DenomUnits[0].Aliases, coin.Denom) {
		// unlock coin after burn

		// burn alias
		return k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
	}

	// reset is target Coin is alias1, coin is alias2
	// burn alias2, mint alias1
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return err
	}
	return k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(targetCoin))
}

func (k Keeper) convertNativeCoin(ctx sdk.Context, targetCoin, coin sdk.Coin, metadata banktypes.Metadata) error {
	if coin.Denom == metadata.Base {
		return k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
	}
	if targetCoin.Denom == metadata.Base {
		return k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(targetCoin))
	}
	// NOTE: convert alias to alias
	return nil
}

func (k Keeper) convertNativeERC20(ctx sdk.Context, targetCoin, coin sdk.Coin, metadata banktypes.Metadata) error {
	if coin.Denom == metadata.Base {
		return k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(targetCoin))
	}
	if targetCoin.Denom == metadata.Base {
		return k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
	}
	return nil
}

func (k Keeper) UpdateParams(c context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(c)
	if err := k.SetParams(ctx, &req.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateParamsResponse{}, nil
}

func (k Keeper) RegisterCoin(c context.Context, req *types.MsgRegisterCoin) (*types.MsgRegisterCoinResponse, error) {
	if k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(c)
	pair, err := k.RegisterNativeCoin(ctx, req.Metadata)
	if err != nil {
		return nil, err
	}
	return &types.MsgRegisterCoinResponse{Pair: *pair}, nil
}

func (k Keeper) RegisterERC20(c context.Context, req *types.MsgRegisterERC20) (*types.MsgRegisterERC20Response, error) {
	if k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(c)
	pair, err := k.RegisterNativeERC20(ctx, common.HexToAddress(req.Erc20Address), req.Aliases...)
	if err != nil {
		return nil, err
	}
	return &types.MsgRegisterERC20Response{Pair: *pair}, nil
}

func (k Keeper) ToggleTokenConversion(c context.Context, req *types.MsgToggleTokenConversion) (*types.MsgToggleTokenConversionResponse, error) {
	if k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(c)
	pair, err := k.ToggleTokenConvert(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &types.MsgToggleTokenConversionResponse{Pair: pair}, nil
}

func (k Keeper) UpdateDenomAlias(c context.Context, req *types.MsgUpdateDenomAlias) (*types.MsgUpdateDenomAliasResponse, error) {
	if k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(c)
	_, err := k.UpdateDenomAliases(ctx, req.Denom, req.Alias)
	if err != nil {
		return nil, err
	}
	return &types.MsgUpdateDenomAliasResponse{}, nil
}

func (k Keeper) ToTargetDenom(ctx sdk.Context, denom, base string, aliases []string, fxTarget fxtypes.FxTarget) string {
	// erc20
	if len(fxTarget.GetTarget()) == 0 || fxTarget.GetTarget() == types.ModuleName {
		return base
	}
	if len(aliases) == 0 {
		return denom
	}

	for _, alias := range aliases {
		if fxTarget.IsIBC() && strings.HasPrefix(alias, ibctransfertypes.DenomPrefix+"/") {
			hexHash := strings.TrimPrefix(alias, ibctransfertypes.DenomPrefix+"/")
			hash, err := ibctransfertypes.ParseHexHash(hexHash)
			if err != nil {
				k.Logger(ctx).Info("invalid ibc denom", "denom", alias)
				continue
			}
			denomTrace, found := k.ibcTransferKeeper.GetDenomTrace(ctx, hash)
			if !found {
				continue
			}
			if !strings.HasPrefix(denomTrace.GetPath(), fmt.Sprintf("%s/%s", fxTarget.SourcePort, fxTarget.SourceChannel)) {
				continue
			}
			return alias
		}

		if strings.HasPrefix(alias, fxTarget.GetTarget()) {
			return alias
		}
	}

	// if not match any alias, return base denom
	return base
}
