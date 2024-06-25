package precompile

import (
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/legacy"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	evmtypes "github.com/functionx/fx-core/v7/x/evm/types"
)

type FIP20CrossChainMethod struct {
	*CrossChainMethod
	abi.Method
}

func NewFIP20CrossChainMethod(keeper *Keeper) *FIP20CrossChainMethod {
	return &FIP20CrossChainMethod{
		CrossChainMethod: NewCrossChainMethod(keeper),
		Method:           crosschaintypes.GetABI().Methods["fip20CrossChain"],
	}
}

func (m *FIP20CrossChainMethod) IsReadonly() bool {
	return false
}

func (m *FIP20CrossChainMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *FIP20CrossChainMethod) RequiredGas() uint64 {
	return 40_000
}

func (m *FIP20CrossChainMethod) Run(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	tokenContract := contract.Caller()
	tokenPair, found := m.erc20Keeper.GetTokenPairByAddress(ctx, tokenContract)
	if !found {
		return nil, fmt.Errorf("token pair not found: %s", tokenContract.String())
	}

	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	amountCoin := sdk.NewCoin(tokenPair.GetDenom(), sdkmath.NewIntFromBigInt(args.Amount))
	feeCoin := sdk.NewCoin(tokenPair.GetDenom(), sdkmath.NewIntFromBigInt(args.Fee))
	totalCoin := sdk.NewCoin(tokenPair.GetDenom(), amountCoin.Amount.Add(feeCoin.Amount))

	// NOTE: if user call evm denom transferCrossChain with msg.value
	// we need transfer msg.value from sender to contract in bank keeper
	if tokenPair.GetDenom() == fxtypes.DefaultDenom {
		balance := m.bankKeeper.GetBalance(ctx, tokenContract.Bytes(), fxtypes.DefaultDenom)
		evmBalance := evm.StateDB.GetBalance(tokenContract)

		cmp := evmBalance.Cmp(balance.Amount.BigInt())
		if cmp == -1 {
			return nil, fmt.Errorf("invalid balance(chain: %s,evm: %s)", balance.Amount.String(), evmBalance.String())
		}
		if cmp == 1 {
			// sender call transferCrossChain with msg.value, the msg.value evm denom should send to contract
			value := big.NewInt(0).Sub(evmBalance, balance.Amount.BigInt())
			valueCoin := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(value)))
			if err := m.bankKeeper.SendCoins(ctx, args.Sender.Bytes(), tokenContract.Bytes(), valueCoin); err != nil {
				return nil, fmt.Errorf("send coin: %s", err.Error())
			}
		}
	}

	// transfer token from evm to local chain
	if err = m.convertERC20(ctx, evm, tokenPair, totalCoin, args.Sender); err != nil {
		return nil, err
	}

	fxTarget := fxtypes.ParseFxTarget(fxtypes.Byte32ToString(args.Target))
	if err = m.handlerCrossChain(ctx, args.Sender.Bytes(), args.Receipt, amountCoin, feeCoin, fxTarget, args.Memo, false); err != nil {
		return nil, err
	}

	data, topic, err := m.NewCrossChainEvent(args.Sender, tokenPair.GetERC20Contract(), tokenPair.GetDenom(), args.Receipt, args.Amount, args.Fee, args.Target, args.Memo)
	if err != nil {
		return nil, err
	}
	EmitEvent(evm, data, topic)

	legacy.Fip20CrossChainEvents(ctx, args.Sender, tokenPair.GetERC20Contract(), args.Receipt,
		fxtypes.Byte32ToString(args.Target), tokenPair.GetDenom(), args.Amount, args.Fee)

	return m.PackOutput(true)
}

func (m *FIP20CrossChainMethod) UnpackInput(data []byte) (*crosschaintypes.FIP20CrossChainArgs, error) {
	args := new(crosschaintypes.FIP20CrossChainArgs)
	if err := evmtypes.ParseMethodArgs(m.Method, args, data[4:]); err != nil {
		return nil, err
	}
	return args, nil
}

func (m *FIP20CrossChainMethod) PackOutput(success bool) ([]byte, error) {
	pack, err := m.Method.Outputs.Pack(success)
	if err != nil {
		return nil, err
	}
	return pack, nil
}
