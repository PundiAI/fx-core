package crosschain

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
	"github.com/functionx/fx-core/v7/x/evm/types"
)

// FIP20CrossChain only for fip20 contract transferCrossChain called
//
//gocyclo:ignore
func (c *Contract) FIP20CrossChain(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("fip20 cross chain method not readonly")
	}

	tokenContract := contract.Caller()
	tokenPair, found := c.erc20Keeper.GetTokenPairByAddress(ctx, tokenContract)
	if !found {
		return nil, fmt.Errorf("token pair not found: %s", tokenContract.String())
	}

	var args FIP20CrossChainArgs
	if err := types.ParseMethodArgs(FIP20CrossChainMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}

	amountCoin := sdk.NewCoin(tokenPair.GetDenom(), sdkmath.NewIntFromBigInt(args.Amount))
	feeCoin := sdk.NewCoin(tokenPair.GetDenom(), sdkmath.NewIntFromBigInt(args.Fee))
	totalCoin := sdk.NewCoin(tokenPair.GetDenom(), amountCoin.Amount.Add(feeCoin.Amount))

	// NOTE: if user call evm denom transferCrossChain with msg.value
	// we need transfer msg.value from sender to contract in bank keeper
	evmDenom := c.evmKeeper.GetParams(ctx).EvmDenom
	if tokenPair.GetDenom() == evmDenom {
		balance := c.bankKeeper.GetBalance(ctx, tokenContract.Bytes(), evmDenom)
		evmBalance := evm.StateDB.GetBalance(tokenContract)

		cmp := evmBalance.Cmp(balance.Amount.BigInt())
		if cmp == -1 {
			return nil, fmt.Errorf("invalid balance(chain: %s,evm: %s)", balance.Amount.String(), evmBalance.String())
		}
		if cmp == 1 {
			// sender call transferCrossChain with msg.value, the msg.value evm denom should send to contract
			value := big.NewInt(0).Sub(evmBalance, balance.Amount.BigInt())
			valueCoin := sdk.NewCoins(sdk.NewCoin(evmDenom, sdkmath.NewIntFromBigInt(value)))
			if err := c.bankKeeper.SendCoins(ctx, args.Sender.Bytes(), tokenContract.Bytes(), valueCoin); err != nil {
				return nil, fmt.Errorf("send coin: %s", err.Error())
			}
		}
	}

	// transfer token from evm to local chain
	if err := c.convertERC20(ctx, evm, tokenPair, totalCoin, args.Sender); err != nil {
		return nil, err
	}

	fxTarget := fxtypes.ParseFxTarget(fxtypes.Byte32ToString(args.Target))
	if err := c.handlerCrossChain(ctx, args.Sender.Bytes(), args.Receipt, amountCoin, feeCoin, fxTarget, args.Memo, false); err != nil {
		return nil, err
	}

	// add event log
	if err := c.AddLog(evm, CrossChainEvent, []common.Hash{args.Sender.Hash(), tokenPair.GetERC20Contract().Hash()},
		tokenPair.GetDenom(), args.Receipt, args.Amount, args.Fee, args.Target, args.Memo); err != nil {
		return nil, err
	}

	// add fip20CrossChain events
	fip20CrossChainEvents(ctx, args.Sender, tokenPair.GetERC20Contract(), args.Receipt,
		fxtypes.Byte32ToString(args.Target), tokenPair.GetDenom(), args.Amount, args.Fee)

	return FIP20CrossChainMethod.Outputs.Pack(true)
}

// CrossChain called at any address(account or contract)
//
//gocyclo:ignore
func (c *Contract) CrossChain(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("cross chain method not readonly")
	}

	// args
	var args CrossChainArgs
	err := types.ParseMethodArgs(CrossChainMethod, &args, contract.Input[4:])
	if err != nil {
		return nil, err
	}

	// call param
	value := contract.Value()
	sender := contract.Caller()

	// cross chain param
	originToken := false
	crossChainDenom := ""

	// cross-chain origin token
	if value.Cmp(big.NewInt(0)) == 1 && args.Token.String() == fxtypes.EmptyEvmAddress {
		totalAmount := big.NewInt(0).Add(args.Amount, args.Fee)
		if totalAmount.Cmp(value) != 0 {
			return nil, errors.New("amount + fee not equal msg.value")
		}

		crossChainDenom, err = c.handlerOriginToken(ctx, evm, sender, totalAmount)
		if err != nil {
			return nil, err
		}

		// origin token flag is true when cross chain evm denom
		originToken = true
	} else {
		crossChainDenom, err = c.handlerERC20Token(ctx, evm, args.Token, sender, big.NewInt(0).Add(args.Amount, args.Fee))
		if err != nil {
			return nil, err
		}
	}

	fxTarget := fxtypes.ParseFxTarget(fxtypes.Byte32ToString(args.Target))
	amountCoin := sdk.NewCoin(crossChainDenom, sdkmath.NewIntFromBigInt(args.Amount))
	feeCoin := sdk.NewCoin(crossChainDenom, sdkmath.NewIntFromBigInt(args.Fee))

	if err := c.handlerCrossChain(ctx, sender.Bytes(), args.Receipt, amountCoin, feeCoin, fxTarget, args.Memo, originToken); err != nil {
		return nil, err
	}

	// add event log
	if err := c.AddLog(evm, CrossChainEvent, []common.Hash{sender.Hash(), args.Token.Hash()},
		crossChainDenom, args.Receipt, args.Amount, args.Fee, args.Target, args.Memo); err != nil {
		return nil, err
	}

	// add cross chain events
	crossChainEvents(ctx, sender, args.Token, args.Receipt, fxtypes.Byte32ToString(args.Target),
		crossChainDenom, args.Memo, args.Amount, args.Fee)

	return CrossChainMethod.Outputs.Pack(true)
}

func (c *Contract) handlerOriginToken(ctx sdk.Context, evm *vm.EVM, sender common.Address, amount *big.Int) (string, error) {
	crossChainDenom := c.evmKeeper.GetParams(ctx).EvmDenom
	// NOTE: stateDB sub sender balance,but bank keeper not update.
	// so mint token to crosschain, end of stateDB commit will sub balance from bank keeper.
	// if only allow depth 1, the sender is origin sender, we can sub balance from bank keeper and not need burn/mint coin
	evm.StateDB.SubBalance(c.Address(), amount)
	totalCoin := sdk.NewCoins(sdk.NewCoin(crossChainDenom, sdkmath.NewIntFromBigInt(amount)))

	if err := c.bankKeeper.MintCoins(ctx, evmtypes.ModuleName, totalCoin); err != nil {
		return "", fmt.Errorf("mint: %s", err.Error())
	}
	if err := c.bankKeeper.SendCoinsFromModuleToAccount(ctx, evmtypes.ModuleName, sender.Bytes(), totalCoin); err != nil {
		return "", fmt.Errorf("send account: %s", err.Error())
	}
	return crossChainDenom, nil
}

func (c *Contract) handlerERC20Token(ctx sdk.Context, evm *vm.EVM, token, sender common.Address, amount *big.Int) (string, error) {
	// contract token
	tokenPair, found := c.erc20Keeper.GetTokenPairByAddress(ctx, token)
	if !found {
		return "", fmt.Errorf("token pair not found: %s", token.String())
	}
	crossChainDenom := tokenPair.GetDenom()
	// transferFrom to erc20 module
	if err := NewContractCall(ctx, evm, c.Address(), token).ERC20TransferFrom(sender, c.erc20Keeper.ModuleAddress(), amount); err != nil {
		return "", err
	}
	if err := c.convertERC20(ctx, evm, tokenPair, sdk.NewCoin(crossChainDenom, sdkmath.NewIntFromBigInt(amount)), sender); err != nil {
		return "", err
	}
	return crossChainDenom, nil
}

func (c *Contract) convertERC20(
	ctx sdk.Context,
	evm *vm.EVM,
	tokenPair erc20types.TokenPair,
	amount sdk.Coin,
	receiver common.Address,
) error {
	if tokenPair.IsNativeCoin() {
		err := NewContractCall(ctx, evm, c.erc20Keeper.ModuleAddress(), tokenPair.GetERC20Contract()).ERC20Burn(amount.Amount.BigInt())
		if err != nil {
			return err
		}
		if tokenPair.GetDenom() == fxtypes.DefaultDenom {
			// cache token contract balance
			evm.StateDB.GetBalance(tokenPair.GetERC20Contract())

			err := c.bankKeeper.SendCoinsFromAccountToModule(ctx, tokenPair.GetERC20Contract().Bytes(), erc20types.ModuleName, sdk.NewCoins(amount))
			if err != nil {
				return fmt.Errorf("send module: %s", err.Error())
			}

			// evm stateDB sub token contract balance
			evm.StateDB.SubBalance(tokenPair.GetERC20Contract(), amount.Amount.BigInt())
		}

	} else if tokenPair.IsNativeERC20() {
		err := c.bankKeeper.MintCoins(ctx, erc20types.ModuleName, sdk.NewCoins(amount))
		if err != nil {
			return fmt.Errorf("mint: %s", err.Error())
		}
	} else {
		return erc20types.ErrUndefinedOwner
	}

	sendAddr := sdk.AccAddress(receiver.Bytes())
	if err := c.bankKeeper.SendCoinsFromModuleToAccount(ctx, erc20types.ModuleName, sendAddr, sdk.NewCoins(amount)); err != nil {
		return fmt.Errorf("send account: %s", err.Error())
	}

	return nil
}

// handlerCrossChain cross chain handler
// originToken is true represent cross chain denom(FX)
// when refund it, will not refund to evm token
// NOTE: fip20CrossChain only use for contract token, so origin token flag always false
func (c *Contract) handlerCrossChain(
	ctx sdk.Context,
	from sdk.AccAddress,
	receipt string,
	amount, fee sdk.Coin,
	fxTarget fxtypes.FxTarget,
	memo string,
	originToken bool,
) error {
	total := sdk.NewCoin(amount.Denom, amount.Amount.Add(fee.Amount))
	// convert denom to target coin
	targetCoin, err := c.erc20Keeper.ConvertDenomToTarget(ctx, from.Bytes(), total, fxTarget)
	if err != nil {
		return fmt.Errorf("convert denom: %s", err.Error())
	}
	amount.Denom = targetCoin.Denom
	fee.Denom = targetCoin.Denom

	if fxTarget.IsIBC() {
		return c.ibcTransfer(ctx, from.Bytes(), receipt, amount, fee, fxTarget, memo, originToken)
	}
	return c.outgoingTransfer(ctx, from.Bytes(), receipt, amount, fee, fxTarget, originToken)
}

func (c *Contract) outgoingTransfer(
	ctx sdk.Context,
	from sdk.AccAddress,
	to string,
	amount, fee sdk.Coin,
	fxTarget fxtypes.FxTarget,
	originToken bool,
) error {
	if c.router == nil {
		return errors.New("cross chain router empty")
	}
	route, has := c.router.GetRoute(fxTarget.GetTarget())
	if !has {
		return errors.New("invalid target")
	}
	if err := route.TransferAfter(ctx, from, to, amount, fee, originToken); err != nil {
		return fmt.Errorf("cross chain error: %s", err.Error())
	}
	return nil
}

func (c *Contract) ibcTransfer(
	ctx sdk.Context,
	from sdk.AccAddress,
	to string,
	amount, fee sdk.Coin,
	fxTarget fxtypes.FxTarget,
	memo string,
	originToken bool,
) error {
	if !fee.IsZero() {
		return fmt.Errorf("ibc transfer fee must be zero: %s", fee.String())
	}
	if strings.ToLower(fxTarget.Prefix) == fxtypes.EthereumAddressPrefix {
		if err := fxtypes.ValidateEthereumAddress(to); err != nil {
			return fmt.Errorf("invalid to address: %s", to)
		}
	} else {
		if _, err := sdk.GetFromBech32(to, fxTarget.Prefix); err != nil {
			return fmt.Errorf("invalid to address: %s", to)
		}
	}

	ibcTimeoutTimestamp := uint64(ctx.BlockTime().UnixNano()) + uint64(c.erc20Keeper.GetIbcTimeout(ctx))
	transferResponse, err := c.ibcTransferKeeper.Transfer(sdk.WrapSDKContext(ctx),
		transfertypes.NewMsgTransfer(
			fxTarget.SourcePort,
			fxTarget.SourceChannel,
			amount,
			from.String(),
			to,
			ibcclienttypes.ZeroHeight(),
			ibcTimeoutTimestamp,
			memo,
		),
	)
	if err != nil {
		return fmt.Errorf("ibc transfer error: %s", err.Error())
	}

	if !originToken {
		c.erc20Keeper.SetIBCTransferRelation(ctx, fxTarget.SourceChannel, transferResponse.GetSequence())
	}
	return nil
}

// transferCrossChainEvents use for fip20 cross chain
// Deprecated
func fip20CrossChainEvents(ctx sdk.Context, from, token common.Address, recipient, target, denom string, amount, fee *big.Int) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		EventTypeRelayTransferCrossChain,
		sdk.NewAttribute(AttributeKeyFrom, from.String()),
		sdk.NewAttribute(AttributeKeyRecipient, recipient),
		sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
		sdk.NewAttribute(sdk.AttributeKeyFee, fee.String()),
		sdk.NewAttribute(AttributeKeyTarget, target),
		sdk.NewAttribute(AttributeKeyTokenAddress, token.String()),
		sdk.NewAttribute(AttributeKeyDenom, denom),
	))

	telemetry.IncrCounterWithLabels(
		[]string{"relay_transfer_cross_chain"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("erc20", token.String()),
			telemetry.NewLabel("denom", denom),
			telemetry.NewLabel("target", target),
		},
	)
}

func crossChainEvents(ctx sdk.Context, from, token common.Address, recipient, target, denom, memo string, amount, fee *big.Int) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		EventTypeCrossChain,
		sdk.NewAttribute(AttributeKeyFrom, from.String()),
		sdk.NewAttribute(AttributeKeyRecipient, recipient),
		sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
		sdk.NewAttribute(sdk.AttributeKeyFee, fee.String()),
		sdk.NewAttribute(AttributeKeyTarget, target),
		sdk.NewAttribute(AttributeKeyTokenAddress, token.String()),
		sdk.NewAttribute(AttributeKeyDenom, denom),
		sdk.NewAttribute(AttributeKeyMemo, memo),
	))
}
