package crosschain

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/x/evm/types"
)

type GetERC20TokenMethod struct {
	*Keeper
	GetERC20TokenABI
}

func NewGetERC20TokenMethod(keeper *Keeper) *GetERC20TokenMethod {
	return &GetERC20TokenMethod{
		Keeper:           keeper,
		GetERC20TokenABI: NewGetERC20TokenABI(),
	}
}

func (m *GetERC20TokenMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *GetERC20TokenMethod) RequiredGas() uint64 {
	return 1_000
}

func (m *GetERC20TokenMethod) IsReadonly() bool {
	return true
}

func (m *GetERC20TokenMethod) Run(evm *vm.EVM, vmContract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(vmContract.Input)
	if err != nil {
		return nil, err
	}
	stateDB := evm.StateDB.(types.ExtStateDB)
	denom := contract.Byte32ToString(args.Denom)
	erc20Token, err := m.Keeper.erc20Keeper.GetERC20Token(stateDB.Context(), denom)
	if err != nil {
		return nil, err
	}
	return m.PackOutput(common.HexToAddress(erc20Token.GetErc20Address()), erc20Token.Enabled)
}

type GetERC20TokenABI struct {
	abi.Method
}

func NewGetERC20TokenABI() GetERC20TokenABI {
	return GetERC20TokenABI{
		Method: crosschainABI.Methods["getERC20Token"],
	}
}

func (m GetERC20TokenABI) PackInput(args contract.GetERC20TokenArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Denom)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m GetERC20TokenABI) UnpackInput(data []byte) (*contract.GetERC20TokenArgs, error) {
	args := new(contract.GetERC20TokenArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m GetERC20TokenABI) PackOutput(tokenAddr common.Address, enable bool) ([]byte, error) {
	return m.Method.Outputs.Pack(tokenAddr, enable)
}
