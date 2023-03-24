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
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	"github.com/functionx/fx-core/v3/x/evm/types"
)

var (
	FIP20CrossChainMethod = abi.NewMethod(
		FIP20CrossChainMethodName,
		FIP20CrossChainMethodName,
		abi.Function, "nonpayable", false, false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: types.TypeAddress},
			abi.Argument{Name: "receipt", Type: types.TypeString},
			abi.Argument{Name: "amount", Type: types.TypeUint256},
			abi.Argument{Name: "fee", Type: types.TypeUint256},
			abi.Argument{Name: "target", Type: types.TypeBytes32},
			abi.Argument{Name: "memo", Type: types.TypeString},
		},
		abi.Arguments{
			abi.Argument{Name: "result", Type: types.TypeBool},
		},
	)
	CrossChainMethod = abi.NewMethod(
		CrossChainMethodName,
		CrossChainMethodName,
		abi.Function, "payable", false, true,
		abi.Arguments{
			abi.Argument{Name: "token", Type: types.TypeAddress},
			abi.Argument{Name: "receipt", Type: types.TypeString},
			abi.Argument{Name: "amount", Type: types.TypeUint256},
			abi.Argument{Name: "fee", Type: types.TypeUint256},
			abi.Argument{Name: "target", Type: types.TypeBytes32},
			abi.Argument{Name: "memo", Type: types.TypeString},
		},
		abi.Arguments{
			abi.Argument{Name: "result", Type: types.TypeBool},
		})

	CrossChainEvent = abi.NewEvent(
		CrossChainEventName,
		CrossChainEventName,
		false,
		abi.Arguments{
			abi.Argument{Name: "sender", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "token", Type: types.TypeAddress, Indexed: true},
			abi.Argument{Name: "denom", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "receipt", Type: types.TypeString, Indexed: false},
			abi.Argument{Name: "amount", Type: types.TypeUint256, Indexed: false},
			abi.Argument{Name: "fee", Type: types.TypeUint256, Indexed: false},
			abi.Argument{Name: "target", Type: types.TypeBytes32, Indexed: false},
			abi.Argument{Name: "memo", Type: types.TypeString, Indexed: false},
		})
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

	args, err := FIP20CrossChainMethod.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("failed to unpack input")
	}
	sender, ok0 := args[0].(common.Address)
	receipt, ok1 := args[1].(string)
	amount, ok2 := args[2].(*big.Int)
	fee, ok3 := args[3].(*big.Int)
	target, ok4 := args[4].([32]byte)
	memo, ok5 := args[5].(string)
	if !ok0 || !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return nil, errors.New("unexpected arg type")
	}

	amountCoin := sdk.NewCoin(tokenPair.GetDenom(), sdkmath.NewIntFromBigInt(amount))
	feeCoin := sdk.NewCoin(tokenPair.GetDenom(), sdkmath.NewIntFromBigInt(fee))
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
			if err := c.bankKeeper.SendCoins(ctx, sender.Bytes(), tokenContract.Bytes(), valueCoin); err != nil {
				return nil, fmt.Errorf("send coin: %s", err.Error())
			}
		}
	}

	// transfer token from evm to local chain
	if err := c.convertERC20(ctx, evm, tokenPair, totalCoin, sender); err != nil {
		return nil, err
	}

	fxTarget := fxtypes.ParseFxTarget(fxtypes.Byte32ToString(target))
	if err := c.handlerCrossChain(ctx, sender.Bytes(), receipt, amountCoin, feeCoin, fxTarget, memo, false); err != nil {
		return nil, err
	}

	// add event log
	if err := crossChainLog(evm, contract.Address(), sender, tokenPair.GetERC20Contract(),
		receipt, tokenPair.GetDenom(), memo, amount, fee, target); err != nil {
		return nil, err
	}

	// add fip20CrossChain events
	fip20CrossChainEvents(ctx, sender, tokenPair.GetERC20Contract(), receipt, fxtypes.Byte32ToString(target), tokenPair.GetDenom(), amount, fee)

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
	args, err := CrossChainMethod.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("failed to unpack input")
	}
	token, ok0 := args[0].(common.Address)
	receipt, ok1 := args[1].(string)
	amount, ok2 := args[2].(*big.Int)
	fee, ok3 := args[3].(*big.Int)
	target, ok4 := args[4].([32]byte)
	memo, ok5 := args[5].(string)
	if !ok0 || !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return nil, errors.New("unexpected arg type")
	}

	// call param
	value := contract.Value()
	sender := contract.Caller()

	// cross chain param
	originToken := false
	crossChainDenom := ""

	// cross-chain origin token
	if value.Cmp(big.NewInt(0)) == 1 && token.String() == fxtypes.EmptyEvmAddress {
		totalAmount := big.NewInt(0).Add(amount, fee)
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
		crossChainDenom, err = c.handlerERC20Token(ctx, evm, token, sender, big.NewInt(0).Add(amount, fee))
		if err != nil {
			return nil, err
		}
	}

	fxTarget := fxtypes.ParseFxTarget(fxtypes.Byte32ToString(target))
	amountCoin := sdk.NewCoin(crossChainDenom, sdkmath.NewIntFromBigInt(amount))
	feeCoin := sdk.NewCoin(crossChainDenom, sdkmath.NewIntFromBigInt(fee))

	if err := c.handlerCrossChain(ctx, sender.Bytes(), receipt, amountCoin, feeCoin, fxTarget, memo, originToken); err != nil {
		return nil, err
	}

	// add event log
	if err := crossChainLog(evm, contract.Address(), sender, token,
		receipt, crossChainDenom, memo, amount, fee, target); err != nil {
		return nil, err
	}

	// add cross chain events
	crossChainEvents(ctx, sender, token, receipt, fxtypes.Byte32ToString(target), crossChainDenom, memo, amount, fee)

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
	if tokenPair.GetContractOwner() == erc20types.OWNER_MODULE {
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

	} else if tokenPair.GetContractOwner() == erc20types.OWNER_EXTERNAL {
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

func crossChainLog(evm *vm.EVM, logAddr, sender, token common.Address, recipient, denom, memo string, amount, fee *big.Int, target [32]byte) error {
	eventData, err := CrossChainEvent.Inputs.NonIndexed().Pack(denom, recipient, amount, fee, target, memo)
	if err != nil {
		return err
	}
	topic := []common.Hash{
		CrossChainEvent.ID,
		sender.Hash(),
		token.Hash(),
	}
	evm.StateDB.AddLog(&ethtypes.Log{
		Address:     logAddr,
		Topics:      topic,
		Data:        eventData,
		BlockNumber: evm.Context.BlockNumber.Uint64(),
	})
	return nil
}
