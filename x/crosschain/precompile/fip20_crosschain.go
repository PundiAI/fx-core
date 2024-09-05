package precompile

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/legacy"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	evmtypes "github.com/functionx/fx-core/v8/x/evm/types"
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

func (m *FIP20CrossChainMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	tokenContract := contract.Caller()

	stateDB := evm.StateDB.(evmtypes.ExtStateDB)
	if err := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		tokenPair, found := m.erc20Keeper.GetTokenPairByAddress(ctx, tokenContract)
		if !found {
			return fmt.Errorf("token pair not found: %s", tokenContract.String())
		}

		args, err := m.UnpackInput(contract.Input)
		if err != nil {
			return err
		}

		amountCoin := sdk.NewCoin(tokenPair.GetDenom(), sdkmath.NewIntFromBigInt(args.Amount))
		feeCoin := sdk.NewCoin(tokenPair.GetDenom(), sdkmath.NewIntFromBigInt(args.Fee))
		totalCoin := sdk.NewCoin(tokenPair.GetDenom(), amountCoin.Amount.Add(feeCoin.Amount))

		// transfer token from evm to local chain
		if err = m.convertERC20(ctx, evm, tokenPair, totalCoin, args.Sender); err != nil {
			return err
		}

		fxTarget := fxtypes.ParseFxTarget(fxtypes.Byte32ToString(args.Target))
		if err = m.handlerCrossChain(ctx, args.Sender.Bytes(), args.Receipt, amountCoin, feeCoin, fxTarget, args.Memo, false); err != nil {
			return err
		}

		data, topic, err := m.NewCrossChainEvent(args.Sender, tokenPair.GetERC20Contract(), tokenPair.GetDenom(), args.Receipt, args.Amount, args.Fee, args.Target, args.Memo)
		if err != nil {
			return err
		}
		EmitEvent(evm, data, topic)

		legacy.Fip20CrossChainEvents(ctx, args.Sender, tokenPair.GetERC20Contract(), args.Receipt,
			fxtypes.Byte32ToString(args.Target), tokenPair.GetDenom(), args.Amount, args.Fee)
		return nil
	}); err != nil {
		return nil, err
	}

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
	return m.Method.Outputs.Pack(success)
}
