package bank

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/precompiles/types"
	evmtypes "github.com/pundiai/fx-core/v8/x/evm/types"
)

var (
	bankAbi     = contract.MustABIJson(contract.IBankMetaData.ABI)
	bankAddress = common.HexToAddress(contract.BankAddress)
)

type Contract struct {
	methods   []contract.PrecompileMethod
	govKeeper types.GovKeeper
}

func NewPrecompiledContract(
	bankKeeper types.BankKeeper,
	erc20Keeper types.Erc20Keeper,
	govKeeper types.GovKeeper,
) *Contract {
	keeper := NewKeeper(bankKeeper, erc20Keeper)
	return &Contract{
		govKeeper: govKeeper,
		methods: []contract.PrecompileMethod{
			NewTransferFromModuleToAccountMethod(keeper),
			NewTransferFromAccountToModuleMethod(keeper),
		},
	}
}

func (c *Contract) Address() common.Address {
	return bankAddress
}

func (c *Contract) IsStateful() bool {
	return true
}

func (c *Contract) RequiredGas(input []byte) uint64 {
	if len(input) <= 4 {
		return 0
	}
	for _, method := range c.methods {
		if bytes.Equal(method.GetMethodId(), input[:4]) {
			return method.RequiredGas()
		}
	}
	return 0
}

func (c *Contract) Run(evm *vm.EVM, vmContract *vm.Contract, readonly bool) (ret []byte, err error) {
	if len(vmContract.Input) <= 4 {
		return contract.PackRetErrV2(errors.New("invalid input"))
	}

	for _, method := range c.methods {
		if bytes.Equal(method.GetMethodId(), vmContract.Input[:4]) {
			if readonly && !method.IsReadonly() {
				return contract.PackRetErrV2(errors.New("write protection"))
			}

			stateDB := evm.StateDB.(evmtypes.ExtStateDB)
			if err = c.govKeeper.CheckDisabledPrecompiles(stateDB.Context(), c.Address(), method.GetMethodId()); err != nil {
				return contract.PackRetErrV2(err)
			}

			ret, err = method.Run(evm, vmContract)
			if err != nil {
				return contract.PackRetErrV2(err)
			}
			return ret, nil
		}
	}
	return contract.PackRetErrV2(errors.New("unknown method"))
}
