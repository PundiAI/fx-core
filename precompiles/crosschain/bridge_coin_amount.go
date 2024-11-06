package crosschain

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	fxcontract "github.com/functionx/fx-core/v8/contract"
	fxtypes "github.com/functionx/fx-core/v8/types"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
	evmtypes "github.com/functionx/fx-core/v8/x/evm/types"
)

type BridgeCoinAmountMethod struct {
	*Keeper
	BridgeCoinAmountABI
}

func NewBridgeCoinAmountMethod(keeper *Keeper) *BridgeCoinAmountMethod {
	return &BridgeCoinAmountMethod{
		Keeper:              keeper,
		BridgeCoinAmountABI: NewBridgeCoinAmountABI(),
	}
}

func (m *BridgeCoinAmountMethod) IsReadonly() bool {
	return true
}

func (m *BridgeCoinAmountMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *BridgeCoinAmountMethod) RequiredGas() uint64 {
	return 10_000
}

func (m *BridgeCoinAmountMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(evmtypes.ExtStateDB)
	ctx := stateDB.Context()
	crosschainKeeper, found := m.router.GetRoute(ethtypes.ModuleName)
	if !found {
		return nil, errors.New("invalid router")
	}

	totalSupply, err := crosschainKeeper.BridgeCoinSupply(ctx, args.Token.String(), fxtypes.Byte32ToString(args.Target))
	if err != nil {
		return nil, err
	}
	return m.PackOutput(totalSupply.Amount.BigInt())
}

type BridgeCoinAmountABI struct {
	abi.Method
}

func NewBridgeCoinAmountABI() BridgeCoinAmountABI {
	return BridgeCoinAmountABI{
		Method: crosschainABI.Methods["bridgeCoinAmount"],
	}
}

func (m BridgeCoinAmountABI) PackInput(args fxcontract.BridgeCoinAmountArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Token, args.Target)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m BridgeCoinAmountABI) UnpackInput(data []byte) (*fxcontract.BridgeCoinAmountArgs, error) {
	args := new(fxcontract.BridgeCoinAmountArgs)
	if err := evmtypes.ParseMethodArgs(m.Method, args, data[4:]); err != nil {
		return nil, err
	}
	return args, nil
}

func (m BridgeCoinAmountABI) PackOutput(amount *big.Int) ([]byte, error) {
	return m.Method.Outputs.Pack(amount)
}

func (m BridgeCoinAmountABI) UnpackOutput(data []byte) (*big.Int, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return nil, err
	}
	return amount[0].(*big.Int), nil
}
