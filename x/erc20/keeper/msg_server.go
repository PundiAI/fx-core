package keeper

import (
	"context"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"

	evmtypes "github.com/functionx/fx-core/x/evm/types"

	"github.com/functionx/fx-core/contracts"
	"github.com/functionx/fx-core/x/erc20/types"
)

var _ types.MsgServer = &Keeper{}

// ConvertCoin converts Cosmos-native Coins into ERC20 tokens for both
// Cosmos-native and ERC20 TokenPair Owners
func (k Keeper) ConvertCoin(goCtx context.Context, msg *types.MsgConvertCoin) (*types.MsgConvertCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Error checked during msg validation
	receiver := common.HexToAddress(msg.Receiver)
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)

	pair, err := k.MintingEnabled(ctx, sender, receiver.Bytes(), msg.Coin.Denom)
	if err != nil {
		return nil, err
	}

	// Check ownership and execute conversion
	switch {
	case pair.IsNativeCoin():
		return k.convertCoinNativeCoin(ctx, pair, msg, receiver, sender) // case 1.1
	case pair.IsNativeERC20():
		return k.convertCoinNativeERC20(ctx, pair, msg, receiver, sender) // case 2.2
	default:
		return nil, types.ErrUndefinedOwner
	}
}

// ConvertERC20 converts ERC20 tokens into Cosmos-native Coins for both
// Cosmos-native and ERC20 TokenPair Owners
func (k Keeper) ConvertERC20(goCtx context.Context, msg *types.MsgConvertERC20) (*types.MsgConvertERC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Error checked during msg validation
	receiver, _ := sdk.AccAddressFromBech32(msg.Receiver)
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

	// Check ownership
	switch {
	case pair.IsNativeCoin():
		return k.convertERC20NativeCoin(ctx, pair, msg, receiver, sender) // case 1.2
	case pair.IsNativeERC20():
		return k.convertERC20NativeToken(ctx, pair, msg, receiver, sender) // case 2.1
	default:
		return nil, types.ErrUndefinedOwner
	}
}

// convertCoinNativeCoin handles the Coin conversion flow for a native coin
// token pair:
//  - Escrow Coins on module account (Coins are not burned)
//  - Mint Tokens and send to receiver
//  - Check if token balance increased by amount
func (k Keeper) convertCoinNativeCoin(ctx sdk.Context, pair types.TokenPair, msg *types.MsgConvertCoin, receiver common.Address, sender sdk.AccAddress) (*types.MsgConvertCoinResponse, error) {
	// NOTE: ignore validation from NewCoin constructor
	coins := sdk.Coins{msg.Coin}
	erc20 := contracts.GetERC20(ctx.BlockHeight()).ABI
	contract := pair.GetERC20Contract()
	balanceToken := k.balanceOf(ctx, erc20, contract, receiver)

	// Escrow Coins on module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to escrow coins")
	}

	// Mint Tokens and send to receiver
	_, err := k.CallEVM(ctx, erc20, types.ModuleAddress, contract, "mint", receiver, msg.Coin.Amount.BigInt())
	if err != nil {
		return nil, err
	}

	evmParams := k.evmKeeper.GetParams(ctx)
	if pair.Denom == evmParams.EvmDenom {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, contract.Bytes(), coins); err != nil {
			return nil, sdkerrors.Wrap(err, "failed to transfer escrow coins to origin denom")
		}
	}

	// Check expected Receiver balance after transfer execution
	tokens := msg.Coin.Amount.BigInt()
	balanceTokenAfter := k.balanceOf(ctx, erc20, contract, receiver)
	exp := big.NewInt(0).Add(balanceToken, tokens)

	if r := balanceTokenAfter.Cmp(exp); r != 0 {
		return nil, sdkerrors.Wrapf(
			types.ErrBalanceInvariance,
			"invalid token balance - expected: %v, actual: %v", exp, balanceTokenAfter,
		)
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, msg.Type()},
			1,
			[]metrics.Label{
				telemetry.NewLabel("erc20", pair.Erc20Address),
				telemetry.NewLabel("denom", pair.Denom),
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
				sdk.NewAttribute(types.AttributeKeyCosmosCoin, msg.Coin.Denom),
				sdk.NewAttribute(types.AttributeKeyERC20Token, pair.Erc20Address),
			),
		},
	)

	return &types.MsgConvertCoinResponse{}, nil
}

// convertERC20NativeCoin handles the erc20 conversion flow for a native coin token pair:
//  - Burn escrowed tokens
//  - Unescrow coins that have been previously escrowed with ConvertCoin
//  - Check if coin balance increased by amount
//  - Check if token balance decreased by amount
func (k Keeper) convertERC20NativeCoin(ctx sdk.Context, pair types.TokenPair, msg *types.MsgConvertERC20, receiver sdk.AccAddress, sender common.Address) (*types.MsgConvertERC20Response, error) {
	// NOTE: coin fields already validated
	coins := sdk.Coins{sdk.Coin{Denom: pair.Denom, Amount: msg.Amount}}

	erc20 := contracts.GetERC20(ctx.BlockHeight()).ABI
	contract := pair.GetERC20Contract()
	balanceCoin := k.bankKeeper.GetBalance(ctx, receiver, pair.Denom)
	balanceToken := k.balanceOf(ctx, erc20, contract, sender)

	// Burn escrowed tokens
	_, err := k.CallEVM(ctx, erc20, types.ModuleAddress, contract, "burn", sender, msg.Amount.BigInt())
	if err != nil {
		return nil, err
	}

	// Transfer origin denom to module
	evmParams := k.evmKeeper.GetParams(ctx)
	if pair.Denom == evmParams.EvmDenom {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, contract.Bytes(), types.ModuleName, coins); err != nil {
			return nil, sdkerrors.Wrap(err, "failed to transfer origin denom to module")
		}
	}

	// Unescrow Coins and send to receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, coins); err != nil {
		return nil, err
	}

	// Check expected Receiver balance after transfer execution
	balanceCoinAfter := k.bankKeeper.GetBalance(ctx, receiver, pair.Denom)
	expCoin := balanceCoin.Add(coins[0])
	if ok := balanceCoinAfter.IsEqual(expCoin); !ok {
		return nil, sdkerrors.Wrapf(
			types.ErrBalanceInvariance,
			"invalid coin balance - expected: %v, actual: %v",
			expCoin, balanceCoinAfter,
		)
	}

	// Check expected Sender balance after transfer execution
	tokens := coins[0].Amount.BigInt()
	balanceTokenAfter := k.balanceOf(ctx, erc20, contract, sender)
	expToken := big.NewInt(0).Sub(balanceToken, tokens)
	if r := balanceTokenAfter.Cmp(expToken); r != 0 {
		return nil, sdkerrors.Wrapf(
			types.ErrBalanceInvariance,
			"invalid token balance - expected: %v, actual: %v",
			expToken, balanceTokenAfter,
		)
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, msg.Type()},
			1,
			[]metrics.Label{
				telemetry.NewLabel("erc20", pair.Erc20Address),
				telemetry.NewLabel("denom", pair.Denom),
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
				sdk.NewAttribute(types.AttributeKeyCosmosCoin, pair.Denom),
				sdk.NewAttribute(types.AttributeKeyERC20Token, msg.ContractAddress),
			),
		},
	)

	return &types.MsgConvertERC20Response{}, nil
}

// convertERC20NativeToken handles the erc20 conversion flow for a native erc20 token pair:
//  - Escrow tokens on module account (Don't burn as module is not contract owner)
//  - Mint coins on module
//  - Send minted coins to the receiver
//  - Check if coin balance increased by amount
//  - Check if token balance decreased by amount
//  - Check for unexpected `appove` event in logs
func (k Keeper) convertERC20NativeToken(ctx sdk.Context, pair types.TokenPair, msg *types.MsgConvertERC20, receiver sdk.AccAddress, sender common.Address) (*types.MsgConvertERC20Response, error) {
	// NOTE: coin fields already validated
	coins := sdk.Coins{sdk.Coin{Denom: pair.Denom, Amount: msg.Amount}}
	erc20 := contracts.GetERC20(ctx.BlockHeight()).ABI
	contract := pair.GetERC20Contract()
	balanceCoin := k.bankKeeper.GetBalance(ctx, receiver, pair.Denom)
	balanceToken := k.balanceOf(ctx, erc20, contract, types.ModuleAddress)

	// Escrow tokens on module account
	transferData, err := erc20.Pack("transfer", types.ModuleAddress, msg.Amount.BigInt())
	if err != nil {
		return nil, err
	}

	res, err := k.CallEVMWithData(ctx, sender, &contract, transferData)
	if err != nil {
		return nil, err
	}

	// Check unpackedRet execution
	var unpackedRet types.ERC20BoolResponse
	if err := erc20.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
		return nil, err
	}

	if !unpackedRet.Value {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "failed to execute transfer")
	}

	// Check expected escrow balance after transfer execution
	tokens := coins[0].Amount.BigInt()
	balanceTokenAfter := k.balanceOf(ctx, erc20, contract, types.ModuleAddress)
	expToken := big.NewInt(0).Add(balanceToken, tokens)

	if r := balanceTokenAfter.Cmp(expToken); r != 0 {
		return nil, sdkerrors.Wrapf(
			types.ErrBalanceInvariance,
			"invalid token balance - expected: %v, actual: %v",
			expToken, balanceTokenAfter,
		)
	}

	// Mint coins
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return nil, err
	}

	// Send minted coins to the receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, coins); err != nil {
		return nil, err
	}

	// Check expected Receiver balance after transfer execution
	balanceCoinAfter := k.bankKeeper.GetBalance(ctx, receiver, pair.Denom)
	expCoin := balanceCoin.Add(coins[0])

	if ok := balanceCoinAfter.IsEqual(expCoin); !ok {
		return nil, sdkerrors.Wrapf(
			types.ErrBalanceInvariance,
			"invalid coin balance - expected: %v, actual: %v",
			expCoin, balanceCoinAfter,
		)
	}

	// Check for unexpected `appove` event in logs
	if err := k.monitorApprovalEvent(res); err != nil {
		return nil, err
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, msg.Type()},
			1,
			[]metrics.Label{
				telemetry.NewLabel("erc20", pair.Erc20Address),
				telemetry.NewLabel("denom", pair.Denom),
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
				sdk.NewAttribute(types.AttributeKeyCosmosCoin, pair.Denom),
				sdk.NewAttribute(types.AttributeKeyERC20Token, msg.ContractAddress),
			),
		},
	)

	return &types.MsgConvertERC20Response{}, nil
}

// convertCoinNativeERC20 handles the Coin conversion flow for a native ERC20
// token pair:
//  - Escrow Coins on module account
//  - Unescrow Tokens that have been previously escrowed with ConvertERC20 and send to receiver
//  - Burn escrowed Coins
//  - Check if token balance increased by amount
//  - Check for unexpected `appove` event in logs
func (k Keeper) convertCoinNativeERC20(ctx sdk.Context, pair types.TokenPair, msg *types.MsgConvertCoin, receiver common.Address, sender sdk.AccAddress) (*types.MsgConvertCoinResponse, error) {
	// NOTE: ignore validation from NewCoin constructor
	coins := sdk.Coins{msg.Coin}

	erc20 := contracts.GetERC20(ctx.BlockHeight()).ABI
	contract := pair.GetERC20Contract()
	balanceToken := k.balanceOf(ctx, erc20, contract, receiver)

	// Escrow Coins on module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to escrow coins")
	}

	// Unescrow Tokens and send to receiver
	res, err := k.CallEVM(ctx, erc20, types.ModuleAddress, contract, "transfer", receiver, msg.Coin.Amount.BigInt())
	if err != nil {
		return nil, err
	}

	// Check unpackedRet execution
	var unpackedRet types.ERC20BoolResponse
	if err := erc20.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
		return nil, err
	}

	if !unpackedRet.Value {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "failed to execute unescrow tokens from user")
	}

	// Check expected Receiver balance after transfer execution
	tokens := msg.Coin.Amount.BigInt()
	balanceTokenAfter := k.balanceOf(ctx, erc20, contract, receiver)
	exp := big.NewInt(0).Add(balanceToken, tokens)

	if r := balanceTokenAfter.Cmp(exp); r != 0 {
		return nil, sdkerrors.Wrapf(
			types.ErrBalanceInvariance,
			"invalid token balance - expected: %v, actual: %v", exp, balanceTokenAfter,
		)
	}

	// Burn escrowed Coins
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to burn coins")
	}

	// Check for unexpected `appove` event in logs
	if err := k.monitorApprovalEvent(res); err != nil {
		return nil, err
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, msg.Type()},
			1,
			[]metrics.Label{
				telemetry.NewLabel("erc20", pair.Erc20Address),
				telemetry.NewLabel("denom", pair.Denom),
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
				sdk.NewAttribute(types.AttributeKeyCosmosCoin, msg.Coin.Denom),
				sdk.NewAttribute(types.AttributeKeyERC20Token, pair.Erc20Address),
			),
		},
	)

	return &types.MsgConvertCoinResponse{}, nil
}

// balanceOf queries an account's balance for a given ERC20 contract
func (k Keeper) balanceOf(ctx sdk.Context, abi abi.ABI, contract, account common.Address) *big.Int {
	res, err := k.CallEVM(ctx, abi, types.ModuleAddress, contract, "balanceOf", account)
	if err != nil {
		return nil
	}

	unpacked, err := abi.Unpack("balanceOf", res.Ret)
	if err != nil || len(unpacked) == 0 {
		return nil
	}

	balance, ok := unpacked[0].(*big.Int)
	if !ok {
		return nil
	}

	return balance
}

// monitorApprovalEvent returns an error if the given transactions logs include
// an unexpected `approve` event
func (k Keeper) monitorApprovalEvent(res *evmtypes.MsgEthereumTxResponse) error {
	if res == nil || len(res.Logs) == 0 {
		return nil
	}

	logApprovalSig := []byte("Approval(address,address,uint256)")
	logApprovalSigHash := crypto.Keccak256Hash(logApprovalSig)

	for _, log := range res.Logs {
		if log.Topics[0] == logApprovalSigHash.Hex() {
			return sdkerrors.Wrapf(
				types.ErrUnexpectedEvent, "unexpected Approval event",
			)
		}
	}

	return nil
}
